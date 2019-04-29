// +build !js

package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
)

const pubsubTopic = "0x-orders"

// messageBuffer is the number of messages to hold in memory at once. It is the
// buffer length of the messages channel.
const messageBuffer = 100

type Node struct {
	host     host.Host
	pubsub   *pubsub.PubSub
	sub      *pubsub.Subscription
	config   Config
	messages chan *Message
}

type Config struct {
	ListenPort int
	Insecure   bool
	RandSeed   int64
}

// New creates a new Node with the given config. The Node will automatically
// and continuously connect to peers and receive new messages until Close is
// called.
func New(config Config) (*Node, error) {
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
	sub, err := ps.Subscribe(pubsubTopic)
	if err != nil {
		return nil, err
	}

	// Create the Node.
	node := &Node{
		host:     basicHost,
		config:   config,
		pubsub:   ps,
		sub:      sub,
		messages: make(chan *Message, messageBuffer),
	}

	// Start listening in the background for messages over the subscription.
	go node.listenForMessages()

	return node, nil
}

func isErrSubCancelled(err error) bool {
	return err != nil && strings.Contains(err.Error(), "subscription cancelled")
}

// listenForMessages continuously listens in the background for new messages on
// n.sub. It is a blocking function, but can be called in a goroutine. It
// returns if n.sub is cancelled.
func (n *Node) listenForMessages() {
	for {
		msg, err := n.sub.Next(context.Background())
		if err != nil {
			if isErrSubCancelled(err) {
				// The sub was cancelled. We can assume Node.Close() was called or
				// there was another error that has already been returned/handled.
				return
			}
			// TODO(albrow): Don't panic here.
			panic(err)
		}
		n.messages <- &Message{Data: msg.Data}
	}
}

// Send sends a message to connected peers in the network.
func (n *Node) Send(msg *Message) error {
	// TODO(albrow): Encode the message properly instead of just using msg.Data.
	// Need to agree on message and order format first.
	return n.pubsub.Publish(pubsubTopic, msg.Data)
}

// Receive returns a read-only channel that can be used to listen for new
// messages.
func (n *Node) Receive() <-chan *Message {
	return n.messages
}

// Evict signals that a message should be evicted.
func (n *Node) Evict(msg *Message) error {
	return errors.New("Not yet implemented")
}

// Close closes the Node and any active connections.
func (n *Node) Close() error {
	n.sub.Cancel()
	// TODO(albrow): We should be closing the host here. Unfortunately there is a
	// bug where the transport panics with "panic: close of closed channel". Needs
	// to be fixed upstream.
	// return n.host.Close()
	return nil
}
