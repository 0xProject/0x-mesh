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
		},
	)
	require.NoError(t, err)
	subp0 := &simpleOrderSyncSubprotocolV0{
		myPeerID: n.ID(),
	}
	subp1 := &simpleOrderSyncSubprotocolV1{
		myPeerID: n.ID(),
		hostSubp: subp0,
	}
	s := New(context.Background(), n, []Subprotocol{subp0, subp1})

	// Simulate encoding the request, sending it over the wire, and decoding
	// it by simply encoding and then decoding the raw request.
	rawReq := &rawRequest{
		Type:         TypeRequest,
		Subprotocols: []string{subp0.Name(), subp1.Name()},
	}
	encodedReq, err := json.Marshal(rawReq)
	require.NoError(t, err)
	decodedReq := &rawRequest{}
	err = json.Unmarshal(encodedReq, decodedReq)
	require.NoError(t, err)
	// This request has multiple subprotocols included and nil metadata. This
	// has the same structure as requests that would have been sent by older
	// versions of Mesh, and allows us to test that newer Mesh nodes provide
	// backwards compatability as ordersync providers.
	res := s.handleRawRequest(decodedReq, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())
	// Ensure that the correct metadata was returned. We expect the
	decodedMetadata := simpleOrderSyncSubprotocolRequestMetadataV0{}
	err = json.Unmarshal(decodedReq.Metadata, &decodedMetadata)
	require.NoError(t, err)
	encodedMetadata, err := json.Marshal(decodedMetadata)
	require.NoError(t, err)
	assert.Equal(t, res.Metadata, json.RawMessage(encodedMetadata))

	// Test handling a request from a node that is using the new first request
	// encoding scheme.
	rawReq, err = s.createFirstRequestForAllSubprotocols()
	require.NoError(t, err)
	encodedReq, err = json.Marshal(rawReq)
	require.NoError(t, err)
	decodedReq = &rawRequest{}
	err = json.Unmarshal(encodedReq, decodedReq)
	require.NoError(t, err)
	res = s.handleRawRequest(decodedReq, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())
	// Ensure that the correct metadata was returned
	encodedMetadata, err = subp0.GenerateFirstRequestMetadata()
	require.NoError(t, err)
	assert.Equal(t, res.Metadata, json.RawMessage(encodedMetadata))
}

var _ p2p.MessageHandler = &simpleMessageHandler{}

type simpleMessageHandler struct{}

func (s *simpleMessageHandler) HandleMessages(context.Context, []*p2p.Message) error {
	return nil
}

var _ Subprotocol = &simpleOrderSyncSubprotocolV0{}

type simpleOrderSyncSubprotocolV0 struct {
	myPeerID peer.ID
}

type simpleOrderSyncSubprotocolRequestMetadataV0 struct {
	SomeValue interface{} `json:"someValue"`
}

type simpleOrderSyncSubprotocolResponseMetadataV0 struct {
	SomeValue interface{} `json:"someValue"`
}

func (s *simpleOrderSyncSubprotocolV0) Name() string {
	return "/simple-order-sync-subprotocol/v0"
}

func (s *simpleOrderSyncSubprotocolV0) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed simpleOrderSyncSubprotocolRequestMetadataV0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocolV0) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed simpleOrderSyncSubprotocolResponseMetadataV0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocolV0) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
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

func (s *simpleOrderSyncSubprotocolV0) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return &Request{
		RequesterID: s.myPeerID,
		Metadata:    res.Metadata,
	}, len(res.Orders), nil
}

func (s *simpleOrderSyncSubprotocolV0) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(simpleOrderSyncSubprotocolRequestMetadataV0{
		SomeValue: 0,
	})
}

var _ Subprotocol = &simpleOrderSyncSubprotocolV1{}

type simpleOrderSyncSubprotocolV1 struct {
	myPeerID peer.ID
	hostSubp *simpleOrderSyncSubprotocolV0
}

type simpleOrderSyncSubprotocolRequestMetadataV1 struct {
	AnotherValue interface{} `json:"anotherValue"`
}

type simpleOrderSyncSubprotocolResponseMetadataV1 struct {
	AnotherValue interface{} `json:"anotherValue"`
}

func (s *simpleOrderSyncSubprotocolV1) Name() string {
	return "/simple-order-sync-subprotocol/v1"
}

func (s *simpleOrderSyncSubprotocolV1) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed simpleOrderSyncSubprotocolRequestMetadataV1
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocolV1) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed simpleOrderSyncSubprotocolResponseMetadataV1
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocolV1) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
	return s.hostSubp.HandleOrderSyncRequest(ctx, req)
}

func (s *simpleOrderSyncSubprotocolV1) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return s.hostSubp.HandleOrderSyncResponse(ctx, res)
}

func (s *simpleOrderSyncSubprotocolV1) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(simpleOrderSyncSubprotocolRequestMetadataV1{
		AnotherValue: 1,
	})
}
