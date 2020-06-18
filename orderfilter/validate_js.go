// +build js, wasm

package orderfilter

import (
	"errors"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/0xProject/0x-mesh/zeroex"
)

// sleepTime is used to force processes that block the event loop to give the
// event loop time to continue. This should be as low as possible.
const sleepTime = 500 * time.Microsecond

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
	resultChan := make(chan js.Value, 1)
	errChan := make(chan error, 1)
	go func() {
		time.Sleep(sleepTime)
		result, err := jsutil.AwaitPromise(f.orderValidator.Invoke(string(orderJSON)))
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()
	var result js.Value
	select {
	case err := <-errChan:
		return nil, err
	case result = <-resultChan:
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
	resultChan := make(chan js.Value, 1)
	errChan := make(chan error, 1)
	go func() {
		time.Sleep(sleepTime)
		result, err := jsutil.AwaitPromise(f.messageValidator.Invoke(string(messageJSON)))
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()
	var result js.Value
	select {
	case err := <-errChan:
		return false, err
	case result = <-resultChan:
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
