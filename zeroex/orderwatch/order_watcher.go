package orderwatch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/expirationwatch"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/slowcounter"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

	// defaultLastUpdatedBuffer specifies how long it must have been since an order was
	// last updated in order to be re-validated by the cleanup worker
	defaultLastUpdatedBuffer = 30 * time.Minute

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

	// maxBlockEventsToHandle is the max number of block events we want to process in a single
	// call to `handleBlockEvents`
	maxBlockEventsToHandle = 500

	// configuration options for the SlowCounter used for increasing max
	// expiration time. Effectively, we will increase every 5 minutes as long as
	// there is enough space in the database for orders. The first increase will
	// be 5 seconds and the amount doubles from there (second increase will be 10
	// seconds, then 20 seconds, then 40, etc.)
	slowCounterOffset   = 5 // seconds
	slowCounterRate     = 2.0
	slowCounterInterval = 5 * time.Minute
)

var errNoBlocksStored = errors.New("no blocks were stored in the database")

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	db                         *db.DB
	blockWatcher               *blockwatch.Watcher
	eventDecoder               *decoder.Decoder
	assetDataDecoder           *zeroex.AssetDataDecoder
	blockSubscription          event.Subscription
	blockEventsChan            chan []*blockwatch.Event
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
	handleBlockEventsMu        sync.RWMutex
	// atLeastOneBlockProcessed is closed to signal that the BlockWatcher has processed at least one
	// block. Validation of orders should block until this has completed
	atLeastOneBlockProcessed   chan struct{}
	atLeastOneBlockProcessedMu sync.Mutex
	didProcessABlock           bool
}

type Config struct {
	DB                *db.DB
	BlockWatcher      *blockwatch.Watcher
	OrderValidator    *ordervalidator.OrderValidator
	ChainID           int
	ContractAddresses ethereum.ContractAddresses
	MaxOrders         int
	MaxExpirationTime *big.Int
}

// New instantiates a new order watcher
func New(config Config) (*Watcher, error) {
	decoder, err := decoder.New()
	if err != nil {
		return nil, err
	}
	assetDataDecoder := zeroex.NewAssetDataDecoder()

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
		db:                         config.DB,
		blockWatcher:               config.BlockWatcher,
		expirationWatcher:          expirationwatch.New(),
		contractAddressToSeenCount: map[common.Address]uint{},
		orderValidator:             config.OrderValidator,
		eventDecoder:               decoder,
		assetDataDecoder:           assetDataDecoder,
		contractAddresses:          config.ContractAddresses,
		maxExpirationTime:          big.NewInt(0).Set(config.MaxExpirationTime),
		maxExpirationCounter:       maxExpirationCounter,
		maxOrders:                  config.MaxOrders,
		blockEventsChan:            make(chan []*blockwatch.Event, 100),
		atLeastOneBlockProcessed:   make(chan struct{}),
		didProcessABlock:           false,
	}

	// Pre-populate the OrderWatcher with all orders already stored in the DB
	orders, err := w.db.FindOrders(nil)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		err := w.setupInMemoryOrderState(order)
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
	w.blockSubscription = w.blockWatcher.Subscribe(w.blockEventsChan)

	for {
		select {
		case <-ctx.Done():
			w.blockSubscription.Unsubscribe()
			close(w.blockEventsChan)
			return nil
		case err := <-w.blockSubscription.Err():
			logger.WithFields(logger.Fields{
				"error": err.Error(),
			}).Error("block subscription error encountered")
		case events := <-w.blockEventsChan:
			// Instead of simply processing the first array of events in the blockEventsChan,
			// we might as well process _all_ events in the channel.
			drainedEvents := drainBlockEventsChan(w.blockEventsChan, maxBlockEventsToHandle)
			events = append(events, drainedEvents...)
			w.handleBlockEventsMu.Lock()
			if err := w.handleBlockEvents(ctx, events); err != nil {
				w.handleBlockEventsMu.Unlock()
				return err
			}
			w.handleBlockEventsMu.Unlock()
		}
	}
}

func drainBlockEventsChan(blockEventsChan chan []*blockwatch.Event, max int) []*blockwatch.Event {
	allEvents := []*blockwatch.Event{}
Loop:
	for {
		select {
		case moreEvents := <-blockEventsChan:
			allEvents = append(allEvents, moreEvents...)
			if len(allEvents) >= max {
				break Loop
			}
		default:
			break Loop
		}
	}
	return allEvents
}

func (w *Watcher) cleanupLoop(ctx context.Context) error {
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(minCleanupInterval - time.Since(start)):
			// Wait minCleanupInterval before calling cleanup again. Since
			// we only start sleeping _after_ cleanup completes, we will never
			// have multiple calls to cleanup running in parallel
			break
		}

		start = time.Now()
		if err := w.Cleanup(ctx, defaultLastUpdatedBuffer); err != nil {
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

// TODO(albrow): Update this comment.
// handleOrderExpirations takes care of generating expired and unexpired order events for orders that do not require re-validation.
// Since expiry is now done according to block timestamp, we can figure out which orders have expired/unexpired statically. We do not
// process blocks that require re-validation, since the validation process will already emit the necessary events and we cannot make
// multiple updates to an order within a single DB transaction.
// latestBlockTimestamp is the latest block timestamp Mesh knows about
// previousLatestBlockTimestamp is the previous latest block timestamp Mesh knew about
// ordersToRevalidate contains all the orders Mesh needs to re-validate given the events emitted by the blocks processed
func (w *Watcher) handleOrderExpirations(latestBlockTimestamp, previousLatestBlockTimestamp time.Time, ordersToRevalidate map[common.Hash]*types.OrderWithMetadata) ([]*zeroex.OrderEvent, error) {
	orderEvents := []*zeroex.OrderEvent{}
	var defaultTime time.Time

	if previousLatestBlockTimestamp == defaultTime || previousLatestBlockTimestamp.Before(latestBlockTimestamp) {
		expiredOrders := w.expirationWatcher.Prune(latestBlockTimestamp)
		for _, expiredOrder := range expiredOrders {
			orderHash := common.HexToHash(expiredOrder.ID)
			// If we will re-validate this order, the revalidation process will discover that
			// it's expired, and an appropriate event will already be emitted
			if _, ok := ordersToRevalidate[orderHash]; ok {
				continue
			}
			order, err := w.db.GetOrder(orderHash)
			if err != nil {
				logger.WithFields(logger.Fields{
					"error":     err.Error(),
					"orderHash": expiredOrder.ID,
				}).Trace("Order expired that was no longer in DB")
				continue
			}
			w.unwatchOrder(order, nil)

			orderEvent := &zeroex.OrderEvent{
				Timestamp:                latestBlockTimestamp,
				OrderHash:                common.HexToHash(expiredOrder.ID),
				SignedOrder:              order.SignedOrder(),
				FillableTakerAssetAmount: big.NewInt(0),
				EndState:                 zeroex.ESOrderExpired,
			}
			orderEvents = append(orderEvents, orderEvent)
		}
	} else if previousLatestBlockTimestamp.After(latestBlockTimestamp) {
		// A block re-org happened resulting in the latest block timestamp being
		// lower than on the previous latest block. We need to "unexpire" any orders
		// that have now become valid again as a result.
		removedOrders, err := w.db.FindOrders(&db.OrderQuery{
			Filters: []db.OrderFilter{
				{
					Field: db.OFIsRemoved,
					Kind:  db.Equal,
					Value: true,
				},
			},
		})
		if err != nil {
			return orderEvents, err
		}
		for _, order := range removedOrders {
			// Orders removed due to expiration have non-zero FillableTakerAssetAmounts
			if order.FillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
				continue
			}
			// If we will re-validate this order, the revalidation process will discover that
			// it's unexpired, and an appropriate event will already be emitted
			if _, ok := ordersToRevalidate[order.Hash]; ok {
				continue
			}
			expiration := time.Unix(order.SignedOrder().ExpirationTimeSeconds.Int64(), 0)
			if latestBlockTimestamp.Before(expiration) {
				w.rewatchOrder(order, order.FillableTakerAssetAmount)
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                latestBlockTimestamp,
					OrderHash:                order.Hash,
					SignedOrder:              order.SignedOrder(),
					FillableTakerAssetAmount: order.FillableTakerAssetAmount,
					EndState:                 zeroex.ESOrderUnexpired,
				}
				orderEvents = append(orderEvents, orderEvent)
			}
		}
	} else {
		// The block timestamp hasn't changed, noop
	}

	return orderEvents, nil
}

// handleBlockEvents processes a set of block events into order events for a set of orders.
// handleBlockEvents MUST only be called after acquiring a lock to the `handleBlockEventsMu` mutex.
func (w *Watcher) handleBlockEvents(
	ctx context.Context,
	events []*blockwatch.Event,
) error {
	if len(events) == 0 {
		return nil
	}

	// TODO(albrow): Consider implementing transactions.
	// miniHeadersColTxn := w.meshDB.MiniHeaders.OpenTransaction()
	// defer func() {
	// 	_ = miniHeadersColTxn.Discard()
	// }()
	// ordersColTxn := w.meshtypes.OrderWithMetadatas.OpenTransaction()
	// defer func() {
	// 	_ = ordersColTxn.Discard()
	// }()

	var previousLatestBlockTimestamp time.Time
	latestBlocks, err := w.db.FindMiniHeaders(&db.MiniHeaderQuery{
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
		Limit: 1,
	})
	if err != nil {
		return err
	}
	if len(latestBlocks) != 0 {
		previousLatestBlockTimestamp = latestBlocks[0].Timestamp
	}
	latestBlockNumber, latestBlockTimestamp := w.getBlockchainState(events)

	// TODO(albrow): Implement this.
	// err = updateBlockHeadersStoredInDB(miniHeadersColTxn, events)
	// if err != nil {
	// 	return err
	// }

	orderHashToDBOrder := map[common.Hash]*types.OrderWithMetadata{}
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{}
	for _, event := range events {
		for _, log := range event.BlockHeader.Logs {
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
			contractEvent := &zeroex.ContractEvent{
				BlockHash: log.BlockHash,
				TxHash:    log.TxHash,
				TxIndex:   log.TxIndex,
				LogIndex:  log.Index,
				IsRemoved: log.Removed,
				Address:   log.Address,
				Kind:      eventType,
			}
			orders := []*types.OrderWithMetadata{}
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
				contractEvent.Parameters = transferEvent
				fromOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.To, log.Address, nil)
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
				contractEvent.Parameters = approvalEvent
				orders, err = w.findOrdersByTokenAddressAndTokenID(approvalEvent.Owner, log.Address, nil)
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
				contractEvent.Parameters = transferEvent
				fromOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.From, log.Address, transferEvent.TokenId)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.To, log.Address, transferEvent.TokenId)
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
				contractEvent.Parameters = approvalEvent
				orders, err = w.findOrdersByTokenAddressAndTokenID(approvalEvent.Owner, log.Address, approvalEvent.TokenId)
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
				contractEvent.Parameters = approvalForAllEvent
				orders, err = w.findOrdersByTokenAddressAndTokenID(approvalForAllEvent.Owner, log.Address, nil)
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
				contractEvent.Parameters = transferEvent
				fromOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.To, log.Address, nil)
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
				contractEvent.Parameters = transferEvent
				fromOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.From, log.Address, nil)
				if err != nil {
					return err
				}
				orders = append(orders, fromOrders...)
				toOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.To, log.Address, nil)
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
				contractEvent.Parameters = approvalForAllEvent
				orders, err = w.findOrdersByTokenAddressAndTokenID(approvalForAllEvent.Owner, log.Address, nil)
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
				contractEvent.Parameters = withdrawalEvent
				orders, err = w.findOrdersByTokenAddressAndTokenID(withdrawalEvent.Owner, log.Address, nil)
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
				contractEvent.Parameters = depositEvent
				orders, err = w.findOrdersByTokenAddressAndTokenID(depositEvent.Owner, log.Address, nil)
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
				contractEvent.Parameters = exchangeFillEvent

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
				contractEvent.Parameters = exchangeCancelEvent
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
				contractEvent.Parameters = exchangeCancelUpToEvent
				cancelledOrders, err := w.db.FindOrders(&db.OrderQuery{
					Filters: []db.OrderFilter{
						{
							Field: db.OFMakerAddress,
							Kind:  db.Equal,
							Value: exchangeCancelUpToEvent.MakerAddress,
						},
						{
							Field: db.OFSalt,
							Kind:  db.LessOrEqual,
							Value: exchangeCancelUpToEvent.OrderEpoch,
						},
					},
				})
				if err != nil {
					logger.WithFields(logger.Fields{
						"error": err.Error(),
					}).Error("unexpected query error encountered")
					return err
				}
				orders = append(orders, cancelledOrders...)

			default:
				logger.WithFields(logger.Fields{
					"eventType": eventType,
					"log":       log,
				}).Error("unknown eventType encountered")
				return err
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

	expirationOrderEvents, err := w.handleOrderExpirations(latestBlockTimestamp, previousLatestBlockTimestamp, orderHashToDBOrder)
	if err != nil {
		return err
	}

	// This timeout of 1min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, done := context.WithTimeout(ctx, 1*time.Minute)
	defer done()
	postValidationOrderEvents, err := w.generateOrderEventsIfChanged(ctx, orderHashToDBOrder, orderHashToEvents, latestBlockNumber, latestBlockTimestamp)
	if err != nil {
		return err
	}

	// TODO(albrow): Consider using transactions.
	// if err := ordersColTxn.Commit(); err != nil {
	// 	logger.WithFields(logger.Fields{
	// 		"error": err.Error(),
	// 	}).Error("Failed to commit orders collection transaction")
	// 	return err
	// }
	// if err := miniHeadersColTxn.Commit(); err != nil {
	// 	logger.WithFields(logger.Fields{
	// 		"error": err.Error(),
	// 	}).Error("Failed to commit miniheaders collection transaction")
	// 	return err
	// }

	orderEvents := append(expirationOrderEvents, postValidationOrderEvents...)
	if len(orderEvents) > 0 {
		w.orderFeed.Send(orderEvents)
	}

	w.atLeastOneBlockProcessedMu.Lock()
	if !w.didProcessABlock {
		w.didProcessABlock = true
		close(w.atLeastOneBlockProcessed)
	}
	w.atLeastOneBlockProcessedMu.Unlock()

	// TODO(albrow): Remove this after implementing miniheader retention limit inside db package.
	// Since we might have added MiniHeaders to the DB, we need to prune any excess MiniHeaders stored
	// in the DB
	// err = w.meshDB.PruneMiniHeadersAboveRetentionLimit()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (w *Watcher) getLatestBlock() (*types.MiniHeader, error) {
	latestBlocks, err := w.db.FindMiniHeaders(&db.MiniHeaderQuery{
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(latestBlocks) == 0 {
		return nil, errNoBlocksStored
	}
	return latestBlocks[0], nil
}

// Cleanup re-validates all orders in DB which haven't been re-validated in
// `lastUpdatedBuffer` time to make sure all orders are still up-to-date
func (w *Watcher) Cleanup(ctx context.Context, lastUpdatedBuffer time.Duration) error {
	// Pause block event processing until we finished cleaning up at current block height
	w.handleBlockEventsMu.RLock()
	defer w.handleBlockEventsMu.RUnlock()

	// TODO(albrow): Do we need a transaction here?
	// ordersColTxn := w.meshtypes.OrderWithMetadatas.OpenTransaction()
	// defer func() {
	// 	_ = ordersColTxn.Discard()
	// }()
	lastUpdatedCutOff := time.Now().Add(-lastUpdatedBuffer)
	// orders, err := w.meshDB.FindOrdersLastUpdatedBefore(lastUpdatedCutOff)
	orders, err := w.db.FindOrders(&db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFLastUpdated,
				Kind:  db.Less,
				Value: lastUpdatedCutOff,
			},
		},
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error":             err.Error(),
			"lastUpdatedCutOff": lastUpdatedCutOff,
		}).Error("Failed to find orders by LastUpdatedBefore")
		return err
	}
	orderHashToDBOrder := map[common.Hash]*types.OrderWithMetadata{}
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

	latestBlock, err := w.getLatestBlock()
	if err != nil {
		return err
	}
	// This timeout of 30min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	orderEvents, err := w.generateOrderEventsIfChanged(ctx, orderHashToDBOrder, orderHashToEvents, latestBlock.Number, latestBlock.Timestamp)
	if err != nil {
		return err
	}

	// if err := ordersColTxn.Commit(); err != nil {
	// 	logger.WithFields(logger.Fields{
	// 		"error": err.Error(),
	// 	}).Error("Failed to commit orders collection transaction")
	// }

	if len(orderEvents) > 0 {
		w.orderFeed.Send(orderEvents)
	}

	return nil
}

func (w *Watcher) permanentlyDeleteStaleRemovedOrders(ctx context.Context) error {
	// TODO(albrow): This could be optimized by using a single query to delete
	// stale orders instead of finding them and deleting one-by-one. Limited by
	// the fact that we need to update in-memory state. When we remove in-memory
	// state we can revisit this.
	//
	// opts := &db.DeleteOrdersOpts{
	// 	Filters: []db.OrderFilter{
	// 		{
	// 			Field: db.OFIsRemoved,
	// 			Kind:  db.Equal,
	// 			Value: true,
	// 		},
	// 		{
	// 			Field: db.OFLastUpdated,
	// 			Kind:  db.Less,
	// 			Value: minLastUpdated,
	// 		},
	// 	},
	// }
	// return w.db.DeleteOrders(opts)

	// Find any orders marked as removed that have not been updated for a
	// long time. The cutoff time is determined by permanentlyDeleteAfter.
	minLastUpdated := time.Now().Add(-permanentlyDeleteAfter)
	opts := &db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFIsRemoved,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OFLastUpdated,
				Kind:  db.Less,
				Value: minLastUpdated,
			},
		},
	}
	ordersToDelete, err := w.db.FindOrders(opts)
	if err != nil {
		return err
	}
	for _, order := range ordersToDelete {
		if err := w.permanentlyDeleteOrder(order); err != nil {
			return err
		}
	}
	return nil
}

// add adds a 0x order to the DB and watches it for changes in fillability. It
// will no-op (and return nil) if the order has already been added. If pinned is
// true, the orders will be marked as pinned. Pinned orders will not be affected
// by any DDoS prevention or incentive mechanisms and will always stay in
// storage until they are no longer fillable.
func (w *Watcher) add(orderInfos []*ordervalidator.AcceptedOrderInfo, validationBlockNumber *big.Int, pinned bool) ([]*zeroex.OrderEvent, error) {
	now := time.Now().UTC()
	orderEvents := []*zeroex.OrderEvent{}
	dbOrders := []*types.OrderWithMetadata{}

	for _, orderInfo := range orderInfos {
		dbOrder, err := w.orderInfoToOrderWithMetadata(orderInfo, pinned, now)
		if err != nil {
			return nil, err
		}
		dbOrders = append(dbOrders, dbOrder)
	}

	// TODO(albrow): Should AddOrders return the new max expiration time?
	// Or is there a better way to do this?
	addedOrders, removedOrders, err := w.db.AddOrders(dbOrders)
	if err != nil {
		return nil, err
	}
	for _, order := range addedOrders {
		err = w.setupInMemoryOrderState(order)
		if err != nil {
			return orderEvents, err
		}
		addedEvent := &zeroex.OrderEvent{
			Timestamp:                now,
			OrderHash:                order.Hash,
			SignedOrder:              order.SignedOrder(),
			FillableTakerAssetAmount: order.FillableTakerAssetAmount,
			EndState:                 zeroex.ESOrderAdded,
		}
		orderEvents = append(orderEvents, addedEvent)
	}
	for _, order := range removedOrders {
		stoppedWatchingEvent := &zeroex.OrderEvent{
			Timestamp:                now,
			OrderHash:                order.Hash,
			SignedOrder:              order.SignedOrder(),
			FillableTakerAssetAmount: order.FillableTakerAssetAmount,
			EndState:                 zeroex.ESStoppedWatching,
		}
		orderEvents = append(orderEvents, stoppedWatchingEvent)

		// Remove in-memory state
		err = w.removeAssetDataAddressFromEventDecoder(order.MakerAssetData)
		if err != nil {
			// This should never happen since the same error would have happened when adding
			// the assetData to the EventDecoder.
			logger.WithFields(logger.Fields{
				"error":       err.Error(),
				"signedOrder": order.SignedOrder(),
			}).Error("Unexpected error when trying to remove an assetData from decoder")
			return nil, err
		}
	}

	// TODO(albrow): How to handle the edge case of orders that were not
	// added due to the max expiration time changing?

	return orderEvents, nil
}

func (w *Watcher) orderInfoToOrderWithMetadata(orderInfo *ordervalidator.AcceptedOrderInfo, pinned bool, now time.Time) (*types.OrderWithMetadata, error) {
	parsedMakerAssetData, err := db.ParseContractAddressesAndTokenIdsFromAssetData(orderInfo.SignedOrder.MakerAssetData, w.contractAddresses)
	if err != nil {
		return nil, err
	}
	parsedMakerFeeAssetData, err := db.ParseContractAddressesAndTokenIdsFromAssetData(orderInfo.SignedOrder.MakerFeeAssetData, w.contractAddresses)
	if err != nil {
		return nil, err
	}
	return &types.OrderWithMetadata{
		Hash:                     orderInfo.OrderHash,
		ChainID:                  orderInfo.SignedOrder.ChainID,
		ExchangeAddress:          orderInfo.SignedOrder.ExchangeAddress,
		MakerAddress:             orderInfo.SignedOrder.MakerAddress,
		MakerAssetData:           orderInfo.SignedOrder.MakerAssetData,
		MakerFeeAssetData:        orderInfo.SignedOrder.MakerFeeAssetData,
		MakerAssetAmount:         orderInfo.SignedOrder.MakerAssetAmount,
		MakerFee:                 orderInfo.SignedOrder.MakerFee,
		TakerAddress:             orderInfo.SignedOrder.TakerAddress,
		TakerAssetData:           orderInfo.SignedOrder.TakerAssetData,
		TakerFeeAssetData:        orderInfo.SignedOrder.TakerFeeAssetData,
		TakerAssetAmount:         orderInfo.SignedOrder.TakerAssetAmount,
		TakerFee:                 orderInfo.SignedOrder.TakerFee,
		SenderAddress:            orderInfo.SignedOrder.SenderAddress,
		FeeRecipientAddress:      orderInfo.SignedOrder.FeeRecipientAddress,
		ExpirationTimeSeconds:    orderInfo.SignedOrder.ExpirationTimeSeconds,
		Salt:                     orderInfo.SignedOrder.Salt,
		Signature:                orderInfo.SignedOrder.Signature,
		IsRemoved:                false,
		IsPinned:                 pinned,
		LastUpdated:              now,
		ParsedMakerAssetData:     parsedMakerAssetData,
		ParsedMakerFeeAssetData:  parsedMakerFeeAssetData,
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
	}, nil
}

// MaxExpirationTime returns the current maximum expiration time for incoming
// orders.
func (w *Watcher) MaxExpirationTime() *big.Int {
	return w.maxExpirationTime
}

// TODO(albrow): All in-memory state can be removed.
func (w *Watcher) setupInMemoryOrderState(order *types.OrderWithMetadata) error {
	w.eventDecoder.AddKnownExchange(order.ExchangeAddress)

	// Add MakerAssetData and MakerFeeAssetData to EventDecoder
	err := w.addAssetDataAddressToEventDecoder(order.MakerAssetData)
	if err != nil {
		return err
	}
	if order.MakerFee.Cmp(big.NewInt(0)) == 1 {
		err = w.addAssetDataAddressToEventDecoder(order.MakerFeeAssetData)
		if err != nil {
			return err
		}
	}

	expirationTimestamp := time.Unix(order.ExpirationTimeSeconds.Int64(), 0)
	w.expirationWatcher.Add(expirationTimestamp, order.Hash.Hex())

	return nil
}

// Subscribe allows one to subscribe to the order events emitted by the OrderWatcher.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*zeroex.OrderEvent) event.Subscription {
	return w.orderScope.Track(w.orderFeed.Subscribe(sink))
}

// TODO(albrow): Do we even still need this method?
func (w *Watcher) findOrder(orderHash common.Hash) *types.OrderWithMetadata {
	order, err := w.db.GetOrder(orderHash)
	if err != nil {
		// TODO(albrow): Implement NotFoundError
		// if _, ok := err.(db.NotFoundError); ok {
		// 	// short-circuit. We expect to receive events from orders we aren't actively tracking
		// 	return nil
		// }
		logger.WithFields(logger.Fields{
			"error":     err.Error(),
			"orderHash": orderHash,
		}).Warning("Unexpected error using FindByID for order")
		return nil
	}
	return order
}

// findOrdersByTokenAddressAndTokenID finds and returns all orders that have
// either a makerAsset or a makerFeeAsset matching the given tokenAddress and
// tokenID.
func (w *Watcher) findOrdersByTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*types.OrderWithMetadata, error) {
	ordersWithAffectedMakerAsset, err := w.db.FindOrders(&db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFMakerAddress,
				Kind:  db.Equal,
				Value: tokenAddress,
			},
			db.IncludesMakerAssetData(tokenAddress, tokenID),
		},
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}
	ordersWithAffectedMakerFeeAsset, err := w.db.FindOrders(&db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFMakerAddress,
				Kind:  db.Equal,
				Value: tokenAddress,
			},
			db.IncludesMakerFeeAssetData(tokenAddress, tokenID)},
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}

	return append(ordersWithAffectedMakerAsset, ordersWithAffectedMakerFeeAsset...), nil
}

func (w *Watcher) convertValidationResultsIntoOrderEvents(
	validationResults *ordervalidator.ValidationResults,
	orderHashToDBOrder map[common.Hash]*types.OrderWithMetadata,
	orderHashToEvents map[common.Hash][]*zeroex.ContractEvent,
	validationBlockTimestamp time.Time,
) ([]*zeroex.OrderEvent, error) {
	orderEvents := []*zeroex.OrderEvent{}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		order, found := orderHashToDBOrder[acceptedOrderInfo.OrderHash]
		if !found {
			logger.WithFields(logger.Fields{
				"unknownOrderHash":   acceptedOrderInfo.OrderHash,
				"validationResults":  validationResults,
				"orderHashToDBOrder": orderHashToDBOrder,
			}).Error("validationResults.Accepted contained unknown order hash")
			continue
		}
		oldFillableAmount := order.FillableTakerAssetAmount
		newFillableAmount := acceptedOrderInfo.FillableTakerAssetAmount
		oldAmountIsMoreThenNewAmount := oldFillableAmount.Cmp(newFillableAmount) == 1

		if oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
			// A previous event caused this order to be removed from DB because it's
			// fillableAmount became 0, but it has now been revived (e.g., block re-org
			// causes order fill txn to get reverted). We need to re-add order and emit an event.
			w.rewatchOrder(order, acceptedOrderInfo.FillableTakerAssetAmount)
			orderEvent := &zeroex.OrderEvent{
				Timestamp:                validationBlockTimestamp,
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder(),
				FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
				EndState:                 zeroex.ESOrderAdded,
				ContractEvents:           orderHashToEvents[order.Hash],
			}
			orderEvents = append(orderEvents, orderEvent)
		} else {
			expiration := time.Unix(order.SignedOrder().ExpirationTimeSeconds.Int64(), 0)

			if oldFillableAmount.Cmp(newFillableAmount) == 0 {
				// If order was previously expired, check if it has become unexpired
				if order.IsRemoved && oldFillableAmount.Cmp(big.NewInt(0)) != 0 && validationBlockTimestamp.Before(expiration) {
					w.rewatchOrder(order, nil)
					orderEvent := &zeroex.OrderEvent{
						Timestamp:                validationBlockTimestamp,
						OrderHash:                order.Hash,
						SignedOrder:              order.SignedOrder(),
						FillableTakerAssetAmount: order.FillableTakerAssetAmount,
						EndState:                 zeroex.ESOrderUnexpired,
					}
					orderEvents = append(orderEvents, orderEvent)
				}
				// No important state-change happened
				continue
			}
			if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && oldAmountIsMoreThenNewAmount {
				// If order was previously expired, check if it has become unexpired
				if order.IsRemoved && oldFillableAmount.Cmp(big.NewInt(0)) != 0 && validationBlockTimestamp.Before(expiration) {
					w.rewatchOrder(order, newFillableAmount)
					orderEvent := &zeroex.OrderEvent{
						Timestamp:                validationBlockTimestamp,
						OrderHash:                order.Hash,
						SignedOrder:              order.SignedOrder(),
						FillableTakerAssetAmount: order.FillableTakerAssetAmount,
						EndState:                 zeroex.ESOrderUnexpired,
					}
					orderEvents = append(orderEvents, orderEvent)
				} else {
					// TODO(albrow): Do we need to do this in a transaction?

				}
				// Order was filled, emit event
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                validationBlockTimestamp,
					OrderHash:                acceptedOrderInfo.OrderHash,
					SignedOrder:              order.SignedOrder(),
					EndState:                 zeroex.ESOrderFilled,
					FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
					ContractEvents:           orderHashToEvents[order.Hash],
				}
				orderEvents = append(orderEvents, orderEvent)
			} else if oldFillableAmount.Cmp(big.NewInt(0)) == 1 && !oldAmountIsMoreThenNewAmount {
				// The order is now fillable for more then it was before. E.g.: A fill txn reverted (block-reorg)
				// If order was previously expired, check if it has become unexpired
				if order.IsRemoved && oldFillableAmount.Cmp(big.NewInt(0)) != 0 && validationBlockTimestamp.Before(expiration) {
					w.rewatchOrder(order, newFillableAmount)
					orderEvent := &zeroex.OrderEvent{
						Timestamp:                validationBlockTimestamp,
						OrderHash:                order.Hash,
						SignedOrder:              order.SignedOrder(),
						FillableTakerAssetAmount: order.FillableTakerAssetAmount,
						EndState:                 zeroex.ESOrderUnexpired,
					}
					orderEvents = append(orderEvents, orderEvent)
				} else {
					w.updateOrderFillableTakerAssetAmount(order, newFillableAmount)
				}
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                validationBlockTimestamp,
					OrderHash:                acceptedOrderInfo.OrderHash,
					SignedOrder:              order.SignedOrder(),
					EndState:                 zeroex.ESOrderFillabilityIncreased,
					FillableTakerAssetAmount: acceptedOrderInfo.FillableTakerAssetAmount,
					ContractEvents:           orderHashToEvents[order.Hash],
				}
				orderEvents = append(orderEvents, orderEvent)
			}
		}
	}
	for _, rejectedOrderInfo := range validationResults.Rejected {
		switch rejectedOrderInfo.Kind {
		case ordervalidator.MeshError:
			// TODO(fabio): Do we want to handle MeshErrors somehow here?
		case ordervalidator.ZeroExValidation:
			order, found := orderHashToDBOrder[rejectedOrderInfo.OrderHash]
			if !found {
				logger.WithFields(logger.Fields{
					"unknownOrderHash":   rejectedOrderInfo.OrderHash,
					"validationResults":  validationResults,
					"orderHashToDBOrder": orderHashToDBOrder,
				}).Error("validationResults.Rejected contained unknown order hash")
				continue
			}
			oldFillableAmount := order.FillableTakerAssetAmount
			if oldFillableAmount.Cmp(big.NewInt(0)) == 0 {
				// If the oldFillableAmount was already 0, this order is already flagged for removal.
			} else {
				// If oldFillableAmount > 0, it got fullyFilled, cancelled, expired or unfunded
				w.unwatchOrder(order, big.NewInt(0))
				endState, ok := ordervalidator.ConvertRejectOrderCodeToOrderEventEndState(rejectedOrderInfo.Status)
				if !ok {
					err := fmt.Errorf("no OrderEventEndState corresponding to RejectedOrderStatus: %q", rejectedOrderInfo.Status)
					logger.WithError(err).WithField("rejectedOrderStatus", rejectedOrderInfo.Status).Error("no OrderEventEndState corresponding to RejectedOrderStatus")
					return nil, err
				}
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                validationBlockTimestamp,
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

func (w *Watcher) generateOrderEventsIfChanged(
	ctx context.Context,
	orderHashToDBOrder map[common.Hash]*types.OrderWithMetadata,
	orderHashToEvents map[common.Hash][]*zeroex.ContractEvent,
	validationBlockNumber *big.Int,
	validationBlockTimestamp time.Time,
) ([]*zeroex.OrderEvent, error) {
	signedOrders := []*zeroex.SignedOrder{}
	for _, order := range orderHashToDBOrder {
		if order.IsRemoved && time.Since(order.LastUpdated) > permanentlyDeleteAfter {
			if err := w.permanentlyDeleteOrder(order); err != nil {
				return nil, err
			}
			continue
		}
		signedOrders = append(signedOrders, order.SignedOrder())
	}
	if len(signedOrders) == 0 {
		return nil, nil
	}
	areNewOrders := false
	validationResults := w.orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, validationBlockNumber)

	return w.convertValidationResultsIntoOrderEvents(
		validationResults, orderHashToDBOrder, orderHashToEvents, validationBlockTimestamp,
	)
}

// ValidateAndStoreValidOrders applies general 0x validation and Mesh-specific validation to
// the given orders and if they are valid, adds them to the OrderWatcher
func (w *Watcher) ValidateAndStoreValidOrders(ctx context.Context, orders []*zeroex.SignedOrder, pinned bool, chainID int) (*ordervalidator.ValidationResults, error) {
	results, validMeshOrders, err := w.meshSpecificOrderValidation(orders, chainID)
	if err != nil {
		return nil, err
	}

	// Lock down the processing of additional block events until we've validated and added these new orders
	w.handleBlockEventsMu.RLock()
	defer w.handleBlockEventsMu.RUnlock()

	validationBlock, zeroexResults, err := w.onchainOrderValidation(ctx, validMeshOrders)
	if err != nil {
		return nil, err
	}
	results.Accepted = append(results.Accepted, zeroexResults.Accepted...)
	results.Rejected = append(results.Rejected, zeroexResults.Rejected...)

	// Filter out only the new orders.
	newOrderInfos := []*ordervalidator.AcceptedOrderInfo{}
	for _, acceptedOrderInfo := range results.Accepted {
		// If the order isn't new, we don't add to OrderWatcher.
		if acceptedOrderInfo.IsNew {
			newOrderInfos = append(newOrderInfos, acceptedOrderInfo)
		}
	}

	// Add the order to the OrderWatcher. This also saves the order in the
	// database.
	allOrderEvents := []*zeroex.OrderEvent{}
	orderEvents, err := w.add(newOrderInfos, validationBlock.Number, pinned)
	if err != nil {
		return nil, err
	}
	allOrderEvents = append(allOrderEvents, orderEvents...)

	if len(allOrderEvents) > 0 {
		// NOTE(albrow): Send can block if the subscriber(s) are slow. Blocking here can cause problems when Mesh is
		// shutting down, so to prevent that, we call Send in a goroutine and return immediately if the context
		// is done.
		done := make(chan interface{})
		go func() {
			w.orderFeed.Send(allOrderEvents)
			done <- struct{}{}
		}()
		select {
		case <-done:
			return results, nil
		case <-ctx.Done():
			return results, nil
		}
	}

	return results, nil
}

func (w *Watcher) onchainOrderValidation(ctx context.Context, orders []*zeroex.SignedOrder) (*types.MiniHeader, *ordervalidator.ValidationResults, error) {
	// HACK(fabio): While we wait for EIP-1898 support in Parity, we have no choice but to do the `eth_call`
	// at the latest known block _number_. As outlined in the `Rationale` section of EIP-1898, this approach cannot account
	// for the block being re-org'd out before the `eth_call` and then back in before the `eth_getBlockByNumber`
	// call (an unlikely but possible situation leading to an incorrect view of the world for these orders).
	// Unfortunately, this is the best we can do until EIP-1898 support in Parity.
	// Source: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1898.md#rationale
	latestMiniHeaders, err := w.db.FindMiniHeaders(&db.MiniHeaderQuery{
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, nil, err
	}
	if len(latestMiniHeaders) == 0 {
		return nil, nil, errors.New("orderwatch: Could not get latest miniheader")
	}
	latestMiniHeader := latestMiniHeaders[0]
	// This timeout of 1min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	areNewOrders := true
	zeroexResults := w.orderValidator.BatchValidate(ctx, orders, areNewOrders, latestMiniHeader.Number)
	return latestMiniHeader, zeroexResults, nil
}

func (w *Watcher) meshSpecificOrderValidation(orders []*zeroex.SignedOrder, chainID int) (*ordervalidator.ValidationResults, []*zeroex.SignedOrder, error) {
	results := &ordervalidator.ValidationResults{}
	validMeshOrders := []*zeroex.SignedOrder{}
	for _, order := range orders {
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			logger.WithField("error", err).Error("could not compute order hash")
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshError,
				Status:      ordervalidator.ROInternalError,
			})
			continue
		}
		if order.ExpirationTimeSeconds.Cmp(w.MaxExpirationTime()) == 1 {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROMaxExpirationExceeded,
			})
			continue
		}
		// Note(albrow): Orders with a sender address can be canceled or invalidated
		// off-chain which is difficult to support since we need to prune
		// canceled/invalidated orders from the database. We can special-case some
		// sender addresses over time. (For example we already have support for
		// validating Coordinator orders. What we're missing is a way to effeciently
		// remove orders that are soft-canceled via the Coordinator API).
		if order.SenderAddress != constants.NullAddress {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROSenderAddressNotAllowed,
			})
			continue
		}
		if order.ChainID.Cmp(big.NewInt(int64(chainID))) != 0 {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROIncorrectChain,
			})
			continue
		}
		if err == nil {
			// Only check the ExchangeAddress if we know the expected address for the
			// given chainID/networkID. If we don't know it, the order could still be
			// valid.
			expectedExchangeAddress := w.contractAddresses.Exchange
			if order.ExchangeAddress != expectedExchangeAddress {
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshValidation,
					Status:      ordervalidator.ROIncorrectExchangeAddress,
				})
				continue
			}
		}

		if err := validateOrderSize(order); err != nil {
			if err == constants.ErrMaxOrderSize {
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshValidation,
					Status:      ordervalidator.ROMaxOrderSizeExceeded,
				})
				continue
			} else {
				logger.WithField("error", err).Error("could not validate order size")
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshError,
					Status:      ordervalidator.ROInternalError,
				})
				continue
			}
		}

		// Check if order is already stored in DB
		dbOrder, err := w.db.GetOrder(orderHash)
		if err != nil {
			// TODO(albrow): Handle not found error.
			// if _, ok := err.(db.NotFoundError); !ok {
			// 	logger.WithField("error", err).Error("could not check if order was already stored")
			// 	return nil, nil, err
			// }
			// If the error is a db.NotFoundError, it just means the order is not currently stored in
			// the database. There's nothing else in the database to check, so we can continue.
		} else {
			// If stored but flagged for removal, reject it
			if dbOrder.IsRemoved {
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshValidation,
					Status:      ordervalidator.ROOrderAlreadyStoredAndUnfillable,
				})
				continue
			} else {
				// If stored but not flagged for removal, accept it without re-validation
				results.Accepted = append(results.Accepted, &ordervalidator.AcceptedOrderInfo{
					OrderHash:                orderHash,
					SignedOrder:              order,
					FillableTakerAssetAmount: dbOrder.FillableTakerAssetAmount,
					IsNew:                    false,
				})
				continue
			}
		}

		validMeshOrders = append(validMeshOrders, order)
	}

	return results, validMeshOrders, nil
}

func validateOrderSize(order *zeroex.SignedOrder) error {
	encoded, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if len(encoded) > constants.MaxOrderSizeInBytes {
		return constants.ErrMaxOrderSize
	}
	return nil
}

// TODO(albrow): Do we need to do this in a transaction?
func (w *Watcher) updateOrderFillableTakerAssetAmount(order *types.OrderWithMetadata, newFillableTakerAssetAmount *big.Int) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.LastUpdated = time.Now().UTC()
		orderToUpdate.FillableTakerAssetAmount = newFillableTakerAssetAmount
		return orderToUpdate, nil
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}
}

// TODO(albrow): Do we need to do this in a transaction?
func (w *Watcher) rewatchOrder(order *types.OrderWithMetadata, newFillableTakerAssetAmount *big.Int) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.IsRemoved = false
		orderToUpdate.LastUpdated = time.Now().UTC()
		if newFillableTakerAssetAmount != nil {
			orderToUpdate.FillableTakerAssetAmount = newFillableTakerAssetAmount
		}
		return orderToUpdate, nil
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}

	// Re-add order to expiration watcher
	// TODO(albrow): ExpirationTimeSeconds can overflow int64. This is probably a bug.
	expirationTimestamp := time.Unix(order.SignedOrder().ExpirationTimeSeconds.Int64(), 0)
	w.expirationWatcher.Add(expirationTimestamp, order.Hash.Hex())
}

// TODO(albrow): Do we need to do this in a transaction?
func (w *Watcher) unwatchOrder(order *types.OrderWithMetadata, newFillableAmount *big.Int) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.IsRemoved = true
		orderToUpdate.LastUpdated = time.Now().UTC()
		if newFillableAmount != nil {
			orderToUpdate.FillableTakerAssetAmount = newFillableAmount
		}
		return orderToUpdate, nil
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}

	expirationTimestamp := time.Unix(order.SignedOrder().ExpirationTimeSeconds.Int64(), 0)
	w.expirationWatcher.Remove(expirationTimestamp, order.Hash.Hex())
}

type orderDeleter interface {
	Delete(id []byte) error
}

// TODO(albrow): Do we need to do this inside a transaction?
func (w *Watcher) permanentlyDeleteOrder(order *types.OrderWithMetadata) error {
	if err := w.db.DeleteOrder(order.Hash); err != nil {
		return err
	}

	// After permanently deleting an order, we also remove it's assetData from the Decoder
	err := w.removeAssetDataAddressFromEventDecoder(order.MakerAssetData)
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
	case "StaticCall":
		var decodedAssetData zeroex.StaticCallAssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		// NOTE(jalextowle): The only staticcall that is currently supported
		// only relies on transaction gas price. This means that we do not need
		// to monitor any new contract addresses because the only supported
		// staticcall doesn't rely on any blockchain state.
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
	case "StaticCall":
		var decodedAssetData zeroex.StaticCallAssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		// NOTE(jalextowle): We aren't adding any contract addresses to the
		// orderwatcher for currently supported staticcalls, so we don't need
		// to remove anything here.
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

// TODO(albrow): Imlement CountOrders method
func (w *Watcher) increaseMaxExpirationTimeIfPossible() error {
	return errors.New("Not yet implemented")
	// if orderCount, err := w.meshtypes.OrderWithMetadatas.Count(); err != nil {
	// 	return err
	// } else if orderCount < w.maxOrders {
	// 	// We have enough space for new orders. Set the new max expiration time to the
	// 	// value of slow counter.
	// 	newMaxExpiration := w.maxExpirationCounter.Count()
	// 	if w.maxExpirationTime.Cmp(newMaxExpiration) != 0 {
	// 		logger.WithFields(logger.Fields{
	// 			"oldMaxExpirationTime": w.maxExpirationTime.String(),
	// 			"newMaxExpirationTime": fmt.Sprint(newMaxExpiration),
	// 		}).Debug("increasing max expiration time")
	// 		w.maxExpirationTime.Set(newMaxExpiration)
	// 		w.saveMaxExpirationTime(newMaxExpiration)
	// 	}
	// }

	// return nil
}

// saveMaxExpirationTime saves the new max expiration time in the database.
func (w *Watcher) saveMaxExpirationTime(maxExpirationTime *big.Int) {
	if err := w.db.UpdateMetadata(func(metadata db.Metadata) db.Metadata {
		metadata.MaxExpirationTime = maxExpirationTime
		return metadata
	}); err != nil {
		logger.WithError(err).Error("could not update max expiration time in database")
	}
}

func (w *Watcher) getBlockchainState(events []*blockwatch.Event) (*big.Int, time.Time) {
	var latestBlockNumber *big.Int
	var latestBlockTimestamp time.Time
	for _, event := range events {
		latestBlockNumber = event.BlockHeader.Number
		latestBlockTimestamp = event.BlockHeader.Timestamp
	}
	return latestBlockNumber, latestBlockTimestamp
}

// WaitForAtLeastOneBlockToBeProcessed waits until the OrderWatcher has processed it's
// first block
func (w *Watcher) WaitForAtLeastOneBlockToBeProcessed(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("Context cancelled")
	case <-w.atLeastOneBlockProcessed:
		return nil
	case <-time.After(60 * time.Second):
		return errors.New("timed out waiting for first block to be processed by Mesh node. Check your backing Ethereum RPC endpoint")
	}
}

type logWithType struct {
	Type string
	Log  ethtypes.Log
}
