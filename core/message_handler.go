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
	// Use a snapshot to make sure state doesn't change between our two queries.
	ordersSnapshot, err := app.db.Orders.GetSnapshot()
	if err != nil {
		return nil, err
	}
	defer ordersSnapshot.Release()
	notRemovedFilter := app.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	count, err := ordersSnapshot.NewQuery(notRemovedFilter).Count()
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
		// If count is greater than max, choose an offset such that we always return
		// at least max orders.
		offset = rand.Intn(count - max)
	}
	var selectedOrders []*meshdb.Order
	err = ordersSnapshot.NewQuery(notRemovedFilter).Offset(offset).Max(max).Run(&selectedOrders)
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

func (app *App) HandleMessages(messages []*p2p.Message) error {
	// First we validate the messages and decode them into orders.
	orders := []*zeroex.SignedOrder{}
	orderHashToMessage := map[common.Hash]*p2p.Message{}

	for _, msg := range messages {
		if err := validateMessageSize(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"error":               err,
				"from":                msg.From,
				"maxOrderSizeInBytes": maxOrderSizeInBytes,
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
		if err := app.orderWatcher.Add(acceptedOrderInfo); err != nil {
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
		case ROInternalError, zeroex.ROEthRPCRequestFailed, zeroex.ROCoordinatorRequestFailed:
			// Don't incur a negative score for these status types (it might not be
			// their fault).
		default:
			// For other status types, we need to update the peer's score
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
		}
	}
	return nil
}
