// +build !js

// We currently don't run these tests in WASM because of an issue in Go. See the header of
// eth_watcher_test.go for more details.
package zeroex

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var unsupportedAssetData = common.Hex2Bytes("a2cb61b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064")
var malformedAssetData = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
var malformedSignature = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
var multiAssetAssetData = common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000204a7cb5fb70000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000002711000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000007d10000000000000000000000000000000000000000000000000000000000004e210000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

var testOrder = Order{
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
}

type testCase struct {
	Field               string
	Value               interface{}
	ExpectedOrderStatus OrderStatus
}

var testCases = []testCase{
	testCase{
		Field:               "MakerAssetAmount",
		Value:               big.NewInt(0),
		ExpectedOrderStatus: InvalidMakerAssetAmount,
	},
	testCase{
		Field:               "MakerAssetAmount",
		Value:               big.NewInt(1000000),
		ExpectedOrderStatus: Fillable,
	},
	testCase{
		Field:               "TakerAssetAmount",
		Value:               big.NewInt(0),
		ExpectedOrderStatus: InvalidTakerAssetAmount,
	},
	testCase{
		Field:               "TakerAssetAmount",
		Value:               big.NewInt(1000000),
		ExpectedOrderStatus: Fillable,
	},
	testCase{
		Field:               "MakerAssetData",
		Value:               multiAssetAssetData,
		ExpectedOrderStatus: InvalidMakerAssetData,
	},
	testCase{
		Field:               "TakerAssetData",
		Value:               multiAssetAssetData,
		ExpectedOrderStatus: InvalidTakerAssetData,
	},
	testCase{
		Field:               "MakerAssetData",
		Value:               malformedAssetData,
		ExpectedOrderStatus: InvalidMakerAssetData,
	},
	testCase{
		Field:               "TakerAssetData",
		Value:               malformedAssetData,
		ExpectedOrderStatus: InvalidTakerAssetData,
	},
	testCase{
		Field:               "MakerAssetData",
		Value:               unsupportedAssetData,
		ExpectedOrderStatus: InvalidMakerAssetData,
	},
	testCase{
		Field:               "TakerAssetData",
		Value:               unsupportedAssetData,
		ExpectedOrderStatus: InvalidTakerAssetData,
	},
	testCase{
		Field:               "ExpirationTimeSeconds",
		Value:               big.NewInt(time.Now().Add(-5 * time.Minute).Unix()),
		ExpectedOrderStatus: Expired,
	},
}

func TestBatchValidateOffChainCases(t *testing.T) {
	for _, testCase := range testCases {
		order := modifyOrder(t, testOrder, testCase.Field, testCase.Value)

		signedOrder, err := SignTestOrder(&order)
		require.NoError(t, err)

		orderHash, err := order.ComputeOrderHash()
		require.NoError(t, err)

		ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
		require.NoError(t, err)

		signedOrders := []*SignedOrder{
			signedOrder,
		}

		orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID)
		require.NoError(t, err)

		orderInfos := orderValidator.BatchValidate(signedOrders)
		assert.Len(t, orderInfos, 1)
		assert.Equal(t, testCase.ExpectedOrderStatus, orderInfos[orderHash].OrderStatus)
		assert.Equal(t, signedOrder, orderInfos[orderHash].SignedOrder)
	}
}

func TestBatchValidateMalformedSignature(t *testing.T) {
	signedOrder := &SignedOrder{
		Order: &testOrder,
	}
	// Incorrectly formatted signature
	signedOrder.Signature = malformedSignature

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	signedOrders := []*SignedOrder{
		signedOrder,
	}

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	orderInfos := orderValidator.BatchValidate(signedOrders)
	assert.Len(t, orderInfos, 1)
	assert.Equal(t, SignatureInvalid, orderInfos[orderHash].OrderStatus)
	assert.Equal(t, signedOrder, orderInfos[orderHash].SignedOrder)
}

func TestBatchValidateSignatureInvalid(t *testing.T) {
	signedOrder := &SignedOrder{
		Order: &testOrder,
	}
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

	orderInfos := orderValidator.BatchValidate(signedOrders)
	assert.Len(t, orderInfos, 1)
	assert.Equal(t, SignatureInvalid, orderInfos[orderHash].OrderStatus)
	assert.Equal(t, signedOrder, orderInfos[orderHash].SignedOrder)
}

func TestCalculateRemainingFillableTakerAmount(t *testing.T) {
	takerAssetAmount := big.NewInt(200000000000000000)
	makerAssetAmount := big.NewInt(100000000000000000)
	makerFee := big.NewInt(10000000000000000)
	order := testOrder
	order.TakerAssetAmount = takerAssetAmount
	order.MakerAssetAmount = makerAssetAmount
	order.MakerFee = makerFee
	signedOrder, err := SignTestOrder(&order)
	require.NoError(t, err)

	orderHash, err := order.ComputeOrderHash()
	require.NoError(t, err)

	orderInfo := wrappers.OrderInfo{
		OrderHash:                   orderHash,
		OrderStatus:                 uint8(Fillable),
		OrderTakerAssetFilledAmount: big.NewInt(0),
	}

	testCases := [...]struct {
		expectedRemainingAmount *big.Int
		traderInfo              wrappers.TraderInfo
	}{
		// No balances or allowances
		{
			expectedRemainingAmount: big.NewInt(0),
			traderInfo: wrappers.TraderInfo{
				MakerBalance:      big.NewInt(0),
				MakerAllowance:    big.NewInt(0),
				TakerBalance:      big.NewInt(0),
				TakerAllowance:    big.NewInt(0),
				MakerZrxBalance:   big.NewInt(0),
				MakerZrxAllowance: big.NewInt(0),
				TakerZrxBalance:   big.NewInt(0),
				TakerZrxAllowance: big.NewInt(0),
			},
		},
		// Sufficient balances & allowances
		{
			expectedRemainingAmount: big.NewInt(200000000000000000),
			traderInfo: wrappers.TraderInfo{
				MakerBalance:      makerAssetAmount,
				MakerAllowance:    makerAssetAmount,
				TakerBalance:      takerAssetAmount,
				TakerAllowance:    takerAssetAmount,
				MakerZrxBalance:   makerFee,
				MakerZrxAllowance: big.NewInt(10000000000000000),
				TakerZrxBalance:   big.NewInt(10000000000000000),
				TakerZrxAllowance: big.NewInt(10000000000000000),
			},
		},
		// Taker only has half the required amount BUT takerAddress is NULL address so it's
		// ignored.
		{
			expectedRemainingAmount: big.NewInt(200000000000000000),
			traderInfo: wrappers.TraderInfo{
				MakerBalance:      makerAssetAmount,
				MakerAllowance:    makerAssetAmount,
				TakerBalance:      new(big.Int).Div(takerAssetAmount, big.NewInt(2)),
				TakerAllowance:    takerAssetAmount,
				MakerZrxBalance:   makerFee,
				MakerZrxAllowance: big.NewInt(10000000000000000),
				TakerZrxBalance:   big.NewInt(10000000000000000),
				TakerZrxAllowance: big.NewInt(10000000000000000),
			},
		},
		// Maker only has half the required balance
		{
			expectedRemainingAmount: big.NewInt(100000000000000000),
			traderInfo: wrappers.TraderInfo{
				MakerBalance:      new(big.Int).Div(makerAssetAmount, big.NewInt(2)),
				MakerAllowance:    makerAssetAmount,
				TakerBalance:      takerAssetAmount,
				TakerAllowance:    takerAssetAmount,
				MakerZrxBalance:   makerFee,
				MakerZrxAllowance: big.NewInt(10000000000000000),
				TakerZrxBalance:   big.NewInt(10000000000000000),
				TakerZrxAllowance: big.NewInt(10000000000000000),
			},
		},
		// Maker only has half the required ZRX balance
		{
			expectedRemainingAmount: big.NewInt(100000000000000000),
			traderInfo: wrappers.TraderInfo{
				MakerBalance:      makerAssetAmount,
				MakerAllowance:    makerAssetAmount,
				TakerBalance:      takerAssetAmount,
				TakerAllowance:    takerAssetAmount,
				MakerZrxBalance:   new(big.Int).Div(makerFee, big.NewInt(2)),
				MakerZrxAllowance: makerFee,
				TakerZrxBalance:   big.NewInt(10000000000000000),
				TakerZrxAllowance: big.NewInt(10000000000000000),
			},
		},
	}

	for _, testCase := range testCases {
		remainingFillableTakerAssetAmount := calculateRemainingFillableTakerAmount(signedOrder, orderInfo, testCase.traderInfo)
		assert.Equal(t, testCase.expectedRemainingAmount, remainingFillableTakerAssetAmount)
	}

	// Order already half filled
	orderInfo = wrappers.OrderInfo{
		OrderHash:                   orderHash,
		OrderStatus:                 uint8(Fillable),
		OrderTakerAssetFilledAmount: new(big.Int).Div(takerAssetAmount, big.NewInt(2)),
	}
	// Sufficient balances & allowances
	traderInfo := wrappers.TraderInfo{
		MakerBalance:      makerAssetAmount,
		MakerAllowance:    makerAssetAmount,
		TakerBalance:      takerAssetAmount,
		TakerAllowance:    takerAssetAmount,
		MakerZrxBalance:   makerFee,
		MakerZrxAllowance: big.NewInt(10000000000000000),
		TakerZrxBalance:   big.NewInt(10000000000000000),
		TakerZrxAllowance: big.NewInt(10000000000000000),
	}
	remainingFillableTakerAssetAmount := calculateRemainingFillableTakerAmount(signedOrder, orderInfo, traderInfo)
	assert.Equal(t, new(big.Int).Div(takerAssetAmount, big.NewInt(2)), remainingFillableTakerAssetAmount)
}

func modifyOrder(t *testing.T, order Order, field string, value interface{}) Order {
	switch field {
	case "MakerAssetData":
		order.MakerAssetData = value.([]byte)
	case "TakerAssetData":
		order.TakerAssetData = value.([]byte)
	case "MakerAssetAmount":
		order.MakerAssetAmount = value.(*big.Int)
	case "TakerAssetAmount":
		order.TakerAssetAmount = value.(*big.Int)
	case "ExpirationTimeSeconds":
		order.ExpirationTimeSeconds = value.(*big.Int)
	default:
		t.Fatal("Unsupported order field: ", field)
	}
	return order
}
