package blockwatch

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// We give up on ETH RPC requests sent for the purpose of block watching after 10 seconds
	requestTimeout = 10 * time.Second
)

// Client defines the methods needed to satisfy the client expected when
// instantiating a Watcher instance.
type Client interface {
	HeaderByNumber(number *big.Int) (*miniheader.MiniHeader, error)
	HeaderByHash(hash common.Hash) (*miniheader.MiniHeader, error)
	FilterLogs(q ethereum.FilterQuery) ([]types.Log, error)
}

// RpcClient is a Client for fetching Ethereum blocks from a specific JSON-RPC endpoint.
type RpcClient struct {
	ethRPCClient ethrpcclient.Client
}

// NewRpcClient returns a new Client for fetching Ethereum blocks using the given
// ethclient.Client.
func NewRpcClient(ethRPCClient ethrpcclient.Client) (*RpcClient, error) {
	return &RpcClient{
		ethRPCClient: ethRPCClient,
	}, nil
}

type GetBlockByNumberResponse struct {
	Hash       common.Hash `json:"hash"`
	ParentHash common.Hash `json:"parentHash"`
	Number     string      `json:"number"`
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByNumber(number *big.Int) (*miniheader.MiniHeader, error) {
	var blockParam string
	if number == nil {
		blockParam = "latest"
	} else {
		blockParam = hexutil.EncodeBig(number)
	}
	shouldIncludeTransactions := false

	// Note(fabio): We use a raw RPC call here instead of `EthClient`'s `BlockByNumber()` method because block
	// hashes are computed differently on Kovan vs. mainnet, resulting in the wrong block hash being returned by
	// `BlockByNumber` when using Kovan. By doing a raw RPC call, we can simply use the blockHash returned in the
	// RPC response rather than re-compute it from the block header.
	// Source: https://github.com/ethereum/go-ethereum/pull/18166
	var header GetBlockByNumberResponse
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	err := rc.ethRPCClient.CallContext(ctx, &header, "eth_getBlockByNumber", blockParam, shouldIncludeTransactions)
	if err != nil {
		return nil, err
	}
	// If it returned an empty struct
	if header.Number == "" {
		return nil, ethereum.NotFound
	}

	blockNum, ok := math.ParseBig256(header.Number)
	if !ok {
		return nil, errors.New("Failed to parse big.Int value from hex-encoded block number returned from eth_getBlockByNumber")
	}
	miniHeader := &miniheader.MiniHeader{
		Hash:   header.Hash,
		Parent: header.ParentHash,
		Number: blockNum,
	}
	return miniHeader, nil
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByHash(hash common.Hash) (*miniheader.MiniHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	header, err := rc.ethRPCClient.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	miniHeader := &miniheader.MiniHeader{
		Hash:      header.Hash(),
		Parent:    header.ParentHash,
		Number:    header.Number,
		Timestamp: time.Unix(int64(header.Time), 0),
	}
	return miniHeader, nil
}

// FilterLogs returns the logs that satisfy the supplied filter query.
func (rc *RpcClient) FilterLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	logs, err := rc.ethRPCClient.FilterLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
