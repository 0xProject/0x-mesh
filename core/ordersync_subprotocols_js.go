// +build js, wasm

package core

import (
	"context"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-peer"
)

const batchSize = 100

func filterOrdersForRequest(f *orderfilter.Filter, orderInfos []*types.OrderInfo, filteredOrders []*zeroex.SignedOrder) ([]*zeroex.SignedOrder, error) {
	nextTick()
	for i, orderInfo := range orderInfos {
		if i%batchSize == batchSize-1 {
			nextTick()
		}
		if matches, err := f.MatchOrder(orderInfo.SignedOrder); err != nil {
			return nil, err
		} else if matches {
			filteredOrders = append(filteredOrders, orderInfo.SignedOrder)
		}
	}
	return filteredOrders, nil
}

func filterOrdersForResponse(a *App, f *orderfilter.Filter, providerID peer.ID, orders []*zeroex.SignedOrder) ([]*zeroex.SignedOrder, error) {
	nextTick()
	var filteredOrders []*zeroex.SignedOrder
	for i, order := range orders {
		if i%batchSize == batchSize-1 {
			nextTick()
		}
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

func nextTick() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	jsutil.NextTick(ctx)
	cancel()
}
