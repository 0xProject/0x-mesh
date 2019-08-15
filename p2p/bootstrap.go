// +build !js

package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const dhtProtocolID = protocol.ID("/0x-mesh-dht/version/1")

// BootstrapPeers is a list of peers to use for bootstrapping the DHT.
var BootstrapPeers []peer.AddrInfo

func init() {
	for _, rawInfo := range []struct {
		addrs  []string
		peerID string
	}{
		{
			addrs: []string{
				"/ip4/3.214.190.67/tcp/60558",
				"/ip4/3.214.190.67/tcp/60559/ws",
				"/dns4/bootstrap-0.mesh.0x.org/tcp/60558",
				"/dns4/bootstrap-0.mesh.0x.org/tcp/60559/ws",
			},
			peerID: "16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF",
		},
		{
			addrs: []string{
				"/ip4/18.200.96.60/tcp/60558",
				"/ip4/18.200.96.60/tcp/60559/ws",
				"/dns4/bootstrap-1.mesh.0x.org/tcp/60558",
				"/dns4/bootstrap-1.mesh.0x.org/tcp/60559/ws",
			},
			peerID: "16Uiu2HAkwsDZk4LzXy2rnWANRsyBjB4fhjnsNeJmjgsBqxPGTL32",
		},
		{
			addrs: []string{
				"/ip4/13.232.193.142/tcp/60558",
				"/ip4/13.232.193.142/tcp/60559/ws",
				"/dns4/bootstrap-2.mesh.0x.org/tcp/60558",
				"/dns4/bootstrap-2.mesh.0x.org/tcp/60559/ws",
			},
			peerID: "16Uiu2HAkykwoBxwyvoEbaEkuKMeKrmJDPZ2uKFPUKtqd2JbGHUNH",
		},

		// These nodes are provided by the libp2p community on a best-effort basis.
		// We're using them as a backup for increased redundancy.
		{
			addrs:  []string{"/ip4/34.201.54.78/tcp/4001"},
			peerID: "12D3KooWHwJDdbx73qiBpSCJfg4RuYyzqnLUwfLBqzn77TSy7kRX",
		},
		{
			addrs:  []string{"/ip4/18.204.221.103/tcp/4001"},
			peerID: "12D3KooWQS6Gsr2kLZvF7DVtoRFtj24aar5jvz88LvJePrawM3EM",
		},
	} {
		peerID, err := peer.IDB58Decode(rawInfo.peerID)
		if err != nil {
			panic(err)
		}
		addrInfo := peer.AddrInfo{
			ID: peerID,
		}
		for _, rawAddr := range rawInfo.addrs {
			addr, err := multiaddr.NewMultiaddr(rawAddr)
			if err != nil {
				panic(err)
			}
			addrInfo.Addrs = append(addrInfo.Addrs, addr)
		}
		BootstrapPeers = append(BootstrapPeers, addrInfo)
	}
}

func ConnectToBootstrapList(ctx context.Context, host host.Host) error {
	log.WithField("BootstrapPeers", BootstrapPeers).Info("connecting to bootstrap peers")
	connectCtx, cancel := context.WithTimeout(ctx, defaultNetworkTimeout)
	defer cancel()
	wg := sync.WaitGroup{}
	for _, peerInfo := range BootstrapPeers {
		if peerInfo.ID == host.ID() {
			// Don't connect to self.
			continue
		}
		wg.Add(1)
		go func(peerInfo peer.AddrInfo) {
			defer wg.Done()
			if err := host.Connect(connectCtx, peerInfo); err != nil {
				log.WithFields(map[string]interface{}{
					"error":    err.Error(),
					"peerInfo": peerInfo,
				}).Warn("failed to connect to bootstrap peer")
			}
		}(peerInfo)
	}
	wg.Wait()

	// It is recommended to wait for 2 seconds after connecting to all the
	// bootstrap peers to give time for the relevant notifees to trigger and the
	// DHT to fully initialize.
	// See: https://github.com/0xProject/0x-mesh/pull/69#discussion_r286849679
	time.Sleep(2 * time.Second)

	return nil
}

// NewDHT returns a new Kademlia DHT instance configured to work with 0x Mesh.
func NewDHT(ctx context.Context, host host.Host) (*dht.IpfsDHT, error) {
	return dht.New(ctx, host, dhtopts.Protocols(dhtProtocolID))
}
