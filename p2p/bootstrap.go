// +build !js

package p2p

import (
	"context"
	"sync"
	"time"

	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const dhtProtocolID = protocol.ID("/0x-mesh-dht/version/1")

// BootstrapPeers is a list of peers to use for bootstrapping the DHT.
var BootstrapPeers []multiaddr.Multiaddr

func init() {
	for _, s := range []string{
		"/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF",
		"/ip4/18.200.96.60/tcp/60558/ipfs/16Uiu2HAkwsDZk4LzXy2rnWANRsyBjB4fhjnsNeJmjgsBqxPGTL32",
		"/ip4/13.232.193.142/tcp/60558/ipfs/16Uiu2HAkykwoBxwyvoEbaEkuKMeKrmJDPZ2uKFPUKtqd2JbGHUNH",
		"/dns4/bootstrap-0.mesh.0x.org/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF",
		"/dns4/bootstrap-1.mesh.0x.org/tcp/60558/ipfs/16Uiu2HAkwsDZk4LzXy2rnWANRsyBjB4fhjnsNeJmjgsBqxPGTL32",
		"/dns4/bootstrap-2.mesh.0x.org/tcp/60558/ipfs/16Uiu2HAkykwoBxwyvoEbaEkuKMeKrmJDPZ2uKFPUKtqd2JbGHUNH",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		BootstrapPeers = append(BootstrapPeers, ma)
	}
}

func ConnectToBootstrapList(ctx context.Context, host host.Host) error {
	log.WithField("BootstrapPeers", BootstrapPeers).Info("connecting to bootstrap peers")
	connectCtx, cancel := context.WithTimeout(ctx, defaultNetworkTimeout)
	defer cancel()
	wg := sync.WaitGroup{}
	for _, addr := range BootstrapPeers {
		peerInfo, err := peerstore.InfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		if peerInfo.ID == host.ID() {
			// Don't connect to self.
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(connectCtx, *peerInfo); err != nil {
				log.WithFields(map[string]interface{}{
					"error":    err.Error(),
					"peerInfo": peerInfo,
				}).Warn("failed to connect to bootstrap peer")
			}
		}()
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
