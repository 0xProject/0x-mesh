package rpc

import (
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
)

// AddOrdersResponse is the response returned to the `AddOrders` request
type AddOrdersResponse struct {
	Added       []*zeroex.SuccinctOrderInfo
	Invalid     []*zeroex.SuccinctOrderInfo
	FailedToAdd []common.Hash
}
