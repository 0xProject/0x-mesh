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

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

// Ensure that we implement the ethclient.RPCClient interface.
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

// CallContext performs a JSON-RPC call with the given arguments. If the context is
// canceled before the call has successfully returned, CallContext returns immediately.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
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
	//     interface JSONRPCResponseError {
	//         message: string;
	//         code: number;
	//     }
	//
	//     interface JSONRPCResponsePayload {
	//         result: any;
	//         id: number;
	//         jsonrpc: string;
	//         error?: JSONRPCResponseError;
	//     }
	//
	//     type JSONRPCErrorCallback = (err: Error | null, result?: JSONRPCResponsePayload) => void;
	//
	//     sendAsync(payload: JSONRPCRequestPayload, callback: JSONRPCErrorCallback): void;
	//

	// Set up payload
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      rand.Intn(math.MaxInt32),
		"method":  method,
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
		return ctx.Err()
	case err := <-errChan:
		return err
	case jsResult := <-resultChan:
		if rpcErr := jsResult.Get("error"); !jsutil.IsNullOrUndefined(rpcErr) {
			return jsErrorToRPCError(rpcErr)
		}
		if result == nil {
			return nil
		}
		if err := jsutil.InefficientlyConvertFromJS(jsResult.Get("result"), result); err != nil {
			return fmt.Errorf("could not decode JSON RPC response: %s", err.Error())
		}
		return nil
	}
}

// BatchCall sends all given requests as a single batch and waits for the server
// to return a response for all of them. The wait duration is bounded by the
// context's deadline.
//
// In contrast to CallContext, BatchCallContext only returns errors that have occurred
// while sending the request. Any error specific to a request is reported through the
// Error field of the corresponding BatchElem.
//
// Note that batch calls may not be executed atomically on the server side.
func (c *RPCClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	log.WithField("batch", b).Error("BatchCallContext was unexpectedly called in the browser")
	return errors.New("BatchCallContext not yet implemented")
}

// EthSubscribe registers a subscripion under the "eth" namespace.
func (c *RPCClient) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	log.WithField("args", args).Error("EthSubscribe was unexpectedly called in the browser")
	return nil, errors.New("EthSubscribe not yet implemented")
}

// Close is a no-op in this implementation.
func (c *RPCClient) Close() {
	// no-op for now.
}

// rpcError is an implementation of rpc.Error from the go-ethereum/rpc package.
type rpcError struct {
	Message string
	Code    int
}

var _ rpc.Error = &rpcError{}

func (e rpcError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("json-rpc error %d", e.Code)
	}
	return e.Message
}

func (e rpcError) ErrorCode() int {
	return e.Code
}

func jsErrorToRPCError(jsError js.Value) rpc.Error {
	return &rpcError{
		Message: jsError.Get("message").String(),
		Code:    jsError.Get("code").Int(),
	}
}
