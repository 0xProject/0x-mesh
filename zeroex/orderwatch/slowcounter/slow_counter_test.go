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
		Rate:               2.0,
		MinDelayBeforeIncr: 10 * time.Millisecond,
		MinTicksBeforeIncr: 3,
		MaxCount:           big.NewInt(1000),
	}
	startingCount := big.NewInt(100)
	sc, err := New(config, startingCount)
	require.NoError(t, err)

	// Since min ticks has not been met, the counter should not be incremented.
	for i := 0; i < config.MinTicksBeforeIncr-1; i++ {
		sc.Tick()
		assert.Equal(t, startingCount.String(), sc.Count().String(), "after %d ticks, count should not yet be incremented", i+1)
	}

	// Since the time hasn't elapsed yet, the next tick should *not* increment the
	// counter.
	sc.Tick()
	assert.Equal(t, startingCount.String(), sc.Count().String(), "count should not yet be incremented because MinDelayBeforeIncr has not passed")

	// Sleep until MinDelayBeforeIncr is satisfied.
	time.Sleep(config.MinDelayBeforeIncr + time.Since(sc.lastIncr))

	// On the next tick, the counter should be incremented once.
	expectedCount := big.NewInt(200)
	for i := 0; i < config.MinTicksBeforeIncr; i++ {
		sc.Tick()
		assert.Equal(t, expectedCount.String(), sc.Count().String(), "after %d ticks, count should be incremented once", i+config.MinTicksBeforeIncr)
	}

	// Sleep until MinDelayBeforeIncr is satisfied.
	time.Sleep(config.MinDelayBeforeIncr + time.Since(sc.lastIncr))

	// On the next tick, the counter should be incremented *twice* total.
	expectedCount = big.NewInt(400)
	sc.Tick()
	assert.Equal(t, expectedCount.String(), sc.Count().String(), "after %d ticks, count should be incremented twice", config.MinTicksBeforeIncr*2)

	// Reset the counter.
	newStart := big.NewInt(150)
	sc.Reset(newStart)
	assert.Equal(t, newStart.String(), sc.Count().String(), "after being reset, count should be the new starting value")

	// Wait for conditions to be met.
	for i := 0; i < config.MinTicksBeforeIncr-1; i++ {
		sc.Tick()
		assert.Equal(t, newStart.String(), sc.Count().String(), "after %d ticks after being reset, count should not yet be incremented", i+1)
	}
	time.Sleep(config.MinDelayBeforeIncr + time.Since(sc.lastIncr))

	// On the next tick, the counter should be incremented *once* from its *new*
	// starting value.
	expectedCount = big.NewInt(300)
	sc.Tick()
	assert.Equal(t, expectedCount.String(), sc.Count().String(), "after being reset, count should be incremented once")
}
