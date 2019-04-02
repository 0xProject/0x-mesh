package blockwatch

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
)

// MiniHeader is a more succinct block header representation then the one returned by go-ethereum.
// It contains all the information necessary to implement Watcher.
type MiniHeader struct {
	Hash   common.Hash `json:"hash"   gencodec:"required"`
	Parent common.Hash `json:"parent" gencodec:"required"`
	Number *big.Int    `json:"number" gencodec:"required"`
	Logs   []types.Log `json:"logs" gencodec:"required"`
}

// NewMiniHeader returns a new MiniHeader.
func NewMiniHeader(hash common.Hash, parent common.Hash, number *big.Int) *MiniHeader {
	miniHeader := MiniHeader{Hash: hash, Parent: parent, Number: number, Logs: []types.Log{}}
	return &miniHeader
}

// EventType describes the types of events emitted by blockwatch.Watcher. A block can be discovered
// and added to our representation of the chain. During a block re-org, a block previously stored
// can be removed from the list. Lastly, if more then blockRetentionLimit blocks have been discovered,
// the oldest block stored will be retired (e.g., no longer tracked, but still considered part of the
// canonical chain).
type EventType int

const (
	Added EventType = iota
	Removed
	Retired
)

// Event describes a block event emitted by a Watcher
type Event struct {
	Type        EventType
	BlockHeader *MiniHeader
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

	mu sync.RWMutex
}

// New creates a new Watcher instance.
func New(pollingInterval time.Duration, startBlockDepth rpc.BlockNumber, blockRetentionLimit int, withLogs bool, client Client) *Watcher {
	stack := NewStack(blockRetentionLimit)
	// Buffer the first 5 errors, if no channel consumer processing the errors, any additional errors are dropped
	errorsChan := make(chan error, 5)
	bs := &Watcher{
		Errors:              errorsChan,
		pollingInterval:     pollingInterval,
		blockRetentionLimit: blockRetentionLimit,
		startBlockDepth:     startBlockDepth,
		stack:               stack,
		client:              client,
		withLogs:            withLogs,
	}
	return bs
}

// Subscribe allows one to subscribe to the block events emitted by the Watcher.
// As soon as the first subscription is registered, the block poller is started.
// To unsubscribe, simply call `Unsubscribe` on the returned subscription. When all
// consumers have unsubscribed, the block polling stops.
// The sink channel should have ample buffer space to avoid blocking other subscribers.
// Slow subscribers are not dropped.
func (w *Watcher) Subscribe(sink chan<- []*Event) event.Subscription {
	// We need the mutex to reliably start/stop the update loop
	w.mu.Lock()
	defer w.mu.Unlock()

	sub := w.blockScope.Track(w.blockFeed.Subscribe(sink))

	if !w.isWatching {
		w.isWatching = true
		w.ticker = time.NewTicker(w.pollingInterval)
		go w.startPolling()
	}

	return sub
}

// InspectRetainedBlocks returns the blocks retained in-memory by the Watcher instance. It is not
// particularly performant and therefore should only be used for debugging and testing purposes.
func (w *Watcher) InspectRetainedBlocks() []*MiniHeader {
	return w.stack.Inspect()
}

func (w *Watcher) startPolling() {
	for {
		<-w.ticker.C

		w.mu.Lock()
		if w.blockScope.Count() == 0 {
			w.isWatching = false
			w.ticker.Stop()
			w.mu.Unlock()
			return
		}
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

// pollNextBlock polls for the next block header to be added to the block stack.
// If there are no blocks on the stack, it fetches the first block at the specified
// `startBlockDepth` supplied at instantiation.
func (w *Watcher) pollNextBlock() error {
	var nextBlockNumber *big.Int
	latestHeader := w.stack.Peek()
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

func (w *Watcher) buildCanonicalChain(nextHeader *MiniHeader, events []*Event) ([]*Event, error) {
	latestHeader := w.stack.Peek()
	// Is the stack empty or is it the next block?
	if latestHeader == nil || nextHeader.Parent == latestHeader.Hash {
		nextHeader, err := w.addLogs(nextHeader)
		if err != nil {
			return events, err
		}
		retiredBlock := w.stack.Push(nextHeader)
		events = append(events, &Event{
			Type:        Added,
			BlockHeader: nextHeader,
		})
		if retiredBlock != nil {
			events = append(events, &Event{
				Type:        Retired,
				BlockHeader: retiredBlock,
			})
		}
		return events, nil
	}

	w.stack.Pop() // Pop latestHeader from the stack. We already have a reference to it.
	events = append(events, &Event{
		Type:        Removed,
		BlockHeader: latestHeader,
	})

	nextParentHeader, err := w.client.HeaderByHash(nextHeader.Parent)
	if err != nil {
		if err == ethereum.NotFound {
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
		return events, err
	}
	retiredBlock := w.stack.Push(nextHeader)
	events = append(events, &Event{
		Type:        Added,
		BlockHeader: nextHeader,
	})
	if retiredBlock != nil {
		events = append(events, &Event{
			Type:        Retired,
			BlockHeader: retiredBlock,
		})
	}

	return events, nil
}

func (w *Watcher) addLogs(header *MiniHeader) (*MiniHeader, error) {
	if !w.withLogs {
		return header, nil
	}
	logs, err := w.client.FilterLogs(ethereum.FilterQuery{
		BlockHash: &header.Hash,
		// TODO(fabio): Add topics (hash of event signatures we care about)
	})
	if err != nil {
		return nil, err
	}
	header.Logs = logs
	return header, nil
}
