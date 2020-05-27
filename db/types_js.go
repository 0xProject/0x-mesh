// +build js,wasm

package db

import (
	"syscall/js"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
)

func (opts *Options) JSValue() js.Value {
	value, _ := jsutil.InefficientlyConvertToJS(opts)
	return value
}

func (query *OrderQuery) JSValue() js.Value {
	if query == nil {
		return js.Null()
	}
	value, _ := jsutil.InefficientlyConvertToJS(query)
	return value
}
