// +build !js

package core

import (
	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-peer"
)

func filterOrdersForRequest(f *orderfilter.Filter, orderInfos []*types.OrderInfo, filteredOrders []*zeroex.SignedOrder) ([]*zeroex.SignedOrder, error) {
	for _, orderInfo := range orderInfos {
		if matches, err := f.MatchOrder(orderInfo.SignedOrder); err != nil {
			return nil, err
		} else if matches {
			filteredOrders = append(filteredOrders, orderInfo.SignedOrder)
		}
	}
	return filteredOrders, nil
}

func filterOrdersForResponse(a *App, f *orderfilter.Filter, providerID peer.ID, orders []*zeroex.SignedOrder) ([]*zeroex.SignedOrder, error) {
	var filteredOrders []*zeroex.SignedOrder
	for _, order := range orders {
		if matches, err := f.MatchOrder(order); err != nil {
			return nil, err
		} else if matches {
			filteredOrders = append(filteredOrders, order)
		} else if !matches {
			a.handlePeerScoreEvent(providerID, psReceivedOrderDoesNotMatchFilter)
		}
	}
	return filteredOrders, nil
}
