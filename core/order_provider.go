package core

import (
	"encoding/json"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/p2p/ordersync"
	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-core/peer"
	log "github.com/sirupsen/logrus"
)

// Ensure that our order provider implements the ordersync.Provider interface.
var _ ordersync.Provider = (*orderProvider)(nil)

type orderProvider struct {
	app *App
	db  *meshdb.MeshDB
}

func newOrderProvider(app *App, db *meshdb.MeshDB) *orderProvider {
	return &orderProvider{
		db:  db,
		app: app,
	}
}

func (p *orderProvider) ProvideOrders(topic string, requestingPeer peer.ID) ([]byte, error) {
	// TODO(albrow): Optimize this.
	// For now we simply get all non-removed orders and return those that match
	// the topic.
	log.WithFields(log.Fields{
		"me":        p.app.peerID.Pretty(),
		"requester": requestingPeer.Pretty(),
	}).Trace("inside ProvideOrders")
	notRemovedFilter := p.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	var nonRemovedOrders []*meshdb.Order
	if err := p.db.Orders.NewQuery(notRemovedFilter).Run(&nonRemovedOrders); err != nil {
		return nil, err
	}
	if len(nonRemovedOrders) == 0 {
		var allOrders []*meshdb.Order
		if err := p.db.Orders.FindAll(&allOrders); err != nil {
			return nil, err
		}
		log.WithFields(log.Fields{
			"me":           p.app.peerID.Pretty(),
			"requester":    requestingPeer.Pretty(),
			"numAllOrders": len(allOrders),
		}).Trace("no non-removed orders found")
	}
	filter, err := orderfilter.NewFromTopic(topic)
	if err != nil {
		return nil, err
	}
	filteredOrders := []*zeroex.SignedOrder{}
	for _, order := range nonRemovedOrders {
		matches, err := filter.MatchOrder(order.SignedOrder)
		if err != nil {
			return nil, err
		}
		if matches {
			filteredOrders = append(filteredOrders, order.SignedOrder)
		}
	}
	if len(filteredOrders) == 0 {
		log.WithFields(log.Fields{
			"me":        p.app.peerID.Pretty(),
			"requester": requestingPeer.Pretty(),
		}).Trace("no orders found that pass filter")
		return nil, nil
	}

	log.WithFields(log.Fields{
		"me":        p.app.peerID.Pretty(),
		"requester": requestingPeer.Pretty(),
		"numOrders": len(filteredOrders),
	}).Trace("provided orders to neighbor")
	return json.Marshal(filteredOrders)
}
