package ordersync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p-core/event"
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
	ProvideOrders(topic string, requestingPeer peer.ID) ([]byte, error)
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
	defer func() {
		_ = stream.Close()
	}()
	remotePeerID := stream.Conn().RemotePeer()
	for {
		// TODO(albrow): Close stream if we haven't received any requests in a
		// while.
		var req GetOrdersRequest
		if err := json.NewDecoder(stream).Decode(&req); err != nil {
			log.WithError(err).WithField("remotePeer", remotePeerID.Pretty()).Warn("received invalid GetOrdersRequest from peer")
			// TODO(albrow): Handle peer scores somewhere else?
			s.host.ConnManager().UpsertTag(remotePeerID, scoreTag, func(current int) int { return current + inavlidMessageScoreDiff })
			return
		}
		if req.Type != "getOrdersRequest" {
			log.WithField("gotType", req.Type).Warn("wrong type for GetOrdersRequest")
		}
		orders, err := s.provider.ProvideOrders(req.Topic, remotePeerID)
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
			// TODO(albrow): Close stream if we couldn't write to it.
			log.WithError(err).Error("could not JSON encode orders returned by ProvideOrders")
			time.Sleep(providerErrorDelay)
			continue
		}
	}
}

func (s *Service) GetOrders(ctx context.Context, topic string) ([]byte, error) {
	log.WithFields(log.Fields{
		"me": s.host.ID().Pretty(),
	}).Trace("inside ordersync.GetOrders")
	// if err := s.waitForPeers(ctx); err != nil {
	// 	return nil, err
	// }
	// peers, err := s.getOrderSyncPeers()
	// if err != nil {
	// 	return nil, err
	// }
	peers := s.host.Network().Peers()
	shufflePeers(peers)
	// TODO(albrow): Do this for loop partly in parallel.
	// TODO(albrow): Add a timeout when waiting for a response.
	for _, peerID := range peers {
		log.WithFields(log.Fields{
			"me":       s.host.ID().Pretty(),
			"provider": peerID.Pretty(),
		}).Trace("requesting orders from neighbor")
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
		// TODO(albrow): Close the stream when we're done.
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

func (s *Service) waitForPeers(ctx context.Context) error {
	// Subscribe to the event bus so we can be notified when we connect to peers
	// that speak new protocols.
	//
	// NOTE(albrow): The ordering here is really important. We need to subscribe
	// to the event bus *and then* check if we are already connected to peers that
	// speak the ordersync protocol. Otherwise we could potentially miss events.
	eventSub, err := s.host.EventBus().Subscribe(new(event.EvtPeerProtocolsUpdated))
	if err != nil {
		return err
	}
	defer eventSub.Close()

	// Check if we already are connected to any peers that speak ordersync.
	orderSyncPeers, err := s.getOrderSyncPeers()
	if err != nil {
		return err
	}
	if len(orderSyncPeers) > 0 {
		// We already are connected to peers that speak the ordersync protocol. No
		// need to wait.
		return nil
	}

	select {
	case <-ctx.Done():
		return nil
	case ev := <-eventSub.Out():
		updatedEvent, ok := ev.(*event.EvtPeerProtocolsUpdated)
		if !ok {
			log.WithField("eventType", fmt.Sprintf("%T", ev)).Error("unexpected event type received from bus")
		}
		for _, added := range updatedEvent.Added {
			// We connected to at least one peer that speaks the ordersync protocol.
			if added == ID {
				log.WithField("peerID", updatedEvent.Peer.Pretty).Info("found peer who speaks ordersync protocol")
				return nil
			}
		}
	}
	return nil
}

func shufflePeers(peers []peer.ID) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })
}

func (s *Service) getOrderSyncPeers() ([]peer.ID, error) {
	allPeers := s.host.Network().Peers()
	orderSyncPeers := []peer.ID{}
	for _, peerID := range allPeers {
		protocols, err := s.host.Peerstore().GetProtocols(peerID)
		if err != nil {
			return nil, err
		}
		for _, protocol := range protocols {
			if protocol == string(ID) {
				orderSyncPeers = append(orderSyncPeers, peerID)
			}
		}
	}
	return orderSyncPeers, nil
}
