package core

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

func encodeOrderMessage(topic string, order *zeroex.SignedOrder) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
		Topics:      []string{topic},
	})
}

func decodeOrderMessage(data []byte) (*zeroex.SignedOrder, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}
