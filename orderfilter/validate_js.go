// +build js, wasm

package orderfilter

import (
	"errors"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
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

// ValidateOrderJSON Validates a JSON encoded signed order using the AJV javascript library.
// This libarary is used to increase the performance of Mesh nodes that run in the browser.
func (f *Filter) ValidateOrderJSON(orderJSON []byte) (*SchemaValidationResult, error) {
	result, err := jsutil.AwaitPromiseContext(ctx, f.orderValidator.Invoke(string(orderJSON)))
	if err != nil {
		return nil, err
	}
	valid := result.Get("success").Bool()
	jsErrors := result.Get("errors")
	var convertedErrors []*SchemaValidationError
	for i := 0; i < jsErrors.Length(); i++ {
		convertedErrors = append(convertedErrors, &SchemaValidationError{errors.New(jsErrors.Index(i).String())})
	}
	return &SchemaValidationResult{valid: valid, errors: convertedErrors}, nil
}

func (f *Filter) MatchOrderMessageJSON(messageJSON []byte) (bool, error) {
	result, err := jsutil.AwaitPromiseContext(ctx, f.messageValidator.Invoke(string(messageJSON)))
	if err != nil {
		return false, err
	}
	return result.Get("success").Bool(), nil
}

func (f *Filter) ValidateOrder(order *zeroex.SignedOrder) (*SchemaValidationResult, error) {
	orderJSON, err := order.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return f.ValidateOrderJSON(orderJSON)
}
