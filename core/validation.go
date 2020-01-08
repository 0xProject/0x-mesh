package core

import (
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
)

func validateMessageSize(message *p2p.Message) error {
	if len(message.Data) > constants.MaxMessageSizeInBytes {
		return constants.ErrMaxMessageSize
	}
	return nil
}

func validateOrderSize(order *zeroex.SignedOrder) error {
	// TODO(albrow): split up max order size and max message size.
	// encoded, err := encodeOrder(order)
	// if err != nil {
	// 	return err
	// }
	// if len(encoded) > constants.MaxOrderSizeInBytes {
	// 	return errMaxSize
	// }
	return nil
}
