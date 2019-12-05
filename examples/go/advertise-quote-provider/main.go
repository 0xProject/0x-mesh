// +build !js

// demo/add_order is a short program that adds an order to 0x Mesh via RPC
package main

import (
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

type clientEnvVars struct {
	// RPCAddress is the address of the 0x Mesh node to communicate with.
	RPCAddress string `envvar:"RPC_ADDRESS"`
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
	fmt.Println("Client created...")

	standard := "zaidan-v1.0"
	assetPair := "WETH/DAI"
	advertisementTTL := 5 * time.Minute
	err = client.AdvertiseAsQuoteProvider(standard, assetPair, advertisementTTL)
	if err != nil {
		log.WithError(err).Fatal("error from AdvertiseAsQuoteProvider")
	}
	fmt.Println("Quote provider advertised")
}
