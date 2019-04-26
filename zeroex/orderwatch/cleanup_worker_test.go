package orderwatch

import (
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/configs"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const ganacheEndpoint = "http://localhost:8545"

var nullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

func TestRevalidateOrder(t *testing.T) {
	order := zeroex.SignedOrder{
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

	orders := []zeroex.SignedOrder{
		order,
	}

	ethClient, err := ethclient.Dial(ganacheEndpoint)
	require.NoError(t, err)

	cleanupWorker, err := NewCleanupWorker(GanacheOrderValidatorAddress, ethClient)
	require.NoError(t, err)

	cleanupWorker.RevalidateOrders(orders)
}
