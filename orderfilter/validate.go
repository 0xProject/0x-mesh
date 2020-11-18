// +build !js

package orderfilter

import (
	"github.com/0xProject/0x-mesh/zeroex"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

func (f *Filter) ValidateOrderJSON(orderJSON []byte) (*jsonschema.Result, error) {
	return f.orderSchema.Validate(jsonschema.NewBytesLoader(orderJSON))
}

func (f *Filter) MatchOrderMessageJSON(messageJSON []byte) (bool, error) {
	result, err := f.messageSchema.Validate(jsonschema.NewBytesLoader(messageJSON))
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) ValidateOrder(order *zeroex.SignedV3Order) (*jsonschema.Result, error) {
	return f.orderSchema.Validate(jsonschema.NewGoLoader(order))
}
