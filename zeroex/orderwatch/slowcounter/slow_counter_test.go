package slowcounter

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlowCounterWithNoDelay(t *testing.T) {
	t.Parallel()

	config := Config{
		StartingOffset:     big.NewInt(10),
		Rate:               2,
		MinDelayBeforeIncr: 0,
		MinTicksBeforeIncr: 3,
		MaxCount:           big.NewInt(1000),
	}

	testCases := []struct {
		startingCount *big.Int
		ticks         int
		expectedCount *big.Int
	}{
		{
			startingCount: big.NewInt(0),
			ticks:         0,
			expectedCount: big.NewInt(0),
		},
		{
			startingCount: big.NewInt(10),
			ticks:         0,
			expectedCount: big.NewInt(10),
		},
		{
			startingCount: big.NewInt(0),
			ticks:         2,
			expectedCount: big.NewInt(0),
		},
		{
			startingCount: big.NewInt(10),
			ticks:         2,
			expectedCount: big.NewInt(10),
		},
		{
			startingCount: big.NewInt(0),
			ticks:         3,
			expectedCount: big.NewInt(20),
		},
		{
			startingCount: big.NewInt(10),
			ticks:         3,
			expectedCount: big.NewInt(30),
		},
		{
			startingCount: big.NewInt(0),
			ticks:         6,
			expectedCount: big.NewInt(40),
		},
		{
			startingCount: big.NewInt(0),
			ticks:         18,
			expectedCount: big.NewInt(640),
		},
		{
			startingCount: big.NewInt(0),
			ticks:         21,
			expectedCount: big.NewInt(1000), // max count
		},
	}

	for _, tc := range testCases {
		counter, err := New(config, tc.startingCount)
		require.NoError(t, err)

		for i := 0; i < tc.ticks; i++ {
			counter.Tick()
		}

		actualCount := counter.Count()
		assert.Equal(t, tc.expectedCount.String(), actualCount.String(), "incorrect count (started at %s and did %d ticks)", tc.startingCount, tc.ticks)
	}
}

func TestSlowCounterWithDelay(t *testing.T) {
	t.Parallel()

	config := Config{
		StartingOffset:     big.NewInt(10),
		Rate:               2,
		MinDelayBeforeIncr: 10 * time.Millisecond,
		MinTicksBeforeIncr: 3,
		MaxCount:           big.NewInt(1000),
	}
	counter, err := New(config, big.NewInt(0))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		wasIncremented := counter.Tick()
		assert.False(t, wasIncremented, "counter should not be incremented before min delay")
	}

	time.Sleep(config.MinDelayBeforeIncr)

	{
		wasIncremented := counter.Tick()
		assert.True(t, wasIncremented, "counter should be incremented after min delay")
		expectedCount := big.NewInt(20)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count after counter was incremented")
	}

	for i := 0; i < 5; i++ {
		wasIncremented := counter.Tick()
		assert.False(t, wasIncremented, "counter should not be incremented before min delay")
	}

	time.Sleep(config.MinDelayBeforeIncr)

	{
		wasIncremented := counter.Tick()
		assert.True(t, wasIncremented, "counter should be incremented after min delay")
		expectedCount := big.NewInt(40)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count after counter was incremented")
	}
}
