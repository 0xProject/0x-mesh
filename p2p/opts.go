// +build !js

package p2p

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/0xProject/0x-mesh/db"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoreds"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
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
	publicIP, err := getPublicIP(config.AdditionalPublicIPSources)
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

	pstore, err := pstoreds.NewPeerstore(ctx, config.DB.PeerStore(), pstoreds.DefaultOpts())
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

func fetchPublicIPFromExternalSource(source string) (string, error) {
	res, err := http.Get(source)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	ipBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(ipBytes), "\n"), nil
}

func getPublicIP(additionalSources []string) (string, error) {
	sources := []string{"https://wtfismyip.com/text", "https://whatismyip.api.0x.org", "https://ifconfig.me/ip"}
	sources = append(additionalSources, sources...)
	for _, source := range sources {
		ip, err := fetchPublicIPFromExternalSource(source)
		if err != nil {
			log.WithField("source", source).Warn("failed to get public ip from source")
			continue
		}

		return ip, nil

	}

	return "", errors.New("failed to get public ip from all provided external sources")
}

// NewDHT returns a new Kademlia DHT instance configured to work with 0x Mesh
// in native (pure Go) environments. Standalone nodes use a SQL key value store
// to persist data and browser nodes use a Dexie key value store.
func NewDHT(ctx context.Context, db *db.DB, host host.Host) (*dht.IpfsDHT, error) {
	return dht.New(ctx, host, dht.Datastore(db.DHTStore()), dht.V1ProtocolOverride(DHTProtocolID), dht.Mode(dht.ModeAutoServer))
}
