// +build !js

package core

import (
	"context"
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

// newTestNode creates and returns a Node which is suitable for testing
// purposes.
func newTestNode(t *testing.T) *Node {
	config := Config{
		ListenPort: 0, // Let OS randomly choose an open port.
		Insecure:   true,
		RandSeed:   atomic.AddInt64(&counter, 1),
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
	connectContext, _ := context.WithTimeout(context.Background(), 10*time.Second)
	node0.host.Connect(connectContext, node1PeerInfo)
	defer node0.Close()
	defer node1.Close()

	// Send ping from node0 to node1
	pingMessage := &Message{Data: []byte("ping\n")}
	require.NoError(t, node0.Send(pingMessage))
	const pingPongTimeout = 15 * time.Second
	expectMessage(t, node1.Receive(), pingMessage, pingPongTimeout)

	// Send pong from node1 to node0
	pongMessage := &Message{Data: []byte("pong\n")}
	require.NoError(t, node1.Send(pongMessage))
	expectMessage(t, node0.Receive(), pongMessage, pingPongTimeout)
}

func expectMessage(t *testing.T, ch <-chan *Message, expected *Message, timeout time.Duration) {
	timeoutChan := time.After(timeout)
	for {
		select {
		case msg := <-ch:
			// We might receive other messages. Ignore anything that doesn't match the
			// expected message.
			if assert.ObjectsAreEqualValues(expected, msg) {
				return
			}
		case <-timeoutChan:
			t.Errorf("Timed out after %s waiting for message: %v\n", timeout, expected)
			return
		}
	}
}
