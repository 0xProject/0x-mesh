package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/zeroex"
)

type orderMessageV4 struct {
	MessageType string                `json:"messageType"`
	Order       *zeroex.SignedOrderV4 `json:"order"`
	Topics      []string              `json:"topics"`
}

func OrderToRawMessageV4(topic string, order *zeroex.SignedOrderV4) ([]byte, error) {
	return json.Marshal(orderMessageV4{
		MessageType: "order",
		Order:       order,
		Topics:      []string{topic},
	})
}

func RawMessageToOrderV4(data []byte) (*zeroex.SignedOrderV4, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}
