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
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	// minCleanupInterval specified the minimum amount of time between orderbook
	// cleanups. These cleanups are meant to catch any stale orders that somehow
	// were not caught by the event watcher process.
	minCleanupInterval = 1 * time.Hour

	// maxDeleteInterval specifies the maximum amount of time between calls to
	// permanentlyDeleteStaleRemovedOrders.
	maxDeleteInterval = 5 * time.Minute

	// checkDatabaseUtilizationThresholdInterval specifies the amount of time
	// between calls to db.CountOrders. If the count is higher than a threshold,
	// then permanentlyDeleteStaleRemovedOrders will be called, so this is the
	// minimum interval between calls to permanentlyDeleteStaleRemovedOrders.
	checkDatabaseUtilizationThresholdInterval = 5 * time.Second

	// If, after a call to db.CountOrders, the database utilization exceeds
	// databaseUtilizationThreshold, then permanentlyDeleteStaleRemovedOrders
	// will be called.
	databaseUtilizationThreshold = 0.5

	// defaultLastUpdatedBuffer specifies how long it must have been since an
	// order was last updated in order to be re-validated by the cleanup worker.
	defaultLastUpdatedBuffer = 30 * time.Minute

	// permanentlyDeleteAfter specifies how long after an order is marked as
	// IsRemoved and not updated that it should be considered for permanent
	// deletion. Blocks get mined on avg. every 12 sec, so 5 minutes
	// corresponds to a block depth of ~25.
	permanentlyDeleteAfter = 5 * time.Minute

	// maxBlockEventsToHandle is the max number of block events we want to
	// process in a single call to `handleBlockEvents`
	maxBlockEventsToHandle = 500
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
	orderFeed                  event.Feed
	orderScope                 event.SubscriptionScope // Subscription scope tracking current live listeners
	contractAddressToSeenCount *contractAddressesSeenCounter
	orderValidator             *ordervalidator.OrderValidator
	wasStartedOnce             bool
	mu                         sync.Mutex
	maxOrders                  int
	handleBlockEventsMu        sync.RWMutex
	// atLeastOneBlockProcessed is closed to signal that the BlockWatcher has processed at least one
	// block. Validation of orders should block until this has completed
	atLeastOneBlockProcessed   chan struct{}
	atLeastOneBlockProcessedMu sync.Mutex
	didProcessABlock           bool
	// recentlyValidatedOrders is a list of orders that were added to the
	// orderwatcher after the most recent call to `handleBlockEvents`. Order
	// events may have been missed by the orderwatcher due to a rare edge case.
	// These orders must be tracked and checked for any missing order events
	// during the next call to `handleBlockEvents`.
	// For more information, refer to this issue:
	// https://github.com/0xProject/0x-mesh/issues/590
	recentlyValidatedOrdersMu sync.RWMutex
	recentlyValidatedOrders   []*types.OrderWithMetadata
}

type Config struct {
	DB                *db.DB
	BlockWatcher      *blockwatch.Watcher
	OrderValidator    *ordervalidator.OrderValidator
	ChainID           int
	ContractAddresses ethereum.ContractAddresses
	MaxOrders         int
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

	w := &Watcher{
		db:                         config.DB,
		blockWatcher:               config.BlockWatcher,
		contractAddressToSeenCount: NewContractAddressesSeenCounter(),
		orderValidator:             config.OrderValidator,
		eventDecoder:               decoder,
		assetDataDecoder:           assetDataDecoder,
		contractAddresses:          config.ContractAddresses,
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

	g, ctx := errgroup.WithContext(ctx)

	namedLoops := []struct {
		loop func(context.Context) error
		name string
	}{
		{w.mainLoop, "mainLoop"},
		{w.cleanupLoop, "cleanupLoop"},
		{w.removedCheckerLoop, "removedCheckerLoop"},
	}
	for _, namedLoop := range namedLoops {
		namedLoop := namedLoop // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			err := namedLoop.loop(ctx)
			if err != nil {
				logger.WithError(err).Errorf("error in orderwatcher %v", namedLoop.name)
			}
			return err
		})
	}

	return g.Wait()
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
			logger.WithField("error", err.Error()).Error("block subscription error encountered")
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
	for {
		select {
		case moreEvents := <-blockEventsChan:
			allEvents = append(allEvents, moreEvents...)
			if len(allEvents) >= max {
				return allEvents
			}
		default:
			return allEvents
		}
	}
}

func (w *Watcher) cleanupLoop(ctx context.Context) error {
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		// Wait minCleanupInterval before calling cleanup again. Since
		// we only start sleeping _after_ cleanup completes, we will never
		// have multiple calls to cleanup running in parallel
		case <-time.After(minCleanupInterval - time.Since(start)):
		}

		start = time.Now()
		if err := w.Cleanup(ctx, defaultLastUpdatedBuffer); err != nil {
			return err
		}
	}
}

func (w *Watcher) removedCheckerLoop(ctx context.Context) error {
	if err := w.permanentlyDeleteStaleRemovedOrders(); err != nil {
		return err
	}
	lastDeleted := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(checkDatabaseUtilizationThresholdInterval):
			count, err := w.db.CountOrders(nil)
			if err != nil {
				return err
			}
			databaseUtilization := float64(count) / float64(w.maxOrders)

			if time.Since(lastDeleted) > maxDeleteInterval || databaseUtilization > databaseUtilizationThreshold {
				if err := w.permanentlyDeleteStaleRemovedOrders(); err != nil {
					return err
				}
				lastDeleted = time.Now()
			}
		}
	}
}

// handleOrderExpirations takes care of generating expired and unexpired order events for orders that do not require re-validation.
// Since expiry is now done according to block timestamp, we can figure out which orders have expired/unexpired statically. We do not
// process orders that require re-validation, since the validation process will already emit the necessary events.
// latestBlockTimestamp is the latest block timestamp Mesh knows about
// ordersToRevalidate contains all the orders Mesh needs to re-validate given the events emitted by the blocks processed
func (w *Watcher) handleOrderExpirations(validationBlock *types.MiniHeader, ordersToRevalidate map[common.Hash]*types.OrderWithMetadata) ([]*zeroex.OrderEvent, map[common.Hash]struct{}, error) {
	orderEvents := []*zeroex.OrderEvent{}

	// Check for any orders that have now expired.
	expiredOrders, err := w.findOrdersToExpire(validationBlock.Timestamp)
	if err != nil {
		return orderEvents, nil, err
	}
	for _, order := range expiredOrders {
		// If we will re-validate this order, the revalidation process will discover that
		// it's expired, and an appropriate event will already be emitted
		if _, ok := ordersToRevalidate[order.Hash]; ok {
			continue
		}
		if order.KeepExpired {
			w.markOrderUnfillable(order, nil, validationBlock)
		} else {
			w.unwatchOrder(order, nil, validationBlock)
		}
		orderEvent := &zeroex.OrderEvent{
			Timestamp:                validationBlock.Timestamp,
			OrderHash:                order.Hash,
			SignedOrder:              order.SignedOrder(),
			FillableTakerAssetAmount: big.NewInt(0),
			EndState:                 zeroex.ESOrderExpired,
		}
		orderEvents = append(orderEvents, orderEvent)
	}

	// Check for any orders which have now unexpired.
	//
	// A block re-org may have happened resulting in the latest block timestamp
	// being lower than on the previous latest block. We need to "unexpire" any
	// orders that have now become valid again as a result.
	unexpiredOrders, err := w.findOrdersToUnexpire(validationBlock.Timestamp)
	if err != nil {
		return orderEvents, nil, err
	}
	for _, order := range unexpiredOrders {
		// If we will re-validate this order, the revalidation process will discover that
		// it's unexpired, and an appropriate event will already be emitted
		if _, ok := ordersToRevalidate[order.Hash]; ok {
			continue
		}
		w.rewatchOrder(order, order.FillableTakerAssetAmount, validationBlock)
		orderEvent := &zeroex.OrderEvent{
			Timestamp:                validationBlock.Timestamp,
			OrderHash:                order.Hash,
			SignedOrder:              order.SignedOrder(),
			FillableTakerAssetAmount: order.FillableTakerAssetAmount,
			EndState:                 zeroex.ESOrderUnexpired,
		}
		orderEvents = append(orderEvents, orderEvent)
	}

	possiblyUnexpiredOrders, err := w.findOrdersToPossiblyUnexpire(validationBlock.Timestamp)
	if err != nil {
		return orderEvents, nil, err
	}
	orderHashToPossiblyUnexpiredOrders := map[common.Hash]struct{}{}
	for _, order := range possiblyUnexpiredOrders {
		// If we will re-validate this order, the revalidation process will discover
		// whether or not it's unexpired, and an appropriate event will already be
		// emitted
		if _, ok := ordersToRevalidate[order.Hash]; ok {
			continue
		}
		ordersToRevalidate[order.Hash] = order
		orderHashToPossiblyUnexpiredOrders[order.Hash] = struct{}{}
	}

	return orderEvents, orderHashToPossiblyUnexpiredOrders, nil
}

// handleBlockEvents processes a set of block events into order events for a set of orders.
// handleBlockEvents MUST only be called after acquiring a lock to the `handleBlockEventsMu` mutex.
func (w *Watcher) handleBlockEvents(ctx context.Context, events []*blockwatch.Event) error {
	if len(events) == 0 {
		return nil
	}

	oldestBlockFromEvents, validationBlock := w.getExtremeBlocksFromEvents(events)
	orderHashToDBOrder := map[common.Hash]*types.OrderWithMetadata{}
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{}

	oldestBlockInDB, err := w.db.GetOldestMiniHeader()
	if err != nil {
		return err
	}

	w.recentlyValidatedOrdersMu.Lock()
	recentlyValidatedOrders := w.recentlyValidatedOrders
	w.recentlyValidatedOrders = []*types.OrderWithMetadata{}
	w.recentlyValidatedOrdersMu.Unlock()

	var oldestRevalidationBlockNumber *big.Int
	revalidationBlockToOrder := map[*big.Int][]*types.OrderWithMetadata{}
	for _, recentlyValidatedOrder := range recentlyValidatedOrders {
		previousValidationBlockNumber := recentlyValidatedOrder.LastValidatedBlockNumber
		// If the oldestBlock in the list of block events is greater then
		// the last validated block of the recently validated orders, we
		// may be missing block events for this order.
		if oldestRevalidationBlockNumber == nil || previousValidationBlockNumber.Cmp(oldestRevalidationBlockNumber) == -1 {
			oldestRevalidationBlockNumber = previousValidationBlockNumber
		}
		if oldestBlockFromEvents.Number.Cmp(previousValidationBlockNumber) == -1 {
			continue
		}
		// If the previous validation block of the order is a predecessor
		// of the oldest block in the blockwatcher, we must revalidate the
		// order because we may be missing relevant block events.
		if oldestBlockInDB.Number.Cmp(previousValidationBlockNumber) == 1 {
			orderHashToDBOrder[recentlyValidatedOrder.Hash] = recentlyValidatedOrder
			orderHashToEvents[recentlyValidatedOrder.Hash] = []*zeroex.ContractEvent{}
		}
		revalidationBlockToOrder[previousValidationBlockNumber] = append(
			revalidationBlockToOrder[previousValidationBlockNumber],
			recentlyValidatedOrder,
		)
	}

	// revalidationMiniHeaders is the set of blocks in between the latest
	// LastValidatedBlockNumber in recentlyValidatedOrders and the oldest
	// block in blockEvents. We need to check if any of these block events
	// affected the orders in recentlyValidatedOrders.
	revalidationMiniHeaders, err := w.db.FindMiniHeaders(&db.MiniHeaderQuery{
		Filters: []db.MiniHeaderFilter{
			{
				Field: db.MFNumber,
				Kind:  db.Less,
				Value: oldestBlockFromEvents.Number,
			},
			{
				Field: db.MFNumber,
				Kind:  db.Greater,
				Value: oldestRevalidationBlockNumber,
			},
		},
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Ascending,
			},
		},
	})
	if err != nil {
		return err
	}

	// Figure out which orders were potentially affected by the block events
	// and need to be re-validated.
	// For recentlyValidatedOrders, we check the list of revalidation miniheaders
	// for any block events that could change the validity of recently validated
	// orders.
	// For orders stored in the database, we check the list of new block events
	// to see if the validity of an order that could be changed.
	eventFilter := map[common.Hash]struct{}{}
	for _, header := range revalidationMiniHeaders {
		for _, order := range revalidationBlockToOrder[header.Number] {
			eventFilter[order.Hash] = struct{}{}
		}
		for _, log := range header.Logs {
			if err := w.findOrdersByEventWithFilter(log, eventFilter, orderHashToDBOrder, orderHashToEvents); err != nil {
				return err
			}
		}
	}

	for _, event := range events {
		for _, log := range event.BlockHeader.Logs {
			if err := w.findOrdersByEventWithFilter(log, nil, orderHashToDBOrder, orderHashToEvents); err != nil {
				return err
			}
		}
	}

	expirationOrderEvents, orderHashToPossiblyUnexpiredOrders, err := w.handleOrderExpirations(validationBlock, orderHashToDBOrder)
	if err != nil {
		return err
	}

	// This timeout of 1min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	postValidationOrderEvents, err := w.generateOrderEventsIfChanged(ctx, orderHashToDBOrder, orderHashToEvents, orderHashToPossiblyUnexpiredOrders, validationBlock)
	if err != nil {
		return err
	}

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

	return nil
}

// RevalidateOrdersForMissingEvents checks all of the orders in the database for
// any events in the miniheaders table that may have been missed. This should only
// be used on startup, as there is a different mechanism that serves this purpose
// during normal operation.
//
// NOTE(jalextowle): This function can miss block events if the blockwatcher was
// behind by more than db.MaxMiniHeaders when `handleBlockEvents` was last called.
// This is extremely unlikely, so we have decided not to implement more costly
// mechanisms to prevent from this possibility from occurring.
func (w *Watcher) RevalidateOrdersForMissingEvents(ctx context.Context) error {
	miniHeaders, err := w.db.FindMiniHeaders(nil)
	if err != nil {
		return err
	} else if len(miniHeaders) == 0 {
		// There is no need to check for missing events if there are no
		// miniheaders in the database.
		return nil
	}
	orderHashToDBOrder := map[common.Hash]*types.OrderWithMetadata{}
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{}
	for _, header := range miniHeaders {
		for _, log := range header.Logs {
			if err := w.findOrdersByEventWithLastValidatedBlockNumber(log, header.Number, orderHashToDBOrder, orderHashToEvents); err != nil {
				return err
			}
		}
	}
	latestMiniHeader, err := w.db.GetLatestMiniHeader()
	if err != nil {
		return err
	}
	orderEvents, err := w.generateOrderEventsIfChanged(ctx, orderHashToDBOrder, orderHashToEvents, map[common.Hash]struct{}{}, latestMiniHeader)
	if err != nil {
		return err
	}
	w.orderFeed.Send(orderEvents)
	return nil
}

// TODO(jalextowle): This could be made more efficient by only using the state from
// memory to check for orders that need to be revalidated. Currently, this will
// query for a number of orders in the database that do not need to be checked.
func (w *Watcher) findOrdersByEventWithFilter(
	log ethtypes.Log,
	filter map[common.Hash]struct{},
	orderHashToDBOrder map[common.Hash]*types.OrderWithMetadata,
	orderHashToEvents map[common.Hash][]*zeroex.ContractEvent,
) error {
	// TODO(jalextowle): This should be optimized by not querying the database
	// and instead just analyzing the list of recently validated orders.
	contractEvent, orders, err := w.findOrdersAffectedByContractEvents(log, db.OrderFilter{})
	if err != nil {
		return err
	}

	for _, order := range orders {
		found := true
		if filter != nil {
			_, found = filter[order.Hash]
		}
		if found {
			orderHashToDBOrder[order.Hash] = order
			if _, ok := orderHashToEvents[order.Hash]; !ok {
				orderHashToEvents[order.Hash] = []*zeroex.ContractEvent{contractEvent}
			} else {
				orderHashToEvents[order.Hash] = append(orderHashToEvents[order.Hash], contractEvent)
			}
		}
	}
	return nil
}

func (w *Watcher) findOrdersByEventWithLastValidatedBlockNumber(
	log ethtypes.Log,
	logBlockNumber *big.Int,
	orderHashToDBOrder map[common.Hash]*types.OrderWithMetadata,
	orderHashToEvents map[common.Hash][]*zeroex.ContractEvent,
) error {
	contractEvent, orders, err := w.findOrdersAffectedByContractEvents(log, db.OrderFilter{
		Field: db.OFLastValidatedBlockNumber,
		Kind:  db.GreaterOrEqual,
		Value: logBlockNumber,
	})
	if err != nil {
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
	return nil
}

// TODO(jalextowle): This could be optimized by taking functions as inputs that
// abstract away functions like `findOrderByTakerAddress`. This could eliminate
// unnecessary calls to the database and allow this function to be used in more
// general settings.
func (w *Watcher) findOrdersAffectedByContractEvents(log ethtypes.Log, filter db.OrderFilter) (*zeroex.ContractEvent, []*types.OrderWithMetadata, error) {
	eventType, err := w.eventDecoder.FindEventType(log)
	if err != nil {
		switch err := err.(type) {
		case decoder.UntrackedTokenError:
			return nil, nil, nil
		case decoder.UnsupportedEventError:
			logger.WithFields(logger.Fields{
				"topics":          err.Topics,
				"contractAddress": err.ContractAddress,
			}).Info("unsupported event found while trying to find its event type")
			return nil, nil, nil
		default:
			logger.WithFields(logger.Fields{
				"error": err.Error(),
			}).Error("unexpected event decoder error encountered")
			return nil, nil, err
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
				return nil, nil, nil
			}
			return nil, nil, err
		}
		contractEvent.Parameters = transferEvent
		fromOrders, err := w.findOrdersByTokenAddress(transferEvent.From, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, fromOrders...)
		toOrders, err := w.findOrdersByTokenAddress(transferEvent.To, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, toOrders...)

	case "ERC20ApprovalEvent":
		var approvalEvent decoder.ERC20ApprovalEvent
		err = w.eventDecoder.Decode(log, &approvalEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		// Ignores approvals set to anyone except the AssetProxy
		if approvalEvent.Spender != w.contractAddresses.ERC20Proxy {
			return nil, nil, nil
		}
		contractEvent.Parameters = approvalEvent
		orders, err = w.findOrdersByTokenAddress(approvalEvent.Owner, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}

	case "ERC721TransferEvent":
		var transferEvent decoder.ERC721TransferEvent
		err = w.eventDecoder.Decode(log, &transferEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		contractEvent.Parameters = transferEvent
		fromOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.From, log.Address, transferEvent.TokenId, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, fromOrders...)
		toOrders, err := w.findOrdersByTokenAddressAndTokenID(transferEvent.To, log.Address, transferEvent.TokenId, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, toOrders...)

	case "ERC721ApprovalEvent":
		var approvalEvent decoder.ERC721ApprovalEvent
		err = w.eventDecoder.Decode(log, &approvalEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		contractEvent.Parameters = approvalEvent
		orders, err = w.findOrdersByTokenAddressAndTokenID(approvalEvent.Owner, log.Address, approvalEvent.TokenId, filter)
		if err != nil {
			return nil, nil, err
		}

	case "ERC721ApprovalForAllEvent":
		var approvalForAllEvent decoder.ERC721ApprovalForAllEvent
		err = w.eventDecoder.Decode(log, &approvalForAllEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		// Ignores approvals set to anyone except the AssetProxy
		if approvalForAllEvent.Operator != w.contractAddresses.ERC721Proxy {
			return nil, nil, nil
		}
		contractEvent.Parameters = approvalForAllEvent
		orders, err = w.findOrdersByTokenAddress(approvalForAllEvent.Owner, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}

	case "ERC1155TransferSingleEvent":
		var transferEvent decoder.ERC1155TransferSingleEvent
		err = w.eventDecoder.Decode(log, &transferEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		// HACK(fabio): Currently we simply revalidate all orders involving assets in this
		// ERC1155 contract from this particular maker. We could however revalidate fewer orders
		// by also taking into account the `ID` of the assets affected. We punt on this for now
		// in order to support Augur's use-case of a dummy ERC1155 contract. In their case, we
		// need to revalidate all maker orders within the single ERC1155 contract and cannot optimize
		// further. In the future, we might want to special-case this broader approach for the Augur
		// contract address specifically.
		contractEvent.Parameters = transferEvent
		fromOrders, err := w.findOrdersByTokenAddress(transferEvent.From, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, fromOrders...)
		toOrders, err := w.findOrdersByTokenAddress(transferEvent.To, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, toOrders...)

	case "ERC1155TransferBatchEvent":
		var transferEvent decoder.ERC1155TransferBatchEvent
		err = w.eventDecoder.Decode(log, &transferEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		contractEvent.Parameters = transferEvent
		fromOrders, err := w.findOrdersByTokenAddress(transferEvent.From, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, fromOrders...)
		toOrders, err := w.findOrdersByTokenAddress(transferEvent.To, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, toOrders...)

	case "ERC1155ApprovalForAllEvent":
		var approvalForAllEvent decoder.ERC1155ApprovalForAllEvent
		err = w.eventDecoder.Decode(log, &approvalForAllEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		// Ignores approvals set to anyone except the AssetProxy
		if approvalForAllEvent.Operator != w.contractAddresses.ERC1155Proxy {
			return nil, nil, nil
		}
		contractEvent.Parameters = approvalForAllEvent
		orders, err = w.findOrdersByTokenAddress(approvalForAllEvent.Owner, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}

	case "WethWithdrawalEvent":
		var withdrawalEvent decoder.WethWithdrawalEvent
		err = w.eventDecoder.Decode(log, &withdrawalEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		contractEvent.Parameters = withdrawalEvent
		orders, err = w.findOrdersByTokenAddress(withdrawalEvent.Owner, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}

	case "WethDepositEvent":
		var depositEvent decoder.WethDepositEvent
		err = w.eventDecoder.Decode(log, &depositEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
		}
		contractEvent.Parameters = depositEvent
		orders, err = w.findOrdersByTokenAddress(depositEvent.Owner, log.Address, filter)
		if err != nil {
			return nil, nil, err
		}

	case "ExchangeFillEvent":
		var exchangeFillEvent decoder.ExchangeFillEvent
		err = w.eventDecoder.Decode(log, &exchangeFillEvent)
		if err != nil {
			if isNonCritical := w.checkDecodeErr(err, eventType); isNonCritical {
				return nil, nil, nil
			}
			return nil, nil, err
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
				return nil, nil, nil
			}
			return nil, nil, err
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
				return nil, nil, nil
			}
			return nil, nil, err
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
			return nil, nil, err
		}
		orders = append(orders, cancelledOrders...)

	default:
		logger.WithFields(logger.Fields{
			"eventType": eventType,
			"log":       log,
		}).Error("unknown eventType encountered")
		return nil, nil, err
	}

	return contractEvent, orders, nil
}

func (w *Watcher) getLatestBlock() (*types.MiniHeader, error) {
	latestBlock, err := w.db.GetLatestMiniHeader()
	if err != nil {
		if err == db.ErrNotFound {
			return nil, errNoBlocksStored
		}
		return nil, err
	}
	return latestBlock, nil
}

// Cleanup re-validates all orders in DB which haven't been re-validated in
// `lastUpdatedBuffer` time to make sure all orders are still up-to-date
func (w *Watcher) Cleanup(ctx context.Context, lastUpdatedBuffer time.Duration) error {
	// Pause block event processing until we finished cleaning up at current block height
	w.handleBlockEventsMu.RLock()
	defer w.handleBlockEventsMu.RUnlock()

	lastUpdatedCutOff := time.Now().Add(-lastUpdatedBuffer)
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
	orderEvents, err := w.generateOrderEventsIfChanged(ctx, orderHashToDBOrder, orderHashToEvents, map[common.Hash]struct{}{}, latestBlock)
	if err != nil {
		return err
	}

	if len(orderEvents) > 0 {
		w.orderFeed.Send(orderEvents)
	}

	return nil
}

func (w *Watcher) permanentlyDeleteStaleRemovedOrders() error {
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
func (w *Watcher) add(orderInfos []*ordervalidator.AcceptedOrderInfo, validationBlock *types.MiniHeader, pinned bool, opts *types.AddOrdersOpts) ([]*zeroex.OrderEvent, error) {
	now := time.Now().UTC()
	orderEvents := []*zeroex.OrderEvent{}
	dbOrders := []*types.OrderWithMetadata{}

	for _, orderInfo := range orderInfos {
		dbOrder, err := w.orderInfoToOrderWithMetadata(orderInfo, pinned, now, validationBlock, opts)
		if err != nil {
			return nil, err
		}
		dbOrders = append(dbOrders, dbOrder)

		// We create an ADDED event for all orders in orderInfos.
		// Some orders might not actually be added, as a workaround we
		// will also emit a STOPPED_WATCHING event in some cases (see
		// below)
		addedEvent := &zeroex.OrderEvent{
			Timestamp:                now,
			OrderHash:                orderInfo.OrderHash,
			SignedOrder:              orderInfo.SignedOrder,
			FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
			EndState:                 zeroex.ESOrderAdded,
		}
		orderEvents = append(orderEvents, addedEvent)
	}

	addedMap := map[common.Hash]*types.OrderWithMetadata{}
	alreadyStored, addedOrders, removedOrders, err := w.db.AddOrders(dbOrders)
	alreadyStoredSet := map[common.Hash]struct{}{}
	if err != nil {
		return nil, err
	}
	for _, hash := range alreadyStored {
		// Add each hash to a set of already stored hashes. This allows for faster
		// lookups later on.
		alreadyStoredSet[hash] = struct{}{}
	}
	for _, order := range addedOrders {
		err = w.setupInMemoryOrderState(order)
		if err != nil {
			return orderEvents, err
		}
		addedMap[order.Hash] = order
		w.recentlyValidatedOrdersMu.Lock()
		w.recentlyValidatedOrders = append(w.recentlyValidatedOrders, order)
		w.recentlyValidatedOrdersMu.Unlock()
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
		err = w.removeAssetDataAddressFromEventDecoder(order.OrderV3.MakerAssetData)
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

	// HACK(albrow): We need to handle orders in the orderInfos argument that
	// were never added due to the max expiration time effectively changing
	// within the database transaction above. In other words, new orders that
	// _were_ added can change the effective max expiration time, meaning some
	// orders in orderInfos were actually not added. This should not happen
	// often. For now, we respond by emitting an ADDED event (above) immediately
	// followed by a STOPPED_WATCHING event. If this order was submitted via
	// GraphQL, the GraphQL client will see a response that indicates the order was
	// successfully added, and then it will look like we immediately stopped
	// watching it. This is not too far off from what really happened but is
	// slightly inefficient.
	//
	// We can detect this by looking for orders that we should have added but
	// are not included in either wasAdded map or the alreadyStored set.
	//
	// TODO(albrow): In the future, we should add an additional return value and
	// then react to that differently depending on whether the order was
	// received via GraphQL or from a peer. In the former case, we should return an
	// GraphQL error response indicating that the order was not in fact added. In
	// the latter case, we should not emit any order events but might potentially
	// want to adjust the peer's score.
	for _, orderToAdd := range orderInfos {
		_, wasAdded := addedMap[orderToAdd.OrderHash]
		_, alreadyStored := alreadyStoredSet[orderToAdd.OrderHash]
		if !wasAdded && !alreadyStored {
			stoppedWatchingEvent := &zeroex.OrderEvent{
				Timestamp:                now,
				OrderHash:                orderToAdd.OrderHash,
				SignedOrder:              orderToAdd.SignedOrder,
				FillableTakerAssetAmount: orderToAdd.FillableTakerAssetAmount,
				EndState:                 zeroex.ESStoppedWatching,
			}
			orderEvents = append(orderEvents, stoppedWatchingEvent)
		}
	}

	if len(removedOrders) > 0 {
		newMaxExpirationTime, err := w.db.GetCurrentMaxExpirationTime()
		if err != nil {
			return nil, err
		}
		logger.WithFields(logger.Fields{
			"ordersRemoved":        len(removedOrders),
			"newMaxExpirationTime": newMaxExpirationTime.String(),
		}).Debug("removed orders due to exceeding max expiration time")
	}

	return orderEvents, nil
}

func (w *Watcher) orderInfoToOrderWithMetadata(orderInfo *ordervalidator.AcceptedOrderInfo, pinned bool, now time.Time, validationBlock *types.MiniHeader, opts *types.AddOrdersOpts) (*types.OrderWithMetadata, error) {
	// V4 Orders
	if orderInfo.SignedOrder == nil && orderInfo.SignedOrderV4 != nil {
		return &types.OrderWithMetadata{
			Hash:                     orderInfo.OrderHash,
			OrderV4:                  &orderInfo.SignedOrderV4.OrderV4,
			SignatureV4:              orderInfo.SignedOrderV4.Signature,
			IsRemoved:                false,
			IsUnfillable:             orderInfo.FillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0,
			IsPinned:                 pinned,
			IsExpired:                big.NewInt(validationBlock.Timestamp.Unix()).Cmp(orderInfo.SignedOrderV4.OrderV4.Expiry) >= 0,
			LastUpdated:              now,
			FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
			LastValidatedBlockNumber: validationBlock.Number,
			LastValidatedBlockHash:   validationBlock.Hash,
			KeepCancelled:            opts.KeepCancelled,
			KeepExpired:              opts.KeepExpired,
			KeepFullyFilled:          opts.KeepFullyFilled,
			KeepUnfunded:             opts.KeepUnfunded,
		}, nil
	}
	if orderInfo.SignedOrder == nil {
		return nil, errors.New("OrderInfo contains neither V3 nor V4 order")
	}

	// V3 Orders
	parsedMakerAssetData, err := db.ParseContractAddressesAndTokenIdsFromAssetData(w.assetDataDecoder, orderInfo.SignedOrder.MakerAssetData, w.contractAddresses)
	if err != nil {
		return nil, err
	}
	parsedMakerFeeAssetData, err := db.ParseContractAddressesAndTokenIdsFromAssetData(w.assetDataDecoder, orderInfo.SignedOrder.MakerFeeAssetData, w.contractAddresses)
	if err != nil {
		return nil, err
	}
	return &types.OrderWithMetadata{
		Hash: orderInfo.OrderHash,
		OrderV3: &zeroex.Order{
			ChainID:               orderInfo.SignedOrder.ChainID,
			ExchangeAddress:       orderInfo.SignedOrder.ExchangeAddress,
			MakerAddress:          orderInfo.SignedOrder.MakerAddress,
			MakerAssetData:        orderInfo.SignedOrder.MakerAssetData,
			MakerFeeAssetData:     orderInfo.SignedOrder.MakerFeeAssetData,
			MakerAssetAmount:      orderInfo.SignedOrder.MakerAssetAmount,
			MakerFee:              orderInfo.SignedOrder.MakerFee,
			TakerAddress:          orderInfo.SignedOrder.TakerAddress,
			TakerAssetData:        orderInfo.SignedOrder.TakerAssetData,
			TakerFeeAssetData:     orderInfo.SignedOrder.TakerFeeAssetData,
			TakerAssetAmount:      orderInfo.SignedOrder.TakerAssetAmount,
			TakerFee:              orderInfo.SignedOrder.TakerFee,
			SenderAddress:         orderInfo.SignedOrder.SenderAddress,
			FeeRecipientAddress:   orderInfo.SignedOrder.FeeRecipientAddress,
			ExpirationTimeSeconds: orderInfo.SignedOrder.ExpirationTimeSeconds,
			Salt:                  orderInfo.SignedOrder.Salt,
		},
		Signature:                orderInfo.SignedOrder.Signature,
		IsRemoved:                false,
		IsUnfillable:             orderInfo.FillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0,
		IsPinned:                 pinned,
		IsExpired:                big.NewInt(validationBlock.Timestamp.Unix()).Cmp(orderInfo.SignedOrder.ExpirationTimeSeconds) >= 0,
		LastUpdated:              now,
		ParsedMakerAssetData:     parsedMakerAssetData,
		ParsedMakerFeeAssetData:  parsedMakerFeeAssetData,
		FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
		LastValidatedBlockNumber: validationBlock.Number,
		LastValidatedBlockHash:   validationBlock.Hash,
		KeepCancelled:            opts.KeepCancelled,
		KeepExpired:              opts.KeepExpired,
		KeepFullyFilled:          opts.KeepFullyFilled,
		KeepUnfunded:             opts.KeepUnfunded,
	}, nil
}

// TODO(albrow): All in-memory state can be removed.
func (w *Watcher) setupInMemoryOrderState(order *types.OrderWithMetadata) error {
	if order.OrderV3 != nil {
		w.eventDecoder.AddKnownExchange(order.OrderV3.ExchangeAddress)
		// Add MakerAssetData and MakerFeeAssetData to EventDecoder
		err := w.addAssetDataAddressToEventDecoder(order.OrderV3.MakerAssetData)
		if err != nil {
			return err
		}
		if order.OrderV3.MakerFee.Cmp(big.NewInt(0)) == 1 {
			err = w.addAssetDataAddressToEventDecoder(order.OrderV3.MakerFeeAssetData)
			if err != nil {
				return err
			}
		}
	}
	if order.OrderV4 != nil {
		w.eventDecoder.AddKnownExchange(order.OrderV4.ExchangeAddress)
		// Add MakerToken to EventDecoder
		w.eventDecoder.AddKnownERC20(order.OrderV4.MakerToken)
		w.contractAddressToSeenCount.Inc(order.OrderV4.MakerToken)
	}
	return nil
}

// Subscribe allows one to subscribe to the order events emitted by the OrderWatcher.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*zeroex.OrderEvent) event.Subscription {
	return w.orderScope.Track(w.orderFeed.Subscribe(sink))
}

func (w *Watcher) findOrder(orderHash common.Hash) *types.OrderWithMetadata {
	// V3
	order, err := w.db.GetOrder(orderHash)
	if err == nil {
		return order
		}
	if err != db.ErrNotFound {
		logger.WithFields(logger.Fields{
			"error":     err.Error(),
			"orderHash": orderHash,
		}).Warning("Unexpected error from db.GetOrder")
		return nil
	}

	// V4
	orderV4, err := w.db.GetOrderV4(orderHash)
	if err == nil {
		return orderV4
}
	if err == db.ErrNotFound {
		return nil
	}
	logger.WithFields(logger.Fields{
		"error":     err.Error(),
		"orderHash": orderHash,
	}).Warning("Unexpected error from db.GetOrderV4")
	return nil
}

// findOrdersByTokenAddressAndTokenID finds and returns all orders that have
// either a makerAsset or a makerFeeAsset matching the given tokenAddress and
// tokenID.
func (w *Watcher) findOrdersByTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int, filter db.OrderFilter) ([]*types.OrderWithMetadata, error) {
	filters := []db.OrderFilter{
		{
			Field: db.OFMakerAddress,
			Kind:  db.Equal,
			Value: makerAddress,
		},
		db.MakerAssetIncludesTokenAddressAndTokenID(tokenAddress, tokenID),
	}
	if filter.Kind != "" {
		filters = append(filters, filter)
	}
	ordersWithAffectedMakerAsset, err := w.db.FindOrders(&db.OrderQuery{Filters: filters})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}
	filters = []db.OrderFilter{
		{
			Field: db.OFMakerAddress,
			Kind:  db.Equal,
			Value: makerAddress,
		},
		db.MakerFeeAssetIncludesTokenAddressAndTokenID(tokenAddress, tokenID),
	}
	if filter.Kind != "" {
		filters = append(filters, filter)
	}
	ordersWithAffectedMakerFeeAsset, err := w.db.FindOrders(&db.OrderQuery{Filters: filters})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}

	// V4 does not support NFTs so has no orders with tokenIds

	return append(ordersWithAffectedMakerAsset, ordersWithAffectedMakerFeeAsset...), nil
}

// findOrdersByTokenAddress finds and returns all orders that have
// either a makerAsset or a makerFeeAsset matching the given tokenAddress and
// any tokenID (including null).
func (w *Watcher) findOrdersByTokenAddress(makerAddress, tokenAddress common.Address, filter db.OrderFilter) ([]*types.OrderWithMetadata, error) {
	filters := []db.OrderFilter{
		{
			Field: db.OFMakerAddress,
			Kind:  db.Equal,
			Value: makerAddress,
		},
		db.MakerAssetIncludesTokenAddress(tokenAddress),
	}
	if filter.Kind != "" {
		filters = append(filters, filter)
	}
	ordersWithAffectedMakerAsset, err := w.db.FindOrders(&db.OrderQuery{Filters: filters})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}
	filters = []db.OrderFilter{
		{
			Field: db.OFMakerAddress,
			Kind:  db.Equal,
			Value: makerAddress,
		},
		db.MakerFeeAssetIncludesTokenAddress(tokenAddress),
	}
	if filter.Kind != "" {
		filters = append(filters, filter)
	}
	ordersWithAffectedMakerFeeAsset, err := w.db.FindOrders(&db.OrderQuery{Filters: filters})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}

	// V4 Orders
	ordersV4, err := w.db.FindOrdersV4(&db.OrderQueryV4{
		Filters: []db.OrderFilterV4{
			{
				Field: db.OV4FMaker,
				Kind:  db.Equal,
				Value: makerAddress,
			},
			{
				Field: db.OV4FMakerToken,
				Kind:  db.Equal,
				Value: tokenAddress,
			},
		},
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
		}).Error("unexpected query error encountered")
		return nil, err
	}

	return append(append(ordersWithAffectedMakerAsset, ordersWithAffectedMakerFeeAsset...), ordersV4...), nil
}

// findOrdersToExpire returns all orders with an expiration time less than or equal to the latest
// block timestamp that have not already been removed.
func (w *Watcher) findOrdersToExpire(latestBlockTimestamp time.Time) ([]*types.OrderWithMetadata, error) {
	ordersV3, err := w.db.FindOrders(&db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFExpirationTimeSeconds,
				Kind:  db.LessOrEqual,
				Value: big.NewInt(latestBlockTimestamp.Unix()),
			},
			{
				Field: db.OFIsUnfillable,
				Kind:  db.Equal,
				Value: false,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ordersV4, err := w.db.FindOrdersV4(&db.OrderQueryV4{
		Filters: []db.OrderFilterV4{
			{
				Field: db.OV4FExpiry,
				Kind:  db.LessOrEqual,
				Value: big.NewInt(latestBlockTimestamp.Unix()),
			},
			{
				Field: db.OV4FIsUnfillable,
				Kind:  db.Equal,
				Value: false,
			},
		},
	})
	if err != nil {
		return nil, err
}
	return append(ordersV3, ordersV4...), nil
}

// findOrdersToUnexpire returns all orders that:
//
//     1. have an expiration time greater than the latest block timestamp
//     2. were previously unfillable
//     3. have a non-zero FillableTakerAssetAmount
//
func (w *Watcher) findOrdersToUnexpire(latestBlockTimestamp time.Time) ([]*types.OrderWithMetadata, error) {
	ordersV3, err := w.db.FindOrders(&db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFExpirationTimeSeconds,
				Kind:  db.Greater,
				Value: big.NewInt(latestBlockTimestamp.Unix()),
			},
			{
				Field: db.OFIsUnfillable,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OFFillableTakerAssetAmount,
				Kind:  db.NotEqual,
				Value: big.NewInt(0),
			},
		},
	})
	if err != nil {
		return nil, err
}
	ordersV4, err := w.db.FindOrdersV4(&db.OrderQueryV4{
		Filters: []db.OrderFilterV4{
			{
				Field: db.OV4FExpiry,
				Kind:  db.Greater,
				Value: big.NewInt(latestBlockTimestamp.Unix()),
			},
			{
				Field: db.OV4FIsUnfillable,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OV4FFillableTakerAssetAmount,
				Kind:  db.NotEqual,
				Value: big.NewInt(0),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return append(ordersV3, ordersV4...), nil
}

// findOrdersToPossiblyUnexpire returns all orders that:
//
//     1. have an expiration time greater than the latest block timestamp
//     2. were previously unfillable
//     3. were previously expired
//     4. have a zero FillableTakerAssetAmount
//
func (w *Watcher) findOrdersToPossiblyUnexpire(latestBlockTimestamp time.Time) ([]*types.OrderWithMetadata, error) {
	ordersV3, err := w.db.FindOrders(&db.OrderQuery{
		Filters: []db.OrderFilter{
			{
				Field: db.OFExpirationTimeSeconds,
				Kind:  db.Greater,
				Value: big.NewInt(latestBlockTimestamp.Unix()),
			},
			{
				Field: db.OFIsUnfillable,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OFIsExpired,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OFFillableTakerAssetAmount,
				Kind:  db.Equal,
				Value: big.NewInt(0),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ordersV4, err := w.db.FindOrdersV4(&db.OrderQueryV4{
		Filters: []db.OrderFilterV4{
			{
				Field: db.OV4FExpiry,
				Kind:  db.Greater,
				Value: big.NewInt(latestBlockTimestamp.Unix()),
			},
			{
				Field: db.OV4FIsUnfillable,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OV4FIsExpired,
				Kind:  db.Equal,
				Value: true,
			},
			{
				Field: db.OV4FFillableTakerAssetAmount,
				Kind:  db.Equal,
				Value: big.NewInt(0),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return append(ordersV3, ordersV4...), nil
}

func (w *Watcher) convertValidationResultsIntoOrderEvents(
	validationResults *ordervalidator.ValidationResults,
	orderHashToDBOrder map[common.Hash]*types.OrderWithMetadata,
	orderHashToEvents map[common.Hash][]*zeroex.ContractEvent,
	orderHashToPossiblyUnexpiredOrder map[common.Hash]struct{},
	validationBlock *types.MiniHeader,
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
			// A previous event caused this order to be removed from DB because its
			// fillableAmount became 0, but it has now been revived (e.g., block re-org
			// causes order fill txn to get reverted). We need to re-add order and emit an event.
			w.rewatchOrder(order, newFillableAmount, validationBlock)
			orderEvent := &zeroex.OrderEvent{
				Timestamp:                validationBlock.Timestamp,
				OrderHash:                acceptedOrderInfo.OrderHash,
				SignedOrder:              order.SignedOrder(),
				FillableTakerAssetAmount: newFillableAmount,
				EndState:                 zeroex.ESOrderAdded,
				ContractEvents:           orderHashToEvents[order.Hash],
			}
			orderEvents = append(orderEvents, orderEvent)
		} else {
			// The order expiration time is valid if it is greater than the latest block timestamp
			// of the validation block.
			validationBlockTimestampSeconds := big.NewInt(validationBlock.Timestamp.Unix())
			expirationTimeIsValid := order.OrderV3.ExpirationTimeSeconds.Cmp(validationBlockTimestampSeconds) == 1
			isOrderUnexpired := order.IsExpired && order.IsUnfillable && expirationTimeIsValid

			// We can tell that an order was previously expired if it was marked as removed with a
			// non-zero fillable amount. There is no other explanation for this database state. The
			// order is considered "unexpired" if it was previously expired but now has a valid
			// expiration time based on the latest block timestamp.
			if isOrderUnexpired {
				w.rewatchOrder(order, newFillableAmount, validationBlock)
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                validationBlock.Timestamp,
					OrderHash:                order.Hash,
					SignedOrder:              order.SignedOrder(),
					FillableTakerAssetAmount: order.FillableTakerAssetAmount,
					EndState:                 zeroex.ESOrderUnexpired,
				}
				orderEvents = append(orderEvents, orderEvent)
			} else {
				w.updateOrderFillableTakerAssetAmountAndBlockInfo(order, newFillableAmount, validationBlock)
			}

			if oldFillableAmount.Cmp(newFillableAmount) == 0 {
				// No important state-change happened. Note that either rewatchOrder or
				// updateOrderFillableTakerAssetAmountAndBlockInfo in the unexpiration logic has already
				// updated lastValidatedBlock.
				continue
			} else {
				// Either the fillable amount has increased, e.g. a fill transaction reverted
				// because of a block reorg, or it has decreased because of a partial or complete
				// fill.
				endState := zeroex.ESOrderFillabilityIncreased
				if oldAmountIsMoreThenNewAmount {
					endState = zeroex.ESOrderFilled
				}
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                validationBlock.Timestamp,
					OrderHash:                acceptedOrderInfo.OrderHash,
					SignedOrder:              order.SignedOrder(),
					EndState:                 endState,
					FillableTakerAssetAmount: newFillableAmount,
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
			// NOTE(jalextowle): It's theoretically possible that an order could become
			// unexpired for a long period of time. With this in mind, it's important
			// that we set IsExpired to false regardless of why the check failed as
			// any future attempts to unexpire the order will be unsuccessful unless other
			// state changes occur. These state changes would be addressed by the core
			// orderwatcher logic so we can safely avoid re-checking for unexpiry.
			if order.IsUnfillable {
				// If the order is already marked as unfillable, no updates are needed
				// other than updating the expiration state
				if _, ok := orderHashToPossiblyUnexpiredOrder[order.Hash]; ok {
					w.updateOrderExpirationState(order, validationBlock)
				}
			} else {
				endState, ok := ordervalidator.ConvertRejectOrderCodeToOrderEventEndState(rejectedOrderInfo.Status)
				if !ok {
					err := fmt.Errorf("no OrderEventEndState corresponding to RejectedOrderStatus: %q", rejectedOrderInfo.Status)
					logger.WithError(err).WithField("rejectedOrderStatus", rejectedOrderInfo.Status).Error("no OrderEventEndState corresponding to RejectedOrderStatus")
					return nil, err
				}
				if (endState == zeroex.ESOrderCancelled && order.KeepCancelled) ||
					(endState == zeroex.ESOrderExpired && order.KeepExpired) ||
					(endState == zeroex.ESOrderFullyFilled && order.KeepFullyFilled) ||
					(endState == zeroex.ESOrderBecameUnfunded && order.KeepUnfunded) {
					w.markOrderUnfillable(order, big.NewInt(0), validationBlock)
				} else {
					w.unwatchOrder(order, big.NewInt(0), validationBlock)
				}
				orderEvent := &zeroex.OrderEvent{
					Timestamp:                validationBlock.Timestamp,
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
	orderHashToPossiblyUnexpiredOrder map[common.Hash]struct{},
	validationBlock *types.MiniHeader,
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
	validationResults := w.orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, validationBlock)

	return w.convertValidationResultsIntoOrderEvents(
		validationResults, orderHashToDBOrder, orderHashToEvents, orderHashToPossiblyUnexpiredOrder, validationBlock,
	)
}

// ValidateAndStoreValidOrders applies general 0x validation and Mesh-specific validation to
// the given orders and if they are valid, adds them to the OrderWatcher
func (w *Watcher) ValidateAndStoreValidOrders(ctx context.Context, orders []*zeroex.SignedOrder, chainID int, pinned bool, opts *types.AddOrdersOpts) (*ordervalidator.ValidationResults, error) {
	if len(orders) == 0 {
		return &ordervalidator.ValidationResults{}, nil
	}
	results, validMeshOrders, err := w.meshSpecificOrderValidation(orders, chainID, pinned)
	if err != nil {
		return nil, err
	}

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

	if opts.KeepCancelled || opts.KeepExpired || opts.KeepFullyFilled || opts.KeepUnfunded {
		for _, rejectedOrderInfo := range zeroexResults.Rejected {
			// NOTE(jalextowle): We can use the rejectedOrderInfo.Status
			// field to see whether or not the order is new or not. If
			// the order has already been stored, the rejectedOrderInfo.Status
			// field will be ordervalidator.ROOrderAlreadyStoredAndUnfillable.
			// If the rejection reason involves on-chain validation, then the
			// order is new.
			if (opts.KeepCancelled && rejectedOrderInfo.Status.Code == ordervalidator.ROCancelled.Code) ||
				(opts.KeepExpired && rejectedOrderInfo.Status.Code == ordervalidator.ROExpired.Code) ||
				(opts.KeepFullyFilled && rejectedOrderInfo.Status.Code == ordervalidator.ROFullyFilled.Code) ||
				(opts.KeepUnfunded && rejectedOrderInfo.Status.Code == ordervalidator.ROUnfunded.Code) {
				newOrderInfos = append(newOrderInfos, &ordervalidator.AcceptedOrderInfo{
					OrderHash:   rejectedOrderInfo.OrderHash,
					SignedOrder: rejectedOrderInfo.SignedOrder,
					// TODO(jalextowle): Verify that this is consistent with the OrderWatcher
					FillableTakerAssetAmount: big.NewInt(0),
					IsNew:                    true,
				})
			}
		}
	}

	// Add the order to the OrderWatcher. This also saves the order in the
	// database.
	allOrderEvents := []*zeroex.OrderEvent{}
	orderEvents, err := w.add(newOrderInfos, validationBlock, pinned, opts)
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
	latestBlock, err := w.getLatestBlock()
	if err != nil {
		return nil, nil, err
	}
	// This timeout of 1min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	areNewOrders := true
	zeroexResults := w.orderValidator.BatchValidate(ctx, orders, areNewOrders, latestBlock)
	return latestBlock, zeroexResults, nil
}

func (w *Watcher) meshSpecificOrderValidation(orders []*zeroex.SignedOrder, chainID int, pinned bool) (*ordervalidator.ValidationResults, []*zeroex.SignedOrder, error) {
	results := &ordervalidator.ValidationResults{}
	validMeshOrders := []*zeroex.SignedOrder{}

	// Calculate max expiration time based on number of orders stored.
	// This value is *exclusive*. Any incoming orders with an expiration time
	// greater or equal to this will be rejected.
	//
	// Note(albrow): Technically speaking this is sub-optimal. We are assuming
	// that we need to have space in the database for the entire slice of orders,
	// but some of them could be invalid and therefore not actually get stored.
	// However, the optimal implementation would be less efficient and could
	// result in sending more ETH RPC requests than necessary. The edge case
	// where potentially valid orders are rejected should be rare in practice, and
	// would affect at most len(orders)/2 orders.
	maxExpirationTime := constants.UnlimitedExpirationTime
	if !pinned {
		orderCount, err := w.db.CountOrders(nil)
		if err != nil {
			return nil, nil, err
		}
		if orderCount+len(orders) > w.maxOrders {
			storedMaxExpirationTime, err := w.db.GetCurrentMaxExpirationTime()
			if err != nil {
				return nil, nil, err
			}
			maxExpirationTime = storedMaxExpirationTime
		}
	}

	validOrderHashes := []common.Hash{}
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
		if !pinned && order.ExpirationTimeSeconds.Cmp(maxExpirationTime) != -1 {
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
		// sender addresses over time.
		if order.SenderAddress != constants.NullAddress {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROSenderAddressNotAllowed,
			})
			continue
		}
		// NOTE(jalextowle): Orders with a taker address are only accessible to
		// one taker, which complicates more sophisticated pruning technology.
		// With this in mind, we only allow whitelisted taker addresses to be
		// propogated throughout the network. This whitelist should only include
		// addresses that correspond to contracts allow anyone to fill these
		// orders.
		// TODO(jalextowle): If any other addresses are whitelisted, create
		// a isTakerAddressWhitelisted function.
		if order.TakerAddress != constants.NullAddress && order.TakerAddress != w.contractAddresses.ExchangeProxyFlashWallet {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROTakerAddressNotAllowed,
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

		validOrderHashes = append(validOrderHashes, orderHash)
		validMeshOrders = append(validMeshOrders, order)
	}

	newValidOrders := []*zeroex.SignedOrder{}
	storedOrderStatuses, err := w.db.GetOrderStatuses(validOrderHashes)
	if err != nil {
		logger.WithField("error", err).Error("could not get stored order statuses")
		return nil, nil, err
	}
	if len(storedOrderStatuses) != len(validOrderHashes) {
		return nil, nil, errors.New("could not get stored order statuses for all orders")
	}
	for i, order := range validMeshOrders {
		orderStatus := storedOrderStatuses[i]
		orderHash := validOrderHashes[i]
		if !orderStatus.IsStored {
			// If not stored, add the order to a set of new orders.
			newValidOrders = append(newValidOrders, order)
		} else if orderStatus.IsMarkedRemoved || orderStatus.IsMarkedUnfillable {
			// If stored but marked as removed or unfillable, reject the order.
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROOrderAlreadyStoredAndUnfillable,
			})
		} else {
			// If stored but not marked as removed or unfillable, accept the order without re-validation
			results.Accepted = append(results.Accepted, &ordervalidator.AcceptedOrderInfo{
				OrderHash:                orderHash,
				SignedOrder:              order,
				FillableTakerAssetAmount: orderStatus.FillableTakerAssetAmount,
				IsNew:                    false,
			})
		}
	}

	return results, newValidOrders, nil
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

func validateOrderSizeV4(order *zeroex.SignedOrderV4) error {
	encoded, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if len(encoded) > constants.MaxOrderSizeInBytes {
		return constants.ErrMaxOrderSize
	}
	return nil
}

// TODO(albrow): Add tests for LastValidatedBlockNumber and LastValidatedBlockHash for
// this and other similar functions.
func (w *Watcher) updateOrderFillableTakerAssetAmountAndBlockInfo(order *types.OrderWithMetadata, newFillableTakerAssetAmount *big.Int, validationBlock *types.MiniHeader) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.LastUpdated = time.Now().UTC()
		orderToUpdate.FillableTakerAssetAmount = newFillableTakerAssetAmount
		orderToUpdate.LastValidatedBlockNumber = validationBlock.Number
		orderToUpdate.LastValidatedBlockHash = validationBlock.Hash
		return orderToUpdate, nil
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}
}

func (w *Watcher) rewatchOrder(order *types.OrderWithMetadata, newFillableTakerAssetAmount *big.Int, validationBlock *types.MiniHeader) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.IsRemoved = false
		orderToUpdate.IsUnfillable = false
		orderToUpdate.IsExpired = false
		orderToUpdate.LastUpdated = time.Now().UTC()
		orderToUpdate.LastValidatedBlockNumber = validationBlock.Number
		orderToUpdate.LastValidatedBlockHash = validationBlock.Hash
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

func (w *Watcher) markOrderUnfillable(order *types.OrderWithMetadata, newFillableAmount *big.Int, validationBlock *types.MiniHeader) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.IsUnfillable = true
		if big.NewInt(validationBlock.Timestamp.Unix()).Cmp(orderToUpdate.OrderV3.ExpirationTimeSeconds) >= 0 {
			orderToUpdate.IsExpired = true
		}
		orderToUpdate.LastUpdated = time.Now().UTC()
		orderToUpdate.LastValidatedBlockNumber = validationBlock.Number
		orderToUpdate.LastValidatedBlockHash = validationBlock.Hash
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
}

func (w *Watcher) unwatchOrder(order *types.OrderWithMetadata, newFillableAmount *big.Int, validationBlock *types.MiniHeader) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.IsRemoved = true
		orderToUpdate.IsUnfillable = true
		if big.NewInt(validationBlock.Timestamp.Unix()).Cmp(orderToUpdate.OrderV3.ExpirationTimeSeconds) >= 0 {
			orderToUpdate.IsExpired = true
		}
		orderToUpdate.LastUpdated = time.Now().UTC()
		orderToUpdate.LastValidatedBlockNumber = validationBlock.Number
		orderToUpdate.LastValidatedBlockHash = validationBlock.Hash
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
}

func (w *Watcher) updateOrderExpirationState(order *types.OrderWithMetadata, validationBlock *types.MiniHeader) {
	err := w.db.UpdateOrder(order.Hash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		if big.NewInt(validationBlock.Timestamp.Unix()).Cmp(orderToUpdate.OrderV3.ExpirationTimeSeconds) >= 0 {
			orderToUpdate.IsExpired = true
		}
		orderToUpdate.LastUpdated = time.Now().UTC()
		orderToUpdate.LastValidatedBlockNumber = validationBlock.Number
		orderToUpdate.LastValidatedBlockHash = validationBlock.Hash
		return orderToUpdate, nil
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"order": order,
		}).Error("Failed to update order")
	}
}

func (w *Watcher) permanentlyDeleteOrder(order *types.OrderWithMetadata) error {
	if err := w.db.DeleteOrder(order.Hash); err != nil {
		return err
	}

	// After permanently deleting an order, we also remove its assetData from the Decoder
	err := w.removeAssetDataAddressFromEventDecoder(order.OrderV3.MakerAssetData)
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
	// NOTE(oskar): What happens here is that some tokens which do not
	// respect the ERC20 specification can throw ABI decoder errors, we
	// should handle them here and return it as a non-critical error.
	// TODO(oskar): Should this be handled the same way for all non-ABI
	// conforming ERC20 implementations?
	if parserError, ok := err.(decoder.AbiTopicParserError); ok {
		logger.WithFields(logger.Fields{
			"eventType":       eventType,
			"topics":          parserError.Topics,
			"contractAddress": parserError.ContractAddress,
		}).Warn("event parsing error")
		return true
	}
	logger.WithFields(logger.Fields{
		"error":     err.Error(),
		"eventType": eventType,
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
		w.contractAddressToSeenCount.Inc(decodedAssetData.Address)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC721(decodedAssetData.Address)
		w.contractAddressToSeenCount.Inc(decodedAssetData.Address)
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC1155(decodedAssetData.Address)
		w.contractAddressToSeenCount.Inc(decodedAssetData.Address)
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
		count := w.contractAddressToSeenCount.Dec(decodedAssetData.Address)
		if count == 0 {
			w.eventDecoder.RemoveKnownERC20(decodedAssetData.Address)
		}
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		count := w.contractAddressToSeenCount.Dec(decodedAssetData.Address)
		if count == 0 {
			w.eventDecoder.RemoveKnownERC721(decodedAssetData.Address)
		}
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		count := w.contractAddressToSeenCount.Dec(decodedAssetData.Address)
		if count == 0 {
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

func (w *Watcher) getExtremeBlocksFromEvents(events []*blockwatch.Event) (oldestBlock *types.MiniHeader, latestBlock *types.MiniHeader) {
	for _, event := range events {
		if latestBlock == nil && oldestBlock == nil {
			latestBlock = event.BlockHeader
			oldestBlock = event.BlockHeader
		} else {
			if event.BlockHeader.Number.Cmp(latestBlock.Number) == 1 {
				latestBlock = event.BlockHeader
			} else if event.BlockHeader.Number.Cmp(oldestBlock.Number) == -1 {
				oldestBlock = event.BlockHeader
			}
		}
	}
	return oldestBlock, latestBlock
}

// WaitForAtLeastOneBlockToBeProcessed waits until the OrderWatcher has processed its
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
