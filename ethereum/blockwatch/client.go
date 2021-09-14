package blockwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/constants"
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
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Number     string `json:"number"`
	Timestamp  string `json:"timestamp"`
}

// UnknownBlockNumberError is the error returned from a filter logs RPC call when the block number
// specified is not recognized
type UnknownBlockNumberError struct {
	Message     string
	BlockNumber *big.Int
}

func (e UnknownBlockNumberError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.BlockNumber)
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
	// rc.ethRPCClient.HeaderByNumber(ctx context.Context, number *big.Int)
	err := rc.ethRPCClient.CallContext(ctx, &header, "eth_getBlockByNumber", blockParam, shouldIncludeTransactions)
	if err != nil {
		return nil, err
	}
	// If it returned an empty struct
	if header.Number == "" {
		// Add blockHash to error so it gets logged
		return nil, UnknownBlockNumberError{
			Message:     ethereum.NotFound.Error(),
			BlockNumber: number,
		}
	}

	blockNum, ok := math.ParseBig256(header.Number)
	if !ok {
		return nil, errors.New("Failed to parse big.Int value from hex-encoded block number returned from eth_getBlockByNumber")
	}
	unixTimestamp, ok := math.ParseBig256(header.Timestamp)
	if !ok {
		return nil, errors.New("Failed to parse big.Int value from hex-encoded block timestamp returned from eth_getBlockByNumber")
	}
	miniHeader := &miniheader.MiniHeader{
		Hash:      common.HexToHash(header.Hash),
		Parent:    common.HexToHash(header.ParentHash),
		Number:    blockNum,
		Timestamp: time.Unix(unixTimestamp.Int64(), 0),
	}
	return miniHeader, nil
}

// UnknownBlockHashError is the error returned from a filter logs RPC call when the blockHash
// specified is not recognized
type UnknownBlockHashError struct {
	BlockHash common.Hash
}

func (e UnknownBlockHashError) Error() string {
	return fmt.Sprintf("%s: %s", ethereum.NotFound.Error(), e.BlockHash)
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (rc *RpcClient) HeaderByHash(hash common.Hash) (*miniheader.MiniHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	header, err := rc.ethRPCClient.HeaderByHash(ctx, hash)
	if err != nil {
		// Add blockHash to error so it gets logged
		if err.Error() == ethereum.NotFound.Error() {
			err = UnknownBlockHashError{
				BlockHash: hash,
			}
		}
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

// FilterUnknownBlockError is the error returned from a filter logs RPC call when the blockHash
// specified is not recognized
type FilterUnknownBlockError struct {
	Message     string
	FilterQuery ethereum.FilterQuery
}

func (e FilterUnknownBlockError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.FilterQuery.BlockHash.String())
}

// FilterLogs returns the logs that satisfy the supplied filter query.
func (rc *RpcClient) FilterLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	logs, err := rc.ethRPCClient.FilterLogs(ctx, q)
	if err != nil {
		// Add the query filter to the error so that it gets logged
		if err.Error() == constants.ParityFilterUnknownBlock || err.Error() == constants.GethFilterUnknownBlock {
			err = FilterUnknownBlockError{
				Message:     err.Error(),
				FilterQuery: q,
			}
		}
		return nil, err
	}
	return logs, nil
}
