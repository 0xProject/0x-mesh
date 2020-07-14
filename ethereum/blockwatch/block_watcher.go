package blockwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum/simplestack"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	log "github.com/sirupsen/logrus"
)

// go-ethereum client `ethereum.NotFound` error type message
const rpcClientNotFoundError = "not found"

// maxBlocksInGetLogsQuery is the max number of blocks to fetch logs for in a single query. There is
// a hard limit of 10,000 logs returned by a single `eth_getLogs` query by Infura's Ethereum nodes so
// we need to try and stay below it. Parity, Geth and Alchemy all have much higher limits (if any) on
// the number of logs returned so Infura is by far the limiting factor.
var maxBlocksInGetLogsQuery = 60

// warningLevelErrorMessages are certain blockwatch.Watch errors that we want to report as warnings
// because they do not represent a bug or issue with Mesh and are expected to happen from time to time
var warningLevelErrorMessages = []string{
	constants.GethFilterUnknownBlock,
	rpcClientNotFoundError,
	"context deadline exceeded",
	constants.ParityFilterUnknownBlock,
}

// EventType describes the types of events emitted by blockwatch.Watcher. A block can be discovered
// and added to our representation of the chain. During a block re-org, a block previously stored
// can be removed from the list.
type EventType int

const (
	Added EventType = iota
	Removed
)

// Event describes a block event emitted by a Watcher
type Event struct {
	Type        EventType
	BlockHeader *types.MiniHeader
}

// TooMayBlocksBehindError is an error returned if the BlockWatcher has fallen too many blocks behind
// the latest block (>128 blocks), and cannot catch back up when connect to a non-archive Ethereum
// node.
type TooMayBlocksBehindError struct {
	blocksMissing int
}

func (e TooMayBlocksBehindError) Error() string {
	return fmt.Sprintf("too many blocks (%d) behind the latest block", e.blocksMissing)
}

// Config holds some configuration options for an instance of BlockWatcher.
type Config struct {
	DB              *db.DB
	PollingInterval time.Duration
	WithLogs        bool
	Topics          []common.Hash
	Client          Client
}

// Watcher maintains a consistent representation of the latest X blocks (where X is enforced by the
// supplied stack) handling block re-orgs and network disruptions gracefully. It can be started from
// any arbitrary block height, and will emit both block added and removed events.
type Watcher struct {
	stack               *simplestack.SimpleStack
	db                  *db.DB
	client              Client
	blockFeed           event.Feed
	blockScope          event.SubscriptionScope // Subscription scope tracking current live listeners
	wasStartedOnce      bool                    // Whether the block watcher has previously been started
	pollingInterval     time.Duration
	withLogs            bool
	topics              []common.Hash
	mu                  sync.RWMutex
	syncToLatestBlockMu sync.Mutex
}

// New creates a new Watcher instance.
func New(retentionLimit int, config Config) *Watcher {
	return &Watcher{
		pollingInterval: config.PollingInterval,
		db:              config.DB,
		stack:           simplestack.New(retentionLimit, []*types.MiniHeader{}),
		client:          config.Client,
		withLogs:        config.WithLogs,
		topics:          config.Topics,
	}
}

// FastSyncToLatestBlock checks if the BlockWatcher is behind the latest block, and if so,
// catches it back up. If less than 128 blocks passed, we are able to fetch all missing
// block events and process them. If more than 128 blocks passed, we cannot catch up
// without an archive Ethereum node (see: http://bit.ly/2D11Hr6) so we instead clear
// previously tracked blocks so BlockWatcher starts again from the latest block. This
// function blocks until complete or the context is  cancelled.
func (w *Watcher) FastSyncToLatestBlock(ctx context.Context) (blocksElapsed int, err error) {
	w.mu.Lock()
	if w.wasStartedOnce {
		w.mu.Unlock()
		return 0, errors.New("Can only fast-sync to latest block before starting BlockWatcher")
	}
	w.mu.Unlock()

	latestBlockProcessed, err := w.stack.Peek()
	if err != nil {
		return 0, err
	}
	// No previously stored block so no blocks have elapsed
	if latestBlockProcessed == nil {
		return 0, nil
	}

	latestBlock, err := w.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	latestBlockProcessedNumber := int(latestBlockProcessed.Number.Int64())
	blocksElapsed = int(latestBlock.Number.Int64()) - latestBlockProcessedNumber
	if blocksElapsed == 0 {
		return blocksElapsed, nil
	} else if blocksElapsed < constants.MaxBlocksStoredInNonArchiveNode {
		log.WithField("blocksElapsed", blocksElapsed).Info("Some blocks have elapsed since last boot. Backfilling block events (this can take a while)...")
		events, err := w.getMissedEventsToBackfill(ctx, blocksElapsed, latestBlockProcessedNumber)
		if err != nil {
			return blocksElapsed, err
		}
		if len(events) > 0 {
			newMiniHeaders, err := w.stack.PeekAll()
			if err != nil {
				return blocksElapsed, err
			}
			if _, _, err := w.db.ResetMiniHeaders(newMiniHeaders); err != nil {
				return blocksElapsed, err
			}
			w.blockFeed.Send(events)
		}
	} else {
		// Clear all block headers from stack and database so BlockWatcher
		// starts again from latest block
		if _, err := w.db.DeleteMiniHeaders(nil); err != nil {
			return blocksElapsed, err
		}
		if err := w.stack.Clear(); err != nil {
			return blocksElapsed, err
		}
	}

	return blocksElapsed, nil
}

// Watch starts the Watcher. It will continuously look for new blocks and blocks
// until there is a critical error or the given context is canceled. Typically,
// you want to call Watch inside a goroutine. For non-critical errors, callers
// must receive them from the Errors channel.
func (w *Watcher) Watch(ctx context.Context) error {
	w.mu.Lock()
	if w.wasStartedOnce {
		w.mu.Unlock()
		return errors.New("Can only start Watcher once per instance")
	}
	w.wasStartedOnce = true
	w.mu.Unlock()

	// Sync immediately when `Watch()` is called instead of waiting for the
	// first Ticker tick
	if err := w.SyncToLatestBlock(ctx); err != nil {
		if err == db.ErrClosed {
			// We can't continue if the database is closed. Stop the watcher and
			// return an error.
			return err
		}
		logMessage := "blockwatch.Watcher error encountered"
		if isWarning(err) {
			log.WithError(err).Warn(logMessage)
		} else {
			log.WithError(err).Error(logMessage)
		}
	}

	ticker := time.NewTicker(w.pollingInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			if err := w.SyncToLatestBlock(ctx); err != nil {
				if err == db.ErrClosed {
					// We can't continue if the database is closed. Stop the watcher and
					// return an error.
					ticker.Stop()
					return err
				}
				if _, ok := err.(TooMayBlocksBehindError); ok {
					// We've fallen too many blocks behind to sync to the latest block.
					// We'd need to start again from the latest block but also require
					// the OrderWatcher to re-validate all orders at the latest block.
					// By returning an error here, we cause Mesh to gracefully shut down.
					// Upon re-booting, it will reset the blocks stored in the DB and
					// re-validate all orders stored.
					ticker.Stop()
					return err
				}
				logMessage := "blockwatch.Watcher error encountered"
				if isWarning(err) {
					log.WithError(err).Warn(logMessage)
				} else {
					log.WithError(err).Error(logMessage)
				}
			}
		}
	}
}

// Subscribe allows one to subscribe to the block events emitted by the Watcher.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*Event) event.Subscription {
	return w.blockScope.Track(w.blockFeed.Subscribe(sink))
}

// SyncToLatestBlock syncs our local state of the chain to the latest block found via
// Ethereum RPC
func (w *Watcher) SyncToLatestBlock(ctx context.Context) error {
	w.syncToLatestBlockMu.Lock()
	defer w.syncToLatestBlockMu.Unlock()

	checkpoint, err := w.stack.Checkpoint()
	if err != nil {
		return err
	}

	latestHeader, err := w.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	latestBlockNumber := latestHeader.Number.Int64()
	lastStoredHeader, err := w.stack.Peek()
	if err != nil {
		return err
	}
	var lastStoredBlockNumber int64
	if lastStoredHeader != nil {
		lastStoredBlockNumber = lastStoredHeader.Number.Int64()
	}

	var numBlocksToFetch int
	// No blocks stored yet, fetch the first
	if lastStoredHeader == nil {
		numBlocksToFetch = 1
	} else {
		// Noop if already caught up or ahead of latest block returned from Ethereum node
		if latestBlockNumber <= lastStoredBlockNumber {
			return nil
		}
		numBlocksToFetch = int(latestBlockNumber - lastStoredBlockNumber)
	}

	if numBlocksToFetch >= constants.MaxBlocksStoredInNonArchiveNode {
		return TooMayBlocksBehindError{
			blocksMissing: numBlocksToFetch,
		}
	}

	allEvents := []*Event{}
	// Syncing to the latest block involves multiple Ethereum RPC requests. If any of them fail, we
	// stop syncing and set the encountered error to `syncErr` to be returned to the caller after we've
	// either reset or persisted the changes gathered up until the point where the error occurred.
	var syncErr error
	for i := 0; i < numBlocksToFetch; i++ {
		// Optimization: If numBlocksToFetch is 1, we already know what the nextHeader is, so avoid
		// fetching it again. If there is more then 1 block to fetch, compute each from the last
		// stored and fetch it
		nextHeader := latestHeader
		if numBlocksToFetch != 1 {
			lastStoredHeader, err := w.stack.Peek()
			if err != nil {
				syncErr = err
				break
			}
			nextBlockNumber := big.NewInt(0).Add(lastStoredHeader.Number, big.NewInt(1))
			nextHeader, err = w.client.HeaderByNumber(ctx, nextBlockNumber)
			if err != nil {
				syncErr = err
				break
			}
		}

		var events []*Event
		events, err = w.buildCanonicalChain(ctx, nextHeader, events)
		allEvents = append(allEvents, events...)
		if err != nil {
			syncErr = err
			break
		}
	}
	if len(allEvents) == 0 {
		return syncErr
	}
	if w.shouldRevertChanges(lastStoredHeader, allEvents) {
		if err := w.stack.Reset(checkpoint); err != nil {
			return err
		}
	} else {
		_, err := w.stack.Checkpoint()
		if err != nil {
			return err
		}
		newMiniHeaders, err := w.stack.PeekAll()
		if _, _, err := w.db.ResetMiniHeaders(newMiniHeaders); err != nil {
			return err
		}
		w.blockFeed.Send(allEvents)
	}

	return syncErr
}

func (w *Watcher) shouldRevertChanges(lastStoredHeader *types.MiniHeader, events []*Event) bool {
	if len(events) == 0 || lastStoredHeader == nil {
		return false
	}
	// If we haven't progressed in terms of block number, revert back to previous "latest" block.
	// This ensures block events always leave the node further ahead, preventing unnecessary thrash
	// during block-reorgs (which tend to cluster)
	newLatestHeader := events[len(events)-1].BlockHeader
	return newLatestHeader.Number.Cmp(lastStoredHeader.Number) <= 0
}

func (w *Watcher) buildCanonicalChain(ctx context.Context, nextHeader *types.MiniHeader, events []*Event) ([]*Event, error) {
	latestHeader, err := w.stack.Peek()
	if err != nil {
		return nil, err
	}
	// Is the stack empty or is it the next block?
	if latestHeader == nil || nextHeader.Parent == latestHeader.Hash {
		nextHeader, err := w.addLogs(ctx, nextHeader)
		if err != nil {
			return events, err
		}
		err = w.stack.Push(nextHeader)
		if err != nil {
			return events, err
		}
		events = append(events, &Event{
			Type:        Added,
			BlockHeader: nextHeader,
		})
		return events, nil
	}

	// Pop latestHeader from the stack. We already have a reference to it.
	if _, err := w.stack.Pop(); err != nil {
		return events, err
	}
	events = append(events, &Event{
		Type:        Removed,
		BlockHeader: latestHeader,
	})

	nextParentHeader, err := w.client.HeaderByHash(ctx, nextHeader.Parent)
	if err != nil {
		return events, err
	}
	events, err = w.buildCanonicalChain(ctx, nextParentHeader, events)
	if err != nil {
		return events, err
	}
	nextHeader, err = w.addLogs(ctx, nextHeader)
	if err != nil {
		return events, err
	}
	err = w.stack.Push(nextHeader)
	if err != nil {
		return events, err
	}
	events = append(events, &Event{
		Type:        Added,
		BlockHeader: nextHeader,
	})

	return events, nil
}

func (w *Watcher) addLogs(ctx context.Context, header *types.MiniHeader) (*types.MiniHeader, error) {
	if !w.withLogs {
		return header, nil
	}
	logs, err := w.client.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &header.Hash,
		Topics:    [][]common.Hash{w.topics},
	})
	if err != nil {
		return header, err
	}
	header.Logs = logs
	return header, nil
}

// getMissedEventsToBackfill finds missed events that might have occurred while the Mesh node was
// offline. It does this by comparing the last block stored with the latest block discoverable via RPC.
// If the stored block is older then the latest block, it batch fetches the events for missing blocks,
// re-sets the stored blocks and returns the block events found.
func (w *Watcher) getMissedEventsToBackfill(ctx context.Context, blocksElapsed int, latestRetainedBlockNumber int) ([]*Event, error) {
	events := []*Event{}

	startBlockNum := latestRetainedBlockNumber + 1
	endBlockNum := latestRetainedBlockNumber + blocksElapsed
	logs, furthestBlockProcessed := w.getLogsInBlockRange(ctx, startBlockNum, endBlockNum)
	if furthestBlockProcessed > latestRetainedBlockNumber {
		// If we have processed blocks further then the latestRetainedBlock in the DB, we
		// want to remove all blocks from the DB and insert the furthestBlockProcessed
		// Doing so will cause the BlockWatcher to start from that furthestBlockProcessed.
		if err := w.stack.Clear(); err != nil {
			return events, err
		}
		// Add furthest block processed into the DB
		latestHeader, err := w.client.HeaderByNumber(ctx, big.NewInt(int64(furthestBlockProcessed)))
		if err != nil {
			return events, err
		}
		err = w.stack.Push(latestHeader)
		if err != nil {
			return events, err
		}

		// If no logs found, noop
		if len(logs) == 0 {
			return events, nil
		}

		// Create the block events from all the logs found by grouping
		// them into blockHeaders
		hashToBlockHeader := map[common.Hash]*types.MiniHeader{}
		for _, log := range logs {
			blockHeader, found := hashToBlockHeader[log.BlockHash]
			if !found {
				blockNumber := big.NewInt(0).SetUint64(log.BlockNumber)
				header, err := w.client.HeaderByNumber(ctx, blockNumber)
				if err != nil {
					return events, err
				}
				blockHeader = &types.MiniHeader{
					Hash:      log.BlockHash,
					Parent:    header.Parent,
					Number:    blockNumber,
					Logs:      []ethtypes.Log{},
					Timestamp: header.Timestamp,
				}
				hashToBlockHeader[log.BlockHash] = blockHeader
			}
			blockHeader.Logs = append(blockHeader.Logs, log)
		}
		for _, blockHeader := range hashToBlockHeader {
			events = append(events, &Event{
				Type:        Added,
				BlockHeader: blockHeader,
			})
		}
		log.Info("Done backfilling block events")
		return events, nil
	}
	return events, nil
}

type logRequestResult struct {
	From int
	To   int
	Logs []ethtypes.Log
	Err  error
}

// getLogsRequestChunkSize is the number of `eth_getLogs` JSON RPC to send concurrently in each batch fetch
const getLogsRequestChunkSize = 3

// getLogsInBlockRange attempts to fetch all logs in the block range supplied. It implements a
// limited-concurrency batch fetch, where all requests in the previous batch must complete for
// the next batch of requests to be sent. If an error is encountered in a batch, all subsequent
// batch requests are not sent. Instead, it returns all the logs it found up until the error was
// encountered, along with the block number after which no further logs were retrieved.
func (w *Watcher) getLogsInBlockRange(ctx context.Context, from, to int) ([]ethtypes.Log, int) {
	blockRanges := w.getSubBlockRanges(from, to, maxBlocksInGetLogsQuery)

	numChunks := 0
	chunkChan := make(chan []*blockRange, 1000000)
	for len(blockRanges) != 0 {
		var chunk []*blockRange
		if len(blockRanges) < getLogsRequestChunkSize {
			chunk = blockRanges
		} else {
			chunk = blockRanges[:getLogsRequestChunkSize]
		}
		chunkChan <- chunk
		blockRanges = blockRanges[len(chunk):]
		numChunks++
	}

	semaphoreChan := make(chan struct{}, 1)
	defer close(semaphoreChan)

	didAPreviousRequestFail := false
	furthestBlockProcessed := from - 1
	allLogs := []ethtypes.Log{}

	for i := 0; i < numChunks; i++ {
		// Add one to the semaphore chan. If it already has a value, the chunk blocks here until one frees up.
		// We deliberately process the chunks sequentially, since if any request results in an error, we
		// do not want to send any further requests.
		semaphoreChan <- struct{}{}

		// If a previous request failed, we stop processing newer requests
		if didAPreviousRequestFail {
			<-semaphoreChan
			continue // Noop
		}

		mu := sync.Mutex{}
		indexToLogResult := map[int]logRequestResult{}
		chunk := <-chunkChan

		wg := &sync.WaitGroup{}
		for i, aBlockRange := range chunk {
			wg.Add(1)
			go func(index int, b *blockRange) {
				defer wg.Done()

				select {
				case <-ctx.Done():
					indexToLogResult[index] = logRequestResult{
						From: b.FromBlock,
						To:   b.ToBlock,
						Err:  errors.New("context was canceled"),
						Logs: []ethtypes.Log{},
					}
					return
				default:
				}

				logs, err := w.filterLogsRecursively(ctx, b.FromBlock, b.ToBlock, []ethtypes.Log{})
				if err != nil {
					log.WithFields(map[string]interface{}{
						"error":     err,
						"fromBlock": b.FromBlock,
						"toBlock":   b.ToBlock,
					}).Trace("Failed to fetch logs for range")
				}
				mu.Lock()
				indexToLogResult[index] = logRequestResult{
					From: b.FromBlock,
					To:   b.ToBlock,
					Err:  err,
					Logs: logs,
				}
				mu.Unlock()
			}(i, aBlockRange)
		}

		// Wait for all log requests to complete
		wg.Wait()

		for i, aBlockRange := range chunk {
			logRequestResult := indexToLogResult[i]
			// Break at first error encountered
			if logRequestResult.Err != nil {
				didAPreviousRequestFail = true
				furthestBlockProcessed = logRequestResult.From - 1
				break
			}
			allLogs = append(allLogs, logRequestResult.Logs...)
			furthestBlockProcessed = aBlockRange.ToBlock
		}
		<-semaphoreChan
	}

	return allLogs, furthestBlockProcessed
}

type blockRange struct {
	FromBlock int
	ToBlock   int
}

// getSubBlockRanges breaks up the block range into smaller block ranges of rangeSize.
// `eth_getLogs` requests are inclusive to both the start and end blocks specified and
// so we need to make the ranges exclusive of one another to avoid fetching the same
// blocks' logs twice.
func (w *Watcher) getSubBlockRanges(from, to, rangeSize int) []*blockRange {
	chunks := []*blockRange{}
	numBlocksLeft := to - from + 1
	if numBlocksLeft < rangeSize {
		chunks = append(chunks, &blockRange{
			FromBlock: from,
			ToBlock:   to,
		})
	} else {
		numChunks := numBlocksLeft / rangeSize
		remainder := numBlocksLeft % rangeSize
		if remainder > 0 {
			numChunks = numChunks + 1
		}

		for i := 0; i < numChunks; i = i + 1 {
			chunkFromBlock := from + i*rangeSize
			chunkToBlock := chunkFromBlock + rangeSize - 1
			if chunkToBlock > to {
				chunkToBlock = to
			}
			blockRange := &blockRange{
				FromBlock: chunkFromBlock,
				ToBlock:   chunkToBlock,
			}
			chunks = append(chunks, blockRange)
		}
	}
	return chunks
}

const infuraTooManyResultsErrMsg = "query returned more than 10000 results"

func (w *Watcher) filterLogsRecursively(ctx context.Context, from, to int, allLogs []ethtypes.Log) ([]ethtypes.Log, error) {
	log.WithFields(map[string]interface{}{
		"from": from,
		"to":   to,
	}).Trace("Fetching block logs")
	numBlocks := to - from
	topics := [][]common.Hash{}
	if len(w.topics) > 0 {
		topics = append(topics, w.topics)
	}
	logs, err := w.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Topics:    topics,
	})
	if err != nil {
		// Infura caps the logs returned to 10,000 per request, if our request exceeds this limit, split it
		// into two requests. Parity, Geth and Alchemy all have much higher limits (if any at all), so no need
		// to expect any similar errors of this nature from them.
		if err.Error() == infuraTooManyResultsErrMsg {
			// HACK(fabio): Infura limits the returned results to 10,000 logs, BUT some single
			// blocks contain more then 10,000 logs. This has supposedly been fixed but we keep
			// this logic here just in case. It helps us avoid infinite recursion.
			// Source: https://community.infura.io/t/getlogs-error-query-returned-more-than-1000-results/358/10
			if from == to {
				return allLogs, fmt.Errorf("Unable to get the logs for block #%d, because it contains too many logs", from)
			}

			r := numBlocks % 2
			firstBatchSize := numBlocks / 2
			secondBatchSize := firstBatchSize
			if r == 1 {
				secondBatchSize = secondBatchSize + 1
			}

			endFirstHalf := from + firstBatchSize
			startSecondHalf := endFirstHalf + 1
			allLogs, err := w.filterLogsRecursively(ctx, from, endFirstHalf, allLogs)
			if err != nil {
				return nil, err
			}
			allLogs, err = w.filterLogsRecursively(ctx, startSecondHalf, to, allLogs)
			if err != nil {
				return nil, err
			}
			return allLogs, nil
		} else {
			return nil, err
		}
	}
	allLogs = append(allLogs, logs...)
	return allLogs, nil
}

func isWarning(err error) bool {
	message := err.Error()
	for _, warningLevelErrorMessage := range warningLevelErrorMessages {
		if strings.Contains(message, warningLevelErrorMessage) {
			return true
		}
	}
	return false
}
