// +build !js

package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
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

// AddOrders is called when an RPC client calls AddOrders.
func (handler *rpcHandler) AddOrders(orders []*zeroex.SignedOrder) (zeroex.OrderHashToSuccinctOrderInfo, error) {
	log.Debug("received AddOrders request via RPC")
	orderHashToSuccinctOrderInfo, err := handler.app.AddOrders(orders)
	if err != nil {
		// We don't want to leak internal error details to the RPC client.
		log.WithField("error", err.Error()).Error("internal error in AddOrders RPC call")
		return nil, errInternal
	}
	return orderHashToSuccinctOrderInfo, nil
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
