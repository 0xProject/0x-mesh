package orderwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	logger "github.com/sirupsen/logrus"
)

// minCleanupInterval specified the minimum amount of time between orderbook cleanup intervals. These
// cleanups are meant to catch any stale orders that somehow were not caught by the event watcher
// process.
var minCleanupInterval = 1 * time.Hour

// lastUpdatedBuffer specifies how long it must have been since an order was last updated in order to
// be re-validated by the cleanup worker
var lastUpdatedBuffer = 30 * time.Minute

// permanentlyDeleteAfter specifies how long after an order is marked as IsRemoved and not updated that
// it should be considered for permanent deletion. Blocks get mined on avg. every 12 sec, so 4 minutes
// corresponds to a block depth of ~20.
var permanentlyDeleteAfter = 4 * time.Minute

// expirationPollingInterval specifies the interval in which the order watcher should check for expired
// orders
var expirationPollingInterval = 50 * time.Millisecond

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	meshDB                     *meshdb.MeshDB
	blockWatcher               *blockwatch.Watcher
	eventDecoder               *Decoder
	assetDataDecoder           *zeroex.AssetDataDecoder
	blockSubscription          event.Subscription
	contractNameToAddress      constants.ContractNameToAddress
	expirationWatcher          *ExpirationWatcher
	cachedOrderEvents          []*zeroex.OrderInfo
	cachedOrderEventsMux       sync.RWMutex
	cleanupCtx                 context.Context
	cleanupCancelFunc          context.CancelFunc
	contractAddressToSeenCount map[common.Address]uint
	orderValidator             *zeroex.OrderValidator
	isSetup                    bool
	setupMux                   sync.RWMutex
}

// New instantiates a new order watcher
func New(meshDB *meshdb.MeshDB, blockWatcher *blockwatch.Watcher, ethClient *ethclient.Client, networkId int) (*Watcher, error) {
	decoder, err := NewDecoder()
	if err != nil {
		return nil, err
	}
	assetDataDecoder, err := zeroex.NewAssetDataDecoder()
	if err != nil {
		return nil, err
	}
	contractNameToAddress := constants.NetworkIDToContractAddresses[networkId]
	orderValidator, err := zeroex.NewOrderValidator(contractNameToAddress["OrderValidator"], ethClient)
	if err != nil {
		return nil, err
	}
	cleanupCtx, cleanupCancelFunc := context.WithCancel(context.Background())
	var expirationBuffer int64 = 0
	w := &Watcher{
		meshDB:                     meshDB,
		blockWatcher:               blockWatcher,
		expirationWatcher:          NewExpirationWatcher(expirationBuffer),
		cachedOrderEvents:          []*zeroex.OrderInfo{},
		cleanupCtx:                 cleanupCtx,
		cleanupCancelFunc:          cleanupCancelFunc,
		contractAddressToSeenCount: map[common.Address]uint{},
		orderValidator:             orderValidator,
		eventDecoder:               decoder,
		assetDataDecoder:           assetDataDecoder,
		contractNameToAddress:      contractNameToAddress,
	}

	// Pre-populate the OrderWatcher with all orders already stored in the DB
	orders := []*meshdb.Order{}
	err = w.meshDB.Orders.FindAll(&orders)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		err := w.setupInMemoryOrderState(order.SignedOrder, order.Hash)
		if err != nil {
			return nil, err
		}
	}

	return w, nil
}

// Start sets up the event & expiration watchers as well as the cleanup worker. Event
// watching will require the blockwatch.Watcher to be started however.
func (w *Watcher) Start() error {
	w.setupMux.Lock()
	defer w.setupMux.Unlock()
	if w.isSetup {
		return errors.New("Setup can only be called once")
	}

	w.setupEventWatcher()

	if err := w.setupExpirationWatcher(); err != nil {
		return err
	}

	w.startCleanupWorker()

	w.isSetup = true
	return nil
}

// Stop closes the block subscription, stops the event, expiration watcher and the cleanup worker.
func (w *Watcher) Stop() error {
	w.setupMux.Lock()
	if !w.isSetup {
		w.setupMux.Unlock()
		return errors.New("Cannot teardown before calling Setup()")
	}
	w.setupMux.Unlock()

	// Stop event subscription
	w.blockSubscription.Unsubscribe()

	// Stop expiration watcher
	w.expirationWatcher.Stop()

	// Stop cleanup worker
	w.stopCleanupWorker()
	return nil
}

// Watch adds a 0x order to the DB and watches it for changes in fillability.
func (w *Watcher) Watch(signedOrder *zeroex.SignedOrder, orderInfo *zeroex.OrderInfo) error {
	if orderInfo.FillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
		logger.WithFields(logger.Fields{
			"signedOrder": signedOrder,
			"orderInfo":   orderInfo,
		}).Panic("Attempted to add unfillable order to OrderWatcher")
	}
	order := meshdb.Order{
		Hash:                     orderInfo.OrderHash,
		SignedOrder:              signedOrder,
		LastUpdated:              time.Now().Truncate(0),
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
		IsRemoved:                false,
	}
	err := w.meshDB.Orders.Insert(order)
	if err != nil {
		return err
	}

	err := w.setupInMemoryOrderState(signedOrder, orderInfo.OrderHash)
	if err != nil {
		return err
	}

	return nil
}

func (w *Watcher) setupInMemoryOrderState(signedOrder *zeroex.SignedOrder, orderHash common.Hash) error {
	w.eventDecoder.AddKnownExchange(signedOrder.ExchangeAddress)

	err := w.addAssetDataAddressToEventDecoder(signedOrder.MakerAssetData)
	if err != nil {
		return err
	}

	w.expirationWatcher.Add(signedOrder.ExpirationTimeSeconds.Int64(), orderHash)

	return nil
}

// GetEvents fetches the latest order events emitted by the OrderWatcher.
func (w *Watcher) GetEvents() []*zeroex.OrderInfo {
	w.cachedOrderEventsMux.Lock()
	cachedEvents := w.cachedOrderEvents
	w.cachedOrderEvents = []*zeroex.OrderInfo{}
	w.cachedOrderEventsMux.Lock()
	return cachedEvents
}

func (w *Watcher) startCleanupWorker() {
	go func() {
		for {
			select {
			case <-w.cleanupCtx.Done():
				return
			default:
			}

			start := time.Now()

			// We do not re-validate orders that have been updated within the last 30mins
			lastUpdatedCutOff := start.Add(-lastUpdatedBuffer)
			orders, err := w.meshDB.FindOrdersLastUpdatedBefore(lastUpdatedCutOff)
			if err != nil {
				logger.WithFields(logger.Fields{
					"error":             err.Error(),
					"lastUpdatedCutOff": lastUpdatedCutOff,
				}).Panic("Failed to find orders by LastUpdatedBefore")
			}

			w.generateOrderEventsIfChanged(orders)

			// Wait MinCleanupInterval before calling ValidateOrders again. Since
			// we only start sleeping _after_ ValidateOrders completes, we will never
			// have multiple calls to ValidateOrders running in parallel
			time.Sleep(minCleanupInterval - time.Since(start))
		}
	}()
}

func (w *Watcher) stopCleanupWorker() {
	w.cleanupCancelFunc()
}

func (w *Watcher) setupExpirationWatcher() error {
	go func() {
		expiredOrders := w.expirationWatcher.Receive()
		for expiredOrders := range expiredOrders {
			for _, expiredOrder := range expiredOrders {
				order := &meshdb.Order{}
				err := w.meshDB.Orders.FindByID(expiredOrder.OrderHash.Bytes(), order)
				if err != nil {
					logger.WithFields(logger.Fields{
						"error":     err.Error(),
						"orderHash": expiredOrder.OrderHash,
					}).Warning("Order expired that was no longer in DB")
					continue
				}
				orderInfo := &zeroex.OrderInfo{
					OrderHash:                expiredOrder.OrderHash,
					SignedOrder:              order.SignedOrder,
					FillableTakerAssetAmount: big.NewInt(0),
					OrderStatus:              zeroex.Cancelled,
				}
				w.unwatchOrder(order)
				w.cachedOrderEventsMux.Lock()
				w.cachedOrderEvents = append(w.cachedOrderEvents, orderInfo)
				w.cachedOrderEventsMux.Unlock()
			}
		}
	}()

	return w.expirationWatcher.Start(expirationPollingInterval)
}

func (w *Watcher) setupEventWatcher() {
	blockEvents := make(chan []*blockwatch.Event, 10)
	w.blockSubscription = w.blockWatcher.Subscribe(blockEvents)

	go func() {
		for {
			select {
			case err, isOpen := <-w.blockSubscription.Err():
				close(blockEvents)
				if !isOpen {
					// event.Subscription closes the Error channel on unsubscribe.
					// We therefore cleanup this goroutine on channel closure.
					return
				}
				logger.WithFields(logger.Fields{
					"error": err.Error(),
				}).Error("subscription error encountered")
				return

			case events := <-blockEvents:
				for _, event := range events {
					for _, log := range event.BlockHeader.Logs {
						eventType, err := w.eventDecoder.FindEventType(log)
						if err != nil {
							switch err.(type) {
							case UntrackedTokenError:
								continue
							case UnsupportedEventError:
								logger.WithFields(logger.Fields{
									"topics":          err.(UnsupportedEventError).Topics,
									"contractAddress": err.(UnsupportedEventError).ContractAddress,
								}).Info("unsupported event found while trying to find its event type")
								continue
							default:
								logger.WithFields(logger.Fields{
									"err": err.Error(),
								}).Panic("unexpected event decoder error encountered")
							}
						}
						switch eventType {
						case "ERC20TransferEvent":
							var transferEvent ERC20TransferEvent
							err = w.eventDecoder.Decode(log, &transferEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(transferEvent.From, log.Address, nil)

						case "ERC20ApprovalEvent":
							var approvalEvent ERC20ApprovalEvent
							err = w.eventDecoder.Decode(log, &approvalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// Ignores approvals set to anyone except the AssetProxy
							if approvalEvent.Spender != w.contractNameToAddress["ERC20Proxy"] {
								continue
							}
							w.findOrdersAndGenerateOrderEvents(approvalEvent.Owner, log.Address, nil)

						case "ERC721TransferEvent":
							var transferEvent ERC721TransferEvent
							err = w.eventDecoder.Decode(log, &transferEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(transferEvent.From, log.Address, transferEvent.TokenId)

						case "ERC721ApprovalEvent":
							var approvalEvent ERC721ApprovalEvent
							err = w.eventDecoder.Decode(log, &approvalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// Ignores approvals set to anyone except the AssetProxy
							if approvalEvent.Approved != w.contractNameToAddress["ERC721Proxy"] {
								continue
							}
							w.findOrdersAndGenerateOrderEvents(approvalEvent.Owner, log.Address, approvalEvent.TokenId)

						case "ERC721ApprovalForAllEvent":
							var approvalForAllEvent ERC721ApprovalForAllEvent
							err = w.eventDecoder.Decode(log, &approvalForAllEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// Ignores approvals set to anyone except the AssetProxy
							if approvalForAllEvent.Operator != w.contractNameToAddress["ERC721Proxy"] {
								continue
							}
							w.findOrdersAndGenerateOrderEvents(approvalForAllEvent.Owner, log.Address, nil)

						case "WethWithdrawalEvent":
							var withdrawalEvent WethWithdrawalEvent
							err = w.eventDecoder.Decode(log, &withdrawalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(withdrawalEvent.Owner, log.Address, nil)

						case "WethDepositEvent":
							var depositEvent WethDepositEvent
							err = w.eventDecoder.Decode(log, &depositEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(depositEvent.Owner, log.Address, nil)

						case "ExchangeFillEvent":
							var exchangeFillEvent ExchangeFillEvent
							err = w.eventDecoder.Decode(log, &exchangeFillEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrderAndGenerateOrderEvents(exchangeFillEvent.OrderHash)

						case "ExchangeCancelEvent":
							var exchangeCancelEvent ExchangeCancelEvent
							err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrderAndGenerateOrderEvents(exchangeCancelEvent.OrderHash)

						case "ExchangeCancelUpToEvent":
							var exchangeCancelUpToEvent ExchangeCancelUpToEvent
							err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							orders, err := w.meshDB.FindOrdersByMakerAddressAndMaxSalt(exchangeCancelUpToEvent.MakerAddress, exchangeCancelUpToEvent.OrderEpoch)
							if err != nil {
								logger.WithFields(logger.Fields{
									"err": err.Error(),
								}).Panic("unexpected query error encountered")
							}
							w.generateOrderEventsIfChanged(orders)

						default:
							logger.WithFields(logger.Fields{
								"eventType": eventType,
								"log":       log,
							}).Panic("unknown eventType encountered")
						}
					}
				}
			}
		}

	}()
}

func (w *Watcher) findOrderAndGenerateOrderEvents(orderHash common.Hash) {
	order := meshdb.Order{}
	err := w.meshDB.Orders.FindByID(orderHash.Bytes(), &order)
	if err != nil {
		if err.Error() == "model not found" {
			return // We will receive events from orders we aren't actively tracking
		}
		logger.WithFields(logger.Fields{
			"error":     err.Error(),
			"orderHash": orderHash,
		}).Warning("Unexpected error using FindByID for order")
	}
	w.generateOrderEventsIfChanged([]*meshdb.Order{&order})
}

func (w *Watcher) findOrdersAndGenerateOrderEvents(makerAddress, tokenAddress common.Address, tokenID *big.Int) {
	orders, err := w.meshDB.FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress, tokenID)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Panic("unexpected query error encountered")
	}
	w.generateOrderEventsIfChanged(orders)
}

func (w *Watcher) generateOrderEventsIfChanged(orders []*meshdb.Order) {
	signedOrders := []*zeroex.SignedOrder{}
	for _, order := range orders {
		if order.IsRemoved && time.Now().Sub(order.LastUpdated) > permanentlyDeleteAfter {
			w.permanentlyDeleteOrder(order)
			continue
		}
		signedOrders = append(signedOrders, order.SignedOrder)
	}
	hashToOrderInfo := w.orderValidator.BatchValidate(signedOrders)
	for _, order := range orders {
		orderInfo, ok := hashToOrderInfo[order.Hash]
		if !ok {
			continue // Skip orders where OrderInfo was not returned
		}
		if order.FillableTakerAssetAmount != orderInfo.FillableTakerAssetAmount {
			isOrderUnfillable := orderInfo.FillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0
			if isOrderUnfillable {
				w.unwatchOrder(order)
			} else {
				w.rewatchOrder(order, orderInfo)
			}
			w.cachedOrderEventsMux.Lock()
			w.cachedOrderEvents = append(w.cachedOrderEvents, orderInfo)
			w.cachedOrderEventsMux.Unlock()
		}
	}
}

func (w *Watcher) rewatchOrder(order *meshdb.Order, orderInfo *zeroex.OrderInfo) {
	order.IsRemoved = false
	order.LastUpdated = time.Now().Truncate(0)
	order.FillableTakerAssetAmount = orderInfo.FillableTakerAssetAmount
	err := w.meshDB.Orders.Update(order)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}

	// Re-add order to expiration watcher
	w.expirationWatcher.Add(order.SignedOrder.ExpirationTimeSeconds.Int64(), order.Hash)
}

func (w *Watcher) unwatchOrder(order *meshdb.Order) {
	order.IsRemoved = true
	order.LastUpdated = time.Now().Truncate(0)
	order.FillableTakerAssetAmount = big.NewInt(0)
	err := w.meshDB.Orders.Update(order)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}

	w.expirationWatcher.Remove(order.SignedOrder.ExpirationTimeSeconds.Int64(), order.Hash)
}

func (w *Watcher) permanentlyDeleteOrder(order *meshdb.Order) {
	err := w.meshDB.Orders.Delete(order.Hash.Bytes())
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Panic("Unexpected error while attempting to delete an order")
	}

	// After permanently deleting an order, we also remove it's assetData from the Decoder
	err = w.removeAssetDataAddressFromEventDecoder(order.SignedOrder.MakerAssetData)
	if err != nil {
		// This should never happen since the same error would have happened when adding
		// the assetData to the EventDecoder.
		logger.WithFields(logger.Fields{
			"error":       err.Error(),
			"signedOrder": order.SignedOrder,
		}).Panic("Unexpected error when trying to remove an assetData from decoder")
	}
}

func (w *Watcher) handleDecodeErr(err error, eventType string) {
	switch err.(type) {
	case UnsupportedEventError:
		logger.WithFields(logger.Fields{
			"eventType":       eventType,
			"topics":          err.(UnsupportedEventError).Topics,
			"contractAddress": err.(UnsupportedEventError).ContractAddress,
		}).Warn("unsupported event found")

	default:
		logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Panic("unexpected event decoder error encountered")
	}
}

func (w *Watcher) addAssetDataAddressToEventDecoder(assetData []byte) error {
	assetDataName, err := w.assetDataDecoder.GetName(assetData)
	if err != nil {
		return err
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC20(decodedAssetData.Address)
		w.contractAddressToSeenCount[decodedAssetData.Address] = w.contractAddressToSeenCount[decodedAssetData.Address] + 1
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC721(decodedAssetData.Address)
		w.contractAddressToSeenCount[decodedAssetData.Address] = w.contractAddressToSeenCount[decodedAssetData.Address] + 1
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		for _, assetData := range decodedAssetData.NestedAssetData {
			if err := w.addAssetDataAddressToEventDecoder(assetData); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return nil
}

func (w *Watcher) removeAssetDataAddressFromEventDecoder(assetData []byte) error {
	assetDataName, err := w.assetDataDecoder.GetName(assetData)
	if err != nil {
		return err
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.contractAddressToSeenCount[decodedAssetData.Address] = w.contractAddressToSeenCount[decodedAssetData.Address] - 1
		if w.contractAddressToSeenCount[decodedAssetData.Address] == 0 {
			w.eventDecoder.RemoveKnownERC20(decodedAssetData.Address)
		}
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.contractAddressToSeenCount[decodedAssetData.Address] = w.contractAddressToSeenCount[decodedAssetData.Address] - 1
		if w.contractAddressToSeenCount[decodedAssetData.Address] == 0 {
			w.eventDecoder.RemoveKnownERC721(decodedAssetData.Address)
		}
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		for _, assetData := range decodedAssetData.NestedAssetData {
			if err := w.removeAssetDataAddressFromEventDecoder(assetData); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return nil
}
