package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
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
	// HandleMessages is called whenever new messages are received. It should only
	// return an error if there was a problem handling the messages. It should not
	// return an error for invalid or duplicate messages.
	HandleMessages(context.Context, []*Message) error
	HandleMessagesV4(context.Context, []*Message) error
}
