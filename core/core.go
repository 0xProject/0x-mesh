package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	direct "github.com/libp2p/go-libp2p-webrtc-direct"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pions/webrtc"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
)

type Node struct {
	host     host.Host
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

	transport := direct.NewTransport(
		webrtc.Configuration{},
		new(mplex.Transport),
	)

	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/http/p2p-webrtc-direct", config.ListenPort))
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
		libp2p.Transport(transport),
	}

	if config.Insecure {
		opts = append(opts, libp2p.NoSecurity)
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return &Node{
		host:   basicHost,
		config: config,
	}, nil
}

// Send sends a message to connected peers in the network.
func (n *Node) Send(msg *Message) error {
	return errors.New("Not yet implemented")
}

// Receive returns a channel that can be used to listen for new messages.
func (n *Node) Receive() <-chan *Message {
	return n.messages
}

// Evict signals that a message should be evicted. The Node will update its
// score for each neighbor appropriately.
func (n *Node) Evict(msg *Message) error {
	return errors.New("Not yet implemented")
}

// Close closes the Node and any active connections.
func (n *Node) Close() error {
	return errors.New("Not yet implemented")
}
