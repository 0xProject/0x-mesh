// +build js, wasm

package orderfilter

import (
	"errors"
	"syscall/js"

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

func (f *Filter) checkForValidator() bool {
	if f.validatorLoaded {
		return true
	}
	if jsutil.IsNullOrUndefined(js.Global().Get("schemaValidator")) {
		return false
	}
	f.validatorLoaded = true
	return f.validatorLoaded

}

// ValidateOrderJSON Validates a JSON encoded signed order using the AJV javascript library.
// This libarary is used to increase the performance of Mesh nodes that run in the browser.
func (f *Filter) ValidateOrderJSON(orderJSON []byte) (*SchemaValidationResult, error) {
	if !f.checkForValidator() {
		return nil, errors.New(`"schemaValidator" has not been set on the "window" object`)
	}
	jsResult := js.Global().Get("schemaValidator").Call("orderValidator", string(orderJSON))
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

func (f *Filter) MatchOrderMessageJSON(messageJSON []byte) (bool, error) {
	if !f.checkForValidator() {
		return false, errors.New(`"schemaValidator" has not been set on the "window" object`)
	}
	jsResult := js.Global().Get("schemaValidator").Call("messageValidator", js.ValueOf(string(messageJSON)))
	fatal := jsResult.Get("fatal")
	if !jsutil.IsNullOrUndefined(fatal) {
		return false, errors.New(fatal.String())
	}
	return jsResult.Get("success").Bool(), nil
}

func (f *Filter) ValidateOrder(order *zeroex.SignedOrder) (*SchemaValidationResult, error) {
	orderJSON, err := order.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return f.ValidateOrderJSON(orderJSON)
}
