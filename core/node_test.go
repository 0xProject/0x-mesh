// +build !js

package core

import (
	"bytes"
	"context"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	p2pnet "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Counter used for config.RandSeed. Atomically incremented each time a new Node
// is created.
var counter int64 = -1

const (
	testTopic            = "0x-mesh-testing"
	testRendezvousString = "0x-mesh-testing-rendezvous"
)

// dummyMessageHandler satisfies the MessageHandler interface but considers all
// messages valid and doesn't actually store or share any messages.
type dummyMessageHandler struct{}

func (*dummyMessageHandler) ValidateAndStore(messages []*Message) ([]*Message, error) {
	return messages, nil
}

func (*dummyMessageHandler) GetMessagesToShare(max int) ([][]byte, error) {
	return nil, nil
}

// newTestNode creates and returns a Node which is suitable for testing
// purposes.
func newTestNode(t *testing.T) *Node {
	config := Config{
		Topic:            testTopic,
		ListenPort:       0, // Let OS randomly choose an open port.
		RandSeed:         atomic.AddInt64(&counter, 1),
		MessageHandler:   &dummyMessageHandler{},
		RendezvousString: testRendezvousString,
	}
	node, err := New(config)
	require.NoError(t, err)
	return node
}

func createTwoConnectedTestNodes(t *testing.T) (*Node, *Node) {
	// Create two nodes and add each one to the other's peerstore and connect
	// them.
	node0 := newTestNode(t)
	node1 := newTestNode(t)
	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	node1PeerInfo := peerstore.PeerInfo{
		ID:    node1.ID(),
		Addrs: node1.Multiaddrs(),
	}
	require.NoError(t, node0.Connect(connectCtx, node1PeerInfo))
	return node0, node1
}

func TestPingPong(t *testing.T) {
	node0, node1 := createTwoConnectedTestNodes(t)
	defer node0.Close()
	defer node1.Close()

	// Send ping from node0 to node1
	pingMessage := &Message{From: node0.host.ID(), Data: []byte("ping\n")}
	require.NoError(t, node0.send(pingMessage.Data))
	const pingPongTimeout = 5 * time.Second
	expectMessage(t, node1, pingMessage, pingPongTimeout)

	// Send pong from node1 to node0
	pongMessage := &Message{From: node1.host.ID(), Data: []byte("pong\n")}
	require.NoError(t, node1.send(pongMessage.Data))
	expectMessage(t, node0, pongMessage, pingPongTimeout)
}

func expectMessage(t *testing.T, node *Node, expected *Message, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Fatal("timed out waiting for message")
			return
		default:
		}
		actual, err := node.receive(ctx)
		require.NoError(t, err)
		// We might receive other messages. Ignore anything that doesn't match the
		// expected message.
		if assert.ObjectsAreEqualValues(expected, actual) {
			return
		}
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

func (mh *inMemoryMessageHandler) ValidateAndStore(messages []*Message) ([]*Message, error) {
	validMessages := []*Message{}
	for _, msg := range messages {
		valid, err := mh.validator(msg)
		if err != nil {
			return nil, err
		}
		if valid {
			validMessages = append(validMessages, msg)
		}
	}
	if err := mh.store(validMessages); err != nil {
		return nil, err
	}
	return validMessages, nil
}

func (mh *inMemoryMessageHandler) store(messages []*Message) error {
	for _, msg := range messages {
		found := false
		for _, existing := range mh.messages {
			if bytes.Compare(existing.Data, msg.Data) == 0 {
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

func (mh *inMemoryMessageHandler) GetMessagesToShare(max int) ([][]byte, error) {
	// Always just return the first messages up to max.
	var toShare []*Message
	if max > len(mh.messages) {
		toShare = mh.messages
	} else {
		toShare = mh.messages[:max]
	}
	data := make([][]byte, len(toShare))
	for i, msg := range toShare {
		data[i] = msg.Data
	}
	return data, nil
}

func TestMessagesAreShared(t *testing.T) {
	node0, node1 := createTwoConnectedTestNodes(t)
	defer node0.Close()
	defer node1.Close()

	// Set up special MessageHandlers for testing purposes.
	// oddMessageHandler only considers messages valid if the first byte is odd.
	oddMessageHandler := newInMemoryMessageHandler(func(msg *Message) (bool, error) {
		if len(msg.Data) < 1 {
			return false, nil
		}
		isFirstByteOdd := msg.Data[0]%2 != 0
		return isFirstByteOdd, nil
	})
	oddMessageHandler.messages = []*Message{
		{
			From: node0.host.ID(),
			Data: []byte{1, 2, 3, 4},
		},
		{
			From: node0.host.ID(),
			Data: []byte{3, 4, 5, 6},
		},
	}
	node0.messageHandler = oddMessageHandler

	// allMessageHandler considers all messages valid.
	allMessageHandler := newInMemoryMessageHandler(func(msg *Message) (bool, error) {
		return true, nil
	})
	allMessageHandler.messages = []*Message{
		{
			From: node1.host.ID(),
			Data: []byte{0, 1, 2, 3},
		},
		{
			From: node1.host.ID(),
			Data: []byte{1, 2, 3, 4},
		},
		{
			From: node1.host.ID(),
			Data: []byte{2, 3, 4, 5},
		},
		{
			From: node1.host.ID(),
			Data: []byte{5, 6, 7, 8},
		},
	}
	node1.messageHandler = allMessageHandler

	// Call runOnce to cause each node to share and receive messages.
	require.NoError(t, node0.runOnce())
	require.NoError(t, node1.runOnce())
	require.NoError(t, node0.runOnce())

	// We expect that all the odd messages have been collected by node0.
	expectedOddMessages := []*Message{
		{
			From: node0.host.ID(),
			Data: []byte{1, 2, 3, 4},
		},
		{
			From: node0.host.ID(),
			Data: []byte{3, 4, 5, 6},
		},
		{
			From: node1.host.ID(),
			Data: []byte{5, 6, 7, 8},
		},
	}
	assert.Equal(t, expectedOddMessages, node0.messageHandler.(*inMemoryMessageHandler).messages, "node0 should be storing all odd messages")

	// We expect that all messages have been collected by node1.
	expectedAllMessages := []*Message{
		{
			From: node1.host.ID(),
			Data: []byte{0, 1, 2, 3},
		},
		{
			From: node1.host.ID(),
			Data: []byte{1, 2, 3, 4},
		},
		{
			From: node1.host.ID(),
			Data: []byte{2, 3, 4, 5},
		},
		{
			From: node0.host.ID(),
			Data: []byte{3, 4, 5, 6},
		},
		{
			From: node1.host.ID(),
			Data: []byte{5, 6, 7, 8},
		},
	}
	assert.Equal(t, expectedAllMessages, node1.messageHandler.(*inMemoryMessageHandler).messages, "node1 should be storing all messages")
}

type testNotifee struct {
	conns chan p2pnet.Conn
}

func (n *testNotifee) Listen(p2pnet.Network, ma.Multiaddr)                   {}
func (n *testNotifee) ListenClose(p2pnet.Network, ma.Multiaddr)              {}
func (n *testNotifee) Disconnected(network p2pnet.Network, conn p2pnet.Conn) {}
func (n *testNotifee) OpenedStream(p2pnet.Network, p2pnet.Stream)            {}
func (n *testNotifee) ClosedStream(p2pnet.Network, p2pnet.Stream)            {}

// Connected is called when a connection opened
func (n *testNotifee) Connected(network p2pnet.Network, conn p2pnet.Conn) {
	go func() {
		n.conns <- conn
	}()
}

func TestPeerDiscovery(t *testing.T) {
	// Create three nodes: 0, 1, and 2.
	//
	//   - node0 is connected to node1
	//   - node1 is ocnnected to node2
	//   - node0 is not initially connected to node2
	//
	node0, node1 := createTwoConnectedTestNodes(t)
	defer node0.Close()
	defer node1.Close()
	node2 := newTestNode(t)
	defer node2.Close()
	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := node2.Connect(connectCtx, peerstore.PeerInfo{
		ID:    node1.host.ID(),
		Addrs: node1.host.Addrs(),
	})
	require.NoError(t, err)

	// Create a test notifee for both node0 and node1 which will be used to detect
	// new connections.
	node0Notif := &testNotifee{
		conns: make(chan p2pnet.Conn),
	}
	node0.host.Network().Notify(node0Notif)
	node2Notif := &testNotifee{
		conns: make(chan p2pnet.Conn),
	}
	node2.host.Network().Notify(node2Notif)

	// Start all the nodes (this also starts the peer discovery process).
	go func() {
		require.NoError(t, node0.Start())
	}()
	go func() {
		require.NoError(t, node1.Start())
	}()
	go func() {
		require.NoError(t, node2.Start())
	}()

	// Wait for node0 ande node2 to find each other
	timeout := time.After(10 * time.Second)
loop:
	for {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for node0 to discover node2")
		case conn := <-node0Notif.conns:
			if conn.RemotePeer() == node2.ID() {
				break loop
			}
		case conn := <-node2Notif.conns:
			if conn.RemotePeer() == node0.ID() {
				break loop
			}
		}
	}
}
