// +build js, wasm

package orderfilter

import (
	"fmt"
	"math/big"
	"strings"
	"syscall/js"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

var (
	standardValidOrderJSON             = []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","chainId":1337,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`)
	orderWithSpecificSenderAddressJSON = []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x00000000000000000000000000000000ba5eba11","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","chainId":1337,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`)
	contractAddresses                  = ethereum.GanacheAddresses
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
			expectedErrors: []string{
				`{"keyword":"enum","dataPath":".exchangeAddress","schemaPath":"http://example.com/exchangeAddress/enum","params":{"allowedValues":["0x48bacb9266a570d521063ef5dd96e61686dbe788"]},"message":"should be equal to one of the allowed values"}`,
			},
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
		},
		{
			note:              "wrong chainID",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			expectedErrors: []string{
				`{"keyword":"const","dataPath":".chainId","schemaPath":"http://example.com/chainId/const","params":{"allowedValue":1337},"message":"should be equal to constant"}`,
			},
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
			expectedErrors: []string{
				`{"keyword":"const","dataPath":".senderAddress","schemaPath":"http://example.com/customOrder/properties/senderAddress/const","params":{"allowedValue":"0x00000000000000000000000000000000ba5eba11"},"message":"should be equal to constant"}`,
			},
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
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s", i, tc.chainID, tc.customOrderSchema)
		normalizedExchangeAddress := strings.ToLower(contractAddresses.Exchange.String())
		js.Global().Call("setSchemaValidator", js.ValueOf(tc.chainID), js.ValueOf(normalizedExchangeAddress), js.ValueOf(tc.customOrderSchema))
		filter, err := New(tc.chainID, tc.customOrderSchema, contractAddresses)
		require.NoError(t, err, tcInfo)
		signedOrder, err := zeroex.SignTestOrder(tc.order)
		require.NoError(t, err)
		actualResult, err := filter.ValidateOrder(signedOrder)
		require.NoError(t, err, tc.customOrderSchema)
		if len(tc.expectedErrors) == 0 {
			require.Len(t, actualResult.Errors(), 0, "expected no errors but received %d: %+v", len(actualResult.Errors()), actualResult.Errors())
		} else {
		loop:
			for _, expectedErr := range tc.expectedErrors {
				for _, actualErr := range actualResult.Errors() {
					if strings.Contains(actualErr.Error(), expectedErr) {
						continue loop
					}
				}
				require.Fail(t, fmt.Sprintf("missing expected error: %q\ngot errors: %v", expectedErr, actualResult.Errors()), tcInfo)
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
			orderJSON:         standardValidOrderJSON,
		},
		{
			note:              "order with mispelled makerAddress",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAdddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","chainId":1337,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				`{"keyword":"required","dataPath":"","schemaPath":"#/required","params":{"missingProperty":"makerAddress"},"message":"should have required property 'makerAddress'"}`,
			},
		},
		{
			note:              "order with missing makerAddress",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","chainId":1337,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				`{"keyword":"required","dataPath":"","schemaPath":"#/required","params":{"missingProperty":"makerAddress"},"message":"should have required property 'makerAddress'"}`,
			},
		},
		{
			note:              "order with invalid taker address",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"hi","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","chainId":1337,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				`{"keyword":"pattern","dataPath":".takerAddress","schemaPath":"http://example.com/address/pattern","params":{"pattern":"^0x[0-9a-fA-F]{40}$"},"message":"should match pattern \"^0x[0-9a-fA-F]{40}$\""}`,
			},
		},
		{
			note:              "order with wrong exchange address",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e6168deadbeef","chainId":1337,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				`{"keyword":"enum","dataPath":".exchangeAddress","schemaPath":"http://example.com/exchangeAddress/enum","params":{"allowedValues":["0x48bacb9266a570d521063ef5dd96e61686dbe788"]},"message":"should be equal to one of the allowed values"}`,
			},
		},
		{
			note:              "order with wrong chain ID",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         []byte(`{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","makerFeeAssetData":"0x","takerFeeAssetData":"0x","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","chainId":42,"signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`),
			expectedErrors: []string{
				`{"keyword":"const","dataPath":".chainId","schemaPath":"http://example.com/chainId/const","params":{"allowedValue":1337},"message":"should be equal to constant"}`,
			},
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s\nnote: %s", i, tc.chainID, tc.customOrderSchema, tc.note)
		// NOTE(jalextowle): Update the `schemaValidator` that is being used to use `tc.customOrderSchema`
		normalizedExchangeAddress := strings.ToLower(contractAddresses.Exchange.String())
		js.Global().Call("setSchemaValidator", js.ValueOf(tc.chainID), js.ValueOf(normalizedExchangeAddress), js.ValueOf(tc.customOrderSchema))
		filter, err := New(tc.chainID, tc.customOrderSchema, contractAddresses)
		require.NoError(t, err)
		actualResult, err := filter.ValidateOrderJSON(tc.orderJSON)
		require.NoError(t, err, tc.customOrderSchema)
		if len(tc.expectedErrors) == 0 {
			require.Len(t, actualResult.Errors(), 0, "expected no errors but received %d: %+v", len(actualResult.Errors()), actualResult.Errors())
		} else {
		loop:
			for _, expectedErr := range tc.expectedErrors {
				for _, actualErr := range actualResult.Errors() {
					if strings.Contains(actualErr.Error(), expectedErr) {
						continue loop
					}
				}
				require.Fail(t, fmt.Sprintf("missing expected error: %q\ngot errors: %s", expectedErr, actualResult.Errors()), tcInfo)
			}
		}
	}
}

func TestFilterMatchOrderMessageJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		note              string
		chainID           int
		customOrderSchema string
		orderMessageJSON  []byte
		expectedResult    bool
	}{
		{
			note:              "happy path",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"order","order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainId":1337,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"},"topics":["/0x-orders/version/3/chain/1337/schema/e30="]}`),
			expectedResult:    true,
		},
		{
			note:              "missing topics",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"order","order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainId":1337,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"}}`),
			expectedResult:    false,
		},
		{
			note:              "empty topics",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"order","order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainId":1337,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"},"topics":[]}`),
			expectedResult:    false,
		},
		{
			note:              "missing message type",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainId":1337,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"},"topics":["/0x-orders/version/3/chain/1337/schema/e30="]}`),
			expectedResult:    false,
		},
		{
			note:              "wrong message type",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"wrong","order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainId":1337,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"},"topics":["/0x-orders/version/3/chain/1337/schema/e30="]}`),
			expectedResult:    false,
		},
		{
			note:              "missing order",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"order","topics":["/0x-orders/version/3/chain/1337/schema/e30="]}`),
			expectedResult:    false,
		},
		{
			note:              "order with wrong exchange address",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"order","order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x0000000000000000000000000000000000000000","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainId":1337,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"},"topics":["/0x-orders/version/3/chain/1337/schema/e30="]}`),
			expectedResult:    false,
		},
		{
			note:              "order with wrong chain ID",
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderMessageJSON:  []byte(`{"messageType":"order","order":{"makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","makerAssetData":"0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c","makerAssetAmount":"100000000000000000000","makerFee":"0","takerAddress":"0x0000000000000000000000000000000000000000","takerAssetData":"0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082","takerAssetAmount":"50000000000000000000","takerFee":"0","senderAddress":"0x0000000000000000000000000000000000000000","exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","feeRecipientAddress":"0xa258b39954cef5cb142fd567a46cddb31a670124","expirationTimeSeconds":"1575499721","salt":"1548619145450","makerFeeAssetData":"0x","takerFeeAssetData":"0x","chainID":42,"signature":"0x1b0d147219c5c92262f0902727a8d72b09ea5165ac2ede14bccbfbf6559343d8305978e22516dc1ea75e10af2c8954cd45da562ec907ce5723a62728272c566a3f02"},"topics":["/0x-orders/version/3/chain/1337/schema/e30="]}`),
			expectedResult:    false,
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s\nnote: %s", i, tc.chainID, tc.customOrderSchema, tc.note)
		// NOTE(jalextowle): Update the `schemaValidator` that is being used to use `tc.customOrderSchema`
		normalizedExchangeAddress := strings.ToLower(contractAddresses.Exchange.String())
		js.Global().Call("setSchemaValidator", js.ValueOf(tc.chainID), js.ValueOf(normalizedExchangeAddress), js.ValueOf(tc.customOrderSchema))
		filter, err := New(tc.chainID, tc.customOrderSchema, contractAddresses)
		require.NoError(t, err)
		actualResult, err := filter.MatchOrderMessageJSON(tc.orderMessageJSON)
		require.NoError(t, err, tc.customOrderSchema)
		require.Equal(t, tc.expectedResult, actualResult, tcInfo)
	}
}

func TestFilterTopic(t *testing.T) {
	testCases := []struct {
		chainID           int
		customOrderSchema string
		// orderJSON must be valid according to the filter
		orderJSON     []byte
		expectedTopic string
	}{
		{
			chainID:           constants.TestChainID,
			customOrderSchema: DefaultCustomOrderSchema,
			orderJSON:         standardValidOrderJSON,
			expectedTopic:     "/0x-orders/version/3/chain/1337/schema/e30=",
		},
		{
			chainID:           constants.TestChainID,
			customOrderSchema: `{"properties":{"senderAddress":{"type":"string","pattern":"0x00000000000000000000000000000000ba5eba11"}}}`,
			orderJSON:         orderWithSpecificSenderAddressJSON,
			expectedTopic:     "/0x-orders/version/3/chain/1337/schema/eyJwcm9wZXJ0aWVzIjp7InNlbmRlckFkZHJlc3MiOnsicGF0dGVybiI6IjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDBiYTVlYmExMSIsInR5cGUiOiJzdHJpbmcifX19",
		},
		{
			// Same as above but the JSON schema has extra whitespace and some
			// properties are in a different order.
			chainID:           constants.TestChainID,
			customOrderSchema: `{"properties": {"senderAddress": {"pattern": "0x00000000000000000000000000000000ba5eba11", "type": "string"}}}`,
			orderJSON:         orderWithSpecificSenderAddressJSON,
			expectedTopic:     "/0x-orders/version/3/chain/1337/schema/eyJwcm9wZXJ0aWVzIjp7InNlbmRlckFkZHJlc3MiOnsicGF0dGVybiI6IjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDBiYTVlYmExMSIsInR5cGUiOiJzdHJpbmcifX19",
		},
	}

	for i, tc := range testCases {
		tcInfo := fmt.Sprintf("test case %d\nchainID: %d\nschema: %s", i, tc.chainID, tc.customOrderSchema)
		// NOTE(jalextowle): Update the `schemaValidator` that is being used to use `tc.customOrderSchema`
		normalizedExchangeAddress := strings.ToLower(contractAddresses.Exchange.String())
		js.Global().Call("setSchemaValidator", js.ValueOf(tc.chainID), js.ValueOf(normalizedExchangeAddress), js.ValueOf(tc.customOrderSchema))
		originalFilter, err := New(tc.chainID, tc.customOrderSchema, contractAddresses)
		require.NoError(t, err, tcInfo)
		result, err := originalFilter.ValidateOrderJSON(tc.orderJSON)
		require.NoError(t, err, tcInfo)
		require.Empty(t, result.Errors(), "original filter should validate the given order\n"+tcInfo)
		require.Equal(t, tc.expectedTopic, originalFilter.Topic(), tcInfo)
		newFilter, err := NewFromTopic(originalFilter.Topic(), contractAddresses)
		require.NoError(t, err, tcInfo)
		require.Equal(t, tc.expectedTopic, newFilter.Topic(), tcInfo)
		result, err = newFilter.ValidateOrderJSON(tc.orderJSON)
		require.NoError(t, err, tcInfo)
		require.Empty(t, result.Errors(), "filter generated from topic should validate the same order\n"+tcInfo)
	}
}

func TestDefaultOrderSchemaTopic(t *testing.T) {
	chainID := 1337
	defaultTopic, err := GetDefaultTopic(chainID, contractAddresses)
	require.NoError(t, err)
	expectedTopic := "/0x-orders/version/3/chain/1337/schema/e30="
	require.Equal(t, expectedTopic, defaultTopic, "the topic for the default filter should not change")
}
