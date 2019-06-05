package rpc

import (
	"context"

	"github.com/0xProject/0x-mesh/telemetry"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/rpc"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

// rpcService is an /ethereum/go-ethereum/rpc compatible service.
type rpcService struct {
	rpcHandler RPCHandler
}

// RPCHandler is used to respond to incoming requests from the client.
type RPCHandler interface {
	// AddOrders is called when the client sends an AddOrders request.
	AddOrders(orders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error)
	// AddPeer is called when the client sends an AddPeer request.
	AddPeer(peerInfo peerstore.PeerInfo) error
	// SubscribeToOrders is called when a client sends a Subscribe to orderStream request
	SubscribeToOrders(ctx context.Context) (*rpc.Subscription, error)
}

// Orders calls rpcHandler.Orders and returns the rpc subscription.
func (s *rpcService) Orders(ctx context.Context) (*rpc.Subscription, error) {
	return s.rpcHandler.SubscribeToOrders(ctx)
}

// AddOrders calls rpcHandler.AddOrders and returns the validation results.
func (s *rpcService) AddOrders(orders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error) {
	telemetry.AddOrdersRequests.Inc()
	return s.rpcHandler.AddOrders(orders)
}

// AddPeer builds PeerInfo out of the given peer ID and multiaddresses and
// calls rpcHandler.AddPeer. If there is an error, it returns it.
func (s *rpcService) AddPeer(peerID string, multiaddrs []string) error {
	telemetry.AddPeerRequests.Inc()
	// Parse peer ID.
	parsedPeerID, err := peer.IDB58Decode(peerID)
	if err != nil {
		return err
	}
	peerInfo := peerstore.PeerInfo{
		ID: parsedPeerID,
	}

	// Parse each given multiaddress.
	parsedMultiaddrs := make([]ma.Multiaddr, len(multiaddrs))
	for i, addr := range multiaddrs {
		parsed, err := ma.NewMultiaddr(addr)
		if err != nil {
			return err
		}
		parsedMultiaddrs[i] = parsed
	}
	peerInfo.Addrs = parsedMultiaddrs

	return s.rpcHandler.AddPeer(peerInfo)
}
