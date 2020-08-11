// +build js,wasm

package main

import (
	"context"
	"encoding/json"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/browserutil"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

const (
	// loadEventName is the name of a global event that will be fired after all
	// WebAssembly is done loading.
	loadEventName = "0xmeshload"
	// orderEventsBufferSize is the buffer size for the orderEvents channel. If
	// the buffer is full, any additional events won't be processed.
	orderEventsBufferSize = 100
)

func main() {
	setGlobals()
	triggerLoadEvent()

	// In order for callback functions to work, we can't allow main to exit.
	// Simply use select to block forever.
	select {}
}

// setGlobals sets the global identifiers that are needed to interact with Mesh
// from the JavaScript world.
func setGlobals() {
	zeroexMesh := map[string]interface{}{
		// newWrapperAsync(config: Config): Promise<MeshWrapper>;
		"newWrapperAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return jsutil.WrapInPromise(func() (interface{}, error) {
				config, err := browserutil.ConvertConfig(args[0])
				if err != nil {
					return nil, err
				}
				return NewMeshWrapper(config)
			})
		}),
	}
	js.Global().Set("zeroExMesh", zeroexMesh)
}

// triggerLoadEvent triggers the global load event to indicate that the Wasm is
// done loading.
func triggerLoadEvent() {
	event := js.Global().Get("document").Call("createEvent", "Event")
	event.Call("initEvent", loadEventName, true, true)
	js.Global().Call("dispatchEvent", event)
}

// MeshWrapper is a wrapper around core.App. It exposes methods with basic,
// JavaScript-compatible types like string and int.
type MeshWrapper struct {
	app                     *core.App
	ctx                     context.Context
	cancel                  context.CancelFunc
	errChan                 chan error
	errHandler              js.Value
	orderEvents             chan []*zeroex.OrderEvent
	orderEventsSubscription event.Subscription
	orderEventsHandler      js.Value
}

// NewMeshWrapper creates a new wrapper from the given config.
func NewMeshWrapper(config core.Config) (*MeshWrapper, error) {
	ctx, cancel := context.WithCancel(context.Background())
	app, err := core.New(ctx, config)
	if err != nil {
		return nil, err
	}
	return &MeshWrapper{
		app:    app,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start starts core.App and sets up some channels. Unlike core.App.Start, it
// *does not* block. Instead, any erorrs that occur while Mesh is running
// will be sent through cw.errHandler.
func (cw *MeshWrapper) Start() error {
	cw.orderEvents = make(chan []*zeroex.OrderEvent, orderEventsBufferSize)
	cw.orderEventsSubscription = cw.app.SubscribeToOrderEvents(cw.orderEvents)
	cw.errChan = make(chan error, 1)

	// cw.app.Start blocks until there is an error or the app is closed, so we
	// need to start it in a goroutine.
	go func() {
		cw.errChan <- cw.app.Start()
	}()

	// Wait up to 1 second to see if cw.app.Start returns an error right away.
	// If it does, it probably indicates a configuration error which we should
	// return immediately.
	startTimeout := 1 * time.Second
	select {
	case err := <-cw.errChan:
		return err
	case <-time.After(startTimeout):
		break
	}

	// Otherwise listen for future events in a goroutine and return nil.
	go func() {
		for {
			select {
			case err := <-cw.errChan:
				// core.App exited with an error. Call errHandler.
				if !jsutil.IsNullOrUndefined(cw.errHandler) {
					cw.errHandler.Invoke(jsutil.ErrorToJS(err))
				}
			case <-cw.ctx.Done():
				return
			case events := <-cw.orderEvents:
				if !jsutil.IsNullOrUndefined(cw.orderEventsHandler) {
					eventsJS := make([]interface{}, len(events))
					for i, event := range events {
						eventsJS[i] = event.JSValue()
					}
					cw.orderEventsHandler.Invoke(eventsJS)
				}
			}
		}
	}()

	return nil
}

// AddOrders converts raw JavaScript orders into the appropriate type, calls
// core.App.AddOrders, converts the result into basic JavaScript types (string,
// int, etc.) and returns it.
func (cw *MeshWrapper) AddOrders(rawOrders js.Value, pinned bool) (js.Value, error) {
	var rawMessages []*json.RawMessage
	if err := jsutil.InefficientlyConvertFromJS(rawOrders, &rawMessages); err != nil {
		return js.Undefined(), err
	}
	results, err := cw.app.AddOrdersRaw(cw.ctx, rawMessages, pinned)
	if err != nil {
		return js.Undefined(), err
	}
	encodedResults, err := json.Marshal(results)
	resultsJS := js.Global().Get("JSON").Call("parse", string(encodedResults))
	return resultsJS, nil
}

// GetStats calls core.GetStats, converts the result to a js.Value and returns
// it.
func (cw *MeshWrapper) GetStats() (js.Value, error) {
	stats, err := cw.app.GetStats()
	if err != nil {
		return js.Undefined(), err
	}
	return js.ValueOf(stats), nil
}

// GetOrders converts raw JavaScript parameters into the appropriate type, calls
// core.App.GetOrders, converts the result into basic JavaScript types (string,
// int, etc.) and returns it.
func (cw *MeshWrapper) GetOrders(perPage int, minOrderHash string) (js.Value, error) {
	ordersResponse, err := cw.app.GetOrders(perPage, common.HexToHash(minOrderHash))
	if err != nil {
		return js.Undefined(), err
	}
	return js.ValueOf(ordersResponse), nil
}

// JSValue satisfies the js.Wrapper interface. The return value is a JavaScript
// object consisting of named functions. They act like methods by capturing the
// MeshWrapper through a closure.
func (cw *MeshWrapper) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		// startAsync(): Promise<void>;
		"startAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return jsutil.WrapInPromise(func() (interface{}, error) {
				return nil, cw.Start()
			})
		}),
		// onError(handler: (error: Error) => void): void;
		"onError": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			handler := args[0]
			cw.errHandler = handler
			return nil
		}),
		// onOrderEvents(handler: (events: Array<OrderEvent>) => void): void;
		"onOrderEvents": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			handler := args[0]
			cw.orderEventsHandler = handler
			return nil
		}),
		// getStatsAsync(): Promise<Stats>
		"getStatsAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return jsutil.WrapInPromise(func() (interface{}, error) {
				return cw.GetStats()
			})
		}),
		// getOrdersForPageAsync(perPage: number, minOrderHash?: string): Promise<GetOrdersResponse>
		"getOrdersForPageAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return jsutil.WrapInPromise(func() (interface{}, error) {
				// minOrderHash is optional in the JavaScript function. Check if it is
				// null or undefined.
				minOrderHash := ""
				if !jsutil.IsNullOrUndefined(args[1]) {
					minOrderHash = args[1].String()
				}
				return cw.GetOrders(args[0].Int(), minOrderHash)
			})
		}),
		// addOrdersAsync(orders: Array<SignedOrder>): Promise<ValidationResults>
		"addOrdersAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return jsutil.WrapInPromise(func() (interface{}, error) {
				return cw.AddOrders(args[0], args[1].Bool())
			})
		}),
	})
}
