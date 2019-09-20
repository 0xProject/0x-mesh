// package p2p is a low-level library responsible for peer discovery and
// sending/receiving messages.
package p2p

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	mathrand "math/rand"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/albrow/stringset"
	lru "github.com/hashicorp/golang-lru"
	libp2p "github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	swarm "github.com/libp2p/go-libp2p-swarm"
	filter "github.com/libp2p/go-maddr-filter"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	// receiveTimeout is the maximum amount of time to wait for receiving new messages.
	receiveTimeout = 1 * time.Second
	// peerGraceDuration is the amount of time a newly opened connection is given
	// before it becomes subject to pruning.
	peerGraceDuration = 10 * time.Second
	// defaultNetworkTimeout is the default timeout for network requests (e.g.
	// connecting to a new peer).
	defaultNetworkTimeout = 10 * time.Second
	// advertiseTTL is the TTL for our announcement to the discovery network.
	advertiseTTL = 5 * time.Minute
	// pubsubProtocolID is the protocol ID to use for pubsub.
	// TODO(albrow): Is there a way to use a custom protocol ID with GossipSub?
	// pubsubProtocolID = protocol.ID("/0x-mesh-gossipsub/0.0.1")
	pubsubProtocolID = pubsub.GossipSubID
	// chanceToCheckBandwidthUsage is the approximate ratio of (number of main
	// loop iterations in which we check bandwidth usage) to (total number of main
	// loop iterations). We check bandwidth non-deterministically in order to
	// prevent spammers from avoiding detection by carefully timing their
	// bandwidth usage. So on each iteration of the main loop we generate a number
	// between 0 and 1. If its less than chanceToCheckBandiwdthUsage, we perform
	// a bandwidth check.
	chanceToCheckBandiwdthUsage = 0.1
	// logBandwidthUsageInterval is how often to log bandwidth usage data.
	logBandwidthUsageInterval = 5 * time.Minute
)

var errProtectedIP = errors.New("cannot ban protected IP address")

// Node is the main type for the p2p package. It represents a particpant in the
// 0x Mesh network who is capable of sending, receiving, validating, and storing
// messages.
type Node struct {
	ctx              context.Context
	config           Config
	messageHandler   MessageHandler
	host             host.Host
	filters          *filter.Filters
	connManager      *connmgr.BasicConnMgr
	dht              *dht.IpfsDHT
	routingDiscovery discovery.Discovery
	pubsub           *pubsub.PubSub
	sub              *pubsub.Subscription
	protectedIPsMut  sync.RWMutex
	protectedIPs     stringset.Set
	bandwidthChecker *bandwidthChecker
}

// Config contains configuration options for a Node.
type Config struct {
	// Topic is a unique string representing the pubsub topic. Only Nodes which
	// have the same topic will share messages with one another.
	Topic string
	// TCPPort is the port on which to listen for incoming TCP connections.
	TCPPort int
	// WebSocketsPort is the port on which to listen for incoming WebSockets
	// connections.
	WebSocketsPort int
	// Insecure controls whether or not messages should be encrypted. It should
	// always be set to false in production.
	Insecure bool
	// PrivateKey is the private key which will be used for signing messages and
	// generating a peer ID.
	PrivateKey p2pcrypto.PrivKey
	// MessageHandler is an interface responsible for validating, storing, and
	// finding new messages to share.
	MessageHandler MessageHandler
	// RendezvousString is a unique identifier for the rendezvous point. This node
	// will attempt to find peers with the same Rendezvous string.
	RendezvousString string
	// UseBootstrapList determines whether or not to use the list of hard-coded
	// peers to bootstrap the DHT for peer discovery.
	UseBootstrapList bool
	// BootstrapList is a list of multiaddress strings to use for bootstrapping
	// the DHT. If empty, the default list will be used.
	BootstrapList []string
	// DataDir is the directory to use for storing data.
	DataDir string
}

func getPeerstoreDir(datadir string) string {
	return filepath.Join(datadir, "peerstore")
}

func getDHTDir(datadir string) string {
	return filepath.Join(datadir, "dht")
}

// New creates a new Node with the given context and config. The Node will stop
// all background operations if the context is canceled.
func New(ctx context.Context, config Config) (*Node, error) {
	if config.MessageHandler == nil {
		return nil, errors.New("config.MessageHandler is required")
	} else if config.RendezvousString == "" {
		return nil, errors.New("config.RendezvousString is required")
	}

	// We need to declare the newDHT function ahead of time so we can use it in
	// the libp2p.Routing option.
	var kadDHT *dht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dhtDir := getDHTDir(config.DataDir)
		kadDHT, err = NewDHT(ctx, dhtDir, h)
		if err != nil {
			log.WithField("error", err).Error("could not create DHT")
		}
		return kadDHT, err
	}

	// Get environment specific host options.
	opts, err := getHostOptions(ctx, config)
	if err != nil {
		return nil, err
	}

	// Initialize filters.
	filters := filter.NewFilters()

	// Set up and append environment agnostic host options.
	bandwidthCounter := metrics.NewBandwidthCounter()
	connManager := connmgr.NewConnManager(peerCountLow, peerCountHigh, peerGraceDuration)
	opts = append(opts, []libp2p.Option{
		libp2p.Routing(newDHT),
		libp2p.ConnectionManager(connManager),
		libp2p.Identity(config.PrivateKey),
		libp2p.EnableAutoRelay(),
		libp2p.EnableRelay(),
		libp2p.BandwidthReporter(bandwidthCounter),
		Filters(filters),
	}...)
	if config.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	}

	// Initialize the host.
	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// Close the host whenever the context is canceled.
	go func() {
		<-ctx.Done()
		_ = basicHost.Close()
	}()

	// Set up the notifee.
	basicHost.Network().Notify(&notifee{
		ctx:         ctx,
		connManager: connManager,
	})

	// Set up DHT for peer discovery.
	routingDiscovery := discovery.NewRoutingDiscovery(kadDHT)

	// Set up pubsub
	pubsubOpts := getPubSubOptions()
	pubsub, err := pubsub.NewGossipSub(ctx, basicHost, pubsubOpts...)
	if err != nil {
		return nil, err
	}

	// Create the Node.
	node := &Node{
		ctx:              ctx,
		config:           config,
		messageHandler:   config.MessageHandler,
		host:             basicHost,
		filters:          filters,
		connManager:      connManager,
		dht:              kadDHT,
		routingDiscovery: routingDiscovery,
		pubsub:           pubsub,
		protectedIPs:     stringset.New(),
	}
	node.bandwidthChecker = newBandwidthChecker(node, bandwidthCounter)

	return node, nil
}

func getPrivateKey(path string) (p2pcrypto.PrivKey, error) {
	if path == "" {
		// If path is empty, generate a new key.
		priv, _, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
		if err != nil {
			return nil, err
		}
		return priv, nil
	}

	// Otherwise parse the key at the path given.
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decodedKey, err := p2pcrypto.ConfigDecodeKey(string(keyBytes))
	if err != nil {
		return nil, err
	}
	priv, err := p2pcrypto.UnmarshalPrivateKey(decodedKey)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// Multiaddrs returns all multi addresses at which the node is dialable.
func (n *Node) Multiaddrs() []ma.Multiaddr {
	return n.host.Addrs()
}

// ID returns the peer id of the node.
func (n *Node) ID() peer.ID {
	return n.host.ID()
}

// Start causes the Node to continuously send messages to and receive messages
// from its peers. It blocks until an error is encountered or `Stop` is called.
func (n *Node) Start() error {
	// Note: dht.Bootstrap has a somewhat confusing name. It doesn't automatically
	// connect to the bootstrap peers. It just starts the background process of
	// searching for new peers.
	if err := n.dht.Bootstrap(n.ctx); err != nil {
		return err
	}

	// Use the default bootstrap list if none was provided.
	if len(n.config.BootstrapList) == 0 {
		n.config.BootstrapList = DefaultBootstrapList
	}

	// If needed, connect to all peers in the bootstrap list.
	if n.config.UseBootstrapList {
		if err := ConnectToBootstrapList(n.ctx, n.host, n.config.BootstrapList); err != nil {
			return err
		}
		// Protect the IP addresses for each bootstrap node.
		bootstrapAddrInfos, err := BootstrapListToAddrInfos(n.config.BootstrapList)
		if err != nil {
			return err
		}
		for _, addrInfo := range bootstrapAddrInfos {
			for _, addr := range addrInfo.Addrs {
				_ = n.ProtectIP(addr)
			}
		}
	}

	// Advertise ourselves for the purposes of peer discovery.
	discovery.Advertise(n.ctx, n.routingDiscovery, n.config.RendezvousString, discovery.TTL(advertiseTTL))

	// Start logging bandwidth in the background.
	go n.bandwidthChecker.logBandwidthUsage(n.ctx)

	return n.mainLoop()
}

// AddPeerScore adds diff to the current score for a given peer. Tag is a unique
// identifier for the score. A peer's total score is the sum of the scores
// associated with each tag. Peers that end up with a low total score will
// eventually be disconnected.
func (n *Node) AddPeerScore(id peer.ID, tag string, diff int) {
	n.connManager.UpsertTag(id, tag, func(current int) int { return current + diff })
}

// SetPeerScore sets the current score for a given peer (overwriting any
// previous value with the same tag). Tag is a unique identifier for the score.
// A peer's total score is the sum of the scores associated with each tag. Peers
// that end up with a low total score will eventually be disconnected.
func (n *Node) SetPeerScore(id peer.ID, tag string, val int) {
	n.connManager.TagPeer(id, tag, val)
}

// UnsetPeerScore removes any scores associated with the given tag for a peer
// (i.e., they will no longer be counted toward the peers total score).
func (n *Node) UnsetPeerScore(id peer.ID, tag string) {
	n.connManager.UntagPeer(id, tag)
}

// GetNumPeers returns the number of peers the node is connected to
func (n *Node) GetNumPeers() int {
	return n.connManager.GetInfo().ConnCount
}

// Connect ensures there is a connection between this host and the peer with
// given peerInfo. If there is not an active connection, Connect will dial the
// peer, and block until a connection is open, timeout is exceeded, or an error
// is returned.
func (n *Node) Connect(peerInfo peer.AddrInfo, timeout time.Duration) error {
	connectCtx, cancel := context.WithTimeout(n.ctx, timeout)
	defer cancel()
	err := n.host.Connect(connectCtx, peerInfo)
	if err != nil {
		return err
	}
	return nil
}

// ProtectIP permanently adds the IP address of the given Multiaddr to a
// list of protected IP addresses. Protected IPs can never be banned and will
// not be added to the blacklist. If the IP address is already on the blacklist,
// it will be removed.
func (n *Node) ProtectIP(maddr ma.Multiaddr) error {
	n.protectedIPsMut.Lock()
	defer n.protectedIPsMut.Unlock()
	ipNet, err := ipNetFromMaddr(maddr)
	if err != nil {
		return err
	}
	n.unbanIPNet(ipNet)
	n.protectedIPs.Add(ipNet.IP.String())
	return nil
}

// BanIP adds the IP address of the given Multiaddr to the blacklist. The
// node will no longer dial or accept connections from this IP address. However,
// if the IP address is protected, calling BanIP will not ban the IP address and
// will instead return errProtectedIP. BanIP does not automatically disconnect
// from the given multiaddress if there is currently an open connection.
func (n *Node) BanIP(maddr ma.Multiaddr) error {
	ipNet, err := ipNetFromMaddr(maddr)
	if err != nil {
		// HACK(albrow) relay addresses don't include the full transport address
		// (IP, port, etc) for older versions of libp2p-circuit. (See
		// https://github.com/libp2p/go-libp2p/issues/723). As a temporary
		// workaround, we no-op for relayed connections. We can remove this after
		// updating our bootstrap nodes to the latest version. We detect relay
		// addresses by looking for the /ipfs prefix.
		if strings.HasPrefix(maddr.String(), "/ipfs") {
			return nil
		}
		return err
	}
	n.protectedIPsMut.RLock()
	defer n.protectedIPsMut.RUnlock()
	if n.protectedIPs.Contains(ipNet.IP.String()) {
		// IP address is protected. no-op.
		return errProtectedIP
	}
	n.filters.AddFilter(ipNet, filter.ActionDeny)
	return nil
}

// UnbanIP removes the IP address of the given Multiaddr from the blacklist. If
// the IP address is not currently on the blacklist this is a no-op.
func (n *Node) UnbanIP(maddr ma.Multiaddr) error {
	ipNet, err := ipNetFromMaddr(maddr)
	if err != nil {
		return err
	}
	n.unbanIPNet(ipNet)
	return nil
}

func (n *Node) unbanIPNet(ipNet net.IPNet) {
	// There is no guarantee in the public API of the filters package that would
	// prevent multiple filters being added for the same IPNet (though it
	// shouldn't happen in practice). We use a for loop here to make sure we
	// remove all possible filters. RemoveLiteral returns false if no filter was
	// removed.
	for n.filters.RemoveLiteral(ipNet) {
	}
}

// mainLoop is where the core logic for a Node is implemented. On each iteration
// of the loop, the node receives new messages and sends messages to its peers.
func (n *Node) mainLoop() error {
	for {
		select {
		case <-n.ctx.Done():
			return nil
		default:
		}
		if err := n.runOnce(); err != nil {
			return err
		}
	}
}

// runOnce runs a single iteration of the main loop.
func (n *Node) runOnce() error {
	peerCount := n.connManager.GetInfo().ConnCount
	if peerCount < peerCountLow {
		if err := n.findNewPeers(peerCountLow - peerCount); err != nil {
			return err
		}
	}

	// Receive up to maxReceiveBatch messages.
	incoming, err := n.receiveBatch()
	if err != nil {
		return err
	}
	if err := n.messageHandler.HandleMessages(incoming); err != nil {
		return fmt.Errorf("could not validate or store messages: %s", err.Error())
	}

	// Check bandwidth usage non-deterministically
	if mathrand.Float64() <= chanceToCheckBandiwdthUsage {
		n.bandwidthChecker.checkUsage()
	}

	// Send up to maxSendBatch messages.
	if err := n.shareBatch(); err != nil {
		return err
	}
	return nil
}

func (n *Node) findNewPeers(max int) error {
	log.WithFields(map[string]interface{}{
		"max": max,
	}).Trace("looking for new peers")
	findPeersCtx, cancel := context.WithTimeout(n.ctx, defaultNetworkTimeout)
	defer cancel()
	peerChan, err := n.routingDiscovery.FindPeers(findPeersCtx, n.config.RendezvousString, discovery.Limit(max))
	if err != nil {
		return err
	}

	connectCtx, cancel := context.WithTimeout(n.ctx, defaultNetworkTimeout)
	defer cancel()
	for peer := range peerChan {
		if peer.ID == n.host.ID() || len(peer.Addrs) == 0 {
			continue
		}
		log.WithFields(map[string]interface{}{
			"peerInfo": peer,
		}).Trace("found peer via rendezvous")
		if err := n.host.Connect(connectCtx, peer); err != nil {
			// We still want to try connecting to the other peers. Log the error and
			// keep going.
			logPeerConnectionError(peer, err)
		}
	}
	return nil
}

// failedPeerConnectionCache keeps track of peer IDs for which we have already
// logged a connection error. lru.New only returns an error if size is <= 0, so
// we can safely ignore it.
var failedPeerConnectionCache, _ = lru.New(peerCountHigh * 2)

func logPeerConnectionError(peerInfo peer.AddrInfo, connectionErr error) {
	// If we fail to connect to a single peer we should still keep trying the
	// others. Log instead of returning the error.
	logMsg := "could not connect to peer"
	logFields := map[string]interface{}{
		"error":        connectionErr.Error(),
		"remotePeerID": peerInfo.ID,
		"remoteAddrs":  peerInfo.Addrs,
	}
	if failedPeerConnectionCache.Contains(peerInfo.ID) {
		// If we have already logged a connection error for this peer ID, log at
		// level "trace".
		log.WithFields(logFields).Trace(logMsg)
	} else if connectionErr == swarm.ErrDialBackoff {
		// ErrDialBackoff means that we dialed the peer too frequently. Logging
		// it leads to too much verbosity and in most cases what we care about
		// is the underlying error. Log at level "trace".
		log.WithFields(logFields).Trace(logMsg)
	} else {
		// For other types of errors, and if we have not already logged a connection
		// error for this peer ID, we log at level "warn".
		failedPeerConnectionCache.Add(peerInfo.ID, nil)
		log.WithFields(logFields).Warn(logMsg)
	}
}

// receiveBatch returns up to maxReceiveBatch messages which are received from
// peers. There is no guarantee that the messages are unique.
func (n *Node) receiveBatch() ([]*Message, error) {
	messages := []*Message{}
	for {
		if len(messages) >= maxReceiveBatch {
			return messages, nil
		}
		select {
		// If the parent context was canceled, return.
		case <-n.ctx.Done():
			return messages, nil
		default:
		}
		receiveCtx, receiveCancel := context.WithTimeout(n.ctx, receiveTimeout)
		msg, err := n.receive(receiveCtx)
		receiveCancel()
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return messages, nil
			}
			return nil, err
		}
		if msg.From == n.host.ID() {
			continue
		}
		messages = append(messages, msg)
	}
}

// shareBatch shares up to maxShareBatch messages (selected via the
// MessageHandler) with all connected peers.
func (n *Node) shareBatch() error {
	// TODO(albrow): This will need to change when we switch to WeijieSub.
	outgoing, err := n.messageHandler.GetMessagesToShare(maxShareBatch)
	if err != nil {
		return err
	}
	for _, data := range outgoing {
		if err := n.send(data); err != nil {
			return err
		}
	}
	return nil
}

// send sends a message continaing the given data to all connected peers.
func (n *Node) send(data []byte) error {
	return n.pubsub.Publish(n.config.Topic, data)
}

// receive returns the next pending message. It blocks if no messages are
// available. If the given context is canceled, it returns nil, ctx.Err().
func (n *Node) receive(ctx context.Context) (*Message, error) {
	if n.sub == nil {
		var err error
		n.sub, err = n.pubsub.Subscribe(n.config.Topic)
		if err != nil {
			return nil, err
		}
	}
	msg, err := n.sub.Next(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{From: msg.GetFrom(), Data: msg.Data}, nil
}

func ipNetFromMaddr(maddr ma.Multiaddr) (ipNet net.IPNet, err error) {
	ip, err := ipFromMaddr(maddr)
	if err != nil {
		return net.IPNet{}, err
	}
	mask := getAllMaskForIP(ip)
	return net.IPNet{
		IP:   ip,
		Mask: mask,
	}, nil
}

func ipFromMaddr(maddr ma.Multiaddr) (net.IP, error) {
	var (
		ip    net.IP
		found bool
	)

	ma.ForEach(maddr, func(c ma.Component) bool {
		switch c.Protocol().Code {
		case ma.P_IP6ZONE:
			return true
		case ma.P_IP6, ma.P_IP4:
			found = true
			ip = net.IP(c.RawValue())
			return false
		default:
			return false
		}
	})

	if !found {
		return net.IP{}, fmt.Errorf("could not parse IP address from multiaddress: %s", maddr)
	}
	return ip, nil
}

var (
	ipv4AllMask = net.IPMask{255, 255, 255, 255}
	ipv6AllMask = net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
)

// getAllMaskForIP returns an IPMask that will match all IP addresses. The size
// of the mask depends on whether the given IP address is an IPv4 or an IPv6
// address.
func getAllMaskForIP(ip net.IP) net.IPMask {
	if ip.To4() != nil {
		// This is an ipv4 address. Return 4 byte mask.
		return ipv4AllMask
	}
	// Assume ipv6 address. Return 16 byte mask.
	return ipv6AllMask
}
