// +build !js

package core

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

type orderMessage struct {
	MessageType string
	Order       *zeroex.SignedOrder
}

func encodeOrder(order *zeroex.SignedOrder) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
	})
}

func decodeOrder(data []byte) (*zeroex.SignedOrder, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}

// Ensure that App implements p2p.MessageHandler.
var _ p2p.MessageHandler = &App{}

func (app *App) GetMessagesToShare(max int) ([][]byte, error) {
	// For now, we just select a random set of orders from those we have stored.
	// TODO(albrow): This could be made more efficient if the db package supported
	// a `Count` method for counting the number of models in a collection for
	// counting the number of models that satisfy some query and an `Offset` field
	// for skipping some number of models.
	// TODO: This will need to change when we add support for WeijieSub.
	notDeletedFilter := app.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	var allOrders []*meshdb.Order
	if err := app.db.Orders.NewQuery(notDeletedFilter).Run(&allOrders); err != nil {
		return nil, err
	}
	if len(allOrders) == 0 {
		return nil, nil
	}
	start := rand.Intn(len(allOrders))
	end := start + max
	if end > len(allOrders) {
		end = len(allOrders)
	}
	selectedOrders := allOrders[start:end]

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
	orderHashToOrderInfo := app.orderValidator.BatchValidate(orders)

	// Filter out the invalid messages (i.e. messages which correspond to invalid
	// orders).
	validMessages := []*p2p.Message{}
	for orderHash, msg := range orderHashToMessage {
		orderInfo, found := orderHashToOrderInfo[orderHash]
		if !found {
			continue
		}
		if zeroex.IsOrderValid(orderInfo) {
			validMessages = append(validMessages, msg)
			alreadyStored, err := app.orderAlreadyStored(orderInfo.OrderHash)
			if err != nil {
				return nil, err
			}
			if alreadyStored {
				log.WithFields(map[string]interface{}{
					"orderInfo": orderInfo,
					"from":      msg.From.String(),
				}).Trace("order received from peer is valid but already stored")
			} else {
				log.WithFields(map[string]interface{}{
					"orderInfo": orderInfo,
					"from":      msg.From.String(),
				}).Debug("storing valid order received from peer")
				// Watch stores the message in the database.
				if err := app.orderWatcher.Watch(orderInfo); err != nil {
					return nil, err
				}
			}
		} else {
			log.WithFields(map[string]interface{}{
				"orderInfo": orderInfo,
				"from":      msg.From.String(),
			}).Warn("not storing invalid order received from peer")
		}
	}
	return validMessages, nil
}
