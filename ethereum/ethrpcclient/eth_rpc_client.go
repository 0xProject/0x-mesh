package ethrpcclient

import (
	"context"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client defines the methods needed to satisfy the subsdet of ETH JSON-RPC client
// methods used by Mesh
type Client interface {
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*miniheader.MiniHeader, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	GetRateLimitDroppedRequests() int64
}

// client is a Client through which _all_ Ethereum JSON-RPC requests should be routed through. It
// enforces a max requestTimeout and also rate-limits requests
type client struct {
	// rpcClient is the underlying RPC client or provider
	rpcClient ethclient.RPCClient
	// client is the higher level Ethereum RPC client with lots of helper methods
	// for converting Go types to JSON and vice versa.
	client         *ethclient.Client
	requestTimeout time.Duration
	rateLimiter    ratelimit.RateLimiter
	// rateLimitDroppedRequests counts the number of requests that had their context cancelled or expire
	// and were therefore never granted
	rateLimitDroppedRequests int64
}

// New returns a new instance of client
func New(rpcClient ethclient.RPCClient, requestTimeout time.Duration, rateLimiter ratelimit.RateLimiter) (Client, error) {
	ethClient := ethclient.NewClient(rpcClient)
	return &client{
		client:         ethClient,
		rpcClient:      rpcClient,
		requestTimeout: requestTimeout,
		rateLimiter:    rateLimiter,
	}, nil
}

// CallContext performs a JSON-RPC call with the given arguments. If the context is
// canceled before the call has successfully returned, CallContext returns immediately.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
func (ec *client) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	err := ec.rateLimiter.Wait(ctx)
	if err != nil {
		atomic.AddInt64(&ec.rateLimitDroppedRequests, 1)
		// Context cancelled or deadline exceeded
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer cancel()
	return ec.rpcClient.CallContext(ctx, &result, method, args...)
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (ec *client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	err := ec.rateLimiter.Wait(ctx)
	if err != nil {
		atomic.AddInt64(&ec.rateLimitDroppedRequests, 1)
		// Context cancelled or deadline exceeded
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer cancel()
	header, err := ec.client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (ec *client) HeaderByNumber(ctx context.Context, number *big.Int) (*miniheader.MiniHeader, error) {
	err := ec.rateLimiter.Wait(ctx)
	if err != nil {
		atomic.AddInt64(&ec.rateLimitDroppedRequests, 1)
		// Context cancelled or deadline exceeded
		return nil, err
	}

	header, err := ec.client.HeaderByNumber(ctx, number)
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

// CodeAt returns the code of the given account. This is needed to differentiate
// between contract internal errors and the local chain being out of sync.
func (ec *client) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	err := ec.rateLimiter.Wait(ctx)
	if err != nil {
		atomic.AddInt64(&ec.rateLimitDroppedRequests, 1)
		// Context cancelled or deadline exceeded
		return []byte{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer cancel()
	return ec.client.CodeAt(ctx, contract, blockNumber)
}

// CallContract executes an Ethereum contract call with the specified data as the input.
func (ec *client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	err := ec.rateLimiter.Wait(ctx)
	if err != nil {
		atomic.AddInt64(&ec.rateLimitDroppedRequests, 1)
		// Context cancelled or deadline exceeded
		return []byte{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer cancel()
	return ec.client.CallContract(ctx, call, blockNumber)
}

// FilterLogs returns the logs that satisfy the supplied filter query.
func (ec *client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	err := ec.rateLimiter.Wait(ctx)
	if err != nil {
		atomic.AddInt64(&ec.rateLimitDroppedRequests, 1)
		// Context cancelled or deadline exceeded
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ec.requestTimeout)
	defer cancel()
	logs, err := ec.client.FilterLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (ec *client) GetRateLimitDroppedRequests() int64 {
	return ec.rateLimitDroppedRequests
}
