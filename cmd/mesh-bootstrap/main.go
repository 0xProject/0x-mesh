// +build !js

// mesh-bootstrap is a separate executable for bootstrap nodes. Bootstrap nodes
// will not share or receive any orders and its sole responsibility is to
// facilitate peer discovery.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/0xProject/0x-mesh/keys"
	"github.com/0xProject/0x-mesh/loghooks"
	"github.com/0xProject/0x-mesh/p2p"
	libp2p "github.com/libp2p/go-libp2p"
	autonat "github.com/libp2p/go-libp2p-autonat-svc"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	p2pcrypto "github.com/libp2p/go-libp2p-crypto"
	p2pnet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
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
	Verbosity int `envvar:"VERBOSITY" default:"6"`
	// P2PListenPort is the port on which to listen for new connections. It can be
	// set to 0 to make the OS automatically choose any available port.
	P2PListenPort int `envvar:"P2P_LISTEN_PORT" default:"0"`
	// PrivateKey path is the path to a private key file which will be used for
	// signing messages and generating a peer ID.
	PrivateKeyPath string `envvar:"PRIVATE_KEY_PATH" default:"0x_mesh/keys/privkey"`
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

	// Parse private key file
	privKey, err := initPrivateKey(config.PrivateKeyPath)
	if err != nil {
		log.WithField("error", err).Fatal("could not initialize private key")
	}

	// Set up the transport and the host.
	// Note: 0.0.0.0 will use all available addresses.
	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", config.P2PListenPort))
	if err != nil {
		log.WithField("error", err).Fatal("could not parse multiaddr")
	}
	connManager := connmgr.NewConnManager(peerCountLow, peerCountHigh, peerGraceDuration)
	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(privKey),
		libp2p.ConnectionManager(connManager),
	}
	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		log.WithField("error", err).Fatal("could not create host")
	}

	// Add the peer ID hook to the logger.
	log.AddHook(loghooks.NewPeerIDHook(basicHost.ID()))

	// Set up the notifee.
	basicHost.Network().Notify(&notifee{})

	// Enable AutoNAT service.
	if _, err := autonat.NewAutoNATService(ctx, basicHost); err != nil {
		log.WithField("error", err).Fatal("could not enable AutoNAT service")
	}

	// Set up DHT for peer discovery.
	kadDHT, err := p2p.NewDHT(ctx, basicHost)
	if err != nil {
		log.WithField("error", err).Fatal("could not create DHT")
	}

	// Initialize the DHT and then connect to the other bootstrap nodes.
	if err := kadDHT.Bootstrap(ctx); err != nil {
		log.WithField("error", err).Fatal("could not bootstrap DHT")
	}
	if err := p2p.ConnectToBootstrapList(ctx, basicHost); err != nil {
		log.WithField("error", err).Fatal("could not connect to bootstrap peers")
	}

	// Protect each other bootstrap peer via the connection manager so that we
	// maintain an active connection to them.
	for _, addr := range p2p.BootstrapPeers {
		idString, err := addr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.WithField("error", err).Fatal("could not extract peer id from bootstrap peer")
		}
		id, err := peer.IDB58Decode(idString)
		if err != nil {
			log.WithField("error", err).Fatal("could not extract peer id from bootstrap peer")
		}
		connManager.Protect(id, "bootstrap-peer")
	}

	log.WithFields(map[string]interface{}{
		"addrs":  basicHost.Addrs(),
		"peerID": basicHost.ID(),
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
		"peerID":       conn.RemotePeer(),
		"multiaddress": conn.RemoteMultiaddr(),
	}).Info("connected to peer")
}

// Disconnected is called when a connection closed
func (n *notifee) Disconnected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"peerID":       conn.RemotePeer(),
		"multiaddress": conn.RemoteMultiaddr(),
	}).Info("disconnected from peer")
}

// OpenedStream is called when a stream opened
func (n *notifee) OpenedStream(network p2pnet.Network, stream p2pnet.Stream) {}

// ClosedStream is called when a stream closed
func (n *notifee) ClosedStream(network p2pnet.Network, stream p2pnet.Stream) {}
