// +build !js

package core

import (
	"flag"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	rpcClient              *ethrpc.Client
	ethClient              *ethclient.Client
	blockchainLifecycle    *ethereum.BlockchainLifecycle
	makerAddress           = constants.GanacheAccount1
	takerAddress           = constants.GanacheAccount2
	tenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	wethAmount             = new(big.Int).Mul(big.NewInt(2), tenDecimalsInBaseUnits)
	zrxAmount              = new(big.Int).Mul(big.NewInt(1), tenDecimalsInBaseUnits)
)

// Since these tests must be run sequentially, we don't want them to run as part of
// the normal testing process. They will only be run if the "--serial" flag is used.
var serialTestsEnabled bool

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
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

func TestMessageSharingIsolated(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	defer meshDB.Close()
	require.NoError(t, err)

	testCases := []struct {
		orderCount int
		nextOffset int
		max        int
	}{
		{
			orderCount: 1,
			nextOffset: 0,
			max:        1,
		},
		{
			orderCount: 2,
			nextOffset: 0,
			max:        1,
		},
		{
			orderCount: 2,
			nextOffset: 1,
			max:        1,
		},
		{
			orderCount: 2,
			nextOffset: 1,
			max:        2,
		},
		{
			orderCount: 10,
			nextOffset: 5,
			max:        2,
		},
		{
			orderCount: 10,
			nextOffset: 1,
			max:        2,
		},
		{
			orderCount: 3,
			nextOffset: 1,
			max:        4,
		},
		{
			orderCount: 3,
			nextOffset: 20,
			max:        4,
		},
	}

	for _, testCase := range testCases {
		orders, err := signedTestOrders(t, testCase.orderCount)

		require.NoError(t, err)

		selector := &orderSelector{
			nextOffset: testCase.nextOffset,
			db:         meshDB,
		}

		// Insert the orders into the database
		insertOrders(t, selector, orders)

		// Ensure that the correct orders are shared
		verifyRoundRobinSharing(t, selector, testCase.nextOffset, testCase.max)

		// Delete the orders from the database
		deleteOrders(t, selector, orders)
	}
}

func TestMessagesSharedSerial(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	var allOrders []*meshdb.Order
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	defer meshDB.Close()
	require.NoError(t, err)

	selector := &orderSelector{
		topic:      "customTopic",
		nextOffset: 0,
		db:         meshDB,
	}

	// Add five orders to the database
	orders, err := signedTestOrders(t, 5)
	require.NoError(t, err)
	allOrders = orders
	insertOrders(t, selector, orders)
	verifyRoundRobinSharing(t, selector, 0, 3)

	// Add seven more orders to the database
	orders, err = signedTestOrders(t, 7)
	require.NoError(t, err)
	allOrders = append(allOrders, orders...)
	insertOrders(t, selector, orders)
	verifyRoundRobinSharing(t, selector, 3, 10)

	// Delete three orders from the database
	deleteOrders(t, selector, allOrders[0:3:3])
	verifyRoundRobinSharing(t, selector, selector.nextOffset, 5)

	// Add 12 more orders to the database
	orders, err = signedTestOrders(t, 12)
	require.NoError(t, err)
	allOrders = append(allOrders, orders...)
	insertOrders(t, selector, orders)
	verifyRoundRobinSharing(t, selector, selector.nextOffset, 7)
}

// Verify that the correct messages are shared by `GetMessagesToShare` given `orders`, a `nextOffset`, and `max`
func verifyRoundRobinSharing(t *testing.T, selector *orderSelector, nextOffset int, max int) {
	notRemovedFilter := selector.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})

	// Get the number of orders in the database
	count, err := selector.db.Orders.NewQuery(notRemovedFilter).Count()
	require.NoError(t, err)

	expectedOrdersLength := min(max, count)
	expectedOrders := make([][]byte, expectedOrdersLength)

	// Get all of the orders in the database.
	var orderList []*meshdb.Order
	err = selector.db.Orders.NewQuery(notRemovedFilter).Run(&orderList)
	require.NoError(t, err)

	// Update `nextOffset` to zero if it is larger than the number of orders that are stored
	if nextOffset > count {
		nextOffset = 0
	}

	// Calculate the orders that we expect to be shared
	for i := 0; i < expectedOrdersLength; i++ {
		encodedOrder, err := encodeOrderMessage(selector.topic, orderList[(nextOffset+i)%count].SignedOrder)
		require.NoError(t, err)

		expectedOrders[i] = encodedOrder
	}

	// Get the actual list of orders that are shared
	actualOrders, err := selector.GetMessagesToShare(max)
	require.NoError(t, err)

	// Ensure that the result from `GetMessagesToShare` matches the expected result.
	assert.Equal(t, expectedOrders, actualOrders)
}

func deleteOrders(t *testing.T, selector *orderSelector, orders []*meshdb.Order) {
	for _, order := range orders {
		err := selector.db.Orders.Delete(order.ID())
		require.NoError(t, err)
	}
}

func insertOrders(t *testing.T, selector *orderSelector, orders []*meshdb.Order) {
	for _, order := range orders {
		err := selector.db.Orders.Insert(order)
		require.NoError(t, err)
	}
}

func signedTestOrders(t *testing.T, orderCount int) ([]*meshdb.Order, error) {
	orders := make([]*meshdb.Order, orderCount)

	for i := range orders {
		order := scenario.CreateZRXForWETHSignedTestOrder(
			t,
			ethClient,
			makerAddress,
			takerAddress,
			new(big.Int).Add(wethAmount, big.NewInt(int64(i))),
			zrxAmount,
		)

		hash, err := order.ComputeOrderHash()

		if err != nil {
			return nil, err
		}

		orders[i] = &meshdb.Order{
			Hash:                     hash,
			SignedOrder:              order,
			LastUpdated:              time.Now(),
			FillableTakerAssetAmount: order.TakerAssetAmount,
			IsRemoved:                false,
		}
	}

	return orders, nil
}
