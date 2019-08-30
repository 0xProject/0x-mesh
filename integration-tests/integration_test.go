// Package integrationtests contains broad integration integrationtests that
// include a bootstrap node, a standalone node, and a browser node.
package integrationtests

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ethereumRPCURL    = "http://localhost:8545"
	ethereumNetworkID = 50

	// Various config options/information for the bootstrap node. The private key
	// for the bootstrap node is checked in to version control so we know it's
	// peer ID ahead of time.
	bootstrapAddr    = "/ip4/127.0.0.1/tcp/60500/ws"
	bootstrapPeerID  = "16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapList    = "/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapDataDir = "./data/bootstrap-0"

	// Various config options/information for the standalone node. Like the
	// bootstrap node, we know the private key/peer ID ahead of time.
	standalonePeerID      = "16Uiu2HAmM9j68mgGGSFkXsuzbGJA8ezVHtQ2H9y6mRJAPhx6xtj9"
	standaloneDataDir     = "./data/standalone-0"
	standaloneRPCEndpoint = "ws://localhost:60501"
	standaloneRPCPort     = 60501

	// Hash of the order that will be sent by the browser node.
	expectedBrowserOrderHash = "0x7292f6e7bee79f117c146c57f207d6a380e888b871ef733ae2608a064c36ef83"
	// Hash of the order that will be sent by the standalone node.
	expectedStandaloneOrderHash = "0x4f43d2126b856ed72e40cd504ea4e6cea1c88cd1adc1eb2ea8c30da412470584"
)

// standaloneOrder is the order that will be sent to the Mesh network by the
// standalone node.
var standaloneOrder = &zeroex.SignedOrder{
	Order: zeroex.Order{
		MakerAddress:          common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631"),
		MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		MakerAssetAmount:      big.NewInt(1000),
		MakerFee:              big.NewInt(0),
		TakerAddress:          common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		TakerAssetAmount:      big.NewInt(2000),
		TakerFee:              big.NewInt(0),
		SenderAddress:         common.HexToAddress("0x0000000000000000000000000000000000000000"),
		ExchangeAddress:       common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		ExpirationTimeSeconds: big.NewInt(1567121010),
		Salt:                  big.NewInt(1548619145450),
	},
	Signature: common.Hex2Bytes("1b15b0edc1cab84e1be2a801cef16cb6da2edc1f17cc3239ff5ebf2c84de8bac7854005c7d85a622732177c7abe69545254a564fcf60e57b21fbdf6cd7ade9078c03"),
}

// Since the tests take so long, we don't want them to run as part of the normal
// testing process. They will only be run if the "--integration" flag is used.
var integrationTestsEnabled bool

func init() {
	flag.BoolVar(&integrationTestsEnabled, "integration", false, "enable integration tests")
	flag.Parse()
}

func TestBrowserIntegration(t *testing.T) {
	if !integrationTestsEnabled {
		t.Skip("Integration tests are disabled. You can enable them with the --integration")
	}

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

	// Start the standalone node in a goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		startStandaloneNode(t, ctx, standaloneLogMessages)
	}()

	// In a separate goroutine, send standaloneOrder through the RPC endpoint for
	// the standalone node.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for the RPC server to start before sending the order.
		_, err := waitForLogSubstring(ctx, standaloneLogMessages, "started RPC server")
		require.NoError(t, err, "RPC server didn't start")
		rpcClient, err := rpc.NewClient(standaloneRPCEndpoint)
		require.NoError(t, err)
		results, err := rpcClient.AddOrders([]*zeroex.SignedOrder{standaloneOrder})
		require.NoError(t, err)
		assert.Len(t, results.Accepted, 1, "Expected 1 order to be accepted over RPC")
		assert.Len(t, results.Rejected, 0, "Expected 0 orders to be rejected over RPC")
	}()

	// Start a sinple HTTP server to serve the web page for the browser node.
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

	// browserPeerIDChan is used to retrive the peer ID of the browser nodes.
	// Unlike the other nodes, we can't know it ahead of time because we have no
	// easy way to manipulate localStorage.
	browserPeerIDChan := make(chan string, 1)

	// messageWG is a separate WaitGroup that will be used to wait for all
	// expected messages to be logged.
	messageWG := &sync.WaitGroup{}

	// Start a goroutine to wait for the log messages we expect from the browser
	// node.
	messageWG.Add(1)
	go func() {
		defer messageWG.Done()

		// Wait for the peer ID to be logged first.
		msg, err := waitForLogSubstring(ctx, browserLogMessages, "myPeerID")
		assert.NoError(t, err, "Could not find browser peer ID in logs. Maybe the browser node didn't start?")
		browserPeerID, err := extractPeerIDFromLog(msg)
		assert.NoError(t, err, "Could not extract brower peer ID from log message.")
		fmt.Println("browser peerID is", browserPeerID)
		browserPeerIDChan <- browserPeerID

		// Next, wait for the order to be received.
		expectedOrderEventLog := orderEventLog{
			OrderHash: expectedStandaloneOrderHash,
			Kind:      "ADDED",
		}
		_, err = waitForOrderEventLog(ctx, browserLogMessages, expectedOrderEventLog)
		assert.NoError(t, err, "Browser node did not receive order sent by standalone node")
	}()

	// Start a goroutine to wait for the log messages we expect from the
	// standalone node.
	messageWG.Add(1)
	go func() {
		defer messageWG.Done()
		browserPeerID := <-browserPeerIDChan
		expectedOrderLog := receivedOrderLog{
			OrderHash: expectedBrowserOrderHash,
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

func removeOldFiles(t *testing.T, ctx context.Context) {
	fmt.Println("Removing old files...")
	require.NoError(t, os.RemoveAll(filepath.Join(standaloneDataDir, "db")))
	require.NoError(t, os.RemoveAll(filepath.Join(standaloneDataDir, "p2p")))
	require.NoError(t, os.RemoveAll(filepath.Join(bootstrapDataDir, "p2p")))
}

func buildForTests(t *testing.T, ctx context.Context) {
	fmt.Println("Building mesh executable...")
	cmd := exec.CommandContext(ctx, "go", "install", ".")
	cmd.Dir = "../cmd/mesh"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "could not build mesh: %s", string(output))

	fmt.Println("Building mesh-bootstrap executable...")
	cmd = exec.CommandContext(ctx, "go", "install", ".")
	cmd.Dir = "../cmd/mesh-bootstrap"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not build mesh-bootstrap: %s", string(output))

	fmt.Println("Installing dependencies for TypeScript bindings...")
	cmd = exec.CommandContext(ctx, "yarn", "install")
	cmd.Dir = "../browser"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not install depedencies for TypeScript bindings: %s", string(output))

	fmt.Println("Building TypeScript bindings...")
	cmd = exec.CommandContext(ctx, "yarn", "build")
	cmd.Dir = "../browser"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not build TypeScript bindings: %s", string(output))

	fmt.Println("Installing dependencies for browser node...")
	cmd = exec.CommandContext(ctx, "yarn", "install", "--force")
	cmd.Dir = "./browser"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not install yarn depedencies: %s", string(output))

	fmt.Println("Running postinstall for browser node...")
	cmd = exec.CommandContext(ctx, "yarn", "postinstall")
	cmd.Dir = "./browser"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not run yarn postinstall: %s", string(output))

	fmt.Println("Building browser node...")
	cmd = exec.CommandContext(ctx, "yarn", "build")
	cmd.Dir = "./browser"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not build browser node: %s", string(output))
	fmt.Println("Done building everything")
}

func startBootstrapNode(t *testing.T, ctx context.Context) {
	cmd := exec.CommandContext(ctx, "mesh-bootstrap")
	cmd.Env = append(
		os.Environ(),
		"P2P_BIND_ADDRS="+bootstrapAddr,
		"P2P_ADVERTISE_ADDRS="+bootstrapAddr,
		"DATA_DIR="+bootstrapDataDir,
		"BOOTSTRAP_LIST="+bootstrapList,
	)
	output, err := cmd.CombinedOutput()
	// Note(albrow): unfortunately we can't get the underlying signal that
	// caused the process to exit. We can only compare the error string.
	if err.Error() == "signal: killed" {
		// If the command was killed, that just means the context was cancelled
		// and the test is over.
		return
	}
	assert.NoError(t, err, "could not run bootstrap node: %s", string(output))
}

func startStandaloneNode(t *testing.T, ctx context.Context, logMessages chan<- string) {
	cmd := exec.CommandContext(ctx, "mesh")
	cmd.Env = append(
		os.Environ(),
		"VERBOSITY=5",
		"DATA_DIR="+standaloneDataDir,
		"BOOTSTRAP_LIST="+bootstrapList,
		"ETHEREUM_RPC_URL="+ethereumRPCURL,
		"ETHEREUM_NETWORK_ID="+strconv.Itoa(ethereumNetworkID),
		"RPC_PORT="+strconv.Itoa(standaloneRPCPort),
	)

	// Pipe messages from stderr through the logMessages channel.
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	scanner := bufio.NewScanner(stderr)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			fmt.Println("[standalone]: " + scanner.Text())
			logMessages <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
	}()

	require.NoError(t, cmd.Start())
	err = cmd.Wait()
	if err, ok := err.(*exec.ExitError); ok {
		// Note(albrow): unfortunately we can't get the underlying signal that
		// caused the process to exit. We can only compare the error string.
		if err.Error() == "signal: killed" {
			// If the command was killed, that just means the context was cancelled
			// and the test is over.
			return
		}
	}
	assert.NoError(t, err, "could not run standalone node: %s", err)
	wg.Wait()
}

func startBrowserNode(t *testing.T, ctx context.Context, url string, browserLogMessages chan<- string) {
	// Use chromedp to visit the web page for the browser node.
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			switch ev.Type {
			case runtime.APITypeLog:
				// Send console.log events through the channel.
				for _, arg := range ev.Args {
					if arg.Type == runtime.TypeString {
						fmt.Println("[browser]: " + string(arg.Value))
						browserLogMessages <- string(arg.Value)
					}
				}
			case runtime.APITypeError:
				// Report any console.error events as test failures.
				for _, arg := range ev.Args {
					t.Errorf("JavaScript console error: (%s) %s", arg.Type, arg.Value)
				}
			}
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// The #jsFinished element is one specifically created to indicate that the
		// JavaScript code is done running.
		chromedp.WaitVisible("#jsFinished", chromedp.ByID),
	); err != nil && err != context.Canceled {
		t.Error(err)
	}
}

// waitForLogMessage blocks until a message is logged that psses the given
// filter or the context is done. If the message is logged before the context is
// done, it will return the entire message. Otherwise it returns an error.
func waitForLogMessage(ctx context.Context, logMessages <-chan string, filter func(string) bool) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("(%s) timed out waiting for message", ctx.Err())
		case msg := <-logMessages:
			if filter(msg) {
				return msg, nil
			}
		}
	}
}

func waitForLogSubstring(ctx context.Context, logMessages <-chan string, substr string) (string, error) {
	return waitForLogMessage(ctx, logMessages, func(msg string) bool {
		return strings.Contains(msg, substr)
	})
}

// A holder type used for parsing the "received new valid order from peer"
// messages that are logged by Mesh when an order is received.
type receivedOrderLog struct {
	OrderHash string `json:"orderHash_string"`
	From      string `json:"from_string"`
}

func waitForReceivedOrderLog(ctx context.Context, logMessages <-chan string, expectedLog receivedOrderLog) (string, error) {
	return waitForLogMessage(ctx, logMessages, func(msg string) bool {
		var foundLog receivedOrderLog
		if err := unquoteAndUnmarshal(msg, &foundLog); err != nil {
			return false
		}
		return foundLog.OrderHash == expectedLog.OrderHash &&
			foundLog.From == expectedLog.From
	})
}

// A holder type for parsing logged OrderEvents. These are received by either
// an RPC subscription or in the TypeScript bindings and are not usually logged
// by Mesh. They need to be explicitly logged.
type orderEventLog struct {
	OrderHash string `json:"orderHash"`
	Kind      string `json:"kind"`
}

func waitForOrderEventLog(ctx context.Context, logMessages <-chan string, expectedLog orderEventLog) (string, error) {
	return waitForLogMessage(ctx, logMessages, func(msg string) bool {
		var foundLog orderEventLog
		if err := unquoteAndUnmarshal(msg, &foundLog); err != nil {
			return false
		}
		return foundLog.OrderHash == expectedLog.OrderHash &&
			foundLog.Kind == expectedLog.Kind
	})
}

func unquoteAndUnmarshal(msg string, holder interface{}) error {
	// Depending on the environment, the message may contain escaped quotes
	// which we need to unescape.
	unquoted, err := strconv.Unquote(msg)
	if err == nil {
		msg = unquoted
	}
	if err := json.Unmarshal([]byte(msg), holder); err != nil {
		return err
	}
	return nil
}

// extractPeerIDFromLog expects a log message that contains a peer ID under the
// JSON field "myPeerID". If the given msg is the correct format, it extracts
// and returns the peerID.
func extractPeerIDFromLog(msg string) (string, error) {
	unquoted, err := strconv.Unquote(msg)
	if err == nil {
		msg = unquoted
	}
	holder := struct {
		PeerID string `json:"myPeerID"`
	}{}
	if err := json.Unmarshal([]byte(msg), &holder); err != nil {
		return "", err
	}
	return holder.PeerID, nil
}
