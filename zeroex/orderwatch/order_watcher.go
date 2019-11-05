package orderwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/expirationwatch"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/slowcounter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	logger "github.com/sirupsen/logrus"
)

const (
	// minCleanupInterval specified the minimum amount of time between orderbook
	// cleanups. These cleanups are meant to catch any stale orders that somehow
	// were not caught by the event watcher process.
	minCleanupInterval = 1 * time.Hour

	// minRemovedCheckInterval specifies the minimum amount of time between checks
	// on whether to remove orders flaggged for removal from the DB
	minRemovedCheckInterval = 5 * time.Minute

	// lastUpdatedBuffer specifies how long it must have been since an order was last updated in order to
	// be re-validated by the cleanup worker
	lastUpdatedBuffer = 30 * time.Minute

	// permanentlyDeleteAfter specifies how long after an order is marked as IsRemoved and not updated that
	// it should be considered for permanent deletion. Blocks get mined on avg. every 12 sec, so 5 minutes
	// corresponds to a block depth of ~25.
	permanentlyDeleteAfter = 5 * time.Minute

	// expirationPollingInterval specifies the interval in which the order watcher should check for expired
	// orders
	expirationPollingInterval = 50 * time.Millisecond

	// maxOrdersTrimRatio affects how many orders are trimmed whenever we reach the
	// maximum number of orders. When order storage is full, Watcher will remove
	// orders until the total number of remaining orders is equal to
	// maxOrdersTrimRatio * maxOrders.
	maxOrdersTrimRatio = 0.9

	// defaultMaxOrders is the default max number of orders in storage.
	defaultMaxOrders = 100000

	// maxExpirationTimeCheckInterval is how often to check whether we can
	// increase the max expiration time.
	maxExpirationTimeCheckInterval = 30 * time.Second

	// configuration options for the SlowCounter used for increasing max
	// expiration time. Effectively, we will increase every 5 minutes as long as
	// there is enough space in the database for orders. The first increase will
	// be 5 seconds and the amount doubles from there (second increase will be 10
	// seconds, then 20 seconds, then 40, etc.)
	slowCounterOffset   = 5 // seconds
	slowCounterRate     = 2.0
	slowCounterInterval = 5 * time.Minute
)

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	meshDB                     *meshdb.MeshDB
	blockWatcher               *blockwatch.Watcher
	eventDecoder               *decoder.Decoder
	assetDataDecoder           *zeroex.AssetDataDecoder
	blockSubscription          event.Subscription
	contractAddresses          ethereum.ContractAddresses
	expirationWatcher          *expirationwatch.Watcher
	orderFeed                  event.Feed
	orderScope                 event.SubscriptionScope // Subscription scope tracking current live listeners
	contractAddressToSeenCount map[common.Address]uint
	orderValidator             *ordervalidator.OrderValidator
	wasStartedOnce             bool
	mu                         sync.Mutex
	maxExpirationTime          *big.Int
	maxExpirationCounter       *slowcounter.SlowCounter
	maxOrders                  int
	latestBlockTimestamp       *time.Time
}

type Config struct {
	MeshDB               *meshdb.MeshDB
	BlockWatcher         *blockwatch.Watcher
	OrderValidator       *ordervalidator.OrderValidator
	ChainID              int
	MaxOrders            int
	MaxExpirationTime    *big.Int
	LatestBlockTimestamp *time.Time
}

// New instantiates a new order watcher
func New(config Config) (*Watcher, error) {
	decoder, err := decoder.New()
	if err != nil {
		return nil, err
	}
	assetDataDecoder := zeroex.NewAssetDataDecoder()
	contractAddresses, err := ethereum.GetContractAddressesForChainID(config.ChainID)
	if err != nil {
		return nil, err
	}

	// Validate config.
	if config.MaxOrders == 0 {
		return nil, errors.New("config.MaxOrders is required and cannot be zero")
	}
	if config.MaxExpirationTime == nil {
		return nil, errors.New("config.MaxExpirationTime is required and cannot be nil")
	} else if big.NewInt(time.Now().Unix()).Cmp(config.MaxExpirationTime) == 1 {
		// MaxExpirationTime should never be in the past.
		config.MaxExpirationTime = big.NewInt(time.Now().Unix())
	}

	// Configure a SlowCounter to be used for increasing max expiration time.
	slowCounterConfig := slowcounter.Config{
		Offset:   big.NewInt(slowCounterOffset),
		Rate:     slowCounterRate,
		Interval: slowCounterInterval,
		MaxCount: constants.UnlimitedExpirationTime,
	}
	maxExpirationCounter, err := slowcounter.New(slowCounterConfig, config.MaxExpirationTime)
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		meshDB:                     config.MeshDB,
		blockWatcher:               config.BlockWatcher,
		expirationWatcher:          expirationwatch.New(),
		contractAddressToSeenCount: map[common.Address]uint{},
		orderValidator:             config.OrderValidator,
		eventDecoder:               decoder,
		assetDataDecoder:           assetDataDecoder,
		contractAddresses:          contractAddresses,
		maxExpirationTime:          big.NewInt(0).Set(config.MaxExpirationTime),
		maxExpirationCounter:       maxExpirationCounter,
		maxOrders:                  config.MaxOrders,
		latestBlockTimestamp:       config.LatestBlockTimestamp,
	}

	// Check if any orders need to be removed right away due to high expiration
	// times.
	if err := w.decreaseMaxExpirationTimeIfNeeded(); err != nil {
		return nil, err
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

	// Start four independent goroutines. The main loop, cleanup loop, removed orders
	// checker and max expirationTime checker. Use four separate channels to communicate errors.
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
	maxExpirationTimeLoopErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		maxExpirationTimeLoopErrChan <- w.maxExpirationTimeLoop(innerCtx)
	}()
	removedCheckerLoopErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		removedCheckerLoopErrChan <- w.removedCheckerLoop(innerCtx)
	}()

	// If any error channel returns a non-nil error, we cancel the inner context
	// and return the error. Note that this means we only return the first error
	// that occurs.
	select {
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
	case err := <-maxExpirationTimeLoopErrChan:
		if err != nil {
			cancel()
			return err
		}
	case err := <-removedCheckerLoopErrChan:
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

func (w *Watcher) maxExpirationTimeLoop(ctx context.Context) error {
	ticker := time.NewTicker(maxExpirationTimeCheckInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			if err := w.increaseMaxExpirationTimeIfPossible(); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) removedCheckerLoop(ctx context.Context) error {
	for {
		start := time.Now()
		if err := w.permanentlyDeleteStaleRemovedOrders(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		// Wait minRemovedCheckInterval before calling permanentlyDeleteStaleRemovedOrders again. Since
		// we only start waiting _after_ permanentlyDeleteStaleRemovedOrders completes, we will never
		// have multiple calls to permanentlyDeleteStaleRemovedOrders running in parallel
		case <-time.After(minRemovedCheckInterval - time.Since(start)):
			continue
		}
	}
}

func (w *Watcher) generateExpirationOrderEventsIfAny(ordersColTxn *db.Transaction, latestBlockTimestamp time.Time) ([]*zeroex.OrderEvent, error) {
	orderEvents := []*zeroex.OrderEvent{}

	if w.latestBlockTimestamp == nil || w.latestBlockTimestamp.Before(latestBlockTimestamp) {
		expiredOrders := w.expirationWatcher.Prune(latestBlockTimestamp)
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
			w.unwatchOrder(ordersColTxn, order, order.FillableTakerAssetAmount)

			orderEvent := &zeroex.OrderEvent{
				OrderHash:                common.HexToHash(expiredOrder.ID),
				SignedOrder:              order.SignedOrder,
				FillableTakerAssetAmount: big.NewInt(0),
				EndState:                 zeroex.ESOrderExpired,
			}
			orderEvents = append(orderEvents, orderEvent)
		}
	} else {
		// A block re-org happened resulting in the latest block timestamp being
		// lower than on the previous latest block. We need to "unexpire" any orders
		// that have now become valid again as a result.
		removedOrders, err := w.meshDB.FindRemovedOrders()
		if err != nil {
			return nil, err
		}
		for _, order := range removedOrders {
			// Orders removed due to expiration have non-zero FillableTakerAssetAmounts
			if order.FillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
				continue
			}
			expiration := time.Unix(order.SignedOrder.ExpirationTimeSeconds.Int64(), 0)
			if latestBlockTimestamp.Before(expiration) {
				w.rewatchOrder(ordersColTxn, order, order.FillableTakerAssetAmount)
				orderEvent := &zeroex.OrderEvent{
					OrderHash:                order.Hash,
					SignedOrder:              order.SignedOrder,
					FillableTakerAssetAmount: order.FillableTakerAssetAmount,
					EndState:                 zeroex.ESOrderUnexpired,
				}
				orderEvents = append(orderEvents, orderEvent)
			}
		}
	}

	return orderEvents, nil
}

func (w *Watcher) handleBlockEvents(events []*blockwatch.Event) error {
	ordersColTxn := w.meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersColTxn.Discard()
	}()
	orderHashToDBOrder := map[common.Hash]*meshdb.Order{}
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{}
	var latestBlockNumber rpc.BlockNumber
	var latestBlockTimestamp time.Time
	for _, event := range events {
		latestBlockNumber = rpc.BlockNumber(event.BlockHeader.Number.Int64())
		latestBlockTimestamp = event.BlockHeader.Timestamp
		for _, log := range event.BlockHeader.Logs {
			var parameters zeroex.ContractEventParameters
			eventType, err := w.eventDecoder.FindEventType(log)
			if err != nil {
				switch err.(type) {
				case decoder.UntrackedTokenError:
					continue
				case decoder.UnsupportedEventError:
					logger.WithFields(logger.Fields{
						"topics":          err.(decoder.UnsupportedEventError).Topics,
						"contractAddress": err.(decoder.UnsupportedEventError).ContractAddress,
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
				var transferEvent decoder.ERC20TransferEvent
				err = w.eventDecoder.Decode(log, &transferEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = transferEvent
				fromOrders, err := w.findOrders(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrders(transferEvent.To, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, toOrders...)

			case "ERC20ApprovalEvent":
				var approvalEvent decoder.ERC20ApprovalEvent
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
				parameters = approvalEvent
				orders, err = w.findOrders(approvalEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "ERC721TransferEvent":
				var transferEvent decoder.ERC721TransferEvent
				err = w.eventDecoder.Decode(log, &transferEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = transferEvent
				fromOrders, err := w.findOrders(transferEvent.From, log.Address, transferEvent.TokenId)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrders(transferEvent.To, log.Address, transferEvent.TokenId)
				if err != nil {
					return err
				}
				orders = append(orders, toOrders...)

			case "ERC721ApprovalEvent":
				var approvalEvent decoder.ERC721ApprovalEvent
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
				parameters = approvalEvent
				orders, err = w.findOrders(approvalEvent.Owner, log.Address, approvalEvent.TokenId)
				if err != nil {
					return err
				}

			case "ERC721ApprovalForAllEvent":
				var approvalForAllEvent decoder.ERC721ApprovalForAllEvent
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
				parameters = approvalForAllEvent
				orders, err = w.findOrders(approvalForAllEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "ERC1155TransferSingleEvent":
				var transferEvent decoder.ERC1155TransferSingleEvent
				err = w.eventDecoder.Decode(log, &transferEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				// HACK(fabio): Currently we simply revalidate all orders involving assets in this
				// ERC1155 contract from this particular maker. We could however revalidate fewer orders
				// by also taking into account the `ID` of the assets affected. We punt on this for now
				// in order to support Augur's use-case of a dummy ERC1155 contract. In their case, we
				// need to revalidate all maker orders within the single ERC1155 contract and cannot optimize
				// further. In the future, we might want to special-case this broader approach for the Augur
				// contract address specifically.
				parameters = transferEvent
				fromOrders, err := w.findOrders(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrders(transferEvent.To, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, toOrders...)

			case "ERC1155TransferBatchEvent":
				var transferEvent decoder.ERC1155TransferBatchEvent
				err = w.eventDecoder.Decode(log, &transferEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = transferEvent
				fromOrders, err := w.findOrders(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrders(transferEvent.To, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, toOrders...)

			case "ERC1155ApprovalForAllEvent":
				var approvalForAllEvent decoder.ERC1155ApprovalForAllEvent
				err = w.eventDecoder.Decode(log, &approvalForAllEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				// Ignores approvals set to anyone except the AssetProxy
				if approvalForAllEvent.Operator != w.contractAddresses.ERC1155Proxy {
					continue
				}
				parameters = approvalForAllEvent
				orders, err = w.findOrders(approvalForAllEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "WethWithdrawalEvent":
				var withdrawalEvent decoder.WethWithdrawalEvent
				err = w.eventDecoder.Decode(log, &withdrawalEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = withdrawalEvent
				orders, err = w.findOrders(withdrawalEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "WethDepositEvent":
				var depositEvent decoder.WethDepositEvent
				err = w.eventDecoder.Decode(log, &depositEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = depositEvent
				orders, err = w.findOrders(depositEvent.Owner, log.Address, nil)
				if err != nil {
					return err
				}

			case "ExchangeFillEvent":
				var exchangeFillEvent decoder.ExchangeFillEvent
				err = w.eventDecoder.Decode(log, &exchangeFillEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = exchangeFillEvent
				order := w.findOrder(exchangeFillEvent.OrderHash)
				if order != nil {
					orders = append(orders, order)
				}

			case "ExchangeCancelEvent":
				var exchangeCancelEvent decoder.ExchangeCancelEvent
				err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = exchangeCancelEvent
				orders = []*meshdb.Order{}
				order := w.findOrder(exchangeCancelEvent.OrderHash)
				if order != nil {
					orders = append(orders, order)
				}

			case "ExchangeCancelUpToEvent":
				var exchangeCancelUpToEvent decoder.ExchangeCancelUpToEvent
				err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
				if err != nil {
					if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
						continue
					}
					return err
				}
				parameters = exchangeCancelUpToEvent
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
			contractEvent := &zeroex.ContractEvent{
				BlockHash:  log.BlockHash,
				TxHash:     log.TxHash,
				TxIndex:    log.TxIndex,
				LogIndex:   log.Index,
				IsRemoved:  log.Removed,
				Address:    log.Address,
				Kind:       eventType,
				Parameters: parameters,
			}
			for _, order := range orders {
				orderHashToDBOrder[order.Hash] = order
				if _, ok := orderHashToEvents[order.Hash]; !ok {
					orderHashToEvents[order.Hash] = []*zeroex.ContractEvent{contractEvent}
				} else {
					orderHashToEvents[order.Hash] = append(orderHashToEvents[order.Hash], contractEvent)
				}
			}
		}
	}

	orderEvents, err := w.generateExpirationOrderEventsIfAny(ordersColTxn, latestBlockTimestamp)
	if err != nil {
		return err
	}
	w.latestBlockTimestamp = &latestBlockTimestamp

	// This timeout of 1min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, done := context.WithTimeout(context.Background(), 1*time.Minute)
	defer done()
	moreOrderEvents, err := w.generateOrderEventsIfChanged(ctx, ordersColTxn, orderHashToDBOrder, orderHashToEvents, latestBlockNumber)
	if err != nil {
		return err
	}

	if err := ordersColTxn.Commit(); err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("Failed to commit orders collection transaction")
	}

	orderEvents = append(orderEvents, moreOrderEvents...)
	if len(orderEvents) > 0 {
		w.orderFeed.Send(orderEvents)
	}

	return nil
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
	orderHashToDBOrder := map[common.Hash]*meshdb.Order{}
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{} // No events when running cleanup job
	for _, order := range orders {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		orderHashToDBOrder[order.Hash] = order
		orderHashToEvents[order.Hash] = []*zeroex.ContractEvent{}
	}
	// This timeout of 30min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	orderEvents, err := w.generateOrderEventsIfChanged(ctx, ordersColTxn, orderHashToDBOrder, orderHashToEvents, rpc.LatestBlockNumber)
	if err != nil {
		return err
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

func (w *Watcher) permanentlyDeleteStaleRemovedOrders(ctx context.Context) error {
	removedOrders, err := w.meshDB.FindRemovedOrders()
	if err != nil {
		return err
	}

	for _, order := range removedOrders {
		if time.Since(order.LastUpdated) > permanentlyDeleteAfter {
			if err := w.permanentlyDeleteOrder(w.meshDB.Orders, order); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

// Add adds a 0x order to the DB and watches it for changes in fillability. It
// will no-op (and return nil) if the order has already been added. If pinned is
// true, the orders will be marked as pinned. Pinned orders will not be affected
// by any DDoS prevention or incentive mechanisms and will always stay in
// storage until they are no longer fillable.
func (w *Watcher) Add(orderInfo *ordervalidator.AcceptedOrderInfo, pinned bool) error {
	if err := w.decreaseMaxExpirationTimeIfNeeded(); err != nil {
		return err
	}

	// TODO(albrow): technically we should count the current number of orders,
	// remove some if needed, and then insert the order in a single transaction to
	// ensure that we don't accidentally exceed the maximum. In practice, and
	// because of the way OrderWatcher works, the distinction shouldn't matter.
	txn := w.meshDB.Orders.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	// Final expiration time check before inserting the order. We might have just
	// changed max expiration time above.
	if !pinned && orderInfo.SignedOrder.ExpirationTimeSeconds.Cmp(w.maxExpirationTime) == 1 {
		// HACK(albrow): This is technically not the ideal way to respond to this
		// situation, but it is a lot easier to implement for the time being. In the
		// future, we should return an error and then react to that error
		// differently depending on whether the order was received via RPC or from a
		// peer. In the former case, we should return an RPC error response
		// indicating that the order was not in fact added. In the latter case, we
		// should effectively no-op, neither penalizing the peer or emitting any
		// order events. For now, we respond by emitting an ADDED event immediately
		// followed by a STOPPED_WATCHING event. If this order was submitted via
		// RPC, the RPC client will see a response that indicates the order was
		// successfully added, and then it will look like we immediately stopped
		// watching it. This is not too far off from what really happened but is
		// slightly inefficient.
		addedEvent := &zeroex.OrderEvent{
			OrderHash:                orderInfo.OrderHash,
			SignedOrder:              orderInfo.SignedOrder,
			FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
			EndState:                 zeroex.ESOrderAdded,
		}
		stoppedWatchingEvent := &zeroex.OrderEvent{
			OrderHash:                orderInfo.OrderHash,
			SignedOrder:              orderInfo.SignedOrder,
			FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
			EndState:                 zeroex.ESStoppedWatching,
		}
		w.orderFeed.Send([]*zeroex.OrderEvent{addedEvent, stoppedWatchingEvent})
		return nil
	}

	order := &meshdb.Order{
		Hash:                     orderInfo.OrderHash,
		SignedOrder:              orderInfo.SignedOrder,
		LastUpdated:              time.Now().UTC(),
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
		IsRemoved:                false,
		IsPinned:                 pinned,
	}
	err := txn.Insert(order)
	if err != nil {
		if _, ok := err.(db.AlreadyExistsError); ok {
			// If we're already watching the order, that's fine in this case. Don't
			// return an error.
			return nil
		}
		return err
	}
	if err := txn.Commit(); err != nil {
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
		EndState:                 zeroex.ESOrderAdded,
	}
	w.orderFeed.Send([]*zeroex.OrderEvent{orderEvent})

	return nil
}

func (w *Watcher) trimOrdersAndFireEvents() error {
	targetMaxOrders := int(maxOrdersTrimRatio * float64(w.maxOrders))
	newMaxExpirationTime, removedOrders, err := w.meshDB.TrimOrdersByExpirationTime(targetMaxOrders)
	if err != nil {
		return err
	}
	if len(removedOrders) > 0 {
		logger.WithFields(logger.Fields{
			"numOrdersRemoved": len(removedOrders),
			"targetMaxOrders":  targetMaxOrders,
		}).Debug("removing orders to make space")
	}
	for _, removedOrder := range removedOrders {
		// Fire a "STOPPED_WATCHING" event for each order that was removed.
		orderEvent := &zeroex.OrderEvent{
			OrderHash:                removedOrder.Hash,
			SignedOrder:              removedOrder.SignedOrder,
			FillableTakerAssetAmount: removedOrder.FillableTakerAssetAmount,
			EndState:                 zeroex.ESStoppedWatching,
		}
		w.orderFeed.Send([]*zeroex.OrderEvent{orderEvent})

		// Remove in-memory state
		expirationTimestamp := time.Unix(removedOrder.SignedOrder.ExpirationTimeSeconds.Int64(), 0)
		w.expirationWatcher.Remove(expirationTimestamp, removedOrder.Hash.Hex())
		err = w.removeAssetDataAddressFromEventDecoder(removedOrder.SignedOrder.MakerAssetData)
		if err != nil {
			// This should never happen since the same error would have happened when adding
			// the assetData to the EventDecoder.
			logger.WithFields(logger.Fields{
				"error":       err.Error(),
				"signedOrder": removedOrder.SignedOrder,
			}).Error("Unexpected error when trying to remove an assetData from decoder")
			return err
		}
	}
	if newMaxExpirationTime.Cmp(w.maxExpirationTime) == -1 {
		// Decrease the max expiration time to account for the fact that orders were
		// removed.
		logger.WithFields(logger.Fields{
			"oldMaxExpirationTime": w.maxExpirationTime.String(),
			"newMaxExpirationTime": newMaxExpirationTime.String(),
		}).Debug("decreasing max expiration time")
		w.maxExpirationTime = newMaxExpirationTime
		w.maxExpirationCounter.Reset(newMaxExpirationTime)
		w.saveMaxExpirationTime(newMaxExpirationTime)
	}

	return nil
}

// MaxExpirationTime returns the current maximum expiration time for incoming
// orders.
func (w *Watcher) MaxExpirationTime() *big.Int {
	return w.maxExpirationTime
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

	expirationTimestamp := time.Unix(signedOrder.ExpirationTimeSeconds.Int64(), 0).UTC()
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

func (w *Watcher) generateOrderEventsIfChanged(ctx context.Context, ordersColTxn *db.Transaction, orderHashToDBOrder map[common.Hash]*meshdb.Order, orderHashToEvents map[common.Hash][]*zeroex.ContractEvent, validationBlockNumber rpc.BlockNumber) ([]*zeroex.OrderEvent, error) {
	signedOrders := []*zeroex.SignedOrder{}
	for _, order := range orderHashToDBOrder {
		if order.IsRemoved && time.Since(order.LastUpdated) > permanentlyDeleteAfter {
			if err := w.permanentlyDeleteOrder(ordersColTxn, order); err != nil {
				return nil, err
			}
			continue
		}
		signedOrders = append(signedOrders, order.SignedOrder)
	}
	if len(signedOrders) == 0 {
		return nil, nil
	}
	areNewOrders := false
	validationResults := w.orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, validationBlockNumber)

	orderEvents := []*zeroex.OrderEvent{}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		order := orderHashToDBOrder[acceptedOrderInfo.OrderHash]
		oldFillableAmount := order.FillableTakerAssetAmount
		newFillableAmount := acceptedOrderInfo.FillableTakerAssetAmount
		oldAmountIsMoreThenNewAmount := oldFillableAmount.Cmp(newFillableAmount) == 1

		expirationTime := time.Unix(order.SignedOrder.ExpirationTimeSeconds.Int64(), 0)
		isExpired := time.Now().After(expirationTime)
		if !isExpired && oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
			// A previous event caused this order to be removed from DB because it's
			// fillableAmount became 0, but it has now been revived (e.g., block re-org
			// causes order fill txn to get reverted). We need to re-add order and emit an event.
			w.rewatchOrder(ordersColTxn, order, acceptedOrderInfo.FillableTakerAssetAmount)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				EndState:                 zeroex.ESOrderAdded,
				ContractEvents:           orderHashToEvents[order.Hash],
			}
			orderEvents = append(orderEvents, orderEvent)
		} else if oldFillableAmount.Cmp(newFillableAmount) == 0 {
			// No important state-change happened
		} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && oldAmountIsMoreThenNewAmount {
			// Order was filled, emit  event and update order in DB
			order.FillableTakerAssetAmount = newFillableAmount
			w.updateOrderDBEntry(ordersColTxn, order)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				EndState:                 zeroex.ESOrderFilled,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				ContractEvents:           orderHashToEvents[order.Hash],
			}
			orderEvents = append(orderEvents, orderEvent)
		} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && !oldAmountIsMoreThenNewAmount {
			// The order is now fillable for more then it was before. E.g.: A fill txn reverted (block-reorg)
			// Update order in DB and emit event
			order.FillableTakerAssetAmount = newFillableAmount
			w.updateOrderDBEntry(ordersColTxn, order)
			orderEvent := &zeroex.OrderEvent{
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder,
				EndState:                 zeroex.ESOrderFillabilityIncreased,
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				ContractEvents:           orderHashToEvents[order.Hash],
			}
			orderEvents = append(orderEvents, orderEvent)
		}
	}
	for _, rejectedOrderInfo := range validationResults.Rejected {
		switch rejectedOrderInfo.Kind {
		case ordervalidator.MeshError:
			// TODO(fabio): Do we want to handle MeshErrors somehow here?
		case ordervalidator.ZeroExValidation:
			order := orderHashToDBOrder[rejectedOrderInfo.OrderHash]
			oldFillableAmount := order.FillableTakerAssetAmount
			if oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
				// If the oldFillableAmount was already 0, this order is already flagged for removal.
			} else {
				// If oldFillableAmount > 0, it got fullyFilled, cancelled, expired or unfunded, emit event
				w.unwatchOrder(ordersColTxn, order, big.NewInt(0))
				endState, ok := ordervalidator.ConvertRejectOrderCodeToOrderEventEndState(rejectedOrderInfo.Status)
				if !ok {
					err := fmt.Errorf("no OrderEventEndState corresponding to RejectedOrderStatus: %q", rejectedOrderInfo.Status)
					logger.WithError(err).WithField("rejectedOrderStatus", rejectedOrderInfo.Status).Error("no OrderEventEndState corresponding to RejectedOrderStatus")
					return nil, err
				}
				orderEvent := &zeroex.OrderEvent{
					OrderHash:                rejectedOrderInfo.OrderHash,
					SignedOrder:              rejectedOrderInfo.SignedOrder,
					FillableTakerAssetAmount: big.NewInt(0),
					EndState:                 endState,
					ContractEvents:           orderHashToEvents[order.Hash],
				}
				orderEvents = append(orderEvents, orderEvent)
			}
		default:
			err := fmt.Errorf("unknown rejectedOrderInfo.Kind: %q", rejectedOrderInfo.Kind)
			logger.WithError(err).Error("encountered unhandled rejectedOrderInfo.Kind value")
			return nil, err
		}
	}

	return orderEvents, nil
}

type orderUpdater interface {
	Update(model db.Model) error
}

func (w *Watcher) updateOrderDBEntry(u orderUpdater, order *meshdb.Order) {
	order.LastUpdated = time.Now().UTC()
	err := u.Update(order)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}
}

func (w *Watcher) rewatchOrder(u orderUpdater, order *meshdb.Order, fillableTakerAssetAmount *big.Int) {
	order.IsRemoved = false
	order.LastUpdated = time.Now().UTC()
	order.FillableTakerAssetAmount = fillableTakerAssetAmount
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

func (w *Watcher) unwatchOrder(u orderUpdater, order *meshdb.Order, newFillableAmount *big.Int) {
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

type orderDeleter interface {
	Delete(id []byte) error
}

func (w *Watcher) permanentlyDeleteOrder(deleter orderDeleter, order *meshdb.Order) error {
	err := deleter.Delete(order.Hash.Bytes())
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
	if _, ok := err.(decoder.UnsupportedEventError); ok {
		logger.WithFields(logger.Fields{
			"eventType":       eventType,
			"topics":          err.(decoder.UnsupportedEventError).Topics,
			"contractAddress": err.(decoder.UnsupportedEventError).ContractAddress,
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
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC1155(decodedAssetData.Address)
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
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.contractAddressToSeenCount[decodedAssetData.Address] = w.contractAddressToSeenCount[decodedAssetData.Address] - 1
		if w.contractAddressToSeenCount[decodedAssetData.Address] == 0 {
			w.eventDecoder.RemoveKnownERC1155(decodedAssetData.Address)
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

func (w *Watcher) decreaseMaxExpirationTimeIfNeeded() error {
	if orderCount, err := w.meshDB.Orders.Count(); err != nil {
		return err
	} else if orderCount+1 > w.maxOrders {
		if err := w.trimOrdersAndFireEvents(); err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) increaseMaxExpirationTimeIfPossible() error {
	if orderCount, err := w.meshDB.Orders.Count(); err != nil {
		return err
	} else if orderCount < w.maxOrders {
		// We have enough space for new orders. Set the new max expiration time to the
		// value of slow counter.
		newMaxExpiration := w.maxExpirationCounter.Count()
		if w.maxExpirationTime.Cmp(newMaxExpiration) != 0 {
			logger.WithFields(logger.Fields{
				"oldMaxExpirationTime": w.maxExpirationTime.String(),
				"newMaxExpirationTime": fmt.Sprint(newMaxExpiration),
			}).Debug("increasing max expiration time")
			w.maxExpirationTime.Set(newMaxExpiration)
			w.saveMaxExpirationTime(newMaxExpiration)
		}
	}

	return nil
}

// saveMaxExpirationTime saves the new max expiration time in the database.
func (w *Watcher) saveMaxExpirationTime(maxExpirationTime *big.Int) {
	if err := w.meshDB.UpdateMetadata(func(metadata meshdb.Metadata) meshdb.Metadata {
		metadata.MaxExpirationTime = maxExpirationTime
		return metadata
	}); err != nil {
		logger.WithError(err).Error("could not update max expiration time in database")
	}
}
