package integrationtests

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Since the tests take so long, we don't want them to run as part of the normal
// testing process. They will only be run if the "--integration" flag is used.
var (
	browserLegacyIntegrationTestsEnabled  bool
	browserGraphQLIntegrationTestsEnabled bool
)

var nodeCount int32

func init() {
	flag.BoolVar(&browserLegacyIntegrationTestsEnabled, "enable-browser-legacy-integration-tests", false, "enable browser legacy integration tests")
	flag.BoolVar(&browserGraphQLIntegrationTestsEnabled, "enable-browser-graphql-integration-tests", false, "enable browser graphql integration tests")
	testing.Init()
	flag.Parse()
}

var ethRPCClient *ethrpc.Client

func init() {
	var err error
	ethRPCClient, err = ethrpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(ethRPCClient)
	if err != nil {
		panic(err)
	}
}

func removeOldFiles(t *testing.T) {
	require.NoError(t, os.RemoveAll(filepath.Join(browserIntegrationTestDataDir, "sqlite-db")))
	require.NoError(t, os.RemoveAll(filepath.Join(browserIntegrationTestDataDir, "p2p")))
	require.NoError(t, os.RemoveAll(filepath.Join(bootstrapDataDir, "p2p")))
}

var hasRunBuildStandalone bool

func buildStandaloneForTests(t *testing.T, ctx context.Context) {
	// Skip building if this function has already been called.
	if hasRunBuildStandalone {
		return
	}
	hasRunBuildStandalone = true

	// Build the mesh executable
	fmt.Println("Building mesh executable...")
	cmd := exec.CommandContext(ctx, "go", "install", ".")
	cmd.Dir = "../cmd/mesh"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "could not build mesh: %s", string(output))
}

var hasRunBuildBootstrap bool

func buildBootstrapForTests(t *testing.T, ctx context.Context) {
	// Skip building if this function has already been called.
	if hasRunBuildBootstrap {
		return
	}
	hasRunBuildBootstrap = true

	// Build the bootstrap executable
	fmt.Println("Building mesh-bootstrap executable...")
	cmd := exec.CommandContext(ctx, "go", "install", ".")
	cmd.Dir = "../cmd/mesh-bootstrap"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "could not build mesh-bootstrap: %s", string(output))
}

var hasRunBuildAll bool

func buildForTests(t *testing.T, ctx context.Context) {
	// Skip building if this function has already been called.
	if hasRunBuildAll {
		return
	}
	hasRunBuildAll = true

	buildStandaloneForTests(t, ctx)
	buildBootstrapForTests(t, ctx)

	// Note(albrow): We have to rebuild the browser package manually in case
	// any Go code was changed. The TypeScript compiler can automatically rebuild
	// for TypeScript code changes only.
	fmt.Println("Building mesh-browser package...")
	cmd := exec.CommandContext(ctx, "yarn", "build")
	cmd.Dir = "../packages/mesh-browser"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "could not build mesh-browser package: %s", string(output))

	fmt.Println("Building integration-tests package...")
	cmd = exec.CommandContext(ctx, "yarn", "build")
	cmd.Dir = "../packages/mesh-integration-tests"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not build integration-tests package: %s", string(output))
	fmt.Println("Done building everything")
}

func startBootstrapNode(t *testing.T, ctx context.Context) {
	cmd := exec.CommandContext(ctx, "mesh-bootstrap")
	cmd.Env = append(
		os.Environ(),
		"P2P_BIND_ADDRS="+bootstrapAddr,
		"P2P_ADVERTISE_ADDRS="+bootstrapAddr,
		"LEVELDB_DATA_DIR="+bootstrapDataDir,
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

func startStandaloneNode(t *testing.T, ctx context.Context, nodeID int, dataDir string, customOrderFilter string, logMessages chan<- string) {
	cmd := exec.CommandContext(ctx, "mesh")
	if dataDir == "" {
		// If dataDir is empty. Set a default data dir to a file in the /tmp directory
		dataDir = filepath.Join("/tmp", "mesh_testing", uuid.New().String())
	}
	cmd.Env = append(
		os.Environ(),
		"VERBOSITY=6",
		"DATA_DIR="+dataDir,
		"BOOTSTRAP_LIST="+bootstrapList,
		"ETHEREUM_RPC_URL="+ethereumRPCURL,
		"ETHEREUM_CHAIN_ID="+strconv.Itoa(ethereumChainID),
		"ENABLE_GRAPHQL_SERVER=true",
		"GRAPHQL_SERVER_ADDR="+graphQLServerAddr,
		"BLOCK_POLLING_INTERVAL="+standaloneBlockPollingInterval,
		"ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC="+standaloneEthereumRPCMaxRequestsPer24HrUtc,
	)
	if customOrderFilter != "" {
		cmd.Env = append(cmd.Env, "CUSTOM_ORDER_FILTER="+customOrderFilter)
	}

	// Pipe messages from stderr through the logMessages channel.
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)

	scanner := bufio.NewScanner(stderr)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	errChan := make(chan error)

	go func() {
		defer wg.Done()

		for scanner.Scan() {
			fmt.Printf("[standalone %d]: %s\n", nodeID, scanner.Text())
			logMessages <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errChan <- err
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
	if err := <-errChan; err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "could not run standalone node: %s", err)
	wg.Wait()
}

func startBrowserNode(t *testing.T, ctx context.Context, url string, browserLogMessages chan<- string, browserErrors chan<- []string) {
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
				stringifiedBrowserErrors := make([]string, len(ev.Args))
				for _, arg := range ev.Args {
					stringifiedBrowserErrors = append(
						stringifiedBrowserErrors,
						fmt.Sprintf("JavaScript console error: (%s) %s %s", arg.Type, arg.Value, arg.Description),
					)
				}
				browserErrors <- stringifiedBrowserErrors
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
// a GraphQL subscription or in the TypeScript bindings and are not usually logged
// by Mesh. They need to be explicitly logged.
type orderEventLog struct {
	OrderHash string `json:"orderHash"`
	EndState  string `json:"endState"`
}

func waitForOrderEventLog(ctx context.Context, logMessages <-chan string, expectedLog orderEventLog) (string, error) {
	return waitForLogMessage(ctx, logMessages, func(msg string) bool {
		var foundLog orderEventLog
		if err := unquoteAndUnmarshal(msg, &foundLog); err != nil {
			return false
		}
		return foundLog.OrderHash == expectedLog.OrderHash &&
			foundLog.EndState == expectedLog.EndState
	})
}

// A holder type used for parsing the "signed order in browser" message that
// comes from the browser node
type signedOrderLog struct {
	Message   string `json:"message"`
	OrderHash string `json:"orderHash"`
}

func waitForSignedOrderLog(ctx context.Context, logMessages <-chan string) (string, error) {
	return waitForLogMessage(ctx, logMessages, func(msg string) bool {
		var foundLog signedOrderLog
		if err := unquoteAndUnmarshal(msg, &foundLog); err != nil {
			return false
		}
		return foundLog.Message == "signed order in browser"
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

// extractOrderHashFromLog expects a log message that contains an order hash
// and a field message that equals "signed order in browser". It extracts and
// returns the order hash.
func extractOrderHashFromLog(msg string) (string, error) {
	unquoted, err := strconv.Unquote(msg)
	if err == nil {
		msg = unquoted
	}
	holder := signedOrderLog{}
	if err := json.Unmarshal([]byte(msg), &holder); err != nil {
		return "", err
	}
	return holder.OrderHash, nil
}

var blockchainLifecycle *ethereum.BlockchainLifecycle

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}
