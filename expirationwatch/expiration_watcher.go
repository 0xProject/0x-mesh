package expirationwatch

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/albrow/stringset"
	"github.com/ocdogan/rbt"
	log "github.com/sirupsen/logrus"
)

// ExpiredItem represents an expired item returned from the Watcher
type ExpiredItem struct {
	ExpirationTimestamp time.Time
	ID                  string
}

// Watcher watches the expiration of items
type Watcher struct {
	expiredItems     chan []ExpiredItem
	rbTreeMu         sync.RWMutex
	rbTree           *rbt.RbTree
	expirationBuffer time.Duration
	ticker           *time.Ticker
	wasStartedOnce   bool
	mu               sync.Mutex
}

// New instantiates a new expiration watcher. An expiration buffer (positive or
// negative) can be specified, causing items to be deemed expired some time
// before or after their expiration reaches current UTC time. A positive
// expirationBuffer will make the item expire sooner then UTC, and a negative
// buffer after.
func New(expirationBuffer time.Duration) *Watcher {
	rbTree := rbt.NewRbTree()
	return &Watcher{
		expiredItems:     make(chan []ExpiredItem, 10),
		rbTree:           rbTree,
		expirationBuffer: expirationBuffer,
	}
}

// Add adds a new item identified by an ID to the expiration watcher
func (w *Watcher) Add(expirationTimestamp time.Time, id string) {
	key := rbt.Int64Key(expirationTimestamp.Unix())
	w.rbTreeMu.Lock()
	defer w.rbTreeMu.Unlock()
	value, ok := w.rbTree.Get(&key)
	var ids stringset.Set
	if !ok {
		ids = stringset.New()
	} else {
		ids = value.(stringset.Set)
	}
	ids.Add(id)
	w.rbTree.Insert(&key, ids)
}

// Remove removes the item with a specified id from the expiration watcher
func (w *Watcher) Remove(expirationTimestamp time.Time, id string) {
	key := rbt.Int64Key(expirationTimestamp.Unix())
	w.rbTreeMu.Lock()
	defer w.rbTreeMu.Unlock()
	value, ok := w.rbTree.Get(&key)
	if !ok {
		// Due to the asynchronous nature of the Watcher and OrderWatcher, there are
		// race-conditions where we try to remove an item from the Watcher after it
		// has already been removed.
		log.WithFields(log.Fields{
			"id": id,
		}).Trace("Attempted to remove item from Watcher that no longer exists")
		return // Noop
	} else {
		ids := value.(stringset.Set)
		ids.Remove(id)
		if len(ids) == 0 {
			w.rbTree.Delete(&key)
		} else {
			w.rbTree.Insert(&key, ids)
		}
	}
}

// Watch starts the expiration watchers poller. It continuously checks all items
// for expiration until there is an error or the given context is canceled. You
// usually want to call Watch inside a goroutine.
func (w *Watcher) Watch(ctx context.Context, pollingInterval time.Duration) error {
	w.mu.Lock()
	if w.wasStartedOnce {
		w.mu.Unlock()
		return errors.New("Can only start Watcher once per instance")
	}
	w.wasStartedOnce = true
	w.mu.Unlock()

	// TODO(fabio): Optimize this poller. We could keep track of soonestExpirationTime as a property of
	// Watcher. Whenever a new item is added via Add, we check if the expiration time is sooner
	// than soonestExpirationTime and if so, we update soonestExpirationTime. Then instead of running the
	// inner for loop at a constant frequency, we adjust the frequency based on the value of
	// soonestExpirationTime (probably by using time.After or time.Sleep).
	ticker := time.NewTicker(pollingInterval)
	for {
		select {
		case <-ctx.Done():
			close(w.expiredItems)
			return nil
		case <-ticker.C:
			expiredItems := w.prune()
			if len(expiredItems) > 0 {
				w.expiredItems <- expiredItems
			}
		}
	}
}

// ExpiredItems returns a read-only channel that can be used to listen for
// expired items. The channel will be closed if/when the watcher is done
// watching.
func (w *Watcher) ExpiredItems() <-chan []ExpiredItem {
	return w.expiredItems
}

// prune checks for any expired items, removes them from the expiration watcher and returns them
// to the caller
func (w *Watcher) prune() []ExpiredItem {
	pruned := []ExpiredItem{}
	for {
		w.rbTreeMu.RLock()
		key, value := w.rbTree.Min()
		w.rbTreeMu.RUnlock()
		if key == nil {
			break
		}
		expirationTimeSeconds := int64(*key.(*rbt.Int64Key))
		expirationTimestamp := time.Unix(expirationTimeSeconds, 0)
		currentTimePlusBuffer := time.Now().Add(w.expirationBuffer)
		if expirationTimestamp.After(currentTimePlusBuffer) {
			break
		}
		ids := value.(stringset.Set)
		for id := range ids {
			pruned = append(pruned, ExpiredItem{
				ExpirationTimestamp: expirationTimestamp,
				ID:                  id,
			})
		}
		w.rbTreeMu.Lock()
		w.rbTree.Delete(key)
		w.rbTreeMu.Unlock()
	}
	return pruned
}
