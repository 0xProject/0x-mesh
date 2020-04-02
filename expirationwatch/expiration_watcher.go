package expirationwatch

import (
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
	expiredItems chan []ExpiredItem
	rbTreeMu     sync.RWMutex
	rbTree       *rbt.RbTree
}

// New instantiates a new expiration watcher
func New() *Watcher {
	rbTree := rbt.NewRbTree()
	return &Watcher{
		expiredItems: make(chan []ExpiredItem, 10),
		rbTree:       rbTree,
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

// Prune checks for any expired items given a timestamp and removes any expired
// items from the expiration watcher and returns them to the caller
func (w *Watcher) Prune(timestamp time.Time) []ExpiredItem {
	pruned := []ExpiredItem{}
	for {
		w.rbTreeMu.RLock()
		key, value := w.rbTree.Min()
		w.rbTreeMu.RUnlock()
		if key == nil {
			break
		}
		expirationTimeSeconds := int64(*key.(*rbt.Int64Key))
		expirationTime := time.Unix(expirationTimeSeconds, 0)
		if timestamp != expirationTime && !timestamp.After(expirationTime) {
			break
		}
		ids := value.(stringset.Set)
		for id := range ids {
			pruned = append(pruned, ExpiredItem{
				ExpirationTimestamp: expirationTime,
				ID:                  id,
			})
		}
		w.rbTreeMu.Lock()
		w.rbTree.Delete(key)
		w.rbTreeMu.Unlock()
	}
	return pruned
}
