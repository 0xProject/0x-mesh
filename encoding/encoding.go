package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/zeroex"
)

type orderMessage struct {
	MessageType string
	Order       *zeroex.SignedOrder
}

// EncodeOrder encodes an order into an order message to be sent over the wire
func EncodeOrder(order *zeroex.SignedOrder) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
	})
}

// DecodeOrder decodes an order message sent over the wire
func DecodeOrder(data []byte) (*zeroex.SignedOrder, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}
