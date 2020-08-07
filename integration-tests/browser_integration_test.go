// +build !js

// Package integrationtests contains broad integration tests that
// include a bootstrap node, a standalone node, and a browser node.
package integrationtests

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	gqlclient "github.com/0xProject/0x-mesh/graphql/client"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrowserIntegration(t *testing.T) {
	if !browserIntegrationTestsEnabled {
		t.Skip("Browser integration tests are disabled. You can enable them with the --enable-browser-integration-tests flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	// Declare a context that will be used for all child processes, servers, and
	// other goroutines.
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	ctx, _ = chromedp.NewContext(ctx, chromedp.WithErrorf(t.Errorf))
	defer cancel()

	removeOldFiles(t, ctx)
	buildForTests(t, ctx)

	// wg is a WaitGroup for the entire tests. We won't exit until wg is done.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		startBootstrapNode(t, ctx)
	}()

	// standaloneLogMessages is a channel through which log messages from the
	// standalone node will be sent. We use a large buffer so it doesn't cause
	// goroutines to block.
	standaloneLogMessages := make(chan string, 1024)
	count := int(atomic.AddInt32(&nodeCount, 1))
	customOrderFilter := `{"properties": { "makerAddress": { "const": "0x6ecbe1db9ef729cbe972c83fb886247691fb6beb" }}}`

	// Start the standalone node in a goroutine.
	// Note(albrow) we need to use a specific data dir because we need to use the same private key for each test.
	// The tests themselves are written in a way that depend on this.
	wg.Add(1)
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, browserIntegrationTestDataDir, customOrderFilter, standaloneLogMessages)
	}()

	// standaloneOrder is an order that will be sent to the network by the
	// standalone node.
	standaloneOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))

	// We also need to set up the maker state for the order that will be created in the browser (we don't care
	// if this order exactly matches the one created in the browser, we just care about makerAddress,
	// makerAssetData, and makerAssetAmount).
	scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAddress(constants.GanacheAccount1),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
		orderopts.MakerAssetAmount(big.NewInt(1000)),
	)

	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
	time.Sleep(500 * time.Millisecond)
	standaloneOrderHash, err := standaloneOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	// In a separate goroutine, send standaloneOrder through the GraphQL API for
	// the standalone node.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for the GraphQL server to start before sending the order.
		_, err := waitForLogSubstring(ctx, standaloneLogMessages, "starting GraphQL server")
		require.NoError(t, err, "GraphQL server didn't start")
		graphQLClient := gqlclient.New(graphQLServerURL)
		require.NoError(t, err)
		results, err := graphQLClient.AddOrders(ctx, []*zeroex.SignedOrder{standaloneOrder})
		require.NoError(t, err)
		assert.Len(t, results.Accepted, 1, "Expected 1 order to be accepted via GraphQL API")
		assert.Len(t, results.Rejected, 0, "Expected 0 orders to be rejected via GraphQL API")
	}()

	// Start a simple HTTP server to serve the web page for the browser node.
	ts := httptest.NewServer(http.FileServer(http.Dir("../packages/mesh-integration-tests/dist")))
	defer ts.Close()

	// browserLogMessages is a channel through which log messages from the
	// standalone node will be sent. We use a large buffer so it doesn't cause
	// goroutines to block.
	browserLogMessages := make(chan string, 1024)

	// Start the browser node.
	wg.Add(1)
	go func() {
		defer wg.Done()
		startBrowserNode(t, ctx, ts.URL, browserLogMessages)
	}()

	// browserPeerIDChan is used to retrieve the peer ID of the browser node.
	// Unlike the other nodes, we can't know it ahead of time because we have no
	// easy way to manipulate localStorage.
	browserPeerIDChan := make(chan string, 1)

	// browserOrderHashChan is used to retrieve the order hash of the order signed
	// in the browser node.
	browserOrderHashChan := make(chan string, 1)

	// messageWG is a separate WaitGroup that will be used to wait for all
	// expected messages to be logged.
	messageWG := &sync.WaitGroup{}

	// Start a goroutine to wait for the log messages we expect from the browser
	// node.
	messageWG.Add(1)
	go func() {
		defer messageWG.Done()

		// Wait for the order hash to be logged.
		msg, err := waitForSignedOrderLog(ctx, browserLogMessages)
		assert.NoError(t, err, "Could not find browser orderHash in logs. Maybe the browser node didn't start?")
		browserOrderHash, err := extractOrderHashFromLog(msg)
		assert.NoError(t, err, "Could not extract brower orderHash from log message.")
		fmt.Println("browser order hash is", browserOrderHash)
		browserOrderHashChan <- browserOrderHash

		// Wait for the peer ID to be logged first.
		msg, err = waitForLogSubstring(ctx, browserLogMessages, "myPeerID")
		assert.NoError(t, err, "Could not find browser peer ID in logs. Maybe the browser node didn't start?")
		browserPeerID, err := extractPeerIDFromLog(msg)
		assert.NoError(t, err, "Could not extract brower peer ID from log message.")
		fmt.Println("browser peer ID is", browserPeerID)
		browserPeerIDChan <- browserPeerID

		// Next, wait for the order to be received.
		expectedOrderEventLog := orderEventLog{
			OrderHash: standaloneOrderHash.Hex(),
			EndState:  "ADDED",
		}
		_, err = waitForOrderEventLog(ctx, browserLogMessages, expectedOrderEventLog)
		assert.NoError(t, err, "Browser node did not receive order sent by standalone node")
	}()

	// Start a goroutine to wait for the log messages we expect from the
	// standalone node.
	messageWG.Add(1)
	go func() {
		defer messageWG.Done()
		browserOrderHash := <-browserOrderHashChan
		browserPeerID := <-browserPeerIDChan
		expectedOrderLog := receivedOrderLog{
			OrderHash: browserOrderHash,
			From:      browserPeerID,
		}
		_, err := waitForReceivedOrderLog(ctx, standaloneLogMessages, expectedOrderLog)
		assert.NoError(t, err, "Standalone node did not receive order sent by browser node")
	}()

	// Wait for all expected messages to be logged.
	messageWG.Wait()

	// Cancel the context and wait for all outstanding goroutines to finish.
	cancel()
	wg.Wait()
}
