package ethereum

import (
	"context"
	"testing"

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
func NewBlockchainLifecycle(rpcURL string) (*BlockchainLifecycle, error) {
	rpcClient, err := rpc.DialContext(context.Background(), rpcURL)
	if err != nil {
		return nil, err
	}
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
