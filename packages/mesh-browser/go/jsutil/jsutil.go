// +build js,wasm

// Package jsutil contains various utility functions for working with
// JavaScript and WebAssembly
package jsutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"syscall/js"
)

func CopyBytesToJS(bytes []byte) (jsBytes js.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = RecoverError(r)
		}
	}()

	jsBytes = js.Global().Get("Uint8Array").New(len(bytes))
	copied := js.CopyBytesToJS(jsBytes, bytes)
	if copied != len(bytes) {
		return js.Undefined(), fmt.Errorf("should have copied %d bytes to JS but only copied %d", len(bytes), copied)
	}
	return jsBytes, nil
}

func CopyBytesToGo(jsBytes js.Value) (bytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = RecoverError(r)
		}
	}()

	jsBytesLength := jsBytes.Get("length").Int()
	bytes = make([]byte, jsBytesLength)
	copied := js.CopyBytesToGo(bytes, jsBytes)
	if copied != jsBytesLength {
		return nil, fmt.Errorf("should have copied %d bytes to Go but only copied %d", jsBytesLength, copied)
	}
	return bytes, nil
}

// ErrorToJS converts a Go error to a JavaScript Error.
func ErrorToJS(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}

// IsNullOrUndefined returns true if the given JavaScript value is either null
// or undefined.
func IsNullOrUndefined(value js.Value) bool {
	return value.IsNull() || value.IsUndefined()
}

// WrapInPromise converts a potentially blocking Go function to a non-blocking
// JavaScript Promise. If the function returns an error, the promise will reject
// with that error. Otherwise, the promise will resolve with the first return
// value.
func WrapInPromise(f func() (interface{}, error)) js.Value {
	var executor js.Func
	executor = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			defer executor.Release()
			if result, err := f(); err != nil {
				reject.Invoke(ErrorToJS(err))
			} else {
				resolve.Invoke(result)
			}
		}()
		return nil
	})
	return js.Global().Get("Promise").New(executor)
}

// AwaitPromiseContext is like AwaitPromise but accepts a context. If the context
// is canceled or times out before the promise resolves, it will return
// (js.Undefined, ctx.Error).
func AwaitPromiseContext(ctx context.Context, promise js.Value) (result js.Value, err error) {
	resultsChan := make(chan js.Value)
	errChan := make(chan js.Error)

	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			resultsChan <- args[0]
		}()
		return js.Undefined()
	})
	defer thenFunc.Release()
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			errChan <- js.Error{Value: args[0]}
		}()
		return js.Undefined()
	})
	defer catchFunc.Release()
	promise.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case <-ctx.Done():
		return js.Undefined(), ctx.Err()
	case result := <-resultsChan:
		return result, nil
	case err := <-errChan:
		return js.Undefined(), err
	}
}

// AwaitPromise accepts a js.Value representing a Promise. If the promise
// resolves, it returns (result, nil). If the promise rejects, it returns
// (js.Undefined, error). AwaitPromise has a synchronous-like API but does not
// block the JavaScript event loop.
func AwaitPromise(promise js.Value) (result js.Value, err error) {
	return AwaitPromiseContext(context.Background(), promise)
}

// InefficientlyConvertToJS converts the given Go value to a JS value by
// encoding to JSON and then decoding it. This function is not very efficient
// and its use should be phased out over time as much as possible.
func InefficientlyConvertToJS(value interface{}) (result js.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("unexpected error: (%T) %s", e, e)
			}
		}
	}()
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return js.Undefined(), err
	}
	return js.Global().Get("JSON").Call("parse", string(jsonBytes)), nil
}

// InefficientlyConvertFromJS converts the given JS value to a Go value and sets
// it. This function is not very efficient and its use should be phased out over
// time as much as possible.
func InefficientlyConvertFromJS(jsValue js.Value, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("unexpected error: (%T) %s", e, e)
			}
		}
	}()
	jsonString := js.Global().Get("JSON").Call("stringify", jsValue)
	return json.Unmarshal([]byte(jsonString.String()), value)
}

// RecoverError allows a function to recover from a thrown Javascript error if
// called inside of a deferred function with a recover statement.
func RecoverError(e interface{}) error {
	if e != nil {
		debug.PrintStack()
	}
	switch e := e.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	default:
		return fmt.Errorf("unexpected JavaScript error: (%T) %v", e, e)
	}
}
