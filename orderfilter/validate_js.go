// +build js, wasm

package orderfilter

import (
	"context"
	"errors"
	"fmt"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

var (
	orderValidatorLoaded   = false
	messageValidatorLoaded = false
	// validatorLoadCheckInterval is frequently to check whether the schema
	// validators have been loaded.
	validatorLoadCheckInterval = 50 * time.Millisecond
)

type SchemaValidationResult struct {
	valid  bool
	errors []error
}

func (s *SchemaValidationResult) Valid() bool {
	return s.valid
}

func (s *SchemaValidationResult) Errors() []error {
	return s.errors
}

// ValidateOrderJSON Validates a JSON encoded signed order using the AJV javascript library.
// This libarary is used to increase the performance of Mesh nodes that run in the browser.
func (f *Filter) ValidateOrderJSON(orderJSON []byte) (*SchemaValidationResult, error) {
	jsResult := js.Global().Get("schemaValidator").Call("orderValidator", js.ValueOf(string(orderJSON)))
	fatal := jsResult.Get("fatal")
	if !jsutil.IsNullOrUndefined(fatal) {
		return nil, errors.New(fatal.String())
	}
	valid := jsResult.Get("success").Bool()
	jsErrors := jsResult.Get("errors")
	var convertedErrors []error
	for i := 0; i < jsErrors.Length(); i++ {
		convertedErrors = append(convertedErrors, errors.New(jsErrors.Get(fmt.Sprintf("%d", i)).String()))
	}
	return &SchemaValidationResult{valid: valid, errors: convertedErrors}, nil
}

// MatchOrder returns true if the order passes the filter. It only returns an
// error if there was a problem with validation. For details about
// orders that do not pass the filter, use ValidateOrder.
func (f *Filter) MatchOrder(order *zeroex.SignedOrder) (bool, error) {
	result, err := f.ValidateOrder(order)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) MatchOrderMessageJSON(messageJSON []byte) (bool, error) {
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

// Dummy declaration to ensure that ValidatePubSubMessage matches the expected
// signature for pubsub.Validator.
var _ pubsub.Validator = (&Filter{}).ValidatePubSubMessage

// ValidatePubSubMessage is an implementation of pubsub.Validator and will
// return true if the contents of the message pass the message JSON Schema.
func (f *Filter) ValidatePubSubMessage(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	isValid, err := f.MatchOrderMessageJSON(msg.Data)
	if err != nil {
		log.WithError(err).Error("MatchOrderMessageJSON returned an error")
		return false
	}
	return isValid
}
