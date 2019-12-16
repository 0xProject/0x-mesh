// Package integrationtests contains broad integration tests that
// include a bootstrap node, a standalone node, and a browser node.
package integrationtests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/chromedp/chromedp"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// Start the standalone node in a goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, count, standaloneLogMessages)
	}()

	// standaloneOrder is an order that will be sent to the network by the
	// standalone node.
	ethClient := ethclient.NewClient(ethRPCClient)
	standaloneOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	standaloneOrderHash, err := standaloneOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")

	// In a separate goroutine, send standaloneOrder through the RPC endpoint for
	// the standalone node.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for the RPC server to start before sending the order.
		_, err := waitForLogSubstring(ctx, standaloneLogMessages, "started RPC server")
		require.NoError(t, err, "RPC server didn't start")
		rpcClient, err := rpc.NewClient(standaloneRPCEndpointPrefix + strconv.Itoa(rpcPort+count))
		require.NoError(t, err)
		results, err := rpcClient.AddOrders([]*zeroex.SignedOrder{standaloneOrder})
		require.NoError(t, err)
		assert.Len(t, results.Accepted, 1, "Expected 1 order to be accepted over RPC")
		assert.Len(t, results.Rejected, 0, "Expected 0 orders to be rejected over RPC")
	}()

	// Start a simple HTTP server to serve the web page for the browser node.
	ts := httptest.NewServer(http.FileServer(http.Dir("./browser/dist")))
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
