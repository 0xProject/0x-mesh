package expirationwatch

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrunesExpiredItems(t *testing.T) {
	watcher := New(0)

	current := time.Now().Truncate(time.Second)
	expiryEntryOne := ExpiredItem{
		ExpirationTimestamp: current.Add(-3 * time.Second),
		ID:                  "0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e",
	}
	watcher.Add(expiryEntryOne.ExpirationTimestamp, expiryEntryOne.ID)

	expiryEntryTwo := ExpiredItem{
		ExpirationTimestamp: current.Add(-1 * time.Second),
		ID:                  "0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521",
	}
	watcher.Add(expiryEntryTwo.ExpirationTimestamp, expiryEntryTwo.ID)

	pruned := watcher.prune()
	assert.Len(t, pruned, 2, "two expired items should get pruned")
	assert.Equal(t, expiryEntryOne, pruned[0])
	assert.Equal(t, expiryEntryTwo, pruned[1])
}

func TestPrunesTwoExpiredItemsWithSameExpiration(t *testing.T) {
	watcher := New(0)

	current := time.Now().Truncate(time.Second)
	expiration := current.Add(-3 * time.Second)
	expiryEntryOne := ExpiredItem{
		ExpirationTimestamp: expiration,
		ID:                  "0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e",
	}
	watcher.Add(expiryEntryOne.ExpirationTimestamp, expiryEntryOne.ID)

	expiryEntryTwo := ExpiredItem{
		ExpirationTimestamp: expiration,
		ID:                  "0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521",
	}
	watcher.Add(expiryEntryTwo.ExpirationTimestamp, expiryEntryTwo.ID)

	pruned := watcher.prune()
	assert.Len(t, pruned, 2, "two expired items should get pruned")
	hashes := map[string]bool{
		expiryEntryOne.ID: true,
		expiryEntryTwo.ID: true,
	}
	for _, expiredItem := range pruned {
		assert.True(t, hashes[expiredItem.ID])
	}
}

func TestKeepsUnexpiredItem(t *testing.T) {
	watcher := New(0)

	id := "0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e"
	current := time.Now().Truncate(time.Second)
	watcher.Add(current.Add(10*time.Second), id)

	pruned := watcher.prune()
	assert.Equal(t, 0, len(pruned), "Doesn't prune unexpired item")
}

func TestReturnsEmptyIfNoItems(t *testing.T) {
	watcher := New(0)

	pruned := watcher.prune()
	assert.Len(t, pruned, 0, "Returns empty array when no items tracked")
}

func TestRemoveOnlyItemWithSpecificExpirationTime(t *testing.T) {
	watcher := New(0)

	current := time.Now().Truncate(time.Second)
	expiryEntryOne := ExpiredItem{
		ExpirationTimestamp: current.Add(-3 * time.Second),
		ID:                  "0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e",
	}
	watcher.Add(expiryEntryOne.ExpirationTimestamp, expiryEntryOne.ID)

	expiryEntryTwo := ExpiredItem{
		ExpirationTimestamp: current.Add(-1 * time.Second),
		ID:                  "0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521",
	}
	watcher.Add(expiryEntryTwo.ExpirationTimestamp, expiryEntryTwo.ID)

	watcher.Remove(expiryEntryTwo.ExpirationTimestamp, expiryEntryTwo.ID)

	pruned := watcher.prune()
	assert.Len(t, pruned, 1, "two expired items should get pruned")
	assert.Equal(t, expiryEntryOne, pruned[0])
}
func TestRemoveItemWhichSharesExpirationTimeWithOtherItems(t *testing.T) {
	watcher := New(0)

	current := time.Now().Truncate(time.Second)
	singleExpirationTimestamp := current.Add(-3 * time.Second)
	expiryEntryOne := ExpiredItem{
		ExpirationTimestamp: singleExpirationTimestamp,
		ID:                  "0x8e209dda7e515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d97385e",
	}
	watcher.Add(expiryEntryOne.ExpirationTimestamp, expiryEntryOne.ID)

	expiryEntryTwo := ExpiredItem{
		ExpirationTimestamp: singleExpirationTimestamp,
		ID:                  "0x12ab7edd34515025d0c34aa61a0d1156a631248a4318576a2ce0fb408d3bee521",
	}
	watcher.Add(expiryEntryTwo.ExpirationTimestamp, expiryEntryTwo.ID)

	watcher.Remove(expiryEntryTwo.ExpirationTimestamp, expiryEntryTwo.ID)

	pruned := watcher.prune()
	assert.Len(t, pruned, 1, "two expired items should get pruned")
	assert.Equal(t, expiryEntryOne, pruned[0])
}

func TestStartAndStopPoller(t *testing.T) {
	watcher := New(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pollingInterval := 50 * time.Millisecond

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := watcher.Watch(ctx, pollingInterval)
		require.NoError(t, err)
	}()

	wg.Add(1)
	expiredCount := int32(0)
	go func() {
		defer wg.Done()
		for {
			select {
			case _, isOpen := <-watcher.ExpiredItems():
				if !isOpen {
					return
				}
				atomic.AddInt32(&expiredCount, 1)
			}
		}
	}()

	// Stop the watcher by canceling the context after a short waiting period.
	time.Sleep(60 * time.Millisecond)
	cancel()

	// Wait for all goroutines to exit.
	wg.Wait()

	assert.Equal(t, true, watcher.wasStartedOnce, "watcher.wasStartedOnce should be true")
	expectedExpiredCount := int32(0)
	assert.Equal(t, expectedExpiredCount, expiredCount, "expected no items to expire but at least one did")
}
