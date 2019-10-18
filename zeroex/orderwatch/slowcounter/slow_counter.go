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
	lastIncr           time.Time
	ticksSinceLastIncr int
	// currentCount tracks the current count. We use a big.Float here to preserve
	// accuracy (i.e., big.Int might not increase after an incr due to rounding).
	currentCount *big.Float
	// placeholder to minimize memory allocations.
	nextCount *big.Float
}

// Config is a set of configuration options for SlowCounter.
type Config struct {
	// Rate controls how much each increment increases the current count.
	// SlowCounter uses the exponential growth formula:
	//
	//     nextCount = currentCount * rate
	//
	Rate float64
	// rateBig is Rate converted to a big.Float in order to make the math easier.
	rateBig *big.Float
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
	MaxCount *big.Int
	// maxCountFloat is CaxCount converted to a big.Float in order to make the
	// math easier.
	maxCountFloat *big.Float
}

// New returns a new SlowCounter with the given start count.
func New(config Config, start *big.Int) (*SlowCounter, error) {
	if config.MaxCount == nil {
		return nil, errors.New("config.MaxCount cannot be nil")
	}
	config.rateBig = big.NewFloat(config.Rate)
	config.maxCountFloat = big.NewFloat(1).SetInt(config.MaxCount)
	return &SlowCounter{
		config:       config,
		lastIncr:     time.Now(),
		currentCount: big.NewFloat(0).SetInt(start),
		nextCount:    big.NewFloat(0),
	}, nil
}

// Tick processes a single tick and may increment the counter if the required
// conditions are met. Returns true if the counter was incremented.
func (sc *SlowCounter) Tick() bool {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	if sc.currentCount.Cmp(sc.config.maxCountFloat) != -1 {
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
	sc.currentCount = sc.currentCount.Mul(sc.currentCount, sc.config.rateBig)

	// If currentCount is greater than MaxCount, set it to MaxCount.
	if sc.currentCount.Cmp(sc.config.maxCountFloat) == 1 {
		sc.currentCount.Set(sc.config.maxCountFloat)
	}
}

// Count returns the current count.
func (sc *SlowCounter) Count() *big.Int {
	countInt := big.NewInt(0)
	sc.currentCount.Int(countInt)
	return countInt
}

// Reset resets the counter to the given count. This also resets the conditions
// for MinDelayBeforeIncr and MinTicksBeforeIncr.
func (sc *SlowCounter) Reset(count *big.Int) {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	sc.lastIncr = time.Now()
	sc.ticksSinceLastIncr = 0
	sc.currentCount.SetInt(count)
}
