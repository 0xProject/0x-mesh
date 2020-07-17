// +build js, wasm

package ordervalidator

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultEthRPCTimeout = 5 * time.Second
)

var (
	signedErc20Orders      []*zeroex.SignedOrder
	signedMultiAssetOrders []*zeroex.SignedOrder
	orderValidator         *OrderValidator
)

func init() {
	rpcClient, err := ethrpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	ethRPCClient, err := ethrpcclient.New(rpcClient, defaultEthRPCTimeout, ratelimit.NewUnlimited())
	if err != nil {
		panic(err)
	}
	orderValidator, err = New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ethereum.GanacheAddresses)
	if err != nil {
		panic(err)
	}

	erc20Data := common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32")
	erc20Order := zeroex.Order{
		ChainID:               big.NewInt(constants.TestChainID),
		MakerAddress:          constants.GanacheAccount1,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   constants.NullAddress,
		MakerAssetData:        erc20Data,
		MakerFeeAssetData:     constants.NullBytes,
		TakerAssetData:        erc20Data,
		TakerFeeAssetData:     constants.NullBytes,
		Salt:                  big.NewInt(int64(time.Now().Nanosecond())),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(100),
		TakerAssetAmount:      big.NewInt(42),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.GanacheAddresses.Exchange,
	}
	signedErc20Order, err := zeroex.SignTestOrder(&erc20Order)
	if err != nil {
		panic(err)
	}
	for i := 1; i < 100; i++ {
		signedErc20Orders = append(signedErc20Orders, signedErc20Order)
	}

	multiAssetData := common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000204a7cb5fb70000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000002711000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000007d10000000000000000000000000000000000000000000000000000000000004e210000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	multiAssetOrder := erc20Order
	multiAssetOrder.MakerAssetData = multiAssetData
	multiAssetOrder.MakerFeeAssetData = multiAssetData
	multiAssetOrder.TakerAssetData = multiAssetData
	multiAssetOrder.TakerFeeAssetData = multiAssetData

	signedMultiAssetOrder, err := zeroex.SignTestOrder(&multiAssetOrder)
	if err != nil {
		panic(err)
	}
	for i := 1; i < 100; i++ {
		signedMultiAssetOrders = append(signedMultiAssetOrders, signedMultiAssetOrder)
	}
}

func BenchmarkErc20ComputeOptimalChunkSizes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		orderValidator.computeOptimalChunkSizes(signedErc20Orders)
	}
}

func BenchmarkMultiAssetComputeOptimalChunkSizes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		orderValidator.computeOptimalChunkSizes(signedMultiAssetOrders)
	}
}
