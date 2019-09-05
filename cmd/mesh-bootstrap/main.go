// +build !js

// mesh-bootstrap is a separate executable for bootstrap nodes. Bootstrap nodes
// will not share or receive any orders and its sole responsibility is to
// facilitate peer discovery.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/keys"
	"github.com/0xProject/0x-mesh/loghooks"
	"github.com/0xProject/0x-mesh/p2p"
	libp2p "github.com/libp2p/go-libp2p"
	autonat "github.com/libp2p/go-libp2p-autonat-svc"
	circuit "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	p2pnet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/host/relay"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

const (
	// peerCountLow is the target number of peers to connect to at any given time.
	peerCountLow = 1000
	// peerCountHigh is the maximum number of peers to be connected to. If the
	// number of connections exceeds this number, we will prune connections until
	// we reach peerCountLow.
	peerCountHigh = 1100
	// peerGraceDuration is the amount of time a newly opened connection is given
	// before it becomes subject to pruning.
	peerGraceDuration = 10 * time.Second
	// defaultNetworkTimeout is the default timeout for network requests (e.g.
	// connecting to a new peer).
	defaultNetworkTimeout = 30 * time.Second
)

// Config contains configuration options for a Node.
type Config struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"5"`
	// P2PBindAddrs is a comma separated list of libp2p multiaddresses which the
	// bootstrap node will bind to.
	P2PBindAddrs string `envvar:"P2P_BIND_ADDRS"`
	// P2PAdvertiseAddrs is a comma separated list of libp2p multiaddresses which the
	// bootstrap node will advertise to peers.
	P2PAdvertiseAddrs string `envvar:"P2P_ADVERTISE_ADDRS"`
	// DataDir is the directory used for storing data.
	DataDir string `envvar:"DATA_DIR" default:"0x_mesh"`
	// BootstrapList is a comma-separated list of multiaddresses to use for
	// bootstrapping the DHT (e.g.,
	// "/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF").
	// If empty, the default bootstrap list will be used.
	BootstrapList string `envvar:"BOOTSTRAP_LIST" default:""`
}

func init() {
	// Since we know that the bootstrap nodes are more stable, we can
	// safely reduce AdvertiseBootDelay. This will allow the bootstrap nodes to
	// advertise themselves as relays sooner.
	relay.AdvertiseBootDelay = 30 * time.Second
}

func getPrivateKeyPath(config Config) string {
	return filepath.Join(config.DataDir, "keys", "privkey")
}

func getDHTDir(config Config) string {
	return filepath.Join(config.DataDir, "p2p", "dht")
}

func getPeerstoreDir(config Config) string {
	return filepath.Join(config.DataDir, "p2p", "peerstore")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse env vars
	var config Config
	if err := envvar.Parse(&config); err != nil {
		panic(fmt.Sprintf("could not parse environment variables: %s", err.Error()))
	}

	// Configure logger to output JSON
	// TODO(albrow): Don't use global settings for logger.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(config.Verbosity))
	log.AddHook(loghooks.NewKeySuffixHook())

	// Parse private key file and add peer ID log hook
	privKey, err := initPrivateKey(getPrivateKeyPath(config))
	if err != nil {
		log.WithField("error", err).Fatal("could not initialize private key")
	}
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		log.Fatal(err)
	}
	log.AddHook(loghooks.NewPeerIDHook(peerID))

	// We need to declare the newDHT function ahead of time so we can use it in
	// the libp2p.Routing option.
	var kadDHT *dht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dhtDir := getDHTDir(config)
		kadDHT, err = p2p.NewDHT(ctx, dhtDir, h)
		if err != nil {
			log.WithField("error", err).Fatal("could not create DHT")
		}
		return kadDHT, err
	}

	// Parse multiaddresses given in the config
	bindAddrs, err := parseAddrs(config.P2PBindAddrs)
	if err != nil {
		log.Fatal(err)
	}
	advertiseAddrs, err := parseAddrs(config.P2PAdvertiseAddrs)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the transport and the host.
	connManager := connmgr.NewConnManager(peerCountLow, peerCountHigh, peerGraceDuration)
	opts := []libp2p.Option{
		libp2p.ListenAddrs(bindAddrs...),
		libp2p.Identity(privKey),
		libp2p.ConnectionManager(connManager),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.EnableAutoRelay(),
		libp2p.Routing(newDHT),
		libp2p.AddrsFactory(newAddrsFactory(advertiseAddrs)),
	}
	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		log.WithField("error", err).Fatal("could not create host")
	}

	// Set up the notifee.
	basicHost.Network().Notify(&notifee{})

	// Enable AutoNAT service.
	if _, err := autonat.NewAutoNATService(ctx, basicHost); err != nil {
		log.WithField("error", err).Fatal("could not enable AutoNAT service")
	}

	// Initialize the DHT and then connect to the other bootstrap nodes.
	if err := kadDHT.Bootstrap(ctx); err != nil {
		log.WithField("error", err).Fatal("could not bootstrap DHT")
	}
	bootstrapList := p2p.DefaultBootstrapList
	if config.BootstrapList != "" {
		bootstrapList = strings.Split(config.BootstrapList, ",")
	}
	if err := p2p.ConnectToBootstrapList(ctx, basicHost, bootstrapList); err != nil {
		log.WithField("error", err).Fatal("could not connect to bootstrap peers")
	}

	// Protect each other bootstrap peer via the connection manager so that we
	// maintain an active connection to them.
	bootstrapAddrInfos, err := p2p.BootstrapListToAddrInfos(bootstrapList)
	if err != nil {
		log.WithField("error", err).Fatal("could not parse bootstrap list")
	}
	for _, addrInfo := range bootstrapAddrInfos {
		connManager.Protect(addrInfo.ID, "bootstrap-peer")
	}

	log.WithFields(map[string]interface{}{
		"addrs":  basicHost.Addrs(),
		"config": config,
	}).Info("started bootstrap node")

	// Sleep until stopped
	select {}
}

func initPrivateKey(path string) (p2pcrypto.PrivKey, error) {
	privKey, err := keys.GetPrivateKeyFromPath(path)
	if err == nil {
		return privKey, nil
	} else if os.IsNotExist(err) {
		// If the private key doesn't exist, generate one.
		log.Info("No private key found. Generating a new one.")
		return keys.GenerateAndSavePrivateKey(path)
	}

	// For any other type of error, return it.
	return nil, err
}

// notifee receives notifications for network-related events.
type notifee struct{}

var _ p2pnet.Notifiee = &notifee{}

// Listen is called when network starts listening on an addr
func (n *notifee) Listen(p2pnet.Network, ma.Multiaddr) {}

// ListenClose is called when network stops listening on an addr
func (n *notifee) ListenClose(p2pnet.Network, ma.Multiaddr) {}

// Connected is called when a connection opened
func (n *notifee) Connected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"remotePeerID":       conn.RemotePeer(),
		"remoteMultiaddress": conn.RemoteMultiaddr(),
	}).Info("connected to peer")
}

// Disconnected is called when a connection closed
func (n *notifee) Disconnected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"remotePeerID":       conn.RemotePeer(),
		"remoteMultiaddress": conn.RemoteMultiaddr(),
	}).Info("disconnected from peer")
}

// OpenedStream is called when a stream opened
func (n *notifee) OpenedStream(network p2pnet.Network, stream p2pnet.Stream) {}

// ClosedStream is called when a stream closed
func (n *notifee) ClosedStream(network p2pnet.Network, stream p2pnet.Stream) {}

func newAddrsFactory(advertiseAddrs []ma.Multiaddr) func([]ma.Multiaddr) []ma.Multiaddr {
	return func([]ma.Multiaddr) []ma.Multiaddr {
		return advertiseAddrs
	}
}

func parseAddrs(commaSeparatedAddrs string) ([]ma.Multiaddr, error) {
	maddrStrings := strings.Split(commaSeparatedAddrs, ",")
	maddrs := make([]ma.Multiaddr, len(maddrStrings))
	for i, maddrString := range maddrStrings {
		ma, err := ma.NewMultiaddr(maddrString)
		if err != nil {
			return nil, err
		}
		maddrs[i] = ma
	}
	return maddrs, nil
}
