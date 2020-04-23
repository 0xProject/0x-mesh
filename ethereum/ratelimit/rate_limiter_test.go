package ratelimit

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultMaxRequestsPer24Hrs  = 100000
	defaultMaxRequestsPerSecond = 10.0
	defaultCheckpointInterval   = 1 * time.Minute

	// grantTimingTolerance is the maximum allowed difference between the expected
	// time for a request to be granted and the actual time it is granted. Used
	// throughout these tests to account for subtle timing differences.
	grantTimingTolerance = 50 * time.Millisecond
)

var contractAddresses = ethereum.GanacheAddresses

// Scenario1: If the 24 hour limit has *not* been hit, requests should be
// granted based on the per second limiter.
func TestScenario1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	meshDB, err := db.New(ctx, "/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
	require.NoError(t, err)
	defer meshDB.Close()
	initMetadata(t, meshDB)

	// Set up some constants for this test.
	const maxRequestsPer24Hrs = 100000
	const maxRequestsPerSecond = 10

	// Set mock clock to UTC midnight
	aClock := clock.NewMock()
	aClock.Set(GetUTCMidnightOfDate(time.Now()))

	rateLimiter, err := New(maxRequestsPer24Hrs, maxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := rateLimiter.Start(ctx, defaultCheckpointInterval)
		require.NoError(t, err)
	}()

	// First maxRequestsPerSecond/2 should be granted pretty much immediately.
	expectRequestsGranted(t, rateLimiter, int(maxRequestsPerSecond/2), 0, grantTimingTolerance)
	// Next 5 requests should be granted after 1s / maxRequestsPerSecond
	expectedDelay := (1 * time.Second) / time.Duration(maxRequestsPerSecond)
	expectRequestsGranted(t, rateLimiter, 5, expectedDelay-grantTimingTolerance, expectedDelay+grantTimingTolerance)

	cancel()
	wg.Wait()
}

// Scenario 2: Max requests per 24 hours used up. Subsequent calls to Wait
// should return an error.
func TestScenario2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	meshDB, err := db.New(ctx, "/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
	require.NoError(t, err)
	defer meshDB.Close()

	now := time.Now()
	startOfCurrentUTCDay := GetUTCMidnightOfDate(now)
	requestsRemainingInCurrentDay := 10
	requestsSentInCurrentDay := defaultMaxRequestsPer24Hrs - requestsRemainingInCurrentDay

	// Set metadata to just short of maximum requests per 24 hours.
	metadata := &db.Metadata{
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
	rateLimiter, err := New(defaultMaxRequestsPer24Hrs, math.MaxFloat64, meshDB, aClock)
	require.NoError(t, err)

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

	// Subsequent requests should result in ErrTooManyRequestsIn24Hours
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	meshDB, err := db.New(ctx, "/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
	require.NoError(t, err)
	defer meshDB.Close()

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayMidnightUTC := GetUTCMidnightOfDate(yesterday)

	// Set metadata to include an outdated `StartOfCurrentUTCDay` and an associated
	// non-zero `EthRPCRequestsSentInCurrentUTCDay`
	metadata := &db.Metadata{
		EthereumChainID:                   1337,
		MaxExpirationTime:                 constants.UnlimitedExpirationTime,
		StartOfCurrentUTCDay:              yesterdayMidnightUTC,
		EthRPCRequestsSentInCurrentUTCDay: 5000,
	}
	err = meshDB.SaveMetadata(metadata)
	require.NoError(t, err)

	aClock := clock.NewMock()
	aClock.Set(now)
	rateLimiter, err := New(defaultMaxRequestsPer24Hrs, defaultMaxRequestsPerSecond, meshDB, aClock)
	require.NoError(t, err)

	// Check that grant count and currentUTCCheckpoint were reset during instantiation
	assert.Equal(t, 0, rateLimiter.getGrantedInLast24hrsUTC())
	expectedCurrentUTCCheckpoint := GetUTCMidnightOfDate(aClock.Now())
	assert.Equal(t, expectedCurrentUTCCheckpoint, rateLimiter.getCurrentUTCCheckpoint())

	// Start the rateLimiter
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

func initMetadata(t *testing.T, meshDB *db.DB) {
	metadata := &db.Metadata{
		EthereumChainID:   1337,
		MaxExpirationTime: constants.UnlimitedExpirationTime,
	}
	err := meshDB.SaveMetadata(metadata)
	require.NoError(t, err)
}

// expectRequestsGranted calls ratelimiter.Wait until the given number of
// reqeusts are allowed. It returns an error if it waits for less than
// minDelay or longer than maxDelay for any single request.
func expectRequestsGranted(t *testing.T, rateLimiter RateLimiter, numRequests int, minDelay time.Duration, maxDelay time.Duration) {
	for i := 0; i < numRequests; i++ {
		requestedAt := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), maxDelay)
		defer cancel()
		err := rateLimiter.Wait(ctx)
		require.NoError(t, err, "unexpected error for request %d: possible context timeout meaning the request took too long (max delay was %s, actual delay was >%s)", i, maxDelay, time.Since(requestedAt))
		actualDelay := time.Since(requestedAt)
		assert.True(t, actualDelay <= maxDelay, "waited too long to grant request %d (max delay was %s, actual delay was %s)", i, maxDelay, actualDelay)
		assert.True(t, actualDelay >= minDelay, "request %d was granted too quickly (min delay was %s, actual delay was %s)", i, minDelay, actualDelay)
	}
}

func TestGetUTCMidnightOfDate(t *testing.T) {
	testCases := []struct {
		input    time.Time
		expected time.Time
	}{
		{
			// Timezone behind UTC. Day stays the same.
			input:    time.Date(1992, time.September, 29, 8, 30, 10, 99, time.FixedZone("Eastern Standard", -5*60*60)),
			expected: time.Date(1992, time.September, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			// Timezone behind UTC. Day changes forward.
			input:    time.Date(1992, time.September, 29, 20, 30, 10, 99, time.FixedZone("Eastern Standard", -5*60*60)),
			expected: time.Date(1992, time.September, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			// Timezone ahead of UTC. Day changes backward.
			input:    time.Date(2019, time.June, 18, 6, 45, 10, 99, time.FixedZone("Japan Standard", 9*60*60)),
			expected: time.Date(2019, time.June, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			// Timezone ahead of UTC. Day stays the same.
			input:    time.Date(2019, time.June, 18, 9, 45, 10, 99, time.FixedZone("Japan Standard", 9*60*60)),
			expected: time.Date(2019, time.June, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			// Timezone at UTC.
			input:    time.Date(2019, time.December, 25, 0, 45, 15, 89, time.UTC),
			expected: time.Date(2019, time.December, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		actual := GetUTCMidnightOfDate(tc.input)
		assert.Equal(t, tc.expected, actual, "input: %s", tc.input)
	}
}
