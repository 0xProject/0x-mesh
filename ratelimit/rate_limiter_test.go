package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	maxRequestsPer24Hrs  = 200000
	maxRequestsPerSecond = 5
)

func TestSingleRequest(t *testing.T) {
	rateLimiter := NewRateLimiter(maxRequestsPer24Hrs, maxRequestsPerSecond)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		rateLimiter.Start(ctx)
	}()
	responseChan := make(chan struct{})
	r := &Request{
		ResponseChan: responseChan,
	}
	rateLimiter.Request(r)
	<-r.ResponseChan
	cancelFunc()
}

func TestProcessesRequestsImmediatelyIfBufferedRequests(t *testing.T) {
	rateLimiter := NewRateLimiter(maxRequestsPer24Hrs, maxRequestsPerSecond)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		rateLimiter.Start(ctx)
	}()

	minTimeBetweenRequests := time.Duration(1000/maxRequestsPerSecond) * time.Millisecond

	regularTimeBetweenRequests := time.Duration(24*60*60*1000/maxRequestsPer24Hrs) * time.Millisecond
	// Wait for two request whitelist slots to buffer
	time.Sleep(regularTimeBetweenRequests * 3)

	wg := &sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		time.Sleep(10 * time.Millisecond)
		go func(i int) {
			defer wg.Done()
			responseChan := make(chan struct{})
			r := &Request{
				ResponseChan: responseChan,
			}
			rateLimiter.Request(r)
			start := time.Now()
			<-r.ResponseChan
			elapsed := time.Since(start)
			expectedElapsed := time.Duration(i) * minTimeBetweenRequests
			diff := elapsed - expectedElapsed
			// Subsequent requests should wait at least minTimeBetweenRequests
			assert.Condition(t, func() bool {
				return diff < 15*time.Millisecond
			})
		}(i)
	}
	wg.Wait()
	cancelFunc()
}

func TestTakesRegularIntervalWithoutBufferedRequests(t *testing.T) {
	rateLimiter := NewRateLimiter(maxRequestsPer24Hrs, maxRequestsPerSecond)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		rateLimiter.Start(ctx)
	}()

	regularTimeBetweenRequests := time.Duration(24*60*60*1000/maxRequestsPer24Hrs) * time.Millisecond

	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			responseChan := make(chan struct{})
			r := &Request{
				ResponseChan: responseChan,
			}
			rateLimiter.Request(r)
			start := time.Now()
			<-r.ResponseChan
			elapsed := time.Since(start)
			expectedElapsed := time.Duration(i+1) * regularTimeBetweenRequests
			diff := elapsed - expectedElapsed
			assert.Condition(t, func() bool {
				return diff < 10*time.Millisecond
			})
		}(i)
	}
	wg.Wait()
	cancelFunc()
}

// TODO:
// - Test resetting after the 24hr period ends, and starts new period
// - Test cancelling the context while requests are pending
// - Add step of requesting permission before every ETH JSON RPC request
