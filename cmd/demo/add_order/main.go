// +build !js

// demo/add_order is a short program that adds an order to 0x Mesh via RPC
package main

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

type clientEnvVars struct {
	// RPCAddress is the address of the 0x Mesh node to communicate with.
	RPCAddress string `envvar:"RPC_ADDRESS"`
	// EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
	// API.
	EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL"`
}

var testOrder = &zeroex.Order{
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

func main() {
	log.SetOutput(os.Stdout)

	env := clientEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	client, err := rpc.NewClient(env.RPCAddress)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}

	ctx := context.Background()
	orderInfosChan := make(chan []*zeroex.OrderInfo, 8000)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderInfosChan)
	if err != nil {
		log.WithError(err).Fatal("Couldn't set up OrderStream subscription")
	}

	ethClient, err := ethrpc.Dial(env.EthereumRPCURL)
	if err != nil {
		log.WithError(err).Fatal("could not create Ethereum rpc client")
	}

	signer := ethereum.NewEthRPCSigner(ethClient)
	signedTestOrder, err := zeroex.SignOrder(signer, testOrder)
	if err != nil {
		log.WithError(err).Fatal("could not sign 0x order")
	}
	orderHash, _ := signedTestOrder.ComputeOrderHash()

	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}
	addOrdersResponse, err := client.AddOrders(signedTestOrders)
	if err != nil {
		log.WithError(err).Fatal("error from AddOrder")
	} else {
		log.Printf("submitted %d orders. Added: %d, Invalid: %d, FailedToAdd: %d", len(signedTestOrders), len(addOrdersResponse.Added), len(addOrdersResponse.Invalid), len(addOrdersResponse.FailedToAdd))
	}

	orderInfos := <-orderInfosChan
	log.Printf("Received order event: %+v\n", orderInfos[0])

	clientSubscription.Unsubscribe()
}
