package core

import (
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/zeroex"
)

type orderMessage struct {
	MessageType string
	Order       *zeroex.SignedOrder
}

func encodeOrder(order *zeroex.SignedOrder) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
	})
}

func decodeOrder(data []byte) (*zeroex.SignedOrder, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}
