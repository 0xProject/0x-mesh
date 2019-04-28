package blockwatch

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

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
	BlockHeader *meshdb.MiniHeader
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
func New(meshDB *meshdb.MeshDB, pollingInterval time.Duration, startBlockDepth rpc.BlockNumber, blockRetentionLimit int, withLogs bool, topics []common.Hash, client Client) *Watcher {
	stack := NewStack(meshDB, blockRetentionLimit)

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
		topics:              topics,
	}
	return bs
}

// StartPolling starts the block poller
func (w *Watcher) StartPolling() error {
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

// StopPolling stops the block poller
func (w *Watcher) StopPolling() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isWatching = false
	if w.ticker != nil {
		w.ticker.Stop()
	}
	w.ticker = nil
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
func (w *Watcher) InspectRetainedBlocks() []*meshdb.MiniHeader {
	return w.stack.Inspect()
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
			log.WithFields(log.Fields{
				"blockNumber": nextBlockNumber,
			}).Info("block header not found")
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
	latestHeader := w.stack.Peek()
	// Is the stack empty or is it the next block?
	if latestHeader == nil || nextHeader.Parent == latestHeader.Hash {
		nextHeader, err := w.addLogs(nextHeader)
		if err != nil {
			return events, err
		}
		retiredBlock, err := w.stack.Push(nextHeader)
		if err != nil {
			return events, err
		}
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

	_, err := w.stack.Pop() // Pop latestHeader from the stack. We already have a reference to it.
	if err != nil {
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
		return events, err
	}
	retiredBlock, err := w.stack.Push(nextHeader)
	if err != nil {
		return events, err
	}
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

func (w *Watcher) addLogs(header *meshdb.MiniHeader) (*meshdb.MiniHeader, error) {
	if !w.withLogs {
		return header, nil
	}
	logs, err := w.client.FilterLogs(ethereum.FilterQuery{
		BlockHash: &header.Hash,
		Topics:    [][]common.Hash{w.topics},
	})
	if err != nil {
		return nil, err
	}
	header.Logs = logs
	return header, nil
}
