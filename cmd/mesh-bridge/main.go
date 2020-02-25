// +build !js

// mesh-bridge is a short program that bridges two Mesh nodes. This is useful in cases where
// we introduce a network-level breaking change but still want the liquidity from one network
//to flow to another
package main

import (
	"context"
	"time"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

const (
	firstWSRPCAddressLabel  = "FirstWSRPCAddress"
	secondWSRPCAddressLabel = "SecondWSRPCAddress"
	maxReceiveBatch         = 100
	receiveTimeout          = 1 * time.Second
	tenThousand             = 10000
)

type clientEnvVars struct {
	FirstWSRPCAddress  string `envvar:"FIRST_WS_RPC_ADDRESS"`
	SecondWSRPCAddress string `envvar:"SECOND_WS_RPC_ADDRESS"`
	Verbosity          int    `envvar:"VERBOSITY"`
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

	go func() {
		pipeOrders(secondClient, firstClient, secondWSRPCAddressLabel, firstWSRPCAddressLabel)
	}()
	pipeOrders(firstClient, secondClient, firstWSRPCAddressLabel, secondWSRPCAddressLabel)
}

func pipeOrders(inClient, outClient *rpc.Client, inLabel, outLabel string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderEventsChan := make(chan []*zeroex.OrderEvent, tenThousand)
	clientSubscription, err := inClient.SubscribeToOrders(ctx, orderEventsChan)
	if err != nil {
		log.WithError(err).Fatal("Couldn't set up OrderStream subscription")
	}
	defer clientSubscription.Unsubscribe()
	for {
		incomingSignedOrders, err := receiveBatch(orderEventsChan, clientSubscription, inLabel, outLabel)
		if err != nil {
			log.Fatal(err)
		}
		validationResults, err := outClient.AddOrders(incomingSignedOrders)
		if err != nil {
			log.Fatal(err)
		}
		log.WithFields(log.Fields{
			"from":        inLabel,
			"to":          outLabel,
			"numSent":     len(incomingSignedOrders),
			"numAccepted": len(validationResults.Accepted),
			"numRejected": len(validationResults.Rejected),
		}).Info("Finished bridging orders")
	}
}

func receiveBatch(inChan chan []*zeroex.OrderEvent, subscription *ethrpc.ClientSubscription, inLabel, outLabel string) ([]*zeroex.SignedOrder, error) {
	signedOrders := []*zeroex.SignedOrder{}
	timeoutChan := time.After(receiveTimeout)
	for {
		if len(signedOrders) >= maxReceiveBatch {
			return signedOrders, nil
		}
		select {
		case <-timeoutChan:
			return signedOrders, nil
		case orderEvents := <-inChan:
			for _, orderEvent := range orderEvents {
				if orderEvent.EndState != zeroex.ESOrderAdded {
					continue
				}
				log.WithFields(log.Fields{
					"from":      inLabel,
					"to":        outLabel,
					"orderHash": orderEvent.OrderHash.Hex(),
				}).Info("Found new order over bridge")
				signedOrders = append(signedOrders, orderEvent.SignedOrder)
			}
		case err := <-subscription.Err():
			log.Fatal(err)
		}
	}
}
