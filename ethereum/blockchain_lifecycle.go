package ethereum

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

// BlockchainLifecycle is a testing utility for taking snapshots of the blockchain
// state on Ganache and reverting those snapshots at a later point in time. Ganache
// supports performing multiple snapshots that can then be reverted in LIFO order.
type BlockchainLifecycle struct {
	rpcClient       *rpc.Client
	snapshotIdStack []string
}

// NewBlockchainLifecycle instantiates a new blockchainLifecycle instance
func NewBlockchainLifecycle(rpcClient *rpc.Client) (*BlockchainLifecycle, error) {
	return &BlockchainLifecycle{
		rpcClient:       rpcClient,
		snapshotIdStack: []string{},
	}, nil
}

// Start creates a snapshot of the blockchain state at that point in time
// and adds it's snapshotId to a stack
func (b *BlockchainLifecycle) Start(t *testing.T) {
	var snapshotId string
	err := b.rpcClient.Call(&snapshotId, "evm_snapshot")
	require.NoError(t, err)
	b.snapshotIdStack = append(b.snapshotIdStack, snapshotId)
}

// Revert reverts the latest snapshot of blockchain state created
func (b *BlockchainLifecycle) Revert(t *testing.T) {
	latestSnapshot := b.snapshotIdStack[len(b.snapshotIdStack)-1]
	b.snapshotIdStack = b.snapshotIdStack[:len(b.snapshotIdStack)-1]
	var didRevert bool
	err := b.rpcClient.Call(&didRevert, "evm_revert", latestSnapshot)
	require.NoError(t, err)
	if !didRevert {
		t.Errorf("Failed to revert snapshot with ID: %s", latestSnapshot)
	}
}

// Mine force-mines a block with the specified block timestamp
func (b *BlockchainLifecycle) Mine(t *testing.T, blockTimestamp time.Time) {
	var didForceMine string
	err := b.rpcClient.Call(&didForceMine, "evm_mine", blockTimestamp.Unix())
	require.NoError(t, err)
	if didForceMine != "0x0" {
		t.Error("Failed to force mine a block")
	}
}
