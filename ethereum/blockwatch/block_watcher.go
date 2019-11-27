package blockwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

// maxBlocksInGetLogsQuery is the max number of blocks to fetch logs for in a single query. There is
// a hard limit of 10,000 logs returned by a single `eth_getLogs` query by Infura's Ethereum nodes so
// we need to try and stay below it. Parity, Geth and Alchemy all have much higher limits (if any) on
// the number of logs returned so Infura is by far the limiting factor.
var maxBlocksInGetLogsQuery = 60

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
	BlockHeader *miniheader.MiniHeader
}

// Stack defines the interface a stack must implement in order to be used by
// OrderWatcher for block header storage
type Stack interface {
	Pop() (*miniheader.MiniHeader, error)
	Push(*miniheader.MiniHeader) error
	Peek() (*miniheader.MiniHeader, error)
	PeekAll() ([]*miniheader.MiniHeader, error)
	Clear() error
}

// Config holds some configuration options for an instance of BlockWatcher.
type Config struct {
	Stack           Stack
	PollingInterval time.Duration
	StartBlockDepth rpc.BlockNumber
	WithLogs        bool
	Topics          []common.Hash
	Client          Client
}

// Watcher maintains a consistent representation of the latest X blocks (where X is enforced by the
// supplied stack) handling block re-orgs and network disruptions gracefully. It can be started from
// any arbitrary block height, and will emit both block added and removed events.
type Watcher struct {
	startBlockDepth rpc.BlockNumber
	stack           Stack
	client          Client
	blockFeed       event.Feed
	blockScope      event.SubscriptionScope // Subscription scope tracking current live listeners
	wasStartedOnce  bool                    // Whether the block watcher has previously been started
	pollingInterval time.Duration
	ticker          *time.Ticker
	withLogs        bool
	topics          []common.Hash
	mu              sync.RWMutex

	didProcessABlock bool
	// AtLeastOneBlockProcessed is closed to signal that the BlockWatcher has processed at least one
	// block. Validation of orders should block until this has completed
	AtLeastOneBlockProcessed chan struct{}
}

// New creates a new Watcher instance.
func New(config Config) *Watcher {
	bs := &Watcher{
		pollingInterval: config.PollingInterval,
		startBlockDepth: config.StartBlockDepth,
		stack:           config.Stack,
		client:          config.Client,
		withLogs:        config.WithLogs,
		topics:          config.Topics,

		didProcessABlock:         false,
		AtLeastOneBlockProcessed: make(chan struct{}),
	}
	return bs
}

// SyncToLatestBlock checks if the BlockWatcher is behind the latest block, and if so,
// catches it back up. If less than 128 blocks passed, we are able to fetch all missing
// block events and process them. If more than 128 blocks passed, we cannot catch up
// without an archive Ethereum node (see: http://bit.ly/2D11Hr6) so we instead clear
// previously tracked blocks so BlockWatcher starts again from the latest block. This
// function blocks until complete or the context is  cancelled.
func (w *Watcher) SyncToLatestBlock(ctx context.Context) (blocksElapsed int, err error) {
	w.mu.Lock()
	if w.wasStartedOnce {
		w.mu.Unlock()
		return 0, errors.New("Can only sync to latest block before starting BlockWatcher")
	}
	w.mu.Unlock()

	latestBlockProcessed, err := w.GetLatestBlockProcessed()
	if err != nil {
		return 0, err
	}
	// No previously stored block so no blocks have elapsed
	if latestBlockProcessed == nil {
		return 0, nil
	}

	latestBlock, err := w.client.HeaderByNumber(nil)
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
			w.mu.Lock()
			if !w.didProcessABlock {
				w.didProcessABlock = true
				close(w.AtLeastOneBlockProcessed)
			}
			w.mu.Unlock()
			w.blockFeed.Send(events)
		}
	} else {
		// Clear all block headers from stack so BlockWatcher starts again from latest block
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

	if err := w.PollNextBlock(); err != nil {
		if err == leveldb.ErrClosed {
			// We can't continue if the database is closed. Stop the watcher and
			// return an error.
			return err
		}
		log.WithError(err).Error("blockwatch.Watcher error encountered")
	}

	ticker := time.NewTicker(w.pollingInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			if err := w.PollNextBlock(); err != nil {
				if err == leveldb.ErrClosed {
					// We can't continue if the database is closed. Stop the watcher and
					// return an error.
					ticker.Stop()
					return err
				}
				log.WithError(err).Error("blockwatch.Watcher error encountered")
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

// GetLatestBlockProcessed returns the latest block processed
func (w *Watcher) GetLatestBlockProcessed() (*miniheader.MiniHeader, error) {
	return w.stack.Peek()
}

// GetAllRetainedBlocks returns the blocks retained in-memory by the Watcher.
func (w *Watcher) GetAllRetainedBlocks() ([]*miniheader.MiniHeader, error) {
	return w.stack.PeekAll()
}

// PollNextBlock polls for the next block header to be added to the block stack.
// If there are no blocks on the stack, it fetches the first block at the specified
// `startBlockDepth` supplied at instantiation.
func (w *Watcher) PollNextBlock() error {
	var nextBlockNumber *big.Int
	latestHeader, err := w.stack.Peek()
	if err != nil {
		return err
	}
	if latestHeader == nil {
		if w.startBlockDepth == rpc.LatestBlockNumber {
			nextBlockNumber = nil // Fetch latest block
		} else {
			nextBlockNumber = big.NewInt(int64(w.startBlockDepth))
		}
	} else {
		nextBlockNumber = big.NewInt(0).Add(latestHeader.Number, big.NewInt(1))
	}
	nextHeader, err := w.client.HeaderByNumber(nextBlockNumber)
	if err != nil {
		if err == ethereum.NotFound {
			log.WithFields(log.Fields{
				"blockNumber": nextBlockNumber,
			}).Trace("block header not found")
			return nil // Noop and wait next polling interval
		}
		return err
	}

	events := []*Event{}
	events, err = w.buildCanonicalChain(nextHeader, events)
	if err != nil {
		// If error encountered after some block events generated
		if len(events) != 0 {
			newLatestHeader := events[len(events)-1].BlockHeader
			// If we haven't progressed in terms of block number before an error was encountered,
			// revert back to previous "latest" block. This ensures block events always leave the
			// node further ahead, preventing unnecessary thrash during block-reorgs (which tend to cluster)
			if newLatestHeader.Number.Int64() <= latestHeader.Number.Int64() {
				for i := len(events) - 1; i >= 0; i-- {
					event := events[i]
					switch event.Type {
					case Added:
						_, err := w.stack.Pop()
						if err != nil {
							return err // Could only be unexpected DB errors
						}
					case Removed:
						if err := w.stack.Push(event.BlockHeader); err != nil {
							return err // Could only be unexpected DB errors
						}
					default:
						log.WithField("Type", event.Type).Panic("Unrecognized event.Type encountered")
					}
				}
			}
		}
		return err
	}
	if len(events) > 0 {
		// Upon processing the first block event, consider the BlockWatcher started
		w.mu.Lock()
		if !w.didProcessABlock {
			w.didProcessABlock = true
			close(w.AtLeastOneBlockProcessed)
		}
		w.mu.Unlock()
		w.blockFeed.Send(events)
	}
	return nil
}

func (w *Watcher) buildCanonicalChain(nextHeader *miniheader.MiniHeader, events []*Event) ([]*Event, error) {
	latestHeader, err := w.stack.Peek()
	if err != nil {
		return nil, err
	}
	// Is the stack empty or is it the next block?
	if latestHeader == nil || nextHeader.Parent == latestHeader.Hash {
		nextHeader, err := w.addLogs(nextHeader)
		if err != nil {
			// Due to block re-orgs & Ethereum node services load-balancing requests across multiple nodes
			// a block header might be returned, but when fetching it's logs, an "unknown block" error is
			// returned. This is expected to happen sometimes, and we simply return the events gathered so
			// far and pick back up where we left off on the next polling interval.
			if isUnknownBlockErr(err) {
				log.WithFields(log.Fields{
					"nextHeader": nextHeader,
				}).Trace("failed to get logs for block")
				return events, nil
			}
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

	nextParentHeader, err := w.client.HeaderByHash(nextHeader.Parent)
	if err != nil {
		if err == ethereum.NotFound {
			log.WithFields(log.Fields{
				"blockNumber": nextHeader.Parent.Hex(),
			}).Info("block header not found")
			// Noop and wait next polling interval
			return events, nil
		}
		return events, err
	}
	events, err = w.buildCanonicalChain(nextParentHeader, events)
	if err != nil {
		return events, err
	}
	nextHeader, err = w.addLogs(nextHeader)
	if err != nil {
		// Due to block re-orgs & Ethereum node services load-balancing requests across multiple nodes
		// a block header might be returned, but when fetching it's logs, an "unknown block" error is
		// returned. This is expected to happen sometimes, and we simply return the events gathered so
		// far and pick back up where we left off on the next polling interval.
		if isUnknownBlockErr(err) {
			log.WithFields(log.Fields{
				"nextHeader": nextHeader,
			}).Trace("failed to get logs for block")
			return events, nil
		}
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

func (w *Watcher) addLogs(header *miniheader.MiniHeader) (*miniheader.MiniHeader, error) {
	if !w.withLogs {
		return header, nil
	}
	logs, err := w.client.FilterLogs(ethereum.FilterQuery{
		BlockHash: &header.Hash,
		Topics:    [][]common.Hash{w.topics},
	})
	if err != nil {
		return header, err
	}
	header.Logs = logs
	return header, nil
}

// getMissedEventsToBackfill finds missed events that might have occured while the Mesh node was
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
		latestHeader, err := w.client.HeaderByNumber(big.NewInt(int64(furthestBlockProcessed)))
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
		hashToBlockHeader := map[common.Hash]*miniheader.MiniHeader{}
		for _, log := range logs {
			blockHeader, ok := hashToBlockHeader[log.BlockHash]
			if !ok {
				blockNumber := big.NewInt(0).SetUint64(log.BlockNumber)
				header, err := w.client.HeaderByNumber(blockNumber)
				if err != nil {
					return events, err
				}
				blockHeader = &miniheader.MiniHeader{
					Hash:      log.BlockHash,
					Parent:    header.Parent,
					Number:    blockNumber,
					Logs:      []types.Log{},
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
	Logs []types.Log
	Err  error
}

// getLogsRequestChunkSize is the number of `eth_getLogs` JSON RPC to send concurrently in each batch fetch
const getLogsRequestChunkSize = 3

// getLogsInBlockRange attempts to fetch all logs in the block range supplied. It implements a
// limited-concurrency batch fetch, where all requests in the previous batch must complete for
// the next batch of requests to be sent. If an error is encountered in a batch, all subsequent
// batch requests are not sent. Instead, it returns all the logs it found up until the error was
// encountered, along with the block number after which no further logs were retrieved.
func (w *Watcher) getLogsInBlockRange(ctx context.Context, from, to int) ([]types.Log, int) {
	blockRanges := w.getSubBlockRanges(from, to, maxBlocksInGetLogsQuery)

	numChunks := 0
	chunkChan := make(chan []*blockRange, 1000000)
	for len(blockRanges) != 0 {
		var chunk []*blockRange
		if len(blockRanges) < getLogsRequestChunkSize {
			chunk = blockRanges[:len(blockRanges)]
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
	allLogs := []types.Log{}

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
						Logs: []types.Log{},
					}
					return
				default:
				}

				logs, err := w.filterLogsRecurisively(b.FromBlock, b.ToBlock, []types.Log{})
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
	numBlocksLeft := to - from
	if numBlocksLeft < rangeSize {
		chunks = append(chunks, &blockRange{
			FromBlock: from,
			ToBlock:   to,
		})
	} else {
		blocks := []int{}
		for i := 0; i <= numBlocksLeft; i++ {
			blocks = append(blocks, from+i)
		}
		numChunks := len(blocks) / rangeSize
		remainder := len(blocks) % rangeSize
		if remainder > 0 {
			numChunks = numChunks + 1
		}

		for i := 0; i < numChunks; i = i + 1 {
			fromIndex := i * rangeSize
			toIndex := fromIndex + rangeSize
			if toIndex > len(blocks) {
				toIndex = len(blocks)
			}
			bs := blocks[fromIndex:toIndex]
			blockRange := &blockRange{
				FromBlock: bs[0],
				ToBlock:   bs[len(bs)-1],
			}
			chunks = append(chunks, blockRange)
		}
	}
	return chunks
}

const infuraTooManyResultsErrMsg = "query returned more than 10000 results"

func (w *Watcher) filterLogsRecurisively(from, to int, allLogs []types.Log) ([]types.Log, error) {
	log.WithFields(map[string]interface{}{
		"from": from,
		"to":   to,
	}).Trace("Fetching block logs")
	numBlocks := to - from
	topics := [][]common.Hash{}
	if len(w.topics) > 0 {
		topics = append(topics, w.topics)
	}
	logs, err := w.client.FilterLogs(ethereum.FilterQuery{
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
			allLogs, err := w.filterLogsRecurisively(from, endFirstHalf, allLogs)
			if err != nil {
				return nil, err
			}
			allLogs, err = w.filterLogsRecurisively(startSecondHalf, to, allLogs)
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

func isUnknownBlockErr(err error) bool {
	// Geth error
	if err.Error() == "unknown block" {
		return true
	}
	// Parity error
	if err.Error() == "One of the blocks specified in filter (fromBlock, toBlock or blockHash) cannot be found" {
		return true
	}
	return false
}
