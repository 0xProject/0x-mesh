// +build js, wasm

package orderfilter

import (
	"errors"

	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	"github.com/0xProject/0x-mesh/zeroex"
)

type SchemaValidationError struct {
	err error
}

func (s *SchemaValidationError) String() string {
	return s.err.Error()
}

type SchemaValidationResult struct {
	valid  bool
	errors []*SchemaValidationError
}

func (s *SchemaValidationResult) Valid() bool {
	return s.valid
}

func (s *SchemaValidationResult) Errors() []*SchemaValidationError {
	return s.errors
}

func (f *Filter) ValidateOrderV3JSON(orderJSON []byte) (*SchemaValidationResult, error) {
	jsResult := f.orderValidatorV3.Invoke(string(orderJSON))
	fatal := jsResult.Get("fatal")
	if !jsutil.IsNullOrUndefined(fatal) {
		return nil, errors.New(fatal.String())
	}
	valid := jsResult.Get("success").Bool()
	jsErrors := jsResult.Get("errors")
	var convertedErrors []*SchemaValidationError
	for i := 0; i < jsErrors.Length(); i++ {
		convertedErrors = append(convertedErrors, &SchemaValidationError{errors.New(jsErrors.Index(i).String())})
	}
	return &SchemaValidationResult{valid: valid, errors: convertedErrors}, nil
}

func (f *Filter) ValidateOrderV4JSON(orderJSON []byte) (*SchemaValidationResult, error) {
	jsResult := f.orderValidatorV4.Invoke(string(orderJSON))
	fatal := jsResult.Get("fatal")
	if !jsutil.IsNullOrUndefined(fatal) {
		return nil, errors.New(fatal.String())
	}
	valid := jsResult.Get("success").Bool()
	jsErrors := jsResult.Get("errors")
	var convertedErrors []*SchemaValidationError
	for i := 0; i < jsErrors.Length(); i++ {
		convertedErrors = append(convertedErrors, &SchemaValidationError{errors.New(jsErrors.Index(i).String())})
	}
	return &SchemaValidationResult{valid: valid, errors: convertedErrors}, nil
}

func (f *Filter) MatchOrderMessageV3JSON(messageJSON []byte) (bool, error) {
	jsResult := f.messageValidatorV3.Invoke(string(messageJSON))
	fatal := jsResult.Get("fatal")
	if !jsutil.IsNullOrUndefined(fatal) {
		return false, errors.New(fatal.String())
	}
	return jsResult.Get("success").Bool(), nil
}

func (f *Filter) MatchOrderMessageV4JSON(messageJSON []byte) (bool, error) {
	jsResult := f.messageValidatorV4.Invoke(string(messageJSON))
	fatal := jsResult.Get("fatal")
	if !jsutil.IsNullOrUndefined(fatal) {
		return false, errors.New(fatal.String())
	}
	return jsResult.Get("success").Bool(), nil
}

func (f *Filter) ValidateOrderV3(order *zeroex.SignedOrder) (*SchemaValidationResult, error) {
	orderJSON, err := order.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return f.ValidateOrderV3JSON(orderJSON)
}

func (f *Filter) ValidateOrderV4(order *zeroex.SignedOrder) (*SchemaValidationResult, error) {
	orderJSON, err := order.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return f.ValidateOrderV4JSON(orderJSON)
}
