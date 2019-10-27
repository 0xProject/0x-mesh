package core

import (
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// Ensure that App implements p2p.MessageHandler.
var _ p2p.MessageHandler = &App{}

type OrderSelector struct {
	nextOffset int
	db         *meshdb.MeshDB
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (app *App) GetMessagesToShare(max int) ([][]byte, error) {
	return app.orderSelector.GetMessagesToShare(max)
}

func (orderSelector *OrderSelector) GetMessagesToShare(max int) ([][]byte, error) {
	// For now, we use a round robin strategy to select a set of orders to share.
	// We might return less than max even if there are max or greater orders
	// currently stored.
	// Use a snapshot to make sure state doesn't change between our two queries.
	ordersSnapshot, err := orderSelector.db.Orders.GetSnapshot()
	if err != nil {
		return nil, err
	}
	defer ordersSnapshot.Release()
	notRemovedFilter := orderSelector.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	count, err := ordersSnapshot.NewQuery(notRemovedFilter).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	// Select up to the maximum number of orders starting at the offset that was
	// calculated the last time this was called with `app`.
	offset := min(orderSelector.nextOffset, count)
	var selectedOrders []*meshdb.Order
	if offset < count {
		err = ordersSnapshot.NewQuery(notRemovedFilter).Offset(offset).Max(max).Run(&selectedOrders)
		if err != nil {
			return nil, err
		}
	}

	// If more orders can be shared than were selected, append the maximum amount of
	// unique (in this round) orders that can be added to the selected orders without
	// exceeding the maximum number to share.
	overflow := min(max-len(selectedOrders), offset)
	if overflow > 0 {
		var overflowSelectedOrders []*meshdb.Order
		err = ordersSnapshot.NewQuery(notRemovedFilter).Offset(0).Max(overflow).Run(&overflowSelectedOrders)
		if err != nil {
			return nil, err
		}
		selectedOrders = append(selectedOrders, overflowSelectedOrders...)
		orderSelector.nextOffset = overflow
	} else {
		// Calculate the next offset and wrap back to 0 if the next offset is larger
		// than or equal to count.
		orderSelector.nextOffset += max
		if orderSelector.nextOffset >= count {
			orderSelector.nextOffset = 0
		}
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

func (app *App) HandleMessages(messages []*p2p.Message) error {
	// First we validate the messages and decode them into orders.
	orders := []*zeroex.SignedOrder{}
	orderHashToMessage := map[common.Hash]*p2p.Message{}

	for _, msg := range messages {
		if err := validateMessageSize(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"error":               err,
				"from":                msg.From,
				"maxOrderSizeInBytes": constants.MaxOrderSizeInBytes,
				"actualSizeInBytes":   len(msg.Data),
			}).Trace("received message that exceeds maximum size")
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
			continue
		}

		result, err := app.schemaValidateMeshMessage(msg.Data)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"error": err,
				"from":  msg.From,
			}).Trace("could not schema validate message")
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
			continue
		}
		if !result.Valid() {
			formattedErrors := make([]string, len(result.Errors()))
			for i, resultError := range result.Errors() {
				formattedErrors[i] = resultError.String()
			}
			log.WithFields(map[string]interface{}{
				"errors": formattedErrors,
				"from":   msg.From,
			}).Trace("order schema validation failed for message")
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
			continue
		}

		order, err := decodeOrder(msg.Data)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"error": err,
				"from":  msg.From,
			}).Trace("could not decode received message")
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
			continue
		}
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			return err
		}
		// Validate doesn't guarantee there are no duplicates so we keep track of
		// which orders we've already seen.
		if _, alreadySeen := orderHashToMessage[orderHash]; alreadySeen {
			continue
		}
		orders = append(orders, order)
		orderHashToMessage[orderHash] = msg
		app.handlePeerScoreEvent(msg.From, psValidMessage)
	}

	// Next, we validate the orders.
	validationResults, err := app.validateOrders(orders)
	if err != nil {
		return err
	}

	// Store any valid orders and update the peer scores.
	for _, acceptedOrderInfo := range validationResults.Accepted {
		if !acceptedOrderInfo.IsNew {
			continue
		}
		msg := orderHashToMessage[acceptedOrderInfo.OrderHash]
		// If we've reached this point, the message is valid and we were able to
		// decode it into an order. Append it to the list of orders to validate and
		// update peer scores accordingly.
		log.WithFields(map[string]interface{}{
			"orderHash": acceptedOrderInfo.OrderHash.Hex(),
			"from":      msg.From.String(),
		}).Info("received new valid order from peer")
		log.WithFields(map[string]interface{}{
			"order":     acceptedOrderInfo.SignedOrder,
			"orderHash": acceptedOrderInfo.OrderHash.Hex(),
			"from":      msg.From.String(),
		}).Trace("all fields for new valid order received from peer")
		// Add stores the message in the database.
		if err := app.orderWatcher.Add(acceptedOrderInfo, false); err != nil {
			if err == meshdb.ErrDBFilledWithPinnedOrders {
				// If the database is full of pinned orders, log and then continue.
				log.WithFields(map[string]interface{}{
					"error":     err.Error(),
					"orderHash": acceptedOrderInfo.OrderHash.Hex(),
					"from":      msg.From.String(),
				}).Error("could not store valid order because database is full")
				continue
			}
			// For any other type of error, return it.
			return err
		}
		app.handlePeerScoreEvent(msg.From, psOrderStored)
	}

	// We don't store invalid orders, but in some cases still need to update peer
	// scores.
	for _, rejectedOrderInfo := range validationResults.Rejected {
		msg := orderHashToMessage[rejectedOrderInfo.OrderHash]
		log.WithFields(map[string]interface{}{
			"rejectedOrderInfo": rejectedOrderInfo,
			"from":              msg.From.String(),
		}).Trace("not storing rejected order received from peer")
		switch rejectedOrderInfo.Status {
		case ordervalidator.ROInternalError, ordervalidator.ROEthRPCRequestFailed, ordervalidator.ROCoordinatorRequestFailed:
			// Don't incur a negative score for these status types (it might not be
			// their fault).
		default:
			// For other status types, we need to update the peer's score
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
		}
	}
	return nil
}
