package rpc

import (
	"context"
	"encoding/json"

	"github.com/0xProject/0x-mesh/zeroex"
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
	AddOrders(orders []*zeroex.SignedOrder) (zeroex.OrderHashToSuccinctOrderInfo, error)
	// AddPeer is called when the client sends an AddPeer request.
	AddPeer(peerInfo peerstore.PeerInfo) error
}

// AddOrders calls rpcHandler.AddOrders and returns the SuccinctOrderInfo for each order.
func (s *rpcService) AddOrders(orders []*zeroex.SignedOrder) (string, error) {
	orderHashToSuccinctOrderInfo, err := s.rpcHandler.AddOrders(orders)
	if err != nil {
		return "", err
	}
	orderHashToSuccinctOrderInfoBytes, err := json.Marshal(orderHashToSuccinctOrderInfo)
	return string(orderHashToSuccinctOrderInfoBytes), nil
}

// AddOrder builds PeerInfo out of the given peer ID and multiaddresses and
// calls rpcHandler.AddPeer. If there is an error, it returns it.
func (s *rpcService) AddPeer(peerID string, multiaddrs []string) error {
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
