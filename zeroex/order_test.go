package zeroex

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateOrderHash(t *testing.T) {
	fakeExchangeContractAddress := common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

	order := SignedOrder{
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

func TestECSignOrder(t *testing.T) {
	fakeExchangeContractAddress := common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

	order := SignedOrder{
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

	signerAddress := common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	signatureBytes, err := order.ECSign(signerAddress, rpcClient)
	require.NoError(t, err)

	expectedSignature := "0x1c5df471fd3ab082ef46b6d258e7f0e76a04aabbab6801b031cd4dc260d2677c13373fad62139cbae182350f91c630ee2125a7b811650229575fcc5492444b7be003"
	actualSignature := fmt.Sprintf("0x%s", common.Bytes2Hex(signatureBytes))
	assert.Equal(t, expectedSignature, actualSignature)
}
