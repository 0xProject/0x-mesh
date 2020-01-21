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

func (r *GetOrdersResponse) JSValue() js.Value {
	// TODO(albrow): Optimize this. Remove other uses of the JSON
	// encoding/decoding hack.
	encodedResponse, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	responseJS := js.Global().Get("JSON").Call("parse", string(encodedResponse))
	return responseJS
}
