package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/zeroex"
)

type orderMessage struct {
	MessageType string                `json:"messageType"`
	Order       *zeroex.SignedV3Order `json:"order"`
	Topics      []string              `json:"topics"`
}

// OrderToRawMessage encodes an order into an order message to be sent over the wire
func OrderToRawMessage(topic string, order *zeroex.SignedV3Order) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
		Topics:      []string{topic},
	})
}

// RawMessageToOrder decodes an order message sent over the wire into an order
func RawMessageToOrder(data []byte) (*zeroex.SignedV3Order, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}
