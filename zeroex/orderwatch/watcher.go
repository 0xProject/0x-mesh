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
	"github.com/0xProject/0x-mesh/db"
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
	orderFeed                  event.Feed
	orderScope                 event.SubscriptionScope // Subscription scope tracking current live listeners
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
	assetDataDecoder := zeroex.NewAssetDataDecoder()
	orderValidator, err := zeroex.NewOrderValidator(ethClient, networkId)
	if err != nil {
		return nil, err
	}
	cleanupCtx, cleanupCancelFunc := context.WithCancel(context.Background())
	var expirationBuffer int64 = 0
	contractNameToAddress := constants.NetworkIDToContractAddresses[networkId]
	w := &Watcher{
		meshDB:                     meshDB,
		blockWatcher:               blockWatcher,
		expirationWatcher:          NewExpirationWatcher(expirationBuffer),
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
		err := w.setupInMemoryOrderState(order.SignedOrder)
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
func (w *Watcher) Watch(orderInfo *zeroex.AcceptedOrderInfo) error {
	order := &meshdb.Order{
		Hash:                     orderInfo.OrderHash,
		SignedOrder:              orderInfo.SignedOrder,
		LastUpdated:              time.Now().UTC(),
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
		IsRemoved:                false,
	}
	err := w.meshDB.Orders.Insert(order)
	if err != nil {
		return err
	}

	err = w.setupInMemoryOrderState(orderInfo.SignedOrder)
	if err != nil {
		return err
	}

	orderEvent := &zeroex.OrderEvent{
		OrderHash:                orderInfo.OrderHash,
		SignedOrder:              orderInfo.SignedOrder,
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
		Kind:                     zeroex.EKOrderAdded,
	}
	w.orderFeed.Send([]*zeroex.OrderEvent{orderEvent})

	return nil
}

func (w *Watcher) setupInMemoryOrderState(signedOrder *zeroex.SignedOrder) error {
	orderHash, err := signedOrder.ComputeOrderHash()
	if err != nil {
		return err
	}

	w.eventDecoder.AddKnownExchange(signedOrder.ExchangeAddress)

	err = w.addAssetDataAddressToEventDecoder(signedOrder.MakerAssetData)
	if err != nil {
		return err
	}

	w.expirationWatcher.Add(signedOrder.ExpirationTimeSeconds.Int64(), orderHash)

	return nil
}

// Subscribe allows one to subscribe to the order events emitted by the OrderWatcher.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*zeroex.OrderEvent) event.Subscription {
	return w.orderScope.Track(w.orderFeed.Subscribe(sink))
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

			// We do not re-validate orders that have been updated within the last `lastUpdatedBuffer` time
			lastUpdatedCutOff := start.Add(-lastUpdatedBuffer)
			orders, err := w.meshDB.FindOrdersLastUpdatedBefore(lastUpdatedCutOff)
			if err != nil {
				logger.WithFields(logger.Fields{
					"error":             err.Error(),
					"lastUpdatedCutOff": lastUpdatedCutOff,
				}).Panic("Failed to find orders by LastUpdatedBefore")
			}

			w.generateOrderEventsIfChanged(orders, common.Hash{})

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
					OrderStatus:              zeroex.OSExpired,
				}
				w.unwatchOrder(order)

				orderEvent := &zeroex.OrderEvent{
					OrderHash:                orderInfo.OrderHash,
					SignedOrder:              orderInfo.SignedOrder,
					FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
					Kind:                     zeroex.EKOrderExpired,
				}
				w.orderFeed.Send([]*zeroex.OrderEvent{orderEvent})
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
									"error": err.Error(),
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
							w.findOrdersAndGenerateOrderEvents(transferEvent.From, log.Address, nil, log.TxHash)

						case "ERC20ApprovalEvent":
							var approvalEvent ERC20ApprovalEvent
							err = w.eventDecoder.Decode(log, &approvalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// Ignores approvals set to anyone except the AssetProxy
							if approvalEvent.Spender != w.contractNameToAddress.ERC20Proxy {
								continue
							}
							w.findOrdersAndGenerateOrderEvents(approvalEvent.Owner, log.Address, nil, log.TxHash)

						case "ERC721TransferEvent":
							var transferEvent ERC721TransferEvent
							err = w.eventDecoder.Decode(log, &transferEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(transferEvent.From, log.Address, transferEvent.TokenId, log.TxHash)

						case "ERC721ApprovalEvent":
							var approvalEvent ERC721ApprovalEvent
							err = w.eventDecoder.Decode(log, &approvalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// Ignores approvals set to anyone except the AssetProxy
							if approvalEvent.Approved != w.contractNameToAddress.ERC721Proxy {
								continue
							}
							w.findOrdersAndGenerateOrderEvents(approvalEvent.Owner, log.Address, approvalEvent.TokenId, log.TxHash)

						case "ERC721ApprovalForAllEvent":
							var approvalForAllEvent ERC721ApprovalForAllEvent
							err = w.eventDecoder.Decode(log, &approvalForAllEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// Ignores approvals set to anyone except the AssetProxy
							if approvalForAllEvent.Operator != w.contractNameToAddress.ERC721Proxy {
								continue
							}
							w.findOrdersAndGenerateOrderEvents(approvalForAllEvent.Owner, log.Address, nil, log.TxHash)

						case "WethWithdrawalEvent":
							var withdrawalEvent WethWithdrawalEvent
							err = w.eventDecoder.Decode(log, &withdrawalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(withdrawalEvent.Owner, log.Address, nil, log.TxHash)

						case "WethDepositEvent":
							var depositEvent WethDepositEvent
							err = w.eventDecoder.Decode(log, &depositEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrdersAndGenerateOrderEvents(depositEvent.Owner, log.Address, nil, log.TxHash)

						case "ExchangeFillEvent":
							var exchangeFillEvent ExchangeFillEvent
							err = w.eventDecoder.Decode(log, &exchangeFillEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrderAndGenerateOrderEvents(exchangeFillEvent.OrderHash, log.TxHash)

						case "ExchangeCancelEvent":
							var exchangeCancelEvent ExchangeCancelEvent
							err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							w.findOrderAndGenerateOrderEvents(exchangeCancelEvent.OrderHash, log.TxHash)

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
									"error": err.Error(),
								}).Panic("unexpected query error encountered")
							}
							w.generateOrderEventsIfChanged(orders, log.TxHash)

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

func (w *Watcher) findOrderAndGenerateOrderEvents(orderHash common.Hash, txHash common.Hash) {
	order := meshdb.Order{}
	err := w.meshDB.Orders.FindByID(orderHash.Bytes(), &order)
	if err != nil {
		if _, ok := err.(db.NotFoundError); ok {
			return // We will receive events from orders we aren't actively tracking
		}
		logger.WithFields(logger.Fields{
			"error":     err.Error(),
			"orderHash": orderHash,
		}).Warning("Unexpected error using FindByID for order")
	}
	w.generateOrderEventsIfChanged([]*meshdb.Order{&order}, txHash)
}

func (w *Watcher) findOrdersAndGenerateOrderEvents(makerAddress, tokenAddress common.Address, tokenID *big.Int, txHash common.Hash) {
	orders, err := w.meshDB.FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress, tokenID)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Panic("unexpected query error encountered")
	}
	w.generateOrderEventsIfChanged(orders, txHash)
}

func (w *Watcher) generateOrderEventsIfChanged(orders []*meshdb.Order, txHash common.Hash) {
	signedOrders := []*zeroex.SignedOrder{}
	orderHashToOrder := map[common.Hash]*meshdb.Order{}
	for _, order := range orders {
		if order.IsRemoved && time.Since(order.LastUpdated) > permanentlyDeleteAfter {
			w.permanentlyDeleteOrder(order)
			continue
		}
		signedOrders = append(signedOrders, order.SignedOrder)
		orderHashToOrder[order.Hash] = order
	}
	validationResults := w.orderValidator.BatchValidate(signedOrders)

	orderEvents := []*zeroex.OrderEvent{}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		order := orderHashToOrder[acceptedOrderInfo.OrderHash]
		oldFillableAmount := order.FillableTakerAssetAmount
		newFillableAmount := acceptedOrderInfo.FillableTakerAssetAmount
		oldAmountIsMoreThenNewAmount := oldFillableAmount.Cmp(newFillableAmount) == 1
		if oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
			// A previous event caused this order to be removed from DB, but it has now
			// been revived (e.g., block re-org causes order fill txn to get reverted)
			// Need to re-add order and emit an event
			w.rewatchOrder(order, acceptedOrderInfo)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				Kind:                     zeroex.EKOrderAdded,
				TxHash:                   txHash,
			}
			orderEvents = append(orderEvents, orderEvent)
		} else if oldFillableAmount.Cmp(newFillableAmount) == 0 {
			// No important state-change happened, ignore
			// Noop
		} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && oldAmountIsMoreThenNewAmount {
			// Order was filled, emit  event
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				Kind:                     zeroex.EKOrderFilled,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				TxHash:                   txHash,
			}
			orderEvents = append(orderEvents, orderEvent)
		} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && !oldAmountIsMoreThenNewAmount {
			// The order is now fillable for more then it was before. E.g.:
			// 1. A fill txn reverted (block-reorg)
			// 2. Traders added missing balance/allowance increasing the order's fillability
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				Kind:                     zeroex.EKOrderFillabilityIncreased,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				TxHash:                   txHash,
			}
			orderEvents = append(orderEvents, orderEvent)
		}
	}
	for _, rejectedOrderInfo := range validationResults.Rejected {
		switch rejectedOrderInfo.Kind {
		case zeroex.MeshError:
			// TODO(fabio): Do we want to handle MeshErrors somehow here?
		case zeroex.ZeroExValidation:
			order := orderHashToOrder[rejectedOrderInfo.OrderHash]
			oldFillableAmount := order.FillableTakerAssetAmount
			if oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
				// If the oldFillableAmount was already 0, this order is already flagged for removal
				// Noop
			} else {
				// If oldFillableAmount > 0, it got fullyFilled, cancelled, expired or unfunded, emit event
				w.unwatchOrder(order)
				kind, ok := zeroex.ConvertRejectOrderCodeToOrderEventKind(rejectedOrderInfo.Status)
				if !ok {
					logger.WithField("rejectedOrderStatus", rejectedOrderInfo.Status).Panic("No OrderEventKind corresponding to RejectedOrderStatus")
				}
				orderEvent := &zeroex.OrderEvent{
					OrderHash:                rejectedOrderInfo.OrderHash,
					SignedOrder:              rejectedOrderInfo.SignedOrder,
					FillableTakerAssetAmount: big.NewInt(0),
					Kind:                     kind,
					TxHash:                   txHash,
				}
				orderEvents = append(orderEvents, orderEvent)
			}
		default:
			logger.WithField("kind", rejectedOrderInfo.Kind).Panic("Encountered unhandled rejectedOrderInfo.Kind value")
		}
	}
	w.orderFeed.Send(orderEvents)
}

func (w *Watcher) rewatchOrder(order *meshdb.Order, orderInfo *zeroex.AcceptedOrderInfo) {
	order.IsRemoved = false
	order.LastUpdated = time.Now().UTC()
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
	order.LastUpdated = time.Now().UTC()
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
		}).Warn("Attempted to delete order that no longer exists")
		// TODO(fabio): With the current way the OrderWatcher is written, it is possible for multiple
		// events to trigger logic that updates the orders in the DB simultaneously. This is mostly
		// benign but is a waste of computation, and causes processes to try and delete orders the
		// have already been deleted. In order to fix this, we need to re-write the event handling logic
		// to queue the processing of events so that they happen sequentially rather then in parallel.
		return // Already deleted. Noop.
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
			"error": err.Error(),
		}).Panic("unexpected event decoder error encountered")
	}
}

// addAssetDataAddressToEventDecoder decodes the supplied AssetData and figures out which
// tokens (address & token standard) it contains. It then registers these token addresses
// with the contract events decoder so that knows how to properly decode events from that
// contract address. This is necessary because different token standards share identical
// event signatures but use different parameter names (see: decoder.go for more context).
// In order to unregister token addresses when the last order involving it is deleted, we
// also keep track of the number of tokens seen referencing a particular token address.
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

// Whenever we delete an order from the DB, we must also decrement the count of orders
// involving a specific token address. We therefore call this method which decrements the
// count, and if it reaches 0 for a given token, it removes the token address from the
// contract event decoder.
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
