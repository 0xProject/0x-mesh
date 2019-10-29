package banner

import (
	"context"

	"github.com/karlseguin/ccache"
	"github.com/libp2p/go-libp2p-core/peer"
)

// violationsTracker is used to count how many times each peer has violated the
// bandwidth limit. It is a workaround for a bug in libp2p's BandwidthCounter.
// See: https://github.com/libp2p/go-libp2p-core/issues/65.
//
// TODO(albrow): Could potentially remove this if the issue is resolved.
type violationsTracker struct {
	cache *ccache.Cache
}

func newViolationsTracker(ctx context.Context) *violationsTracker {
	cache := ccache.New(ccache.Configure().MaxSize(violationsCacheSize).ItemsToPrune(violationsCacheSize / 10))
	go func() {
		// Stop the cache when the context is done. This prevents goroutine leaks
		// since ccache spawns a new goroutine as part of its implementation.
		select {
		case <-ctx.Done():
			cache.Stop()
		}
	}()
	return &violationsTracker{
		cache: cache,
	}
}

// add increments the number of bandwidth violations by the given peer. It
// returns the new count.
func (v *violationsTracker) add(peerID peer.ID) int {
	newCount := 1
	if item := v.cache.Get(peerID.String()); item != nil {
		newCount = item.Value().(int) + 1
	}
	v.cache.Set(peerID.String(), newCount, violationsTTL)
	return newCount
}
