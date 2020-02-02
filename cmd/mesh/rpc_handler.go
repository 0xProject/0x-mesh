// +build !js

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	log "github.com/sirupsen/logrus"
)

// orderEventsBufferSize is the buffer size for the orderEvents channel. If
// the buffer is full, any additional events won't be processed.
const orderEventsBufferSize = 8000

type rpcHandler struct {
	app *core.App
	ctx context.Context
}

// instantiateServer instantiates a new RPC server with the rpcHandler.
func instantiateServer(ctx context.Context, app *core.App, rpcAddr string) *rpc.Server {
	// Initialize the JSON RPC WebSocket server (but don't start it yet).
	rpcHandler := &rpcHandler{
		app: app,
		ctx: ctx,
	}
	rpcServer, err := rpc.NewServer(rpcAddr, rpcHandler)
	if err != nil {
		return nil
	}
	go func() {
		// Wait for the server to start listening and select an address.
		for rpcServer.Addr() == nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			time.Sleep(10 * time.Millisecond)
		}
		log.WithField("address", rpcServer.Addr().String()).Info("started RPC server")
	}()
	return rpcServer
}

// GetOrders is called when an RPC client calls GetOrders.
func (handler *rpcHandler) GetOrders(page, perPage int, snapshotID string) (result *types.GetOrdersResponse, err error) {
	log.WithFields(map[string]interface{}{
		"page":       page,
		"perPage":    perPage,
		"snapshotID": snapshotID,
	}).Debug("received GetOrders request via RPC")
	// Catch panics, log stack trace and return RPC error message
	defer func() {
		if r := recover(); r != nil {
			internalErr, ok := r.(error)
			if !ok {
				// If r is not of type error, convert it.
				internalErr = fmt.Errorf("Recovered from non-error: (%T) %v", r, r)
			}
			log.WithFields(log.Fields{
				"error":      internalErr,
				"method":     "GetOrders",
				"stackTrace": string(debug.Stack()),
			}).Error("RPC method handler crashed")
			err = errors.New("method handler crashed in GetOrders RPC call (check logs for stack trace)")
		}
	}()
	getOrdersResponse, err := handler.app.GetOrders(page, perPage, snapshotID)
	if err != nil {
		if _, ok := err.(core.ErrSnapshotNotFound); ok {
			return nil, err
		}
		if _, ok := err.(core.ErrPerPageZero); ok {
			return nil, err
		}
		// We don't want to leak internal error details to the RPC client.
		log.WithField("error", err.Error()).Error("internal error in GetOrders RPC call")
		return nil, constants.ErrInternal
	}
	return getOrdersResponse, nil
}

// AddOrders is called when an RPC client calls AddOrders.
func (handler *rpcHandler) AddOrders(signedOrdersRaw []*json.RawMessage, opts types.AddOrdersOpts) (results *ordervalidator.ValidationResults, err error) {
	log.WithFields(log.Fields{
		"count":  len(signedOrdersRaw),
		"pinned": opts.Pinned,
	}).Info("received AddOrders request via RPC")
	// Catch panics, log stack trace and return RPC error message
	defer func() {
		if r := recover(); r != nil {
			internalErr, ok := r.(error)
			if !ok {
				// If r is not of type error, convert it.
				internalErr = fmt.Errorf("Recovered from non-error: (%T) %v", r, r)
			}
			log.WithFields(log.Fields{
				"error":      internalErr,
				"method":     "AddOrders",
				"stackTrace": string(debug.Stack()),
			}).Error("RPC method handler crashed")
			err = errors.New("method handler crashed in AddOrders RPC call (check logs for stack trace)")
		}
	}()
	validationResults, err := handler.app.AddOrders(handler.ctx, signedOrdersRaw, opts.Pinned)
	if err != nil {
		// We don't want to leak internal error details to the RPC client.
		log.WithField("error", err.Error()).Error("internal error in AddOrders RPC call")
		return nil, constants.ErrInternal
	}
	return validationResults, nil
}

// AddPeer is called when an RPC client calls AddPeer,
func (handler *rpcHandler) AddPeer(peerInfo peerstore.PeerInfo) (err error) {
	log.Debug("received AddPeer request via RPC")
	// Catch panics, log stack trace and return RPC error message
	defer func() {
		if r := recover(); r != nil {
			internalErr, ok := r.(error)
			if !ok {
				// If r is not of type error, convert it.
				internalErr = fmt.Errorf("Recovered from non-error: (%T) %v", r, r)
			}
			log.WithFields(log.Fields{
				"error":      internalErr,
				"method":     "AddPeer",
				"stackTrace": string(debug.Stack()),
			}).Error("RPC method handler crashed")
			err = errors.New("method handler crashed in AddPeer RPC call (check logs for stack trace)")
		}
	}()
	if err := handler.app.AddPeer(peerInfo); err != nil {
		log.WithField("error", err.Error()).Error("internal error in AddPeer RPC call")
		return constants.ErrInternal
	}
	return nil
}

// GetStats is called when an RPC client calls GetStats,
func (handler *rpcHandler) GetStats() (result *types.Stats, err error) {
	log.Debug("received GetStats request via RPC")
	// Catch panics, log stack trace and return RPC error message
	defer func() {
		if r := recover(); r != nil {
			internalErr, ok := r.(error)
			if !ok {
				// If r is not of type error, convert it.
				internalErr = fmt.Errorf("Recovered from non-error: (%T) %v", r, r)
			}
			log.WithFields(log.Fields{
				"error":      internalErr,
				"method":     "GetStats",
				"stackTrace": string(debug.Stack()),
			}).Error("RPC method handler crashed")
			err = errors.New("method handler crashed in GetStats RPC call (check logs for stack trace)")
		}
	}()
	getStatsResponse, err := handler.app.GetStats()
	if err != nil {
		log.WithField("error", err.Error()).Error("internal error in GetStats RPC call")
		return nil, constants.ErrInternal
	}
	return getStatsResponse, nil
}

// SubscribeToOrders is called when an RPC client sends a `mesh_subscribe` request with the `orders` topic parameter
func (handler *rpcHandler) SubscribeToOrders(ctx context.Context) (result *ethrpc.Subscription, err error) {
	log.Debug("received order event subscription request via RPC")
	// Catch panics, log stack trace and return RPC error message
	defer func() {
		if r := recover(); r != nil {
			internalErr, ok := r.(error)
			if !ok {
				// If r is not of type error, convert it.
				internalErr = fmt.Errorf("Recovered from non-error: (%T) %v", r, r)
			}
			log.WithFields(log.Fields{
				"error":      internalErr,
				"method":     "SubscribeToOrders",
				"stackTrace": string(debug.Stack()),
			}).Error("RPC method handler crashed")
			err = errors.New("method handler crashed in SubscribeToOrders RPC call (check logs for stack trace)")
		}
	}()
	subscription, err := SetupOrderStream(ctx, handler.app)
	if err != nil {
		log.WithField("error", err.Error()).Error("internal error in `mesh_subscribe` to `orders` RPC call")
		return nil, constants.ErrInternal
	}
	return subscription, nil
}

// SetupOrderStream sets up the order stream for a subscription
func SetupOrderStream(ctx context.Context, app *core.App) (*ethrpc.Subscription, error) {
	notifier, supported := ethrpc.NotifierFromContext(ctx)
	if !supported {
		return &ethrpc.Subscription{}, ethrpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		orderEventsChan := make(chan []*zeroex.OrderEvent, orderEventsBufferSize)
		orderWatcherSub := app.SubscribeToOrderEvents(orderEventsChan)
		defer orderWatcherSub.Unsubscribe()

		for {
			select {
			case orderEvents := <-orderEventsChan:
				err := notifier.Notify(rpcSub.ID, orderEvents)
				if err != nil {
					// TODO(fabio): The current implementation of `notifier.Notify` returns a
					// `write: broken pipe` error when it is called _after_ the client has
					// disconnected but before the corresponding error is received on the
					// `rpcSub.Err()` channel. This race-condition is not problematic beyond
					// the unnecessary computation and log spam resulting from it. Once this is
					// fixed upstream, give all logs an `Error` severity.
					logEntry := log.WithFields(map[string]interface{}{
						"error":            err.Error(),
						"subscriptionType": "orders",
						"orderEvents":      len(orderEvents),
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
			case err := <-rpcSub.Err():
				if err != nil {
					log.WithField("err", err).Error("rpcSub returned an error")
				} else {
					log.Debug("rpcSub was closed without error")
				}
				return
			case <-notifier.Closed():
				return
			}
		}
	}()

	return rpcSub, nil
}
