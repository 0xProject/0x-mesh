// +build js,wasm

package zeroex

import (
	"math/big"
	"syscall/js"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

var (
	rpcClient              *ethrpc.Client
	ethClient              *ethclient.Client
	makerAddress           = constants.GanacheAccount1
	takerAddress           = constants.GanacheAccount2
	tenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	wethAmount             = new(big.Int).Mul(big.NewInt(2), tenDecimalsInBaseUnits)
	zrxAmount              = new(big.Int).Mul(big.NewInt(1), tenDecimalsInBaseUnits)
)

func init() {
	var err error
	rpcClient, err = ethrpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
}

func TestContractEvent(t *testing.T) {
	contracts, err := ethereum.GetContractAddressesForChainID(1337)
	require.NoError(t, err)

	for _, event := range []ContractEvent{
		ContractEvent{
			BlockHash: common.HexToHash("0x1"),
			TxHash:    common.HexToHash("0x2"),
			TxIndex:   1,
			LogIndex:  2,
			IsRemoved: true,
			Address:   contracts.WETH9,
			Kind:      "ERC20TransferEvent",
			Parameters: decoder.ERC20TransferEvent{
				From:  constants.GanacheAccount0,
				To:    constants.GanacheAccount1,
				Value: big.NewInt(3),
			},
		},
		ContractEvent{
			BlockHash: constants.GanacheAccount1.Hash(),
			TxHash:    constants.GanacheAccount2.Hash(),
			TxIndex:   342424,
			LogIndex:  1000,
			IsRemoved: false,
			// NOTE(jalextowle): The ERC1155Proxy doesn't actually emit this event,
			// but that isn't important in the context of this test.
			Address: contracts.ERC1155Proxy,
			Kind:    "ERC1155TransferSingleEvent",
			Parameters: decoder.ERC1155TransferSingleEvent{
				Operator: constants.GanacheAccount1,
				From:     constants.GanacheAccount2,
				To:       constants.GanacheAccount3,
				Id:       big.NewInt(32423),
				Value:    big.NewInt(10000),
			},
		},
	} {
		jsEvent := event.JSValue()
		require.Equal(t, jsEvent.Get("address").String(), event.Address.Hex())
		require.Equal(t, jsEvent.Get("blockHash").String(), event.BlockHash.Hex())
		require.Equal(t, jsEvent.Get("txHash").String(), event.TxHash.Hex())
		require.Equal(t, jsEvent.Get("txIndex").Int(), int(event.TxIndex))
		require.Equal(t, jsEvent.Get("logIndex").Int(), int(event.LogIndex))
		require.Equal(t, jsEvent.Get("isRemoved").Bool(), event.IsRemoved)
		require.Equal(t, jsEvent.Get("kind").String(), event.Kind)
		// TODO(jalextowle): Check that the parameters are equal.
	}
}

func TestSignedOrder(t *testing.T) {
	orderCount := 10
	orders := signedTestOrders(t, orderCount)
	for _, order := range orders {
		jsOrder := order.JSValue()
		require.Equal(t, jsOrder.Get("chainId").Int64(), order.ChainID)
		require.Equal(t, jsOrder.Get("exchangeAddress").String(), order.ExchangeAddress.Hex())
		require.Equal(t, jsOrder.Get("senderAddress").String(), order.SenderAddress.Hex())
		require.Equal(t, jsOrder.Get("feeRecipientAddress").String(), order.FeeRecipientAddress.Hex())
		require.Equal(t, jsOrder.Get("expirationTimeSeconds").String(), order.ExpirationTimeSeconds.String())
		require.Equal(t, jsOrder.Get("salt").String(), order.Salt.String())
		require.Equal(t, jsOrder.Get("signature").String(), order.Signature)
		require.Equal(t, jsOrder.Get("makerAddress").String(), order.MakerAddress.Hex())
		require.Equal(t, jsOrder.Get("makerAssetAmount").String(), order.MakerAssetAmount.String())
		require.Equal(t, jsOrder.Get("makerAssetData").String(), string(order.MakerAssetData))
		require.Equal(t, jsOrder.Get("makerFee").String(), order.MakerFee.String())
		require.Equal(t, jsOrder.Get("makerFeeAssetData").String(), string(order.MakerFeeAssetData))
		require.Equal(t, jsOrder.Get("takerAddress").String(), order.TakerAddress.Hex())
		require.Equal(t, jsOrder.Get("takerAssetAmount").String(), order.TakerAssetAmount.String())
		require.Equal(t, jsOrder.Get("takerAssetData").String(), string(order.TakerAssetData))
		require.Equal(t, jsOrder.Get("takerFee").String(), order.TakerFee.String())
		require.Equal(t, jsOrder.Get("takerFeeAssetData").String(), string(order.TakerFeeAssetData))
	}
}

// TODO(jalextowle): Copied from core/message_handler
func signedTestOrders(t *testing.T, orderCount int) []*SignedOrder {
	orders := make([]*SignedOrder, orderCount)

	for i := range orders {
		orders[i] = scenario.CreateZRXForWETHSignedTestOrder(
			t,
			ethClient,
			makerAddress,
			takerAddress,
			new(big.Int).Add(wethAmount, big.NewInt(int64(i))),
			zrxAmount,
		)
	}

	return orders
}
