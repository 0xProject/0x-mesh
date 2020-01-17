// +build js,wasm

package types

import (
	"encoding/json"
	"syscall/js"
)

func (s *Stats) JSValue() js.Value {
	// TODO(albrow): Optimize this. Remove other uses of the JSON
	// encoding/decoding hack.
	encodedStats, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	statsJS := js.Global().Get("JSON").Call("parse", string(encodedStats))
	return statsJS
}

// errorToJS converts a Go error to a JavaScript Error.
func ErrorToJS(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}

func IsNullOrUndefined(value js.Value) bool {
	return value == js.Null() || value == js.Undefined()
}

// wrapInPromise converts a potentially blocking Go function to a non-blocking
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
