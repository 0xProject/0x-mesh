package main

import (
	"context"
	"time"

	"github.com/0xProject/0x-mesh/core"
)

func main() {
	appConfig := core.Config{
		Verbosity:                   4,
		DataDir:                     "zeroex-mesh",
		P2PListenPort:               0,
		EthereumRPCURL:              "https://mainnet.infura.io/v3/af2e590be00f463fbfd0b546784065ad",
		EthereumNetworkID:           1,
		UseBootstrapList:            true,
		OrderExpirationBuffer:       10 * time.Second,
		BlockPollingInterval:        5 * time.Second,
		EthereumRPCMaxContentLength: 524288,
	}
	app, err := core.New(appConfig)
	if err != nil {
		panic(err)
	}

	go func() {
		if err = app.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Block forever.
	select {}
}
