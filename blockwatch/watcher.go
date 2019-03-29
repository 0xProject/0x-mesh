package blockwatch

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// MiniBlockHeader is a more succinct block header representation then the one returned by go-ethereum.
// It contains all the information necessary to implement Watcher.
type MiniBlockHeader struct {
	Hash   common.Hash `json:"hash"   gencodec:"required"`
	Parent common.Hash `json:"parent" gencodec:"required"`
	Number *big.Int    `json:"number" gencodec:"required"`
}

// NewMiniBlockHeader returns a new MiniBlockHeader.
func NewMiniBlockHeader(hash common.Hash, parent common.Hash, number *big.Int) *MiniBlockHeader {
	miniBlockHeader := MiniBlockHeader{Hash: hash, Parent: parent, Number: number}
	return &miniBlockHeader
}

// Event describes a block event emitted by a Watcher
type Event struct {
	WasRemoved  bool
	BlockHeader *MiniBlockHeader
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
}

// New creates a new Watcher instance.
func New(startBlockDepth rpc.BlockNumber, blockRetentionLimit uint, client Client) *Watcher {
	stack := NewStack()
	events := make(chan []*Event)
	errors := make(chan error)
	bs := &Watcher{startBlockDepth, blockRetentionLimit, stack, client, events, errors, nil}
	return bs
}

// StartPolling starts the Watcher block poller.
func (bs *Watcher) StartPolling(ctx context.Context, pollingInterval time.Duration) {
	bs.ticker = time.NewTicker(pollingInterval)
	go func() {
		for _ = range bs.ticker.C {
			err := bs.PollNextBlock(ctx)
			if err != nil {
				bs.Errors <- err
			}
		}
	}()
}

// StopPolling stops the Watcher block poller.
func (bs *Watcher) StopPolling() {
	bs.ticker.Stop()
}

// PollNextBlock lets you manually poll for the next block header to be added to the block stack.
// If there are no blocks on the stack, it fetches the first block at the specified
// `startBlockDepth` supplied at instantiation.
func (bs *Watcher) PollNextBlock(ctx context.Context) error {
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
	nextHeader, err := bs.client.HeaderByNumber(ctx, nextBlockNumber)
	if err != nil {
		if err == ethereum.NotFound {
			return nil // Noop and wait next polling interval
		}
		return err
	}

	events := []*Event{}
	events, err = bs.buildCanonicalChain(ctx, nextHeader, events)
	// Even if an error occurred, we still want to emit the events gathered since we might have
	// popped blocks off the Stack and they won't be re-added
	if len(events) != 0 {
		go func() {
			bs.Events <- events
		}()
	}
	if err != nil {
		return err
	}
	return nil
}

func (bs *Watcher) buildCanonicalChain(ctx context.Context, nextHeader *MiniBlockHeader, events []*Event) ([]*Event, error) {
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

	nextParentHeader, err := bs.client.HeaderByHash(ctx, nextHeader.Parent)
	if err != nil {
		if err == ethereum.NotFound {
			// Noop and wait next polling interval. We remove the popped blocks
			// and refetch them on the next polling interval.
			return events, nil
		}
		return events, err
	}
	events, err = bs.buildCanonicalChain(ctx, nextParentHeader, events)
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
func (bs *Watcher) InspectRetainedBlocks() []*MiniBlockHeader {
	return bs.stack.Inspect()
}
