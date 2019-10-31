package zeroex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeExchangeContractAddress = common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

var testOrder = &Order{
	MakerAddress:          constants.GanacheAccount0,
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   constants.NullAddress,
	MakerAssetData:        constants.NullAddress.Bytes(),
	TakerAssetData:        constants.NullAddress.Bytes(),
	ExchangeAddress:       fakeExchangeContractAddress,
	Salt:                  big.NewInt(200),
	MakerFee:              big.NewInt(201),
	TakerFee:              big.NewInt(202),
	MakerAssetAmount:      big.NewInt(203),
	TakerAssetAmount:      big.NewInt(204),
	ExpirationTimeSeconds: big.NewInt(205),
}

func TestGenerateOrderHash(t *testing.T) {
	// expectedOrderHash copied over from canonical order hashing test in Typescript library
	expectedOrderHash := common.HexToHash("0x3fcd58a6613265e2b0deba902d7ff693f330a0af6e5b04805b44bbffd8a415d3")
	actualOrderHash, err := testOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, actualOrderHash)
}

func TestSignOrder(t *testing.T) {
	signedOrder, err := SignTestOrder(testOrder)
	require.NoError(t, err)

	expectedSignature := "0x1c54a8db1f96a1851886d966b46ceff87c4c2bba6e1b2e3da7e183912b9a328334633f08385016ee5f94dcb3c7a9ca80c2de59de14b73607b2c4eaf64f3d89915103"
	actualSignature := fmt.Sprintf("0x%s", common.Bytes2Hex(signedOrder.Signature))
	assert.Equal(t, expectedSignature, actualSignature)
}

func TestMarshalUnmarshalOrderEvent(t *testing.T) {
	signedOrder, err := SignTestOrder(testOrder)
	require.NoError(t, err)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	orderEvent := OrderEvent{
		OrderHash:                orderHash,
		SignedOrder:              signedOrder,
		EndState:                 ESOrderAdded,
		FillableTakerAssetAmount: big.NewInt(2000),
		ContractEvents:           []*ContractEvent{},
	}

	buf := &bytes.Buffer{}
	require.NoError(t, json.NewEncoder(buf).Encode(orderEvent))
	var decoded OrderEvent

	// We need to call ResetHash so that unexported hash field is equal in later
	// assertions.
	signedOrder.ResetHash()

	require.NoError(t, json.NewDecoder(buf).Decode(&decoded))
	assert.Equal(t, orderEvent, decoded)
}
