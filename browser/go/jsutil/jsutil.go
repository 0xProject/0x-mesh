// +build js,wasm

// Package jsutil contains various utility functions for working with
// JavaScript and WebAssemblysysa
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

func IsNullOrUndefined(value js.Value) bool {
	return value == js.Null() || value == js.Undefined()
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
