// +build !js

package integrationtests

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrdersSuccess(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait until the rpc server has been started, and then create an rpc client
	// that connects to the rpc server.
	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Create a new valid order.
	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)

	// Send the "AddOrders" request to the rpc server.
	validationResponse, err := client.AddOrders([]*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)

	// Ensure that the validation results contain only the order that was
	// sent to the rpc server and that the order was marked as valid.
	require.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)
	acceptedOrderInfo := validationResponse.Accepted[0]
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")
	signedTestOrder.ResetHash()
	assert.Equal(t, expectedFillableTakerAssetAmount, acceptedOrderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")
	assert.Equal(t, expectedOrderHash, acceptedOrderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, acceptedOrderInfo.SignedOrder, "signedOrder did not match")

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")

	client, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Create 10 new valid orders.
	ethClient := ethclient.NewClient(ethRPCClient)
	// NOTE(jalextowle): The default balances are not sufficient to create 10 valid
	//                   orders, so we modify the zrx and weth amounts for this test
	newWethAmount := new(big.Int).Div(wethAmount, big.NewInt(10))
	newZrxAmount := new(big.Int).Div(zrxAmount, big.NewInt(10))
	signedTestOrders := make([]*zeroex.SignedOrder, 10)
	for i := 0; i < 10; i++ {
		signedTestOrders[i] = scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, newWethAmount, newZrxAmount)
	}

	// Send the newly created order to "AddOrders." The order is valid, and this should
	// be reflected in the validation results.
	validationResponse, err := client.AddOrders(signedTestOrders)
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 10)
	assert.Len(t, validationResponse.Rejected, 0)

	// Send an initial "GetOrders" request through the rpc client. This request will
	// get all of the orders in the database after the "AddOrders" request. We can
	// test pagination by comparing to this list.
	expectedPage := 0
	expectedPerPage := 10
	expectedSnapshotID := ""
	initialGetOrdersResponse, err := client.GetOrders(expectedPage, expectedPerPage, expectedSnapshotID)
	require.NoError(t, err)
	assert.Len(t, initialGetOrdersResponse.OrdersInfos, expectedPerPage)

	// Ensure that all of the orders that we added to the mesh node are represented in the
	// get orders request.
	for _, signedTestOrder := range signedTestOrders {
		foundMatchingOrder := false
		expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
		require.NoError(t, err)
		signedTestOrder.ResetHash()

		for _, orderInfo := range initialGetOrdersResponse.OrdersInfos {
			if orderInfo.OrderHash.Hex() == expectedOrderHash.Hex() {
				foundMatchingOrder = true
				assert.Equal(t, signedTestOrder, orderInfo.SignedOrder, "signedOrder did not match")
				assert.Equal(t, signedTestOrder.TakerAssetAmount, orderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")
				break
			}
		}

		assert.True(t, foundMatchingOrder, "found no matching entry in the getOrdersResponse")
	}

	// Make a new "GetOrders" request with different pagination parameters.
	expectedPage = 1
	expectedPerPage = 5
	getOrdersResponse, err := client.GetOrders(expectedPage, expectedPerPage, expectedSnapshotID)
	require.NoError(t, err)
	assert.Len(t, getOrdersResponse.OrdersInfos, expectedPerPage)

	// Ensure that the getOrdersResponse has the correct pagination
	for i := range getOrdersResponse.OrdersInfos {
		assert.Equal(t, initialGetOrdersResponse.OrdersInfos[i+5*expectedPage], getOrdersResponse.OrdersInfos[i], "Incorrect order of second getOrdersResponse")
	}

	cancel()
	wg.Wait()
}

func TestAddPeer(t *testing.T) {
	t.Skip("The AddPeer test is currently skipped because of nondeterministic behavior that causes it to intermittently fail")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start two standalone nodes so that one can add the other as a peer
	wg := &sync.WaitGroup{}
	wg.Add(2)
	logMessages1 := make(chan string, 1024)
	logMessages2 := make(chan string, 1024)
	count2 := int(atomic.AddInt32(&nodeCount, 2))
	count1 := count2 - 1
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count1, logMessages1)
	}()
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count2, logMessages2)
	}()

	// Wait for the "starting p2p node" log to be emitted by both nodes.
	// Scrape the logs to get the peer Ids of both nodes and the multiaddresses
	// of the second node.
	var startingP2PLog struct {
		PeerId    string   `json:"myPeerID"`
		Addresses []string `json:"addresses_array"`
	}
	log, err := waitForLogSubstring(ctx, logMessages1, "starting p2p node")
	require.NoError(t, err, "p2p node didn't start")
	err = json.Unmarshal([]byte(log), &startingP2PLog)
	require.NoError(t, err)
	parsedPeerID1, err := peer.IDB58Decode(startingP2PLog.PeerId)
	require.NoError(t, err)
	log, err = waitForLogSubstring(ctx, logMessages2, "starting p2p node")
	require.NoError(t, err, "p2p node didn't start")
	err = json.Unmarshal([]byte(log), &startingP2PLog)
	require.NoError(t, err)
	parsedPeerID2, err := peer.IDB58Decode(startingP2PLog.PeerId)
	require.NoError(t, err)
	multiaddrs := startingP2PLog.Addresses
	parsedMultiaddrs := make([]ma.Multiaddr, len(multiaddrs))
	for i, addr := range multiaddrs {
		parsed, err := ma.NewMultiaddr(addr)
		require.NoError(t, err)
		parsedMultiaddrs[i] = parsed
	}

	client, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count1))
	require.NoError(t, err)

	// Send the "AddPeer" request
	expectedPeerInfo := peerstore.PeerInfo{
		ID:    parsedPeerID2,
		Addrs: parsedMultiaddrs,
	}
	require.NoError(t, client.AddPeer(expectedPeerInfo))

	// Wait for the "found peer who speaks our protocol" log to be emitted by
	// both nodes. Ensure that the peer IDs of the node that was found match
	// the peer IDs of the nodes created in the test.
	var foundPeerLog struct {
		PeerId   string `json:"remotePeerID_string"`
		Protocol string `json:"protocol_string"`
	}
	log, err = waitForLogSubstring(ctx, logMessages1, "found peer who speaks our protocol")
	require.NoError(t, err, "didn't find peer")
	err = json.Unmarshal([]byte(log), &foundPeerLog)
	require.NoError(t, err)
	parsedFoundPeerID2, err := peer.IDB58Decode(foundPeerLog.PeerId)
	require.NoError(t, err)
	assert.Equal(t, parsedFoundPeerID2, parsedPeerID2)
	log, err = waitForLogSubstring(ctx, logMessages2, "found peer who speaks our protocol")
	require.NoError(t, err, "didn't find peer")
	err = json.Unmarshal([]byte(log), &foundPeerLog)
	require.NoError(t, err)
	parsedFoundPeerID1, err := peer.IDB58Decode(foundPeerLog.PeerId)
	require.NoError(t, err)
	assert.Equal(t, parsedFoundPeerID1, parsedPeerID1)

	cancel()
	wg.Wait()
}

func TestGetStats(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait for the rpc server to start and get the peer ID of the node. Start the
	// rpc client after the server has been started,
	var jsonLog struct {
		PeerID string `json:"myPeerID"`
	}
	log, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	err = json.Unmarshal([]byte(log), &jsonLog)
	require.NoError(t, err)
	client, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	getStatsResponse, err := client.GetStats()
	require.NoError(t, err)

	// Ensure that the "LatestBlock" in the stats response is non-nil and has a nonzero block number.
	assert.NotNil(t, getStatsResponse.LatestBlock)
	assert.True(t, getStatsResponse.LatestBlock.Number > 0)

	// NOTE(jalextowle): Since this test uses an actual mesh node, we can't know in advance which block
	//                   should be the latest block.
	getStatsResponse.LatestBlock = rpc.LatestBlock{}

	// Ensure that the correct response was logged by "GetStats"
	expectedGetStatsResponse := &rpc.GetStatsResponse{
		Version:              "development",
		PubSubTopic:          "/0x-orders/network/1337/version/1",
		Rendezvous:           "/0x-mesh/network/1337/version/1",
		PeerID:               jsonLog.PeerID,
		EthereumChainID:      1337,
		LatestBlock:          rpc.LatestBlock{},
		NumOrders:            0,
		NumPeers:             0,
		MaxExpirationTime:    constants.UnlimitedExpirationTime.String(),
		StartOfCurrentUTCDay: ratelimit.GetUTCMidnightOfDate(time.Now()),
	}
	require.Equal(t, expectedGetStatsResponse, getStatsResponse)

	cancel()
	wg.Wait()
}

func TestOrdersSubscription(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait for the rpc server to start and then start the rpc client.
	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Subscribe to order events through the rpc client and ensure that the subscription
	// is valid.
	orderEventChan := make(chan []*zeroex.OrderEvent, 1)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventChan)
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	// Create a valid order and send it to the rpc client's "AddOrders" endpoint.
	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")
	_, err = client.AddOrders([]*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)

	// Ensure that the "AddOrders" request triggered an order event that was
	// passed through the subscription.
	orderEvent := <-orderEventChan
	signedTestOrder.ResetHash()
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	assert.EqualValues(t,
		[]*zeroex.OrderEvent{
			&zeroex.OrderEvent{
				OrderHash:                expectedOrderHash,
				SignedOrder:              signedTestOrder,
				EndState:                 zeroex.ESOrderAdded,
				FillableTakerAssetAmount: expectedFillableTakerAssetAmount,
				ContractEvents:           []*zeroex.ContractEvent{},
			},
		},
		orderEvent,
	)
}

func TestHeartbeatSubscription(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait for the rpc server to start and then start the rpc client
	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Send the "SubscribeToHeartbeat" request through the rpc client and assert
	// that the subscription is not nil.
	heartbeatChan := make(chan string)
	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
	defer clientSubscription.Unsubscribe()
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	// Ensure that a valid heartbeat was received
	heartbeat := <-heartbeatChan
	assert.Equal(t, "tick", heartbeat)
}
