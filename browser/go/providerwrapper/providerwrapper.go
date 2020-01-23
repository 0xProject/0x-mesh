// +build js,wasm

// Package providerwrapper wraps a web3 provider in order to implement the
// RPCClient interface.
package providerwrapper

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"syscall/js"

	"github.com/0xProject/0x-mesh/browser/go/jsutil"
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
	// Notable type definitions from Web3:
	//
	//     interface JSONRPCRequestPayload {
	//         params: any[];
	//         method: string;
	//         id: number;
	//         jsonrpc: string;
	//     }
	//
	//     type JSONRPCErrorCallback = (err: Error | null, result?: JSONRPCResponsePayload) => void;
	//
	//     sendAsync(payload: JSONRPCRequestPayload, callback: JSONRPCErrorCallback): void;
	//

	// Set up payload
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		// TODO(albrow): Do we need to do something special for the id here?
		"id":     rand.Intn(math.MaxInt32),
		"method": method,
	}
	if len(args) > 0 {
		// Convert args to a value that is compatible with syscall/js. Since we don't
		// know the underlying type of args, the only reliable way to do this is to
		// convert to and from JSON.
		convertedParams, err := jsutil.InefficientlyConvertToJS(args)
		if err != nil {
			return fmt.Errorf("invalid args for JSON payload: %s", err.Error())
		}
		payload["params"] = convertedParams
	}

	// Set up the callback function
	resultChan := make(chan js.Value, 1)
	errChan := make(chan error, 1)
	var callback js.Func
	callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer callback.Release()
		go func() {
			if len(args) == 0 {
				errChan <- errors.New("JSONRPCErrorCallback called with no arguments")
				return
			}
			jsError := args[0]
			if !jsutil.IsNullOrUndefined(jsError) {
				errChan <- js.Error{
					Value: jsError,
				}
				return
			}
			if len(args) < 2 {
				errChan <- errors.New("JSONRPCErrorCallback called with null/undefined error but no results")
				return
			}
			resultChan <- args[1]
			return
		}()
		return nil
	})

	// Call sendAsync and use select to wait for the results.
	c.provider.Call("sendAsync", payload, callback)
	select {
	case <-ctx.Done():
		return context.Canceled
	case err := <-errChan:
		return err
	case jsResult := <-resultChan:
		if err := jsutil.InefficientlyConvertFromJS(jsResult.Get("result"), result); err != nil {
			return fmt.Errorf("could not decode JSON RPC response: %s", err.Error())
		}
		return nil
	}
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
