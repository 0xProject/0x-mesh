package core

import (
	peer "github.com/libp2p/go-libp2p-peer"
)

type Message struct {
	// From is the peer ID of the peer who sent the message.
	From peer.ID
	// Data is the underlying data for the message.
	Data []byte
}

// MessageHandler is an interface responsible for validating and storing
// messages as well as selecting messages which are ready to be shared.
type MessageHandler interface {
	// Validate filters out any invalid messages in the given slice of messages
	// and returns only those that are valid. It should only return an error if
	// there was a problem validating the messages.
	Validate([]*Message) ([]*Message, error)
	// Store stores the given messages. There is no guarantee that each message is
	// unique. Store may no-op if a message is already stored, in which case it
	// should not return an error.
	Store([]*Message) error
	// GetMessagesToShare returns up to max messages to be shared with peers.
	GetMessagesToShare(max int) ([][]byte, error)
}
