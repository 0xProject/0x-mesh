// +build !js

package core

import (
	"strings"

	p2pnet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

// localConnTag is the tag used with the Connection Manager to mark local
// connections.
const localConnTag = "local-connection"

type notifee struct {
	node *Node
}

var _ p2pnet.Notifiee = &notifee{}

// Listen is called when network starts listening on an addr
func (n *notifee) Listen(p2pnet.Network, ma.Multiaddr) {}

// ListenClose is called when network stops listening on an addr
func (n *notifee) ListenClose(p2pnet.Network, ma.Multiaddr) {}

// Connected is called when a connection opened
func (n *notifee) Connected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"peerID":       conn.RemotePeer(),
		"multiaddress": conn.RemoteMultiaddr(),
	}).Trace("connected to peer")
	if isLocalConn(conn) {
		// Protect local connections. This is a temporary measure which helps us do
		// ad hoc tests in a local environment by ensuring we don't disconnect from
		// any local peers.
		// TODO(albrow): Remove this once we have proper peer scoring/tagging in
		// place.
		log.WithFields(map[string]interface{}{
			"peerID":       conn.RemotePeer(),
			"multiaddress": conn.RemoteMultiaddr(),
		}).Debug("protecting local connection")
		n.node.connManager.Protect(conn.RemotePeer(), localConnTag)
	}
}

// Disconnected is called when a connection closed
func (n *notifee) Disconnected(network p2pnet.Network, conn p2pnet.Conn) {
	log.WithFields(map[string]interface{}{
		"peerID":       conn.RemotePeer(),
		"multiaddress": conn.RemoteMultiaddr(),
	}).Trace("disconnected from peer")
}

// OpenedStream is called when a stream opened
func (n *notifee) OpenedStream(p2pnet.Network, p2pnet.Stream) {}

// ClosedStream is called when a stream closed
func (n *notifee) ClosedStream(p2pnet.Network, p2pnet.Stream) {}

func isLocalConn(conn p2pnet.Conn) bool {
	ipv4Addr, err := conn.RemoteMultiaddr().ValueForProtocol(ma.P_IP4)
	if err != nil {
		return false
	}
	return strings.Contains(ipv4Addr, "localhost") || strings.Contains(ipv4Addr, "127.0.0.1")
}
