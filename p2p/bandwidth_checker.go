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
	// real world bandwidth should be.
	defaultMaxBytesPerSecond = 104857600 // 100 MiB.
	// logBandwidthUsageInterval is how often to log bandwidth usage data.
	logBandwidthUsageInterval = 5 * time.Minute
)

type bandwidthChecker struct {
	node              *Node
	counter           *metrics.BandwidthCounter
	maxBytesPerSecond float64
}

func newBandwidthChecker(node *Node, counter *metrics.BandwidthCounter) *bandwidthChecker {
	return &bandwidthChecker{
		node:              node,
		counter:           counter,
		maxBytesPerSecond: defaultMaxBytesPerSecond,
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
			log.WithFields(log.Fields{
				"remotePeerID":      remotePeerID.String(),
				"bytesPerSecondIn":  stats.RateIn,
				"maxBytesPerSecond": checker.maxBytesPerSecond,
			}).Warn("would ban peer due to high bandwidth usage")
			// There are possibly multiple connections to each peer. We ban the IP
			// address associated with each connection.
			for _, conn := range checker.node.host.Network().ConnsToPeer(remotePeerID) {
				// TODO(albrow): We don't actually ban for now due to an apparent bug in
				// libp2p's BandwidthCounter. Uncomment this once the issue is resolved.
				// See: https://github.com/libp2p/go-libp2p-core/issues/65
				//
				// if err := checker.node.BanIP(conn.RemoteMultiaddr()); err != nil {
				// 	if err == errProtectedIP {
				// 		continue
				// 	}
				// 	log.WithFields(log.Fields{
				// 		"remotePeerID":    remotePeerID.String(),
				// 		"remoteMultiaddr": conn.RemoteMultiaddr().String(),
				// 		"error":           err.Error(),
				// 	}).Error("could not ban peer")
				// }
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
		}
	}
}
