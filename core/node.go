// +build !js

package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
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
)

// Node is the main type for the core package. It represents a particpant in the
// 0x Mesh network who is capable of sending, receiving, validating, and storing
// messages.
type Node struct {
	ctx            context.Context
	cancel         context.CancelFunc
	host           host.Host
	pubsub         *pubsub.PubSub
	sub            *pubsub.Subscription
	config         Config
	messageHandler MessageHandler
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
}

// New creates a new Node with the given config.
func New(config Config) (*Node, error) {
	if config.MessageHandler == nil {
		return nil, errors.New("config.MessageHandler is required")
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
		return nil, err
	}

	// Set up the transport and the host.
	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", config.ListenPort))
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}
	if config.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	}
	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Set up pubsub.
	// TODO: Replace with WeijieSub. Using FloodSub for now.
	ps, err := pubsub.NewFloodSub(context.Background(), basicHost)
	if err != nil {
		return nil, err
	}
	sub, err := ps.Subscribe(config.Topic)
	if err != nil {
		return nil, err
	}

	// Create the Node.
	ctx, cancel := context.WithCancel(context.Background())
	node := &Node{
		ctx:            ctx,
		cancel:         cancel,
		host:           basicHost,
		config:         config,
		pubsub:         ps,
		sub:            sub,
		messageHandler: config.MessageHandler,
	}

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
	return n.mainLoop()
}

// Connect ensures there is a connection between this host and the peer with
// given peerInfo. If there is not an active connection, Connect will dial the
// peer, and block until a connection is open, or an error is returned.
func (n *Node) Connect(ctx context.Context, peerInfo peerstore.PeerInfo) error {
	err := n.host.Connect(ctx, peerInfo)
	if err != nil {
		return err
	}
	log.WithField("peerInfo", peerInfo).Debug("connected to peer")
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
	n.sub.Cancel()
	return n.host.Close()
}
