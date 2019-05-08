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
	// Validate returns true if the message is valid and false if it is not. It
	// should only return an error if there was a problem validating the message.
	Validate(*Message) (bool, error)
	// Store stores the message. There is no guarantee that the given message is
	// unique. Store may no-op if the message is already stored, in which case it
	// should not return an error.
	Store(*Message) error
	// GetMessagesToShare returns up to max messages to be shared with peers.
	GetMessagesToShare(max int) ([][]byte, error)
}
