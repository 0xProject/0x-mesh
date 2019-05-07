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

const ganacheEndpoint = "http://localhost:8545"

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

	ethClient, err := ethclient.Dial(ganacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(GanacheOrderValidatorAddress, ethClient)
	require.NoError(t, err)

	orderInfos := orderValidator.BatchValidate(signedOrders)
	assert.Equal(t, len(orderInfos), 1)
	assert.Equal(t, Expired, orderInfos[orderHash].OrderStatus)
	assert.Equal(t, signedOrder, orderInfos[orderHash].SignedOrder)
}

func TestCalculateRemainingFillableTakerAmount(t *testing.T) {
	signedOrder := &SignedOrder{
		MakerAddress:          common.HexToAddress("0x6924a03bb710eaf199ab6ac9f2bb148215ae9b5d"),
		TakerAddress:          nullAddress,
		SenderAddress:         nullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(10000000000000000),
		TakerFee:              big.NewInt(10000000000000000),
		MakerAssetAmount:      big.NewInt(100000000000000000),
		TakerAssetAmount:      big.NewInt(200000000000000000),
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
			MakerBalance:      big.NewInt(100000000000000000),
			MakerAllowance:    big.NewInt(100000000000000000),
			TakerBalance:      big.NewInt(200000000000000000),
			TakerAllowance:    big.NewInt(200000000000000000),
			MakerZrxBalance:   big.NewInt(10000000000000000),
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
		// Taker only has half the required amount BUT takerAddress is NULL address so it's
		// ignored.
		big.NewInt(200000000000000000): wrappers.TraderInfo{
			MakerBalance:      big.NewInt(100000000000000000),
			MakerAllowance:    big.NewInt(100000000000000000),
			TakerBalance:      big.NewInt(100000000000000000),
			TakerAllowance:    big.NewInt(200000000000000000),
			MakerZrxBalance:   big.NewInt(10000000000000000),
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
		// Maker only has half the required balance
		big.NewInt(100000000000000000): wrappers.TraderInfo{
			MakerBalance:      big.NewInt(50000000000000000),
			MakerAllowance:    big.NewInt(100000000000000000),
			TakerBalance:      big.NewInt(200000000000000000),
			TakerAllowance:    big.NewInt(200000000000000000),
			MakerZrxBalance:   big.NewInt(10000000000000000),
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
		// Maker only has half the required ZRX balance
		big.NewInt(100000000000000000): wrappers.TraderInfo{
			MakerBalance:      big.NewInt(100000000000000000),
			MakerAllowance:    big.NewInt(100000000000000000),
			TakerBalance:      big.NewInt(200000000000000000),
			TakerAllowance:    big.NewInt(200000000000000000),
			MakerZrxBalance:   big.NewInt(5000000000000000),
			MakerZrxAllowance: big.NewInt(10000000000000000),
			TakerZrxBalance:   big.NewInt(10000000000000000),
			TakerZrxAllowance: big.NewInt(10000000000000000),
		},
	}

	ethClient, err := ethclient.Dial(ganacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(GanacheOrderValidatorAddress, ethClient)
	require.NoError(t, err)

	for expectedRemainingnFillableTakerAssetAmount, traderInfo := range expectedRemainingAmountToTraderInfo {
		remainingFillableTakerAssetAmount := orderValidator.calculateRemainingFillableTakerAmount(signedOrder, orderInfo, traderInfo)
		assert.Equal(t, expectedRemainingnFillableTakerAssetAmount, remainingFillableTakerAssetAmount)
	}

}
