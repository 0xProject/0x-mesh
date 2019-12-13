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

const (
	defaultMaxRequestsPer24HrsWithoutBuffer = 100000
	defaultMaxRequestsPerSecond             = 10.0
	defaultCheckpointInterval               = 1 * time.Minute
	grantTimingTolerance                    = 15 * time.Millisecond
)

// Scenario1: If the 24 hour limit has not been hit, requests should be granted
// based on the per second limiter.
func TestScenario1(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()
	initMetadata(t, meshDB)

	// Set up some constants for this test.
	const maxRequestsPer24HrsWithoutBuffer = 100000
	const maxRequestsPerSecond = 10

	// Set mock clock to UTC midnight
	aClock := clock.NewMock()
	aClock.Set(GetUTCMidnightOfDate(time.Now()))

	rateLimiter, err := New(maxRequestsPer24HrsWithoutBuffer, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, defaultCheckpointInterval)
		require.NoError(t, err)
	}()

	// First maxRequestsPerSecond should be granted pretty much immediately.
	expectRequestsGranted(t, rateLimiter, int(maxRequestsPerSecond), 0, grantTimingTolerance)
	// Next 5 requests should be granted after 1s / maxRequestsPerSecond
	expectedDelay := (1 * time.Second) / time.Duration(maxRequestsPerSecond)
	expectRequestsGranted(t, rateLimiter, 5, expectedDelay-grantTimingTolerance, expectedDelay+grantTimingTolerance)

	cancel()
	wg.Wait()
}

// Scenario 2: Max requests per 24 hours used up. Subsequent calls to Wait
// should return an error.
func TestScenario2(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	now := time.Now()
	startOfCurrentUTCDay := GetUTCMidnightOfDate(now)
	requestsRemainingInCurrentDay := 10
	requestsSentInCurrentDay := defaultMaxRequestsPer24HrsWithoutBuffer - maxRequestsPer24HrsBuffer - requestsRemainingInCurrentDay

	// Set metadata to just short of maximum requests per 24 hours.
	metadata := &meshdb.Metadata{
		EthereumChainID:                   1337,
		MaxExpirationTime:                 constants.UnlimitedExpirationTime,
		StartOfCurrentUTCDay:              startOfCurrentUTCDay,
		EthRPCRequestsSentInCurrentUTCDay: requestsSentInCurrentDay,
	}
	err = meshDB.SaveMetadata(metadata)
	require.NoError(t, err)

	// Start a new rate limiter and set time to a few hours past midnight.
	// We set the max requests per second extremely high here because it's not
	// what we're trying to test.
	aClock := clock.NewMock()
	aClock.Set(startOfCurrentUTCDay.Add(3 * time.Hour))
	rateLimiter, err := New(defaultMaxRequestsPer24HrsWithoutBuffer, math.MaxFloat64, meshDB, aClock)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, defaultCheckpointInterval)
		require.NoError(t, err)
	}()

	// Up until we reach the 24 hour limit, requests should be granted pretty much
	// immediately.
	expectRequestsGranted(t, rateLimiter, requestsRemainingInCurrentDay, 0, grantTimingTolerance)

	// Subsequent reuqests should result in ErrTooManyRequestsIn24Hours
	for i := 0; i < 5; i++ {
		waitCtx, waitCancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer waitCancel()
		err := rateLimiter.Wait(waitCtx)
		require.Error(t, err, "expected ErrTooManyRequestsIn24Hours")
		assert.Equal(t, ErrTooManyRequestsIn24Hours.Error(), err.Error(), "wrong error message")
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
	yesterdayMidnightUTC := GetUTCMidnightOfDate(yesterday)

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
	aClock.Set(now)
	rateLimiter, err := New(defaultMaxRequestsPer24HrsWithoutBuffer, defaultMaxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)

	// Check that grant count and currentUTCCheckpoint were reset during instantiation
	assert.Equal(t, 0, rateLimiter.getGrantedInLast24hrsUTC())
	expectedCurrentUTCCheckpoint := GetUTCMidnightOfDate(aClock.Now())
	assert.Equal(t, expectedCurrentUTCCheckpoint, rateLimiter.getCurrentUTCCheckpoint())

	// Start the rateLimiter
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, defaultCheckpointInterval)
		require.NoError(t, err)
	}()

	// Grant a request. It should be granted immediately.
	expectRequestsGranted(t, rateLimiter, 1, 0, 5*time.Millisecond)

	// Wait for rate-limiter background process to start.
	time.Sleep(10 * time.Millisecond)

	// Advance time past the checkpointInterval
	aClock.Add(defaultCheckpointInterval + 1*time.Millisecond)

	// Wait for the metadata to be updated.
	time.Sleep(50 * time.Millisecond)

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
	_, err = New(defaultMaxRequestsPer24HrsWithoutBuffer, defaultMaxRequestsPerSecond, meshDB, aClock)
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

func expectRequestsGranted(t *testing.T, rateLimiter RateLimiter, numRequests int, minDelay time.Duration, maxDelay time.Duration) {
	for i := 0; i < numRequests; i++ {
		now := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), maxDelay)
		defer cancel()
		err := rateLimiter.Wait(ctx)
		require.NoError(t, err, "waited too long to grant request %d", i)
		elapsed := time.Since(now)
		assert.True(t, elapsed <= maxDelay, "waited too long to grant request %d", i)
		assert.True(t, elapsed >= minDelay, "request %d was granted too quickly", i)
	}
}
