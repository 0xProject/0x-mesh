// Package ordersync contains the ordersync protocol, which is
// used for sharing existing orders between two peers, typically
// during initialization. The protocol consists of a requester
// (the peer requesting orders) and a provider (the peer providing
// them).
package ordersync_v4

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/albrow/stringset"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jpillora/backoff"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	// ID is the ID for the ordersync protocol.
	ID = protocol.ID("/0x-mesh/order-sync/version/4")
	// requestResponseTimeout is the amount of time to wait for a response/request
	// from the other side of the connection. It is used for both waiting for a request
	// on a newly opened stream and waiting for a response after sending a request.
	requestResponseTimeout = 30 * time.Second
	// maxRequestsPerSecond is the maximum number of ordersync requests to allow per
	// second. If this limit is exceeded, requests will be dropped.
	maxRequestsPerSecond = 30
	// requestsBurst is the maximum number of requests to allow at once.
	requestsBurst = 10
	// ordersyncJitterAmount is the amount of random jitter to add to the delay before
	// each run of ordersync in PeriodicallyGetOrders. It is bound by:
	//
	//    approxDelay * (1 - jitter) <= actualDelay < approxDelay * (1 + jitter)
	//
	ordersyncJitterAmount = 0.1
)

var (
	// ErrNoOrders is returned whenever the orders we are looking for cannot be
	// found anywhere on the network. This can mean that we aren't connected to any
	// peers on the same topic, that there are no orders for the topic throughout
	// the entire network, or that there are peers that have the orders we're
	// looking for, but they are refusing to give them to us.
	ErrNoOrders = errors.New("no orders where received from any known peers")
	// ErrNoOrdersFromPeer is returned when a peer returns no orders during ordersync.
	ErrNoOrdersFromPeer = errors.New("no orders received from peer")
)

type App interface {
	Node() *p2p.Node
	OrderWatcher() *orderwatch.Watcher
	DB() *db.DB
}

// Request is a V4 ordersync request
type Request struct {
	MinOrderHash common.Hash `json:"minOrderHash"`
}

// Response is a V4 ordersync request
type Response struct {
	Orders []*zeroex.SignedOrderV4 `json:"orders"`
}

// Service is the main entrypoint for running the ordersync protocol. It handles
// responding to and sending ordersync requests.
type Service struct {
	ctx context.Context
	app App
	// requestRateLimiter is a rate limiter for incoming ordersync requests. It's
	// shared between all peers.
	requestRateLimiter *rate.Limiter
	perPage            int
}

// New creates and returns a new ordersync service, which is used for both
// requesting orders from other peers and providing orders to peers who request
// them. New expects an array of subprotocols which the service will support, in the
// order of preference. The service will automatically pick the most preferred protocol
// that is supported by both peers for each request/response.
func New(ctx context.Context, app App) *Service {
	s := &Service{
		ctx:                ctx,
		app:                app,
		requestRateLimiter: rate.NewLimiter(maxRequestsPerSecond, requestsBurst),
	}
	s.app.Node().SetStreamHandler(ID, s.HandleStream)
	return s
}

// HandleStream is a stream handler that is used to handle incoming ordersync requests.
func (s *Service) HandleStream(stream network.Stream) {
	if !s.requestRateLimiter.Allow() {
		// Pre-emptively close the stream if we can't accept anymore requests.
		log.WithFields(log.Fields{
			"requester": stream.Conn().RemotePeer().Pretty(),
		}).Warn("closing ordersync stream because rate limiter is backed up")
		_ = stream.Reset()
		return
	}
	log.WithFields(log.Fields{
		"requester": stream.Conn().RemotePeer().Pretty(),
	}).Trace("handling ordersync stream")
	defer func() {
		_ = stream.Close()
	}()
	requesterID := stream.Conn().RemotePeer()

	for {
		if err := s.requestRateLimiter.Wait(s.ctx); err != nil {
			log.WithFields(log.Fields{
				"requester": stream.Conn().RemotePeer().Pretty(),
			}).Warn("ordersync rate limiter returned error")
			return
		}
		request, err := waitForRequest(s.ctx, stream)
		if err != nil {
			log.WithError(err).Warn("waitForRequest returned error")
			return
		}
		log.WithFields(log.Fields{
			"requester": stream.Conn().RemotePeer().Pretty(),
		}).Trace("received ordersync V4 request")
		response := s.handleRequest(request, requesterID)
		if response == nil {
			return
		}
		if err := json.NewEncoder(stream).Encode(response); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": requesterID.Pretty(),
			}).Warn("could not encode ordersync V4 response")
			s.handlePeerScoreEvent(requesterID, psUnexpectedDisconnect)
			return
		}

		// Stop stream if exchange is complete
		if len(response.Orders) == 0 {
			return
		}
	}
}

// GetOrders iterates through every peer the node is currently connected to
// and attempts to perform the ordersync protocol. It keeps trying until
// ordersync has been completed with minPeers, using an exponential backoff
// strategy between retries.
func (s *Service) GetOrders(ctx context.Context, minPeers int) error {
	successfullySyncedPeers := stringset.New()

	// retryBackoff defines how long to wait before trying again if we didn't get
	// orders from enough peers during the ordersync process.
	retryBackoff := &backoff.Backoff{
		Min:    250 * time.Millisecond, // First back-off length
		Max:    1 * time.Minute,        // Longest back-off length
		Factor: 2,                      // Factor to multiple each successive back-off
	}

	// nextRequestForPeer tracks the last meaningful "next request" that was
	// provided by a peer during ordersync. This allows us to pick up where
	// we left off if a peer disconnects rather than starting to ordersync
	// from the beginning of the peer's database.
	nextRequestForPeer := map[peer.ID]*Request{}
	for len(successfullySyncedPeers) < minPeers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// NOTE(jalextowle): m, wg, and semaphore are used to synchronize
		// requests to get orders from other peers during ordersync. m is
		// used to guard the successfullySyncedPeers stringset from concurrent
		// access. wg is used to ensure that all of the request for orders
		// from other peers has ended before asking from new peers. Finally,
		// semaphore is used to ensure that there are only ever minPeers
		// requests being made at a given time.
		m := &sync.RWMutex{}
		wg := &sync.WaitGroup{}
		semaphore := make(chan struct{}, minPeers)

		currentNeighbors := s.app.Node().Neighbors()
		shufflePeers(currentNeighbors)
		innerCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		for _, peerID := range currentNeighbors {
			// The loop will only advance when a new element can be
			// added to the semaphore. This ensures that no more than
			// minPeers goroutines will be active at a given time
			// because the channel only has a capacity of minPeers.
			select {
			case <-innerCtx.Done():
				break
			case semaphore <- struct{}{}:
			}

			m.RLock()
			successfullySyncedPeerLength := len(successfullySyncedPeers)
			successfullySynced := successfullySyncedPeers.Contains(peerID.Pretty())
			nextRequest := nextRequestForPeer[peerID]
			m.RUnlock()
			if successfullySyncedPeerLength >= minPeers {
				return nil
			}
			if successfullySynced {
				continue
			}

			log.WithFields(log.Fields{
				"provider": peerID.Pretty(),
			}).Trace("requesting orders from neighbor via ordersync")

			wg.Add(1)
			go func(id peer.ID) {
				defer func() {
					wg.Done()
					<-semaphore
				}()
				if nextFirstRequest, err := s.getOrdersFromPeer(innerCtx, id, nextRequest); err != nil {
					log.WithFields(log.Fields{
						"error":    err.Error(),
						"provider": id.Pretty(),
					}).Debug("could not get orders from peer via ordersync")
					m.Lock()
					if nextFirstRequest != nil {
						nextRequestForPeer[id] = nextFirstRequest
					}
					m.Unlock()
				} else {
					log.WithFields(log.Fields{
						"provider": id.Pretty(),
					}).Trace("successfully got orders from peer via ordersync")
					m.Lock()
					successfullySyncedPeers.Add(id.Pretty())
					delete(nextRequestForPeer, id)
					m.Unlock()
				}
			}(peerID)
		}

		wg.Wait()
		cancel()

		m.RLock()
		successfullySyncedPeerLength := len(successfullySyncedPeers)
		m.RUnlock()

		if successfullySyncedPeerLength < minPeers {
			delayBeforeNextRetry := retryBackoff.Duration()
			log.WithFields(log.Fields{
				"delayBeforeNextRetry":    delayBeforeNextRetry.String(),
				"minPeers":                minPeers,
				"successfullySyncedPeers": successfullySyncedPeerLength,
			}).Debug("ordersync could not get orders from enough peers (trying again soon)")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delayBeforeNextRetry):
				continue
			}
		}
	}
	log.WithFields(log.Fields{
		"minPeers":                minPeers,
		"successfullySyncedPeers": len(successfullySyncedPeers),
	}).Info("completed a round of ordersync")
	return nil
}

// PeriodicallyGetOrders periodically calls GetOrders. It waits a minimum of
// approxDelay (with some random jitter) between each call. It will block until
// there is a critical error or the given context is canceled.
func (s *Service) PeriodicallyGetOrders(ctx context.Context, minPeers int, approxDelay time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.GetOrders(ctx, minPeers); err != nil {
			return err
		}

		// Note(albrow): The random jitter here helps smooth out the frequency of ordersync
		// requests and helps prevent a situation where a large number of nodes are requesting
		// orders at the same time.
		delay := calculateDelayWithJitter(approxDelay, ordersyncJitterAmount)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func calculateDelayWithJitter(approxDelay time.Duration, jitterAmount float64) time.Duration {
	jitterBounds := int(float64(approxDelay) * jitterAmount * 2)
	delta := rand.Intn(jitterBounds) - jitterBounds/2
	return approxDelay + time.Duration(delta)
}

// shufflePeers randomizes the order of the given list of peers.
func shufflePeers(peers []peer.ID) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })
}

func waitForRequest(parentCtx context.Context, stream network.Stream) (*Request, error) {
	ctx, cancel := context.WithTimeout(parentCtx, requestResponseTimeout)
	defer cancel()
	reqChan := make(chan *Request, 1)
	errChan := make(chan error, 1)
	go func() {
		var req Request
		if err := json.NewDecoder(stream).Decode(&req); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": stream.Conn().RemotePeer().Pretty(),
			}).Warn("could not decode ordersync request")
			errChan <- err
			return
		}
		reqChan <- &req
	}()

	select {
	case <-ctx.Done():
		log.WithFields(log.Fields{
			"error":     ctx.Err(),
			"requester": stream.Conn().RemotePeer().Pretty(),
		}).Warn("timed out waiting for ordersync request")
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case req := <-reqChan:
		return req, nil
	}
}

func waitForResponse(parentCtx context.Context, stream network.Stream) (*Response, error) {
	ctx, cancel := context.WithTimeout(parentCtx, requestResponseTimeout)
	defer cancel()
	resChan := make(chan *Response, 1)
	errChan := make(chan error, 1)
	go func() {
		var res Response
		if err := json.NewDecoder(stream).Decode(&res); err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"provider": stream.Conn().RemotePeer().Pretty(),
			}).Warn("could not decode ordersync response")
			errChan <- err
			return
		}
		resChan <- &res
	}()

	select {
	case <-ctx.Done():
		log.WithFields(log.Fields{
			"error":    ctx.Err(),
			"provider": stream.Conn().RemotePeer().Pretty(),
		}).Warn("timed out waiting for ordersync response")
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case rawRes := <-resChan:
		return rawRes, nil
	}
}

func (s *Service) handleRequest(request *Request, requesterID peer.ID) *Response {
	// Early exit if channel closed?
	select {
	case <-s.ctx.Done():
		log.WithError(s.ctx.Err()).Warn("handleRequest v4 error")
		return nil
	default:
	}

	// Get the orders for this page.
	ordersResp, err := s.GetOrdersV4(s.perPage, request.MinOrderHash)
	if err != nil {
		log.WithError(err).Warn("handleRequest v4 error")
		return nil
	}
	orders := []*zeroex.SignedOrderV4{}
	for _, orderInfo := range ordersResp.OrdersInfos {
		orders = append(orders, orderInfo.SignedOrderV4)
	}

	s.handlePeerScoreEvent(requesterID, psValidMessage)
	return &Response{
		Orders: orders,
	}
}

// Returns the next request if any, or nil, the number of received orders or err.
func (s *Service) handleOrderSyncResponse(res *Response, peer peer.ID) (*Request, int, error) {
	validationResults, err := s.app.orderWatcher.ValidateAndStoreValidOrdersV4(s.ctx, res.Orders, s.app.chainID, false, &types.AddOrdersOpts{})
	if err != nil {
		return nil, len(res.Orders), err
	}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		if acceptedOrderInfo.IsNew {
			log.WithFields(map[string]interface{}{
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      peer.Pretty(),
				"protocol":  "ordersync",
			}).Info("received new valid order from peer")
			log.WithFields(map[string]interface{}{
				"order":     acceptedOrderInfo.SignedOrderV4,
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      peer.Pretty(),
				"protocol":  "ordersync",
			}).Trace("all fields for new valid order received from peer")
		}
	}

	// Calculate the next min order hash to send in our next request.
	// This is equal to the maximum order hash we have received so far.
	if len(res.Orders) > 0 {
		hash, err := res.Orders[len(res.Orders)-1].ComputeOrderHash()
		if err != nil {
			return nil, len(res.Orders), err
		}
		return &Request{
			MinOrderHash: hash,
		}, len(res.Orders), nil
	} else {
		return nil, len(res.Orders), nil
	}
}

func (s *Service) getOrdersFromPeer(ctx context.Context, providerID peer.ID, nextReq *Request) (*Request, error) {
	stream, err := s.node.NewStream(ctx, providerID, ID)
	if err != nil {
		s.handlePeerScoreEvent(providerID, psUnexpectedDisconnect)
		return nil, err
	}
	defer func() {
		_ = stream.Close()
	}()

	totalValidOrders := 0
	for {
		select {
		case <-ctx.Done():
			return nextReq, ctx.Err()
		default:
		}

		// Create initial request if not provided one
		if nextReq == nil {
			nextReq = &Request{
				MinOrderHash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			}
		}

		// Send the request JSON encoded
		if err := json.NewEncoder(stream).Encode(nextReq); err != nil {
			s.handlePeerScoreEvent(providerID, psUnexpectedDisconnect)
			return nil, err
		}

		// Wait for response
		response, err := waitForResponse(ctx, stream)
		if err != nil {
			return nil, err
		}
		s.handlePeerScoreEvent(providerID, psValidMessage)

		// Handle response
		req, numValidOrders, err := s.handleOrderSyncResponse(response, stream.Conn().RemotePeer())
		if err != nil {
			if totalValidOrders == 0 {
				return nil, err
			} else {
				// Likely connection failure, retry where we left
				return nextReq, err
			}
		}
		s.handlePeerScoreEvent(providerID, receivedOrders)

		totalValidOrders += numValidOrders
		if req == nil { // Indicates sync complete
			if totalValidOrders == 0 {
				return nil, ErrNoOrdersFromPeer
			} else {
				return nil, nil
			}
		}
		nextReq = req
	}
}

// ErrPerPageZero is the error returned when a GetOrders request specifies perPage to 0
type ErrPerPageZero struct{}

func (e ErrPerPageZero) Error() string {
	return "perPage cannot be zero"
}

func (s *Service) GetOrdersV4(perPage int, minOrderHash common.Hash) ([]*zeroex.SignedOrderV4, error) {
	if perPage <= 0 {
		return nil, ErrPerPageZero{}
	}
	ordersWithMeta, err := s.app.DB().FindOrdersV4(&db.OrderQueryV4{
		Filters: []db.OrderFilterV4{
			{
				Field: db.OV4FIsRemoved,
				Kind:  db.Equal,
				Value: false,
			},
			{
				Field: db.OV4FHash,
				Kind:  db.Greater,
				Value: minOrderHash,
			},
		},
		Sort: []db.OrderSortV4{
			{
				Field:     db.OV4FHash,
				Direction: db.Ascending,
			},
		},
		Limit: uint(perPage),
	})
	if err != nil {
		return nil, err
	}
	var orders []*zeroex.SignedOrderV4
	for _, order := range ordersWithMeta {
		orders = append(orders, &order.SignedOrderV4())
	}

	getOrdersResponse := &types.GetOrdersResponse{
		Timestamp:   time.Now(),
		OrdersInfos: ordersInfos,
	}

	return getOrdersResponse, nil
}
