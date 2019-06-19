package blockwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
	rpcClient      *rpc.Client
	client         *ethclient.Client
	requestTimeout time.Duration
}

// NewRpcClient returns a new Client for fetching Ethereum blocks using the given
// ethclient.Client.
func NewRpcClient(rpcURL string, requestTimeout time.Duration) (*RpcClient, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &RpcClient{rpcClient, ethClient, requestTimeout}, nil
}

type GetBlockByNumberResponse struct {
	Hash       common.Hash `json:"hash"`
	ParentHash common.Hash `json:"parentHash"`
	Number     string      `json:"number"`
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByNumber(number *big.Int) (*meshdb.MiniHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rc.requestTimeout)
	defer cancel()

	var blockParam string
	if number == nil {
		blockParam = "latest"
	} else {
		blockParam = fmt.Sprintf("0x%s", common.Bytes2Hex(number.Bytes()))
	}
	shouldIncludeTransactions := false

	// Note(fabio): We use a raw RPC call here instead of `EthClient`'s `BlockByNumber()` method because block
	// hashes are computed differently on Kovan vs. mainnet, resulting in the wrong block hash being returned by
	// `BlockByNumber` when using Kovan. By doing a raw RPC call, we can simply use the blockHash returned in the
	// RPC response rather then re-compute it from the block header.
	// Source: https://github.com/ethereum/go-ethereum/pull/18166
	var header GetBlockByNumberResponse
	err := rc.rpcClient.CallContext(ctx, &header, "eth_getBlockByNumber", blockParam, shouldIncludeTransactions)
	if err != nil {
		return nil, err
	}

	if len(header.Number) < 2 {
		return nil, errors.New("Block number returned by eth_getBlockByNumber too short")
	}
	blockNum, ok := math.ParseBig256(header.Number)
	if !ok {
		return nil, errors.New("Failed to parse big.Int value from hex-encoded block number returned from eth_getBlockByNumber")
	}
	miniHeader := &meshdb.MiniHeader{
		Hash:   header.Hash,
		Parent: header.ParentHash,
		Number: blockNum,
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
