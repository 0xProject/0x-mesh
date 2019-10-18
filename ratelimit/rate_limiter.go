package ratelimit

import (
	"context"
	"time"
)

// Request represents a request to the RateLimiter
type Request struct {
	ResponseChan chan struct{}
}

// RateLimiter is a rate-limiter for requests
type RateLimiter struct {
	maxRequestsPer24Hrs    int
	periodStart            time.Time
	lastRequestApprovedAt  time.Time
	ticker                 *time.Ticker
	requestsChan           chan *Request
	whitelistChan          chan struct{}
	minTimeBetweenRequests time.Duration
}

// NewRateLimiter instantiates a new RateLimiter
func NewRateLimiter(maxRequestsPer24Hrs, maxRequestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		maxRequestsPer24Hrs:    maxRequestsPer24Hrs,
		minTimeBetweenRequests: time.Duration(1000/maxRequestsPerSecond) * time.Millisecond,
		requestsChan:           make(chan *Request, 50000),
		whitelistChan:          make(chan struct{}, maxRequestsPer24Hrs),
	}
}

// Start starts the rate limiter
func (r *RateLimiter) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case request := <-r.requestsChan:
				select {
				case <-r.whitelistChan:
					timePassed := time.Since(r.lastRequestApprovedAt)
					if timePassed > r.minTimeBetweenRequests {
						r.lastRequestApprovedAt = time.Now()
						request.ResponseChan <- struct{}{}
					} else {
						timeLeft := r.minTimeBetweenRequests - timePassed
						time.Sleep(timeLeft)
						r.lastRequestApprovedAt = time.Now()
						request.ResponseChan <- struct{}{}
					}
				case <-ctx.Done():
					r.ticker.Stop()
					request.ResponseChan <- struct{}{}
					return
				}
			case <-ctx.Done():
				r.ticker.Stop()
				return
			}
		}
	}()

	r.periodStart = time.Now()
	regularTimeBetweenRequests := 24 * 60 * 60 * 1000 / r.maxRequestsPer24Hrs
	r.ticker = time.NewTicker(time.Duration(regularTimeBetweenRequests) * time.Millisecond)
	for range r.ticker.C {
		if time.Since(r.periodStart) > 24*time.Hour {
			// Drain the whitelistChan
			for {
				select {
				case <-r.whitelistChan:
				default:
					break
				}
			}
			// Restart the period
			r.periodStart = time.Now()
		}
		// Whitelist a request
		r.whitelistChan <- struct{}{}
	}
}

// Request requests to RateLimiter for permission to send a request
func (r *RateLimiter) Request(request *Request) {
	r.requestsChan <- request
}
