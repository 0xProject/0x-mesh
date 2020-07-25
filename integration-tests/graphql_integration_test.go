// +build !js

package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/graphql/client"
	gqlclient "github.com/0xProject/0x-mesh/graphql/client"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
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
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	// Wait until the rpc server has been started, and then create an rpc client
	// that connects to the rpc server.
	_, err := waitForLogSubstring(ctx, logMessages, fmt.Sprintf("starting GraphQL server"))
	require.NoError(t, err, fmt.Sprintf("GraphQL server didn't start"))
	client := gqlclient.New(graphQLServerURL)

	// Create a new valid order.
	signedTestOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
	time.Sleep(500 * time.Millisecond)

	// Send the "AddOrders" request to the rpc server.
	validationResponse, err := client.AddOrders(ctx, []*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)

	// Ensure that the validation results contain only the order that was
	// sent to the rpc server and that the order was marked as valid.
	require.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)
	accepted := validationResponse.Accepted[0]
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")
	expectedAcceptedOrder := &gqlclient.OrderWithMetadata{
		ChainID:                  signedTestOrder.ChainID,
		ExchangeAddress:          signedTestOrder.ExchangeAddress,
		MakerAddress:             signedTestOrder.MakerAddress,
		MakerAssetData:           signedTestOrder.MakerAssetData,
		MakerAssetAmount:         signedTestOrder.MakerAssetAmount,
		MakerFeeAssetData:        signedTestOrder.MakerFeeAssetData,
		MakerFee:                 signedTestOrder.MakerFee,
		TakerAddress:             signedTestOrder.TakerAddress,
		TakerAssetData:           signedTestOrder.TakerAssetData,
		TakerAssetAmount:         signedTestOrder.TakerAssetAmount,
		TakerFeeAssetData:        signedTestOrder.TakerFeeAssetData,
		TakerFee:                 signedTestOrder.TakerFee,
		SenderAddress:            signedTestOrder.SenderAddress,
		FeeRecipientAddress:      signedTestOrder.FeeRecipientAddress,
		ExpirationTimeSeconds:    signedTestOrder.ExpirationTimeSeconds,
		Salt:                     signedTestOrder.Salt,
		Signature:                signedTestOrder.Signature,
		Hash:                     expectedOrderHash,
		FillableTakerAssetAmount: expectedFillableTakerAssetAmount,
	}
	assert.Equal(t, expectedAcceptedOrder, accepted.Order, "accepted.Order")
	assert.Equal(t, true, accepted.IsNew, "accepted.IsNew")

	cancel()
	wg.Wait()
}

func TestGetOrders(t *testing.T) {
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
	_, err := waitForLogSubstring(ctx, logMessages, fmt.Sprintf("starting GraphQL server"))
	require.NoError(t, err, fmt.Sprintf("GraphQL server didn't start"))
	client := client.New(graphQLServerURL)

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
	validationResponse, err := client.AddOrders(ctx, signedTestOrders)
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, numOrders)
	assert.Len(t, validationResponse.Rejected, 0)

	// Get orders without any options.
	actualOrders, err := client.GetOrders(ctx)
	require.NoError(t, err)
	require.Len(t, actualOrders, 10)
	expectedOrders := make([]*gqlclient.OrderWithMetadata, len(signedTestOrders))
	for i, signedOrder := range signedTestOrders {
		hash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err)
		expectedOrders[i] = &gqlclient.OrderWithMetadata{
			ChainID:                  signedOrder.ChainID,
			ExchangeAddress:          signedOrder.ExchangeAddress,
			MakerAddress:             signedOrder.MakerAddress,
			MakerAssetData:           signedOrder.MakerAssetData,
			MakerAssetAmount:         signedOrder.MakerAssetAmount,
			MakerFeeAssetData:        signedOrder.MakerFeeAssetData,
			MakerFee:                 signedOrder.MakerFee,
			TakerAddress:             signedOrder.TakerAddress,
			TakerAssetData:           signedOrder.TakerAssetData,
			TakerAssetAmount:         signedOrder.TakerAssetAmount,
			TakerFeeAssetData:        signedOrder.TakerFeeAssetData,
			TakerFee:                 signedOrder.TakerFee,
			SenderAddress:            signedOrder.SenderAddress,
			FeeRecipientAddress:      signedOrder.FeeRecipientAddress,
			ExpirationTimeSeconds:    signedOrder.ExpirationTimeSeconds,
			Salt:                     signedOrder.Salt,
			Signature:                signedOrder.Signature,
			Hash:                     hash,
			FillableTakerAssetAmount: signedOrder.TakerAssetAmount,
		}
	}
	assertOrdersAreUnsortedEqual(t, expectedOrders, actualOrders)

	// Get orders with filter, sort, and limit.
	opts := gqlclient.GetOrdersOpts{
		Filters: []gqlclient.OrderFilter{
			{
				Field: gqlclient.OrderFieldChainID,
				Kind:  gqlclient.FilterKindEqual,
				Value: signedTestOrders[0].ChainID,
			},
			{
				Field: gqlclient.OrderFieldExpirationTimeSeconds,
				Kind:  gqlclient.FilterKindGreaterOrEqual,
				Value: big.NewInt(0),
			},
		},
		Sort: []gqlclient.OrderSort{
			{
				Field:     gqlclient.OrderFieldHash,
				Direction: gqlclient.SortDirectionDesc,
			},
		},
		Limit: 5,
	}
	actualOrdersWithOpts, err := client.GetOrders(ctx, opts)
	require.NoError(t, err)
	require.Len(t, actualOrdersWithOpts, 5)
	sortOrdersByHashDesc(expectedOrders)
	expectedOrdersWithOpts := expectedOrders[:5]
	assertOrderSlicesAreEqual(t, expectedOrdersWithOpts, actualOrdersWithOpts)

	cancel()
	wg.Wait()
}

// func TestWSGetStats(t *testing.T) {
// 	runGetStatsTest(t, standaloneWSRPCEndpointPrefix, "WS", wsRPCPort)
// }

// func TestHTTPGetStats(t *testing.T) {
// 	runGetStatsTest(t, standaloneHTTPRPCEndpointPrefix, "HTTP", httpRPCPort)
// }

// func runGetStatsTest(t *testing.T, rpcEndpointPrefix, rpcServerType string, rpcPort int) {
// 	teardownSubTest := setupSubTest(t)
// 	defer teardownSubTest(t)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	removeOldFiles(t, ctx)
// 	buildStandaloneForTests(t, ctx)

// 	// Start a standalone node with a wait group that is completed when the goroutine completes.
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)
// 	logMessages := make(chan string, 1024)
// 	count := int(atomic.AddInt32(&nodeCount, 1))
// 	go func() {
// 		defer wg.Done()
// 		startStandaloneNode(t, ctx, count, "", "", logMessages)
// 	}()

// 	// Wait for the rpc server to start and get the peer ID of the node. Start the
// 	// rpc client after the server has been started,
// 	var jsonLog struct {
// 		PeerID string `json:"myPeerID"`
// 	}
// 	log, err := waitForLogSubstring(ctx, logMessages, fmt.Sprintf("started %s RPC server", rpcServerType))
// 	require.NoError(t, err, fmt.Sprintf("%s RPC server didn't start", rpcServerType))
// 	err = json.Unmarshal([]byte(log), &jsonLog)
// 	require.NoError(t, err)
// 	client, err := rpc.NewClient(rpcEndpointPrefix + strconv.Itoa(rpcPort+count))
// 	require.NoError(t, err)

// 	getStatsResponse, err := client.GetStats()
// 	require.NoError(t, err)

// 	// Ensure that the "LatestBlock" in the stats response is non-nil and has a nonzero block number.
// 	assert.NotNil(t, getStatsResponse.LatestBlock)
// 	assert.True(t, getStatsResponse.LatestBlock.Number != "")

// 	// NOTE(jalextowle): Since this test uses an actual mesh node, we can't know in advance which block
// 	//                   should be the latest block.
// 	getStatsResponse.LatestBlock = types.LatestBlock{}

// 	// Ensure that the correct response was logged by "GetStats"
// 	require.Equal(t, "/0x-orders/version/3/chain/1337/schema/e30=", getStatsResponse.PubSubTopic, "PubSubTopic")
// 	require.Equal(t, "/0x-mesh/network/1337/version/2", getStatsResponse.Rendezvous, "Rendezvous")
// 	require.Equal(t, []string{}, getStatsResponse.SecondaryRendezvous, "SecondaryRendezvous")
// 	require.Equal(t, jsonLog.PeerID, getStatsResponse.PeerID, "PeerID")
// 	require.Equal(t, 1337, getStatsResponse.EthereumChainID, "EthereumChainID")
// 	require.Equal(t, types.LatestBlock{}, getStatsResponse.LatestBlock, "LatestBlock")
// 	require.Equal(t, 0, getStatsResponse.NumOrders, "NumOrders")
// 	require.Equal(t, 0, getStatsResponse.NumPeers, "NumPeers")
// 	require.Equal(t, constants.UnlimitedExpirationTime.String(), getStatsResponse.MaxExpirationTime, "MaxExpirationTime")
// 	require.Equal(t, ratelimit.GetUTCMidnightOfDate(time.Now()), getStatsResponse.StartOfCurrentUTCDay, "StartOfCurrentDay")

// 	cancel()
// 	wg.Wait()
// }

// func TestOrdersSubscription(t *testing.T) {
// 	teardownSubTest := setupSubTest(t)
// 	defer teardownSubTest(t)

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	removeOldFiles(t, ctx)
// 	buildStandaloneForTests(t, ctx)

// 	// Start a standalone node with a wait group that is completed when the goroutine completes.
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)
// 	logMessages := make(chan string, 1024)
// 	count := int(atomic.AddInt32(&nodeCount, 1))
// 	go func() {
// 		defer wg.Done()
// 		startStandaloneNode(t, ctx, count, "", "", logMessages)
// 	}()

// 	// Wait for the rpc server to start and then start the rpc client.
// 	_, err := waitForLogSubstring(ctx, logMessages, "started WS RPC server")
// 	require.NoError(t, err, "WS RPC server didn't start")
// 	client, err := rpc.NewClient(standaloneWSRPCEndpointPrefix + strconv.Itoa(wsRPCPort+count))
// 	require.NoError(t, err)

// 	// Subscribe to order events through the rpc client and ensure that the subscription
// 	// is valid.
// 	orderEventChan := make(chan []*zeroex.OrderEvent, 1)
// 	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventChan)
// 	require.NoError(t, err)
// 	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

// 	// Create a valid order and send it to the rpc client's "AddOrders" endpoint.
// 	signedTestOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
// 	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
// 	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
// 	// in order for the order validation run at order submission to occur at a block number equal or higher then
// 	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
// 	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
// 	time.Sleep(500 * time.Millisecond)
// 	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
// 	require.NoError(t, err, "could not compute order hash for standalone order")
// 	_, err = client.AddOrders([]*zeroex.SignedOrder{signedTestOrder})
// 	require.NoError(t, err)

// 	// Ensure that the "AddOrders" request triggered an order event that was
// 	// passed through the subscription.
// 	orderEvents := <-orderEventChan
// 	signedTestOrder.ResetHash()
// 	expectedFillableTakerAssetAmount := signedTestOrder.TakerAssetAmount
// 	assert.Len(t, orderEvents, 1)
// 	orderEvent := orderEvents[0]
// 	assert.Equal(t, expectedOrderHash, orderEvent.OrderHash)
// 	assert.Equal(t, signedTestOrder, orderEvent.SignedOrder)
// 	assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)
// 	assert.Equal(t, expectedFillableTakerAssetAmount, orderEvent.FillableTakerAssetAmount)
// 	assert.Equal(t, []*zeroex.ContractEvent{}, orderEvent.ContractEvents)
// }

// func TestHeartbeatSubscription(t *testing.T) {
// 	teardownSubTest := setupSubTest(t)
// 	defer teardownSubTest(t)

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	removeOldFiles(t, ctx)
// 	buildStandaloneForTests(t, ctx)

// 	// Start a standalone node with a wait group that is completed when the goroutine completes.
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)
// 	logMessages := make(chan string, 1024)
// 	count := int(atomic.AddInt32(&nodeCount, 1))
// 	go func() {
// 		defer wg.Done()
// 		startStandaloneNode(t, ctx, count, "", "", logMessages)
// 	}()

// 	// Wait for the rpc server to start and then start the rpc client
// 	_, err := waitForLogSubstring(ctx, logMessages, "started WS RPC server")
// 	require.NoError(t, err, "WS RPC server didn't start")
// 	client, err := rpc.NewClient(standaloneWSRPCEndpointPrefix + strconv.Itoa(wsRPCPort+count))
// 	require.NoError(t, err)

// 	// Send the "SubscribeToHeartbeat" request through the rpc client and assert
// 	// that the subscription is not nil.
// 	heartbeatChan := make(chan string)
// 	clientSubscription, err := client.SubscribeToHeartbeat(ctx, heartbeatChan)
// 	defer clientSubscription.Unsubscribe()
// 	require.NoError(t, err)
// 	assert.NotNil(t, clientSubscription, "clientSubscription not nil")

// 	// Ensure that a valid heartbeat was received
// 	heartbeat := <-heartbeatChan
// 	assert.Equal(t, "tick", heartbeat)
// }

func assertOrdersAreUnsortedEqual(t *testing.T, expected, actual []*gqlclient.OrderWithMetadata) {
	// Make a copy of the given orders so we don't mess up the original when sorting them.
	expectedCopy := make([]*gqlclient.OrderWithMetadata, len(expected))
	copy(expectedCopy, expected)
	sortOrdersByHashAsc(expectedCopy)
	actualCopy := make([]*gqlclient.OrderWithMetadata, len(actual))
	copy(actualCopy, actual)
	sortOrdersByHashAsc(actualCopy)
	assertOrderSlicesAreEqual(t, expectedCopy, actualCopy)
}

func assertOrderSlicesAreEqual(t *testing.T, expected, actual []*gqlclient.OrderWithMetadata) {
	assert.Equal(t, len(expected), len(actual), "wrong number of orders")
	for i, expectedOrder := range expected {
		if i >= len(actual) {
			break
		}
		actualOrder := actual[i]
		assert.Equal(t, expectedOrder, actualOrder)
	}
	if t.Failed() {
		expectedJSON, err := json.MarshalIndent(expected, "", "  ")
		require.NoError(t, err)
		actualJSON, err := json.MarshalIndent(actual, "", "  ")
		require.NoError(t, err)
		t.Logf("\nexpected:\n%s\n\n", string(expectedJSON))
		t.Logf("\nactual:\n%s\n\n", string(actualJSON))
		assert.Equal(t, string(expectedJSON), string(actualJSON))
	}
}

func sortOrdersByHashAsc(orders []*gqlclient.OrderWithMetadata) {
	sort.SliceStable(orders, func(i, j int) bool {
		return bytes.Compare(orders[i].Hash.Bytes(), orders[j].Hash.Bytes()) == -1
	})
}

func sortOrdersByHashDesc(orders []*gqlclient.OrderWithMetadata) {
	sort.SliceStable(orders, func(i, j int) bool {
		return bytes.Compare(orders[i].Hash.Bytes(), orders[j].Hash.Bytes()) == 1
	})
}
