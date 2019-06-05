package orderwatch

import (
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrunesExpiredOrders(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	current := time.Now().Unix()
	expiryEntryOne := ExpiredOrder{
		ExpirationTimeSeconds: current - 3,
		OrderHash:             common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e"),
	}
	watcher.Add(expiryEntryOne.ExpirationTimeSeconds, expiryEntryOne.OrderHash)

	expiryEntryTwo := ExpiredOrder{
		ExpirationTimeSeconds: current - 1,
		OrderHash:             common.HexToHash("0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521"),
	}
	watcher.Add(expiryEntryTwo.ExpirationTimeSeconds, expiryEntryTwo.OrderHash)

	pruned := watcher.prune()
	assert.Len(t, pruned, 2, "Pruned the expired order")
	assert.Equal(t, expiryEntryOne, pruned[0])
	assert.Equal(t, expiryEntryTwo, pruned[1])
}

func TestPrunesTwoExpiredOrdersWithSameExpiration(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	current := time.Now().Unix()
	expiration := current - 3
	expiryEntryOne := ExpiredOrder{
		ExpirationTimeSeconds: expiration,
		OrderHash:             common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e"),
	}
	watcher.Add(expiryEntryOne.ExpirationTimeSeconds, expiryEntryOne.OrderHash)

	expiryEntryTwo := ExpiredOrder{
		ExpirationTimeSeconds: expiration,
		OrderHash:             common.HexToHash("0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521"),
	}
	watcher.Add(expiryEntryTwo.ExpirationTimeSeconds, expiryEntryTwo.OrderHash)

	pruned := watcher.prune()
	assert.Len(t, pruned, 2, "Pruned the expired order")
	orderHashes := map[common.Hash]bool{
		expiryEntryOne.OrderHash: true,
		expiryEntryTwo.OrderHash: true,
	}
	for _, expiredOrder := range pruned {
		assert.True(t, orderHashes[expiredOrder.OrderHash])
	}
}

func TestKeepsUnexpiredOrder(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	orderHash := common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e")
	current := time.Now().Unix()
	watcher.Add(current+10, orderHash)

	pruned := watcher.prune()
	assert.Equal(t, 0, len(pruned), "Doesn't prune unexpired order")
}

func TestReturnsEmptyIfNoOrders(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	pruned := watcher.prune()
	assert.Len(t, pruned, 0, "Returns empty array when no orders tracked")
}

func TestRemoveOnlyOrderWithSpecificExpirationTime(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	current := time.Now().Unix()
	expiryEntryOne := ExpiredOrder{
		ExpirationTimeSeconds: current - 3,
		OrderHash:             common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e"),
	}
	watcher.Add(expiryEntryOne.ExpirationTimeSeconds, expiryEntryOne.OrderHash)

	expiryEntryTwo := ExpiredOrder{
		ExpirationTimeSeconds: current - 1,
		OrderHash:             common.HexToHash("0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521"),
	}
	watcher.Add(expiryEntryTwo.ExpirationTimeSeconds, expiryEntryTwo.OrderHash)

	watcher.Remove(expiryEntryTwo.ExpirationTimeSeconds, expiryEntryTwo.OrderHash)

	pruned := watcher.prune()
	assert.Len(t, pruned, 1, "Pruned the expired order")
	assert.Equal(t, expiryEntryOne, pruned[0])
}
func TestRemoveOrderWhichSharesExpirationTimeWithOtherOrders(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	current := time.Now().Unix()
	singleExpirationTimeSeconds := current - 3
	expiryEntryOne := ExpiredOrder{
		ExpirationTimeSeconds: singleExpirationTimeSeconds,
		OrderHash:             common.HexToHash("0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e"),
	}
	watcher.Add(expiryEntryOne.ExpirationTimeSeconds, expiryEntryOne.OrderHash)

	expiryEntryTwo := ExpiredOrder{
		ExpirationTimeSeconds: singleExpirationTimeSeconds,
		OrderHash:             common.HexToHash("0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521"),
	}
	watcher.Add(expiryEntryTwo.ExpirationTimeSeconds, expiryEntryTwo.OrderHash)

	watcher.Remove(expiryEntryTwo.ExpirationTimeSeconds, expiryEntryTwo.OrderHash)

	pruned := watcher.prune()
	assert.Len(t, pruned, 1, "Pruned the expired order")
	assert.Equal(t, expiryEntryOne, pruned[0])
}

func TestStartsAndStopsPoller(t *testing.T) {
	expirationBuffer := 0 * time.Second
	watcher := NewExpirationWatcher(expirationBuffer)

	pollingInterval := 50 * time.Millisecond
	require.NoError(t, watcher.Start(pollingInterval))

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

	time.Sleep(60 * time.Millisecond)
	watcher.Stop()
	expectedIsWatching = false
	assert.Equal(t, expectedIsWatching, watcher.isWatching)

	countMux.Lock()
	expectedChannelCount := 1
	assert.Equal(t, expectedChannelCount, channelCount)
	countMux.Unlock()
}
