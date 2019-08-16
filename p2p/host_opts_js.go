// +build js,wasm

package p2p

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	ws "github.com/libp2p/go-ws-transport"
)

func getOptionsForCurrentEnvironment(ctx context.Context, config Config) ([]libp2p.Option, error) {
	return []libp2p.Option{
		libp2p.Transport(ws.New),
		// Don't listen on any addresses by default. We can't accept incoming
		// connections in the browser.
		libp2p.ListenAddrs(),
	}, nil
}
