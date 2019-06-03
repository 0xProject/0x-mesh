package blockwatch

import (
	"context"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client defines the methods needed to satisfy the client expected when
// instantiating a Watcher instance.
type Client interface {
	HeaderByNumber(number *big.Int) (*meshdb.MiniHeader, error)
	HeaderByHash(hash common.Hash) (*meshdb.MiniHeader, error)
	FilterLogs(q ethereum.FilterQuery) ([]types.Log, error)
}

// RpcClient is a Client for fetching Ethereum blocks from a specific JSON-RPC endpoint.
type RpcClient struct {
	client         *ethclient.Client
	requestTimeout time.Duration
}

// NewRpcClient returns a new Client for fetching Ethereum blocks using the given
// ethclient.Client.
func NewRpcClient(ethClient *ethclient.Client, requestTimeout time.Duration) (*RpcClient, error) {
	return &RpcClient{ethClient, requestTimeout}, nil
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByNumber(number *big.Int) (*meshdb.MiniHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rc.requestTimeout)
	defer cancel()
	header, err := rc.client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	miniHeader := &meshdb.MiniHeader{
		Hash:   header.Hash(),
		Parent: header.ParentHash,
		Number: header.Number,
	}
	return miniHeader, nil
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByHash(hash common.Hash) (*meshdb.MiniHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rc.requestTimeout)
	defer cancel()
	header, err := rc.client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	miniHeader := &meshdb.MiniHeader{
		Hash:   header.Hash(),
		Parent: header.ParentHash,
		Number: header.Number,
	}
	return miniHeader, nil
}

// FilterLogs returns the logs that satisfy the supplied filter query.
func (rc *RpcClient) FilterLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rc.requestTimeout)
	defer cancel()
	logs, err := rc.client.FilterLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
