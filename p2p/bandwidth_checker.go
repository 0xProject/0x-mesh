package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/metrics"
	log "github.com/sirupsen/logrus"
)

const (
	// defaultMaxBytesPerSecond is the maximum number of bytes per second that a
	// peer is allowed to send before failing the bandwidth check.
	// TODO(albrow): Reduce this limit once we have a better picture of what
	// real-world bandwidth should be.
	defaultMaxBytesPerSecond = 104857600 // 100 MiB.
)

type bandwidthChecker struct {
	node              *Node
	counter           *metrics.BandwidthCounter
	maxBytesPerSecond float64
	// TODO(albrow): We'll use these later.
	// lastSnapshot     map[peer.ID]metrics.Stats
	// lastSnapshotTime time.Time
}

func newBandwidthChecker(node *Node, counter *metrics.BandwidthCounter) *bandwidthChecker {
	return &bandwidthChecker{
		node:              node,
		counter:           counter,
		maxBytesPerSecond: defaultMaxBytesPerSecond,
	}
}

func (checker *bandwidthChecker) logBandwidthUsage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

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

		time.Sleep(logBandwidthUsageInterval)
	}
}

// checkUsage checks the amount of data sent by each connected peer and bans
// (via IP address) any peers which have exceeded the bandwidth limit.
func (checker *bandwidthChecker) checkUsage() {
	for _, remotePeerID := range checker.node.host.Network().Peers() {
		stats := checker.counter.GetBandwidthForPeer(remotePeerID)
		// If the peer is sending is data at a higher rate than is allowed, ban
		// them.
		if stats.RateIn > checker.maxBytesPerSecond {
			log.WithFields(log.Fields{
				"remotePeerID":     remotePeerID.String(),
				"bytesPerSecondIn": stats.RateIn,
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
					"remotePeerID":    remotePeerID.String(),
					"remoteMultiaddr": conn.RemoteMultiaddr().String(),
					"rateIn":          stats.RateIn,
				}).Trace("banning IP/multiaddress due to high bandwidth usage")
			}
			// Banning the IP doesn't close the connection, so we do that
			// separately. ClosePeer closes all connections to the given peer.
			_ = checker.node.host.Network().ClosePeer(remotePeerID)
		}
	}
}
