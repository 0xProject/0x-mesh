// +build !js

package core

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// blockProcessingWaitTime is the amount of time to wait for Mesh to process
	// new blocks that have been mined.
	blockProcessingWaitTime = 1 * time.Second
)

func TestEthereumChainDetection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	// simulate starting up on mainnet
	err = initMetadata(1, database)
	require.NoError(t, err)

	// simulate restart on same chain
	err = initMetadata(1, database)
	require.NoError(t, err)

	// should error when attempting to start on different chain
	err = initMetadata(2, database)
	assert.Error(t, err)
}

func TestConfigChainIDAndRPCMatchDetection(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	dataDir := "/tmp/test_node/" + uuid.New().String()
	config := Config{
		Verbosity:                        5,
		DataDir:                          dataDir,
		P2PTCPPort:                       0,
		P2PWebSocketsPort:                0,
		EthereumRPCURL:                   constants.GanacheEndpoint,
		EthereumChainID:                  42, // RPC has chain id 1337
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
	app, err := New(ctx, config)
	require.NoError(t, err)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.Start()
		require.Error(t, err)
		require.Contains(t, err.Error(), "ChainID mismatch")
	}()

	// Wait for nodes to exit without error.
	wg.Wait()
}

func newTestAppWithPrivateConfig(t *testing.T, ctx context.Context, customOrderFilter string, pConfig privateConfig) *App {
	if customOrderFilter == "" {
		customOrderFilter = defaultOrderFilter
	}
	dataDir := "/tmp/test_node/" + uuid.New().String()
	config := Config{
		Verbosity:                        2,
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
		CustomOrderFilter:                customOrderFilter,
	}
	app, err := newWithPrivateConfig(ctx, config, pConfig)
	require.NoError(t, err)
	return app
}

var (
	rpcClient           *ethrpc.Client
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
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(rpcClient)
	if err != nil {
		panic(err)
	}
}

func TestRepeatedAppInitialization(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dataDir := "/tmp/test_node/" + uuid.New().String()
	config := Config{
		Verbosity:                        2,
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
		CustomContractAddresses:          `{"exchange":"0x48bacb9266a570d521063ef5dd96e61686dbe788","devUtils":"0x38ef19fdf8e8415f18c307ed71967e19aac28ba1","erc20Proxy":"0x1dc4c1cefef38a777b15aa20260a54e584b16c48","erc721Proxy":"0x1d7022f5b17d2f8b695918fb48fa1089c9f85401","erc1155Proxy":"0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"}`,
	}
	_, err := New(ctx, config)
	require.NoError(t, err)
	_, err = New(ctx, config)
	require.NoError(t, err)
}

func TestOrderSync(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	testCases := []ordersyncTestCase{
		{
			name: "FilteredPaginationSubprotocol version 0",
			pConfig: privateConfig{
				paginationSubprotocolPerPage: 10,
				paginationSubprotocols: []ordersyncSubprotocolFactory{
					NewFilteredPaginationSubprotocolV0,
				},
			},
		},
		{
			name: "FilteredPaginationSubprotocol version 1",
			pConfig: privateConfig{
				paginationSubprotocolPerPage: 10,
				paginationSubprotocols: []ordersyncSubprotocolFactory{
					NewFilteredPaginationSubprotocolV1,
				},
			},
		},
		{
			name: "FilteredPaginationSubprotocol version 1 and version 0",
			pConfig: privateConfig{
				paginationSubprotocolPerPage: 10,
				paginationSubprotocols: []ordersyncSubprotocolFactory{
					NewFilteredPaginationSubprotocolV1,
					NewFilteredPaginationSubprotocolV0,
				},
			},
		},
		{
			name:              "makerAssetAmount orderfilter - match all orders",
			customOrderFilter: `{"properties":{"makerAssetAmount":{"pattern":"^1$","type":"string"}}}`,
			orderOptionsForIndex: func(_ int) []orderopts.Option {
				return []orderopts.Option{orderopts.MakerAssetAmount(big.NewInt(1))}
			},
			pConfig: privateConfig{
				paginationSubprotocolPerPage: 10,
				paginationSubprotocols: []ordersyncSubprotocolFactory{
					NewFilteredPaginationSubprotocolV1,
					NewFilteredPaginationSubprotocolV0,
				},
			},
		},
		{
			name:              "makerAssetAmount OrderFilter - matches one order",
			customOrderFilter: `{"properties":{"makerAssetAmount":{"pattern":"^1$","type":"string"}}}`,
			orderOptionsForIndex: func(i int) []orderopts.Option {
				if i == 0 {
					return []orderopts.Option{orderopts.MakerAssetAmount(big.NewInt(1))}
				}
				return []orderopts.Option{}
			},
			pConfig: privateConfig{
				paginationSubprotocolPerPage: 10,
				paginationSubprotocols: []ordersyncSubprotocolFactory{
					NewFilteredPaginationSubprotocolV1,
					NewFilteredPaginationSubprotocolV0,
				},
			},
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("%s (test case %d)", testCase.name, i)
		t.Run(testCaseName, runOrdersyncTestCase(testCase))
	}
}

type ordersyncTestCase struct {
	name                 string
	customOrderFilter    string
	orderOptionsForIndex func(int) []orderopts.Option
	pConfig              privateConfig
}

const defaultOrderFilter = "{}"

func runOrdersyncTestCase(testCase ordersyncTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		teardownSubTest := setupSubTest(t)
		defer teardownSubTest(t)

		// Set up two Mesh nodes. originalNode starts with some orders. newNode enters
		// the network without any orders.
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		wg := &sync.WaitGroup{}
		originalNode := newTestAppWithPrivateConfig(t, ctx, defaultOrderFilter, testCase.pConfig)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := originalNode.Start(); err != nil && err != context.Canceled {
				// context.Canceled is expected. For any other error, fail the test.
				panic(fmt.Sprintf("%s %s", testCase.name, err))
			}
		}()

		// Manually add some orders to originalNode.
		orderOptionsForIndex := func(i int) []orderopts.Option {
			orderOptions := []orderopts.Option{orderopts.SetupMakerState(true)}
			if testCase.orderOptionsForIndex != nil {
				return append(testCase.orderOptionsForIndex(i), orderOptions...)
			}
			return orderOptions
		}
		numOrders := testCase.pConfig.paginationSubprotocolPerPage*3 + 1
		originalOrders := scenario.NewSignedTestOrdersBatch(t, numOrders, orderOptionsForIndex)

		// We have to wait for latest block to be processed by the Mesh node.
		time.Sleep(blockProcessingWaitTime)

		results, err := originalNode.orderWatcher.ValidateAndStoreValidOrders(ctx, originalOrders, true, constants.TestChainID)
		require.NoError(t, err)
		require.Empty(t, results.Rejected, "tried to add orders but some were invalid: \n%s\n", spew.Sdump(results))

		newNode := newTestAppWithPrivateConfig(t, ctx, testCase.customOrderFilter, defaultPrivateConfig())
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := newNode.Start(); err != nil && err != context.Canceled {
				// context.Canceled is expected. For any other error, fail the test.
				panic(fmt.Sprintf("%s %s", testCase.name, err))
			}
		}()
		<-newNode.started

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

		// Only the orders that satisfy the new node's orderfilter should
		// be received during ordersync.
		filteredOrders := []*zeroex.SignedOrder{}
		for _, order := range originalOrders {
			matches, err := newNode.orderFilter.MatchOrder(order)
			require.NoError(t, err)
			if matches {
				filteredOrders = append(filteredOrders, order)
			}
		}

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
				if len(receivedAddedEvents) >= len(filteredOrders) {
					break OrderEventLoop
				}
			}
		}

		// Test that the orders are actually in the database and are returned by
		// GetOrders.
		newNodeOrdersResp, err := newNode.GetOrders(len(filteredOrders), common.Hash{})
		require.NoError(t, err)
		assert.Len(t, newNodeOrdersResp.OrdersInfos, len(filteredOrders), "new node should have %d orders", len(originalOrders))
		for _, expectedOrder := range filteredOrders {
			orderHash, err := expectedOrder.ComputeOrderHash()
			require.NoError(t, err)
			expectedOrder.ResetHash()
			dbOrder, err := newNode.db.GetOrder(orderHash)
			require.NoError(t, err)
			actualOrder := dbOrder.SignedOrder()
			assert.Equal(t, expectedOrder, actualOrder, "correct order was not stored in new node database")
		}

		// Wait for nodes to exit without error.
		cancel()
		wg.Wait()
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}
