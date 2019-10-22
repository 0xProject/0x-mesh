package slowcounter

import (
	"errors"
	"math/big"
	"sync"
	"time"
)

// SlowCounter is an exponentially increasing counter that is slowly incremented
// after a certain time interval, unless reset. It has a few configuration
// options to control the rate of increase. SlowCounter uses the following
// exponential growth formula:
//
//    currentCount = startingCount + offset * (rate ^ n)
//
// where n is the number of increments that have occurred. And the number of
// increments is calculated as:
//
//    n = math.Floor(time.Since(startTime) / interval)
//
type SlowCounter struct {
	mut           sync.Mutex
	config        Config
	startingCount *big.Int
	// startingTime is the time the counter was started or reset.
	startingTime time.Time
	// isMax is a boolean cache which is used to prevent any computation from
	// occurring if the counter has already hit MaxCount.
	isMax bool
}

// Config is a set of configuration options for SlowCounter.
type Config struct {
	// Offset affects how much the count is increased on the first
	// increment.
	Offset *big.Int
	// Rate controls how fast the offset increases after each increment.
	Rate float64
	// Interval is the amount of time to wait before each time the counter is
	// incremented.
	Interval time.Duration
	// MaxCount is the maximum value for the counter. After reaching MaxCount, the
	// counter will stop incrementing until reset.
	MaxCount *big.Int

	// maxCountFloat is CaxCount converted to a big.Float in order to make the
	// math easier.
	maxCountFloat *big.Float
}

// New returns a new SlowCounter with the given starting count.
func New(config Config, startingCount *big.Int) (*SlowCounter, error) {
	if config.MaxCount == nil {
		return nil, errors.New("config.MaxCount cannot be nil")
	} else if config.Interval == 0 {
		return nil, errors.New("config.Interval cannot be 0")
	}
	config.maxCountFloat = big.NewFloat(1).SetInt(config.MaxCount)
	return &SlowCounter{
		config:        config,
		startingCount: big.NewInt(0).Set(startingCount),
		startingTime:  time.Now(),
	}, nil
}

// Count returns the current count.
func (sc *SlowCounter) Count() *big.Int {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	if sc.isMax {
		currentCount := big.NewInt(0).Set(sc.config.MaxCount)
		return currentCount
	}

	// TODO(albrow): Could be further optimized to reduce memory allocations and
	// math/big operations.
	//
	// currentCount = startingCount + offset * (rate ^ numIncrements)
	//
	numIncrements := time.Since(sc.startingTime) / sc.config.Interval
	if numIncrements == 0 {
		currentCount := big.NewInt(0).Set(sc.startingCount)
		return currentCount
	}
	currentCount := big.NewFloat(0).SetInt(sc.startingCount)
	offset := big.NewFloat(0).SetInt(sc.config.Offset)
	rate := big.NewFloat(sc.config.Rate)
	for i := 0; i < int(numIncrements)-1; i++ {
		offset.Mul(offset, rate)
	}
	currentCount.Add(currentCount, offset)
	currentCountInt := big.NewInt(0)
	currentCount.Int(currentCountInt)

	// If current count exceeds max, return max.
	if currentCountInt.Cmp(sc.config.MaxCount) == 1 {
		currentCountInt.Set(sc.config.MaxCount)
		sc.isMax = true
	}

	return currentCountInt
}

// Reset resets the counter to the given count.
func (sc *SlowCounter) Reset(count *big.Int) {
	sc.mut.Lock()
	defer sc.mut.Unlock()

	sc.isMax = false
	sc.startingCount.Set(count)
	sc.startingTime = time.Now()
}
