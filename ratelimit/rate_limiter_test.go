package ratelimit

import (
	"context"
	"math"
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
	maxRequestsPer24Hrs       = 300000
	maxRequestsPerSecond      = 10.0
	twentyFourHrs             = (24 * 60 * 60 * 1000 * time.Millisecond)
	maxExpectedDelay          = twentyFourHrs / time.Duration(maxRequestsPer24Hrs)
	minExpectedDelay          = time.Duration(1000) / time.Duration(maxRequestsPerSecond) * time.Millisecond
	defaultCheckpointInterval = 1 * time.Minute
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

	// Set mock clock to start of UTC day
	aClock := clock.NewMock()
	now := time.Now()
	midnightUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	threeGrantsPastUTCMidnight := midnightUTC.Add(maxExpectedDelay * 3)
	aClock.Set(threeGrantsPastUTCMidnight)

	rateLimiter, err := New(maxRequestsPer24Hrs, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
	ctx := context.Background()
	go func() {
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
			assert.Condition(t, func() bool {
				return elapsed < 1*time.Millisecond
			})
		} else if i == 1 || i == 2 || i == 3 {
			// Subsequent requests take 1sec / maxRequestsPerSecond
			// Note: Despite initially waiting for 3 grants to accrue, by
			// the time we request the 4th, another has accrued.
			delta := math.Abs(float64(minExpectedDelay - elapsed))
			assert.Condition(t, func() bool {
				return time.Duration(delta) < 5*time.Millisecond
			})
		} else {
			// Subsequent requests take 24hrs / maxRequestsPer24hrs
			delta := math.Abs(float64(maxExpectedDelay - elapsed))
			assert.Condition(t, func() bool {
				return time.Duration(delta) < 15*time.Millisecond
			})
		}
	}

	ctx.Done()
}

// Scenario 2: Request grants have accrued but after 12am UTC, they get cleared
// and a new day starts
func TestScenario2(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	initMetadata(t, meshDB)

	// Set mock clock to start of UTC day
	aClock := clock.NewMock()
	now := time.Now()
	midnightUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	rightBeforeMidNight := midnightUTC.Add(twentyFourHrs - (500 * time.Millisecond))
	aClock.Set(rightBeforeMidNight)

	rateLimiter, err := New(maxRequestsPer24Hrs, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
	ctx := context.Background()
	go func() {
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
			assert.Condition(t, func() bool {
				return elapsed < 1*time.Millisecond
			})
		} else {
			// Subsequent requests take 1sec / maxRequestsPerSecond
			// Note: Despite initially waiting for 3 grants to accrue, by
			// the time we request the 4th, another has accrued.
			delta := math.Abs(float64(minExpectedDelay - elapsed))
			assert.Condition(t, func() bool {
				return time.Duration(delta) < 5*time.Millisecond
			})
		}
	}

	// Move time forward by 500ms
	// NOTE: This does not move time forward within the rate.Limiter instances
	// we use. They unfortunately don't expose the clock they use
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

		// First request takes 1 extra second because the clock within rate.Limiter
		// is actually 500ms behind. This means the RateLimiter will attempt to
		// empty the bucket AS IF the second has passed, but because it hasn't, we will
		// wait 500ms for those last grants to accrue before the bucket can clear.
		// The remaining time is close to what we'd expect from an empty bucket the needs
		// refilling.
		if i == 0 {
			assert.Condition(t, func() bool {
				expectedDuration := (1 * time.Second) + maxExpectedDelay
				delta := elapsed - expectedDuration
				return delta < 55*time.Millisecond
			})
		} else {
			// Subsequent requests take 24hrs / maxRequestsPer24hrs
			delta := math.Abs(float64(maxExpectedDelay - elapsed))
			assert.Condition(t, func() bool {
				return time.Duration(delta) < 15*time.Millisecond
			})
		}
	}

	ctx.Done()
}

// Scenario 3: DB has outdated metadata values. These get overwritten when
// RateLimiter is instantiated. They then get updated after the checkpoint
// interval elapses.
func TestScenario3(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	now := time.Now()
	yesterdayMidnightUTC := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)

	// Set metadata to include an outdated `StartOfCurrentUTCDay` and an associated
	// non-zero `EthRPCRequestsSentInCurrentUTCDay`
	metadata := &meshdb.Metadata{
		EthereumNetworkID:                 50,
		MaxExpirationTime:                 constants.UnlimitedExpirationTime,
		StartOfCurrentUTCDay:              yesterdayMidnightUTC,
		EthRPCRequestsSentInCurrentUTCDay: 5000,
	}
	err = meshDB.SaveMetadata(metadata)
	require.NoError(t, err)

	aClock := clock.NewMock()
	rateLimiter, err := New(maxRequestsPer24Hrs, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)

	// Check that grant count and currentUTCCheckpoint were reset during instantiation
	assert.Equal(t, 0, rateLimiter.grantedInLast24hrsUTC)
	now = aClock.Now()
	expectedCurrentUTCCheckpoint := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedCurrentUTCCheckpoint, rateLimiter.currentUTCCheckpoint)

	ctx := context.Background()
	checkpointInterval := 200 * time.Millisecond
	go func() {
		rateLimiter.Start(ctx, checkpointInterval)
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

	ctx.Done()
}

func initMetadata(t *testing.T, meshDB *meshdb.MeshDB) {
	metadata := &meshdb.Metadata{
		EthereumNetworkID: 50,
		MaxExpirationTime: constants.UnlimitedExpirationTime,
	}
	err := meshDB.SaveMetadata(metadata)
	require.NoError(t, err)
}
