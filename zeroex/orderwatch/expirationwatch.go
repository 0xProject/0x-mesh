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
	expirationTimeSeconds int64
	orderHash             common.Hash
}

// ExpirationWatcher watches the expiration of 0x orders
type ExpirationWatcher struct {
	ExpiredOrders    chan []ExpiredOrder
	rbTree           *rbt.RbTree
	expirationBuffer int64
	ticker           *time.Ticker
	isWatching       bool
	wasStartedOnce   bool
	mu               sync.RWMutex
}

// NewExpirationWatcher instantiates a new expiration watcher. An expiration buffer (either positive or negative)
// can be specified, causing orders to be deemed expired some time before or after their expiration reaches current
// UTC time.
func NewExpirationWatcher(expirationBuffer int64) *ExpirationWatcher {
	rbTree := rbt.NewRbTree()
	return &ExpirationWatcher{
		ExpiredOrders:    make(chan []ExpiredOrder),
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
	e.wasStartedOnce = true

	ticker := time.NewTicker(pollingInterval)
	go func() {
		for {
			<-ticker.C

			e.mu.Lock()
			if !e.isWatching {
				ticker.Stop()
				close(e.ExpiredOrders)
				e.mu.Unlock()
				return
			}
			e.mu.Unlock()

			expiredOrders := e.prune()
			go func() {
				e.ExpiredOrders <- expiredOrders
			}()
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
			expirationTimeSeconds: expirationTimeSeconds,
			orderHash:             value.(common.Hash),
		})
		e.rbTree.Delete(key)
	}
	return pruned
}
