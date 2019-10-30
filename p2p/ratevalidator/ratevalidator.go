package ratevalidator

import (
	"context"

	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"golang.org/x/time/rate"
)

// Dummy declaration to ensure that Validate can be used as a pubsub.Validator
var _ pubsub.Validator = (&Validator{}).Validate

// Validator is a rate-limiting pubsub validator that only allows messages to be
// sent at a certain rate.
type Validator struct {
	globalLimiter *rate.Limiter
}

// New creates and returns a new rate-limiting validator. limit is the number of
// messages that can be sent per second and burst is the maximum amount of
// messages that can be sent at once.
func New(limit rate.Limit, burst int) *Validator {
	return &Validator{
		globalLimiter: rate.NewLimiter(limit, burst),
	}
}

func (v *Validator) Validate(ctx context.Context, peerID peer.ID, msg *pubsub.Message) bool {
	// TOOD(albrow): Implement per-peer limits based on From address.
	return v.globalLimiter.Allow()
}
