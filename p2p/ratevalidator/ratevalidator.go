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
	peerLimiterSize = 500
	peerLimiterTTL  = 5 * time.Minute
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

// TODO(albrow): Document this.
type Config struct {
	GlobalLimit  rate.Limit
	GlobalBurst  int
	PerPeerLimit rate.Limit
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

func (v *Validator) Validate(ctx context.Context, peerID peer.ID, msg *pubsub.Message) bool {
	select {
	case <-v.ctx.Done():
		// If the context was canceled, don't propogate any more messages. This also
		// prevents a nil pointer exception if the cache is stopped.
		return false
	default:
	}
	if !v.globalLimiter.Allow() {
		return false
	}
	peerLimiter := v.getOrCreateLimiterForPeer(peerID)
	return peerLimiter.Allow()
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
