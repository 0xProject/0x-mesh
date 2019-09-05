package core

import (
	peer "github.com/libp2p/go-libp2p-peer"
	log "github.com/sirupsen/logrus"
)

type peerScoreEvent uint

const (
	psInvalidMessage peerScoreEvent = iota
	psValidMessage
	psOrderStored
)

func (app *App) handlePeerScoreEvent(id peer.ID, event peerScoreEvent) {
	// Note: for some events, we use `SetPeerScore` instead of `AddPeerScore` in
	// order to limit the maximum positive score associated with that event.
	// Without this, peers could be incentivized to artificially increase their
	// score in a way that doesn't benefit the network. (For example, they could
	// spam the network with valid messages).
	switch event {
	case psInvalidMessage:
		app.node.AddPeerScore(id, "invalid-message", -5)
	case psValidMessage:
		app.node.SetPeerScore(id, "valid-message", 5)
	case psOrderStored:
		app.node.SetPeerScore(id, "order-stored", 10)
	default:
		log.WithField("event", event).Error("unknown peerScoreEvent")
	}
}
