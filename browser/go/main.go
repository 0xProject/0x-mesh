// +build js,wasm

package main

import (
	"context"
	"errors"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/core"
)

const loadEventName = "0xmeshload"

func main() {
	setGlobals()
	triggerLoadEvent()

	// In order for callback functions to work, we can't allow main to exit.
	// Simply use select to block forever.
	select {}
}

func setGlobals() {
	zeroexMesh := map[string]interface{}{
		// newWrapperAsync(config: Config): Promise<MeshWrapper>;
		"newWrapperAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return wrapInPromise(func() (interface{}, error) {
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

func triggerLoadEvent() {
	event := js.Global().Get("document").Call("createEvent", "Event")
	event.Call("initEvent", loadEventName, true, true)
	js.Global().Call("dispatchEvent", event)
}

type MeshWrapper struct {
	app    *core.App
	ctx    context.Context
	cancel context.CancelFunc
}

func convertConfig(jsConfig js.Value) (core.Config, error) {
	if isNullOrUndefined(jsConfig) {
		return core.Config{}, errors.New("config is required")
	}

	// Default config options. Some might be overridden.
	config := core.Config{
		Verbosity:                   5,
		DataDir:                     "0x-mesh",
		P2PTCPPort:                  0,
		P2PWebSocketsPort:           0,
		UseBootstrapList:            true,
		BlockPollingInterval:        5 * time.Second,
		EthereumRPCMaxContentLength: 524288,
	}

	// Required conig options
	if ethereumRPCURL := jsConfig.Get("ethereumRPCURL"); isNullOrUndefined(ethereumRPCURL) || ethereumRPCURL.String() == "" {
		return core.Config{}, errors.New("ethereumRPCURL is required")
	} else {
		config.EthereumRPCURL = ethereumRPCURL.String()
	}
	if ethereumNetworkID := jsConfig.Get("ethereumNetworkID"); isNullOrUndefined(ethereumNetworkID) {
		return core.Config{}, errors.New("ethereumNetworkID is required")
	} else {
		config.EthereumNetworkID = ethereumNetworkID.Int()
	}

	// Optional config options
	if useBootstrapList := jsConfig.Get("useBootstrapList"); !isNullOrUndefined(useBootstrapList) {
		config.UseBootstrapList = useBootstrapList.Bool()
	}
	if orderExpirationBufferSeconds := jsConfig.Get("orderExpirationBufferSeconds"); !isNullOrUndefined(orderExpirationBufferSeconds) {
		config.OrderExpirationBuffer = time.Duration(orderExpirationBufferSeconds.Int()) * time.Second
	}
	if blockPollingIntervalSeconds := jsConfig.Get("blockPollingIntervalSeconds"); !isNullOrUndefined(blockPollingIntervalSeconds) {
		config.BlockPollingInterval = time.Duration(blockPollingIntervalSeconds.Int()) * time.Second
	}
	if ethereumRPCMaxContentLength := jsConfig.Get("ethereumRPCMaxContentLength"); !isNullOrUndefined(ethereumRPCMaxContentLength) {
		config.EthereumRPCMaxContentLength = ethereumRPCMaxContentLength.Int()
	}

	return config, nil
}

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

func (cw *MeshWrapper) Start() error {
	return cw.app.Start(cw.ctx)
}

func (cw *MeshWrapper) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		// startAsync(): Promise<void>;
		"startAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return wrapInPromise(func() (interface{}, error) {
				return nil, cw.Start()
			})
		}),
	})
}

func errorToJS(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}

func isNullOrUndefined(value js.Value) bool {
	return value == js.Null() || value == js.Undefined()
}

func wrapInPromise(f func() (interface{}, error)) js.Value {
	var executor js.Func
	executor = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			defer executor.Release()
			if result, err := f(); err != nil {
				reject.Invoke(errorToJS(err))
			} else {
				resolve.Invoke(result)
			}
		}()
		return nil
	})
	return js.Global().Get("Promise").New(executor)
}
