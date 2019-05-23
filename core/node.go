// +build !js

package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	crypto "github.com/libp2p/go-libp2p-crypto"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	p2pnet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	// receiveTimeout is the maximum amount of time to wait for receiving new messages.
	receiveTimeout = 1 * time.Second
	// maxReceiveBatch is the maximum number of new messages to receive at once.
	maxReceiveBatch = 100
	// maxShareBatch is the maximum number of messages to share at once.
	maxShareBatch = 10
	// peerCountLow is the target number of peers to connect to at any given time.
	peerCountLow = 10
	// peerCountHigh is the maximum number of peers to be connected to. If the
	// number of connections exceeds this number, we will prune connections until
	// we reach peerCountLow.
	peerCountHigh = 12
	// peerGraceDuration is the amount of time a newly opened connection is given
	// before it becomes subject to pruning.
	peerGraceDuration = 10 * time.Second
	// defaultNetworkTimeout is the default timeout for network requests (e.g.
	// connecting to a new peer).
	defaultNetworkTimeout = 30 * time.Second
	// advertiseTTL is the TTL for our announcement to the discovery network.
	advertiseTTL = 5 * time.Minute
)

// Node is the main type for the core package. It represents a particpant in the
// 0x Mesh network who is capable of sending, receiving, validating, and storing
// messages.
type Node struct {
	ctx              context.Context
	cancel           context.CancelFunc
	host             host.Host
	dht              *dht.IpfsDHT
	routingDiscovery discovery.Discovery
	connManager      *connmgr.BasicConnMgr
	pubsub           *pubsub.PubSub
	sub              *pubsub.Subscription
	config           Config
	messageHandler   MessageHandler
	notifee          p2pnet.Notifiee
}

// Config contains configuration options for a Node.
type Config struct {
	// Topic is a unique string representing the pubsub topic. Only Nodes which
	// have the same topic will share messages with one another.
	Topic string
	// ListenPort is the port on which to listen for new connections. It can be
	// set to 0 to make the OS automatically choose any available port.
	ListenPort int
	// Insecure controls whether or not messages should be encrypted. It should
	// always be set to false in production.
	Insecure bool
	// RandSeed is a random seed used to generate a private key. If the same seed
	// is used, the same private key will be generated. It can be set to 0 to use
	// no seed and generate a cryptographically secure private key.
	RandSeed int64
	// MessageHandler is an interface responsible for validating, storing, and
	// finding new messages to share.
	MessageHandler MessageHandler
	// RendezvousString is a unique identifier for the rendezvous point. This node
	// will attempt to find peers with the same Rendezvous string.
	RendezvousString string
	// UseBootstrapList determines whether or not to use the list of predetermined
	// peers to bootstrap the DHT for peer discovery.
	UseBootstrapList bool
}

// New creates a new Node with the given config.
func New(config Config) (*Node, error) {
	nodeCtx, cancel := context.WithCancel(context.Background())

	if config.MessageHandler == nil {
		cancel()
		return nil, errors.New("config.MessageHandler is required")
	} else if config.RendezvousString == "" {
		cancel()
		return nil, errors.New("config.RendezvousString is required")
	}

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if config.RandSeed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(config.RandSeed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		cancel()
		return nil, err
	}

	// Set up the transport and the host.
	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", config.ListenPort))
	if err != nil {
		cancel()
		return nil, err
	}
	connManager := connmgr.NewConnManager(peerCountLow, peerCountHigh, peerGraceDuration)
	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
		libp2p.ConnectionManager(connManager),
	}
	if config.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	}
	basicHost, err := libp2p.New(nodeCtx, opts...)
	if err != nil {
		cancel()
		return nil, err
	}

	// Set up DHT for peer discovery.
	kadDHT, err := dht.New(nodeCtx, basicHost)
	if err != nil {
		cancel()
		return nil, err
	}
	routingDiscovery := discovery.NewRoutingDiscovery(kadDHT)

	// Set up pubsub.
	// TODO: Replace with WeijieSub. Using FloodSub for now.
	ps, err := pubsub.NewFloodSub(nodeCtx, basicHost)
	if err != nil {
		cancel()
		return nil, err
	}
	sub, err := ps.Subscribe(config.Topic)
	if err != nil {
		cancel()
		return nil, err
	}

	// Create the Node.
	node := &Node{
		ctx:              nodeCtx,
		cancel:           cancel,
		host:             basicHost,
		dht:              kadDHT,
		routingDiscovery: routingDiscovery,
		connManager:      connManager,
		config:           config,
		pubsub:           ps,
		sub:              sub,
		messageHandler:   config.MessageHandler,
	}

	// Set up the notifee.
	basicHost.Network().Notify(&notifee{node: node})

	return node, nil
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

	// Advertise ourselves for the purposes of peer discovery.
	discovery.Advertise(n.ctx, n.routingDiscovery, n.config.RendezvousString, discovery.TTL(advertiseTTL))

	if n.config.UseBootstrapList {
		if err := n.connectToBootstrapList(); err != nil {
			return err
		}
	}

	return n.mainLoop()
}

func (n *Node) connectToBootstrapList() error {
	log.WithField("bootstrapPeers", dht.DefaultBootstrapPeers).Info("connecting to bootstrap peers")
	connectCtx, cancel := context.WithTimeout(n.ctx, defaultNetworkTimeout)
	defer cancel()
	wg := sync.WaitGroup{}
	// TODO(albrow): Replace default list with our own bootstrap peer list.
	for _, addr := range dht.DefaultBootstrapPeers {
		peerInfo, err := peerstore.InfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := n.host.Connect(connectCtx, *peerInfo); err != nil {
				log.WithFields(map[string]interface{}{
					"error":    err.Error(),
					"peerInfo": peerInfo,
				}).Warn("failed to connect to bootstrap peer")
			}
		}()
	}
	wg.Wait()
	return nil
}

// Connect ensures there is a connection between this host and the peer with
// given peerInfo. If there is not an active connection, Connect will dial the
// peer, and block until a connection is open, or an error is returned.
func (n *Node) Connect(ctx context.Context, peerInfo peerstore.PeerInfo) error {
	err := n.host.Connect(ctx, peerInfo)
	if err != nil {
		return err
	}
	return nil
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
	_, err = n.messageHandler.ValidateAndStore(incoming)
	if err != nil {
		return fmt.Errorf("could not validate or store messages: %s", err.Error())
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
	}).Debug("looking for new peers")
	findPeersCtx, cancel := context.WithTimeout(n.ctx, defaultNetworkTimeout)
	defer cancel()
	peerChan, err := n.routingDiscovery.FindPeers(findPeersCtx, n.config.RendezvousString, discovery.Limit(max))
	if err != nil {
		return err
	}

	connectCtx, cancel := context.WithTimeout(n.ctx, defaultNetworkTimeout)
	defer cancel()
	for peer := range peerChan {
		if peer.ID == n.host.ID() {
			continue
		}
		log.WithFields(map[string]interface{}{
			"peerInfo": peer,
		}).Debug("found peer via rendezvous")
		if err := n.host.Connect(connectCtx, peer); err != nil {
			// If we fail to connect to a single peer we should still keep trying the
			// others. Log instead of returning the error.
			log.WithFields(map[string]interface{}{
				"error":    err.Error(),
				"peerInfo": peer,
			}).Warn("could not connect to peer")
		}
	}
	return nil
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
	msg, err := n.sub.Next(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{From: msg.GetFrom(), Data: msg.Data}, nil
}

// Close closes the Node and any active connections.
func (n *Node) Close() error {
	n.cancel()
	return n.host.Close()
}
