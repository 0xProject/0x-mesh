package loghooks

import (
	peer "github.com/libp2p/go-libp2p-peer"
	log "github.com/sirupsen/logrus"
)

// PeerIDHook is a logger hook that injects the peer ID in all logs when
// possible.
type PeerIDHook struct {
	peerID string
}

// NewPeerIDHook creates and returns a new PeerIDHook with the given peer ID.
func NewPeerIDHook(peerID peer.ID) *PeerIDHook {
	return &PeerIDHook{peerID: peer.IDB58Encode(peerID)}
}

// Ensure that PeerIDHook implements log.Hook.
var _ log.Hook = &PeerIDHook{}

func (h *PeerIDHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *PeerIDHook) Fire(entry *log.Entry) error {
	entry.Data["myPeerID"] = h.peerID
	return nil
}
