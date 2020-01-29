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
	WSRPCAddress string `envvar:"WS_RPC_ADDR"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	env := clientEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	client, err := rpc.NewClient(env.WSRPCAddress)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}

	ctx := context.Background()
	orderEventsChan := make(chan []*zeroex.OrderEvent, 8000)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventsChan)
	if err != nil {
		log.WithError(err).Fatal("Couldn't set up OrderStream subscription")
	}
	defer clientSubscription.Unsubscribe()

	for {
		select {
		case orderEvents := <-orderEventsChan:
			for _, orderEvent := range orderEvents {
				log.WithFields(log.Fields{
					"event": orderEvent,
				}).Printf("received order event")
			}
		case err := <-clientSubscription.Err():
			log.Fatal(err)
		}
	}
}
