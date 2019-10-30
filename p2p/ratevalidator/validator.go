package ratevalidator

import (
	"context"
	"sync"
	"time"

	"github.com/karlseguin/ccache"
	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
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

// Validator is a rate limiting pubsub validator that only allows messages to be
// sent at a certain rate.
type Validator struct {
	ctx           context.Context
	mut           sync.RWMutex
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

// New creates and returns a new rate limiting validator.
func New(ctx context.Context, config Config) *Validator {
	validator := &Validator{
		ctx:           ctx,
		config:        config,
		globalLimiter: rate.NewLimiter(config.GlobalLimit, config.GlobalBurst),
		peerLimiters:  ccache.New(ccache.Configure().MaxSize(peerLimiterSize)),
	}
	go func() {
		// Stop and clear the cache when the context is canceled.
		select {
		case <-ctx.Done():
			// Hack(albrow): ccache.Cache.Stop() is not threadsafe so we need to
			// protect calls to stop by a mutex.
			validator.mut.Lock()
			defer validator.mut.Unlock()
			validator.peerLimiters.Clear()
			validator.peerLimiters.Stop()
		}
	}()
	return validator
}

// Validate validates a pubsub message based solely on the rate of messages
// received. If either the global or per-peer limits are exceeded, the message
// is considered "invalid" and will be dropped.
func (v *Validator) Validate(ctx context.Context, peerID peer.ID, msg *pubsub.Message) bool {
	// Hack(albrow): ccache.Cache.Stop() is not threadsafe so we need to
	// protect Get/Set/Fetch with a mutex.
	v.mut.RLock()
	defer v.mut.RUnlock()
	if v.isClosed() {
		return false
	}
	// Note: We check the per-peer rate limiter first so that peers who are
	// exceeding the limit do not contribute toward the global rate limit.
	peerLimiter, err := v.getOrCreateLimiterForPeer(peerID)
	if err != nil {
		log.WithError(err).Error("unexpected error in getOrCreateLimiterForPeer")
		return false
	}
	if !peerLimiter.Allow() {
		return false
	}

	return v.globalLimiter.Allow()
}

func (v *Validator) getOrCreateLimiterForPeer(peerID peer.ID) (*rate.Limiter, error) {
	item, err := v.peerLimiters.Fetch(peerID.String(), peerLimiterTTL, func() (interface{}, error) {
		limiter := rate.NewLimiter(v.config.PerPeerLimit, v.config.PerPeerBurst)
		return limiter, nil
	})
	if err != nil {
		return nil, err
	}
	return item.Value().(*rate.Limiter), nil
}

// isClosed returns true if the context is done and false otherwise.
func (v *Validator) isClosed() bool {
	select {
	case <-v.ctx.Done():
		return true
	default:
		return false
	}
}
