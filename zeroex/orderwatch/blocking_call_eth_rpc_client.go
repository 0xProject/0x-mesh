package orderwatch

import (
	"context"
	"math/big"

	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// callBlockingClient is identical to ethrpcclient.client but contract calls block until
// released. It is useful for testing high latency scenarios involving order validation
type callBlockingClient struct {
	client       ethrpcclient.Client
	blockingChan chan struct{}
}

// NewBlockingCallEthRPCClient returns a new instance of client
func NewBlockingCallEthRPCClient(ethRPCClient ethrpcclient.Client, blockingChan chan struct{}) (ethrpcclient.Client, error) {

	return &callBlockingClient{
		client:       ethRPCClient,
		blockingChan: blockingChan,
	}, nil
}

// CallContext performs a JSON-RPC call with the given arguments
func (cbc *callBlockingClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return cbc.client.CallContext(ctx, result, method, args...)
}

// HeaderByHash fetches a block header by its block hash.
func (cbc *callBlockingClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return cbc.client.HeaderByHash(ctx, hash)
}

// HeaderByNumber fetches a block header by block number.
func (cbc *callBlockingClient) HeaderByNumber(ctx context.Context, number *big.Int) (*miniheader.MiniHeader, error) {
	return cbc.client.HeaderByNumber(ctx, number)
}

// CodeAt returns the code of the given account.
func (cbc *callBlockingClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return cbc.client.CodeAt(ctx, contract, blockNumber)
}

// CallContract executes an Ethereum contract call with the specified data as the input.
func (cbc *callBlockingClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	<-cbc.blockingChan
	return cbc.client.CallContract(ctx, call, blockNumber)
}

// FilterLogs returns the logs that satisfy the supplied filter query.
func (cbc *callBlockingClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return cbc.client.FilterLogs(ctx, q)
}

func (cbc *callBlockingClient) GetRateLimitDroppedRequests() int64 {
	return 0
}
