package ws

import (
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
)

// rpcService is an /ethereum/go-ethereum/rpc compatible service.
type rpcService struct {
	orderHandler OrderHandler
}

// OrderHandler is used to respond to incoming requests from the client.
type OrderHandler interface {
	// AddOrder is called when the client sends an AddOrder request.
	AddOrder(order *zeroex.SignedOrder) error
	// RemoveOrder is called when the client sends a RemoveOrder request.
	RemoveOrder(orderHash common.Hash) error
}

// AddOrder calls orderHandler.AddOrder and returns the computed order hash.
func (s *rpcService) AddOrder(order *zeroex.SignedOrder) (orderHashHex string, err error) {
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		return "", err
	}
	if err := s.orderHandler.AddOrder(order); err != nil {
		return "", err
	}
	return orderHash.Hex(), nil
}

// RemoveOrder calls orderHandler.RemoveOrder.
func (s *rpcService) RemoveOrder(orderHashHex string) error {
	orderHash := common.HexToHash(orderHashHex)
	return s.orderHandler.RemoveOrder(orderHash)
}
