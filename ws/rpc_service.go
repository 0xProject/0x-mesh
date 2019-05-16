package ws

import (
	"github.com/0xProject/0x-mesh/zeroex"
)

// rpcService is an /ethereum/go-ethereum/rpc compatible service.
type rpcService struct {
	rpcHandler RPCHandler
}

// RPCHandler is used to respond to incoming requests from the client.
type RPCHandler interface {
	// AddOrder is called when the client sends an AddOrder request.
	AddOrder(order *zeroex.SignedOrder) error
}

// AddOrder calls rpcHandler.AddOrder and returns the computed order hash.
// TODO(albrow): Add the ability to send multiple orders at once.
func (s *rpcService) AddOrder(order *zeroex.SignedOrder) (orderHashHex string, err error) {
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		return "", err
	}
	if err := s.rpcHandler.AddOrder(order); err != nil {
		return "", err
	}
	return orderHash.Hex(), nil
}
