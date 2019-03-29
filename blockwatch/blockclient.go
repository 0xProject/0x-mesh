package blockwatch

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockClient defines the methods needed to satisfy the blockClient expected when
// instantiating a BlockWatch instance.
type BlockClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*SuccinctBlock, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*SuccinctBlock, error)
}

// RpcBlockClient is a Client for fetching Ethereum blocks from a specific JSON-RPC endpoint.
type RpcBlockClient struct {
	client *ethclient.Client
}

// NewRpcBlockClient returns a new Client for fetching Ethereum blocks from a supplied JSON-RPC endpoint.
func NewRpcBlockClient(rpcURL string) (*RpcBlockClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &RpcBlockClient{client}, nil
}

// BlockByNumber fetches a block by its number. If no `number` is supplied, it will return the latest block.
// If not block exists with this number it will return a `ethereum.NotFound` error.
func (rbc *RpcBlockClient) BlockByNumber(ctx context.Context, number *big.Int) (*SuccinctBlock, error) {
	block, err := rbc.client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	succintBlock := NewSuccintBlock(block.Hash(), block.ParentHash(), block.Number())
	return succintBlock, nil
}

// BlockByHash fetches a block by its block hash. If not block exists with this number it will return a `ethereum.NotFound` error.
func (rbc *RpcBlockClient) BlockByHash(ctx context.Context, hash common.Hash) (*SuccinctBlock, error) {
	block, err := rbc.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	succintBlock := NewSuccintBlock(block.Hash(), block.ParentHash(), block.Number())
	return succintBlock, nil
}
