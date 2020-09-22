// +build !js

package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
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
	wg := &sync.WaitGroup{}
	client, _ := buildAndStartGraphQLServer(t, ctx, wg)

	// Create a new valid order.
	signedTestOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	time.Sleep(blockProcessingWaitTime)

	// Send the "AddOrders" request to the GraphQL server.
	validationResponse, err := client.AddOrders(ctx, []*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)

	// Ensure that the validation results contain only the order that was
	// sent to the GraphQL server and that the order was marked as valid.
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

func TestGetOrder(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	client, _ := buildAndStartGraphQLServer(t, ctx, wg)

	orderOptions := orderopts.SetupMakerState(true)
	signedTestOrder := scenario.NewSignedTestOrder(t, orderOptions)
	time.Sleep(blockProcessingWaitTime)

	validationResponse, err := client.AddOrders(ctx, []*zeroex.SignedOrder{signedTestOrder})
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	expectedHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	expectedOrder := &gqlclient.OrderWithMetadata{
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
		Hash:                     expectedHash,
		FillableTakerAssetAmount: signedTestOrder.TakerAssetAmount,
	}
	actualOrder, err := client.GetOrder(ctx, expectedHash)
	require.NoError(t, err)
	require.Equal(t, expectedOrder, actualOrder)

	cancel()
	wg.Wait()
}

func TestFindOrders(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	client, _ := buildAndStartGraphQLServer(t, ctx, wg)

	// Create 10 new valid orders.
	numOrders := 10
	orderOptions := scenario.OptionsForAll(orderopts.SetupMakerState(true))
	signedTestOrders := scenario.NewSignedTestOrdersBatch(t, numOrders, orderOptions)
	time.Sleep(blockProcessingWaitTime)

	// Send the newly created order to "AddOrders." The order is valid, and this should
	// be reflected in the validation results.
	validationResponse, err := client.AddOrders(ctx, signedTestOrders)
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, numOrders)
	assert.Len(t, validationResponse.Rejected, 0)

	// Get orders without any options.
	actualOrders, err := client.FindOrders(ctx)
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
	opts := gqlclient.FindOrdersOpts{
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
	actualOrdersWithOpts, err := client.FindOrders(ctx, opts)
	require.NoError(t, err)
	require.Len(t, actualOrdersWithOpts, 5)
	sortOrdersByHashDesc(expectedOrders)
	expectedOrdersWithOpts := expectedOrders[:5]
	assertOrderSlicesAreEqual(t, expectedOrdersWithOpts, actualOrdersWithOpts)

	cancel()
	wg.Wait()
}

func TestGetStats(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	client, peerID := buildAndStartGraphQLServer(t, ctx, wg)

	actualStats, err := client.GetStats(ctx)
	require.NoError(t, err)

	// Ensure that the "LatestBlock" in the stats response is non-nil and has a nonzero block number.
	assert.NotNil(t, actualStats.LatestBlock)
	assert.True(t, actualStats.LatestBlock.Number.String() != "0", "stats.LatestBlock.Number should not be 0")
	assert.NotEmpty(t, actualStats.LatestBlock.Hash, "stats.LatestBlock.Hash should not be empty")
	assert.NotEqual(t, actualStats.Version, "")
	actualStats.Version = ""
	expectedStats := &gqlclient.Stats{
		PubSubTopic:     "/0x-orders/version/3/chain/1337/schema/e30=",
		Rendezvous:      "/0x-mesh/network/1337/version/2",
		PeerID:          peerID,
		EthereumChainID: 1337,
		// NOTE(jalextowle): Since this test uses an actual mesh node, we can't know in advance which block
		//                   should be the latest block.
		LatestBlock:                       actualStats.LatestBlock,
		NumOrders:                         0,
		NumOrdersIncludingRemoved:         0,
		NumPeers:                          0,
		MaxExpirationTime:                 constants.UnlimitedExpirationTime,
		StartOfCurrentUTCDay:              ratelimit.GetUTCMidnightOfDate(time.Now()),
		EthRPCRequestsSentInCurrentUTCDay: 0,
		EthRPCRateLimitExpiredRequests:    0,
	}
	assert.Equal(t, expectedStats, actualStats)

	cancel()
	wg.Wait()
}

type rawResponse struct {
	Stats statsWithJustNumOrders `json:"stats"`
}

type statsWithJustNumOrders struct {
	NumOrders int `json:"numOrders"`
}

func TestRawQuery(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	client, _ := buildAndStartGraphQLServer(t, ctx, wg)

	query := `{
		stats {
			numOrders
		}
	}`
	var actualResponse rawResponse
	require.NoError(t, client.RawQuery(ctx, query, &actualResponse))
	expectedResponse := rawResponse{
		Stats: statsWithJustNumOrders{
			NumOrders: 0,
		},
	}
	require.Equal(t, expectedResponse, actualResponse)

	cancel()
	wg.Wait()
}

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

func buildAndStartGraphQLServer(t *testing.T, ctx context.Context, wg *sync.WaitGroup) (client *gqlclient.Client, peerID string) {
	removeOldFiles(t)
	buildStandaloneForTests(t, ctx)

	// Start a standalone node with a wait group that is completed when the goroutine completes.
	wg.Add(1)
	logMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, "", "", logMessages)
	}()

	// Wait for the GraphQL server to start and extract the PeerID from the log.
	var jsonLog struct {
		PeerID string `json:"myPeerID"`
	}
	log, err := waitForLogSubstring(ctx, logMessages, "starting GraphQL server")
	require.NoError(t, err, "GraphQL server didn't start")
	err = json.Unmarshal([]byte(log), &jsonLog)
	require.NoError(t, err)

	time.Sleep(serverStartWaitTime)
	return gqlclient.New(graphQLServerURL), jsonLog.PeerID
}
