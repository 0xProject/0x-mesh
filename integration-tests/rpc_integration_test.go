// +build !js

package integrationtests

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait until the rpc server has been started, and then create an rpc client
	// that connects to the rpc server.
	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Create a new valid order.
	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	// Send the "AddOrders" request to the rpc server.
	validationResponse, err := client.AddOrders([]*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)

	// Ensure that the validation validation results contain only the order that was
	// sent to the rpc server and that the order was marked as valid.
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)
	signedTestOrder.ResetHash()
	acceptedOrderInfo := validationResponse.Accepted[0]
	assert.Equal(t, expectedOrderHash, acceptedOrderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, acceptedOrderInfo.SignedOrder, "signedOrder did not match")
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	assert.Equal(t, expectedFillableTakerAssetAmount, acceptedOrderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

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

	expectedPage := 0
	expectedPerPage := 5
	expectedSnapshotID := ""

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")

	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Create a new valid order.
	ethClient := ethclient.NewClient(ethRPCClient)
	signedTestOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	// Send the newly created order to "AddOrders." The order is valid, and this should
	// be reflected in the validation results.
	validationResponse, err := client.AddOrders([]*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	// Send the "GetOrders" request through the rpc client.
	getOrdersResponse, err := client.GetOrders(expectedPage, expectedPerPage, expectedSnapshotID)
	require.NoError(t, err)
	assert.Len(t, getOrdersResponse.OrdersInfos, 1)

	// Ensure that the orders returned by the rpc server match the orders that were
	// sent through the "AddOrders" endpoint.
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	expectedOrderHash, err = signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	signedTestOrder.ResetHash()
	orderInfo := getOrdersResponse.OrdersInfos[0]
	assert.Equal(t, expectedOrderHash, orderInfo.OrderHash, "orderHashes did not match")
	assert.Equal(t, signedTestOrder, orderInfo.SignedOrder, "signedOrder did not match")
	assert.Equal(t, expectedFillableTakerAssetAmount, orderInfo.FillableTakerAssetAmount, "fillableTakerAssetAmount did not match")

	cancel()
	wg.Wait()
}

func TestAddPeer(t *testing.T) {
	t.Skip("The AddPeer test is currently skipped because of nondeterministic behavior that causes it to intermittently fail")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start two standalone nodes so that one can add the other as a peer
	wg := &sync.WaitGroup{}
	wg.Add(2)
	logMessages1 := make(chan string, 1024)
	logMessages2 := make(chan string, 1024)
	count2 := int(atomic.AddInt32(nodeCount, 2))
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

	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+count1))
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
	assert.Equal(t, foundPeerLog.Protocol, protocolString)
	log, err = waitForLogSubstring(ctx, logMessages2, "found peer who speaks our protocol")
	require.NoError(t, err, "didn't find peer")
	err = json.Unmarshal([]byte(log), &foundPeerLog)
	require.NoError(t, err)
	parsedFoundPeerID1, err := peer.IDB58Decode(foundPeerLog.PeerId)
	require.NoError(t, err)
	assert.Equal(t, parsedFoundPeerID1, parsedPeerID1)
	assert.Equal(t, foundPeerLog.Protocol, protocolString)

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
	count := int(atomic.AddInt32(nodeCount, 1))
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
	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	getStatsResponse, err := client.GetStats()
	require.NoError(t, err)

	// HACK(jalextowle): Zeroing the Number and Hash fields of LatestBlock
	//                   allows us to test more of the "GetStats" response
	//                   without being too restrictive about the blockchain
	//                   that is being used.
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
		MaxExpirationTime:    "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		StartOfCurrentUTCDay: getUTCMidnightOfDate(time.Now()),
	}
	require.Equal(t, expectedGetStatsResponse, getStatsResponse)

	cancel()
	wg.Wait()
}

func TestOrdersSubscription(t *testing.T) {
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
	count := int(atomic.AddInt32(nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait for the rpc server to start and then start the rpc client.
	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Subscribe to order events through the rpc client and ensure that the subscription
	// is valid.
	orderEventChan := make(chan []*zeroex.OrderEvent)
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
	signedTestOrder.ResetHash()
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	orderEvent := <-orderEventChan
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	removeOldFiles(t, ctx)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, logMessages)
	}()

	// Wait for the rpc server to start and then start the rpc client
	_, err := waitForLogSubstring(ctx, logMessages, "started RPC server")
	require.NoError(t, err, "RPC server didn't start")
	client, err := rpc.NewClient(standaloneRPCEndpoint + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Send the "SubscribeToHeartbeat" request through the rpc client and assert
	// that the subscription is not nil.
	heartbeatChan := make(chan string)
	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
	defer clientSubscription.Unsubscribe()
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	heartbeat := <-heartbeatChan
	assert.Equal(t, "tick", heartbeat)
}
