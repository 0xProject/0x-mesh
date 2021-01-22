package ordersync_v4

import (
	peer "github.com/libp2p/go-libp2p-core/peer"
	log "github.com/sirupsen/logrus"
)

type peerScoreEvent uint

const (
	psInvalidMessage peerScoreEvent = iota
	psValidMessage
	psUnexpectedDisconnect
	receivedOrders
)

func (s *Service) handlePeerScoreEvent(id peer.ID, event peerScoreEvent) {
	// Note: for some events, we use `SetPeerScore` instead of `AddPeerScore` in
	// order to limit the maximum positive score associated with that event.
	// Without this, peers could be incentivized to artificially increase their
	// score in a way that doesn't benefit the network. (For example, they could
	// spam the network with valid messages).
	switch event {
	case psInvalidMessage:
		s.app.Node().AddPeerScore(id, "ordersync/invalid-message", -5)
	case psValidMessage:
		s.app.Node().SetPeerScore(id, "ordersync/valid-message", 5)
	case psUnexpectedDisconnect:
		s.app.Node().AddPeerScore(id, "ordersync/unexpected-disconnect", -1)
	case receivedOrders:
		s.app.Node().UnsetPeerScore(id, "ordersync/unexpected-disconnect")
		s.app.Node().SetPeerScore(id, "ordersync/received-orders", 10)
	default:
		log.WithField("event", event).Error("unknown ordersync peerScoreEvent")
	}
}
