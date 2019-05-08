// +build !js

package core

import (
	"bytes"
	"context"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Counter used for config.RandSeed. Atomically incremented each time a new Node
// is created.
var counter int64 = -1

// dummyMessageHandler satisfies the MessageHandler interface but considers all
// messages valid and doesn't actually store or share any messages.
type dummyMessageHandler struct{}

func (*dummyMessageHandler) Validate(*Message) (bool, error) {
	return true, nil
}

func (*dummyMessageHandler) Store(*Message) error {
	return nil
}

func (*dummyMessageHandler) GetMessagesToShare(max int) ([][]byte, error) {
	return nil, nil
}

// newTestNode creates and returns a Node which is suitable for testing
// purposes.
func newTestNode(t *testing.T) *Node {
	config := Config{
		Topic:          "0x-mesh-testing",
		ListenPort:     0, // Let OS randomly choose an open port.
		Insecure:       true,
		RandSeed:       atomic.AddInt64(&counter, 1),
		MessageHandler: &dummyMessageHandler{},
	}
	node, err := New(config)
	require.NoError(t, err)
	return node
}

func TestPingPong(t *testing.T) {
	// Create two nodes and add each one to the other's peerstore and connect
	// them. Note that the connection is symmetrical so we only need to establish
	// one connection.
	node0 := newTestNode(t)
	node1 := newTestNode(t)
	node0.host.Peerstore().AddAddrs(node1.host.ID(), node1.host.Addrs(), peerstore.PermanentAddrTTL)
	node1.host.Peerstore().AddAddrs(node0.host.ID(), node0.host.Addrs(), peerstore.PermanentAddrTTL)
	node1PeerInfo := node0.host.Peerstore().PeerInfo(node1.host.ID())
	connectContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, node0.host.Connect(connectContext, node1PeerInfo))
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

func (mh *inMemoryMessageHandler) Validate(msg *Message) (bool, error) {
	return mh.validator(msg)
}

func (mh *inMemoryMessageHandler) Store(msg *Message) error {
	for _, existing := range mh.messages {
		if bytes.Compare(existing.Data, msg.Data) == 0 {
			// Don't need to store. Already in existing messages.
			return nil
		}
	}
	// append the new message and keep the list of messages sorted for easy
	// testing.
	mh.messages = append(mh.messages, msg)
	sort.Slice(mh.messages, func(i int, j int) bool {
		return bytes.Compare(mh.messages[i].Data, mh.messages[j].Data) == -1
	})
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
	// Create two nodes and add each one to the other's peerstore and connect
	// them. Note that the connection is symmetrical so we only need to establish
	// one connection.
	node0 := newTestNode(t)
	node1 := newTestNode(t)
	node0.host.Peerstore().AddAddrs(node1.host.ID(), node1.host.Addrs(), peerstore.PermanentAddrTTL)
	node1.host.Peerstore().AddAddrs(node0.host.ID(), node0.host.Addrs(), peerstore.PermanentAddrTTL)
	node1PeerInfo := node0.host.Peerstore().PeerInfo(node1.host.ID())
	connectContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, node0.host.Connect(connectContext, node1PeerInfo))
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
