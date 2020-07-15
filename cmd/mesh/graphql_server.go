package main

import (
	"context"
	"net/http"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/graphql"
	"github.com/graph-gophers/graphql-go/relay"
)

func serveGraphQL(ctx context.Context, app *core.App, addr string, enableGraphiQL bool) error {
	schema, err := graphql.NewSchema(app)
	if err != nil {
		return err
	}

	handler := http.NewServeMux()
	if enableGraphiQL {
		handler.Handle("/graphiql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(graphiQLPage)
		}))
	}
	handler.Handle("/", &relay.Handler{Schema: schema})

	server := &http.Server{Addr: addr, Handler: handler}

	go func() {
		select {
		case <-ctx.Done():
			// TODO(albrow): Graceful shutdowns.
			server.Close()
		}
	}()

	return server.ListenAndServe()
}

// TODO(albrow): Update this to use the latest GraphiQL version.
var graphiQLPage = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
