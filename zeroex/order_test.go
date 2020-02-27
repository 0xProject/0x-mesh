package zeroex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOrder = &Order{
	ChainID:               big.NewInt(constants.TestChainID),
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
	ChainID:               big.NewInt(constants.TestChainID),
	ExchangeAddress:       common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
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

func init() {
	contractAddresses, err := ethereum.NewContractAddressesForChainID(constants.TestChainID)
	if err != nil {
		panic(err)
	}

	testOrder.ExchangeAddress = contractAddresses.Exchange
}

func TestGenerateOrderHash(t *testing.T) {
	// expectedOrderHash copied over from canonical order hashing test in Typescript library
	expectedOrderHash := common.HexToHash("0xcb36e4fedb36508fb707e2c05e21bffc7a72766ccae93f8ff096693fff7f1714")
	actualOrderHash, err := testHashOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, actualOrderHash)
}

func TestSignOrder(t *testing.T) {
	signedOrder, err := SignTestOrder(testOrder)
	require.NoError(t, err)

	expectedSignature := "0x1befcf4b6b1da4d207067a4b06e9bfbf21f85e2b6644f3ecf3a15f009e484756f251e3e00e909447ce45a16c620d14920a9acf516d9f4fe45bc36c914be6c9ec2703"
	actualSignature := fmt.Sprintf("0x%s", common.Bytes2Hex(signedOrder.Signature))
	assert.Equal(t, expectedSignature, actualSignature)
}

func TestMarshalUnmarshalOrderEvent(t *testing.T) {
	signedOrder, err := SignTestOrder(testOrder)
	require.NoError(t, err)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	orderEvent := OrderEvent{
		Timestamp:                time.Now().UTC(),
		OrderHash:                orderHash,
		SignedOrder:              signedOrder,
		EndState:                 ESOrderAdded,
		FillableTakerAssetAmount: big.NewInt(2000),
		ContractEvents: []*ContractEvent{
			{
				BlockHash: common.HexToHash("0x3fcd58a6613265e2b0deba902d7ff693f330a0af6e5b04805b44bbffd8a415d4"),
				TxHash:    common.HexToHash("0x3fcd58a6613265e2b0deba902d7ff693f330a0af6e5b04805b44bbffd8a415d5"),
				TxIndex:   42,
				LogIndex:  1337,
				IsRemoved: true,
				Address:   common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c49"),
				Kind:      "ERC20TransferEvent",
				Parameters: decoder.ERC20TransferEvent{
					From:  common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c50"),
					To:    common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c51"),
					Value: big.NewInt(120),
				},
			},
		},
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
