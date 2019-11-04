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

// BUG(albrow): newViolationsTracker currently leaks goroutines due to a
// limitation of the caching library used under the hood.
func newViolationsTracker(ctx context.Context) *violationsTracker {
	cache := ccache.New(ccache.Configure().MaxSize(violationsCacheSize).ItemsToPrune(violationsCacheSize / 10))
	// TODO(albrow): We should be calling Stop to cleanup any goroutines
	// started by ccache, but doing so now results in a race condition. Figure
	// out a workaround or use a different library, possibly one we write
	// ourselves.
	// go func() {
	// 	// Stop the cache when the context is done.
	// 	select {
	// 	case <-ctx.Done():
	// 		cache.Stop()
	// 	}
	// }()
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
