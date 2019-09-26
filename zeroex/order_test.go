package zeroex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOrder = &Order{
	ChainID:               big.NewInt(constants.TestNetworkID),
	ExchangeAddress:       ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
	MakerAddress:          constants.GanacheAccount0,
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   constants.NullAddress,
	MakerAssetData:        constants.NullAddress.Bytes(),
	MakerFeeAssetData:     constants.NullBytes,
	TakerAssetData:        constants.NullAddress.Bytes(),
	TakerFeeAssetData:     constants.NullBytes,
	Salt:                  big.NewInt(200),
	MakerFee:              big.NewInt(201),
	TakerFee:              big.NewInt(202),
	MakerAssetAmount:      big.NewInt(203),
	TakerAssetAmount:      big.NewInt(204),
	ExpirationTimeSeconds: big.NewInt(205),
}

var testHashOrder = &Order{
	ChainID:               big.NewInt(constants.TestNetworkID),
	ExchangeAddress:       ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
	MakerAddress:          constants.NullAddress,
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   constants.NullAddress,
	MakerAssetData:        constants.NullAddress.Bytes(),
	MakerFeeAssetData:     constants.NullAddress.Bytes(),
	TakerAssetData:        constants.NullAddress.Bytes(),
	TakerFeeAssetData:     constants.NullAddress.Bytes(),
	Salt:                  big.NewInt(0),
	MakerFee:              big.NewInt(0),
	TakerFee:              big.NewInt(0),
	MakerAssetAmount:      big.NewInt(0),
	TakerAssetAmount:      big.NewInt(0),
	ExpirationTimeSeconds: big.NewInt(0),
}

func TestGenerateOrderHash(t *testing.T) {
	// expectedOrderHash copied over from canonical order hashing test in Typescript library
	expectedOrderHash := common.HexToHash("0x224d188cf6709804d400e3ad6bd2311a99c25a3d24372a8a53fda5c71866f657")
	actualOrderHash, err := testHashOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, actualOrderHash)
}

func TestSignOrder(t *testing.T) {
	signedOrder, err := SignTestOrder(testOrder)
	require.NoError(t, err)

	expectedSignature := "0x1b66bd3722e2225d9ecfd2e086a199c1131e46fe4f892430fc717e9a7b4f74ffed10ef3c8068863a0a91e6454dc5fd4222284aaeac1df37cc79ede60067d0cdddc03"
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
		EndState:                     ESOrderAdded,
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
