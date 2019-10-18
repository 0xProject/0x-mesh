package slowcounter

import (
	"errors"
	"math/big"
	"sync"
	"time"
)

// SlowCounter is an exponentially increasing counter that is only incremented
// after a certain number of "ticks" and/or a minimum time duration. It has a
// few configuration  options to control the rate of increase.
type SlowCounter struct {
	mut                sync.Mutex
	config             Config
	rateRat            *big.Rat
	lastIncr           time.Time
	ticksSinceLastIncr int
	interval           float64
	currentCount       *big.Rat
	// placeholder to minimize memory allocations.
	nextCount *big.Rat
}

// Config is a set of configuration options for SlowCounter.
type Config struct {
	// Rate controls how much each increment increases the current count.
	// SlowCounter uses the exponential growth formula:
	//
	//     nextCount = currentCount * rate
	//
	Rate *big.Rat
	// MinDelayBeforeIncr is the minum amount of time to wait before the counter
	// is incremented. Both MinDelayBeforeIncr and MinTicksBeforeIncr conditions
	// must be satisfied in order for the counter to be incremented.
	MinDelayBeforeIncr time.Duration
	// MinTicksBeforeIncr is the minimum number of ticks befer the counter is
	// incremented. Both MinDelayBeforeIncr and MinTicksBeforeIncr conditions
	// must be satisfied in order for the counter to be incremented.
	MinTicksBeforeIncr int
	// MaxCount is the maximum value for the counter. After reaching MaxCount, the
	// counter will stop incrementing.
	MaxCount *big.Rat
}

// New returns a new SlowCounter with the given start count.
func New(config Config, start *big.Rat) (*SlowCounter, error) {
	if config.MaxCount == nil {
		return nil, errors.New("config.MaxCount cannot be nil")
	}
	return &SlowCounter{
		config:       config,
		lastIncr:     time.Now(),
		currentCount: start,
		nextCount:    big.NewRat(1, 1),
	}, nil
}

// Tick processes a single tick and may increment the counter if the required
// conditions are met. Returns true if the counter was incremented.
func (sc *SlowCounter) Tick() bool {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	if sc.currentCount.Cmp(sc.config.MaxCount) == 0 {
		// Count is already at maximum. Don't need to do anything.
		return false
	}

	sc.ticksSinceLastIncr += 1
	minTicksHaveOccurred := sc.ticksSinceLastIncr >= sc.config.MinTicksBeforeIncr
	minTimeHasPassed := time.Now().After(sc.lastIncr.Add(sc.config.MinDelayBeforeIncr))
	if minTicksHaveOccurred && minTimeHasPassed {
		sc.incr()
		return true
	}

	return false
}

// incr increments the counter.
func (sc *SlowCounter) incr() {
	sc.lastIncr = time.Now()
	sc.ticksSinceLastIncr = 0

	// Use the exponential growth forumula to determine nextCount. The following
	// is written to minimize new memory allocations.
	sc.currentCount = sc.currentCount.Mul(sc.currentCount, sc.config.Rate)

	// If currentCount is greater than MaxCount, set it to MaxCount.
	if sc.currentCount.Cmp(sc.config.MaxCount) == 1 {
		sc.currentCount.Set(sc.config.MaxCount)
	}
}

// Count returns the current count.
func (sc *SlowCounter) Count() *big.Rat {
	return sc.currentCount
}

// Reset resets the counter to the given count. This also resets the conditions
// for MinDelayBeforeIncr and MinTicksBeforeIncr.
func (sc *SlowCounter) Reset(count *big.Rat) {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	sc.lastIncr = time.Now()
	sc.ticksSinceLastIncr = 0
	sc.interval = 0
	sc.currentCount.Set(count)
}
