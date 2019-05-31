// +build !js

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
	switch event {
	case psInvalidMessage:
		app.node.UpdatePeerScore(id, -5)
	case psValidMessage, psOrderStored:
		// For now we don't update the score for these events. Might change later.
	default:
		log.WithField("event", event).Error("unknown peerScoreEvent")
	}
}
