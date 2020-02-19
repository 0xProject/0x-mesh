package ratevalidator

import (
	"context"
	"errors"
	"time"

	"github.com/karlseguin/ccache"
	peer "github.com/libp2p/go-libp2p-core/peer"
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
	ctx           context.Context
	config        Config
	globalLimiter *trackingRateLimiter
	peerLimiters  *ccache.Cache
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
	// MaxMessageSize is the maximum size (in bytes) for a message. Any messages
	// that exceed this size will be considered invalid.
	MaxMessageSize int
}

// New creates and returns a new rate limiting validator.
// BUG(albrow): New currently leaks goroutines due to a limitation of the
// caching library used under the hood.
func New(ctx context.Context, config Config) (*Validator, error) {
	if config.MyPeerID.String() == "" {
		return nil, errors.New("config.MyPeerID is required")
	} else if config.MaxMessageSize == 0 {
		return nil, errors.New("config.MaxMessageSize is required")
	}
	validator := &Validator{
		ctx:           ctx,
		config:        config,
		globalLimiter: newTrackingRateLimiter(config.GlobalLimit, config.GlobalBurst),
		peerLimiters:  ccache.New(ccache.Configure().MaxSize(peerLimiterCacheSize)),
	}
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
	go validator.periodicallyLogStats(ctx)
	return validator, nil
}

// Validate validates a pubsub message based solely on the rate of messages
// received. If either the global or per-peer limits are exceeded, the message
// is considered "invalid" and will be dropped.
func (v *Validator) Validate(ctx context.Context, peerID peer.ID, msg *pubsub.Message) bool {
	if v.isClosed() {
		return false
	}

	if peerID == v.config.MyPeerID {
		// Don't rate limit our own messages.
		return true
	}

	if data := msg.GetData(); data != nil && len(data) > v.config.MaxMessageSize {
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

// isClosed returns true if the context is done and false otherwise.
func (v *Validator) isClosed() bool {
	select {
	case <-v.ctx.Done():
		return true
	default:
		return false
	}
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
