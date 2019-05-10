package orderwatch

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
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

// expirationPollingInterval specifies the interval in which the order watcher should check for expired
// orders
var expirationPollingInterval = 50 * time.Millisecond

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	blockWatcher            *blockwatch.Watcher
	eventDecoder            *Decoder
	assetDataDecoder        *zeroex.AssetDataDecoder
	blockSubscription       event.Subscription
	expirationWatcher       *ExpirationWatcher
	orderHashToWatchedOrder map[common.Hash]*zeroex.SignedOrder
	watchedOrdersMux        sync.RWMutex
	cleanupCtx              context.Context
	cleanupCancelFunc       context.CancelFunc
	orderValidator          *zeroex.OrderValidator
	isSetup                 bool
	setupMux                sync.RWMutex
}

// New instantiates a new order watcher
func New(blockWatcher *blockwatch.Watcher, ethClient *ethclient.Client, orderValidatorAddress common.Address) (*Watcher, error) {
	decoder, err := NewDecoder()
	if err != nil {
		return nil, err
	}
	assetDataDecoder, err := zeroex.NewAssetDataDecoder()
	if err != nil {
		return nil, err
	}
	orderValidator, err := zeroex.NewOrderValidator(orderValidatorAddress, ethClient)
	if err != nil {
		return nil, err
	}
	cleanupCtx, cleanupCancelFunc := context.WithCancel(context.Background())
	var expirationBuffer int64 = 0
	return &Watcher{
		blockWatcher:            blockWatcher,
		expirationWatcher:       NewExpirationWatcher(expirationBuffer),
		orderHashToWatchedOrder: map[common.Hash]*zeroex.SignedOrder{},
		cleanupCtx:              cleanupCtx,
		cleanupCancelFunc:       cleanupCancelFunc,
		orderValidator:          orderValidator,
		eventDecoder:            decoder,
		assetDataDecoder:        assetDataDecoder,
	}, nil
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

// Watch adds a 0x order to the ones being tracked for order-relevant state changes
func (w *Watcher) Watch(signedOrder *zeroex.SignedOrder, orderHash common.Hash) error {
	w.eventDecoder.AddKnownExchange(signedOrder.ExchangeAddress)

	err := w.addAssetDataAddressToEventDecoder(signedOrder.MakerAssetData)
	if err != nil {
		return err
	}
	err = w.addAssetDataAddressToEventDecoder(signedOrder.TakerAssetData)
	if err != nil {
		return err
	}

	w.expirationWatcher.Add(signedOrder.ExpirationTimeSeconds.Int64(), orderHash)

	w.watchedOrdersMux.Lock()
	w.orderHashToWatchedOrder[orderHash] = signedOrder
	w.watchedOrdersMux.Unlock()

	return nil
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

			// TODO(fabio): Once orders are stored in DB, only fetch orders where lastUpdated field is > X
			signedOrders := []*zeroex.SignedOrder{}
			w.watchedOrdersMux.Lock()
			for _, signedOrder := range w.orderHashToWatchedOrder {
				signedOrders = append(signedOrders, signedOrder)
			}
			w.watchedOrdersMux.Unlock()
			orderInfos := w.orderValidator.BatchValidate(signedOrders)
			for orderHash, orderInfo := range orderInfos {
				// TODO(fabio): Emit all change events
				_ = orderHash
				_ = orderInfo
			}

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
				// TODO(fabio): Emit an event for the expired order
				w.watchedOrdersMux.Lock()
				delete(w.orderHashToWatchedOrder, expiredOrder.OrderHash)
				w.watchedOrdersMux.Unlock()
				panic(fmt.Sprintf("Handling expired orders is not implemented yet: %+v\n", expiredOrder))
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
							// TODO(fabio): Handle this event
						case "ERC20ApprovalEvent":
							var approvalEvent ERC20ApprovalEvent
							err = w.eventDecoder.Decode(log, &approvalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "ERC721TransferEvent":
							var transferEvent ERC721TransferEvent
							err = w.eventDecoder.Decode(log, &transferEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "ERC721ApprovalEvent":
							var approvalEvent ERC721ApprovalEvent
							err = w.eventDecoder.Decode(log, &approvalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "ERC721ApprovalForAllEvent":
							var approvalForAllEvent ERC721ApprovalForAllEvent
							err = w.eventDecoder.Decode(log, &approvalForAllEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "WethWithdrawalEvent":
							var withdrawalEvent WethWithdrawalEvent
							err = w.eventDecoder.Decode(log, &withdrawalEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "WethDepositEvent":
							var depositEvent WethDepositEvent
							err = w.eventDecoder.Decode(log, &depositEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "ExchangeFillEvent":
							var exchangeFillEvent ExchangeFillEvent
							err = w.eventDecoder.Decode(log, &exchangeFillEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "ExchangeCancelEvent":
							var exchangeCancelEvent ExchangeCancelEvent
							err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
						case "ExchangeCancelUpToEvent":
							var exchangeCancelUpToEvent ExchangeCancelUpToEvent
							err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
							if err != nil {
								w.handleDecodeErr(err, eventType)
								continue
							}
							// TODO(fabio): Handle this event
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
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC721(decodedAssetData.Address)
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
		return errors.New(fmt.Sprintf("Unrecognized assetData type name found: %s\n", assetDataName))
	}
	return nil
}
