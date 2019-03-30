package blockwatch

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// MiniHeader is a more succinct block header representation then the one returned by go-ethereum.
// It contains all the information necessary to implement Watcher.
type MiniHeader struct {
	Hash   common.Hash `json:"hash"   gencodec:"required"`
	Parent common.Hash `json:"parent" gencodec:"required"`
	Number *big.Int    `json:"number" gencodec:"required"`
}

// NewMiniHeader returns a new MiniHeader.
func NewMiniHeader(hash common.Hash, parent common.Hash, number *big.Int) *MiniHeader {
	miniHeader := MiniHeader{Hash: hash, Parent: parent, Number: number}
	return &miniHeader
}

// Event describes a block event emitted by a Watcher
type Event struct {
	WasRemoved  bool
	BlockHeader *MiniHeader
}

// Watcher maintains a consistent representation of the latest `BlockRetentionLimit` blocks,
// handling block re-orgs and network disruptions gracefully. It can be started from any arbitrary
// block height, and will emit both block added and removed events.
type Watcher struct {
	startBlockDepth     rpc.BlockNumber
	BlockRetentionLimit uint
	stack               *Stack
	client              Client
	Events              chan []*Event
	Errors              chan error
	ticker              *time.Ticker
	tickerCancelFunc    context.CancelFunc
	tickerMut           sync.Mutex
}

// New creates a new Watcher instance.
func New(startBlockDepth rpc.BlockNumber, blockRetentionLimit uint, client Client) *Watcher {
	stack := NewStack()
	events := make(chan []*Event)
	errors := make(chan error)
	bs := &Watcher{startBlockDepth: startBlockDepth, BlockRetentionLimit: blockRetentionLimit, stack: stack, client: client, Events: events, Errors: errors, ticker: nil, tickerCancelFunc: nil}
	return bs
}

// StartPolling starts the Watcher block poller.
func (bs *Watcher) StartPolling(pollingInterval time.Duration) error {
	bs.tickerMut.Lock()
	defer bs.tickerMut.Unlock()
	if bs.ticker != nil {
		return errors.New("cannot start polling more than once")
	}
	bs.ticker = time.NewTicker(pollingInterval)
	ctx, cancelFn := context.WithCancel(context.Background())
	bs.tickerCancelFunc = cancelFn
	go func() {
		for {
			select {
			case <-bs.ticker.C:
				err := bs.PollNextBlock()
				if err != nil {
					bs.Errors <- err
				}
			case <-ctx.Done():
				// The context was cancelled or timed out, etc. End the goroutine by returning.
				return
			}
		}
	}()
	return nil
}

// StopPolling stops the Watcher block poller.
func (bs *Watcher) StopPolling() {
	bs.tickerMut.Lock()
	defer bs.tickerMut.Unlock()
	bs.ticker.Stop()
	bs.tickerCancelFunc()
}

// PollNextBlock lets you manually poll for the next block header to be added to the block stack.
// If there are no blocks on the stack, it fetches the first block at the specified
// `startBlockDepth` supplied at instantiation.
func (bs *Watcher) PollNextBlock() error {
	var nextBlockNumber *big.Int
	latestHeader := bs.stack.Peek()
	if latestHeader == nil {
		if bs.startBlockDepth == rpc.LatestBlockNumber {
			nextBlockNumber = nil
		} else {
			nextBlockNumber = big.NewInt(int64(bs.startBlockDepth))
		}
	} else {
		nextBlockNumber = big.NewInt(0).Add(latestHeader.Number, big.NewInt(1))
	}
	nextHeader, err := bs.client.HeaderByNumber(nextBlockNumber)
	if err != nil {
		if err == ethereum.NotFound {
			return nil // Noop and wait next polling interval
		}
		return err
	}

	events := []*Event{}
	events, err = bs.buildCanonicalChain(nextHeader, events)
	// Even if an error occurred, we still want to emit the events gathered since we might have
	// popped blocks off the Stack and they won't be re-added
	if len(events) != 0 {
		// TODO(fabio): This could be a memory leak if the channel consumer does not receive all
		// events from the bs.Events channel. To fix this, we would need to `select` on sending
		// to the Events channel and returning if a value is received on chan `ctx.Done()`
		// BUG(fabio): Although channels preserve the ordering with which values are sent into the
		// channel, go-routines make no guarentees about when they are run. Since we use blocking
		// channels, and use go routines to queue values, we are not guarenteed they will be sent in
		// the order produced. This means we are guarenteed to receive the events of a particular polling
		// interval in the correct order, but have no ordering guarentees on the events we receive from
		// multiple polling intervals. This "bug" does not cause issues for our current use-cases of BlockWatch.
		go func() {
			bs.Events <- events
		}()
	}
	if err != nil {
		return err
	}
	return nil
}

func (bs *Watcher) buildCanonicalChain(nextHeader *MiniHeader, events []*Event) ([]*Event, error) {
	latestHeader := bs.stack.Peek()
	// Is the stack empty or is it the next block?
	if latestHeader == nil || nextHeader.Parent == latestHeader.Hash {
		bs.stack.Push(nextHeader)
		events = append(events, &Event{
			WasRemoved:  false,
			BlockHeader: nextHeader,
		})
		return events, nil
	}

	poppedBlockHeader := bs.stack.Pop()
	events = append(events, &Event{
		WasRemoved:  true,
		BlockHeader: poppedBlockHeader,
	})

	nextParentHeader, err := bs.client.HeaderByHash(nextHeader.Parent)
	if err != nil {
		if err == ethereum.NotFound {
			// Noop and wait next polling interval. We remove the popped blocks
			// and refetch them on the next polling interval.
			return events, nil
		}
		return events, err
	}
	events, err = bs.buildCanonicalChain(nextParentHeader, events)
	if err != nil {
		return events, err
	}
	bs.stack.Push(nextHeader)
	events = append(events, &Event{
		WasRemoved:  false,
		BlockHeader: nextHeader,
	})

	return events, nil
}

// InspectRetainedBlocks returns the blocks retained in-memory by the Watcher instance. It is not
// particularly performant and therefore should only be used for debugging and testing purposes.
func (bs *Watcher) InspectRetainedBlocks() []*MiniHeader {
	return bs.stack.Inspect()
}
