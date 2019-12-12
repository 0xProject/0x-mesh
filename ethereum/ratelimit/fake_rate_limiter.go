package ratelimit

import (
	"context"
	"sync"
	"time"
)

// fakeLimiter is a fake RateLimiter that always allows a request through
type fakeLimiter struct {
	currentUTCCheckpoint  time.Time // Start of current UTC 24hr period
	grantedInLast24hrsUTC int       // Number of granted requests issued in last 24hr UTC
	mu                    sync.Mutex
}

// NewUnlimited returns a new RateLimiter without any limits. It will always
// allow requests immediately. It still keeps track of the number of requests
// that are allowed.
func NewUnlimited() RateLimiter {
	return &fakeLimiter{
		currentUTCCheckpoint:  getUTCMidnightOfDate(time.Now()),
		grantedInLast24hrsUTC: 0,
	}
}

// Start starts the fake rateLimiter
func (f *fakeLimiter) Start(ctx context.Context, checkpointInterval time.Duration) error {
	return nil
}

// Wait blocks until the rateLimiter allows for another request to be sent
func (f *fakeLimiter) Wait(ctx context.Context) error {
	f.mu.Lock()
	f.grantedInLast24hrsUTC++
	f.mu.Unlock()
	return nil
}

func (f *fakeLimiter) getGrantedInLast24hrsUTC() int {
	return f.grantedInLast24hrsUTC
}

func (f *fakeLimiter) getCurrentUTCCheckpoint() time.Time {
	return f.currentUTCCheckpoint
}
