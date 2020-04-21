package ratelimit

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/time/rate"
)

var ErrTooManyRequestsIn24Hours = errors.New("too many Ethereum RPC requests have been sent this 24 hour period")

// RateLimiter is the interface one must satisfy to be considered a RateLimiter
type RateLimiter interface {
	Wait(ctx context.Context) error
	Start(ctx context.Context, checkpointInterval time.Duration) error
	getCurrentUTCCheckpoint() time.Time
	getGrantedInLast24hrsUTC() int
}

// rateLimiter is a rate-limiter for requests
type rateLimiter struct {
	maxRequestsPer24Hrs   int
	perSecondLimiter      *rate.Limiter
	currentUTCCheckpoint  time.Time // Start of current UTC 24hr period
	grantedInLast24hrsUTC int       // Number of granted requests issued in last 24hr UTC
	meshDB                *db.DB
	aClock                clock.Clock
	wasStartedOnce        bool       // Whether the rate limiter has previously been started
	startMutex            sync.Mutex // Mutex around the start check
	mu                    sync.Mutex
}

// New instantiates a new RateLimiter
func New(maxRequestsPer24Hrs int, maxRequestsPerSecond float64, meshDB *db.DB, aClock clock.Clock) (RateLimiter, error) {
	metadata, err := meshDB.GetMetadata()
	if err != nil {
		return nil, err
	}

	// Check if stored checkpoint in DB is still relevant
	now := aClock.Now()
	currentUTCCheckpoint := GetUTCMidnightOfDate(now)
	storedUTCCheckpoint := metadata.StartOfCurrentUTCDay
	storedGrantedInLast24HrsUTC := metadata.EthRPCRequestsSentInCurrentUTCDay
	// Update DB if current values are from previous 24hr period and therefore no longer relevant
	if currentUTCCheckpoint != storedUTCCheckpoint {
		storedUTCCheckpoint = currentUTCCheckpoint
		storedGrantedInLast24HrsUTC = 0
		if err := meshDB.UpdateMetadata(func(metadata db.Metadata) db.Metadata {
			metadata.StartOfCurrentUTCDay = storedUTCCheckpoint
			metadata.EthRPCRequestsSentInCurrentUTCDay = storedGrantedInLast24HrsUTC
			return metadata
		}); err != nil {
			return nil, err
		}
	}

	// Instantiate limiter with a bucketsize of maxRequestsPerSecond/2 and a limit
	// of `maxRequestsPerSecond` requests per second. This does a pretty good job
	// of limiting the number of requests we send per second while still allowing
	// for some bursts.
	limit := rate.Limit(maxRequestsPerSecond)
	perSecondLimiter := rate.NewLimiter(limit, int(math.Max(1, maxRequestsPerSecond/2)))

	return &rateLimiter{
		aClock:                aClock,
		maxRequestsPer24Hrs:   maxRequestsPer24Hrs,
		perSecondLimiter:      perSecondLimiter,
		meshDB:                meshDB,
		currentUTCCheckpoint:  storedUTCCheckpoint,
		grantedInLast24hrsUTC: storedGrantedInLast24HrsUTC,
	}, nil
}

// Start starts two background processes required for the RateLimiter to function. One that
// stores it's state to the DB at a checkpoint interval, and another that clears accrued
// grants when the UTC day time window elapses.
func (r *rateLimiter) Start(ctx context.Context, checkpointInterval time.Duration) error {
	r.startMutex.Lock()
	if r.wasStartedOnce {
		r.startMutex.Unlock()
		return errors.New("Can only start RateLimiter once per instance")
	}
	r.wasStartedOnce = true
	r.startMutex.Unlock()

	// Start 24hr UTC accrued grants resetter
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			now := r.aClock.Now()
			currentUTCCheckpoint := GetUTCMidnightOfDate(now)
			nextUTCCheckpoint := time.Date(currentUTCCheckpoint.Year(), currentUTCCheckpoint.Month(), currentUTCCheckpoint.Day()+1, 0, 0, 0, 0, time.UTC)
			untilNextUTCCheckpoint := nextUTCCheckpoint.Sub(r.aClock.Now())
			select {
			case <-ctx.Done():
				return
			case <-r.aClock.After(untilNextUTCCheckpoint):
				// Reset the number of requests granted and set the next UTC
				// checkpoint.
				r.mu.Lock()
				r.currentUTCCheckpoint = nextUTCCheckpoint
				r.grantedInLast24hrsUTC = 0
				r.mu.Unlock()
			}
		}
	}()

	ticker := r.aClock.Ticker(checkpointInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			wg.Wait()
			return nil
		case <-ticker.C:
			// Store grants issued and current UTC checkpoint to DB
			r.mu.Lock()
			err := r.meshDB.UpdateMetadata(func(metadata db.Metadata) db.Metadata {
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
					wg.Wait()
					return err
				}
				log.WithError(err).Error("rateLimiter.Start() error encountered while updating metadata in DB")
			}
		}
	}
}

// Wait blocks until the rateLimiter allows for another request to be sent. It
// returns an error if the deadline of the given context is before the request
// would be granted. It also returns an error if too many requests have been
// sent during this 24 hour period.
func (r *rateLimiter) Wait(ctx context.Context) error {
	r.mu.Lock()
	if r.grantedInLast24hrsUTC >= r.maxRequestsPer24Hrs {
		r.mu.Unlock()
		return ErrTooManyRequestsIn24Hours
	}
	r.mu.Unlock()
	if err := r.perSecondLimiter.Wait(ctx); err != nil {
		return err
	}
	r.mu.Lock()
	r.grantedInLast24hrsUTC++
	r.mu.Unlock()
	return nil
}

func (r *rateLimiter) getCurrentUTCCheckpoint() time.Time {
	return r.currentUTCCheckpoint
}

func (r *rateLimiter) getGrantedInLast24hrsUTC() int {
	return r.grantedInLast24hrsUTC
}

// Rounds the current date and time to midnight of the current day.
func GetUTCMidnightOfDate(date time.Time) time.Time {
	utcDate := date.UTC()
	return time.Date(utcDate.Year(), utcDate.Month(), utcDate.Day(), 0, 0, 0, 0, time.UTC)
}
