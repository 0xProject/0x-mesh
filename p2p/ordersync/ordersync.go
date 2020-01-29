package ordersync

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	log "github.com/sirupsen/logrus"
)

// ErrNoOrders is returned whenever the orders we are looking for cannot be
// found anywhere on the network. This can mean that we aren't connected to any
// peers on the same topic, that there are no orders for the topic throughout
// the entire network, or that there are peers that have the orders we're
// looking for, but they are refusing to give them to us.
var ErrNoOrders = errors.New("no orders where received from any known peers")

const (
	// ID is the ID for the ordersync protocol.
	ID = protocol.ID("/0x-mesh/order-sync/version/0")

	providerErrorDelay      = 10 * time.Second
	scoreTag                = "/0x-mesh/order-sync"
	inavlidMessageScoreDiff = -10
	validMessageScoreDiff   = 1
)

type GetOrdersRequest struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	// TODO(albrow): Add arbitrary metadata to make the request/response more
	// efficient.
}

type GetOrdersResponse struct {
	Type   string `json:"type"`
	Topic  string `json:"topic"`
	Orders []byte `json:"orders"`
}

// Service is the main entrypoint for running the ordersync protocol. It handles
// responding to and sending ordersync requests.
type Service struct {
	host     host.Host
	provider Provider
}

type Provider interface {
	// TODO(albrow): Add arbitrary metadata to make the request/response more
	// efficient.
	ProvideOrders(topic string) ([]byte, error)
}

func New(h host.Host, provider Provider) *Service {
	s := &Service{
		host:     h,
		provider: provider,
	}
	s.host.SetStreamHandler(ID, s.handleStream)
	return s
}

func (s *Service) handleStream(stream network.Stream) {
	for {
		var req GetOrdersRequest
		if err := json.NewDecoder(stream).Decode(&req); err != nil {
			remotePeerID := stream.Conn().RemotePeer()
			log.WithError(err).WithField("remotePeer", remotePeerID.Pretty()).Warn("received invalid GetOrdersRequest from peer")
			// TODO(albrow): Handle peer scores somewhere else?
			s.host.ConnManager().UpsertTag(remotePeerID, scoreTag, func(current int) int { return current + inavlidMessageScoreDiff })
			_ = stream.Close()
		}
		if req.Type != "getOrdersRequest" {
			log.WithField("gotType", req.Type).Warn("wrong type for GetOrdersRequest")
		}
		orders, err := s.provider.ProvideOrders(req.Topic)
		if err != nil {
			log.WithError(err).Error("ProvideOrders returned error")
			time.Sleep(providerErrorDelay)
			continue
		}
		res := GetOrdersResponse{
			Type:   "getOrdersResponse",
			Topic:  req.Topic,
			Orders: orders,
		}
		if err := json.NewEncoder(stream).Encode(res); err != nil {
			log.WithError(err).Error("could not JSON encode orders returned by ProvideOrders")
			time.Sleep(providerErrorDelay)
			continue
		}
	}
}

func (s *Service) GetOrders(ctx context.Context, topic string) ([]byte, error) {
	peers := s.host.Network().Peers()
	shufflePeers(peers)
	// TODO(albrow): Do this for loop partly in parallel.
	// TODO(albrow): Add a timeout when waiting for a response.
	for _, peerID := range peers {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		stream, err := s.host.NewStream(ctx, peerID, ID)
		if err != nil {
			// TODO(albrow): Detect the type of error. Do we want to return it?
			continue
		}
		req := GetOrdersRequest{
			Type:  "getOrdersRequest",
			Topic: topic,
		}
		if err := json.NewEncoder(stream).Encode(req); err != nil {
			// TODO(albrow): Detect the type of error. Do we want to return it?
			continue
		}
		var res GetOrdersResponse
		if err := json.NewDecoder(stream).Decode(&res); err != nil {
			// TODO(albrow): Detect the type of error. Do we want to return it?
			continue
		}
		if res.Type != "getOrdersResponse" {
			log.WithField("gotType", res.Type).Warn("wrong type for GetOrdersResponse")
		}
		if res.Topic != topic {
			log.WithFields(log.Fields{
				"expectedTopic": topic,
				"gotTopic":      res.Topic,
			}).Warn("wrong topic for GetOrdersResponse")
		}
		if len(res.Orders) != 0 {
			// TODO(albrow): Handle peer scores somewhere else?
			s.host.ConnManager().UpsertTag(peerID, scoreTag, func(current int) int { return current + validMessageScoreDiff })
			return res.Orders, nil
		}
	}

	return nil, ErrNoOrders
}

func shufflePeers(peers []peer.ID) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })
}
