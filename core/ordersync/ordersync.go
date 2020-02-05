package ordersync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/0xProject/0x-mesh/p2p"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/albrow/stringset"
	network "github.com/libp2p/go-libp2p-core/network"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	peer "github.com/libp2p/go-libp2p-peer"
	log "github.com/sirupsen/logrus"
)

const (
	TypeRequest  = "Request"
	TypeResponse = "Response"
)

var (
	// ErrNoOrders is returned whenever the orders we are looking for cannot be
	// found anywhere on the network. This can mean that we aren't connected to any
	// peers on the same topic, that there are no orders for the topic throughout
	// the entire network, or that there are peers that have the orders we're
	// looking for, but they are refusing to give them to us.
	ErrNoOrders = errors.New("no orders where received from any known peers")
)

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

	providerErrorDelay      = 10 * time.Second
	scoreTag                = "/0x-mesh/order-sync"
	inavlidMessageScoreDiff = -10
	validMessageScoreDiff   = 1
)

type Request struct {
	Metadata interface{} `json:"metadata"`
}

type rawRequest struct {
	Type         string          `json:"type"`
	Subprotocols []string        `json:"subprotocols"`
	Metadata     json.RawMessage `json:"metadata"`
}

type Response struct {
	Orders   []*zeroex.SignedOrder `json:"orders"`
	Complete bool                  `json:"complete"`
	Metadata interface{}           `json:"metadata"`
}

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
	node         *p2p.Node
	subprotocols map[string]Subprotocol
}

func (s *Service) SupportedSubprotocols() []string {
	sids := []string{}
	for sid := range s.subprotocols {
		sids = append(sids, sid)
	}
	return sids
}

type Subprotocol interface {
	Name() string
	GetOrders(*Request) (*Response, error)
	HandleOrders(*Response) (*Request, error)
	ParseRequestMetadata(metadata json.RawMessage) (interface{}, error)
	ParseResponseMetadata(metadata json.RawMessage) (interface{}, error)
}

func New(node *p2p.Node, subprotocols []Subprotocol) *Service {
	supportedSubprotocols := map[string]Subprotocol{}
	for _, subp := range subprotocols {
		supportedSubprotocols[subp.Name()] = subp
	}
	return &Service{
		node:         node,
		subprotocols: supportedSubprotocols,
	}
}

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

func (s *Service) HandleStream(stream network.Stream) {
	defer func() {
		_ = stream.Close()
	}()
	remotePeerID := stream.Conn().RemotePeer()
	for {
		// TODO(albrow): Close stream if we haven't received any requests in a
		// while.
		var rawReq rawRequest
		if err := json.NewDecoder(stream).Decode(&rawReq); err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"requester": remotePeerID.Pretty(),
			}).Warn("could not encode ordersync request")
			// TODO(albrow): Handle peer scores somewhere else?
			// s.host.ConnManager().UpsertTag(remotePeerID, scoreTag, func(current int) int { return current + inavlidMessageScoreDiff })
			return
		}
		if rawReq.Type != TypeRequest {
			log.WithField("gotType", rawReq.Type).Warn("wrong type for Request")
			return
		}
		subprotocol, err := s.GetMatchingSubprotocol(&rawReq)
		if err != nil {
			log.WithError(err).Warn("GetMatchingSubprotocol returned error")
			return
		}
		res, err := handleRequestWithSubprotocol(subprotocol, &rawReq)
		if err != nil {
			log.WithError(err).Warn("subprotocol returned error")
			return
		}
		encodedMetadata, err := json.Marshal(res.Metadata)
		if err != nil {
			log.WithError(err).Error("could not encode raw metadata")
			return
		}
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
		}
		if res.Complete {
			return
		}
	}
}

func handleRequestWithSubprotocol(subprotocol Subprotocol, rawReq *rawRequest) (*Response, error) {
	req, err := parseRequestWithSubprotocol(subprotocol, rawReq)
	if err != nil {
		return nil, err
	}
	return subprotocol.GetOrders(req)
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

func (s *Service) GetOrders(ctx context.Context, minPeers int) error {
	allPeers := s.node.Neighbors()
	successfullySyncedPeers := stringset.New()

	for len(successfullySyncedPeers) < minPeers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// TODO(albrow): Do this for loop partly in parallel.
		// TODO(albrow): Add a timeout when waiting for a response.
		shufflePeers(allPeers)
		for _, peerID := range allPeers {
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
			} else {
				log.WithFields(log.Fields{
					"provider": peerID.Pretty(),
				}).Trace("succesfully got orders from peer via ordersync")
				successfullySyncedPeers.Add(peerID.Pretty())
			}
		}
	}

	return nil
}

func (s *Service) getOrdersFromPeer(ctx context.Context, providerID peer.ID) error {
	stream, err := s.node.NewStream(ctx, providerID, ID)
	if err != nil {
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
			return err
		}
		var rawRes rawResponse
		if err := json.NewDecoder(stream).Decode(&rawRes); err != nil {
			return err
		}
		subprotocol, found := s.subprotocols[rawRes.Subprotocol]
		if !found {
			return fmt.Errorf("unsupported subprotocol: %s", subprotocol)
		}
		selectedSubprotocol = subprotocol
		res, err := parseResponseWithSubprotocol(subprotocol, &rawRes)
		if err != nil {
			return err
		}
		nextReq, err = subprotocol.HandleOrders(res)
		if err != nil {
			return err
		}

		if rawRes.Complete {
			return nil
		}
	}
}

func shufflePeers(peers []peer.ID) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })
}
