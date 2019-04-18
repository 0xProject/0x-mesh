package zeroex

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

func TestGenerateOrderHash(t *testing.T) {
	fakeExchangeContractAddress := common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

	order := SignedOrder{
		MakerAddress:          nullAddress,
		TakerAddress:          nullAddress,
		SenderAddress:         nullAddress,
		FeeRecipientAddress:   nullAddress,
		MakerAssetData:        nullAddress.Bytes(),
		TakerAssetData:        nullAddress.Bytes(),
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
