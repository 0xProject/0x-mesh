// Package ordersync contains the ordersync protocol, which is
// used for sharing existing orders between two peers, typically
// during initialization. The protocol consists of a requester
// (the peer requesting orders) and a provider (the peer providing
// them).
package ordersync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/metrics"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/albrow/stringset"
	"github.com/jpillora/backoff"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	// TypeRequest is used to identify a JSON message as an ordersync request.
	TypeRequest = "Request"
	// TypeResponse is used to identify a JSON message as an ordersync response.
	TypeResponse = "Response"
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
	// ErrNoOrderFromPeer is returned when a peer returns no orders during ordersync.
	ErrNoOrdersFromPeer = errors.New("no orders received from peer")
)

// NoMatchingSubprotocolsError is returned whenever two peers attempting to use
// the ordersync protocol cannot agree on a subprotocol to use.
type NoMatchingSubprotocolsError struct {
	Requested []string
	Supported []string
}

func (e NoMatchingSubprotocolsError) Error() string {
	return fmt.Sprintf("could not agree on an ordersync subprotocol (requested: %v, supported: %s)", e.Requested, e.Supported)
}

const (
	// ID is the ID for the ordersync protocol.
	ID = protocol.ID("/0x-mesh/order-sync/version/0")
)

// Request represents a high-level ordersync request. It abstracts away some
// of the details of subprotocol negotiation and encoding/decoding.
type Request struct {
	RequesterID peer.ID     `json:"requesterID"`
	Metadata    interface{} `json:"metadata"`
}

// rawRequest contains all the details we need at the lowest level to encode/decode
// the request and perform subprotocol negoatiation.
type rawRequest struct {
	Type         string          `json:"type"`
	Subprotocols []string        `json:"subprotocols"`
	Metadata     json.RawMessage `json:"metadata"`
}

// Response represents a high-level ordersync response. It abstracts away some
// of the details of subprotocol negotiation and encoding/decoding.
type Response struct {
	ProviderID peer.ID               `json:"providerID"`
	Orders     []*zeroex.SignedOrder `json:"orders"`
	Complete   bool                  `json:"complete"`
	Metadata   interface{}           `json:"metadata"`
}

// rawResponse contains all the details we need at the lowest level to encode/decode
// the response, perform subprotocol negoatiation, and more.
type rawResponse struct {
	Type        string                `json:"type"`
	Subprotocol string                `json:"subprotocol"`
	Orders      []*zeroex.SignedOrder `json:"orders"`
	Complete    bool                  `json:"complete"`
	Metadata    json.RawMessage       `json:"metadata"`
}

// Service is the main entrypoint for running the ordersync protocol. It handles
// responding to and sending ordersync requests.
type Service struct {
	ctx  context.Context
	node *p2p.Node
	// preferredSubprotocols is the list of supported subprotocol IDs in order of preference.
	preferredSubprotocols []string
	subprotocolSet        map[string]Subprotocol
	// requestRateLimiter is a rate limiter for incoming ordersync requests. It's
	// shared between all peers.
	requestRateLimiter *rate.Limiter
}

// Subprotocol is a lower-level protocol which defines the details for the
// request/response metadata. While the ordersync protocol supports sending
// requests and responses in order to synchronize orders between two peers
// in general, a subprotocol defines exactly what those requests and responses
// should look like and how each peer is expected to respond to them.
type Subprotocol interface {
	// Name is the name of the subprotocol. Must be unique.
	Name() string
	// HandleOrderSyncRequest returns a Response based on the given Request. It is the
	// implementation for the "provider" side of the subprotocol.
	HandleOrderSyncRequest(context.Context, *Request) (*Response, error)
	// HandleOrderSyncResponse handles a response (e.g. typically by saving orders to
	// the database), returns the number of valid orders that were received,
	// and, if needed, creates and returns the next request that should be sent.
	// If nextRequest is nil, the ordersync protocol is considered finished.
	// HandleOrderSyncResponse is the implementation for the "requester" side
	// of the subprotocol.
	HandleOrderSyncResponse(context.Context, *Response) (nextRequest *Request, numValidOrders int, err error)
	// ParseRequestMetadata converts raw request metadata into a concrete type
	// that the subprotocol expects.
	ParseRequestMetadata(metadata json.RawMessage) (interface{}, error)
	// ParseResponseMetadata converts raw response metadata into a concrete type
	// that the subprotocol expects.
	ParseResponseMetadata(metadata json.RawMessage) (interface{}, error)
	// GenerateFirstRequestMetadata generates the metadata for the first request
	// that should be made with this subprotocol.
	GenerateFirstRequestMetadata() (json.RawMessage, error)
}

// New creates and returns a new ordersync service, which is used for both
// requesting orders from other peers and providing orders to peers who request
// them. New expects an array of subprotocols which the service will support, in the
// order of preference. The service will automatically pick the most preferred protocol
// that is supported by both peers for each request/response.
func New(ctx context.Context, node *p2p.Node, subprotocols []Subprotocol) *Service {
	sids := []string{}
	supportedSubprotocols := map[string]Subprotocol{}
	for _, subp := range subprotocols {
		sids = append(sids, subp.Name())
		supportedSubprotocols[subp.Name()] = subp
	}
	// TODO(jalextowle): We should ensure that there were no duplicates -- there
	// is no reason to support this.
	s := &Service{
		ctx:                   ctx,
		node:                  node,
		subprotocolSet:        supportedSubprotocols,
		preferredSubprotocols: sids,
		requestRateLimiter:    rate.NewLimiter(maxRequestsPerSecond, requestsBurst),
	}
	s.node.SetStreamHandler(ID, s.HandleStream)
	return s
}

// GetMatchingSubprotocol returns the most preferred subprotocol to use
// based on the given request.
func (s *Service) GetMatchingSubprotocol(rawReq *rawRequest) (Subprotocol, int, error) {
	for i, protoID := range rawReq.Subprotocols {
		subprotocol, found := s.subprotocolSet[protoID]
		if found {
			return subprotocol, i, nil
		}
	}

	err := NoMatchingSubprotocolsError{
		Requested: rawReq.Subprotocols,
		Supported: s.preferredSubprotocols,
	}
	return nil, 0, err
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
		rawReq, err := waitForRequest(s.ctx, stream)
		if err != nil {
			log.WithError(err).Warn("waitForRequest returned error")
			return
		}
		log.WithFields(log.Fields{
			"requester": stream.Conn().RemotePeer().Pretty(),
		}).Trace("received ordersync request")
		metrics.OrdersyncRequestsReceived.WithLabelValues(metrics.ProtocolV3).Inc()
		rawRes := s.handleRawRequest(rawReq, requesterID)
		if rawRes == nil {
			return
		}
		if err := json.NewEncoder(stream).Encode(rawRes); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": requesterID.Pretty(),
			}).Warn("could not encode ordersync response")
			s.handlePeerScoreEvent(requesterID, psUnexpectedDisconnect)
			return
		}
		if rawRes.Complete {
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
	nextRequestForPeer := map[peer.ID]*rawRequest{}
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

		currentNeighbors := s.node.Neighbors()
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
					metrics.OrdersyncRequestsSent.WithLabelValues(metrics.ProtocolV3, metrics.OrdersyncSuccess).Inc()
					m.Lock()
					if nextFirstRequest != nil {
						nextRequestForPeer[id] = nextFirstRequest
					}
					m.Unlock()
				} else {
					log.WithFields(log.Fields{
						"provider": id.Pretty(),
					}).Trace("successfully got orders from peer via ordersync")
					metrics.OrdersyncRequestsSent.WithLabelValues(metrics.ProtocolV3, metrics.OrdersyncFailure).Inc()
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

func (s *Service) handleRawRequest(rawReq *rawRequest, requesterID peer.ID) *rawResponse {
	if rawReq.Type != TypeRequest {
		log.WithField("gotType", rawReq.Type).Warn("wrong type for Request")
		s.handlePeerScoreEvent(requesterID, psInvalidMessage)
		return nil
	}
	subprotocol, i, err := s.GetMatchingSubprotocol(rawReq)
	if err != nil {
		log.WithError(err).Warn("GetMatchingSubprotocol returned error")
		s.handlePeerScoreEvent(requesterID, psSubprotocolNegotiationFailed)
		return nil
	}
	if len(rawReq.Subprotocols) > 1 {
		firstRequests := FirstRequestsForSubprotocols{}
		err := json.Unmarshal(rawReq.Metadata, &firstRequests)

		// NOTE(jalextowle): Older versions of Mesh did not include
		// metadata in the first ordersync request. In order to handle
		// this in a backwards compatible way, we simply avoid updating
		// the request metadata if there was an error decoding the
		// metadata from the request or if the length of the
		// MetadataForSubprotocol is too small (or empty). This latter
		// check also ensures that the array is long enough for us
		// to access the i-th element.
		if err == nil && len(firstRequests.MetadataForSubprotocol) > i {
			rawReq.Metadata = firstRequests.MetadataForSubprotocol[i]
		}
	}
	res, err := handleRequestWithSubprotocol(s.ctx, subprotocol, requesterID, rawReq)
	if err != nil {
		log.WithError(err).Warn("subprotocol returned error")
		return nil
	}
	encodedMetadata, err := json.Marshal(res.Metadata)
	if err != nil {
		log.WithError(err).Error("could not encode raw metadata")
		return nil
	}
	s.handlePeerScoreEvent(requesterID, psValidMessage)
	return &rawResponse{
		Type:        TypeResponse,
		Subprotocol: subprotocol.Name(),
		Orders:      res.Orders,
		Complete:    res.Complete,
		Metadata:    encodedMetadata,
	}
}

func handleRequestWithSubprotocol(ctx context.Context, subprotocol Subprotocol, requesterID peer.ID, rawReq *rawRequest) (*Response, error) {
	req, err := parseRequestWithSubprotocol(subprotocol, requesterID, rawReq)
	if err != nil {
		return nil, err
	}
	return subprotocol.HandleOrderSyncRequest(ctx, req)
}

func parseRequestWithSubprotocol(subprotocol Subprotocol, requesterID peer.ID, rawReq *rawRequest) (*Request, error) {
	metadata, err := subprotocol.ParseRequestMetadata(rawReq.Metadata)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequesterID: requesterID,
		Metadata:    metadata,
	}, nil
}

func parseResponseWithSubprotocol(subprotocol Subprotocol, providerID peer.ID, rawRes *rawResponse) (*Response, error) {
	metadata, err := subprotocol.ParseResponseMetadata(rawRes.Metadata)
	if err != nil {
		return nil, err
	}
	return &Response{
		ProviderID: providerID,
		Orders:     rawRes.Orders,
		Complete:   rawRes.Complete,
		Metadata:   metadata,
	}, nil
}

type FirstRequestsForSubprotocols struct {
	MetadataForSubprotocol []json.RawMessage `json:"metadata"`
}

// createFirstRequestForAllSubprotocols creates an initial ordersync request that
// contains metadata for all of the ordersync subprotocols.
func (s *Service) createFirstRequestForAllSubprotocols() (*rawRequest, error) {
	metadata := []json.RawMessage{}
	for _, sid := range s.preferredSubprotocols {
		subp := s.subprotocolSet[sid]
		m, err := subp.GenerateFirstRequestMetadata()
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, m)
	}
	encodedMetadata, err := json.Marshal(FirstRequestsForSubprotocols{
		MetadataForSubprotocol: metadata,
	})
	if err != nil {
		return nil, err
	}
	return &rawRequest{
		Type:         TypeRequest,
		Subprotocols: s.preferredSubprotocols,
		Metadata:     encodedMetadata,
	}, nil
}

func (s *Service) getOrdersFromPeer(ctx context.Context, providerID peer.ID, firstRequest *rawRequest) (*rawRequest, error) {
	stream, err := s.node.NewStream(ctx, providerID, ID)
	if err != nil {
		s.handlePeerScoreEvent(providerID, psUnexpectedDisconnect)
		return nil, err
	}
	defer func() {
		_ = stream.Close()
	}()

	totalValidOrders := 0
	var nextReq *rawRequest
	if firstRequest != nil {
		nextReq = firstRequest
	} else {
		nextReq, err = s.createFirstRequestForAllSubprotocols()
		if err != nil {
			return nil, err
		}
	}
	var numValidOrders int
	nextReq, numValidOrders, err = s.makeOrderSyncRequest(ctx, nextReq, stream, providerID)
	if err != nil {
		return nil, err
	}
	totalValidOrders += numValidOrders
	if totalValidOrders == 0 {
		return nil, ErrNoOrdersFromPeer
	} else if nextReq == nil {
		return nil, nil
	}
	nextFirstReq := nextReq

	for {
		select {
		case <-ctx.Done():
			return nextFirstReq, ctx.Err()
		default:
		}
		nextReq, numValidOrders, err = s.makeOrderSyncRequest(ctx, nextReq, stream, providerID)
		if err != nil {
			return nextFirstReq, err
		}
		totalValidOrders += numValidOrders
		if nextReq == nil {
			err = nil
			if totalValidOrders == 0 {
				err = ErrNoOrdersFromPeer
			}
			return nextFirstReq, err
		}
		nextFirstReq = nextReq
	}
}

// makeOrderSyncRequest sends an ordersync request with the given subprotocol
// to the provider, decodes the response, and returns the next raw request (if applicable).
func (s *Service) makeOrderSyncRequest(
	ctx context.Context,
	rawReq *rawRequest,
	stream network.Stream,
	providerID peer.ID,
) (*rawRequest, int, error) {
	if err := json.NewEncoder(stream).Encode(rawReq); err != nil {
		s.handlePeerScoreEvent(providerID, psUnexpectedDisconnect)
		return nil, 0, err
	}

	rawRes, err := waitForResponse(ctx, stream)
	if err != nil {
		return nil, 0, err
	}
	s.handlePeerScoreEvent(providerID, psValidMessage)

	subprotocol, found := s.subprotocolSet[rawRes.Subprotocol]
	if !found {
		s.handlePeerScoreEvent(providerID, psSubprotocolNegotiationFailed)
		return nil, 0, fmt.Errorf("unsupported subprotocol: %s", subprotocol)
	}
	selectedSubprotocol := subprotocol
	res, err := parseResponseWithSubprotocol(subprotocol, providerID, rawRes)
	if err != nil {
		s.handlePeerScoreEvent(providerID, psInvalidMessage)
		return nil, 0, err
	}

	nextReq, numValidOrders, err := subprotocol.HandleOrderSyncResponse(ctx, res)
	if err != nil {
		return nil, 0, err
	}
	s.handlePeerScoreEvent(providerID, receivedOrders)

	// If the result is marked as complete, no more requests should be made.
	if rawRes.Complete {
		return nil, numValidOrders, nil
	}

	encodedMetadata, err := json.Marshal(nextReq.Metadata)
	if err != nil {
		return nil, numValidOrders, err
	}
	return &rawRequest{
		Type:         TypeRequest,
		Subprotocols: []string{selectedSubprotocol.Name()},
		Metadata:     encodedMetadata,
	}, numValidOrders, nil
}

// shufflePeers randomizes the order of the given list of peers.
func shufflePeers(peers []peer.ID) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })
}

func waitForRequest(parentCtx context.Context, stream network.Stream) (*rawRequest, error) {
	ctx, cancel := context.WithTimeout(parentCtx, requestResponseTimeout)
	defer cancel()
	reqChan := make(chan *rawRequest, 1)
	errChan := make(chan error, 1)
	go func() {
		var rawReq rawRequest
		if err := json.NewDecoder(stream).Decode(&rawReq); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": stream.Conn().RemotePeer().Pretty(),
			}).Warn("could not encode ordersync request")
			errChan <- err
			return
		}
		reqChan <- &rawReq
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
	case rawReq := <-reqChan:
		return rawReq, nil
	}
}

func waitForResponse(parentCtx context.Context, stream network.Stream) (*rawResponse, error) {
	ctx, cancel := context.WithTimeout(parentCtx, requestResponseTimeout)
	defer cancel()
	resChan := make(chan *rawResponse, 1)
	errChan := make(chan error, 1)
	go func() {
		var rawRes rawResponse
		if err := json.NewDecoder(stream).Decode(&rawRes); err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"provider": stream.Conn().RemotePeer().Pretty(),
			}).Warn("could not encode ordersync response")
			errChan <- err
			return
		}
		resChan <- &rawRes
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
