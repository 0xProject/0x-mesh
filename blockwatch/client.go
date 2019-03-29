package blockwatch

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client defines the methods needed to satisfy the client expected when
// instantiating a Watcher instance.
type Client interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*MiniBlockHeader, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*MiniBlockHeader, error)
}

// RpcClient is a Client for fetching Ethereum blocks from a specific JSON-RPC endpoint.
type RpcClient struct {
	client *ethclient.Client
}

// NewRpcClient returns a new Client for fetching Ethereum blocks from a supplied JSON-RPC endpoint.
func NewRpcClient(rpcURL string) (*RpcClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &RpcClient{client}, nil
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByNumber(ctx context.Context, number *big.Int) (*MiniBlockHeader, error) {
	header, err := rc.client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	miniBlockHeader := NewMiniBlockHeader(header.Hash(), header.ParentHash, header.Number)
	return miniBlockHeader, nil
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByHash(ctx context.Context, hash common.Hash) (*MiniBlockHeader, error) {
	header, err := rc.client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	miniBlockHeader := NewMiniBlockHeader(header.Hash(), header.ParentHash, header.Number)
	return miniBlockHeader, nil
}
