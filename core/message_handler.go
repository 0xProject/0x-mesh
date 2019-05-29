// +build !js

package core

import (
	"math/rand"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// Ensure that App implements p2p.MessageHandler.
var _ p2p.MessageHandler = &App{}

func (app *App) GetMessagesToShare(max int) ([][]byte, error) {
	// For now, we just select a random set of orders from those we have stored.
	// We might return less than max even if there are max or greater orders
	// currently stored.
	// TODO(albrow): If the db package supported transactions, we could gaurantee
	// that max orders are returned if we have stored at least max orders. As it
	// currently stands, orders could be added or removed in between when we call
	// Count and when we call Run.
	notRemovedFilter := app.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	count, err := app.db.Orders.NewQuery(notRemovedFilter).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	// If count is less than max, we use an offset of 0 to simply return all
	// orders we have stored.
	offset := 0
	if count > max {
		// If count is greater than max, choose an offset such that we always try to
		// return at least max orders.
		offset = rand.Intn(count - max)
	}
	var selectedOrders []*meshdb.Order
	err = app.db.Orders.NewQuery(notRemovedFilter).Offset(offset).Max(max).Run(&selectedOrders)
	if err != nil {
		return nil, err
	}

	log.WithFields(map[string]interface{}{
		"maxNumberToShare":    max,
		"actualNumberToShare": len(selectedOrders),
	}).Trace("preparing to share orders with peers")

	// After we have selected all the orders to share, we need to encode them to
	// the message data format.
	messageData := make([][]byte, len(selectedOrders))
	for i, order := range selectedOrders {
		log.WithFields(map[string]interface{}{
			"order": order,
		}).Trace("selected order to share")
		encoded, err := encodeOrder(order.SignedOrder)
		if err != nil {
			return nil, err
		}
		messageData[i] = encoded
	}
	return messageData, nil
}

func (app *App) ValidateAndStore(messages []*p2p.Message) ([]*p2p.Message, error) {
	orders := []*zeroex.SignedOrder{}
	orderHashToMessage := map[common.Hash]*p2p.Message{}
	for _, msg := range messages {
		if err := validateMessageSize(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"error":             err,
				"from":              msg.From,
				"maxSizeInBytes":    maxSizeInBytes,
				"actualSizeInBytes": len(msg.Data),
			}).Trace("received message that exceeds maximum size")
			// TODO(albrow): Update peer scores via the Connection Manager. This
			// incur a negative score.
			continue
		}
		order, err := decodeOrder(msg.Data)
		if err != nil {
			return nil, err
		}
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			return nil, err
		}
		// Validate doesn't guarantee there are no duplicates so we keep track of
		// which orders we've already seen.
		if _, alreadySeen := orderHashToMessage[orderHash]; alreadySeen {
			continue
		}
		log.WithFields(map[string]interface{}{
			"order":     order,
			"orderHash": orderHash,
			"from":      msg.From.String(),
		}).Trace("received order from peer")
		orders = append(orders, order)
		orderHashToMessage[orderHash] = msg
	}

	// Validate the orders in a single batch.
	validationResults := app.orderValidator.BatchValidate(orders)

	validMessages := []*p2p.Message{}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		msg := orderHashToMessage[acceptedOrderInfo.OrderHash]
		validMessages = append(validMessages, msg)

		alreadyStored, err := app.orderAlreadyStored(acceptedOrderInfo.OrderHash)
		if err != nil {
			return nil, err
		}
		if alreadyStored {
			log.WithFields(map[string]interface{}{
				"acceptedOrderInfo": acceptedOrderInfo,
				"from":              msg.From.String(),
			}).Trace("order received from peer is valid but already stored")
		} else {
			log.WithFields(map[string]interface{}{
				"acceptedOrderInfo": acceptedOrderInfo,
				"from":              msg.From.String(),
			}).Debug("storing valid order received from peer")
			// Watch stores the message in the database.
			if err := app.orderWatcher.Watch(acceptedOrderInfo); err != nil {
				return nil, err
			}
		}
	}
	for _, rejectedOrderInfo := range validationResults.Rejected {
		// TODO(fabio): What should we do with orders that we fail to validate
		// because of a MeshError (e.g., network disruption) while attempting to validate
		// them? Currently we simply drop them, but perhaps we should re-try validation
		// at a later time?
		msg := orderHashToMessage[rejectedOrderInfo.OrderHash]
		log.WithFields(map[string]interface{}{
			"rejectedOrderInfo": rejectedOrderInfo,
			"from":              msg.From.String(),
		}).Warn("not storing rejected order received from peer")
	}
	return validMessages, nil
}
