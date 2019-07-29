// +build !js

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	log "github.com/sirupsen/logrus"
)

type rpcHandler struct {
	app *core.App
}

// listenRPC starts the RPC server and listens on config.RPCPort. It blocks
// until there is an error or the RPC server is closed.
func listenRPC(app *core.App, config standaloneConfig, ctx context.Context) error {
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
	return rpcServer.Listen(ctx)
}

// GetOrders is called when an RPC client calls GetOrders.
func (handler *rpcHandler) GetOrders(page, perPage int, snapshotID string) (*rpc.GetOrdersResponse, error) {
	log.WithFields(map[string]interface{}{
		"page":       page,
		"perPage":    perPage,
		"snapshotID": snapshotID,
	}).Debug("received GetOrders request via RPC")
	getOrdersResponse, err := handler.app.GetOrders(page, perPage, snapshotID)
	if err != nil {
		if _, ok := err.(core.ErrSnapshotNotFound); ok {
			return nil, err
		}
		// We don't want to leak internal error details to the RPC client.
		log.WithField("error", err.Error()).Error("internal error in AddOrders RPC call")
		return nil, constants.ErrInternal
	}
	return getOrdersResponse, nil
}

// AddOrders is called when an RPC client calls AddOrders.
func (handler *rpcHandler) AddOrders(signedOrdersRaw []*json.RawMessage) (*zeroex.ValidationResults, error) {
	log.WithField("count", len(signedOrdersRaw)).Debug("received AddOrders request via RPC")
	validationResults, err := handler.app.AddOrders(signedOrdersRaw)
	if err != nil {
		// We don't want to leak internal error details to the RPC client.
		log.WithField("error", err.Error()).Error("internal error in AddOrders RPC call")
		return nil, constants.ErrInternal
	}
	return validationResults, nil
}

// AddPeer is called when an RPC client calls AddPeer,
func (handler *rpcHandler) AddPeer(peerInfo peerstore.PeerInfo) error {
	log.Debug("received AddPeer request via RPC")
	if err := handler.app.AddPeer(peerInfo); err != nil {
		log.WithField("error", err.Error()).Error("internal error in AddPeer RPC call")
		return constants.ErrInternal
	}
	return nil
}

// SubscribeToOrders is called when an RPC client sends a `mesh_subscribe` request with the `orders` topic parameter
func (handler *rpcHandler) SubscribeToOrders(ctx context.Context) (*ethRpc.Subscription, error) {
	log.Debug("received order event subscription request via RPC")
	subscription, err := SetupOrderStream(ctx, handler.app)
	if err != nil {
		log.WithField("error", err.Error()).Error("internal error in `mesh_subscribe` to `orders` RPC call")
		return nil, constants.ErrInternal
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
		orderEventsChan := make(chan []*zeroex.OrderEvent)
		orderWatcherSub := app.SubscribeToOrderEvents(orderEventsChan)

		for {
			select {
			case orderEvents := <-orderEventsChan:
				err := notifier.Notify(rpcSub.ID, orderEvents)
				if err != nil {
					// TODO(fabio): The current implementation of `notifier.Notify` returns a
					// `write: broken pipe` error when it is called _after_ the client has
					// disconnected but before the corresponding error is received on the
					// `rpcSub.Err()` channel. This race-condition is not problematic beyond
					// the unnecessary computation and log spam resulting from it. Once this is
					// fixed upstream, give all logs an `Error` severity.
					logEntry := log.WithFields(map[string]interface{}{
						"error":            err.Error(),
						"subscriptionType": "orders",
						"orderEvents":      len(orderEvents),
					})
					message := "error while calling notifier.Notify"
					// If the network connection disconnects for longer then ~2mins and then comes
					// back up, we've noticed the call to `notifier.Notify` return `i/o timeout`
					// `net.OpError` errors everytime it's called and no values are sent over
					// `rpcSub.Err()` nor `notifier.Closed()`. In order to stop the error from
					// endlessly re-occuring, we unsubscribe and return for encountering this type of
					// error.
					if _, ok := err.(*net.OpError); ok {
						logEntry.Trace(message)
						orderWatcherSub.Unsubscribe()
						return
					}
					if strings.Contains(err.Error(), "write: broken pipe") {
						logEntry.Trace(message)
					} else {
						logEntry.Error(message)
					}
				}
			case err := <-rpcSub.Err():
				log.WithField("err", err).Error("rpcSub returned an error")
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
