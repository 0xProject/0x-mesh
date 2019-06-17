package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

type BlockchainLifecycle struct {
	rpcClient       *rpc.Client
	snapshotIdStack []string
}

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

func (b *BlockchainLifecycle) Start() {
	var snapshotId string
	b.rpcClient.Call(&snapshotId, "evm_snapshot")
	b.snapshotIdStack = append(b.snapshotIdStack, snapshotId)
}

func (b *BlockchainLifecycle) Revert() error {
	latestSnapshot := b.snapshotIdStack[len(b.snapshotIdStack)-1]
	b.snapshotIdStack = b.snapshotIdStack[:len(b.snapshotIdStack)-1]
	var didRevert bool
	b.rpcClient.Call(didRevert, "evm_revert", latestSnapshot)
	if !didRevert {
		return fmt.Errorf("Failed to revert snapshot with ID: %s", latestSnapshot)
	}
	return nil
}
