package banner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/albrow/stringset"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	filter "github.com/libp2p/go-maddr-filter"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	// logBandwidthUsageInterval is how often to log bandwidth usage data.
	logBandwidthUsageInterval = 5 * time.Minute
	// violationsCacheSize is the size of the cache (number of entries) used for
	// tracking bandwidth violations over time.
	violationsCacheSize = 1000
	// violationsBeforeBan is the number of times a peer is allowed to violate the
	// bandwidth limits before being banned.
	violationsBeforeBan = 4
	// violationsTTL is the TTL for bandwidth violations. If a peer does not have
	// any violations during this timespan, their violation count will be reset.
	violationsTTL = 6 * time.Hour
)

var ErrProtectedIP = errors.New("cannot ban protected IP address")

type Banner struct {
	config          Config
	protectedIPsMut sync.RWMutex
	protectedIPs    stringset.Set
	violations      *violationsTracker
}

type Config struct {
	Host                   host.Host
	Filters                *filter.Filters
	BandwidthCounter       *metrics.BandwidthCounter
	MaxBytesPerSecond      float64
	LogBandwidthUsageStats bool
}

func New(ctx context.Context, config Config) *Banner {
	banner := &Banner{
		config:       config,
		protectedIPs: stringset.New(),
		violations:   newViolationsTracker(ctx),
	}
	if config.LogBandwidthUsageStats {
		go banner.continuouslyLogBandwidthUsage(ctx)
	}
	return banner
}

// ProtectIP permanently adds the IP address of the given Multiaddr to a
// list of protected IP addresses. Protected IPs can never be banned and will
// not be added to the blacklist. If the IP address is already on the blacklist,
// it will be removed.
func (banner *Banner) ProtectIP(maddr ma.Multiaddr) error {
	banner.protectedIPsMut.Lock()
	defer banner.protectedIPsMut.Unlock()
	ipNet, err := ipNetFromMaddr(maddr)
	if err != nil {
		return err
	}
	banner.unbanIPNet(ipNet)
	banner.protectedIPs.Add(ipNet.IP.String())
	return nil
}

// BanIP adds the IP address of the given Multiaddr to the blacklist. The
// node will no longer dial or accept connections from this IP address. However,
// if the IP address is protected, calling BanIP will not ban the IP address and
// will instead return errProtectedIP. BanIP does not automatically disconnect
// from the given multiaddress if there is currently an open connection.
func (banner *Banner) BanIP(maddr ma.Multiaddr) error {
	ipNet, err := ipNetFromMaddr(maddr)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"maddr": maddr.String(),
		}).Error("could not get IP address from multiaddress")
		return err
	}
	banner.protectedIPsMut.RLock()
	defer banner.protectedIPsMut.RUnlock()
	if banner.protectedIPs.Contains(ipNet.IP.String()) {
		// IP address is protected. no-op.
		return ErrProtectedIP
	}
	banner.config.Filters.AddFilter(ipNet, filter.ActionDeny)
	return nil
}

// UnbanIP removes the IP address of the given Multiaddr from the blacklist. If
// the IP address is not currently on the blacklist this is a no-op.
func (banner *Banner) UnbanIP(maddr ma.Multiaddr) error {
	ipNet, err := ipNetFromMaddr(maddr)
	if err != nil {
		return err
	}
	banner.unbanIPNet(ipNet)
	return nil
}

func (banner *Banner) IsAddrBanned(maddr ma.Multiaddr) bool {
	return banner.config.Filters.AddrBlocked(maddr)
}

func (banner *Banner) SetMaxBytesPerSecond(limit float64) {
	banner.config.MaxBytesPerSecond = limit
}

func (banner *Banner) unbanIPNet(ipNet net.IPNet) {
	// There is no guarantee in the public API of the filters package that would
	// prevent multiple filters being added for the same IPNet (though it
	// shouldn't happen in practice). We use a for loop here to make sure we
	// remove all possible filters. RemoveLiteral returns false if no filter was
	// removed.
	for banner.config.Filters.RemoveLiteral(ipNet) {
	}
}

func (banner *Banner) continuouslyLogBandwidthUsage(ctx context.Context) {
	logTicker := time.Tick(logBandwidthUsageInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-logTicker:
			// Log the bandwidth used by each peer.
			for _, remotePeerID := range banner.config.Host.Network().Peers() {
				stats := banner.config.BandwidthCounter.GetBandwidthForPeer(remotePeerID)
				log.WithFields(log.Fields{
					"remotePeerID":      remotePeerID.String(),
					"bytesPerSecondIn":  stats.RateIn,
					"totalBytesIn":      stats.TotalIn,
					"bytesPerSecondOut": stats.RateOut,
					"totalBytesOut":     stats.TotalOut,
				}).Debug("bandwidth used by peer")
			}

			// Log the bandwidth used by each protocol.
			for protocolID, stats := range banner.config.BandwidthCounter.GetBandwidthByProtocol() {
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

// CheckBandwidthUsage checks the amount of data sent by each connected peer and
// bans (via IP address) any peers which have exceeded the bandwidth limit.
func (banner *Banner) CheckBandwidthUsage() {
	for _, remotePeerID := range banner.config.Host.Network().Peers() {
		stats := banner.config.BandwidthCounter.GetBandwidthForPeer(remotePeerID)
		// If the peer is sending data at a higher rate than is allowed, ban
		// them.
		if stats.RateIn > banner.config.MaxBytesPerSecond {
			numViolations := banner.violations.add(remotePeerID)

			// Check if the number of violations exceeds violationsBeforeBan.
			if numViolations >= violationsBeforeBan {
				log.WithFields(log.Fields{
					"remotePeerID":      remotePeerID.String(),
					"bytesPerSecondIn":  stats.RateIn,
					"maxBytesPerSecond": banner.config.MaxBytesPerSecond,
					"numViolations":     numViolations,
				}).Warn("banning peer due to high bandwidth usage")
				// There are possibly multiple connections to each peer. We ban the IP
				// address associated with each connection.
				for _, conn := range banner.config.Host.Network().ConnsToPeer(remotePeerID) {
					if err := banner.BanIP(conn.RemoteMultiaddr()); err != nil {
						if err == ErrProtectedIP {
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
						"maxBytesPerSecond": banner.config.MaxBytesPerSecond,
					}).Error("banning IP/multiaddress due to high bandwidth usage")
				}
				// Banning the IP doesn't close the connection, so we do that
				// separately. ClosePeer closes all connections to the given peer.
				_ = banner.config.Host.Network().ClosePeer(remotePeerID)
			} else {
				// Log that high bandwidth usage occurred but don't yet ban the peer.
				log.WithFields(log.Fields{
					"remotePeerID":      remotePeerID.String(),
					"bytesPerSecondIn":  stats.RateIn,
					"maxBytesPerSecond": banner.config.MaxBytesPerSecond,
					"numViolations":     numViolations,
				}).Warn("detected high bandwidth usage")
			}
		}
	}
}

func ipNetFromMaddr(maddr ma.Multiaddr) (ipNet net.IPNet, err error) {
	ip, err := ipFromMaddr(maddr)
	if err != nil {
		return net.IPNet{}, err
	}
	mask := getAllMaskForIP(ip)
	return net.IPNet{
		IP:   ip,
		Mask: mask,
	}, nil
}

func ipFromMaddr(maddr ma.Multiaddr) (net.IP, error) {
	var (
		ip    net.IP
		found bool
	)

	ma.ForEach(maddr, func(c ma.Component) bool {
		switch c.Protocol().Code {
		case ma.P_IP6ZONE:
			return true
		case ma.P_IP6, ma.P_IP4:
			found = true
			ip = net.IP(c.RawValue())
			return false
		default:
			return false
		}
	})

	if !found {
		return net.IP{}, fmt.Errorf("could not parse IP address from multiaddress: %s", maddr)
	}
	return ip, nil
}

var (
	ipv4AllMask = net.IPMask{255, 255, 255, 255}
	ipv6AllMask = net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
)

// getAllMaskForIP returns an IPMask that will match all IP addresses. The size
// of the mask depends on whether the given IP address is an IPv4 or an IPv6
// address.
func getAllMaskForIP(ip net.IP) net.IPMask {
	if ip.To4() != nil {
		// This is an ipv4 address. Return 4 byte mask.
		return ipv4AllMask
	}
	// Assume ipv6 address. Return 16 byte mask.
	return ipv6AllMask
}
