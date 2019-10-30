package ratevalidator

import (
	"context"
	"time"

	"github.com/karlseguin/ccache"
	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"golang.org/x/time/rate"
)

const (
	// peerLimiterSize is the maximum number of peers to keep track of at once. It
	// controls the size of a cache that holds a rate limiter for each peer.
	peerLimiterSize = 500
	// peerLimiterTTL is the TTL for rate limiters for each peer. If a peer does
	// not send any messages for this duration, they will be removed from the
	// cache and their rate limiter will be reset.
	peerLimiterTTL = 5 * time.Minute
)

// Dummy declaration to ensure that Validate can be used as a pubsub.Validator
var _ pubsub.Validator = (&Validator{}).Validate

// Validator is a rate-limiting pubsub validator that only allows messages to be
// sent at a certain rate.
type Validator struct {
	ctx           context.Context
	config        Config
	globalLimiter *rate.Limiter
	peerLimiters  *ccache.Cache
}

// Config is a set of configuration options for the validator.
type Config struct {
	// GlobalLimit is the maximum rate of messages per second across all peers.
	GlobalLimit rate.Limit
	// GlobalBurst is the maximum number of messages that can be received at once
	// from all peers.
	GlobalBurst int
	// PerPeerLimit is the maximum rate of messages for each peer.
	PerPeerLimit rate.Limit
	// PerPeerBurst is the maximum number of messages that can be received at once
	// from each peer.
	PerPeerBurst int
}

// New creates and returns a new rate-limiting validator.
func New(ctx context.Context, config Config) *Validator {
	validator := &Validator{
		ctx:           ctx,
		config:        config,
		globalLimiter: rate.NewLimiter(config.GlobalLimit, config.GlobalBurst),
		peerLimiters:  ccache.New(ccache.Configure().MaxSize(peerLimiterSize)),
	}
	go func() {
		// Stop the cache when the context is canceled.
		select {
		case <-ctx.Done():
			validator.peerLimiters.Stop()
		}
	}()
	return validator
}

// Validate validates a pubsub message based solely on the rate of messages
// received. If either the global or per-peer limits are exceeded, the message
// is considered "invalid" and will be dropped.
func (v *Validator) Validate(ctx context.Context, peerID peer.ID, msg *pubsub.Message) bool {
	select {
	case <-v.ctx.Done():
		// If the context was canceled, don't propogate any more messages. This also
		// prevents a nil pointer exception if the cache is stopped.
		return false
	default:
	}
	// Note: We check the per-peer rate limiter first so that peers who are
	// exceeding the limit do not contribute toward the global rate-limit.
	peerLimiter := v.getOrCreateLimiterForPeer(peerID)
	if !peerLimiter.Allow() {
		return false
	}

	return v.globalLimiter.Allow()
}

func (v *Validator) getOrCreateLimiterForPeer(peerID peer.ID) *rate.Limiter {
	cacheItem := v.peerLimiters.Get(peerID.String())
	if cacheItem != nil {
		return cacheItem.Value().(*rate.Limiter)
	} else {
		limiter := rate.NewLimiter(v.config.PerPeerLimit, v.config.PerPeerBurst)
		v.peerLimiters.Set(peerID.String(), limiter, peerLimiterTTL)
		return limiter
	}
}
