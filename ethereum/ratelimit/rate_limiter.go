package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/time/rate"
)

const (
	// maxRequestsPer24HrsBuffer is the buffer subtracted from the operator supplied
	// maxRequestsPer24Hrs. This buffer helps ensure that we don't overstep the desired
	// max number of requests.
	maxRequestsPer24HrsBuffer         = 1000
	lowestPossibleMaxRequestsPer24Hrs = 40000
)

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
	twentyFourHourLimiter *rate.Limiter
	perSecondLimiter      *rate.Limiter
	currentUTCCheckpoint  time.Time // Start of current UTC 24hr period
	grantedInLast24hrsUTC int       // Number of granted requests issued in last 24hr UTC
	meshDB                *meshdb.MeshDB
	aClock                clock.Clock
	wasStartedOnce        bool         // Whether the rate limiter has previously been started
	startMut              sync.Mutex   // Mutex around the start check
	countMut              sync.Mutex   // Mutex around grantedInLast24hrsUTC
	waitMut               sync.RWMutex // Mutex for calls to Wait. Also held when swapping out rate limiters at midnight.
}

// New instantiates a new RateLimiter
func New(maxRequestsPer24HrsWithoutBuffer int, maxRequestsPerSecond float64, meshDB *meshdb.MeshDB, aClock clock.Clock) (RateLimiter, error) {
	if maxRequestsPer24HrsWithoutBuffer < lowestPossibleMaxRequestsPer24Hrs {
		return nil, fmt.Errorf("EthereumRPCMaxRequestsPer24HrUTC too low. Should be at least %d", lowestPossibleMaxRequestsPer24Hrs)
	}
	// Reduce the requested maxRequestsPer24Hrs by maxRequestsPer24HrsBuffer out of extra precaution
	maxRequestsPer24Hrs := maxRequestsPer24HrsWithoutBuffer - maxRequestsPer24HrsBuffer

	metadata, err := meshDB.GetMetadata()
	if err != nil {
		return nil, err
	}

	// Check if stored checkpoint in DB is still relevant
	now := aClock.Now()
	currentUTCCheckpoint := getUTCMidnightOfDate(now)
	storedUTCCheckpoint := metadata.StartOfCurrentUTCDay
	storedGrantedInLast24HrsUTC := metadata.EthRPCRequestsSentInCurrentUTCDay
	// Update DB if current values are from previous 24hr period and therefore no longer relevant
	if currentUTCCheckpoint != storedUTCCheckpoint {
		storedUTCCheckpoint = currentUTCCheckpoint
		storedGrantedInLast24HrsUTC = 0
		if err := meshDB.UpdateMetadata(func(metadata meshdb.Metadata) meshdb.Metadata {
			metadata.StartOfCurrentUTCDay = storedUTCCheckpoint
			metadata.EthRPCRequestsSentInCurrentUTCDay = storedGrantedInLast24HrsUTC
			return metadata
		}); err != nil {
			return nil, err
		}
	}

	// Compute the number of grants accrued since 12am UTC that have not been used. We will than
	// instantiate the rate limiter to start with the accrued grants remaining available for
	// immediate use

	timePassedSinceCheckpoint := aClock.Since(currentUTCCheckpoint)
	// Translate time passed into theoretical # grants accrued
	// (timePassed / 24hrs) * maxRequestsPer24hrs
	theoreticalGrantsAccrued := int((float64(timePassedSinceCheckpoint.Nanoseconds()) / float64((24 * time.Hour).Nanoseconds())) * float64(maxRequestsPer24Hrs))
	// theoreticalGrants - storedGrantedInLast24HrsUTC = accruedGrants
	accruedGrants := theoreticalGrantsAccrued - storedGrantedInLast24HrsUTC

	// Instantiate limiter with `maxRequestsPer24Hrs` bucketsize and a limit
	// that results in `maxRequestsPer24Hrs` requests being whitelisted in a 24hr period
	limit := rate.Limit(float64(maxRequestsPer24Hrs) / (24 * 60 * 60))
	twentyFourHourLimiter := rate.NewLimiter(limit, maxRequestsPer24Hrs)

	// Since Limiter begins initially full, we drain it before use. i.e., We do not want 100k
	// requests to already be queued up, instead we only want the number of accrued grants that
	// have gone unused to be available at startup
	amountToDrain := maxRequestsPer24Hrs - accruedGrants
	ctx := context.Background()
	err = twentyFourHourLimiter.WaitN(ctx, amountToDrain)
	if err != nil {
		return nil, err
	}

	// Instantiate limiter with a bucketsize of one and a limit that results
	// in no more than `maxRequestsPerSecond` requests per second.
	limit = rate.Limit(maxRequestsPerSecond)
	perSecondLimiter := rate.NewLimiter(limit, 1)

	return &rateLimiter{
		aClock:                aClock,
		maxRequestsPer24Hrs:   maxRequestsPer24Hrs,
		twentyFourHourLimiter: twentyFourHourLimiter,
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
	r.startMut.Lock()
	if r.wasStartedOnce {
		r.startMut.Unlock()
		return errors.New("Can only start RateLimiter once per instance")
	}
	r.wasStartedOnce = true
	r.startMut.Unlock()

	// Start 24hr UTC accrued grants resetter
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			now := r.aClock.Now()
			currentUTCCheckpoint := getUTCMidnightOfDate(now)
			nextUTCCheckpoint := time.Date(currentUTCCheckpoint.Year(), currentUTCCheckpoint.Month(), currentUTCCheckpoint.Day()+1, 0, 0, 0, 0, time.UTC)
			untilNextUTCCheckpoint := nextUTCCheckpoint.Sub(r.aClock.Now())
			select {
			case <-ctx.Done():
				return
			case <-r.aClock.After(untilNextUTCCheckpoint):
				// Create a fresh 24 hour rate limiter and drain all tokens from it.
				r.waitMut.Lock()
				r.countMut.Lock()
				limit := rate.Limit(float64(r.maxRequestsPer24Hrs) / (24 * 60 * 60))
				r.twentyFourHourLimiter = rate.NewLimiter(limit, r.maxRequestsPer24Hrs)
				if err := r.twentyFourHourLimiter.WaitN(ctx, r.maxRequestsPer24Hrs); err != nil {
					// Since we never set n to exceed the burst size, an error will only
					// occur if the context is cancelled or it's deadline is exceeded. In
					// these cases, we simply return so that this go-routine exits.
					// From docs: "It returns an error if n exceeds the Limiter's burst
					// size, the Context is canceled, or the expected wait time exceeds the
					// Context's Deadline."
					// Source: https://godoc.org/golang.org/x/time/rate#Limiter.WaitN
					r.waitMut.Unlock()
					r.countMut.Unlock()
					return
				}
				// Reset the checkpoint and current count.
				r.currentUTCCheckpoint = nextUTCCheckpoint
				r.grantedInLast24hrsUTC = 0
				r.waitMut.Unlock()
				r.countMut.Unlock()
			}
		}
	}()

	ticker := time.NewTicker(checkpointInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			wg.Wait()
			return nil
		case <-ticker.C:
			// Store grants issued and current UTC checkpoint to DB
			r.countMut.Lock()
			err := r.meshDB.UpdateMetadata(func(metadata meshdb.Metadata) meshdb.Metadata {
				metadata.StartOfCurrentUTCDay = r.currentUTCCheckpoint
				metadata.EthRPCRequestsSentInCurrentUTCDay = r.grantedInLast24hrsUTC
				return metadata
			})
			r.countMut.Unlock()
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

// Wait blocks until the rateLimiter allows for another request to be sent
func (r *rateLimiter) Wait(ctx context.Context) error {
	r.waitMut.RLock()
	defer r.waitMut.RUnlock()
	if err := r.twentyFourHourLimiter.Wait(ctx); err != nil {
		return err
	}
	if err := r.perSecondLimiter.Wait(ctx); err != nil {
		return err
	}
	r.countMut.Lock()
	r.grantedInLast24hrsUTC++
	r.countMut.Unlock()
	return nil
}

func (r *rateLimiter) getCurrentUTCCheckpoint() time.Time {
	return r.currentUTCCheckpoint
}

func (r *rateLimiter) getGrantedInLast24hrsUTC() int {
	return r.grantedInLast24hrsUTC
}

func getUTCMidnightOfDate(date time.Time) time.Time {
	utcDate := date.UTC()
	return time.Date(utcDate.Year(), utcDate.Month(), utcDate.Day(), 0, 0, 0, 0, time.UTC)
}
