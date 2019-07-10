package blockwatch

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

// Number of concurrent JSON RPC requests to allow while fast-syncing logs
const concurrencyLimit = 3

// maxBlocksInGetLogsQuery is the max number of blocks to fetch logs for in a single query. There is
// a hard limit of 10,000 logs returned by a single `eth_getLogs` query by Infura's Ethereum nodes so
// we need to try and stay below it. Parity, Geth and Alchemy all have much higher limits (if any) on
// the number of logs returned so Infura is the limiting factor here.
var maxBlocksInGetLogsQuery = int64(60)

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
	BlockHeader *meshdb.MiniHeader
}

// Config holds some configuration options for an instance of BlockWatcher.
type Config struct {
	MeshDB              *meshdb.MeshDB
	PollingInterval     time.Duration
	StartBlockDepth     rpc.BlockNumber
	BlockRetentionLimit int
	WithLogs            bool
	Topics              []common.Hash
	Client              Client
}

// Watcher maintains a consistent representation of the latest `blockRetentionLimit` blocks,
// handling block re-orgs and network disruptions gracefully. It can be started from any arbitrary
// block height, and will emit both block added and removed events.
type Watcher struct {
	Errors              chan error
	blockRetentionLimit int
	startBlockDepth     rpc.BlockNumber
	stack               *Stack
	client              Client
	blockFeed           event.Feed
	blockScope          event.SubscriptionScope // Subscription scope tracking current live listeners
	isWatching          bool                    // Whether the block poller is running
	pollingInterval     time.Duration
	ticker              *time.Ticker
	withLogs            bool
	topics              []common.Hash
	mu                  sync.RWMutex
}

// New creates a new Watcher instance.
func New(config Config) *Watcher {
	stack := NewStack(config.MeshDB, config.BlockRetentionLimit)

	// Buffer the first 5 errors, if no channel consumer processing the errors, any additional errors are dropped
	errorsChan := make(chan error, 5)
	bs := &Watcher{
		Errors:              errorsChan,
		pollingInterval:     config.PollingInterval,
		blockRetentionLimit: config.BlockRetentionLimit,
		startBlockDepth:     config.StartBlockDepth,
		stack:               stack,
		client:              config.Client,
		withLogs:            config.WithLogs,
		topics:              config.Topics,
	}
	return bs
}

// Start starts the BlockWatcher
func (w *Watcher) Start() error {
	err := w.backfillMissedEventsIfNeeded()
	if err != nil {
		return err
	}

	// We need the mutex to reliably start/stop the update loop
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.isWatching {
		return errors.New("Polling already started")
	}

	w.isWatching = true
	if w.ticker == nil {
		w.ticker = time.NewTicker(w.pollingInterval)
	}
	go w.startPollingLoop()
	return nil
}

func (w *Watcher) startPollingLoop() {
	for {
		w.mu.Lock()
		if !w.isWatching {
			w.mu.Unlock()
			return
		}
		<-w.ticker.C
		w.mu.Unlock()

		err := w.pollNextBlock()
		if err != nil {
			// Attempt to send errors but if buffered channel is full, we assume there is no
			// interested consumer and drop them. The Watcher recovers gracefully from errors.
			select {
			case w.Errors <- err:
			default:
			}
		}
	}
}

// stopPolling stops the block poller
func (w *Watcher) stopPolling() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isWatching = false
	if w.ticker != nil {
		w.ticker.Stop()
	}
	w.ticker = nil
}

// Stop stops the BlockWatcher
func (w *Watcher) Stop() {
	if w.isWatching {
		w.stopPolling()
	}
	close(w.Errors)
}

// Subscribe allows one to subscribe to the block events emitted by the Watcher.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*Event) event.Subscription {
	return w.blockScope.Track(w.blockFeed.Subscribe(sink))
}

// InspectRetainedBlocks returns the blocks retained in-memory by the Watcher instance. It is not
// particularly performant and therefore should only be used for debugging and testing purposes.
func (w *Watcher) InspectRetainedBlocks() ([]*meshdb.MiniHeader, error) {
	return w.stack.Inspect()
}

// pollNextBlock polls for the next block header to be added to the block stack.
// If there are no blocks on the stack, it fetches the first block at the specified
// `startBlockDepth` supplied at instantiation.
func (w *Watcher) pollNextBlock() error {
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
	// Even if an error occurred, we still want to emit the events gathered since we might have
	// popped blocks off the Stack and they won't be re-added
	if len(events) != 0 {
		w.blockFeed.Send(events)
	}
	if err != nil {
		return err
	}
	return nil
}

func (w *Watcher) buildCanonicalChain(nextHeader *meshdb.MiniHeader, events []*Event) ([]*Event, error) {
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
			if err.Error() == "unknown block" {
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
			// Noop and wait next polling interval. We remove the popped blocks
			// and refetch them on the next polling interval.
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
		if err.Error() == "unknown block" {
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

func (w *Watcher) addLogs(header *meshdb.MiniHeader) (*meshdb.MiniHeader, error) {
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

// backfillMissedEventsIfNeeded backfills missed events that might have occured while the Mesh node was
// offline. It does this by comparing the last block stored with the latest block discoverable via RPC.
// If the stored block is older then the latest block, it batch backfills the events for missing blocks,
// emits them, and re-sets the stored blocks to start at the latest block.
func (w *Watcher) backfillMissedEventsIfNeeded() error {
	latestRetainedBlock, err := w.stack.Peek()
	if err != nil {
		return err
	}
	// No blocks stored, nowhere to backfill to
	if latestRetainedBlock == nil {
		return nil
	}
	latestBlock, err := w.client.HeaderByNumber(nil)
	if err != nil {
		return err
	}
	blocksElapsed := big.NewInt(0).Sub(latestBlock.Number, latestRetainedBlock.Number)
	if blocksElapsed.Int64() == 0 {
		return nil
	}

	log.Info(blocksElapsed.Int64(), " blocks elapsed since last boot. Backfilling events...")
	startBlockNum := latestRetainedBlock.Number.Int64() + 1
	endBlockNum := latestRetainedBlock.Number.Int64() + blocksElapsed.Int64()
	logs, furthestBlockNumProcessed := w.getLogsInBlockRange(startBlockNum, endBlockNum)
	if furthestBlockNumProcessed > latestRetainedBlock.Number.Int64() {
		// Remove all blocks from the DB
		headers, err := w.InspectRetainedBlocks()
		if err != nil {
			return err
		}
		for i := 0; i < len(headers); i++ {
			_, err := w.stack.Pop()
			if err != nil {
				return err
			}
		}
		// Add furthest block processed into the DB
		latestHeader, err := w.client.HeaderByNumber(big.NewInt(furthestBlockNumProcessed))
		if err != nil {
			return err
		}
		err = w.stack.Push(latestHeader)
		if err != nil {
			return err
		}

		// Emit events for all the logs
		if len(logs) > 0 {
			hashToBlockHeader := map[common.Hash]*meshdb.MiniHeader{}
			for _, log := range logs {
				blockHeader, ok := hashToBlockHeader[log.BlockHash]
				if !ok {
					blockHeader = &meshdb.MiniHeader{
						Hash:   log.BlockHash,
						Number: big.NewInt(0).SetUint64(log.BlockNumber),
						Logs:   []types.Log{},
						// TODO(fabio): What about `Parent`?
					}
					hashToBlockHeader[log.BlockHash] = blockHeader
				}
				blockHeader.Logs = append(blockHeader.Logs, log)
			}
			events := []*Event{}
			for _, blockHeader := range hashToBlockHeader {
				events = append(events, &Event{
					Type:        Added,
					BlockHeader: blockHeader,
				})
			}
			w.blockFeed.Send(events)
		}
	}
	return nil
}

// getLogsInBlockRange attempts to fetch all logs in the block range specified. If it retrieves
// all logs in the range, it simply returns them. If it fails to retrieve some of the blocks,
// it returns all the logs it did find, along with the block number after which no further logs
// were retrieved.
func (w *Watcher) getLogsInBlockRange(from, to int64) ([]types.Log, int64) {
	chunks := w.getBlockRangeChunks(from, to, maxBlocksInGetLogsQuery)

	mu := sync.Mutex{}
	orderIndexTologChunk := map[int][]types.Log{}
	var furthestBlockNumProcessed int64

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for i, chunk := range chunks {
		wg.Add(1)
		go func(index int, chunk *blockRange) {
			defer wg.Done()

			// Add one to the semaphore chan. If it already has concurrencyLimit values,
			// the request blocks here until one frees up.
			semaphoreChan <- struct{}{}

			logs, err := w.filterLogsRecurisively(chunk.FromBlock, chunk.ToBlock, []types.Log{})
			mu.Lock()
			if err != nil {
				log.WithField("error", err).Warn("Failed to fast-sync blocks, falling back to polling them individually")
				if furthestBlockNumProcessed == 0 || furthestBlockNumProcessed < chunk.FromBlock-1 {
					furthestBlockNumProcessed = chunk.FromBlock - 1
				}
			} else {
				if furthestBlockNumProcessed == 0 || furthestBlockNumProcessed < chunk.ToBlock {
					furthestBlockNumProcessed = chunk.ToBlock
				}
			}
			orderIndexTologChunk[index] = logs
			mu.Unlock()
			<-semaphoreChan
		}(i, chunk)
	}

	// Wait for all log requests to complete
	wg.Wait()

	logs := []types.Log{}
	for i := range chunks {
		logs = append(logs, orderIndexTologChunk[i]...)
	}

	return logs, furthestBlockNumProcessed
}

type blockRange struct {
	FromBlock int64
	ToBlock   int64
}

// getBlockRangeChunks breaks up the block range into chunks of chunkSize. `eth_getLogs`
// requests are inclusive to both the start and end blocks specified and so we need to
// make the ranges exclusive of one another to avoid fetching the same blocks' logs twice.
func (w *Watcher) getBlockRangeChunks(from, to, chunkSize int64) []*blockRange {
	chunks := []*blockRange{}
	numBlocksLeft := int64(to - from)
	if numBlocksLeft < chunkSize {
		chunks = append(chunks, &blockRange{
			FromBlock: from,
			ToBlock:   to,
		})
	} else {
		blocks := []int64{}
		for i := int64(0); i <= numBlocksLeft; i++ {
			blocks = append(blocks, from+i)
		}
		numChunks := numBlocksLeft / chunkSize
		remainder := numBlocksLeft % chunkSize
		if remainder > 0 {
			numChunks = numChunks + 1
		}

		for i := int64(0); i < numChunks; i = i + 1 {
			fromIndex := i * chunkSize
			toIndex := fromIndex + chunkSize
			if toIndex >= int64(len(blocks)-1) {
				toIndex = int64(len(blocks))
			}
			bs := blocks[fromIndex:toIndex]
			chunks = append(chunks, &blockRange{
				FromBlock: bs[0],
				ToBlock:   bs[len(bs)-1],
			})
		}
	}
	return chunks
}

func (w *Watcher) filterLogsRecurisively(from, to int64, allLogs []types.Log) ([]types.Log, error) {
	log.Info("Fetch block logs from ", from, " to ", to)
	numBlocks := to - from
	topics := [][]common.Hash{}
	if len(w.topics) > 0 {
		topics = append(topics, w.topics)
	}
	logs, err := w.client.FilterLogs(ethereum.FilterQuery{
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(to),
		Topics:    topics,
	})
	if err != nil {
		// If the request context deadline exceeded, the range was probably too large to send as a single request,
		// split it. Likewise, Infura caps the logs returned to 10,000 per request, if our request exceeds this limit,
		// split it. Parity, Geth and Alchemy all have much higher limits (if any at all), so no need to expect any
		// similar errors of this nature from them.
		if err.Error() == "context deadline exceeded" || err.Error() == "query returned more than 10000 results" {
			// HACK(fabio): Infura limits the returned results to 10,000 logs, BUT some single
			// blocks contain more then 10,000 logs. This has supposedly been fixed but we keep
			// this logic here just in case.
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
				return allLogs, err
			}
			allLogs, err = w.filterLogsRecurisively(startSecondHalf, to, allLogs)
			return allLogs, err
		} else {
			return allLogs, err
		}
	}
	allLogs = append(allLogs, logs...)
	return allLogs, nil
}
