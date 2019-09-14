package orderwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/expirationwatch"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidate"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	logger "github.com/sirupsen/logrus"
)

// minCleanupInterval specified the minimum amount of time between orderbook
// cleanups. These cleanups are meant to catch any stale orders that somehow
// were not caught by the event watcher process.
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
	contractAddresses          ethereum.ContractAddresses
	expirationBuffer           time.Duration
	expirationWatcher          *expirationwatch.Watcher
	orderFeed                  event.Feed
	orderScope                 event.SubscriptionScope // Subscription scope tracking current live listeners
	contractAddressToSeenCount map[common.Address]uint
	orderValidator             *ordervalidate.OrderValidator
	wasStartedOnce             bool
	mu                         sync.Mutex
}

// New instantiates a new order watcher
func New(meshDB *meshdb.MeshDB, blockWatcher *blockwatch.Watcher, orderValidator *ordervalidate.OrderValidator, networkID int, expirationBuffer time.Duration) (*Watcher, error) {
	decoder, err := NewDecoder()
	if err != nil {
		return nil, err
	}
	assetDataDecoder := zeroex.NewAssetDataDecoder()
	contractAddresses, err := ethereum.GetContractAddressesForNetworkID(networkID)
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		meshDB:                     meshDB,
		blockWatcher:               blockWatcher,
		expirationBuffer:           expirationBuffer,
		expirationWatcher:          expirationwatch.New(expirationBuffer),
		contractAddressToSeenCount: map[common.Address]uint{},
		orderValidator:             orderValidator,
		eventDecoder:               decoder,
		assetDataDecoder:           assetDataDecoder,
		contractAddresses:          contractAddresses,
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

// Watch sets up the event & expiration watchers as well as the cleanup worker.
// Event watching will require the blockwatch.Watcher to be started first. Watch
// will block until there is a critical error or the given context is canceled.
func (w *Watcher) Watch(ctx context.Context) error {
	w.mu.Lock()
	if w.wasStartedOnce {
		w.mu.Unlock()
		return errors.New("Can only start Watcher once per instance")
	}
	w.wasStartedOnce = true
	w.mu.Unlock()

	// Create a child context so that we can preemptively cancel if there is an
	// error.
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// A waitgroup lets us wait for all goroutines to exit.
	wg := &sync.WaitGroup{}

	// Start three independent goroutines. The expiration watcher, main loop, and
	// cleanup loop. Use three separate channels to communicate errors.
	expirationErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		expirationErrChan <- w.expirationWatcher.Watch(innerCtx, expirationPollingInterval)
	}()
	mainLoopErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainLoopErrChan <- w.mainLoop(innerCtx)
	}()
	cleanupLoopErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		cleanupLoopErrChan <- w.cleanupLoop(innerCtx)
	}()

	// If any error channel returns a non-nil error, we cancel the inner context
	// and return the error. Note that this means we only return the first error
	// that occurs.
	select {
	case err := <-expirationErrChan:
		if err != nil {
			cancel()
			return err
		}
	case err := <-mainLoopErrChan:
		if err != nil {
			cancel()
			return err
		}
	case err := <-cleanupLoopErrChan:
		if err != nil {
			cancel()
			return err
		}
	}

	// Wait for all goroutines to exit. If we reached here it means we are done
	// and there are no errors.
	wg.Wait()
	return nil
}

func (w *Watcher) mainLoop(ctx context.Context) error {
	// Set up the channel used for subscribing to block events.
	blockEvents := make(chan []*blockwatch.Event, 10)
	w.blockSubscription = w.blockWatcher.Subscribe(blockEvents)

	for {
		select {
		case <-ctx.Done():
			w.blockSubscription.Unsubscribe()
			close(blockEvents)
			return nil
		case expiredOrders := <-w.expirationWatcher.ExpiredItems():
			w.handleExpiration(expiredOrders)
		case err := <-w.blockSubscription.Err():
			logger.WithFields(logger.Fields{
				"error": err.Error(),
			}).Error("block subscription error encountered")
		case events := <-blockEvents:
			if err := w.handleBlockEvents(events); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) cleanupLoop(ctx context.Context) error {
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// Wait minCleanupInterval before calling cleanup again. Since
		// we only start sleeping _after_ cleanup completes, we will never
		// have multiple calls to cleanup running in parallel
		time.Sleep(minCleanupInterval - time.Since(start))
		start = time.Now()
		if err := w.cleanup(ctx); err != nil {
			return err
		}
	}
}

func (w *Watcher) handleExpiration(expiredOrders []expirationwatch.ExpiredItem) {
	orderEvents := []*zeroex.OrderEvent{}
	for _, expiredOrder := range expiredOrders {
		order := &meshdb.Order{}
		err := w.meshDB.Orders.FindByID(common.HexToHash(expiredOrder.ID).Bytes(), order)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error":     err.Error(),
				"orderHash": expiredOrder.ID,
			}).Trace("Order expired that was no longer in DB")
			continue
		}
		w.unwatchOrder(w.meshDB.Orders, order, order.FillableTakerAssetAmount)

		orderEvent := &zeroex.OrderEvent{
			OrderHash:                common.HexToHash(expiredOrder.ID),
			SignedOrder:              order.SignedOrder,
			FillableTakerAssetAmount: big.NewInt(0),
			Kind:                     zeroex.EKOrderExpired,
		}
		orderEvents = append(orderEvents, orderEvent)
	}
	w.orderFeed.Send(orderEvents)
}

func (w *Watcher) handleBlockEvents(events []*blockwatch.Event) error {
	ordersColTxn := w.meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersColTxn.Discard()
	}()
	hashToOrderWithTxHashes := map[common.Hash]*OrderWithTxHashes{}
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
					}).Error("unexpected event decoder error encountered")
					return err
				}
			}
			var orders []*meshdb.Order
			switch eventType {
			case "ERC20TransferEvent":
				var transferEvent ERC20TransferEvent
				err = w.eventDecoder.Decode(log, &transferEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				orders, err = w.findOrders(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}

			case "ERC20ApprovalEvent":
				var approvalEvent ERC20ApprovalEvent
				err = w.eventDecoder.Decode(log, &approvalEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				// Ignores approvals set to anyone except the AssetProxy
				if approvalEvent.Spender != w.contractAddresses.ERC20Proxy {
					continue
				}
				orders, err = w.findOrders(approvalEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "ERC721TransferEvent":
				var transferEvent ERC721TransferEvent
				err = w.eventDecoder.Decode(log, &transferEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				orders, err = w.findOrders(transferEvent.From, log.Address, transferEvent.TokenId)
				if err != nil {
					return err
				}

			case "ERC721ApprovalEvent":
				var approvalEvent ERC721ApprovalEvent
				err = w.eventDecoder.Decode(log, &approvalEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				// Ignores approvals set to anyone except the AssetProxy
				if approvalEvent.Approved != w.contractAddresses.ERC721Proxy {
					continue
				}
				orders, err = w.findOrders(approvalEvent.Owner, log.Address, approvalEvent.TokenId)
				if err != nil {
					return err
				}

			case "ERC721ApprovalForAllEvent":
				var approvalForAllEvent ERC721ApprovalForAllEvent
				err = w.eventDecoder.Decode(log, &approvalForAllEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				// Ignores approvals set to anyone except the AssetProxy
				if approvalForAllEvent.Operator != w.contractAddresses.ERC721Proxy {
					continue
				}
				orders, err = w.findOrders(approvalForAllEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "WethWithdrawalEvent":
				var withdrawalEvent WethWithdrawalEvent
				err = w.eventDecoder.Decode(log, &withdrawalEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				orders, err = w.findOrders(withdrawalEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "WethDepositEvent":
				var depositEvent WethDepositEvent
				err = w.eventDecoder.Decode(log, &depositEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				orders, err = w.findOrders(depositEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "ExchangeFillEvent":
				var exchangeFillEvent ExchangeFillEvent
				err = w.eventDecoder.Decode(log, &exchangeFillEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				order := w.findOrder(exchangeFillEvent.OrderHash)
				if order != nil {
					orders = append(orders, order)
				}

			case "ExchangeCancelEvent":
				var exchangeCancelEvent ExchangeCancelEvent
				err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				orders = []*meshdb.Order{}
				order := w.findOrder(exchangeCancelEvent.OrderHash)
				if order != nil {
					orders = append(orders, order)
				}

			case "ExchangeCancelUpToEvent":
				var exchangeCancelUpToEvent ExchangeCancelUpToEvent
				err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				orders, err = w.meshDB.FindOrdersByMakerAddressAndMaxSalt(exchangeCancelUpToEvent.MakerAddress, exchangeCancelUpToEvent.OrderEpoch)
				if err != nil {
					logger.WithFields(logger.Fields{
						"error": err.Error(),
					}).Error("unexpected query error encountered")
					return err
				}
			default:
				logger.WithFields(logger.Fields{
					"eventType": eventType,
					"log":       log,
				}).Error("unknown eventType encountered")
				return err
			}
			for _, order := range orders {
				orderWithTxHashes, ok := hashToOrderWithTxHashes[order.Hash]
				if !ok {
					hashToOrderWithTxHashes[order.Hash] = &OrderWithTxHashes{
						Order: order,
						TxHashes: map[common.Hash]interface{}{
							log.TxHash: struct{}{},
						},
					}
				} else {
					orderWithTxHashes.TxHashes[log.TxHash] = struct{}{}
				}
			}
		}
	}
	return w.generateOrderEventsIfChanged(ordersColTxn, hashToOrderWithTxHashes)
}

func (w *Watcher) cleanup(ctx context.Context) error {
	ordersColTxn := w.meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersColTxn.Discard()
	}()
	lastUpdatedCutOff := time.Now().Add(-lastUpdatedBuffer)
	orders, err := w.meshDB.FindOrdersLastUpdatedBefore(lastUpdatedCutOff)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error":             err.Error(),
			"lastUpdatedCutOff": lastUpdatedCutOff,
		}).Error("Failed to find orders by LastUpdatedBefore")
		return err
	}
	hashToOrderWithTxHashes := map[common.Hash]*OrderWithTxHashes{}
	for _, order := range orders {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		hashToOrderWithTxHashes[order.Hash] = &OrderWithTxHashes{
			Order:    order,
			TxHashes: map[common.Hash]interface{}{},
		}
	}
	return w.generateOrderEventsIfChanged(ordersColTxn, hashToOrderWithTxHashes)
}

// Add adds a 0x order to the DB and watches it for changes in fillability. It
// will no-op (and return nil) if the order has already been added.
func (w *Watcher) Add(orderInfo *ordervalidate.AcceptedOrderInfo) error {
	order := &meshdb.Order{
		Hash:                     orderInfo.OrderHash,
		SignedOrder:              orderInfo.SignedOrder,
		LastUpdated:              time.Now().UTC(),
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
		IsRemoved:                false,
	}
	err := w.meshDB.Orders.Insert(order)
	if err != nil {
		if _, ok := err.(db.AlreadyExistsError); ok {
			// If we're already watching the order, that's fine in this case. Don't
			// return an error.
			return nil
		}
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

	expirationTimestamp := time.Unix(signedOrder.ExpirationTimeSeconds.Int64(), 0)
	w.expirationWatcher.Add(expirationTimestamp, orderHash.Hex())

	return nil
}

// Subscribe allows one to subscribe to the order events emitted by the OrderWatcher.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*zeroex.OrderEvent) event.Subscription {
	return w.orderScope.Track(w.orderFeed.Subscribe(sink))
}

type OrderWithTxHashes struct {
	Order    *meshdb.Order
	TxHashes map[common.Hash]interface{}
}

func (w *Watcher) findOrder(orderHash common.Hash) *meshdb.Order {
	order := meshdb.Order{}
	err := w.meshDB.Orders.FindByID(orderHash.Bytes(), &order)
	if err != nil {
		if _, ok := err.(db.NotFoundError); ok {
			// short-circuit. We expect to receive events from orders we aren't actively tracking
			return nil
		}
		logger.WithFields(logger.Fields{
			"error":     err.Error(),
			"orderHash": orderHash,
		}).Warning("Unexpected error using FindByID for order")
		return nil
	}
	return &order
}

func (w *Watcher) findOrders(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*meshdb.Order, error) {
	orders, err := w.meshDB.FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress, tokenID)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}
	return orders, nil
}

func (w *Watcher) generateOrderEventsIfChanged(ordersColTxn *db.Transaction, hashToOrderWithTxHashes map[common.Hash]*OrderWithTxHashes) error {
	signedOrders := []*zeroex.SignedOrder{}
	for _, orderWithTxHashes := range hashToOrderWithTxHashes {
		order := orderWithTxHashes.Order
		if order.IsRemoved && time.Since(order.LastUpdated) > permanentlyDeleteAfter {
			if err := w.permanentlyDeleteOrder(ordersColTxn, order); err != nil {
				return err
			}
			continue
		}
		signedOrders = append(signedOrders, order.SignedOrder)
	}
	if len(signedOrders) == 0 {
		return nil
	}
	areNewOrders := false
	validationResults := w.orderValidator.BatchValidate(signedOrders, areNewOrders)

	orderEvents := []*zeroex.OrderEvent{}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		orderWithTxHashes := hashToOrderWithTxHashes[acceptedOrderInfo.OrderHash]
		txHashes := make([]common.Hash, len(orderWithTxHashes.TxHashes))
		txHashIndex := 0
		for txHash := range orderWithTxHashes.TxHashes {
			txHashes[txHashIndex] = txHash
			txHashIndex++
		}
		order := orderWithTxHashes.Order
		oldFillableAmount := order.FillableTakerAssetAmount
		newFillableAmount := acceptedOrderInfo.FillableTakerAssetAmount
		oldAmountIsMoreThenNewAmount := oldFillableAmount.Cmp(newFillableAmount) == 1

		expirationTime := time.Unix(order.SignedOrder.ExpirationTimeSeconds.Int64(), 0)
		isExpired := ordervalidate.IsExpired(expirationTime, w.expirationBuffer)
		if !isExpired && oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
			// A previous event caused this order to be removed from DB because it's
			// fillableAmount became 0, but it has now been revived (e.g., block re-org
			// causes order fill txn to get reverted). We need to re-add order and emit an event.
			w.rewatchOrder(ordersColTxn, order, acceptedOrderInfo)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				Kind:                     zeroex.EKOrderAdded,
				TxHashes:                 txHashes,
			}
			orderEvents = append(orderEvents, orderEvent)
		} else if oldFillableAmount.Cmp(newFillableAmount) == 0 {
			// No important state-change happened, simply update lastUpdated timestamp in DB
			w.updateOrderDBEntry(ordersColTxn, order)
		} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && oldAmountIsMoreThenNewAmount {
			// Order was filled, emit  event and update order in DB
			w.updateOrderDBEntry(ordersColTxn, order)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				Kind:                     zeroex.EKOrderFilled,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				TxHashes:                 txHashes,
			}
			orderEvents = append(orderEvents, orderEvent)
		} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && !oldAmountIsMoreThenNewAmount {
			// The order is now fillable for more then it was before. E.g.: A fill txn reverted (block-reorg)
			// Update order in DB and emit event
			w.updateOrderDBEntry(ordersColTxn, order)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				Kind:                     zeroex.EKOrderFillabilityIncreased,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				TxHashes:                 txHashes,
			}
			orderEvents = append(orderEvents, orderEvent)
		}
	}
	for _, rejectedOrderInfo := range validationResults.Rejected {
		switch rejectedOrderInfo.Kind {
		case ordervalidate.MeshError:
			// TODO(fabio): Do we want to handle MeshErrors somehow here?
		case ordervalidate.ZeroExValidation:
			orderWithTxHashes := hashToOrderWithTxHashes[rejectedOrderInfo.OrderHash]
			order := orderWithTxHashes.Order
			oldFillableAmount := order.FillableTakerAssetAmount
			if oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
				// If the oldFillableAmount was already 0, this order is already flagged for removal.
				// Update it's lastUpdated timestamp in DB
				w.updateOrderDBEntry(ordersColTxn, order)
			} else {
				// If oldFillableAmount > 0, it got fullyFilled, cancelled, expired or unfunded, emit event
				w.unwatchOrder(ordersColTxn, order, big.NewInt(0))
				kind, ok := ordervalidate.ConvertRejectOrderCodeToOrderEventKind(rejectedOrderInfo.Status)
				if !ok {
					err := fmt.Errorf("no OrderEventKind corresponding to RejectedOrderStatus: %q", rejectedOrderInfo.Status)
					logger.WithError(err).WithField("rejectedOrderStatus", rejectedOrderInfo.Status).Error("no OrderEventKind corresponding to RejectedOrderStatus")
					return err
				}
				txHashes := make([]common.Hash, len(orderWithTxHashes.TxHashes))
				txHashIndex := 0
				for txHash := range orderWithTxHashes.TxHashes {
					txHashes[txHashIndex] = txHash
					txHashIndex++
				}
				orderEvent := &zeroex.OrderEvent{
					OrderHash:                rejectedOrderInfo.OrderHash,
					SignedOrder:              rejectedOrderInfo.SignedOrder,
					FillableTakerAssetAmount: big.NewInt(0),
					Kind:                     kind,
					TxHashes:                 txHashes,
				}
				orderEvents = append(orderEvents, orderEvent)
			}
		default:
			err := fmt.Errorf("unknown rejectedOrderInfo.Kind: %q", rejectedOrderInfo.Kind)
			logger.WithError(err).Error("encountered unhandled rejectedOrderInfo.Kind value")
			return err
		}
	}
	if err := ordersColTxn.Commit(); err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("Failed to commit orders collection transaction")
	}

	if len(orderEvents) > 0 {
		w.orderFeed.Send(orderEvents)
	}
	return nil
}

type updatableOrdersCol interface {
	Update(model db.Model) error
}

func (w *Watcher) updateOrderDBEntry(u updatableOrdersCol, order *meshdb.Order) {
	order.LastUpdated = time.Now().UTC()
	err := u.Update(order)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}
}

func (w *Watcher) rewatchOrder(u updatableOrdersCol, order *meshdb.Order, orderInfo *ordervalidate.AcceptedOrderInfo) {
	order.IsRemoved = false
	order.LastUpdated = time.Now().UTC()
	order.FillableTakerAssetAmount = orderInfo.FillableTakerAssetAmount
	err := u.Update(order)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}

	// Re-add order to expiration watcher
	expirationTimestamp := time.Unix(order.SignedOrder.ExpirationTimeSeconds.Int64(), 0)
	w.expirationWatcher.Add(expirationTimestamp, order.Hash.Hex())
}

func (w *Watcher) unwatchOrder(u updatableOrdersCol, order *meshdb.Order, newFillableAmount *big.Int) {
	order.IsRemoved = true
	order.LastUpdated = time.Now().UTC()
	order.FillableTakerAssetAmount = newFillableAmount
	err := u.Update(order)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}

	expirationTimestamp := time.Unix(order.SignedOrder.ExpirationTimeSeconds.Int64(), 0)
	w.expirationWatcher.Remove(expirationTimestamp, order.Hash.Hex())
}

func (w *Watcher) permanentlyDeleteOrder(ordersColTxn *db.Transaction, order *meshdb.Order) error {
	err := ordersColTxn.Delete(order.Hash.Bytes())
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
		return nil // Already deleted. Noop.
	}

	// After permanently deleting an order, we also remove it's assetData from the Decoder
	err = w.removeAssetDataAddressFromEventDecoder(order.SignedOrder.MakerAssetData)
	if err != nil {
		// This should never happen since the same error would have happened when adding
		// the assetData to the EventDecoder.
		logger.WithFields(logger.Fields{
			"error":       err.Error(),
			"signedOrder": order.SignedOrder,
		}).Error("Unexpected error when trying to remove an assetData from decoder")
		return err
	}
	return nil
}

// Logs the error and returns true if the error is non-critical.
func (w *Watcher) checkDecodeErr(err error, eventType string) bool {
	if _, ok := err.(UnsupportedEventError); ok {
		logger.WithFields(logger.Fields{
			"eventType":       eventType,
			"topics":          err.(UnsupportedEventError).Topics,
			"contractAddress": err.(UnsupportedEventError).ContractAddress,
		}).Warn("unsupported event found")
		return true
	}
	logger.WithFields(logger.Fields{
		"error": err.Error(),
	}).Error("unexpected event decoder error encountered")
	return false
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
