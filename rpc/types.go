package rpc

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// OrderInfo represents an fillable order and how much it could be filled for
type OrderInfo struct {
	OrderHash                common.Hash         `json:"orderHash"`
	SignedOrder              *zeroex.SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount *big.Int            `json:"fillableTakerAssetAmount"`
}

type orderInfoJSON struct {
	OrderHash                string              `json:"orderHash"`
	SignedOrder              *zeroex.SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount string              `json:"fillableTakerAssetAmount"`
}

// MarshalJSON is a custom Marshaler for OrderInfo
func (o OrderInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder,
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
	})
}

// UnmarshalJSON implements a custom JSON unmarshaller for the OrderEvent type
func (o *OrderInfo) UnmarshalJSON(data []byte) error {
	var orderInfoJSON orderInfoJSON
	err := json.Unmarshal(data, &orderInfoJSON)
	if err != nil {
		return err
	}

	o.OrderHash = common.HexToHash(orderInfoJSON.OrderHash)
	o.SignedOrder = orderInfoJSON.SignedOrder
	var ok bool
	o.FillableTakerAssetAmount, ok = math.ParseBig256(orderInfoJSON.FillableTakerAssetAmount)
	if !ok {
		return errors.New("Invalid uint256 number encountered for FillableTakerAssetAmount")
	}
	return nil
}
