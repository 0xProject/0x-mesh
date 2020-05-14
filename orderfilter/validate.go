// +build !js

package orderfilter

import (
	"context"

	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

func (f *Filter) ValidateOrderJSON(orderJSON []byte) (*jsonschema.Result, error) {
	return f.orderSchema.Validate(jsonschema.NewBytesLoader(orderJSON))
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
	result, err := f.messageSchema.Validate(jsonschema.NewBytesLoader(messageJSON))
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func (f *Filter) ValidateOrder(order *zeroex.SignedOrder) (*jsonschema.Result, error) {
	return f.orderSchema.Validate(jsonschema.NewGoLoader(order))
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
