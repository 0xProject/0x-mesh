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
//
// OrderToRawMessage encodes an order into an order message to be sent over the wire
func OrderToRawMessage(topic string, order *zeroex.SignedOrder) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
		Topics:      []string{topic},
	})
}

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
