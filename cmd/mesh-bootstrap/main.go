// +build !js

// mesh-bootstrap is a separate executable for bootstrap nodes. Bootstrap nodes
// will not share or receive any orders and its sole responsibility is to
// facilitate peer discovery.
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/loghooks"
	"github.com/0xProject/0x-mesh/p2p"
	libp2p "github.com/libp2p/go-libp2p"
	autonat "github.com/libp2p/go-libp2p-autonat-svc"
	circuit "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/host"
	p2pnet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/host/relay"
	ma "github.com/multiformats/go-multiaddr"
	sqlds "github.com/opaolini/go-ds-sql"
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
	// DataStoreType constants
	leveldbDataStore  = "leveldb"
	postgresDataStore = "postgres"
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
	// DataStoreType is the data store which will be used to store DHT data
	// for the bootstrap node.
	// DataStoreType can be either: leveldb or postgres
	DataStoreType string `envvar:"DATA_STORE_TYPE" default:"leveldb"`
	// DataDir is the directory used for storing data when using leveldb as data store type.
	DataDir string `envvar:"DATA_DIR" default:"0x_mesh"`
	// DataDBHost is the database host used to connect to the database when
	// using postgres as data store type.
	DataDBHost string `envvar:"DATA_DB_HOST" default:"localhost"`
	// DataDBPort is the database port used to connect to the database when
	// using postgres as data store type.
	DataDBPort string `envvar:"DATA_DB_PORT" default:"5432"`
	// DataDBUser is the database user used to connect to the database when
	// using postgres as data store type.
	DataDBUser string `envvar:"DATA_DB_USER" default:"postgres"`
	// DataDBPassword is the database password used to connect to the database when
	// using postgres as data store type.
	DataDBPassword string `envvar:"DATA_DB_PASSWORD" default:""`
	// DataDBDatabaseName is the database name to connect to when using
	// postgres as data store type.
	DataDBDatabaseName string `envvar:"DATA_DB_NAME" default:"datastore"`
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
	var newDHT func(h host.Host) (routing.PeerRouting, error)

	// TODO(oskar) - Figure out why returning an anonymous function from
	// getNewDHT() is making kadDHT.runBootstrap panicing.
	// When solved this big switch case can be removed from main()
	switch config.DataStoreType {
	case leveldbDataStore:
		newDHT = func(h host.Host) (routing.PeerRouting, error) {
			var err error
			dhtDir := getDHTDir(config)
			kadDHT, err = p2p.NewDHT(ctx, dhtDir, h)
			if err != nil {
				log.WithField("error", err).Fatal("could not create DHT")
			}
			return kadDHT, err
		}

	case postgresDataStore:
		newDHT = func(h host.Host) (routing.PeerRouting, error) {
			var err error
			sqlOpts := &sqlds.Options{
				Host:     config.DataDBHost,
				Port:     config.DataDBPort,
				User:     config.DataDBUser,
				Password: config.DataDBPassword,
				Database: config.DataDBDatabaseName,
				Table:    "dhtkv",
			}
			store, err := sqlOpts.CreatePostgres()
			if err != nil {
				log.WithField("error", err).Fatal("could not create postgres datastore")
			}

			kadDHT, err = p2p.NewDHTWithDatastore(ctx, store, h)
			if err != nil {
				log.WithField("error", err).Fatal("could not create DHT")
			}

			return kadDHT, err
		}

	default:
		log.Fatalf("invalid datastore configured: %s. Expected either %s or %s", config.DataStoreType, leveldbDataStore, postgresDataStore)

	}

	pstore, err := getNewPeerstore(ctx, config)
	if err != nil {
		log.Fatal(err)
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
		libp2p.Peerstore(pstore),
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
