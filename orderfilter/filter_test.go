package orderfilter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterValidateOrder(t *testing.T) {
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
			expectedErrors: []string{
				"exchangeAddress: Does not match pattern",
			},
			order: &zeroex.Order{
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
		},
		{
			note:              "happy path w/ custom sender address",
			chainID:           constants.TestChainID,
			customOrderSchema: `{"properties":{"senderAddress":{"type":"string","pattern":"0x00000000000000000000000000000000ba5eba11"}}}`,
			order: &zeroex.Order{
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
			customOrderSchema: `{"properties":{"senderAddress":{"type":"string","pattern":"0x00000000000000000000000000000000ba5eba11"}}}`,
			expectedErrors: []string{
				"senderAddress: Does not match pattern",
			},
			order: &zeroex.Order{
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
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s", i, tc.chainID, tc.customOrderSchema)
		filter, err := New(tc.chainID, tc.customOrderSchema)
		require.NoError(t, err, tcInfo)
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
				assert.Fail(t, fmt.Sprintf("missing expected error: %q", expectedErr), tcInfo)
			}
		}
	}
}

func TestFilterValidateOrderJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		note              string
		chainID           int
		customOrderSchema string
		orderJSON         []byte
		expectedErrors    []string
	}{
		{
			note:              "happy path",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
		},
		{
			note:              "order with mispelled makerAddress",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAdress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				"makerAddress is required",
			},
		},
		{
			note:              "order with missing makerAddress",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				"makerAddress is required",
			},
		},
		{
			note:              "order with invalid taker address",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"hi","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				"takerAddress: Does not match pattern '^0x[0-9a-fA-F]{40}$'",
			},
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s", i, tc.chainID, tc.customOrderSchema)
		filter, err := New(tc.chainID, tc.customOrderSchema)
		require.NoError(t, err)
		actualResult, err := filter.ValidateOrderJSON(tc.orderJSON)
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
				assert.Fail(t, fmt.Sprintf("missing expected error: %q", expectedErr), tcInfo)
			}
		}
	}
}
