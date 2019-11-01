package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/time/rate"
)

// RateLimiter is a rate-limiter for requests
type RateLimiter struct {
	aClock                clock.Clock
	maxRequestsPer24Hrs   int
	twentyFourHourLimiter *rate.Limiter
	perSecondLimiter      *rate.Limiter
	meshDB                *meshdb.MeshDB
	currentUTCCheckpoint  time.Time
	grantedInLast24hrsUTC int
	wasStartedOnce        bool // Whether the rate limiter has previously been started
	startMutex            sync.Mutex
	mu                    sync.Mutex
}

// New instantiates a new RateLimiter
func New(maxRequestsPer24Hrs int, maxRequestsPerSecond float64, meshDB *meshdb.MeshDB, aClock clock.Clock) (*RateLimiter, error) {
	metadata, err := meshDB.GetMetadata()
	if err != nil {
		return nil, err
	}

	// Check if stored checkpoint in DB still relevant
	now := aClock.Now()
	currentUTCCheckpoint := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	storedUTCCheckpoint := metadata.StartOfCurrentUTCDay
	storedGrantedInLast24HrsUTC := metadata.EthRPCRequestsSentInCurrentUTCDay
	// Checkpoint in DB if from previous 24hr period and therefore no longer relevant
	if currentUTCCheckpoint != storedUTCCheckpoint {
		storedUTCCheckpoint = currentUTCCheckpoint
		storedGrantedInLast24HrsUTC = 0
		// Re-set DB state to current checkpoint and reset grants to 0
		if err := meshDB.UpdateMetadata(func(metadata meshdb.Metadata) meshdb.Metadata {
			metadata.StartOfCurrentUTCDay = storedUTCCheckpoint
			metadata.EthRPCRequestsSentInCurrentUTCDay = storedGrantedInLast24HrsUTC
			return metadata
		}); err != nil {
			return nil, err
		}
	}

	// Compute the number of grants accrued since 12am UTC that have not been used. We will than
	// instantiate the rate limiter to start with the accrued grants available for immediate use

	// compute time past since last 12am UTC
	timePassedSinceCheckpoint := aClock.Since(currentUTCCheckpoint)
	// Translate time passed into theoretical # grants accrued
	// (timePassed / 24hrs) * maxRequestsPer24hrs
	theoreticalGrantsUsed := int((float64(timePassedSinceCheckpoint.Nanoseconds()) / float64((24 * time.Hour).Nanoseconds())) * float64(maxRequestsPer24Hrs))
	// theoreticalGrants - storedGrantedInLast24HrsUTC = bufferedGrants
	bufferedGrants := theoreticalGrantsUsed - storedGrantedInLast24HrsUTC

	twentyFourHourLimiter, err := instantiateTwentyFourHourLimiter(maxRequestsPer24Hrs, bufferedGrants)
	if err != nil {
		return nil, err
	}

	// Instantiate limiter with a bucketsize of one and an eventsPerSecond rate that
	// results in no more than 30 requests per second.
	perSecondLimiter := rate.NewLimiter(rate.Limit(maxRequestsPerSecond), 1)

	return &RateLimiter{
		aClock:                aClock,
		maxRequestsPer24Hrs:   maxRequestsPer24Hrs,
		twentyFourHourLimiter: twentyFourHourLimiter,
		perSecondLimiter:      perSecondLimiter,
		meshDB:                meshDB,
		currentUTCCheckpoint:  storedUTCCheckpoint,
		grantedInLast24hrsUTC: storedGrantedInLast24HrsUTC,
	}, nil
}

// Start starts the rateLimiter
func (r *RateLimiter) Start(ctx context.Context, checkpointInterval time.Duration) error {
	r.startMutex.Lock()
	if r.wasStartedOnce {
		r.startMutex.Unlock()
		return errors.New("Can only start RateLimiter once per instance")
	}
	r.wasStartedOnce = true
	r.startMutex.Unlock()

	// Start 24hr UTC bucket resetter
	go func() {
		for {
			now := r.aClock.Now()
			currentUTCCheckpoint := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			nextUTCCheckpoint := currentUTCCheckpoint.Add(24 * time.Hour)
			untilNextUTCCheckpoint := nextUTCCheckpoint.Sub(r.aClock.Now())
			select {
			case <-ctx.Done():
				return
			case <-r.aClock.After(untilNextUTCCheckpoint):
				// Compute how many grants have buffered and remove that many from the bucket
				// to clear it for the next 24hr period
				r.mu.Lock()
				bufferedGrants := maxRequestsPer24Hrs - r.grantedInLast24hrsUTC
				if err := r.twentyFourHourLimiter.WaitN(ctx, bufferedGrants); err != nil {
					// Since we never set n to exceed the burst size, an error will only
					// occur if the context is cancelled or it's deadline is exceeded. In
					// these cases, we simply return.
					// From docs: "It returns an error if n exceeds the Limiter's burst
					// size, the Context is canceled, or the expected wait time exceeds the
					// Context's Deadline."
					// Source: https://godoc.org/golang.org/x/time/rate#Limiter.WaitN
					r.mu.Unlock()
					return
				}
				r.currentUTCCheckpoint = nextUTCCheckpoint
				r.grantedInLast24hrsUTC = 0
				r.mu.Unlock()
			}
		}
	}()

	ticker := time.NewTicker(checkpointInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			// Store grant count and current UTC checkpoint to DB
			r.mu.Lock()
			err := r.meshDB.UpdateMetadata(func(metadata meshdb.Metadata) meshdb.Metadata {
				metadata.StartOfCurrentUTCDay = r.currentUTCCheckpoint
				metadata.EthRPCRequestsSentInCurrentUTCDay = r.grantedInLast24hrsUTC
				return metadata
			})
			r.mu.Unlock()
			if err != nil {
				if err == leveldb.ErrClosed {
					// We can't continue if the database is closed. Stop the rateLimiter and
					// return an error.
					ticker.Stop()
					return err
				}
				log.WithError(err).Error("rateLimiter.Start() error encountered while updating metadata in DB")
			}
		}
	}
}

func instantiateTwentyFourHourLimiter(maxRequestsPer24Hrs, bufferedGrants int) (*rate.Limiter, error) {
	// Instantiate limiter with 100k bucketsize and an eventsPerSecond rate that
	// results in 100k requests being whitelisted in a 24hr period. This represents
	// the request per 24 UTC period
	eventsPerSecond := float64(maxRequestsPer24Hrs) / (24 * 60 * 60)
	twentyFourHourLimiter := rate.NewLimiter(rate.Limit(eventsPerSecond), maxRequestsPer24Hrs)

	// Since Limiter begins initially full, we drain it before use. i.e., We do not want 100k
	// requests to already be queued up, instead we only want the number of buffered grants that
	// have gone unused to be available at startup
	amountToDrain := maxRequestsPer24Hrs - bufferedGrants
	ctx := context.Background()
	err := twentyFourHourLimiter.WaitN(ctx, amountToDrain)
	if err != nil {
		return nil, err
	}

	return twentyFourHourLimiter, nil
}

// Wait blocks until the rateLimiter allows for another request to be sent
func (r *RateLimiter) Wait(ctx context.Context) error {
	if err := r.twentyFourHourLimiter.Wait(ctx); err != nil {
		return err
	}
	if err := r.perSecondLimiter.Wait(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	r.grantedInLast24hrsUTC++
	r.mu.Unlock()
	return nil
}
