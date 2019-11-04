package ratevalidator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/karlseguin/ccache"
	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	// peerLimiterCacheSize is the maximum number of peers to keep track of at
	// once. It controls the size of a cache that holds a rate limiter for each
	// peer.
	peerLimiterCacheSize = 500
	// peerLimiterCacheTTL is the TTL for rate limiters for each peer. If a peer
	// does not send any messages for this duration, they will be removed from the
	// cache and their rate limiter will be reset.
	peerLimiterCacheTTL = 5 * time.Minute
	// logStatsInterval is how often to log stats about rate limiting.
	logStatsInterval = 1 * time.Hour
)

// Dummy declaration to ensure that Validate can be used as a pubsub.Validator
var _ pubsub.Validator = (&Validator{}).Validate

// Validator is a rate limiting pubsub validator that only allows messages to be
// sent at a certain rate.
type Validator struct {
	mu             sync.Mutex
	config         Config
	globalLimiter  *trackingRateLimiter
	peerLimiters   *ccache.Cache
	wasStartedOnce bool
	startSignal    chan struct{}
}

// Config is a set of configuration options for the validator.
type Config struct {
	// MyPeerID is the peer ID of the host. Messages where From == MyPeerID will
	// not be rate limited and will not be counted toward the global or per-peer
	// limits.
	MyPeerID peer.ID
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
// BUG(albrow): New currently leaks goroutines due to a limitation of the
// caching library used under the hood.
func New(config Config) (*Validator, error) {
	if config.MyPeerID.String() == "" {
		return nil, errors.New("config.MyPeerID is required")
	}
	validator := &Validator{
		config:        config,
		globalLimiter: newTrackingRateLimiter(config.GlobalLimit, config.GlobalBurst),
		startSignal:   make(chan struct{}),
	}
	return validator, nil
}

// Start starts the background goroutines associated with the validator. It
// blocks until the given context is canceled, at which point it shuts down all
// goroutines and then returns.
func (v *Validator) Start(ctx context.Context) error {
	v.mu.Lock()
	if v.wasStartedOnce {
		v.mu.Unlock()
		return errors.New("Can only start Validator once per instance")
	}
	v.wasStartedOnce = true
	v.peerLimiters = ccache.New(ccache.Configure().MaxSize(peerLimiterCacheSize))
	v.mu.Unlock()
	// TODO(albrow): We should be calling Stop to cleanup any goroutines
	// started by ccache, but doing so now results in a race condition. Figure
	// out a workaround or use a different library, possibly one we write
	// ourselves.
	// go func() {
	// 	// Stop and clear the cache when the context is canceled.
	// 	select {
	// 	case <-ctx.Done():
	// 		// validator.peerLimiters.Clear()
	// 		// validator.peerLimiters.Stop()
	// 	}
	// }()
	close(v.startSignal)
	v.periodicallyLogStats(ctx)
	return nil
}

// Validate validates a pubsub message based solely on the rate of messages
// received. If either the global or per-peer limits are exceeded, the message
// is considered "invalid" and will be dropped.
func (v *Validator) Validate(ctx context.Context, peerID peer.ID, msg *pubsub.Message) bool {
	v.mu.Lock()
	if !v.wasStartedOnce {
		// Prevents nil pointer exceptions if the Validator hasn't been started yet.
		// We can't return an error here because we need to adhere to the
		// pubsub.Validator interface.
		v.mu.Unlock()
		return false
	}
	v.mu.Unlock()

	if peerID == v.config.MyPeerID {
		// Don't rate limit our own messages.
		return true
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

	return v.globalLimiter.allow()
}

func (v *Validator) getOrCreateLimiterForPeer(peerID peer.ID) (*rate.Limiter, error) {
	item, err := v.peerLimiters.Fetch(peerID.String(), peerLimiterCacheTTL, func() (interface{}, error) {
		limiter := rate.NewLimiter(v.config.PerPeerLimit, v.config.PerPeerBurst)
		return limiter, nil
	})
	if err != nil {
		return nil, err
	}
	return item.Value().(*rate.Limiter), nil
}

func (v *Validator) periodicallyLogStats(ctx context.Context) {
	ticker := time.NewTicker(logStatsInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.WithFields(log.Fields{
				"violationsCount": v.globalLimiter.violations,
			}).Debug("global PubSub rate limit violations (since last log)")
			v.globalLimiter.resetViolations()
		}
	}
}

func (v *Validator) waitForStart(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case <-v.startSignal:
		return nil
	}
}
