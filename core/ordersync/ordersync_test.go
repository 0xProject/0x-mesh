package ordersync

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
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
	db, err := db.New(context.Background(), db.TestOptions())
	require.NoError(t, err)
	n, err := p2p.New(
		context.Background(),
		p2p.Config{
			MessageHandler:   &noopMessageHandler{},
			RendezvousPoints: []string{"/test-rendezvous-point"},
			DB:               db,
		},
	)
	require.NoError(t, err)
	subp0 := &oneOrderSubprotocol{
		myPeerID: n.ID(),
	}
	subp1 := &hostedSubprotocol{
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
	// backwards compatibility as ordersync providers.
	res := s.handleRawRequest(rawReq, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())
	// NOTE(jalextowle): Because of the way that nil interfaces are encoded
	// in JSON, the value of `res.Metadata` will not be equal to `rawReq.Metadata`.
	// We simply ensure that `res.Metadata` unmarshals to an empty request metadata
	// object.
	var metadata oneOrderSubprotocolRequestMetadata
	err = json.Unmarshal(res.Metadata, &metadata)
	require.NoError(t, err)
	assert.Equal(t, oneOrderSubprotocolRequestMetadata{}, metadata)

	// Test handling a request from a node that is using the new first request
	// encoding scheme.
	rawReq, err = s.createFirstRequestForAllSubprotocols()
	require.NoError(t, err)
	res = s.handleRawRequest(rawReq, n.ID())
	require.NotNil(t, res)
	assert.True(t, res.Complete)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, res.Subprotocol, subp0.Name())
	assert.Equal(t, res.Metadata, rawReq.Metadata)
}

var _ p2p.MessageHandler = &noopMessageHandler{}

// noopMessageHandler is a dummy message handler that allows a p2p node to be
// instantiated easily in these tests.
type noopMessageHandler struct{}

func (*noopMessageHandler) HandleMessages(context.Context, []*p2p.Message) error {
	return nil
}

var _ Subprotocol = &oneOrderSubprotocol{}

// oneOrderSubprotocol is an ordersync subprotocol that is used for testing
// ordersync. This subprotocol will always respond with a single random test order
// and will duplicate the request metadata in the response.
type oneOrderSubprotocol struct {
	myPeerID peer.ID
}

type oneOrderSubprotocolRequestMetadata struct {
	SomeValue interface{} `json:"someValue"`
}

type oneOrderSubprotocolResponseMetadata struct {
	SomeValue interface{} `json:"someValue"`
}

func (s *oneOrderSubprotocol) Name() string {
	return "/simple-order-sync-subprotocol/v0"
}

func (s *oneOrderSubprotocol) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	if metadata == nil {
		return nil, nil
	}
	var parsed oneOrderSubprotocolRequestMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *oneOrderSubprotocol) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	if metadata == nil {
		return nil, nil
	}
	var parsed oneOrderSubprotocolResponseMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *oneOrderSubprotocol) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
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

func (s *oneOrderSubprotocol) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return &Request{
		RequesterID: s.myPeerID,
		Metadata:    res.Metadata,
	}, len(res.Orders), nil
}

func (s *oneOrderSubprotocol) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(oneOrderSubprotocolRequestMetadata{
		SomeValue: 0,
	})
}

var _ Subprotocol = &hostedSubprotocol{}

// hostedSubprotocol is an ordersync subprotocol that is used for testing
// ordersync. This subprotocol uses oneOrderSubprotocol as a "host" subprotocol
// and delegates the handling of requests and responses to this host.
type hostedSubprotocol struct {
	myPeerID peer.ID
	hostSubp *oneOrderSubprotocol
}

type hostedSubprotocolRequestMetadata struct {
	AnotherValue interface{} `json:"anotherValue"`
}

type hostedSubprotocolResponseMetadata struct {
	AnotherValue interface{} `json:"anotherValue"`
}

func (s *hostedSubprotocol) Name() string {
	return "/simple-order-sync-subprotocol/v1"
}

func (s *hostedSubprotocol) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed hostedSubprotocolRequestMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *hostedSubprotocol) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed hostedSubprotocolResponseMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (s *hostedSubprotocol) HandleOrderSyncRequest(ctx context.Context, req *Request) (*Response, error) {
	return s.hostSubp.HandleOrderSyncRequest(ctx, req)
}

func (s *hostedSubprotocol) HandleOrderSyncResponse(ctx context.Context, res *Response) (*Request, int, error) {
	return s.hostSubp.HandleOrderSyncResponse(ctx, res)
}

func (s *hostedSubprotocol) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(hostedSubprotocolRequestMetadata{
		AnotherValue: 1,
	})
}
