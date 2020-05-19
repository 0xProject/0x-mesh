// +build !js

package integrationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWSAddOrdersSuccess(t *testing.T) {
	runAddOrdersSuccessTest(t, standaloneWSRPCEndpointPrefix, "WS", wsRPCPort)
}

func TestHTTPAddOrdersSuccess(t *testing.T) {
	runAddOrdersSuccessTest(t, standaloneHTTPRPCEndpointPrefix, "HTTP", httpRPCPort)
}

func runAddOrdersSuccessTest(t *testing.T, rpcEndpointPrefix, rpcServerType string, rpcPort int) {
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
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	// Wait until the rpc server has been started, and then create an rpc client
	// that connects to the rpc server.
	_, err := waitForLogSubstring(ctx, logMessages, fmt.Sprintf("started %s RPC server", rpcServerType))
	require.NoError(t, err, fmt.Sprintf("%s RPC server didn't start", rpcServerType))
	client, err := rpc.NewClient(rpcEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Create a new valid order.
	signedTestOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
	time.Sleep(500 * time.Millisecond)

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

func TestWSGetOrders(t *testing.T) {
	runGetOrdersTest(t, standaloneWSRPCEndpointPrefix, "WS", wsRPCPort)
}

func TestHTTPGetOrders(t *testing.T) {
	runGetOrdersTest(t, standaloneHTTPRPCEndpointPrefix, "HTTP", httpRPCPort)
}

// TODO(jalextowle): Since the uuid creation process is inherently random, we
//                   can't meaningfully sanity check the returnedSnapshotID in
//                   this test. Unit testing should be implemented to verify that
//                   this logic is correct, if necessary.
func runGetOrdersTest(t *testing.T, rpcEndpointPrefix, rpcServerType string, rpcPort int) {
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
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	_, err := waitForLogSubstring(ctx, logMessages, fmt.Sprintf("started %s RPC server", rpcServerType))
	require.NoError(t, err, fmt.Sprintf("%s RPC server didn't start", rpcServerType))

	client, err := rpc.NewClient(rpcEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	// Create 10 new valid orders.
	numOrders := 10
	orderOptions := scenario.OptionsForAll(orderopts.SetupMakerState(true))
	signedTestOrders := scenario.NewSignedTestOrdersBatch(t, numOrders, orderOptions)
	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
	time.Sleep(500 * time.Millisecond)

	// Send the newly created order to "AddOrders." The order is valid, and this should
	// be reflected in the validation results.
	validationResponse, err := client.AddOrders(signedTestOrders)
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, numOrders)
	assert.Len(t, validationResponse.Rejected, 0)

	fixmeGetOrdersResponse, err := client.GetOrders(0, 10, "")
	require.NoError(t, err)
	// NOTE(jalextowle) This statement holds true for many pagination algorithms, but it may be necessary
	//                  to drop this requirement if the `GetOrders` endpoint changes dramatically.
	require.Len(t, fixmeGetOrdersResponse.OrdersInfos, 10)

	// Make a new "GetOrders" request with different pagination parameters.
	snapshotID := ""
	for _, testCase := range []struct {
		ordersPerPage int
	}{
		{
			ordersPerPage: -1,
		},
		{
			ordersPerPage: 0,
		},
		{
			ordersPerPage: 3,
		},
		{
			ordersPerPage: 5,
		},
	} {
		if testCase.ordersPerPage <= 0 {
			_, err := client.GetOrders(0, testCase.ordersPerPage, snapshotID)
			require.EqualError(t, err, "perPage cannot be zero")
		} else {

			// If numOrders % testCase.ordersPerPage is nonzero, then we must increment the number of pages to
			// iterate through because the numOrder / testCase.ordersPerPage calculation rounds down.
			highestPageNumber := numOrders / testCase.ordersPerPage
			if numOrders%testCase.ordersPerPage > 0 {
				highestPageNumber++
			}

			// Iterate through enough pages to get all of the orders in the mesh nodes database. Compare the
			// responses to the orders that we expect to be in the database.
			var responseOrders []*types.OrderInfo
			for pageNumber := 0; pageNumber < highestPageNumber; pageNumber++ {
				expectedTimestamp := time.Now().UTC()
				getOrdersResponse, err := client.GetOrders(pageNumber, testCase.ordersPerPage, snapshotID)
				assert.WithinDuration(t, expectedTimestamp, getOrdersResponse.SnapshotTimestamp, time.Second)
				require.NoError(t, err)
				// NOTE(jalextowle) This statement holds true for many pagination algorithms, but it may be necessary
				//                  to drop this requirement if the `GetOrders` endpoint changes dramatically.
				require.Len(t, getOrdersResponse.OrdersInfos, min(testCase.ordersPerPage, numOrders-pageNumber*testCase.ordersPerPage))
				responseOrders = append(responseOrders, getOrdersResponse.OrdersInfos...)
			}
			assertSignedOrdersMatch(t, signedTestOrders, responseOrders)
		}
	}

	cancel()
	wg.Wait()
}

func TestWSGetStats(t *testing.T) {
	runGetStatsTest(t, standaloneWSRPCEndpointPrefix, "WS", wsRPCPort)
}

func TestHTTPGetStats(t *testing.T) {
	runGetStatsTest(t, standaloneHTTPRPCEndpointPrefix, "HTTP", httpRPCPort)
}

func runGetStatsTest(t *testing.T, rpcEndpointPrefix, rpcServerType string, rpcPort int) {
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
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	// Wait for the rpc server to start and get the peer ID of the node. Start the
	// rpc client after the server has been started,
	var jsonLog struct {
		PeerID string `json:"myPeerID"`
	}
	log, err := waitForLogSubstring(ctx, logMessages, fmt.Sprintf("started %s RPC server", rpcServerType))
	require.NoError(t, err, fmt.Sprintf("%s RPC server didn't start", rpcServerType))
	err = json.Unmarshal([]byte(log), &jsonLog)
	require.NoError(t, err)
	client, err := rpc.NewClient(rpcEndpointPrefix + strconv.Itoa(rpcPort+count))
	require.NoError(t, err)

	getStatsResponse, err := client.GetStats()
	require.NoError(t, err)

	// Ensure that the "LatestBlock" in the stats response is non-nil and has a nonzero block number.
	assert.NotNil(t, getStatsResponse.LatestBlock)
	assert.True(t, getStatsResponse.LatestBlock.Number > 0)

	// NOTE(jalextowle): Since this test uses an actual mesh node, we can't know in advance which block
	//                   should be the latest block.
	getStatsResponse.LatestBlock = types.LatestBlock{}

	// Ensure that the correct response was logged by "GetStats"
	require.Equal(t, "/0x-orders/version/3/chain/1337/schema/e30=", getStatsResponse.PubSubTopic, "PubSubTopic")
	require.Equal(t, "/0x-mesh/network/1337/version/2", getStatsResponse.Rendezvous, "Rendezvous")
	require.Equal(t, []string{}, getStatsResponse.SecondaryRendezvous, "SecondaryRendezvous")
	require.Equal(t, jsonLog.PeerID, getStatsResponse.PeerID, "PeerID")
	require.Equal(t, 1337, getStatsResponse.EthereumChainID, "EthereumChainID")
	require.Equal(t, types.LatestBlock{}, getStatsResponse.LatestBlock, "LatestBlock")
	require.Equal(t, 0, getStatsResponse.NumOrders, "NumOrders")
	require.Equal(t, 0, getStatsResponse.NumPeers, "NumPeers")
	require.Equal(t, constants.UnlimitedExpirationTime.String(), getStatsResponse.MaxExpirationTime, "MaxExpirationTime")
	require.Equal(t, ratelimit.GetUTCMidnightOfDate(time.Now()), getStatsResponse.StartOfCurrentUTCDay, "StartOfCurrentDay")

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
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	// Wait for the rpc server to start and then start the rpc client.
	_, err := waitForLogSubstring(ctx, logMessages, "started WS RPC server")
	require.NoError(t, err, "WS RPC server didn't start")
	client, err := rpc.NewClient(standaloneWSRPCEndpointPrefix + strconv.Itoa(wsRPCPort+count))
	require.NoError(t, err)

	// Subscribe to order events through the rpc client and ensure that the subscription
	// is valid.
	orderEventChan := make(chan []*zeroex.OrderEvent, 1)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventChan)
	require.NoError(t, err)
	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

	// Create a valid order and send it to the rpc client's "AddOrders" endpoint.
	signedTestOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
	time.Sleep(500 * time.Millisecond)
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")
	_, err = client.AddOrders([]*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)

	// Ensure that the "AddOrders" request triggered an order event that was
	// passed through the subscription.
	orderEvents := <-orderEventChan
	signedTestOrder.ResetHash()
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	assert.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, expectedOrderHash, orderEvent.OrderHash)
	assert.Equal(t, signedTestOrder, orderEvent.SignedOrder)
	assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)
	assert.Equal(t, expectedFillableTakerAssetAmount, orderEvent.FillableTakerAssetAmount)
	assert.Equal(t, []*zeroex.ContractEvent{}, orderEvent.ContractEvents)
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
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	// Wait for the rpc server to start and then start the rpc client
	_, err := waitForLogSubstring(ctx, logMessages, "started WS RPC server")
	require.NoError(t, err, "WS RPC server didn't start")
	client, err := rpc.NewClient(standaloneWSRPCEndpointPrefix + strconv.Itoa(wsRPCPort+count))
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
