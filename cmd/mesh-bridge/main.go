// +build !js

// mesh-bridge is a short program that bridges two Mesh nodes. This is useful in cases where
// we introduce a network-level breaking change but still want the liquidity from one network
//to flow to another
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

const (
	firstWSRPCAddressLabel  = "FirstWSRPCAddress"
	secondWSRPCAddressLabel = "SecondWSRPCAddress"
)

type clientEnvVars struct {
	FirstWSRPCAddress  string `envvar:"FIRST_WS_RPC_ADDRESS"`
	SecondWSRPCAddress string `envvar:"SECOND_WS_RPC_ADDRESS"`
	Verbosity          int    `envvar:"VERBOSE"`
}

func main() {
	env := clientEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	log.SetLevel(log.Level(env.Verbosity))

	firstClient, err := rpc.NewClient(env.FirstWSRPCAddress)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}
	stats, err := firstClient.GetStats()
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("stats", stats).Info("Spun up first client")

	secondClient, err := rpc.NewClient(env.SecondWSRPCAddress)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}
	stats, err = secondClient.GetStats()
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("stats", stats).Info("Spun up second client")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		pipeOrders(secondClient, firstClient, secondWSRPCAddressLabel, firstWSRPCAddressLabel)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		go pipeOrders(firstClient, secondClient, firstWSRPCAddressLabel, secondWSRPCAddressLabel)
	}()

	wg.Wait()
}

func pipeOrders(inClient, outClient *rpc.Client, inLabel, outLabel string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderEventsChan := make(chan []*zeroex.OrderEvent, 8000)
	clientSubscription, err := inClient.SubscribeToOrders(ctx, orderEventsChan)
	if err != nil {
		log.WithError(err).Fatal("Couldn't set up OrderStream subscription")
	}
	defer clientSubscription.Unsubscribe()
	signedOrdersCache := []*zeroex.SignedOrder{}
	for {
		select {
		case orderEvents := <-orderEventsChan:
			for _, orderEvent := range orderEvents {
				if orderEvent.EndState != zeroex.ESOrderAdded {
					continue
				}
				log.WithField("orderHash", orderEvent.OrderHash.Hex()).Info(fmt.Sprintf("Found new order %s -> %s", inLabel, outLabel))
				signedOrdersCache = append(signedOrdersCache, orderEvent.SignedOrder)
				// If more than 10 orders discovered, pipe them to the outClient
				if len(signedOrdersCache) >= 10 {
					validationResults, err := outClient.AddOrders(signedOrdersCache)
					if err != nil {
						log.Fatal(err)
					}
					log.Info(fmt.Sprintf("Sent %d orders from %s -> %s. Accepted: %d Rejected: %d", len(signedOrdersCache), inLabel, outLabel, len(validationResults.Accepted), len(validationResults.Rejected)))
					// Clear orders from cache
					signedOrdersCache = []*zeroex.SignedOrder{}
				}
			}
		case err := <-clientSubscription.Err():
			log.Fatal(err)
		}
	}
}
