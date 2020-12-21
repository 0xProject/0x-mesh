package main

import (
	"context"
	"net/http"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/graphql"
	"github.com/0xProject/0x-mesh/graphql/generated"

	gqlserver "github.com/99designs/gqlgen/graphql/handler"
)

// gracefulShutdownTimeout is the maximum amount of time to allow
// responding to any incoming requests after the server receives
// the signal to shutdown.
const gracefulShutdownTimeout = 10 * time.Second

func serveGraphQL(ctx context.Context, app *core.App, config *standaloneConfig) error {
	handler := http.NewServeMux()

	// Set up handler for GraphiQL
	if config.EnableGraphQLServer {
		handler.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(graphiQLPage)
		}))
	}

	// Set up handler for GrqphQL queries
	resolver := graphql.NewResolver(app, &graphql.ResolverConfig{
		SlowSubscriberTimeout: config.GraphQLSlowSubscriberTimeout,
	})
	graphQLServer := gqlserver.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	handler.Handle("/graphql", graphQLServer)

	// Start the server
	server := &http.Server{Addr: config.GraphQLServerAddr, Handler: handler}
	go func() {
		<-ctx.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
	}()
	return server.ListenAndServe()
}

var graphiQLPage = []byte(`
<html>
  <head>
    <title>0x Mesh GraphQL Playground</title>
    <link href="https://unpkg.com/graphiql@1.0.3/graphiql.min.css" rel="stylesheet" />
  </head>
  <body style="margin: 0;">
    <div id="graphiql" style="height: 100vh;"></div>

    <script
      crossorigin
      src="https://unpkg.com/react@16.13.1/umd/react.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/react-dom@16.13.1/umd/react-dom.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/graphiql@1.0.3/graphiql.min.js"
    ></script>

    <script>
      const graphQLFetcher = graphQLParams =>
        fetch('/graphql', {
          method: 'post',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(graphQLParams),
        })
          .then(response => response.json())
          .catch(() => response.text());
      ReactDOM.render(
        React.createElement(GraphiQL, { fetcher: graphQLFetcher }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>
`)
