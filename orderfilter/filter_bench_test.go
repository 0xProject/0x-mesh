package orderfilter

import (
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func BenchmarkValidateOrder(b *testing.B) {
	order := &zeroex.Order{
		ChainID:               big.NewInt(constants.TestChainID),
		MakerAddress:          common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631"),
		MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		MakerAssetAmount:      math.MustParseBig256("1000"),
		MakerFee:              math.MustParseBig256("0"),
		TakerAddress:          common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TakerAssetData:        common.FromHex("0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		TakerAssetAmount:      math.MustParseBig256("2000"),
		TakerFee:              math.MustParseBig256("0"),
		SenderAddress:         common.HexToAddress("0x0000000000000000000000000000000000000000"),
		ExchangeAddress:       common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		ExpirationTimeSeconds: math.MustParseBig256("1574532801"),
		Salt:                  math.MustParseBig256("1548619145450"),
	}

	filter, err := New(constants.TestChainID, DefaultCustomOrderSchema, contractAddresses)
	require.NoError(b, err)
	signedOrder, err := zeroex.SignTestOrder(order)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = filter.ValidateOrder(signedOrder)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}
