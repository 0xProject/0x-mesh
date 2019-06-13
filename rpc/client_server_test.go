// +build !js

package rpc

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dummyRPCHandler is used for testing purposes. It allows declaring handlers
// for some requests or all of them, depending on testing needs.
type dummyRPCHandler struct {
	addOrdersHandler            func(orders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error)
	addPeerHandler              func(peerInfo peerstore.PeerInfo) error
	subscribeToOrdersHandler    func(ctx context.Context) (*rpc.Subscription, error)
	subscribeToHeartbeatHandler func(ctx context.Context) (*rpc.Subscription, error)
}

func (d *dummyRPCHandler) AddOrders(orders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error) {
	if d.addOrdersHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for AddOrder")
	}
	return d.addOrdersHandler(orders)
}

func (d *dummyRPCHandler) AddPeer(peerInfo peerstore.PeerInfo) error {
	if d.addPeerHandler == nil {
		return errors.New("dummyRPCHandler: no handler set for AddPeer")
	}
	return d.addPeerHandler(peerInfo)
}

func (d *dummyRPCHandler) SubscribeToOrders(ctx context.Context) (*rpc.Subscription, error) {
	if d.subscribeToOrdersHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for Orders")
	}
	return d.subscribeToOrdersHandler(ctx)
}

func (d *dummyRPCHandler) SubscribeToHeartbeat(ctx context.Context) (*rpc.Subscription, error) {
	if d.subscribeToHeartbeatHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for Heartbeat")
	}
	return d.subscribeToHeartbeatHandler(ctx)
}

// newTestServerAndClient returns a server and client which have been connected
// to one another on the local network. The server will use the given
// orderHandler to handle incoming requests. Useful for testing purposes. Will
// block until both the server and client are running and connected to one
// another.
func newTestServerAndClient(t *testing.T, rpcHandler *dummyRPCHandler) (*Server, *Client) {
	// Start a new server.
	server, err := NewServer(":0", rpcHandler)
	require.NoError(t, err)
	go func() {
		_ = server.Listen()
	}()

	// We need to wait for the OS to choose an available port and for server.Addr
	// to return a non-nil value.
	for server.Addr() == nil {
		time.Sleep(10 * time.Millisecond)
	}

	// Create a new client which is connected to the server.
	client, err := NewClient("ws://" + server.Addr().String())
	require.NoError(t, err)

	return server, client
}

var testOrder = &zeroex.Order{
	MakerAddress:          constants.GanacheAccount0,
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
	MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
	TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
	Salt:                  big.NewInt(1548619145450),
	MakerFee:              big.NewInt(0),
	TakerFee:              big.NewInt(0),
	MakerAssetAmount:      big.NewInt(3551808554499581700),
	TakerAssetAmount:      big.NewInt(300000000000000),
	ExpirationTimeSeconds: big.NewInt(1548619325),
	ExchangeAddress:       constants.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
}

func TestAddOrdersSuccess(t *testing.T) {
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err)
	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		addOrdersHandler: func(signedOrders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error) {
			assert.Equal(t, signedTestOrders, signedOrders, "AddOrders was called with an unexpected orders argument")
			validationResponse := &zeroex.ValidationResults{}
			for _, signedOrder := range signedOrders {
				orderHash, err := signedOrder.ComputeOrderHash()
				require.NoError(t, err)
				validationResponse.Accepted = append(validationResponse.Accepted, &zeroex.AcceptedOrderInfo{
					OrderHash:                orderHash,
					SignedOrder:              signedOrder,
					FillableTakerAssetAmount: signedOrder.TakerAssetAmount,
				})
			}
			wg.Done()
			return validationResponse, nil
		},
	}

	server, client := newTestServerAndClient(t, rpcHandler)
	defer server.Close()

	validationResponse, err := client.AddOrders(signedTestOrders)
	require.NoError(t, err)
	expectedOrderHash, err := testOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	acceptedOrderInfo := validationResponse.Accepted[0]
	assert.Equal(t, expectedOrderHash, acceptedOrderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, acceptedOrderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, acceptedOrderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

	// The WaitGroup signals that AddOrders was called on the server-side.
	wg.Wait()
}

func TestAddPeer(t *testing.T) {
	// Create the expected PeerInfo
	addr0, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/1234")
	require.NoError(t, err)
	addr1, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5678")
	require.NoError(t, err)
	peerID, err := peer.IDB58Decode("QmagLpXZHNrTraqWpY49xtFmZMTLBWctx2PF96s4aFrj9f")
	require.NoError(t, err)
	expectedPeerInfo := peerstore.PeerInfo{
		ID:    peerID,
		Addrs: []ma.Multiaddr{addr0, addr1},
	}

	// Set up the dummy handler with an addPeerHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		addPeerHandler: func(peerInfo peerstore.PeerInfo) error {
			assert.Equal(t, expectedPeerInfo, peerInfo, "AddPeer was called with an unexpected peerInfo argument")
			wg.Done()
			return nil
		},
	}

	server, client := newTestServerAndClient(t, rpcHandler)
	defer server.Close()

	require.NoError(t, client.AddPeer(expectedPeerInfo))

	// The WaitGroup signals that AddPeer was called on the server-side.
	wg.Wait()
}

func TestOrdersSubscription(t *testing.T) {
	ctx := context.Background()

	// Set up the dummy handler with a subscribeToOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		subscribeToOrdersHandler: func(ctx context.Context) (*rpc.Subscription, error) {
			wg.Done()
			return nil, nil
		},
	}

	server, client := newTestServerAndClient(t, rpcHandler)
	defer server.Close()

	orderEventChan := make(chan []*zeroex.OrderEvent)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventChan)
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	// The WaitGroup signals that AddOrder was called on the server-side.
	wg.Wait()
}

func TestHeartbeatSubscription(t *testing.T) {
	ctx := context.Background()

	// Set up the dummy handler with a subscribeToHeartbeatHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		subscribeToHeartbeatHandler: func(ctx context.Context) (*rpc.Subscription, error) {
			wg.Done()
			return nil, nil
		},
	}

	server, client := newTestServerAndClient(t, rpcHandler)
	defer server.Close()

	heartbeatChan := make(chan string)
	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	// The WaitGroup signals that Heartbeat was called on the server-side.
	wg.Wait()
}
