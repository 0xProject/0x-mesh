// +build !js

package orderfilter

import (
	"github.com/0xProject/0x-mesh/zeroex"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

func (f *Filter) ValidateOrderV3JSON(orderJSON []byte) (*jsonschema.Result, error) {
	return f.orderSchemaV3.Validate(jsonschema.NewBytesLoader(orderJSON))
}

func (f *Filter) MatchOrderMessageV3JSON(messageJSON []byte) (bool, error) {
	result, err := f.messageSchemaV3.Validate(jsonschema.NewBytesLoader(messageJSON))
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) ValidateOrderV3(order *zeroex.SignedOrder) (*jsonschema.Result, error) {
	return f.orderSchemaV3.Validate(jsonschema.NewGoLoader(order))
}

func (f *Filter) ValidateOrderV4JSON(orderJSON []byte) (*jsonschema.Result, error) {
	return f.orderSchemaV4.Validate(jsonschema.NewBytesLoader(orderJSON))
}

func (f *Filter) MatchOrderMessageV4JSON(messageJSON []byte) (bool, error) {
	result, err := f.messageSchemaV4.Validate(jsonschema.NewBytesLoader(messageJSON))
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) ValidateOrderV4(order *zeroex.SignedOrder) (*jsonschema.Result, error) {
	return f.orderSchemaV4.Validate(jsonschema.NewGoLoader(order))
}
