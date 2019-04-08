// +build !js

package core

import (
	"bufio"
	"context"
	"io/ioutil"
	"sync/atomic"
	"testing"

	net "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
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
	// Create two nodes and add each one to the other's peerstore. This allows
	// them to connect to one another.
	node0 := newTestNode(t)
	node1 := newTestNode(t)
	node0.host.Peerstore().AddAddrs(node1.host.ID(), node1.host.Addrs(), peerstore.PermanentAddrTTL)
	node1.host.Peerstore().AddAddrs(node0.host.ID(), node0.host.Addrs(), peerstore.PermanentAddrTTL)

	ping := []byte("ping\n")
	pong := []byte("pong\n")
	protocol := protocol.ID("/test/0.0.1")

	// Set a stream handler on node0. This will be called when the stream is open.
	node0.host.SetStreamHandler(protocol, func(stream net.Stream) {
		// Send the "ping" message.
		_, err := stream.Write(ping)
		require.NoError(t, err)
		defer stream.Close()
		// Expect to receive "pong".
		res, err := ioutil.ReadAll(stream)
		require.NoError(t, err)
		assert.Equal(t, res, pong, "node0 did not receive pong from node1")
	})

	// Create the stream.
	stream, err := node1.host.NewStream(context.Background(), node0.host.ID(), protocol)
	require.NoError(t, err)
	// Expect to receive the "ping" message.
	buf := bufio.NewReader(stream)
	req, err := buf.ReadBytes('\n')
	assert.Equal(t, ping, req, "node1 did not receive ping from node0")
	// Send the "pong" message.
	_, err = stream.Write([]byte("pong\n"))
	stream.Close()
	require.NoError(t, err)
}
