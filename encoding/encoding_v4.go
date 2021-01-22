package encoding

import (
	"encoding/json"

	"github.com/0xProject/0x-mesh/zeroex"
)

func OrderToRawMessageV4(topic string, order *zeroex.SignedOrderV4) ([]byte, error) {
	return json.Marshal(order)
}

func RawMessageToOrderV4(data []byte) (*zeroex.SignedOrderV4, error) {
	var order zeroex.SignedOrderV4
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return &order, nil
}
