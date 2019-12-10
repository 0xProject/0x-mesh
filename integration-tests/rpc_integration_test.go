package integrationtests

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrdersSuccess(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Remove the old database and p2p files.
	removeOldFiles(t, ctx)

	buildMeshForTests(t, ctx)

	// logMessages is a channel through which log messages from the
	// node will be sent. We use a large buffer so it doesn't cause
	// goroutines to block.
	logMessages := make(chan string, 1024)

	// count is a channel through which the node count that is being used by
	// a particular standalone node process will be communicated.
	count := make(chan int)

	// Start the node in a goroutine.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")

	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	// Block for the count value and then close the channel
	nodeCount := <-count
	close(count)

	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+nodeCount))
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

// TODO(jalextowle): Since the uuid creation process is inherently random, we
//                   can't meaningfully sanity check the returnedSnapshotID in
//                   this test. Unit testing should be implemented to verify that
//                   this logic is correct, if necessary.
func TestGetOrdersSuccess(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	// FIXME(jalextowle): Cancel in case execution stops before the other
	//                    cancel is called
	defer cancel()

	// Remove the old database and p2p files.
	removeOldFiles(t, ctx)

	buildMeshForTests(t, ctx)

	// logMessages is a channel through which log messages from the
	// node will be sent. We use a large buffer so it doesn't cause
	// goroutines to block.
	logMessages := make(chan string, 1024)

	count := make(chan int)

	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount

	expectedPage := 0
	expectedPerPage := 5
	expectedSnapshotID := ""

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	_, err = waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")

	nodeCount := <-count

	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+nodeCount))
	require.NoError(t, err)

	signedTestOrders := []*zeroex.SignedOrder{signedTestOrder}
	validationResponse, err := client.AddOrders(signedTestOrders)
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	getOrdersResponse, err := client.GetOrders(expectedPage, expectedPerPage, expectedSnapshotID)
	require.NoError(t, err)
	expectedOrderHash, err = signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Len(t, getOrdersResponse.OrdersInfos, 1)

	// We need to call ResetHash so that unexported hash field is equal in later
	// assertions.
	signedTestOrder.ResetHash()

	orderInfo := getOrdersResponse.OrdersInfos[0]
	assert.Equal(t, expectedOrderHash, orderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, orderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, orderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

	// Cancel the context and wait for all outstanding goroutines to finish.
	cancel()
	wg.Wait()
}

// FIXME - A good strategy here might involve spinning up two standalone nodes,
//         listen to the logs of each, and get peer information from there. Then
//         have one node add the other node as a peer. It would also be good to
//         test the case in which the other node doesn't actually exist.
/*
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

	ctx, cancel := context.WithCancel(context.Background())

	// Remove the old database and p2p files.
	removeOldFiles(t, ctx)

	buildMeshForTests(t, ctx)

	logMessages := make(chan string, 1024)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, logMessages)
	}()

	_, err = waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err)

	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint)
	require.NoError(t, err)

	require.NoError(t, client.AddPeer(expectedPeerInfo))

	// Cancel the context and wait for all outstanding goroutines to finish.
	cancel()
	wg.Wait()
}
*/

func TestGetStats(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())

	// Remove the old database and p2p files.
	removeOldFiles(t, ctx)

	buildMeshForTests(t, ctx)

	logMessages := make(chan string, 1024)
	count := make(chan int)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	log, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	var jsonLog map[string]interface{}
	err = json.Unmarshal([]byte(log), &jsonLog)
	require.NoError(t, err)

	expectedGetStatsResponse := &rpc.GetStatsResponse{
		Version:              "development",
		PubSubTopic:          "/0x-orders/network/1337/version/1",
		Rendezvous:           "/0x-mesh/network/1337/version/1",
		PeerID:               jsonLog["myPeerID"].(string),
		EthereumChainID:      1337,
		LatestBlock:          rpc.LatestBlock{},
		NumOrders:            0,
		NumPeers:             0,
		MaxExpirationTime:    "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		StartOfCurrentUTCDay: getUTCMidnightOfDate(time.Now()),
	}

	nodeCount := <-count
	close(count)

	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+nodeCount))
	require.NoError(t, err)

	getStatsResponse, err := client.GetStats()

	// HACK(jalextowle): Zeroing the Number and Hash fields of LatestBlock
	//                   allows us to test more of the "GetStats" response
	//                   without being too restrictive about the blockchain
	//                   that is being used.
	getStatsResponse.LatestBlock = rpc.LatestBlock{}

	require.NoError(t, err)
	require.Equal(t, expectedGetStatsResponse, getStatsResponse)

	// Cancel the context and wait for all outstanding goroutines to finish.
	cancel()
	wg.Wait()
}

// FIXME - This needs some work and I'll need to use the new count trick
/*
func TestOrdersSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Remove the old database and p2p files.
	removeOldFiles(t, ctx)

	buildMeshForTests(t, ctx)

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
*/

func TestHeartbeatSubscription(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Remove the old database and p2p files.
	removeOldFiles(t, ctx)

	buildMeshForTests(t, ctx)

	logMessages := make(chan string, 1024)
	count := make(chan int)

	// Set up the dummy handler with an addOrdersHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")

	nodeCount := <-count
	close(count)

	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+nodeCount))
	require.NoError(t, err)

	heartbeatChan := make(chan string)
	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
	defer clientSubscription.Unsubscribe()
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	heartbeat := <-heartbeatChan
	assert.Equal(t, "tick", heartbeat)
}
