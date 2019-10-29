// +build !js

package p2p

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func TestBandwidthChecker(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node0 := newTestNode(t, ctx, nil)
	node1 := newTestNode(t, ctx, nil)

	// For the purposes of this test, we use only the first multiaddress for each
	// peer.
	node0AddrInfo := peer.AddrInfo{
		ID:    node0.ID(),
		Addrs: node0.Multiaddrs()[0:1],
	}
	node1AddrInfo := peer.AddrInfo{
		ID:    node1.ID(),
		Addrs: node1.Multiaddrs()[0:1],
	}

	// At first, each node should be able to connect to the other.
	require.NoError(t, node0.Connect(node1AddrInfo, testConnectionTimeout))
	require.NoError(t, node1.Connect(node0AddrInfo, testConnectionTimeout))

	// Repeatedly send messages from node0 to node1 that would exceed the
	// bandwidth limit.
	newMaxBytesPerSecond := float64(1)
	message := make([]byte, int(newMaxBytesPerSecond*100))
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Break the loop and exit goroutine when context is canceled.
				return
			case <-ticker.C:
				require.NoError(t, node0.Send(message))
			}
		}
	}()

	// Wait for node1 to receive the message.
	expectedMessage := &Message{
		From: node0.ID(),
		Data: message,
	}
	expectMessage(t, node1, expectedMessage, 15*time.Second)

	// Manually change the bandwidth limit for node1.
	node1.banner.SetMaxBytesPerSecond(newMaxBytesPerSecond)

	// Wait for node1 to block node0 and for the connection to close.
	waitForNodeToBlockAddr(t, node1, node0AddrInfo.Addrs[0], 5*time.Second)
	waitForNodesToDisconect(t, node0, node1, 5*time.Second)
}

func waitForNodeToBlockAddr(t *testing.T, blocker *Node, addressToBlock ma.Multiaddr, timeout time.Duration) {
	blockedCheckTimeout := time.After(timeout)
	blockedCheckInterval := 250 * time.Millisecond

	for {
		select {
		case <-blockedCheckTimeout:
			t.Fatal("timed out waiting for node to block the given address")
			return
		default:
		}

		blocker.banner.CheckBandwidthUsage()
		isBlocked := blocker.banner.IsAddrBanned(addressToBlock)
		if isBlocked {
			// This is what we want. Return and continue the test.
			return
		}

		// Otherwise wait a bit and then check again.
		time.Sleep(blockedCheckInterval)
		continue
	}
}

func waitForNodesToDisconect(t *testing.T, node0 *Node, node1 *Node, timeout time.Duration) {
	disconnectTimeout := time.After(timeout)
	disconnectCheckInterval := 250 * time.Millisecond

	for {
		select {
		case <-disconnectTimeout:
			t.Fatal("timed out waiting for node0 and node1 to disconnect")
			return
		default:
		}

		// Check if node0 is connected to node1
		if node0.host.Network().Connectedness(node1.ID()) != network.NotConnected {
			time.Sleep(disconnectCheckInterval)
			continue
		}

		// Check if node1 is connected to node0
		if node1.host.Network().Connectedness(node0.ID()) != network.NotConnected {
			time.Sleep(disconnectCheckInterval)
			continue
		}

		// If neither node is connected to the other, we're done.
		return
	}
}
