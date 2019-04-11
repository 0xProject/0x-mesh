package orderwatch

import (
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestPrunesExpiredOrders(t *testing.T) {
	var expirationBuffer int64 = 0
	watcher := NewExpirationWatcher(expirationBuffer)

	current := time.Now().Unix()
	expiryEntryOne := ExpiredOrder{
		expirationTimeSeconds: current - 3,
		orderHash:             common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e"),
	}
	watcher.Add(expiryEntryOne.expirationTimeSeconds, expiryEntryOne.orderHash)

	expiryEntryTwo := ExpiredOrder{
		expirationTimeSeconds: current - 1,
		orderHash:             common.HexToHash("0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521"),
	}
	watcher.Add(expiryEntryTwo.expirationTimeSeconds, expiryEntryTwo.orderHash)

	pruned := watcher.prune()
	assert.Equal(t, 2, len(pruned), "Pruned the expired order")
	assert.Equal(t, expiryEntryOne, pruned[0])
	assert.Equal(t, expiryEntryTwo, pruned[1])
}
func TestKeepsUnexpiredOrder(t *testing.T) {
	var expirationBuffer int64 = 0
	watcher := NewExpirationWatcher(expirationBuffer)

	orderHash := common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e")
	current := time.Now().Unix()
	watcher.Add(current+10, orderHash)

	pruned := watcher.prune()
	assert.Equal(t, 0, len(pruned), "Doesn't prune unexpired order")
}

func TestReturnsEmptyIfNoOrders(t *testing.T) {
	var expirationBuffer int64 = 0
	watcher := NewExpirationWatcher(expirationBuffer)

	pruned := watcher.prune()
	assert.Equal(t, 0, len(pruned), "Returns empty array when no orders tracked")
}

func TestStartsAndStopsPoller(t *testing.T) {
	var expirationBuffer int64 = 0
	watcher := NewExpirationWatcher(expirationBuffer)

	pollingInterval := 50 * time.Millisecond
	watcher.Start(pollingInterval)

	var countMux sync.Mutex
	channelCount := 0
	go func() {
		for {
			select {
			case _, isOpen := <-watcher.Receive():
				if !isOpen {
					return
				}
				countMux.Lock()
				channelCount++
				countMux.Unlock()
			}
		}
	}()

	expectedIsWatching := true
	assert.Equal(t, expectedIsWatching, watcher.isWatching)

	<-time.Tick(60 * time.Millisecond)
	watcher.Stop()
	expectedIsWatching = false
	assert.Equal(t, expectedIsWatching, watcher.isWatching)

	countMux.Lock()
	expectedChannelCount := 1
	assert.Equal(t, expectedChannelCount, channelCount)
	countMux.Unlock()
}
