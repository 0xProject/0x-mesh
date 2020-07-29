package orderfilter

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This uses the same test cases as TestFilterValidateOrder, but it marshals and
// then unmarshals the orderfilter before performing validation. This provides a
// sanity check that order filters work properly after being encoded and then
// decoded. More rigorous testing of these properties is tested in integration
// tests (FIXME(jalextowle): Name the test once they are written).
func TestMarshalAndUnmarshalFilter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		note              string
		chainID           int
		customOrderSchema string
		order             *zeroex.Order
		expectedErrors    []string
	}{
		{
			note:              "happy path",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			order: &zeroex.Order{
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
			},
		},
		{
			note:              "wrong exchangeAddress",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			order: &zeroex.Order{
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
				ExchangeAddress:       common.HexToAddress("0xdeadfa11"),
				FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
				ExpirationTimeSeconds: math.MustParseBig256("1574532801"),
				Salt:                  math.MustParseBig256("1548619145450"),
			},
			expectedErrors: []string{
				wrongExchangeAddressError,
			},
		},
		{
			note:              "wrong chainID",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			order: &zeroex.Order{
				ChainID:               big.NewInt(42),
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
			},
			expectedErrors: []string{
				mismatchedChainIDError,
			},
		},
		{
			note:              "happy path w/ custom sender address",
			chainID:           constants.TestChainID,
			customOrderSchema: `{"properties":{"senderAddress":{"const":"0x00000000000000000000000000000000ba5eba11"}}}`,
			order: &zeroex.Order{
				ChainID:               big.NewInt(constants.TestChainID),
				MakerAddress:          common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631"),
				MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
				MakerAssetAmount:      math.MustParseBig256("1000"),
				MakerFee:              math.MustParseBig256("0"),
				TakerAddress:          common.HexToAddress("0x0000000000000000000000000000000000000000"),
				TakerAssetData:        common.FromHex("0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
				TakerAssetAmount:      math.MustParseBig256("2000"),
				TakerFee:              math.MustParseBig256("0"),
				SenderAddress:         common.HexToAddress("0x00000000000000000000000000000000ba5eba11"),
				ExchangeAddress:       common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
				FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
				ExpirationTimeSeconds: math.MustParseBig256("1574532801"),
				Salt:                  math.MustParseBig256("1548619145450"),
			},
		},
		{
			note:              "wrong custom sender address",
			chainID:           constants.TestChainID,
			customOrderSchema: `{"properties":{"senderAddress":{"const":"0x00000000000000000000000000000000ba5eba11"}}}`,
			order: &zeroex.Order{
				ChainID:               big.NewInt(constants.TestChainID),
				MakerAddress:          common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631"),
				MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
				MakerAssetAmount:      math.MustParseBig256("1000"),
				MakerFee:              math.MustParseBig256("0"),
				TakerAddress:          common.HexToAddress("0x0000000000000000000000000000000000000000"),
				TakerAssetData:        common.FromHex("0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
				TakerAssetAmount:      math.MustParseBig256("2000"),
				TakerFee:              math.MustParseBig256("0"),
				SenderAddress:         common.HexToAddress("0x00000000000000000000000000000000defea7ed"),
				ExchangeAddress:       common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
				FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
				ExpirationTimeSeconds: math.MustParseBig256("1574532801"),
				Salt:                  math.MustParseBig256("1548619145450"),
			},
			expectedErrors: []string{
				mismatchedSenderAddressError,
			},
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s", i, tc.chainID, tc.customOrderSchema)
		filter, err := New(tc.chainID, tc.customOrderSchema, contractAddresses)
		require.NoError(t, err, tcInfo)
		marshalledFilter, err := filter.MarshalJSON()
		require.NoError(t, err)

		newFilter := &Filter{}
		err = newFilter.UnmarshalJSON(marshalledFilter)
		require.NoError(t, err)

		signedOrder, err := zeroex.SignTestOrder(tc.order)
		require.NoError(t, err)
		actualResult, err := filter.ValidateOrder(signedOrder)
		require.NoError(t, err, tc.customOrderSchema)
		if len(tc.expectedErrors) == 0 {
			assert.Len(t, actualResult.Errors(), 0, "expected no errors but received %d: %+v", len(actualResult.Errors()), actualResult.Errors())
		} else {
		loop:
			for _, expectedErr := range tc.expectedErrors {
				for _, actualErr := range actualResult.Errors() {
					if strings.Contains(actualErr.String(), expectedErr) {
						continue loop
					}
				}
				assert.Fail(t, fmt.Sprintf("missing expected error: %q\ngot errors: %v", expectedErr, actualResult.Errors()), tcInfo)
			}
		}
	}
}
