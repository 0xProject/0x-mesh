// +build !js

package zeroex

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateOrderHash(t *testing.T) {
	fakeExchangeContractAddress := common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

	order := Order{
		MakerAddress:          constants.NullAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   constants.NullAddress,
		MakerAssetData:        constants.NullAddress.Bytes(),
		TakerAssetData:        constants.NullAddress.Bytes(),
		ExchangeAddress:       fakeExchangeContractAddress,
		Salt:                  big.NewInt(0),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(0),
		TakerAssetAmount:      big.NewInt(0),
		ExpirationTimeSeconds: big.NewInt(0),
	}

	// expectedOrderHash copied over from canonical order hashing test in Typescript library
	expectedOrderHash := common.HexToHash("0x434c6b41e2fb6dfcfe1b45c4492fb03700798e9c1afc6f801ba6203f948c1fa7")
	actualOrderHash, err := order.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, actualOrderHash)
}

func TestSignOrder(t *testing.T) {
	fakeExchangeContractAddress := common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

	order := &Order{
		MakerAddress:          constants.GanacheAccount0,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   constants.NullAddress,
		MakerAssetData:        constants.NullAddress.Bytes(),
		TakerAssetData:        constants.NullAddress.Bytes(),
		ExchangeAddress:       fakeExchangeContractAddress,
		Salt:                  big.NewInt(0),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(0),
		TakerAssetAmount:      big.NewInt(0),
		ExpirationTimeSeconds: big.NewInt(0),
	}
	ethClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	signer := ethereum.NewEthRPCSigner(ethClient)
	signedOrder, err := SignOrder(signer, order)
	require.NoError(t, err)

	expectedSignature := "0x1c3582f06356a1314dbf1c0e534c4d8e92e59b056ee607a7ff5a825f5f2cc5e6151c5cc7fdd420f5608e4d5bef108e42ad90c7a4b408caef32e24374cf387b0d7603"
	actualSignature := fmt.Sprintf("0x%s", common.Bytes2Hex(signedOrder.Signature))
	assert.Equal(t, expectedSignature, actualSignature)
}
