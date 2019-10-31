// +build !js

package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
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
	addOrdersHandler         func(signedOrdersRaw []*json.RawMessage, opts AddOrdersOpts) (*ordervalidator.ValidationResults, error)
	getOrdersHandler         func(page, perPage int, snapshotID string) (*GetOrdersResponse, error)
	addPeerHandler           func(peerInfo peerstore.PeerInfo) error
	getStatsHandler          func() (*GetStatsResponse, error)
	subscribeToOrdersHandler func(ctx context.Context) (*rpc.Subscription, error)
}

func (d *dummyRPCHandler) AddOrders(signedOrdersRaw []*json.RawMessage, opts AddOrdersOpts) (*ordervalidator.ValidationResults, error) {
	if d.addOrdersHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for AddOrder")
	}
	return d.addOrdersHandler(signedOrdersRaw, opts)
}

func (d *dummyRPCHandler) GetOrders(page, perPage int, snapshotID string) (*GetOrdersResponse, error) {
	if d.getOrdersHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for GetOrders")
	}
	return d.getOrdersHandler(page, perPage, snapshotID)
}

func (d *dummyRPCHandler) AddPeer(peerInfo peerstore.PeerInfo) error {
	if d.addPeerHandler == nil {
		return errors.New("dummyRPCHandler: no handler set for AddPeer")
	}
	return d.addPeerHandler(peerInfo)
}

func (d *dummyRPCHandler) GetStats() (*GetStatsResponse, error) {
	if d.getStatsHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for GetStats")
	}
	return d.getStatsHandler()
}

func (d *dummyRPCHandler) SubscribeToOrders(ctx context.Context) (*rpc.Subscription, error) {
	if d.subscribeToOrdersHandler == nil {
		return nil, errors.New("dummyRPCHandler: no handler set for Orders")
	}
	return d.subscribeToOrdersHandler(ctx)
}

// newTestServerAndClient returns a server and client which have been connected
// to one another on the local network. The server will use the given
// orderHandler to handle incoming requests. Useful for testing purposes. Will
// block until both the server and client are running and connected to one
// another.
func newTestServerAndClient(t *testing.T, rpcHandler *dummyRPCHandler, ctx context.Context) (*Server, *Client) {
	// Start a new server.
	server, err := NewServer(":0", rpcHandler)
	require.NoError(t, err)
	go func() {
		err := server.Listen(ctx)
		if err != nil {
			panic(err)
		}
		require.NoError(t, err)
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
	ChainID:               big.NewInt(constants.TestChainID),
	ExchangeAddress:       ethereum.ChainIDToContractAddresses[constants.TestChainID].Exchange,
	MakerAddress:          constants.GanacheAccount0,
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
	MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
	MakerFeeAssetData:     constants.NullBytes,
	TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
	TakerFeeAssetData:     constants.NullBytes,
	Salt:                  big.NewInt(1548619145450),
	MakerFee:              big.NewInt(0),
	TakerFee:              big.NewInt(0),
	MakerAssetAmount:      big.NewInt(3551808554499581700),
	TakerAssetAmount:      big.NewInt(300000000000000),
	ExpirationTimeSeconds: big.NewInt(1548619325),
}

func TestAddOrdersSuccess(t *testing.T) {
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err)

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		addOrdersHandler: func(signedOrdersRaw []*json.RawMessage, opts AddOrdersOpts) (*ordervalidator.ValidationResults, error) {
			require.Len(t, signedOrdersRaw, 1)
			validationResponse := &ordervalidator.ValidationResults{}
			for _, signedOrderRaw := range signedOrdersRaw {
				signedOrder := &zeroex.SignedOrder{}
				err := signedOrder.UnmarshalJSON([]byte(*signedOrderRaw))
				require.NoError(t, err)
				orderHash, err := signedOrder.ComputeOrderHash()
				require.NoError(t, err)
				validationResponse.Accepted = append(validationResponse.Accepted, &ordervalidator.AcceptedOrderInfo{
					OrderHash:                orderHash,
					SignedOrder:              signedOrder,
					FillableTakerAssetAmount: signedOrder.TakerAssetAmount,
					IsNew:                    true,
				})
			}
			wg.Done()
			return validationResponse, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, client := newTestServerAndClient(t, rpcHandler, ctx)

	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}
	validationResponse, err := client.AddOrders(signedTestOrders)
	require.NoError(t, err)
	expectedOrderHash, err := testOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	// We need to call ResetHash so that unexported hash field is equal in later
	// assertions.
	signedTestOrder.ResetHash()

	acceptedOrderInfo := validationResponse.Accepted[0]
	assert.Equal(t, expectedOrderHash, acceptedOrderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, acceptedOrderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, acceptedOrderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

	// The WaitGroup signals that AddOrders was called on the server-side.
	wg.Wait()
}

func TestGetOrdersSuccess(t *testing.T) {
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err)

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	expectedPage := 0
	expectedPerPage := 5
	expectedSnapshotID := ""
	returnedSnapshotID := "0x123"

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		getOrdersHandler: func(page, perPage int, snapshotID string) (*GetOrdersResponse, error) {
			assert.Equal(t, expectedPage, page)
			assert.Equal(t, expectedPerPage, perPage)
			assert.Equal(t, expectedSnapshotID, snapshotID)
			orderHash, err := signedTestOrder.ComputeOrderHash()
			require.NoError(t, err)
			ordersInfos := []*OrderInfo{
				&OrderInfo{
					OrderHash:                orderHash,
					SignedOrder:              signedTestOrder,
					FillableTakerAssetAmount: expectedFillableTakerAssetAmount,
				},
			}
			wg.Done()
			return &GetOrdersResponse{
				SnapshotID:  returnedSnapshotID,
				OrdersInfos: ordersInfos,
			}, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, client := newTestServerAndClient(t, rpcHandler, ctx)

	getOrdersResponse, err := client.GetOrders(expectedPage, expectedPerPage, expectedSnapshotID)
	require.NoError(t, err)
	expectedOrderHash, err := testOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Len(t, getOrdersResponse.OrdersInfos, 1)

	assert.Equal(t, returnedSnapshotID, getOrdersResponse.SnapshotID, "SnapshotID did not match")

	// We need to call ResetHash so that unexported hash field is equal in later
	// assertions.
	signedTestOrder.ResetHash()

	orderInfo := getOrdersResponse.OrdersInfos[0]
	assert.Equal(t, expectedOrderHash, orderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, orderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, orderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, client := newTestServerAndClient(t, rpcHandler, ctx)

	require.NoError(t, client.AddPeer(expectedPeerInfo))

	// The WaitGroup signals that AddPeer was called on the server-side.
	wg.Wait()
}

func TestGetStats(t *testing.T) {
	expectedGetStatsResponse := &GetStatsResponse{
		Version:           "development",
		PubSubTopic:       "/0x-orders/network/development/version/1",
		Rendezvous:        "/0x-mesh/network/development/version/1",
		PeerID:            "16Uiu2HAmJ827EAibLvJxGMj6BvT1tr2e2ssW4cMtpP15qoQqZGSA",
		EthereumChainID: 42,
		LatestBlock: LatestBlock{
			Number: 1,
			Hash:   common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		NumOrders: 0,
		NumPeers:  0,
	}

	// Set up the dummy handler with a getStatsHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rpcHandler := &dummyRPCHandler{
		getStatsHandler: func() (*GetStatsResponse, error) {
			wg.Done()
			return expectedGetStatsResponse, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, client := newTestServerAndClient(t, rpcHandler, ctx)

	getStatsResponse, err := client.GetStats()
	require.NoError(t, err)
	require.Equal(t, expectedGetStatsResponse, getStatsResponse)

	// The WaitGroup signals that GetStats was called on the server-side.
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, client := newTestServerAndClient(t, rpcHandler, ctx)

	orderEventChan := make(chan []*zeroex.OrderEvent)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventChan)
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	// The WaitGroup signals that AddOrder was called on the server-side.
	wg.Wait()
}

func TestHeartbeatSubscription(t *testing.T) {
	ctx := context.Background()

	rpcHandler := &dummyRPCHandler{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, client := newTestServerAndClient(t, rpcHandler, ctx)

	heartbeatChan := make(chan string)
	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
	defer clientSubscription.Unsubscribe()
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	heartbeat := <-heartbeatChan
	assert.Equal(t, "tick", heartbeat)
}
