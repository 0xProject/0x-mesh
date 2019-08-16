// +build js,wasm

package p2p

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
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

// NewDHT returns a new Kademlia DHT instance configured to work with 0x Mesh
// in browser environments.
func NewDHT(ctx context.Context, storageDir string, host host.Host) (*dht.IpfsDHT, error) {
	return dht.New(ctx, host, dhtopts.Client(true), dhtopts.Protocols(dhtProtocolID))
}
