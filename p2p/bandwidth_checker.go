package p2p

import (
	"context"
	"time"

	"github.com/karlseguin/ccache"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	log "github.com/sirupsen/logrus"
)

const (
	// defaultMaxBytesPerSecond is the maximum number of bytes per second that a
	// peer is allowed to send before failing the bandwidth check.
	// TODO(albrow): Reduce this limit once we have a better picture of what
	// real world bandwidth should be.
	defaultMaxBytesPerSecond = 104857600 // 100 MiB.
	// logBandwidthUsageInterval is how often to log bandwidth usage data.
	logBandwidthUsageInterval = 5 * time.Minute
	// violationsCacheSize is the size of the cache (number of entries) used for
	// tracking bandwidth violations over time.
	violationsCacheSize = 100
	// violationsBeforeBan is the number of times a peer is allowed to violate the
	// bandwidth limits before being banned.
	violationsBeforeBan = 4
	// violationsTTL is the TTL for bandwidth violations. If a peer does not have
	// any violations during this timespan, their violation count will be reset.
	violationsTTL = 6 * time.Hour
)

type bandwidthChecker struct {
	node              *Node
	counter           *metrics.BandwidthCounter
	maxBytesPerSecond float64
	violations        *violationsTracker
}

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

func newBandwidthChecker(node *Node, counter *metrics.BandwidthCounter) *bandwidthChecker {
	return &bandwidthChecker{
		node:              node,
		counter:           counter,
		maxBytesPerSecond: defaultMaxBytesPerSecond,
		violations:        newViolationsTracker(node.ctx),
	}
}

func (checker *bandwidthChecker) continuouslyLogBandwidthUsage(ctx context.Context) {
	logTicker := time.Tick(logBandwidthUsageInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-logTicker:
			// Log the bandwidth used by each peer.
			for _, remotePeerID := range checker.node.host.Network().Peers() {
				stats := checker.counter.GetBandwidthForPeer(remotePeerID)
				log.WithFields(log.Fields{
					"remotePeerID":      remotePeerID.String(),
					"bytesPerSecondIn":  stats.RateIn,
					"totalBytesIn":      stats.TotalIn,
					"bytesPerSecondOut": stats.RateOut,
					"totalBytesOut":     stats.TotalOut,
				}).Debug("bandwidth used by peer")
			}

			// Log the bandwidth used by each protocol.
			for protocolID, stats := range checker.counter.GetBandwidthByProtocol() {
				log.WithFields(log.Fields{
					"protocolID":        protocolID,
					"bytesPerSecondIn":  stats.RateIn,
					"totalBytesIn":      stats.TotalIn,
					"bytesPerSecondOut": stats.RateOut,
					"totalBytesOut":     stats.TotalOut,
				}).Debug("bandwidth used by protocol")
			}
		}
	}
}

// checkUsage checks the amount of data sent by each connected peer and bans
// (via IP address) any peers which have exceeded the bandwidth limit.
func (checker *bandwidthChecker) checkUsage() {
	for _, remotePeerID := range checker.node.host.Network().Peers() {
		stats := checker.counter.GetBandwidthForPeer(remotePeerID)
		// If the peer is sending data at a higher rate than is allowed, ban
		// them.
		if stats.RateIn > checker.maxBytesPerSecond {
			numViolations := checker.violations.add(remotePeerID)

			// Check if the number of violations exceeds violationsBeforeBan.
			if numViolations >= violationsBeforeBan {
				log.WithFields(log.Fields{
					"remotePeerID":      remotePeerID.String(),
					"bytesPerSecondIn":  stats.RateIn,
					"maxBytesPerSecond": checker.maxBytesPerSecond,
					"numViolations":     numViolations,
				}).Warn("banning peer due to high bandwidth usage")
				// There are possibly multiple connections to each peer. We ban the IP
				// address associated with each connection.
				for _, conn := range checker.node.host.Network().ConnsToPeer(remotePeerID) {
					if err := checker.node.BanIP(conn.RemoteMultiaddr()); err != nil {
						if err == errProtectedIP {
							continue
						}
						log.WithFields(log.Fields{
							"remotePeerID":    remotePeerID.String(),
							"remoteMultiaddr": conn.RemoteMultiaddr().String(),
							"error":           err.Error(),
						}).Error("could not ban peer")
					}
					log.WithFields(log.Fields{
						"remotePeerID":      remotePeerID.String(),
						"remoteMultiaddr":   conn.RemoteMultiaddr().String(),
						"rateIn":            stats.RateIn,
						"maxBytesPerSecond": checker.maxBytesPerSecond,
					}).Trace("would ban IP/multiaddress due to high bandwidth usage")
				}
				// Banning the IP doesn't close the connection, so we do that
				// separately. ClosePeer closes all connections to the given peer.
				_ = checker.node.host.Network().ClosePeer(remotePeerID)
			} else {
				// Log that high bandwidth usage occurred but don't yet ban the peer.
				log.WithFields(log.Fields{
					"remotePeerID":      remotePeerID.String(),
					"bytesPerSecondIn":  stats.RateIn,
					"maxBytesPerSecond": checker.maxBytesPerSecond,
					"numViolations":     numViolations,
				}).Warn("detected high bandwidth usage")
			}
		}
	}
}
