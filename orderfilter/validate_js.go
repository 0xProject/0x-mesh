// +build js, wasm

package orderfilter

import (
	"context"
	"errors"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	jsutil.NextTick(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	jsutil.NextTick(ctx)
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
