package orderwatch

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ocdogan/rbt"
)

// ExpiredOrder represents an expired order returned from the ExpirationWatcher
type ExpiredOrder struct {
	ExpirationTimeSeconds int64
	OrderHash             common.Hash
}

// ExpirationWatcher watches the expiration of 0x orders
type ExpirationWatcher struct {
	expiredOrders    chan []ExpiredOrder
	rbTree           *rbt.RbTree
	expirationBuffer int64
	ticker           *time.Ticker
	isWatching       bool
	wasStartedOnce   bool
	mu               sync.Mutex
}

// NewExpirationWatcher instantiates a new expiration watcher. An expiration buffer (positive or negative) can
// be specified, causing orders to be deemed expired some time before or after their expiration reaches current
// UTC time. A positive expirationBuffer will make the order expire sooner then UTC, and a negative buffer after.
// A relayer might want to use a positive buffer to ensure all orders on their orderbook are fillable, and a market
// maker might use a negative buffer when tracking their orders to make sure expired orders are truly unfillable.
func NewExpirationWatcher(expirationBuffer int64) *ExpirationWatcher {
	rbTree := rbt.NewRbTree()
	return &ExpirationWatcher{
		expiredOrders:    make(chan []ExpiredOrder, 10),
		rbTree:           rbTree,
		expirationBuffer: expirationBuffer,
	}
}

// Add adds a new order to the expiration watcher
func (e *ExpirationWatcher) Add(expirationTimeSeconds int64, orderHash common.Hash) {
	key := rbt.Int64Key(expirationTimeSeconds)
	e.rbTree.Insert(&key, orderHash)
}

// Start starts the expiration watchers poller
func (e *ExpirationWatcher) Start(pollingInterval time.Duration) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.isWatching {
		return errors.New("Expiration watcher already started")
	}
	if e.wasStartedOnce {
		return errors.New("Can only start ExpirationWatcher once per instance")
	}
	e.isWatching = true
	e.wasStartedOnce = true

	go func() {
		// TODO(fabio): Optimize this poller. We could keep track of soonestExpirationTime as a property of
		// ExpirationWatcher. Whenever a new order is added via Add, we check if the expiration time is sooner
		// than soonestExpirationTime and if so, we update soonestExpirationTime. Then instead of running the
		// inner for loop at a constant frequency, we adjust the frequency based on the value of
		// soonestExpirationTime (probably by using time.After or time.Sleep).
		ticker := time.NewTicker(pollingInterval)
		for {
			<-ticker.C

			e.mu.Lock()
			if !e.isWatching {
				ticker.Stop()
				close(e.expiredOrders)
				e.mu.Unlock()
				return
			}
			e.mu.Unlock()

			expiredOrders := e.prune()
			e.expiredOrders <- expiredOrders
		}
	}()
	return nil
}

// Stop stops the expiration watchers poller
func (e *ExpirationWatcher) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isWatching = false
}

// Receive returns a read-only channel that can be used to listen for expired orders
func (e *ExpirationWatcher) Receive() <-chan []ExpiredOrder {
	return e.expiredOrders
}

// prune checks for any expired orders, removes them from the expiration watcher and returns them
// to the caller
func (e *ExpirationWatcher) prune() []ExpiredOrder {
	pruned := []ExpiredOrder{}
	currentTimestamp := time.Now().Unix()
	for {
		key, value := e.rbTree.Min()
		if key == nil {
			break
		}
		expirationTimeSeconds := reflect.ValueOf(key).Elem().Int()
		if expirationTimeSeconds > currentTimestamp+e.expirationBuffer {
			break
		}
		pruned = append(pruned, ExpiredOrder{
			ExpirationTimeSeconds: expirationTimeSeconds,
			OrderHash:             value.(common.Hash),
		})
		e.rbTree.Delete(key)
	}
	return pruned
}
