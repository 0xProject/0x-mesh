package zeroex

import (
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
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
