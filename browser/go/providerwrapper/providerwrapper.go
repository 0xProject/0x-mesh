// +build js,wasm

// Package providerwrapper wraps a web3 provider in order to implement the
// RPCClient interface.
package providerwrapper

import (
	"context"
	"errors"
	"syscall/js"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var _ ethclient.RPCClient = &RPCClient{}

type RPCClient struct {
	// provider is the underlying Web3 provider which will be used for sending
	// requests.
	provider js.Value
}

func NewRPCClient(provider js.Value) *RPCClient {
	return &RPCClient{
		provider: provider,
	}
}

func (c *RPCClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return errors.New("CallContext not yet implemented")
}

func (c *RPCClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return errors.New("BatchCallContext not yet implemented")
}

func (c *RPCClient) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	return nil, errors.New("EthSubscribe not yet implemented")
}

func (c *RPCClient) Close() {
	// no-op for now.
}
