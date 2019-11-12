package ratelimit

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	maxRequestsPer24HrsWithoutBuffer      = 301000
	maxRequestsPer24hrs                   = maxRequestsPer24HrsWithoutBuffer - maxRequestsPer24HrsBuffer
	maxRequestsPerSecond                  = 10.0
	twentyFourHrs                         = (24 * 60 * 60 * 1000 * time.Millisecond)
	maxExpectedDelay                      = twentyFourHrs / time.Duration(maxRequestsPer24hrs)
	minExpectedDelay                      = time.Duration(1000) / time.Duration(maxRequestsPerSecond) * time.Millisecond
	defaultCheckpointInterval             = 1 * time.Minute
	expectedMaxElapsedTimeForFirstRequest = 1 * time.Millisecond
	expectedDeltaMinExpectedDelay         = 15 * time.Millisecond
	expectedDeltaMaxExpectedDelay         = 15 * time.Millisecond
)

// Scenario1: Mesh starts X seconds after UTC midnight (start of next UTC day) and
// therefore there are Y request grants that have accrued. This test verifies that
// the first request is granted immediately, the next Y - 1 grants are issued at the
// expected minimum rate imposed by the per second rate, and all subsequent requests
// are issued at the max rate imposed by the per 24hr rate limit.
func TestScenario1(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	initMetadata(t, meshDB)

	// Set mock clock to three grants past UTC midnight
	aClock := clock.NewMock()
	now := time.Now()
	midnightUTC := getUTCMidnightOfDate(now)
	threeGrantsPastUTCMidnight := midnightUTC.Add(maxExpectedDelay * 3)
	aClock.Set(threeGrantsPastUTCMidnight)

	rateLimiter, err := New(maxRequestsPer24HrsWithoutBuffer, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, defaultCheckpointInterval)
		require.NoError(t, err)
	}()

	for i := 0; i < 5; i++ {
		now := time.Now()
		ctx := context.Background()
		err = rateLimiter.Wait(ctx)
		require.NoError(t, err)
		elapsed := time.Since(now)

		// First request goes through immediately
		if i == 0 {
			assert.True(t, elapsed < expectedMaxElapsedTimeForFirstRequest, "First request did not get granted immediately")
		} else if i > 0 && i <= 3 {
			// Subsequent requests take 1sec / maxRequestsPerSecond
			// Note: Despite initially waiting for 3 grants to accrue, by
			// the time we request the 4th, another has accrued.
			delta := math.Abs(float64(minExpectedDelay - elapsed))
			assert.True(t, time.Duration(delta) < expectedDeltaMinExpectedDelay, "Delta between minExpectedDelay and rate limit delay too large")
		} else {
			// Subsequent requests take 24hrs / maxRequestsPer24hrs
			delta := math.Abs(float64(maxExpectedDelay - elapsed))
			assert.True(t, time.Duration(delta) < expectedDeltaMaxExpectedDelay, "Delta between maxExpectedDelay and rate limit delay too large")
		}
	}

	cancel()
	wg.Wait()
}

// Scenario 2: Request grants have accrued but after 12am UTC, they get cleared
// and a new day starts
func TestScenario2(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	initMetadata(t, meshDB)

	// Set mock clock close to end of current UTC day
	aClock := clock.NewMock()
	now := time.Now()
	midnightUTC := getUTCMidnightOfDate(now)
	rightBeforeMidnight := midnightUTC.Add(twentyFourHrs - (500 * time.Millisecond))
	aClock.Set(rightBeforeMidnight)

	rateLimiter, err := New(maxRequestsPer24HrsWithoutBuffer, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, defaultCheckpointInterval)
		require.NoError(t, err)
	}()

	// Prior to re-setting the bucket for the next day, grant requests
	// should happen according to the second rate (since many are queued)
	for i := 0; i < 2; i++ {
		now := time.Now()
		ctx := context.Background()
		err = rateLimiter.Wait(ctx)
		require.NoError(t, err)
		elapsed := time.Since(now)

		// First request goes through immediately
		if i == 0 {
			assert.True(t, elapsed < expectedMaxElapsedTimeForFirstRequest, "First request did not get granted immediately")
		} else {
			// Subsequent requests take 1sec / maxRequestsPerSecond
			// Note: Despite initially waiting for 3 grants to accrue, by
			// the time we request the 4th, another has accrued.
			delta := math.Abs(float64(minExpectedDelay - elapsed))
			assert.True(t, time.Duration(delta) < expectedDeltaMinExpectedDelay, "Delta between minExpectedDelay and rate limit delay too large")
		}
	}

	// Move time forward by 500ms
	// NOTE: This does not move time forward within the rate.Limiter instances
	// we use. They unfortunately don't expose their internal clock to us
	aClock.Add(500 * time.Millisecond)

	// After moving into the next UTC day, the accrued grant requests should have been
	// cleared, causing subsequent requests to happen according to the 24hr rate (since
	// none are queued)
	for i := 0; i < 3; i++ {
		now := time.Now()
		ctx := context.Background()
		err = rateLimiter.Wait(ctx)
		require.NoError(t, err)
		elapsed := time.Since(now)

		// First request takes 500 extra miliseconds because the clock within rate.Limiter
		// is actually 500ms behind. This means the RateLimiter will attempt to
		// empty the bucket AS IF the 500ms has passed, but because it hasn't, we will
		// wait 500ms for those last grants to accrue before the bucket can clear.
		// The remaining time is close to what we'd expect from an empty bucket the needs
		// refilling.
		if i == 0 {
			expectedDuration := (500 * time.Millisecond) + maxExpectedDelay
			delta := elapsed - expectedDuration
			assert.True(t, delta < expectedDuration/10, "Delta between expected and elapsed duration too large")
		} else {
			// Subsequent requests take 24hrs / maxRequestsPer24hrs
			delta := math.Abs(float64(maxExpectedDelay - elapsed))
			assert.True(t, time.Duration(delta) < expectedDeltaMaxExpectedDelay, "Delta between maxExpectedDelay and rate limit delay too large")
		}
	}

	cancel()
	wg.Wait()
}

// Scenario 3: DB has outdated metadata values. These get overwritten when
// RateLimiter is instantiated. They then get updated after the checkpoint
// interval elapses.
func TestScenario3(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayMidnightUTC := getUTCMidnightOfDate(yesterday)

	// Set metadata to include an outdated `StartOfCurrentUTCDay` and an associated
	// non-zero `EthRPCRequestsSentInCurrentUTCDay`
	metadata := &meshdb.Metadata{
		EthereumChainID:                   1337,
		MaxExpirationTime:                 constants.UnlimitedExpirationTime,
		StartOfCurrentUTCDay:              yesterdayMidnightUTC,
		EthRPCRequestsSentInCurrentUTCDay: 5000,
	}
	err = meshDB.SaveMetadata(metadata)
	require.NoError(t, err)

	aClock := clock.NewMock()
	rateLimiter, err := New(maxRequestsPer24HrsWithoutBuffer, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)

	// Check that grant count and currentUTCCheckpoint were reset during instantiation
	assert.Equal(t, 0, rateLimiter.getGrantedInLast24hrsUTC())
	now = aClock.Now()
	expectedCurrentUTCCheckpoint := getUTCMidnightOfDate(now)
	assert.Equal(t, expectedCurrentUTCCheckpoint, rateLimiter.getCurrentUTCCheckpoint())

	ctx, cancel := context.WithCancel(context.Background())
	checkpointInterval := 200 * time.Millisecond
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, checkpointInterval)
		require.NoError(t, err)
	}()

	// Grant a request
	err = rateLimiter.Wait(ctx)
	require.NoError(t, err)

	time.Sleep(checkpointInterval + 50*time.Millisecond)

	// Check metadata was stored in DB
	metadata, err = meshDB.GetMetadata()
	require.NoError(t, err)

	assert.Equal(t, expectedCurrentUTCCheckpoint, metadata.StartOfCurrentUTCDay)
	assert.Equal(t, 1, metadata.EthRPCRequestsSentInCurrentUTCDay)

	cancel()
	wg.Wait()
}

// Scenario 4: Regression to test make sure that if local time is one day behind UTC time, the
// rate limiter still functions as expected.
func TestScenario4(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	initMetadata(t, meshDB)

	aClock := clock.NewMock()
	// Set timezone of Mock clock to `Pacific/Majuro` so that it's on the day before UTC
	now := aClock.Now()
	loc := time.FixedZone("UTC+12", 12*60*60)
	aTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	aClock.Set(aTime)

	// If we are not properly converting times to UTC, instantiation will fail with err
	// `Wait(n=450000) exceeds limiter's burst 300000`
	_, err = New(maxRequestsPer24HrsWithoutBuffer, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
}

func initMetadata(t *testing.T, meshDB *meshdb.MeshDB) {
	metadata := &meshdb.Metadata{
		EthereumChainID:   1337,
		MaxExpirationTime: constants.UnlimitedExpirationTime,
	}
	err := meshDB.SaveMetadata(metadata)
	require.NoError(t, err)
}
