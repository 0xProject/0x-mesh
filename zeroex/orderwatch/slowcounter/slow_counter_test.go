package slowcounter

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlowCounter(t *testing.T) {
	t.Parallel()

	config := Config{
		StartingOffset: big.NewInt(10),
		Rate:           2,
		Interval:       250 * time.Millisecond,
		MaxCount:       big.NewInt(1000),
	}
	counter, err := New(config, big.NewInt(0))
	require.NoError(t, err)

	{
		expectedCount := big.NewInt(0)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count before any increments")
	}

	time.Sleep(config.Interval)

	{
		expectedCount := big.NewInt(10)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count after 1 increment")
	}

	time.Sleep(config.Interval)

	{
		expectedCount := big.NewInt(20)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count after 2 increments")
	}
}

func TestSlowCounterReset(t *testing.T) {
	t.Parallel()

	config := Config{
		StartingOffset: big.NewInt(10),
		Rate:           2,
		Interval:       250 * time.Millisecond,
		MaxCount:       big.NewInt(1000),
	}
	counter, err := New(config, big.NewInt(20))
	require.NoError(t, err)

	time.Sleep(config.Interval)

	// Reset the counter and check that the count was correctly reset.
	counter.Reset(big.NewInt(30))
	{
		expectedCount := big.NewInt(30)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count after counter was reset")
	}

	time.Sleep(config.Interval)

	// Check the counter was incremented once from the new value after reset.
	{
		expectedCount := big.NewInt(40)
		actualCount := counter.Count()
		assert.Equal(t, expectedCount, actualCount, "wrong count after counter was reset and then incremented")
	}
}

func TestSlowCounterMaxCount(t *testing.T) {
	t.Parallel()

	config := Config{
		StartingOffset: big.NewInt(10),
		Rate:           2,
		// Note(albrow): For this test, we're okay with a much faster interval since
		// we don't need to be precise. We only need to check that *at least* N
		// increments have occurred. It is okay if more than N have occurred.
		Interval: 1 * time.Millisecond,
		MaxCount: big.NewInt(100),
	}

	counter, err := New(config, big.NewInt(0))
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		time.Sleep(config.Interval)
		actualCount := counter.Count()
		assert.False(t, actualCount.Cmp(config.MaxCount) == 1, "count should never exceed max count")
	}
}
