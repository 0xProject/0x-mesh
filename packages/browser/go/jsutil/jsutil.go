// +build js,wasm

// Package jsutil contains various utility functions for working with
// JavaScript and WebAssembly
package jsutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"
)

// ErrorToJS converts a Go error to a JavaScript Error.
func ErrorToJS(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}

// IsNullOrUndefined returns true if the given JavaScript value is either null
// or undefined.
func IsNullOrUndefined(value js.Value) bool {
	return value.Equal(js.Null()) || value.Equal(js.Undefined())
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

// InefficientlyConvertToJS converts the given Go value to a JS value by
// encoding to JSON and then decoding it. This function is not very efficient
// and its use should be phased out over time as much as possible.
func InefficientlyConvertToJS(value interface{}) (js.Value, error) {
	var jsValue interface{}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return js.Undefined(), err
	}
	if err := json.NewDecoder(&buf).Decode(&jsValue); err != nil {
		return js.Undefined(), err
	}
	return js.ValueOf(jsValue), nil
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
