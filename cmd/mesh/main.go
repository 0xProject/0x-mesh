// +build !js

// package mesh is a standalone 0x Mesh node that can be run from the command
// line. It uses environment variables for configuration and exposes a JSON RPC
// endpoint over WebSockets.
package main

import (
	"context"
	"os"
	"sync"

	"github.com/0xProject/0x-mesh/core"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

// standaloneConfig contains configuration options specific to running 0x Mesh
// in standalone mode (i.e. not in a browser).
type standaloneConfig struct {
	// WSRPCAddr is the interface and port to use for the JSON-RPC API over
	// WebSockets. By default, 0x Mesh will listen on localhost and port 60557.
	WSRPCAddr string `envvar:"WS_RPC_ADDR" default:"localhost:60557"`
	// HTTPRPCAddr is the interface and port to use for the JSON-RPC API over
	// HTTP. By default, 0x Mesh will listen on localhost and port 60556.
	HTTPRPCAddr string `envvar:"HTTP_RPC_ADDR" default:"localhost:60556"`
}

func main() {
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
		log.WithField("error", err.Error()).Fatal("could not initialize app")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Below, we will start several independent goroutines. We use separate
	// channels to communicate errors and a waitgroup to wait for all goroutines
	// to exit.
	wg := &sync.WaitGroup{}

	coreErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Start(ctx); err != nil {
			coreErrChan <- err
		}
	}()

	// Start RPC server.
	rpcErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.WithFields(log.Fields{"ws_rpc_addr": config.WSRPCAddr, "http_rpc_addr": config.HTTPRPCAddr}).Info("starting RPC server")
		if err := listenRPC(ctx, app, config); err != nil {
			rpcErrChan <- err
		}
	}()

	// Block until there is an error or the app is closed.
	select {
	case <-ctx.Done():
		// We exited without error. Wait for all goroutines to finish and then
		// exit the process with a status code of 0.
		wg.Wait()
		os.Exit(0)
	case err := <-coreErrChan:
		cancel()
		log.WithField("error", err.Error()).Error("core app exited with error")
	case err := <-rpcErrChan:
		cancel()
		log.WithField("error", err.Error()).Error("RPC server returned error")
	}

	// If we reached here it means there was an error. Wait for all goroutines
	// to finish and then exit with non-zero status code.
	wg.Wait()
	os.Exit(1)
}
