package ratelimit

import (
	"context"
	"time"
)

// FakeLimiter is a fake RateLimiter that always allows a request through
type FakeLimiter struct{}

// Wait blocks until the rateLimiter allows for another request to be sent
func (f *FakeLimiter) Wait(ctx context.Context) error {
	return nil
}

// Start starts the fake rateLimiter
func (f *FakeLimiter) Start(ctx context.Context, checkpointInterval time.Duration) error {
	return nil
}

// NewFakeLimiter returns a new FakeLimiter
func NewFakeLimiter() *FakeLimiter {
	return &FakeLimiter{}
}
