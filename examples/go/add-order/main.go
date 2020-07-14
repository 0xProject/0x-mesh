// +build !js

// demo/add_order is a short program that adds an order to 0x Mesh via RPC
package main

// TODO(albrow): Update this to use the new GraphQL API.

// type clientEnvVars struct {
// 	// RPCAddress is the address of the 0x Mesh node to communicate with.
// 	WSRPCAddress string `envvar:"WS_RPC_ADDR"`
// 	// EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
// 	// API.
// 	EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL"`
// }

// var contractAddresses = ethereum.GanacheAddresses

// var testOrder = &zeroex.Order{
// 	ChainID:               big.NewInt(constants.TestChainID),
// 	MakerAddress:          constants.GanacheAccount0,
// 	TakerAddress:          constants.NullAddress,
// 	SenderAddress:         constants.NullAddress,
// 	FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 	MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
// 	MakerFeeAssetData:     constants.NullBytes,
// 	TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
// 	TakerFeeAssetData:     constants.NullBytes,
// 	Salt:                  big.NewInt(1548619145450),
// 	MakerFee:              big.NewInt(0),
// 	TakerFee:              big.NewInt(0),
// 	MakerAssetAmount:      big.NewInt(1000),
// 	TakerAssetAmount:      big.NewInt(2000),
// 	ExpirationTimeSeconds: big.NewInt(time.Now().Add(48 * time.Hour).Unix()),
// 	ExchangeAddress:       contractAddresses.Exchange,
// }

func main() {
	// 	env := clientEnvVars{}
	// 	if err := envvar.Parse(&env); err != nil {
	// 		panic(err)
	// 	}

	// 	client, err := rpc.NewClient(env.WSRPCAddress)
	// 	if err != nil {
	// 		log.WithError(err).Fatal("could not create client")
	// 	}

	// 	ethClient, err := ethrpc.Dial(env.EthereumRPCURL)
	// 	if err != nil {
	// 		log.WithError(err).Fatal("could not create Ethereum rpc client")
	// 	}

	// 	signer := signer.NewEthRPCSigner(ethClient)
	// 	signedTestOrder, err := zeroex.SignOrder(signer, testOrder)
	// 	if err != nil {
	// 		log.WithError(err).Fatal("could not sign 0x order")
	// 	}

	// 	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}
	// 	validationResults, err := client.AddOrders(signedTestOrders)
	// 	if err != nil {
	// 		log.WithError(err).Fatal("error from AddOrder")
	// 	} else {
	// 		log.Printf("submitted %d orders. Accepted: %d, Rejected: %d", len(signedTestOrders), len(validationResults.Accepted), len(validationResults.Rejected))
	// 	}
}
