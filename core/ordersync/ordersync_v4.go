// Package ordersync contains the ordersync protocol, which is
// used for sharing existing orders between two peers, typically
// during initialization. The protocol consists of a requester
// (the peer requesting orders) and a provider (the peer providing
// them).
package ordersync

import (
	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

type ResponseV4 struct {
	ProviderID peer.ID                 `json:"providerID"`
	Orders     []*zeroex.SignedOrderV4 `json:"orders"`
}
