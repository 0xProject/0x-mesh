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
	HeaderByNumber(ctx context.Context, number *big.Int) (*MiniBlockHeader, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*MiniBlockHeader, error)
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

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (rbc *RpcBlockClient) HeaderByNumber(ctx context.Context, number *big.Int) (*MiniBlockHeader, error) {
	header, err := rbc.client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	miniBlockHeader := NewMiniBlockHeader(header.Hash(), header.ParentHash, header.Number)
	return miniBlockHeader, nil
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (rbc *RpcBlockClient) HeaderByHash(ctx context.Context, hash common.Hash) (*MiniBlockHeader, error) {
	header, err := rbc.client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	miniBlockHeader := NewMiniBlockHeader(header.Hash(), header.ParentHash, header.Number)
	return miniBlockHeader, nil
}
