package blockwatch

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// MiniBlockHeader is a more succinct block representation then the one returned by go-ethereum.
// It contains all the information necessary to implement BlockWatch.
type MiniBlockHeader struct {
	Hash   common.Hash `json:"hash"   gencodec:"required"`
	Parent common.Hash `json:"parent" gencodec:"required"`
	Number *big.Int    `json:"number" gencodec:"required"`
}

// NewSuccintBlock returns a new MiniBlockHeader.
func NewSuccintBlock(hash common.Hash, parent common.Hash, number *big.Int) *MiniBlockHeader {
	succintBlock := MiniBlockHeader{Hash: hash, Parent: parent, Number: number}
	return &succintBlock
}

// BlockEvent describes a block event emitted by a BlockWatch
type BlockEvent struct {
	WasRemoved bool
	Block      *MiniBlockHeader
}

// BlockWatch maintains a consistent representation of the latest `BlockRetentionLimit` blocks,
// handling block re-orgs and network disruptions gracefully. It can be started from any arbitrary
// block height, and will emit both block added and removed events.
type BlockWatch struct {
	startBlockDepth     rpc.BlockNumber
	BlockRetentionLimit uint
	blockStack          *BlockStack
	client              BlockClient
	Events              chan *BlockEvent
	Errors              chan error
	ticker              *time.Ticker
}

// NewBlockWatch creates a new BlockWatch instance.
func NewBlockWatch(startBlockDepth rpc.BlockNumber, blockRetentionLimit uint, client BlockClient) *BlockWatch {
	blockStack := &BlockStack{}
	events := make(chan *BlockEvent)
	errors := make(chan error)
	bs := &BlockWatch{startBlockDepth, blockRetentionLimit, blockStack, client, events, errors, nil}
	return bs
}

// StartPolling starts the BlockWatch block poller.
func (bs *BlockWatch) StartPolling(ctx context.Context, pollingInterval time.Duration) {
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

// StopPolling stops the BlockWatch block poller.
func (bs *BlockWatch) StopPolling() {
	bs.ticker.Stop()
}

// PollNextBlock lets you manually poll for the next block to be added to the block stack.
// If there are no blocks on the stack, it fetches the first block at the specified
// `startBlockDepth` supplied at instantiation.
func (bs *BlockWatch) PollNextBlock(ctx context.Context) error {
	var nextBlockNumber *big.Int
	latestBlock := bs.blockStack.Peek()
	if latestBlock == nil {
		if bs.startBlockDepth == rpc.LatestBlockNumber {
			nextBlockNumber = nil
		} else {
			nextBlockNumber = big.NewInt(int64(bs.startBlockDepth))
		}
	} else {
		nextBlockNumber = big.NewInt(0).Add(latestBlock.Number, big.NewInt(1))
	}
	nextBlock, err := bs.client.HeaderByNumber(ctx, nextBlockNumber)
	if err != nil {
		if err == ethereum.NotFound {
			return nil // Noop and wait next polling interval
		}
		return err
	}

	return bs.buildCanonicalChain(ctx, nextBlock)
}

func (bs *BlockWatch) buildCanonicalChain(ctx context.Context, nextBlock *MiniBlockHeader) error {
	latestBlock := bs.blockStack.Peek()
	// Is the blockStack empty or is it the next block?
	if latestBlock == nil || nextBlock.Parent == latestBlock.Hash {
		bs.blockStack.Push(nextBlock)
		bs.Events <- &BlockEvent{
			WasRemoved: false,
			Block:      nextBlock,
		}
		return nil
	}

	poppedBlock := bs.blockStack.Pop()
	bs.Events <- &BlockEvent{
		WasRemoved: true,
		Block:      poppedBlock,
	}

	nextBlockParent, err := bs.client.HeaderByHash(ctx, nextBlock.Parent)
	if err != nil {
		if err == ethereum.NotFound {
			// Noop and wait next polling interval. We remove the popped blocks
			// and refetch them on the next polling interval.
			return nil
		}
		return err
	}
	bs.buildCanonicalChain(ctx, nextBlockParent)
	bs.blockStack.Push(nextBlock)
	bs.Events <- &BlockEvent{
		WasRemoved: false,
		Block:      nextBlock,
	}

	return nil
}

// GetRetainedBlocks returns the blocks retained in-memory by the BlockWatch instance
func (bs *BlockWatch) GetRetainedBlocks() []*MiniBlockHeader {
	return bs.blockStack.data
}
