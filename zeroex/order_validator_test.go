// +build !js

// We currently don't run these tests in WASM because of an issue in Go. See the header of
// eth_watcher_test.go for more details.
package zeroex

import (
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/configs"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchValidate(t *testing.T) {
	signedOrder := &SignedOrder{
		MakerAddress:          common.HexToAddress("0x6924a03bb710eaf199ab6ac9f2bb148215ae9b5d"),
		TakerAddress:          nullAddress,
		SenderAddress:         nullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(3551808554499581700),
		TakerAssetAmount:      big.NewInt(300000000000000),
		ExpirationTimeSeconds: big.NewInt(1548619325),
		ExchangeAddress:       configs.GanacheExchangeAddress,
	}

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	signedOrders := []*SignedOrder{
		signedOrder,
	}

	ethClient, err := ethclient.Dial(configs.GanacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(GanacheOrderValidatorAddress, ethClient)
	require.NoError(t, err)

	orderInfos := orderValidator.BatchValidate(signedOrders)
	assert.Equal(t, len(orderInfos), 1)
	assert.Equal(t, Expired, orderInfos[orderHash].OrderStatus)
	assert.Equal(t, signedOrder, orderInfos[orderHash].SignedOrder)
}

func TestCalculateRemainingFillableTakerAmount(t *testing.T) {
	takerAssetAmount := big.NewInt(200000000000000000)
	makerAssetAmount := big.NewInt(100000000000000000)
	makerFee := big.NewInt(10000000000000000)
	signedOrder := &SignedOrder{
		MakerAddress:          common.HexToAddress("0x6924a03bb710eaf199ab6ac9f2bb148215ae9b5d"),
		TakerAddress:          nullAddress,
		SenderAddress:         nullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              makerFee,
		TakerFee:              big.NewInt(10000000000000000),
		MakerAssetAmount:      makerAssetAmount,
		TakerAssetAmount:      takerAssetAmount,
		ExpirationTimeSeconds: big.NewInt(99548619325),
		ExchangeAddress:       configs.GanacheExchangeAddress,
	}

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	orderInfo := wrappers.OrderInfo{
		OrderHash:                   orderHash,
		OrderStatus:                 uint8(Fillable),
		OrderTakerAssetFilledAmount: big.NewInt(0),
	}

	expectedRemainingAmountToTraderInfo := map[*big.Int]wrappers.TraderInfo{
		// No balances or allowances
		big.NewInt(0): wrappers.TraderInfo{
			MakerBalance:      big.NewInt(0),
			MakerAllowance:    big.NewInt(0),
			TakerBalance:      big.NewInt(0),
			TakerAllowance:    big.NewInt(0),
			MakerZrxBalance:   big.NewInt(0),
			MakerZrxAllowance: big.NewInt(0),
			TakerZrxBalance:   big.NewInt(0),
			TakerZrxAllowance: big.NewInt(0),
		},
		// Sufficient balances & allowances
		big.NewInt(200000000000000000): wrappers.TraderInfo{
			MakerBalance:      makerAssetAmount,
			MakerAllowance:    makerAssetAmount,
			TakerBalance:      takerAssetAmount,
			TakerAllowance:    takerAssetAmount,
			MakerZrxBalance:   makerFee,
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
		// Taker only has half the required amount BUT takerAddress is NULL address so it's
		// ignored.
		big.NewInt(200000000000000000): wrappers.TraderInfo{
			MakerBalance:      makerAssetAmount,
			MakerAllowance:    makerAssetAmount,
			TakerBalance:      new(big.Int).Div(takerAssetAmount, big.NewInt(2)),
			TakerAllowance:    takerAssetAmount,
			MakerZrxBalance:   makerFee,
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
		// Maker only has half the required balance
		big.NewInt(100000000000000000): wrappers.TraderInfo{
			MakerBalance:      new(big.Int).Div(makerAssetAmount, big.NewInt(2)),
			MakerAllowance:    makerAssetAmount,
			TakerBalance:      takerAssetAmount,
			TakerAllowance:    takerAssetAmount,
			MakerZrxBalance:   makerFee,
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
		// Maker only has half the required ZRX balance
		big.NewInt(100000000000000000): wrappers.TraderInfo{
			MakerBalance:      makerAssetAmount,
			MakerAllowance:    makerAssetAmount,
			TakerBalance:      takerAssetAmount,
			TakerAllowance:    takerAssetAmount,
			MakerZrxBalance:   new(big.Int).Div(makerFee, big.NewInt(2)),
			MakerZrxAllowance: makerFee,
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
	}

	for expectedRemainingnFillableTakerAssetAmount, traderInfo := range expectedRemainingAmountToTraderInfo {
		remainingFillableTakerAssetAmount := calculateRemainingFillableTakerAmount(signedOrder, orderInfo, traderInfo)
		assert.Equal(t, expectedRemainingnFillableTakerAssetAmount, remainingFillableTakerAssetAmount)
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
