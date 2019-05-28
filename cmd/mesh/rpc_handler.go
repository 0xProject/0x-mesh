// +build !js

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	log "github.com/sirupsen/logrus"
)

var errInternal = errors.New("internal error")

type rpcHandler struct {
	app *core.App
}

// listenRPC starts the RPC server and listens on config.RPCPort. It blocks
// until there is an error or the RPC server is closed.
func listenRPC(app *core.App, config standaloneConfig) error {
	// Initialize the JSON RPC WebSocket server (but don't start it yet).
	rpcAddr := fmt.Sprintf(":%d", config.RPCPort)
	rpcHandler := &rpcHandler{
		app: app,
	}
	rpcServer, err := rpc.NewServer(rpcAddr, rpcHandler)
	if err != nil {
		return nil
	}
	go func() {
		// Wait for the server to start listening and select an address.
		for rpcServer.Addr() == nil {
			time.Sleep(10 * time.Millisecond)
		}
		log.WithField("address", rpcServer.Addr().String()).Info("started RPC server")
	}()
	return rpcServer.Listen()
}

// AddOrder is called when an RPC client calls AddOrder.
func (handler *rpcHandler) AddOrder(order *zeroex.SignedOrder) error {
	log.Debug("received AddOrder request via RPC")
	if err := handler.app.AddOrder(order); err != nil {
		if err == core.ErrInvalidOrder {
			// If the order is invalid, we want to communicate the error back to the
			// RPC client.
			log.WithField("order", order).Warn("received invalid order via RPC")
			return err
		} else {
			// All other error types are considered internal errors. We don't need to
			// leak details to the RPC client.
			log.WithField("error", err.Error()).Error("internal error in AddOrder RPC call")
			return errInternal
		}
	}
	return nil
}

// AddPeer is called when an RPC client calls AddPeer,
func (handler *rpcHandler) AddPeer(peerInfo peerstore.PeerInfo) error {
	log.Debug("received AddPeer request via RPC")
	if err := handler.app.AddPeer(peerInfo); err != nil {
		log.WithField("error", err.Error()).Error("internal error in AddPeer RPC call")
		return errInternal
	}
	return nil
}

// Orders is called when an RPC client sends a `mesh_subscribe` request with the `orders` topic parameter
func (handler *rpcHandler) Orders(ctx context.Context) (*ethRpc.Subscription, error) {
	log.Debug("received order event subscription request via RPC")
	subscription, err := SetupOrderStream(ctx, handler.app)
	if err != nil {
		log.WithField("error", err.Error()).Error("internal error in `mesh_subscribe` to `orders` RPC call")
		return nil, errInternal
	}
	return subscription, nil
}

// SetupOrderStream sets up the order stream for a subscription
func SetupOrderStream(ctx context.Context, app *core.App) (*ethRpc.Subscription, error) {
	notifier, supported := ethRpc.NotifierFromContext(ctx)
	if !supported {
		return &ethRpc.Subscription{}, ethRpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		orderInfosChan := make(chan []*zeroex.OrderInfo)
		orderWatcherSub := app.SubscribeToOrderEvents(orderInfosChan)

		for {
			select {
			case orderInfos := <-orderInfosChan:
				err := notifier.Notify(rpcSub.ID, orderInfos)
				if err != nil {
					log.WithField("error", err.Error()).Error("error while calling notifier.Notify")
				}
			case <-rpcSub.Err():
				orderWatcherSub.Unsubscribe()
				return
			case <-notifier.Closed():
				orderWatcherSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}
