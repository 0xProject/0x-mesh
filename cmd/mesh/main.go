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
	// GraphQLServerAddr is the interface and port to use for the GraphQL API.
	// By default, 0x Mesh will listen on localhost and port 60557.
	GraphQLServerAddr string `envvar:"GRAPHQL_SERVER_ADDR" default:"localhost:60557"`
	// EnableGraphiQL determines whether or not to enable GraphiQL, an interactive
	// GraphQL IDE which can be accessed by visiting /graphiql in a browser. See
	// https://github.com/graphql/graphiql for more information. By default, GraphiQL
	// is disabled.
	EnableGraphiQL bool `envvar:"ENABLE_GRAPHIQL" default:"false"`
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

	// Initialize core.App.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := core.New(ctx, coreConfig)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("could not initialize app")
	}

	// Below, we will start several independent goroutines. We use separate
	// channels to communicate errors and a waitgroup to wait for all goroutines
	// to exit.
	wg := &sync.WaitGroup{}

	coreErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Start(); err != nil {
			coreErrChan <- err
		}
	}()

	// Start GraphQL server.
	graphQLErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.WithField("graphql_server_addr", config.GraphQLServerAddr).Info("starting GraphQL server")
		if err := serveGraphQL(ctx, app, config.GraphQLServerAddr, config.EnableGraphiQL); err != nil {
			graphQLErrChan <- err
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
	case err := <-graphQLErrChan:
		cancel()
		log.WithField("error", err.Error()).Error("GraphQL server returned error")
	}

	// If we reached here it means there was an error. Wait for all goroutines
	// to finish and then exit with non-zero status code.
	wg.Wait()
	os.Exit(1)
}
