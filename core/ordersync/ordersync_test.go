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
	// FIXME - How to get this to not set up any directories?
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

	// This request has multiple subprotocols included and nil metadata. This
	// has the same structure as requests that would have been sent by older
	// versions of Mesh, and allows us to test that newer Mesh nodes provide
	// backwards compatability as ordersync providers.
	var metadata simpleOrderSyncSubprotocolRequestMetadata
	encodedMetadata, err := json.Marshal(metadata)
	require.NoError(t, err)
	res := s.handleRawRequest(&rawRequest{
		Type:         TypeRequest,
		Subprotocols: []string{subp0.Name(), subp1.Name()},
		Metadata:     encodedMetadata,
	}, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())

	// FIXME - Add more test cases for current nodes
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
type simpleOrderSyncSubprotocolRequestMetadata struct {
	SomeValue interface{} `json:"someValue"`
}

type simpleOrderSyncSubprotocolResponseMetadata struct {
	SomeValue interface{} `json:"someValue"`
}

func (s *simpleOrderSyncSubprotocolV0) Name() string {
	return "/simple-order-sync-subprotocol/v0"
}

func (s *simpleOrderSyncSubprotocolV0) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	if metadata == nil {
		return nil, nil
	}
	var parsed simpleOrderSyncSubprotocolRequestMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *simpleOrderSyncSubprotocolV0) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	if metadata == nil {
		return nil, nil
	}
	var parsed simpleOrderSyncSubprotocolRequestMetadata
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
		Metadata:   nil,
	}, nil
}

func (s *simpleOrderSyncSubprotocolV0) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return &Request{
		RequesterID: s.myPeerID,
		Metadata:    nil,
	}, len(res.Orders), nil
}

func (s *simpleOrderSyncSubprotocolV0) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(simpleOrderSyncSubprotocolRequestMetadata{})
}

var _ Subprotocol = &simpleOrderSyncSubprotocolV1{}

type simpleOrderSyncSubprotocolV1 struct {
	myPeerID peer.ID
	hostSubp *simpleOrderSyncSubprotocolV0
}

func (s *simpleOrderSyncSubprotocolV1) Name() string {
	return "/simple-order-sync-subprotocol/v1"
}

func (s *simpleOrderSyncSubprotocolV1) ParseRequestMetadata(encodedMetadata json.RawMessage) (interface{}, error) {
	return s.ParseRequestMetadata(encodedMetadata)
}

func (s *simpleOrderSyncSubprotocolV1) ParseResponseMetadata(encodedMetadata json.RawMessage) (interface{}, error) {
	return s.ParseResponseMetadata(encodedMetadata)
}

func (s *simpleOrderSyncSubprotocolV1) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
	return s.hostSubp.HandleOrderSyncRequest(ctx, req)
}

func (s *simpleOrderSyncSubprotocolV1) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return s.hostSubp.HandleOrderSyncResponse(ctx, res)
}

func (s *simpleOrderSyncSubprotocolV1) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return s.hostSubp.GenerateFirstRequestMetadata()
}
