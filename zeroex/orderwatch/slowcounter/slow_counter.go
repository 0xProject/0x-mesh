package slowcounter

import (
	"errors"
	"math/big"
	"sync"
	"time"
)

// SlowCounter is an exponentially increasing counter that is only incremented
// after a certain number of "ticks" and/or a minimum time duration. It has a
// few configuration  options to control the rate of increase. SlowCounter
// uses the following exponential growth formula:
//
//    currentCount(n) = startingCount + startingOffset * (rate ^ n)
//
// where n is the number of increments that have occurred.
type SlowCounter struct {
	mut                sync.Mutex
	config             Config
	startingCount      *big.Int
	lastIncr           time.Time
	ticksSinceLastIncr int
	// currentOffset tracks the current offset. We use a big.Float here to
	// preserve accuracy (i.e., big.Int might not increase after an incr due to
	// rounding).
	currentOffset *big.Float
	currentCount  *big.Float
}

// Config is a set of configuration options for SlowCounter.
type Config struct {
	// StartingOffset affects how much the count is increased on the first
	// increment.
	StartingOffset *big.Int
	// Rate controls how fast the offset increases after each increment.
	Rate float64
	// MinDelayBeforeIncr is the minum amount of time to wait before the counter
	// is incremented. Both MinDelayBeforeIncr and MinTicksBeforeIncr conditions
	// must be satisfied in order for the counter to be incremented.
	MinDelayBeforeIncr time.Duration
	// MinTicksBeforeIncr is the minimum number of ticks befer the counter is
	// incremented. Both MinDelayBeforeIncr and MinTicksBeforeIncr conditions
	// must be satisfied in order for the counter to be incremented.
	MinTicksBeforeIncr int
	// MaxCount is the maximum value for the counter. After reaching MaxCount, the
	// counter will stop incrementing until reset.
	MaxCount *big.Int

	// rateBig is Rate converted to a big.Float in order to make the math easier.
	rateBig *big.Float
	// maxCountFloat is CaxCount converted to a big.Float in order to make the
	// math easier.
	maxCountFloat *big.Float
}

// New returns a new SlowCounter with the given starting count.
func New(config Config, startingCount *big.Int) (*SlowCounter, error) {
	if config.MaxCount == nil {
		return nil, errors.New("config.MaxCount cannot be nil")
	}
	config.rateBig = big.NewFloat(config.Rate)
	config.maxCountFloat = big.NewFloat(1).SetInt(config.MaxCount)
	return &SlowCounter{
		config:        config,
		startingCount: big.NewInt(0).Set(startingCount),
		lastIncr:      time.Now(),
		currentOffset: big.NewFloat(0).SetInt(config.StartingOffset),
		currentCount:  big.NewFloat(0).SetInt(startingCount),
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

	// Use the exponential growth forumula to adjust currentOffset.
	//
	//    currentOffset = currentOffset * rate
	//    currentCount = startingCount + currentOffset
	//
	sc.currentOffset.Mul(sc.currentOffset, sc.config.rateBig)
	sc.currentCount.SetInt(sc.startingCount)
	sc.currentCount.Add(sc.currentCount, sc.currentOffset)

	// If currentCount is greater than MaxCount, set it to MaxCount.
	if sc.currentCount.Cmp(sc.config.maxCountFloat) == 1 {
		sc.currentCount.Set(sc.config.maxCountFloat)
	}
}

// Count returns the current count.
func (sc *SlowCounter) Count() *big.Int {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	currentCountInt := big.NewInt(0)
	sc.currentCount.Int(currentCountInt)
	return currentCountInt
}

// Reset resets the counter to the given count. This also sets startingCount to
// count, sets currentOffset to config.StartingOffset, and resets the conditions
// for MinDelayBeforeIncr and MinTicksBeforeIncr.
func (sc *SlowCounter) Reset(count *big.Int) {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	sc.startingCount.Set(count)
	sc.lastIncr = time.Now()
	sc.ticksSinceLastIncr = 0
	sc.currentOffset.SetInt(sc.config.StartingOffset)
	sc.currentCount.SetInt(count)
}
