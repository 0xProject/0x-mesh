package ratevalidator

import (
	"sync/atomic"

	"golang.org/x/time/rate"
)

// trackingRateLimiter is a wrapper around rate.Limiter that tracks the number
// of violations (i.e. the number of times that a request is not allowed).
type trackingRateLimiter struct {
	limiter    *rate.Limiter
	violations uint64
}

func newTrackingRateLimiter(r rate.Limit, b int) *trackingRateLimiter {
	return &trackingRateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

func (l *trackingRateLimiter) resetViolations() {
	atomic.StoreUint64(&l.violations, 0)
}

func (l *trackingRateLimiter) allow() bool {
	allowed := l.limiter.Allow()
	if !allowed {
		atomic.AddUint64(&l.violations, 1)
	}
	return allowed
}
