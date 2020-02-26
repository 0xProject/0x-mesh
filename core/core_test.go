// +build !js

package core

import (
	"context"
	"flag"
	"math/big"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// blockProcessingWaitTime is the amount of time to wait for Mesh to process
	// new blocks that have been mined.
	blockProcessingWaitTime = 1 * time.Second
	// ordersyncWaitTime is the amount of time to wait for ordersync to run.
	ordersyncWaitTime = 2 * time.Second
)

func TestEthereumChainDetection(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	// simulate starting up on mainnet
	_, err = initMetadata(1, meshDB)
	require.NoError(t, err)

	// simulate restart on same chain
	_, err = initMetadata(1, meshDB)
	require.NoError(t, err)

	// should error when attempting to start on different chain
	_, err = initMetadata(2, meshDB)
	assert.Error(t, err)
}

func newTestApp(t *testing.T) *App {
	dataDir := "/tmp/test_node/" + uuid.New().String()
	config := Config{
		Verbosity:                        5,
		DataDir:                          dataDir,
		P2PTCPPort:                       0,
		P2PWebSocketsPort:                0,
		EthereumRPCURL:                   constants.GanacheEndpoint,
		EthereumChainID:                  constants.TestChainID,
		UseBootstrapList:                 false,
		BootstrapList:                    "",
		BlockPollingInterval:             250 * time.Millisecond,
		EthereumRPCMaxContentLength:      524288,
		EnableEthereumRPCRateLimiting:    false,
		EthereumRPCMaxRequestsPer24HrUTC: 99999999999999,
		EthereumRPCMaxRequestsPerSecond:  99999999999999,
		MaxOrdersInStorage:               100000,
		CustomOrderFilter:                "{}",
	}
	app, err := New(config)
	require.NoError(t, err)
	return app
}

var (
	rpcClient           *ethrpc.Client
	ethClient           *ethclient.Client
	blockchainLifecycle *ethereum.BlockchainLifecycle
)

// Since these tests must be run sequentially, we don't want them to run as part of
// the normal testing process. They will only be run if the "--serial" flag is used.
var serialTestsEnabled bool

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
	testing.Init()
	flag.Parse()

	var err error
	rpcClient, err = ethrpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(rpcClient)
	if err != nil {
		panic(err)
	}
}

func TestNewGoRoutineLeak(t *testing.T) {
	beforeGoRoutineCount := runtime.NumGoroutine()
	app := newTestApp(t)
	app.db.Close()
	time.Sleep(5 * time.Second)
	afterGoRoutineCount := runtime.NumGoroutine()
	require.Equal(t, beforeGoRoutineCount, afterGoRoutineCount)
}

func TestOrderSync(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	// Set up two Mesh nodes. originalNode starts with some orders. newNode enters
	// the network without any orders.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	originalNode := newTestApp(t)
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, originalNode.Start(ctx))
	}()
	newNode := newTestApp(t)
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, newNode.Start(ctx))
	}()

	// Manually add some orders to originalNode.
	originalOrders := make([]*zeroex.SignedOrder, 10)
	for i := range originalOrders {
		originalOrders[i] = scenario.CreateWETHForZRXSignedTestOrder(t, ethClient, constants.GanacheAccount1, constants.GanacheAccount2, big.NewInt(20), big.NewInt(5))
	}

	// We have to wait for latest block to be processed by the Mesh node.
	time.Sleep(blockProcessingWaitTime)

	results, err := originalNode.orderWatcher.ValidateAndStoreValidOrders(ctx, originalOrders, true, constants.TestChainID)
	require.NoError(t, err)
	require.Empty(t, results.Rejected, "tried to add orders but some were invalid: \n%s\n", spew.Sdump(results))

	orderEventsChan := make(chan []*zeroex.OrderEvent)
	orderEventsSub := newNode.SubscribeToOrderEvents(orderEventsChan)
	defer orderEventsSub.Unsubscribe()

	// Connect the two nodes *after* adding orders to one of them. This should
	// trigger the ordersync protocol.
	err = originalNode.AddPeer(peer.AddrInfo{
		ID:    newNode.node.ID(),
		Addrs: newNode.node.Multiaddrs(),
	})
	require.NoError(t, err)

	// Wait for newNode to get the orders via ordersync.
	receivedAddedEvents := []*zeroex.OrderEvent{}
OrderEventLoop:
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for %d order added events (received %d so far)", len(originalOrders), len(receivedAddedEvents))
		case orderEvents := <-orderEventsChan:
			for _, orderEvent := range orderEvents {
				if orderEvent.EndState == zeroex.ESOrderAdded {
					receivedAddedEvents = append(receivedAddedEvents, orderEvent)
				}
			}
			if len(receivedAddedEvents) >= len(originalOrders) {
				break OrderEventLoop
			}
		}
	}

	// Test that the orders are actually in the database and are returned by
	// GetOrders.
	newNodeOrdersResp, err := newNode.GetOrders(0, len(originalOrders), "")
	require.NoError(t, err)
	assert.Len(t, newNodeOrdersResp.OrdersInfos, len(originalOrders), "new node should have %d orders", len(originalOrders))
	for _, expectedOrder := range originalOrders {
		orderHash, err := expectedOrder.ComputeOrderHash()
		require.NoError(t, err)
		expectedOrder.ResetHash()
		var dbOrder meshdb.Order
		require.NoError(t, newNode.db.Orders.FindByID(orderHash.Bytes(), &dbOrder))
		actualOrder := dbOrder.SignedOrder
		assert.Equal(t, expectedOrder, actualOrder, "correct order was not stored in new node database")
	}

	// Wait for nodes to exit without error.
	cancel()
	wg.Wait()
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}
