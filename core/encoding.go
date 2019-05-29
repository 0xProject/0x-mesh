package core

import (
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
)

// maxSizeInBytes is the maximum number of bytes allowed for encoded orders. It
// is more than 10x the size of a typical ERC20 order.
const maxSizeInBytes = 8192

var errMaxSize = fmt.Errorf("message exceeds maximum size of %d bytes", maxSizeInBytes)

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

func validateMessageSize(message *p2p.Message) error {
	if len(message.Data) > maxSizeInBytes {
		return errMaxSize
	}
	return nil
}

func validateOrderSize(order *zeroex.SignedOrder) error {
	encoded, err := encodeOrder(order)
	if err != nil {
		return err
	}
	if len(encoded) > maxSizeInBytes {
		return errMaxSize
	}
	return nil
}

func filterOrdersBySize(orders []*zeroex.SignedOrder) (valid []*zeroex.SignedOrder, invalid []*zeroex.SignedOrder, err error) {
	valid = []*zeroex.SignedOrder{}
	invalid = []*zeroex.SignedOrder{}
	for _, order := range orders {
		err = validateOrderSize(order)
		if err == nil {
			valid = append(valid, order)
		} else if err == errMaxSize {
			invalid = append(invalid, order)
		} else {
			return nil, nil, err
		}
	}
	return valid, invalid, nil
}
