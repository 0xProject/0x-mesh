// +build !js

package p2p

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/p2p/banner"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	p2pnet "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTopic             = "0x-mesh-testing"
	testConnectionTimeout = 1 * time.Second
	testStreamTimeout     = 10 * time.Second
)

var (
	testRendezvousPoints = []string{"0x-mesh-testing-rendezvous"}
)

// dummyMessageHandler satisfies the MessageHandler interface but considers all
// messages valid and doesn't actually store or share any messages.
type dummyMessageHandler struct{}

func (*dummyMessageHandler) HandleMessages(ctx context.Context, messages []*Message) error {
	return nil
}

// testNotifee can be used to listen for new connections and new streams.
type testNotifee struct {
	conns   chan p2pnet.Conn
	streams chan p2pnet.Stream
}

func (n *testNotifee) Listen(p2pnet.Network, ma.Multiaddr)                   {}
func (n *testNotifee) ListenClose(p2pnet.Network, ma.Multiaddr)              {}
func (n *testNotifee) Disconnected(network p2pnet.Network, conn p2pnet.Conn) {}
func (n *testNotifee) ClosedStream(p2pnet.Network, p2pnet.Stream)            {}

// Connected is called when a connection opened
func (n *testNotifee) Connected(network p2pnet.Network, conn p2pnet.Conn) {
	if n.conns == nil {
		return
	}
	go func() {
		n.conns <- conn
	}()
}

// OpenedStream is called when a stream is opened.
func (n *testNotifee) OpenedStream(network p2pnet.Network, stream p2pnet.Stream) {
	if n.streams == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		waitForStreamProtocol(ctx, stream)
		n.streams <- stream
	}()
}

// newTestNode creates and returns a Node which is suitable for testing
// purposes.
func newTestNode(t *testing.T, ctx context.Context, notifee p2pnet.Notifiee) *Node {
	privKey, _, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
	require.NoError(t, err)
	db, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		PrivateKey:       privKey,
		MessageHandler:   &dummyMessageHandler{},
		RendezvousPoints: testRendezvousPoints,
		UseBootstrapList: false,
		DB:               db,
	}

	return newTestNodeWithConfig(t, ctx, notifee, config)
}

// newTestNodeWithConfig creates and returns a Node which is suitable for testing
// purposes and allows a custom config object.
func newTestNodeWithConfig(t *testing.T, ctx context.Context, notifee p2pnet.Notifiee, config Config) *Node {
	if config.PrivateKey == nil {
		privKey, _, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
		require.NoError(t, err)
		config.PrivateKey = privKey
	}
	node, err := New(ctx, config)
	require.NoError(t, err)

	// If notifee is not nil, set up *both* hosts to use it.
	if notifee != nil {
		node.host.Network().Notify(notifee)
	}

	return node
}

func connectTestNodes(t *testing.T, node0, node1 *Node) {
	node1PeerInfo := peer.AddrInfo{
		ID:    node1.ID(),
		Addrs: node1.Multiaddrs(),
	}
	require.NoError(t, node0.Connect(node1PeerInfo, testConnectionTimeout))
}

func startNodeAndCheckError(t *testing.T, node *Node) {
	if err := node.Start(); err != nil {
		if err == context.DeadlineExceeded || err == context.Canceled {
			// This is an expected part of the tear down process.
			return
		}
		// For other types of errors we should fail the test.
		require.NoError(t, err)
	}
}

type inMemoryMessageHandler struct {
	validator func(*Message) (bool, error)
	messages  []*Message
}

func newInMemoryMessageHandler(validator func(*Message) (bool, error)) *inMemoryMessageHandler {
	return &inMemoryMessageHandler{
		validator: validator,
	}
}

// count returns the current number of messages stored.
func (mh *inMemoryMessageHandler) count() int {
	return len(mh.messages)
}

func (mh *inMemoryMessageHandler) HandleMessages(ctx context.Context, messages []*Message) error {
	validMessages := []*Message{}
	for _, msg := range messages {
		valid, err := mh.validator(msg)
		if err != nil {
			return err
		}
		if valid {
			validMessages = append(validMessages, msg)
		}
	}
	if err := mh.store(validMessages); err != nil {
		return err
	}
	return nil
}

func (mh *inMemoryMessageHandler) store(messages []*Message) error {
	for _, msg := range messages {
		found := false
		for _, existing := range mh.messages {
			if bytes.Equal(existing.Data, msg.Data) {
				found = true
				break
			}
		}
		if found {
			// Don't need to store. Already in existing messages.
			continue
		}
		// append the new message and keep the list of messages sorted for easy
		// testing.
		mh.messages = append(mh.messages, msg)
		sort.Slice(mh.messages, func(i int, j int) bool {
			return bytes.Compare(mh.messages[i].Data, mh.messages[j].Data) == -1
		})
	}
	return nil
}

// func TestPingPong(t *testing.T) {
// 	t.Parallel()
//
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	// Create a test notifee which will be used to detect new streams.
// 	notifee := &testNotifee{
// 		streams: make(chan p2pnet.Stream),
// 	}
//
// 	// Create two test nodes and connect them.
// 	node0 := newTestNode(t, ctx, notifee)
// 	node1 := newTestNode(t, ctx, notifee)
// 	connectTestNodes(t, node0, node1)
//
// 	// Wait for a total of 2 x 2 = 4 GossipSub streams to open (2 streams per
// 	// connection; 2 connections).
// 	waitForGossipSubStreams(t, ctx, notifee, 4, testStreamTimeout)
//
// 	// HACK(albrow): Even though the stream for GossipSub has already been
// 	// opened on both sides, the ping message might *still* not be received by the
// 	// other peer. Waiting for 1 second gives each peer enough time to finish
// 	// setting up GossipSub. I couldn't find any way to avoid this hack :(
// 	time.Sleep(5 * time.Second)
//
// 	// Send ping from node0 to node1
// 	pingMessage := &Message{From: node0.host.ID(), Data: []byte("ping\n")}
// 	require.NoError(t, node0.SendV3(pingMessage.Data))
// 	const pingPongTimeout = 20 * time.Second
// 	expectMessage(t, node1, testTopic, pingMessage, pingPongTimeout)
//
// 	// Send pong from node1 to node0
// 	pongMessage := &Message{From: node1.host.ID(), Data: []byte("pong\n")}
// 	require.NoError(t, node1.SendV3(pongMessage.Data))
// 	expectMessage(t, node0, testTopic, pongMessage, pingPongTimeout)
// }

// func expectMessage(t *testing.T, node *Node, topic string, expected *Message, timeout time.Duration) {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			t.Fatal("timed out waiting for message")
// 			return
// 		default:
// 		}
// 		actual, err := receive(ctx, topic)
// 		require.NoError(t, err)
// 		// We might receive other messages. Ignore anything that doesn't match the
// 		// expected message.
// 		if assert.ObjectsAreEqualValues(expected, actual) {
// 			return
// 		}
// 	}
// }

func TestPeerDiscovery(t *testing.T) {
	t.Parallel()
	// Create a test notifee which will be used to detect new connections.
	notifee := &testNotifee{
		conns: make(chan p2pnet.Conn),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create three nodes: 0, 1, and 2.
	//
	//   - node0 is connected to node1
	//   - node1 is connected to node2
	//   - node0 is not initially connected to node2
	//
	node0 := newTestNode(t, ctx, notifee)
	node1 := newTestNode(t, ctx, notifee)
	node2 := newTestNode(t, ctx, notifee)
	go startNodeAndCheckError(t, node0)
	go startNodeAndCheckError(t, node1)
	go startNodeAndCheckError(t, node2)
	connectTestNodes(t, node0, node1)
	connectTestNodes(t, node1, node2)

	// Wait for node0 and node2 to find each other
loop:
	for {
		select {
		case <-ctx.Done():
			t.Fatal("timed out waiting for node0 to discover node2")
		case conn := <-notifee.conns:
			if conn.LocalPeer() == node0.ID() && conn.RemotePeer() == node2.ID() {
				break loop
			}
		}
	}
}

func TestBanIP(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node0 := newTestNode(t, ctx, nil)
	node1 := newTestNode(t, ctx, nil)
	go startNodeAndCheckError(t, node0)
	go startNodeAndCheckError(t, node1)

	node0AddrInfo := peer.AddrInfo{
		ID:    node0.ID(),
		Addrs: node0.Multiaddrs(),
	}
	node1AddrInfo := peer.AddrInfo{
		ID:    node1.ID(),
		Addrs: node1.Multiaddrs(),
	}

	// Ban all node1 IP addresses.
	for _, maddr := range node1.Multiaddrs() {
		require.NoError(t, node0.banner.BanIP(maddr))
	}

	// node0 should not be able to connect to node1 and vice versa.
	// Unfortunately, libp2p swallows the error and creates a new one so there is
	// no way for us to guarantee that the error we got is the one that we expect.
	err := node0.Connect(node1AddrInfo, testConnectionTimeout)
	require.Error(t, err, "node0 should not be able to connect to node1")
	err = node1.Connect(node0AddrInfo, testConnectionTimeout)
	require.Error(t, err, "node1 should not be able to connect to node0")
}

func TestUnbanIP(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node0 := newTestNode(t, ctx, nil)
	node1 := newTestNode(t, ctx, nil)
	go startNodeAndCheckError(t, node0)
	go startNodeAndCheckError(t, node1)

	node0AddrInfo := peer.AddrInfo{
		ID:    node0.ID(),
		Addrs: node0.Multiaddrs(),
	}
	node1AddrInfo := peer.AddrInfo{
		ID:    node1.ID(),
		Addrs: node1.Multiaddrs(),
	}

	// Ban all node1 IP addresses.
	for _, maddr := range node1.Multiaddrs() {
		require.NoError(t, node0.banner.BanIP(maddr))
	}

	// Unban all node1 IP addresses.
	for _, maddr := range node1.Multiaddrs() {
		require.NoError(t, node0.banner.UnbanIP(maddr))
	}

	// Each node should now be able to connect to the other.
	require.NoError(t, node0.Connect(node1AddrInfo, testConnectionTimeout))
	require.NoError(t, node1.Connect(node0AddrInfo, testConnectionTimeout))
}

func TestProtectIP(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node0 := newTestNode(t, ctx, nil)
	node1 := newTestNode(t, ctx, nil)
	go startNodeAndCheckError(t, node0)
	go startNodeAndCheckError(t, node1)

	node0AddrInfo := peer.AddrInfo{
		ID:    node0.ID(),
		Addrs: node0.Multiaddrs(),
	}
	node1AddrInfo := peer.AddrInfo{
		ID:    node1.ID(),
		Addrs: node1.Multiaddrs(),
	}

	// Protect all node1 IP addresses.
	for _, maddr := range node1.Multiaddrs() {
		require.NoError(t, node0.banner.ProtectIP(maddr))
	}

	// Ban all node1 IP addresses (this should have no effect since the IP
	// addresses are protected).
	for _, maddr := range node1.Multiaddrs() {
		require.EqualError(t, node0.banner.BanIP(maddr), banner.ErrProtectedIP.Error())
	}

	// Each node should now be able to connect to the other.
	require.NoError(t, node0.Connect(node1AddrInfo, testConnectionTimeout))
	require.NoError(t, node1.Connect(node0AddrInfo, testConnectionTimeout))
}

func TestRateValidatorGlobal(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a test notifee which will be used to detect new streams.
	notifee := &testNotifee{
		streams: make(chan p2pnet.Stream),
	}

	db0, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	node0Config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		MessageHandler: newInMemoryMessageHandler(func(*Message) (bool, error) {
			return true, nil
		}),
		RendezvousPoints:         testRendezvousPoints,
		UseBootstrapList:         false,
		DB:                       db0,
		GlobalPubSubMessageLimit: 1,
		GlobalPubSubMessageBurst: 5,
	}
	db1, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	node1Config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		MessageHandler: newInMemoryMessageHandler(func(*Message) (bool, error) {
			return true, nil
		}),
		RendezvousPoints:         testRendezvousPoints,
		UseBootstrapList:         false,
		DB:                       db1,
		GlobalPubSubMessageLimit: 1,
		GlobalPubSubMessageBurst: 5,
	}
	db2, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	node2Config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		MessageHandler: newInMemoryMessageHandler(func(*Message) (bool, error) {
			return true, nil
		}),
		RendezvousPoints:         testRendezvousPoints,
		UseBootstrapList:         false,
		DB:                       db2,
		GlobalPubSubMessageLimit: 1,
		GlobalPubSubMessageBurst: 5,
	}

	// Create three test nodes. node0 is connected to node1. node1 is connected to
	// node2.
	node0 := newTestNodeWithConfig(t, ctx, notifee, node0Config)
	node1 := newTestNodeWithConfig(t, ctx, notifee, node1Config)
	node2 := newTestNodeWithConfig(t, ctx, notifee, node2Config)
	connectTestNodes(t, node0, node1)
	connectTestNodes(t, node1, node2)

	// Wait for a total of 2 x 3 = 6 GossipSub streams to open (2 streams per
	// connection; 3 connections).
	waitForGossipSubStreams(t, ctx, notifee, 6, testStreamTimeout)

	// HACK(albrow): Wait for GossipSub to finish initializing.
	time.Sleep(2 * time.Second)

	require.NoError(t, node1.receiveAndHandleMessages(ctx))
	require.NoError(t, node2.receiveAndHandleMessages(ctx))

	// node0 sends config.GlobalPubSubMessageBurst*2 messages to node1.
	for i := 0; i < node0.config.GlobalPubSubMessageBurst*2; i++ {
		msg := []byte(fmt.Sprintf("message_%d", i))
		require.NoError(t, node0.SendV3(msg))
	}

	// HACK(albrow): Wait for GossipSub messages to fully propagate.
	time.Sleep(1 * time.Second)

	require.NoError(t, node1.receiveAndHandleMessages(ctx))
	require.NoError(t, node2.receiveAndHandleMessages(ctx))

	// node1 and node2 should only have config.GlobalPubSubMessageBurst messages.
	// The others are expected to have been dropped.
	expectedMessageCount := node0.config.GlobalPubSubMessageBurst
	node1MessageCount := node1.messageHandler.(*inMemoryMessageHandler).count()
	assert.Equal(t, expectedMessageCount, node1MessageCount, "node1 received and stored the wrong number of messages")
	node2MessageCount := node2.messageHandler.(*inMemoryMessageHandler).count()
	assert.Equal(t, expectedMessageCount, node2MessageCount, "node2 received and stored the wrong number of messages")
}

func TestRateValidatorPerPeer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a test notifee which will be used to detect new streams.
	notifee := &testNotifee{
		streams: make(chan p2pnet.Stream),
	}

	db0, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	node0Config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		MessageHandler: newInMemoryMessageHandler(func(*Message) (bool, error) {
			return true, nil
		}),
		RendezvousPoints:          testRendezvousPoints,
		UseBootstrapList:          false,
		DB:                        db0,
		PerPeerPubSubMessageLimit: 1,
		PerPeerPubSubMessageBurst: 5,
	}
	db1, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	node1Config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		MessageHandler: newInMemoryMessageHandler(func(*Message) (bool, error) {
			return true, nil
		}),
		RendezvousPoints:          testRendezvousPoints,
		UseBootstrapList:          false,
		DB:                        db1,
		PerPeerPubSubMessageLimit: 1,
		PerPeerPubSubMessageBurst: 5,
	}
	db2, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	node2Config := Config{
		SubscribeTopicV3: testTopic,
		PublishTopicsV3:  []string{testTopic},
		MessageHandler: newInMemoryMessageHandler(func(*Message) (bool, error) {
			return true, nil
		}),
		RendezvousPoints:          testRendezvousPoints,
		UseBootstrapList:          false,
		DB:                        db2,
		PerPeerPubSubMessageLimit: 1,
		PerPeerPubSubMessageBurst: 5,
	}

	// Create three test nodes. Both node0 and node1 are connected to node2.
	node0 := newTestNodeWithConfig(t, ctx, notifee, node0Config)
	node1 := newTestNodeWithConfig(t, ctx, notifee, node1Config)
	node2 := newTestNodeWithConfig(t, ctx, notifee, node2Config)
	connectTestNodes(t, node0, node2)
	connectTestNodes(t, node1, node2)

	// Wait for a total of 2 x 3 = 6 GossipSub streams to open (2 streams per
	// connection; 3 connections).
	waitForGossipSubStreams(t, ctx, notifee, 6, testStreamTimeout)

	// HACK(albrow): Wait for GossipSub to finish initializing.
	time.Sleep(2 * time.Second)

	require.NoError(t, node0.receiveAndHandleMessages(ctx))
	require.NoError(t, node1.receiveAndHandleMessages(ctx))
	require.NoError(t, node2.receiveAndHandleMessages(ctx))

	// node0 sends config.PeerPeerPubSubMessageBurst*2 messages to node2.
	for i := 0; i < node0.config.PerPeerPubSubMessageBurst*2; i++ {
		msg := []byte(fmt.Sprintf("node0_message_%d", i))
		require.NoError(t, node0.SendV3(msg))
	}
	// node1 sends config.PeerPeerPubSubMessageBurst*2 messages to node2.
	for i := 0; i < node1.config.PerPeerPubSubMessageBurst*2; i++ {
		msg := []byte(fmt.Sprintf("node1_message_%d", i))
		require.NoError(t, node1.SendV3(msg))
	}

	// HACK(albrow): Wait for GossipSub messages to fully propagate.
	time.Sleep(1 * time.Second)

	require.NoError(t, node0.receiveAndHandleMessages(ctx))
	require.NoError(t, node1.receiveAndHandleMessages(ctx))
	require.NoError(t, node2.receiveAndHandleMessages(ctx))

	// node2 should only have config.PerPeerPubSubMessageBurst*2 messages.
	// The others are expected to have been dropped.
	expectedMessageCount := node0.config.PerPeerPubSubMessageBurst * 2
	node2MessageCount := node2.messageHandler.(*inMemoryMessageHandler).count()
	assert.Equal(t, expectedMessageCount, node2MessageCount, "node2 received and stored the wrong number of messages")
}

func waitForGossipSubStreams(t *testing.T, ctx context.Context, notifee *testNotifee, count int, timeout time.Duration) {
	streamCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	streamCount := 0
loop:
	for {
		select {
		case <-streamCtx.Done():
			t.Fatal("timed out waiting for pubsub stream to open")
		case stream := <-notifee.streams:
			if stream.Protocol() == pubsubProtocolID {
				streamCount += 1
				if streamCount == count {
					break loop
				}
			}
		}
	}
}
