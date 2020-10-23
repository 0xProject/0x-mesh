package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/zeroex"
)

type orderMessage struct {
	MessageType string              `json:"messageType"`
	Order       *zeroex.SignedOrder `json:"order"`
	Topics      []string            `json:"topics"`
}

// To implement this, we'll need to implement a custom JSON Marshaler and Unmarshaler
// FIXME(jalextowle): This encoding function has to be implemented for message_handler.go
// FIXME(jalextowle): We'll also need to select the correct topic for this order when
// we send it out.
//
// OrderToRawMessage encodes an order into an order message to be sent over the wire
func OrderToRawMessage(topics []string, order *zeroex.SignedOrder) ([]byte, error) {
	// FIXME(jalextowle): What is the best way to make this topic selection
	// nice?
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
		Topics:      []string{topics[0]},
	})
}

// FIXME(jalextowle): This encoding function has to be implemented for message_handler.go
//
// RawMessageToOrder decodes an order message sent over the wire into an order
func RawMessageToOrder(data []byte) (*zeroex.SignedOrder, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}
