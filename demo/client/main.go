// demo/client is a short program that can be used for ad hoc integration
// testing.
package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ws"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

type clientEnvVars struct {
	// RPCPort is the port to use for the JSON RPC API over WebSockets. By
	// default, 0x Mesh will let the OS select a randomly available port.
	RPCPort int `envvar:"RPC_PORT"`
}

var testOrder = &zeroex.SignedOrder{
	MakerAddress:          common.HexToAddress("0x6924a03bb710eaf199ab6ac9f2bb148215ae9b5d"),
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
	MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
	TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
	Salt:                  big.NewInt(1548619145450),
	MakerFee:              big.NewInt(0),
	TakerFee:              big.NewInt(0),
	MakerAssetAmount:      big.NewInt(3551808554499581700),
	TakerAssetAmount:      big.NewInt(300000000000000),
	ExpirationTimeSeconds: big.NewInt(time.Now().Add(48 * time.Hour).Unix()),
	ExchangeAddress:       constants.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
}

func main() {
	env := clientEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	// TODO(albrow): Set up allowances and balances.

	rpcAddr := fmt.Sprintf("ws://localhost:%d", env.RPCPort)
	client, err := ws.NewClient(rpcAddr)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}

	for {
		// TODO(albrow): Create a valid signed order. The current one doesn't have
		// a signature or allowances/balances.
		hash, err := client.AddOrder(testOrder)
		if err != nil {
			log.WithError(err).Fatal("error from AddOrder")
		}
		log.Printf("added order: %s", hash.Hex())
		time.Sleep(10 * time.Second)
	}
}
