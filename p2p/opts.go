// +build !js

package p2p

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	leveldbStore "github.com/ipfs/go-ds-leveldb"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	"github.com/libp2p/go-libp2p-peerstore/pstoreds"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	// maxReceiveBatch is the maximum number of new messages to receive at once.
	maxReceiveBatch = 500
	// maxShareBatch is the maximum number of messages to share at once.
	maxShareBatch = 100
	// peerCountLow is the target number of peers to connect to at any given time.
	peerCountLow = 100
	// peerCountHigh is the maximum number of peers to be connected to. If the
	// number of connections exceeds this number, we will prune connections until
	// we reach peerCountLow.
	peerCountHigh = 110
)

func getHostOptions(ctx context.Context, config Config) ([]libp2p.Option, error) {
	// Note: 0.0.0.0 will use all available addresses.
	tcpBindAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", config.TCPPort))
	if err != nil {
		return nil, err
	}
	wsBindAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/ws", config.WebSocketsPort))
	if err != nil {
		return nil, err
	}

	// HACK(albrow): As a workaround for AutoNAT issues, ping ifconfig.me to
	// determine our public IP address on boot. This will work for nodes that
	// would be reachable via a public IP address but don't know what it is (e.g.
	// because they are running in a Docker container).
	publicIP, err := getPublicIP()
	if err != nil {
		return nil, fmt.Errorf("could not get public IP address: %s", err.Error())
	}
	tcpAdvertiseAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", publicIP, config.TCPPort))
	if err != nil {
		return nil, err
	}
	wsAdvertiseAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/ws", publicIP, config.WebSocketsPort))
	if err != nil {
		return nil, err
	}
	advertiseAddrs := []ma.Multiaddr{tcpAdvertiseAddr, wsAdvertiseAddr}

	// Set up the peerstore to use LevelDB.
	store, err := leveldbStore.NewDatastore(getPeerstoreDir(config.DataDir), nil)
	if err != nil {
		return nil, err
	}
	pstore, err := pstoreds.NewPeerstore(ctx, store, pstoreds.DefaultOpts())
	if err != nil {
		return nil, err
	}

	// Set up the WebSocket transport to ignore TLS verification. We use secio so
	// it is not necessary.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	newWebsocketTransport := ws.NewWithOptions(ws.TLSClientConfig(tlsConfig))

	return []libp2p.Option{
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(newWebsocketTransport),
		libp2p.ListenAddrs(tcpBindAddr, wsBindAddr),
		libp2p.AddrsFactory(newAddrsFactory(advertiseAddrs)),
		libp2p.Peerstore(pstore),
	}, nil
}

func getPubSubOptions() []pubsub.Option {
	// Use the default options.
	return nil
}

func newAddrsFactory(advertiseAddrs []ma.Multiaddr) func([]ma.Multiaddr) []ma.Multiaddr {
	return func(addrs []ma.Multiaddr) []ma.Multiaddr {
		// Note that we append the advertiseAddrs here just in case we are not
		// actually reachable at our public IP address (and are reachable at one of
		// the other addresses).
		return append(addrs, advertiseAddrs...)
	}
}

func getPublicIP() (string, error) {
	res, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	ipBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(ipBytes), nil
}

// NewDHT returns a new Kademlia DHT instance configured to work with 0x Mesh
// in native (pure Go) environments. storageDir is the directory to use for
// persisting the data with LevelDB.
func NewDHT(ctx context.Context, storageDir string, host host.Host) (*dht.IpfsDHT, error) {
	// Set up the DHT to use LevelDB.
	store, err := leveldbStore.NewDatastore(storageDir, nil)
	if err != nil {
		return nil, err
	}

	return dht.New(ctx, host, dhtopts.Datastore(store), dhtopts.Protocols(DHTProtocolID))
}
