// +build !js

package p2p

import (
	"context"
	"time"

	connmgr "github.com/libp2p/go-libp2p-connmgr"
	p2pnet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	// pubsubProtocolTag is the tag used for peers who speak our pubsub protocol.
	pubsubProtocolTag = "pubsub-protocol"
	// pubsubProtocolScore is the score to add to peers who speak our pubsub
	// protocol.
	pubsubProtocolScore = 10
)

// notifee receives notifications for network-related events.
type notifee struct {
	ctx         context.Context
	connManager *connmgr.BasicConnMgr
}

var _ p2pnet.Notifiee = &notifee{}

// Listen is called when network starts listening on an addr
func (n *notifee) Listen(p2pnet.Network, ma.Multiaddr) {}

// ListenClose is called when network stops listening on an addr
func (n *notifee) ListenClose(p2pnet.Network, ma.Multiaddr) {}

// Connected is called when a connection opened
func (n *notifee) Connected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"remotePeerID":       conn.RemotePeer(),
		"remoteMultiaddress": conn.RemoteMultiaddr(),
	}).Trace("connected to peer")
}

// Disconnected is called when a connection closed
func (n *notifee) Disconnected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"remotePeerID":       conn.RemotePeer(),
		"remoteMultiaddress": conn.RemoteMultiaddr(),
	}).Trace("disconnected from peer")
}

// OpenedStream is called when a stream opened
func (n *notifee) OpenedStream(network p2pnet.Network, stream p2pnet.Stream) {
	go func() {
		ctx, cancel := context.WithTimeout(n.ctx, 5*time.Second)
		defer cancel()
		waitForStreamProtocol(ctx, stream)

		if stream.Protocol() == pubsubProtocolID {
			// When we find a peer who speaks our protocol, we give them a slight
			// positive score so the Connection Manager will be less likely to
			// disconnect them.
			log.WithFields(map[string]interface{}{
				"remotePeerID": stream.Conn().RemotePeer(),
				"protocol":     stream.Protocol(),
				"direction":    stream.Stat().Direction,
			}).Debug("found peer who speaks our protocol")
			n.connManager.TagPeer(stream.Conn().RemotePeer(), pubsubProtocolTag, pubsubProtocolScore)
		}
	}()
}

// ClosedStream is called when a stream closed
func (n *notifee) ClosedStream(network p2pnet.Network, stream p2pnet.Stream) {}

// waitForStreamProtocol blocks until the context is canceled or stream.Protocol
// is not empty.
func waitForStreamProtocol(ctx context.Context, stream p2pnet.Stream) {
	// HACK(albrow): When the stream is initially opened, the protocol is not
	// set. For now, we have to manually poll until it is set.
	// https://github.com/libp2p/go-libp2p/issues/467 mentions an internal
	// event bus which could potentially be used to detect when the protocol is
	// set or offer a different way to detect peers who speak the protocol we're
	// looking for.
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for stream.Protocol() == "" {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			continue
		}
	}
}
