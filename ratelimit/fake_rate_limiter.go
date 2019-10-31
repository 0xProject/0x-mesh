package ratelimit

import (
	"context"
	"time"
)

// FakeRateLimiter is a fake RateLimiter that always allows a request through
type FakeRateLimiter struct{}

// Wait blocks until the rateLimiter allows for another request to be sent
func (f *FakeRateLimiter) Wait(ctx context.Context) error {
	return nil
}

// Start starts the fake rateLimiter
func (f *FakeRateLimiter) Start(ctx context.Context, checkpointInterval time.Duration) error {
	return nil
}

// NewFakeRateLimiter returns a new FakeRateLimiter
func NewFakeRateLimiter() *FakeRateLimiter {
	return &FakeRateLimiter{}
}
