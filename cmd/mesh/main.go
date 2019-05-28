// +build !js

// package mesh is a standalone 0x Mesh node that can be run from the command
// line. It uses environment variables for configuration and exposes a JSON RPC
// endpoint over WebSockets.
package main

import (
	"github.com/0xProject/0x-mesh/core"
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
