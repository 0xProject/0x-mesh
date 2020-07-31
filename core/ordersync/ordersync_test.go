package ordersync

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateDelayWithJitters(t *testing.T) {
	numCalls := 100
	approxDelay := 10 * time.Second
	jitterAmount := 0.1
	for i := 0; i < numCalls; i++ {
		actualDelay := calculateDelayWithJitter(approxDelay, jitterAmount)
		// 0.1 * 10 seconds is 1 second. So we assert that the actual delay is within 1 second
		// of the approximate delay.
		assert.InDelta(t, approxDelay, actualDelay, float64(1*time.Second), "actualDelay: %s", actualDelay)
	}
}

func TestHandleRawRequest(t *testing.T) {
	n, err := p2p.New(
		context.Background(),
		p2p.Config{
			MessageHandler:   &simpleMessageHandler{},
			RendezvousPoints: []string{"/test-rendezvous-point"},
			DataDir:          "/tmp",
		},
	)
	require.NoError(t, err)
	subp0 := &simpleOrderSyncSubprotocol0{
		myPeerID: n.ID(),
	}
	subp1 := &simpleOrderSyncSubprotocol1{
		myPeerID: n.ID(),
		hostSubp: subp0,
	}
	s := New(context.Background(), n, []Subprotocol{subp0, subp1})

	rawReq := &rawRequest{
		Type:         TypeRequest,
		Subprotocols: []string{subp0.Name(), subp1.Name()},
	}
	// This request has multiple subprotocols included and nil metadata. This
	// has the same structure as requests that would have been sent by older
	// versions of Mesh, and allows us to test that newer Mesh nodes provide
	// backwards compatability as ordersync providers.
	res := s.handleRawRequest(rawReq, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())
	// NOTE(jalextowle): Because of the way that nil interfaces are encoded
	// in JSON, the value of `res.Metadata` will not be equal to `rawReq.Metadata`.
	// We simply ensure that `res.Metadata` unmarshals to an empty request metadata
	// object.
	var metadata simpleOrderSyncSubprotocolRequestMetadata0
	err = json.Unmarshal(res.Metadata, &metadata)
	assert.Equal(t, simpleOrderSyncSubprotocolRequestMetadata0{}, metadata)

	// Test handling a request from a node that is using the new first request
	// encoding scheme.
	rawReq, err = s.createFirstRequestForAllSubprotocols()
	res = s.handleRawRequest(rawReq, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())
	assert.Equal(t, res.Metadata, rawReq.Metadata)
}

var _ p2p.MessageHandler = &simpleMessageHandler{}

// simpleMessageHandler is a dummy message handler that allows a p2p node to be
// instantiated easily in these tests.
type simpleMessageHandler struct{}

func (s *simpleMessageHandler) HandleMessages(context.Context, []*p2p.Message) error {
	return nil
}

var _ Subprotocol = &simpleOrderSyncSubprotocol0{}

// simpleOrderSyncSubprotocol0 is an ordersync subprotocol that is used for testing
// ordersync. This subprotocol will always respond with a single random test order
// and will duplicate the request metadata in the response.
type simpleOrderSyncSubprotocol0 struct {
	myPeerID peer.ID
}

type simpleOrderSyncSubprotocolRequestMetadata0 struct {
	SomeValue interface{} `json:"someValue"`
}

type simpleOrderSyncSubprotocolResponseMetadata0 struct {
	SomeValue interface{} `json:"someValue"`
}

func (s *simpleOrderSyncSubprotocol0) Name() string {
	return "/simple-order-sync-subprotocol/v0"
}

func (s *simpleOrderSyncSubprotocol0) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	if metadata == nil {
		return nil, nil
	}
	var parsed simpleOrderSyncSubprotocolRequestMetadata0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocol0) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	if metadata == nil {
		return nil, nil
	}
	var parsed simpleOrderSyncSubprotocolResponseMetadata0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocol0) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
	order := &zeroex.Order{
		ChainID:               big.NewInt(constants.TestChainID),
		MakerAddress:          constants.GanacheAccount1,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   constants.NullAddress,
		MakerAssetData:        scenario.ZRXAssetData,
		MakerFeeAssetData:     constants.NullBytes,
		TakerAssetData:        scenario.WETHAssetData,
		TakerFeeAssetData:     constants.NullBytes,
		Salt:                  big.NewInt(int64(time.Now().Nanosecond())),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(100),
		TakerAssetAmount:      big.NewInt(42),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.GanacheAddresses.Exchange,
	}
	signedOrder, err := zeroex.SignTestOrder(order)
	if err != nil {
		return nil, err
	}
	return &Response{
		ProviderID: s.myPeerID,
		Orders:     []*zeroex.SignedOrder{signedOrder},
		Complete:   true,
		Metadata:   req.Metadata,
	}, nil
}

func (s *simpleOrderSyncSubprotocol0) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return &Request{
		RequesterID: s.myPeerID,
		Metadata:    res.Metadata,
	}, len(res.Orders), nil
}

func (s *simpleOrderSyncSubprotocol0) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(simpleOrderSyncSubprotocolRequestMetadata0{
		SomeValue: 0,
	})
}

var _ Subprotocol = &simpleOrderSyncSubprotocol1{}

// simpleOrderSyncSubprotocol1 is an ordersync subprotocol that is used for testing
// ordersync. This subprotocol uses simpleOrderSyncSubprotocol0 as a "host" subprotocol
// and delegates the handling of requests and responses to this host.
type simpleOrderSyncSubprotocol1 struct {
	myPeerID peer.ID
	hostSubp *simpleOrderSyncSubprotocol0
}

type simpleOrderSyncSubprotocolRequestMetadata1 struct {
	AnotherValue interface{} `json:"anotherValue"`
}

type simpleOrderSyncSubprotocolResponseMetadata1 struct {
	AnotherValue interface{} `json:"anotherValue"`
}

func (s *simpleOrderSyncSubprotocol1) Name() string {
	return "/simple-order-sync-subprotocol/v1"
}

func (s *simpleOrderSyncSubprotocol1) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed simpleOrderSyncSubprotocolRequestMetadata1
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocol1) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed simpleOrderSyncSubprotocolResponseMetadata1
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocol1) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
	return s.hostSubp.HandleOrderSyncRequest(ctx, req)
}

func (s *simpleOrderSyncSubprotocol1) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return s.hostSubp.HandleOrderSyncResponse(ctx, res)
}

func (s *simpleOrderSyncSubprotocol1) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(simpleOrderSyncSubprotocolRequestMetadata1{
		AnotherValue: 1,
	})
}
