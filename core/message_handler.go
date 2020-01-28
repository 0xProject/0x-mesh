package core

import (
	"context"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/encoding"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// Ensure that App implements p2p.MessageHandler.
var _ p2p.MessageHandler = &App{}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (app *App) HandleMessages(ctx context.Context, messages []*p2p.Message) error {
	// First we validate the messages and decode them into orders.
	orders := []*zeroex.SignedOrder{}
	orderHashToMessage := map[common.Hash]*p2p.Message{}

	for _, msg := range messages {
		if err := validateMessageSize(msg); err != nil {
			log.WithFields(map[string]interface{}{
				"error":                 err,
				"from":                  msg.From,
				"maxMessageSizeInBytes": constants.MaxMessageSizeInBytes,
				"actualSizeInBytes":     len(msg.Data),
			}).Trace("received message that exceeds maximum size")
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
			continue
		}

		order, err := encoding.RawMessageToOrder(msg.Data)
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
	validationResults, err := app.orderWatcher.ValidateAndStoreValidOrders(ctx, orders, false, app.chainID)
	if err != nil {
		return err
	}

	// Store any valid orders and update the peer scores.
	for _, acceptedOrderInfo := range validationResults.Accepted {
		// If the order isn't new, we don't log it's receipt or adjust peer scores
		if !acceptedOrderInfo.IsNew {
			continue
		}
		msg := orderHashToMessage[acceptedOrderInfo.OrderHash]
		// If we've reached this point, the message is valid, we were able to
		// decode it into an order and check that this order is valid. Update
		// peer scores accordingly.
		log.WithFields(map[string]interface{}{
			"orderHash": acceptedOrderInfo.OrderHash.Hex(),
			"from":      msg.From.String(),
		}).Info("received new valid order from peer")
		log.WithFields(map[string]interface{}{
			"order":     acceptedOrderInfo.SignedOrder,
			"orderHash": acceptedOrderInfo.OrderHash.Hex(),
			"from":      msg.From.String(),
		}).Trace("all fields for new valid order received from peer")
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
		case ordervalidator.ROInternalError, ordervalidator.ROEthRPCRequestFailed, ordervalidator.ROCoordinatorRequestFailed, ordervalidator.RODatabaseFullOfOrders:
			// Don't incur a negative score for these status types (it might not be
			// their fault).
		default:
			// For other status types, we need to update the peer's score
			app.handlePeerScoreEvent(msg.From, psInvalidMessage)
		}
	}
	return nil
}

func validateMessageSize(message *p2p.Message) error {
	if len(message.Data) > constants.MaxMessageSizeInBytes {
		return constants.ErrMaxMessageSize
	}
	return nil
}
