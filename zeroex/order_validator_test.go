// +build !js

// We currently don't run these tests in WASM because of an issue in Go. See the header of
// eth_watcher_test.go for more details.
package zeroex

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var unsupportedAssetData = common.Hex2Bytes("a2cb61b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064")
var malformedAssetData = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
var malformedSignature = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
var multiAssetAssetData = common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000204a7cb5fb70000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000002711000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000007d10000000000000000000000000000000000000000000000000000000000004e210000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

var testSignedOrder = SignedOrder{
	Order: Order{
		MakerAddress:          constants.GanacheAccount0,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1000),
		TakerAssetAmount:      big.NewInt(2000),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(48 * time.Hour).Unix()),
		ExchangeAddress:       constants.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
	},
}

type testCase struct {
	SignedOrder                 SignedOrder
	IsValid                     bool
	ExpectedRejectedOrderStatus RejectedOrderStatus
}

func TestBatchValidateOffChainCases(t *testing.T) {
	var testCases = []testCase{
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetAmount(t, testSignedOrder, big.NewInt(0)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetAmount,
		},
		testCase{
			SignedOrder: signedOrderWithCustomMakerAssetAmount(t, testSignedOrder, big.NewInt(1000000)),
			IsValid:     true,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetAmount(t, testSignedOrder, big.NewInt(0)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetAmount,
		},
		testCase{
			SignedOrder: signedOrderWithCustomTakerAssetAmount(t, testSignedOrder, big.NewInt(1000000)),
			IsValid:     true,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetData(t, testSignedOrder, multiAssetAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetData(t, testSignedOrder, multiAssetAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetData(t, testSignedOrder, malformedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetData(t, testSignedOrder, malformedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetData(t, testSignedOrder, unsupportedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetData(t, testSignedOrder, unsupportedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomExpirationTimeSeconds(t, testSignedOrder, big.NewInt(time.Now().Add(-5*time.Minute).Unix())),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROExpired,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomSignature(t, testSignedOrder, malformedSignature),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidSignature,
		},
	}

	for _, testCase := range testCases {
		ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
		require.NoError(t, err)

		signedOrders := []*SignedOrder{
			&testCase.SignedOrder,
		}

		orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID)
		require.NoError(t, err)

		validationResults := orderValidator.BatchValidate(signedOrders)
		isValid := len(validationResults.Accepted) == 1
		assert.Equal(t, testCase.IsValid, isValid)
		if !isValid {
			assert.Equal(t, testCase.ExpectedRejectedOrderStatus, validationResults.Rejected[0].Status)
		}
	}
}

func TestBatchValidateSignatureInvalid(t *testing.T) {
	signedOrder := &testSignedOrder
	// Add a correctly formatted signature that does not correspond to this order
	signedOrder.Signature = common.Hex2Bytes("1c3582f06356a1314dbf1c0e534c4d8e92e59b056ee607a7ff5a825f5f2cc5e6151c5cc7fdd420f5608e4d5bef108e42ad90c7a4b408caef32e24374cf387b0d7603")

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	signedOrders := []*SignedOrder{
		signedOrder,
	}

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	validationResults := orderValidator.BatchValidate(signedOrders)
	assert.Len(t, validationResults.Accepted, 0)
	assert.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROInvalidSignature, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

func signedOrderWithCustomMakerAssetAmount(t *testing.T, signedOrder SignedOrder, makerAssetAmount *big.Int) SignedOrder {
	signedOrder.MakerAssetAmount = makerAssetAmount
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomTakerAssetAmount(t *testing.T, signedOrder SignedOrder, takerAssetAmount *big.Int) SignedOrder {
	signedOrder.TakerAssetAmount = takerAssetAmount
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomMakerAssetData(t *testing.T, signedOrder SignedOrder, makerAssetData []byte) SignedOrder {
	signedOrder.MakerAssetData = makerAssetData
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomTakerAssetData(t *testing.T, signedOrder SignedOrder, takerAssetData []byte) SignedOrder {
	signedOrder.TakerAssetData = takerAssetData
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomExpirationTimeSeconds(t *testing.T, signedOrder SignedOrder, expirationTimeSeconds *big.Int) SignedOrder {
	signedOrder.ExpirationTimeSeconds = expirationTimeSeconds
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomSignature(t *testing.T, signedOrder SignedOrder, signature []byte) SignedOrder {
	signedOrder.Signature = signature
	return signedOrder
}

func copyOrder(order Order) Order {
	return order
}
