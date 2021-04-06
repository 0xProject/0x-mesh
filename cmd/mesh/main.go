// +build !js

// package mesh is a standalone 0x Mesh node that can be run from the command
// line. It uses environment variables for configuration and optionally exposes
// a GraphQL API for developers to interact with.
package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/metrics"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

// standaloneConfig contains configuration options specific to running 0x Mesh
// in standalone mode (i.e. not in a browser).
type standaloneConfig struct {
	// EnableGraphQLServer determines whether or not to enable the GraphQL server.
	// If enabled, GraphQL queries can be sent to GraphQLServerAddr at the /graphql
	// URL. By default, the GraphQL server is disabled. Please be aware that the GraphQL
	// API is intended to be a *private* API. If you enable the GraphQL server in
	// production we recommend using a firewall/VPC or an authenticated proxy to restrict
	// public access.
	EnableGraphQLServer bool `envvar:"ENABLE_GRAPHQL_SERVER" default:"false"`
	// GraphQLServerAddr is the interface and port to use for the GraphQL API.
	// By default, 0x Mesh will listen on 0.0.0.0 (all available addresses) and
	// port 60557.
	GraphQLServerAddr string `envvar:"GRAPHQL_SERVER_ADDR" default:"0.0.0.0:60557"`
	// GraphQLSlowSubscriberTimeout is the maximum amount of time subscriber has to
	// accept events before being dropped.
	GraphQLSlowSubscriberTimeout time.Duration `envvar:"GRAPHQL_SLOW_SUBSCRIBER_TIMEOUT" default:"2s"`
	// EnableGraphQLPlayground determines whether or not to enable GraphiQL, an interactive
	// GraphQL playground which can be accessed by visiting GraphQLServerAddr in a browser.
	// See https://github.com/graphql/graphiql for more information. By default, GraphiQL
	// is disabled.
	EnableGraphQLPlayground bool `envvar:"ENABLE_GRAPHQL_PLAYGROUND" default:"false"`
	// EnablePrometheusMoniitoring determines whether or not to enable
	// prometheus monitoring. The metrics are accessed by scraping
	// {PrometheusMonitoringServerAddr}/metrics, prometheus is disabled.
	EnablePrometheusMonitoring bool `envvar:"ENABLE_PROMETHEUS_MONITORING" default:"false"`
	// PrometheusMonitoringServerAddr is the interface and port to use for
	// prometheus server metrics endpoint.
	PrometheusMonitoringServerAddr string `envvar:"PROMETHEUS_SERVER_ADDR" default:"0.0.0.0:8080"`
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

	graphQLErrChan := make(chan error, 1)
	if config.EnableGraphQLServer {
		// Start GraphQL server.
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.WithField("graphql_server_addr", config.GraphQLServerAddr).Info("starting GraphQL server")
			if err := serveGraphQL(ctx, app, &config); err != nil {
				graphQLErrChan <- err
			}
		}()
	}

	// NOTE: Prometehus is not an essential service to run.
	if config.EnablePrometheusMonitoring {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.WithField("prometheus_server_addr", config.PrometheusMonitoringServerAddr).Info("starting Prometheus metrics server")
			if err := metrics.ServeMetrics(ctx, config.PrometheusMonitoringServerAddr); err != nil {
				log.Error(err)
			}
		}()
	}

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
