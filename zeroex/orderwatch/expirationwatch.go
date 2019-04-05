package orderwatch

import (
	"reflect"
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
	rbTree           *rbt.RbTree
	expirationBuffer int64
}

// NewExpirationWatcher instantiates a new expiration watcher. An expiration buffer (either positive or negative)
// can be specified, causing orders to be deemed expired some time before or after their expiration reaches current
// UTC time.
func NewExpirationWatcher(expirationBuffer int64) *ExpirationWatcher {
	rbTree := rbt.NewRbTree()
	return &ExpirationWatcher{
		rbTree:           rbTree,
		expirationBuffer: expirationBuffer,
	}
}

// Add adds a new order to the expiration watcher
func (e *ExpirationWatcher) Add(expirationTimeSeconds int64, orderHash common.Hash) {
	key := rbt.Int64Key(expirationTimeSeconds)
	e.rbTree.Insert(&key, orderHash)
}

// Prune checks for any expired orders, removes them from the expiration watcher and returns them
// to the caller
func (e *ExpirationWatcher) Prune() []ExpiredOrder {
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
