// +build !js

// package mesh is a standalone 0x Mesh node that can be run from the command
// line. It uses environment variables for configuration and exposes a JSON RPC
// endpoint over WebSockets.
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/ws"
	"github.com/0xProject/0x-mesh/zeroex"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

// standaloneConfig contains configuration options specific to running 0x Mesh
// in standalone mode (i.e. not in a browser).
type standaloneConfig struct {
	// RPCPort is the port to use for the JSON RPC API over WebSockets. By
	// default, 0x Mesh will let the OS select a randomly available port.
	RPCPort int `envvar:"RPC_PORT" default:"0"`
}

func main() {
	// Configure logger to output JSON
	// TODO(albrow): Don't use global settings for logger.
	log.SetFormatter(&log.JSONFormatter{})

	// Parse env vars
	var coreConfig core.Config
	if err := envvar.Parse(&coreConfig); err != nil {
		log.WithField("error", err.Error()).Fatal("could not parse environment variables")
	}
	var config standaloneConfig
	if err := envvar.Parse(&config); err != nil {
		log.WithField("error", err.Error()).Fatal("could not parse environment variables")
	}

	// Start core.App.
	app, err := core.New(coreConfig)
	if err != nil {
		log.WithField("err", err.Error()).Fatal("could not initialize app")
	}
	if err := app.Start(); err != nil {
		log.WithField("err", err.Error()).Fatal("fatal error while starting app")
	}
	defer app.Close()

	// Start RPC server.
	go func() {
		err := listenRPC(app, config)
		if err != nil {
			app.Close()
			log.WithField("error", err.Error()).Fatal("RPC server returned error")
		}
	}()

	// Block forever or until the app is closed.
	select {}
}

// listenRPC starts the RPC server and listens on config.RPCPort. It blocks
// until there is an error or the RPC server is closed.
func listenRPC(app *core.App, config standaloneConfig) error {
	// Initialize the JSON RPC WebSocket server (but don't start it yet).
	rpcAddr := fmt.Sprintf(":%d", config.RPCPort)
	rpcHandler := &rpcHandler{
		app: app,
	}
	wsServer, err := ws.NewServer(rpcAddr, rpcHandler)
	if err != nil {
		return nil
	}
	go func() {
		// Wait for the server to start listening and select an address.
		for wsServer.Addr() == nil {
			time.Sleep(10 * time.Millisecond)
		}
		log.WithField("address", wsServer.Addr().String()).Info("started RPC server")
	}()
	return wsServer.Listen()
}

var errInternal = errors.New("internal error")

type rpcHandler struct {
	app *core.App
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
