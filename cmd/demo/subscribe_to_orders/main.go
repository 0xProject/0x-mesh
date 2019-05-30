// +build !js

// demo/add_order is a short program that adds an order to 0x Mesh via RPC
package main

import (
	"context"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
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

func main() {
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
	defer clientSubscription.Unsubscribe()

	for orderInfos := range orderInfosChan {
		for _, orderInfo := range orderInfos {
			log.Printf("Received order event: %+v\n", orderInfo)
		}
	}
}
