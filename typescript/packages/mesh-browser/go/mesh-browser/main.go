// +build js,wasm

package main

import (
	"context"
	"encoding/json"
	"errors"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/typescript/packages/mesh-browser/go/jsutil"
	"github.com/0xProject/0x-mesh/typescript/packages/mesh-browser/go/providerwrapper"
	"github.com/0xProject/0x-mesh/zeroex"
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
				config, err := convertConfig(args[0])
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

// convertConfig converts a JavaScript config object into a core.Config. It also
// adds default values for any that are missing in the JavaScript object.
func convertConfig(jsConfig js.Value) (core.Config, error) {
	if jsutil.IsNullOrUndefined(jsConfig) {
		return core.Config{}, errors.New("config is required")
	}

	// Default config options. Some might be overridden.
	config := core.Config{
		Verbosity:                        2,
		DataDir:                          "0x-mesh",
		P2PTCPPort:                       0,
		P2PWebSocketsPort:                0,
		UseBootstrapList:                 true,
		BlockPollingInterval:             5 * time.Second,
		EthereumRPCMaxContentLength:      524288,
		EthereumRPCMaxRequestsPer24HrUTC: 100000,
		EthereumRPCMaxRequestsPerSecond:  30,
		EnableEthereumRPCRateLimiting:    true,
		MaxOrdersInStorage:               100000,
		CustomOrderFilter:                orderfilter.DefaultCustomOrderSchema,
	}

	// Required config options
	if ethereumChainID := jsConfig.Get("ethereumChainID"); jsutil.IsNullOrUndefined(ethereumChainID) {
		return core.Config{}, errors.New("ethereumChainID is required")
	} else {
		config.EthereumChainID = ethereumChainID.Int()
	}

	// Optional config options
	if verbosity := jsConfig.Get("verbosity"); !jsutil.IsNullOrUndefined(verbosity) {
		config.Verbosity = verbosity.Int()
	}
	if useBootstrapList := jsConfig.Get("useBootstrapList"); !jsutil.IsNullOrUndefined(useBootstrapList) {
		config.UseBootstrapList = useBootstrapList.Bool()
	}
	if bootstrapList := jsConfig.Get("bootstrapList"); !jsutil.IsNullOrUndefined(bootstrapList) {
		config.BootstrapList = bootstrapList.String()
	}
	if blockPollingIntervalSeconds := jsConfig.Get("blockPollingIntervalSeconds"); !jsutil.IsNullOrUndefined(blockPollingIntervalSeconds) {
		config.BlockPollingInterval = time.Duration(blockPollingIntervalSeconds.Int()) * time.Second
	}
	if ethereumRPCMaxContentLength := jsConfig.Get("ethereumRPCMaxContentLength"); !jsutil.IsNullOrUndefined(ethereumRPCMaxContentLength) {
		config.EthereumRPCMaxContentLength = ethereumRPCMaxContentLength.Int()
	}
	if ethereumRPCMaxRequestsPer24HrUTC := jsConfig.Get("ethereumRPCMaxRequestsPer24HrUTC"); !jsutil.IsNullOrUndefined(ethereumRPCMaxRequestsPer24HrUTC) {
		config.EthereumRPCMaxRequestsPer24HrUTC = ethereumRPCMaxRequestsPer24HrUTC.Int()
	}
	if ethereumRPCMaxRequestsPerSecond := jsConfig.Get("ethereumRPCMaxRequestsPerSecond"); !jsutil.IsNullOrUndefined(ethereumRPCMaxRequestsPerSecond) {
		config.EthereumRPCMaxRequestsPerSecond = ethereumRPCMaxRequestsPerSecond.Float()
	}
	if enableEthereumRPCRateLimiting := jsConfig.Get("enableEthereumRPCRateLimiting"); !jsutil.IsNullOrUndefined(enableEthereumRPCRateLimiting) {
		config.EnableEthereumRPCRateLimiting = enableEthereumRPCRateLimiting.Bool()
	}
	if customContractAddresses := jsConfig.Get("customContractAddresses"); !jsutil.IsNullOrUndefined(customContractAddresses) {
		config.CustomContractAddresses = customContractAddresses.String()
	}
	if maxOrdersInStorage := jsConfig.Get("maxOrdersInStorage"); !jsutil.IsNullOrUndefined(maxOrdersInStorage) {
		config.MaxOrdersInStorage = maxOrdersInStorage.Int()
	}
	if customOrderFilter := jsConfig.Get("customOrderFilter"); !jsutil.IsNullOrUndefined(customOrderFilter) {
		config.CustomOrderFilter = customOrderFilter.String()
	}
	if ethereumRPCURL := jsConfig.Get("ethereumRPCURL"); !jsutil.IsNullOrUndefined(ethereumRPCURL) && ethereumRPCURL.String() != "" {
		config.EthereumRPCURL = ethereumRPCURL.String()
	}
	if web3Provider := jsConfig.Get("web3Provider"); !jsutil.IsNullOrUndefined(web3Provider) {
		config.EthereumRPCClient = providerwrapper.NewRPCClient(web3Provider)
	}

	return config, nil
}

// NewMeshWrapper creates a new wrapper from the given config.
func NewMeshWrapper(config core.Config) (*MeshWrapper, error) {
	app, err := core.New(config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
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
		cw.errChan <- cw.app.Start(cw.ctx)
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
	// HACK(albrow): There is a more effecient way to do this, but for now,
	// just use JSON to convert to the Go type.
	encodedOrders := js.Global().Get("JSON").Call("stringify", rawOrders).String()
	var rawMessages []*json.RawMessage
	if err := json.Unmarshal([]byte(encodedOrders), &rawMessages); err != nil {
		return js.Undefined(), err
	}
	results, err := cw.app.AddOrders(cw.ctx, rawMessages, pinned)
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
func (cw *MeshWrapper) GetOrders(page int, perPage int, snapshotID string) (js.Value, error) {
	ordersResponse, err := cw.app.GetOrders(page, perPage, snapshotID)
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
		// getOrdersForPageAsync(page: number, perPage: number, snapshotID?: string): Promise<GetOrdersResponse>
		"getOrdersForPageAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return jsutil.WrapInPromise(func() (interface{}, error) {
				// snapshotID is optional in the JavaScript function. Check if it is
				// null or undefined.
				snapshotID := ""
				if !jsutil.IsNullOrUndefined(args[2]) {
					snapshotID = args[2].String()
				}
				return cw.GetOrders(args[0].Int(), args[1].Int(), snapshotID)
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
