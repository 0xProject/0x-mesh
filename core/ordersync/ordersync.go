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

	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/albrow/stringset"
	"github.com/jpillora/backoff"
	network "github.com/libp2p/go-libp2p-core/network"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	peer "github.com/libp2p/go-libp2p-peer"
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
)

var (
	// retryBackoff defines how long to wait before trying again if we didn't get
	// orders from enough peers during the ordersync process.
	retryBackoff = &backoff.Backoff{
		Min:    250 * time.Millisecond, // First back-off length
		Max:    1 * time.Minute,        // Longest back-off length
		Factor: 2,                      // Factor to multiple each successive back-off
	}
	// backoffMut is a mutex around retryBackoff, which otherwise appears to not
	// be goroutine-safe.
	backoffMut = &sync.Mutex{}
	// ErrNoOrders is returned whenever the orders we are looking for cannot be
	// found anywhere on the network. This can mean that we aren't connected to any
	// peers on the same topic, that there are no orders for the topic throughout
	// the entire network, or that there are peers that have the orders we're
	// looking for, but they are refusing to give them to us.
	ErrNoOrders = errors.New("no orders where received from any known peers")
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
	Metadata interface{} `json:"metadata"`
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
	Orders   []*zeroex.SignedOrder `json:"orders"`
	Complete bool                  `json:"complete"`
	Metadata interface{}           `json:"metadata"`
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
	ctx          context.Context
	node         *p2p.Node
	subprotocols map[string]Subprotocol
	// requestRateLimiter is a rate limiter for incoming ordersync requests. It's
	// shared between all peers.
	requestRateLimiter *rate.Limiter
}

// SupportedSubprotocols returns the subprotocols that are supported by the service.
func (s *Service) SupportedSubprotocols() []string {
	sids := []string{}
	for sid := range s.subprotocols {
		sids = append(sids, sid)
	}
	return sids
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
	// the database) and if needed creates and returns the next request that
	// should be sent. If nextRequest is nil, the ordersync protocol is
	// considered finished. HandleOrderSyncResponse is the implementation for the
	// "requester" side of the subprotocol.
	HandleOrderSyncResponse(context.Context, *Response) (nextRequest *Request, err error)
	// ParseRequestMetadata converts raw request metadata into a concrete type
	// that the subprotocol expects.
	ParseRequestMetadata(metadata json.RawMessage) (interface{}, error)
	// ParseResponseMetadata converts raw response metadata into a concrete type
	// that the subprotocol expects.
	ParseResponseMetadata(metadata json.RawMessage) (interface{}, error)
}

// New creates and returns a new ordersync service, which is used for both
// requesting orders from other peers and providing orders to peers who request
// them. New expects an array of subprotocols which the service will support, in the
// order of preference. The service will automatically pick the most preferred protocol
// that is supported by both peers for each request/response.
func New(ctx context.Context, node *p2p.Node, subprotocols []Subprotocol) *Service {
	supportedSubprotocols := map[string]Subprotocol{}
	for _, subp := range subprotocols {
		supportedSubprotocols[subp.Name()] = subp
	}
	s := &Service{
		ctx:                ctx,
		node:               node,
		subprotocols:       supportedSubprotocols,
		requestRateLimiter: rate.NewLimiter(maxRequestsPerSecond, requestsBurst),
	}
	s.node.SetStreamHandler(ID, s.HandleStream)
	return s
}

// GetMatchingSubprotocol returns the most preferred subprotocol to use
// based on the given request.
func (s *Service) GetMatchingSubprotocol(rawReq *rawRequest) (Subprotocol, error) {
	for _, protoID := range rawReq.Subprotocols {
		subprotocol, found := s.subprotocols[protoID]
		if found {
			return subprotocol, nil
		}
	}

	err := NoMatchingSubprotocolsError{
		Requested: rawReq.Subprotocols,
		Supported: s.SupportedSubprotocols(),
	}
	return nil, err
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
	remotePeerID := stream.Conn().RemotePeer()

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
		if rawReq.Type != TypeRequest {
			log.WithField("gotType", rawReq.Type).Warn("wrong type for Request")
			s.handlePeerScoreEvent(remotePeerID, psInvalidMessage)
			return
		}
		subprotocol, err := s.GetMatchingSubprotocol(rawReq)
		if err != nil {
			log.WithError(err).Warn("GetMatchingSubprotocol returned error")
			s.handlePeerScoreEvent(remotePeerID, psSubprotocolNegotiationFailed)
			return
		}
		res, err := handleRequestWithSubprotocol(s.ctx, subprotocol, rawReq)
		if err != nil {
			log.WithError(err).Warn("subprotocol returned error")
			return
		}
		encodedMetadata, err := json.Marshal(res.Metadata)
		if err != nil {
			log.WithError(err).Error("could not encode raw metadata")
			return
		}
		s.handlePeerScoreEvent(remotePeerID, psValidMessage)
		rawRes := rawResponse{
			Type:        TypeResponse,
			Subprotocol: subprotocol.Name(),
			Orders:      res.Orders,
			Complete:    res.Complete,
			Metadata:    encodedMetadata,
		}
		if err := json.NewEncoder(stream).Encode(rawRes); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": remotePeerID.Pretty(),
			}).Warn("could not encode ordersync response")
			s.handlePeerScoreEvent(remotePeerID, psUnexpectedDisconnect)
			return
		}
		if res.Complete {
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

	for len(successfullySyncedPeers) < minPeers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// TODO(albrow): Do this for loop partly in parallel.
		currentNeighbors := s.node.Neighbors()
		shufflePeers(currentNeighbors)
		for _, peerID := range currentNeighbors {
			if len(successfullySyncedPeers) >= minPeers {
				return nil
			}
			if successfullySyncedPeers.Contains(peerID.Pretty()) {
				continue
			}

			log.WithFields(log.Fields{
				"provider": peerID.Pretty(),
			}).Trace("requesting orders from neighbor via ordersync")
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err := s.getOrdersFromPeer(ctx, peerID); err != nil {
				log.WithFields(log.Fields{
					"error":    err.Error(),
					"provider": peerID.Pretty(),
				}).Warn("could not get orders from peer via ordersync")
				continue
			} else {
				// TODO(albrow): Handle case where no orders were returned from this
				// peer. We need to not try them again, but also not count them toward
				// the number of peers we have successfully synced with.
				log.WithFields(log.Fields{
					"provider": peerID.Pretty(),
				}).Trace("succesfully got orders from peer via ordersync")
				successfullySyncedPeers.Add(peerID.Pretty())
			}
		}

		backoffMut.Lock()
		delayBeforeNextRetry := retryBackoff.Duration()
		backoffMut.Unlock()
		log.WithFields(log.Fields{
			"delayBeforeNextRetry":    delayBeforeNextRetry.String(),
			"minPeers":                minPeers,
			"successfullySyncedPeers": len(successfullySyncedPeers),
		}).Debug("ordersync could not get orders from enough peers (trying again soon)")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delayBeforeNextRetry):
			continue
		}
	}

	return nil
}

func handleRequestWithSubprotocol(ctx context.Context, subprotocol Subprotocol, rawReq *rawRequest) (*Response, error) {
	req, err := parseRequestWithSubprotocol(subprotocol, rawReq)
	if err != nil {
		return nil, err
	}
	return subprotocol.HandleOrderSyncRequest(ctx, req)
}

func parseRequestWithSubprotocol(subprotocol Subprotocol, rawReq *rawRequest) (*Request, error) {
	metadata, err := subprotocol.ParseRequestMetadata(rawReq.Metadata)
	if err != nil {
		return nil, err
	}
	return &Request{
		Metadata: metadata,
	}, nil
}

func parseResponseWithSubprotocol(subprotocol Subprotocol, rawRes *rawResponse) (*Response, error) {
	metadata, err := subprotocol.ParseResponseMetadata(rawRes.Metadata)
	if err != nil {
		return nil, err
	}
	return &Response{
		Orders:   rawRes.Orders,
		Complete: rawRes.Complete,
		Metadata: metadata,
	}, nil
}

func (s *Service) getOrdersFromPeer(ctx context.Context, providerID peer.ID) error {
	stream, err := s.node.NewStream(ctx, providerID, ID)
	if err != nil {
		s.handlePeerScoreEvent(providerID, psUnexpectedDisconnect)
		return err
	}
	defer func() {
		_ = stream.Close()
	}()

	var nextReq *Request
	var selectedSubprotocol Subprotocol
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var rawReq *rawRequest
		if nextReq == nil {
			// First request
			rawReq = &rawRequest{
				Type:         TypeRequest,
				Subprotocols: s.SupportedSubprotocols(),
				Metadata:     nil,
			}
		} else {
			encodedMetadata, err := json.Marshal(nextReq.Metadata)
			if err != nil {
				return err
			}
			rawReq = &rawRequest{
				Type:         TypeRequest,
				Subprotocols: []string{selectedSubprotocol.Name()},
				Metadata:     encodedMetadata,
			}
		}

		if err := json.NewEncoder(stream).Encode(rawReq); err != nil {
			s.handlePeerScoreEvent(providerID, psUnexpectedDisconnect)
			return err
		}

		rawRes, err := waitForResponse(ctx, stream)
		if err != nil {
			return err
		}
		s.handlePeerScoreEvent(providerID, psValidMessage)

		subprotocol, found := s.subprotocols[rawRes.Subprotocol]
		if !found {
			s.handlePeerScoreEvent(providerID, psSubprotocolNegotiationFailed)
			return fmt.Errorf("unsupported subprotocol: %s", subprotocol)
		}
		selectedSubprotocol = subprotocol
		res, err := parseResponseWithSubprotocol(subprotocol, rawRes)
		if err != nil {
			s.handlePeerScoreEvent(providerID, psInvalidMessage)
			return err
		}

		nextReq, err = subprotocol.HandleOrderSyncResponse(ctx, res)
		if err != nil {
			return err
		}
		s.handlePeerScoreEvent(providerID, receivedOrders)

		if rawRes.Complete {
			return nil
		}
	}
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
	go func() {
		var rawReq rawRequest
		if err := json.NewDecoder(stream).Decode(&rawReq); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": stream.Conn().RemotePeer().Pretty(),
			}).Warn("could not encode ordersync request")
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
	case rawReq := <-reqChan:
		return rawReq, nil
	}
}

func waitForResponse(parentCtx context.Context, stream network.Stream) (*rawResponse, error) {
	ctx, cancel := context.WithTimeout(parentCtx, requestResponseTimeout)
	defer cancel()
	resChan := make(chan *rawResponse, 1)
	go func() {
		var rawRes rawResponse
		if err := json.NewDecoder(stream).Decode(&rawRes); err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"provider": stream.Conn().RemotePeer().Pretty(),
			}).Warn("could not encode ordersync response")
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
	case rawRes := <-resChan:
		return rawRes, nil
	}
}
