package integrationtests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrdersSuccess(t *testing.T) {
	setupSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())

	// logMessages is a channel through which log messages from the
	// node will be sent. We use a large buffer so it doesn't cause
	// goroutines to block.
	logMessages := make(chan string, 1024)

	// Start the node in a goroutine.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, logMessages)
	}()

	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	_, err = waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint)
	require.NoError(t, err)

	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}
	validationResponse, err := client.AddOrders(signedTestOrders)
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	acceptedOrderInfo := validationResponse.Accepted[0]

	// Reset the hash so that the orders can be compared
	signedTestOrder.ResetHash()

	assert.Equal(t, expectedOrderHash, acceptedOrderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, acceptedOrderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, acceptedOrderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

	// Cancel the context and wait for all outstanding goroutines to finish.
	cancel()
	wg.Wait()
}

/*
func TestGetOrdersSuccess(t *testing.T) {
	setupSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())

	// logMessages is a channel through which log messages from the
	// node will be sent. We use a large buffer so it doesn't cause
	// goroutines to block.
	logMessages := make(chan string, 1024)

	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	expectedPage := 0
	expectedPerPage := 5
	expectedSnapshotID := ""
	returnedSnapshotID := "0x123"

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, logMessages)
	}()

	_, err = waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint)
	require.NoError(t, err)

	getOrdersResponse, err := client.GetOrders(expectedPage, expectedPerPage, expectedSnapshotID)
	require.NoError(t, err)
	expectedOrderHash, err = signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Len(t, getOrdersResponse.OrdersInfos, 1)

	fmt.Printf("%+v\n", signedTestOrder)
	fmt.Printf("%+v\n", getOrdersResponse.OrdersInfos[0].SignedOrder)

	assert.Equal(t, returnedSnapshotID, getOrdersResponse.SnapshotID, "SnapshotID did not match")

	// We need to call ResetHash so that unexported hash field is equal in later
	// assertions.
	signedTestOrder.ResetHash()

	orderInfo := getOrdersResponse.OrdersInfos[0]
	assert.Equal(t, expectedOrderHash, orderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, orderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, orderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

	// The WaitGroup signals that AddOrders was called on the server-side.
	cancel()
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
*/

func TestGetStats(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logMessages := make(chan string, 1024)

	expectedGetStatsResponse := &rpc.GetStatsResponse{
		Version:         "development",
		PubSubTopic:     "/0x-orders/network/development/version/1",
		Rendezvous:      "/0x-mesh/network/development/version/1",
		PeerID:          "16Uiu2HAmJ827EAibLvJxGMj6BvT1tr2e2ssW4cMtpP15qoQqZGSA",
		EthereumChainID: 1337,
		LatestBlock: rpc.LatestBlock{
			Number: 1,
			Hash:   common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		NumOrders: 0,
		NumPeers:  0,
	}

	// Set up the dummy handler with a getStatsHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint)
	require.NoError(t, err)

	getStatsResponse, err := client.GetStats()
	require.NoError(t, err)
	require.Equal(t, expectedGetStatsResponse, getStatsResponse)

	// The WaitGroup signals that GetStats was called on the server-side.
	wg.Wait()
}

func TestOrdersSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logMessages := make(chan string, 1024)

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint)
	require.NoError(t, err)

	orderEventChan := make(chan []*zeroex.OrderEvent)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventChan)
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	// expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	// expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}
	client.AddOrders(signedTestOrders)
	require.NoError(t, err)

	orderEvent := <-orderEventChan

	assert.Equal(t, len(orderEvent), 1)
	fmt.Printf("%+v\n", orderEvent[0].SignedOrder)
}

func TestHeartbeatSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logMessages := make(chan string, 1024)

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint)
	require.NoError(t, err)

	heartbeatChan := make(chan string)
	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
	defer clientSubscription.Unsubscribe()
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	heartbeat := <-heartbeatChan
	assert.Equal(t, "tick", heartbeat)
}
