package rpc

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/rpc"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

// minHeartbeatInterval specifies the interval at which to emit heartbeat events to a subscriber
var minHeartbeatInterval = 5 * time.Second

// rpcService is an /ethereum/go-ethereum/rpc compatible service.
type rpcService struct {
	rpcHandler RPCHandler
}

// RPCHandler is used to respond to incoming requests from the client.
type RPCHandler interface {
	// AddOrders is called when the client sends an AddOrders request.
	AddOrders(signedOrdersRaw []*json.RawMessage) (*zeroex.ValidationResults, error)
	// GetOrders is called when the clients sends a GetOrders request
	GetOrders(page, perPage int, snapshotID string) (*GetOrdersResponse, error)
	// GetStats is called when the client sends an GetStats request.
	GetStats() (*GetStatsResponse, error)
	// SubscribeToOrders is called when a client sends a Subscribe to `orders` request
	SubscribeToOrders(ctx context.Context) (*rpc.Subscription, error)
}

// Orders calls rpcHandler.SubscribeToOrders and returns the rpc subscription.
func (s *rpcService) Orders(ctx context.Context) (*rpc.Subscription, error) {
	return s.rpcHandler.SubscribeToOrders(ctx)
}

// Heartbeat calls rpcHandler.SubscribeToHeartbeat and returns the rpc subscription.
func (s *rpcService) Heartbeat(ctx context.Context) (*rpc.Subscription, error) {
	log.Debug("received heartbeat subscription request via RPC")
	subscription, err := SetupHeartbeat(ctx)
	if err != nil {
		log.WithField("error", err.Error()).Error("internal error in `mesh_subscribe` to `heartbeat` RPC call")
		return nil, constants.ErrInternal
	}
	return subscription, nil
}

// SetupHeartbeat sets up the heartbeat for a subscription
func SetupHeartbeat(ctx context.Context) (*ethRpc.Subscription, error) {
	notifier, supported := ethRpc.NotifierFromContext(ctx)
	if !supported {
		return &ethRpc.Subscription{}, ethRpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		for {
			select {
			case err := <-rpcSub.Err():
				log.WithField("err", err).Error("rpcSub returned an error")
				return
			case <-notifier.Closed():
				return
			default:
				// Continue
			}

			start := time.Now()

			err := notifier.Notify(rpcSub.ID, "tick")
			if err != nil {
				// TODO(fabio): The current implementation of `notifier.Notify` returns a
				// `write: broken pipe` error when it is called _after_ the client has
				// disconnected but before the corresponding error is received on the
				// `rpcSub.Err()` channel. This race-condition is not problematic beyond
				// the unnecessary computation and log spam resulting from it. Once this is
				// fixed upstream, give all logs an `Error` severity.
				logEntry := log.WithFields(map[string]interface{}{
					"error":            err.Error(),
					"subscriptionType": "heartbeat",
				})
				message := "error while calling notifier.Notify"
				// If the network connection disconnects for longer then ~2mins and then comes
				// back up, we've noticed the call to `notifier.Notify` return `i/o timeout`
				// `net.OpError` errors everytime it's called and no values are sent over
				// `rpcSub.Err()` nor `notifier.Closed()`. In order to stop the error from
				// endlessly re-occuring, we unsubscribe and return for encountering this type of
				// error.
				if _, ok := err.(*net.OpError); ok {
					logEntry.Trace(message)
					return
				}
				if strings.Contains(err.Error(), "write: broken pipe") {
					logEntry.Trace(message)
				} else {
					logEntry.Error(message)
				}
			}

			// Wait MinCleanupInterval before emitting the next heartbeat.
			time.Sleep(minHeartbeatInterval - time.Since(start))

		}
	}()

	return rpcSub, nil
}

// AddOrders calls rpcHandler.AddOrders and returns the validation results.
func (s *rpcService) AddOrders(signedOrdersRaw []*json.RawMessage) (*zeroex.ValidationResults, error) {
	return s.rpcHandler.AddOrders(signedOrdersRaw)
}

// GetOrders calls rpcHandler.GetOrders and returns the validation results.
func (s *rpcService) GetOrders(page, perPage int, snapshotID string) (*GetOrdersResponse, error) {
	return s.rpcHandler.GetOrders(page, perPage, snapshotID)
}

// GetStats calls rpcHandler.GetStats. If there is an error, it returns it.
func (s *rpcService) GetStats() (*GetStatsResponse, error) {
	return s.rpcHandler.GetStats()
}
