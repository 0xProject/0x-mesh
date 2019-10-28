package core

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageSharingIsolated(t *testing.T) {
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
		orders := randomOrders(t, testCase.orderCount)

		selector := &OrderSelector{
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
	var allOrders []*meshdb.Order
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	defer meshDB.Close()
	require.NoError(t, err)

	selector := &OrderSelector{
		nextOffset: 0,
		db:         meshDB,
	}

	// Add five orders to the database
	orders := randomOrders(t, 5)
	allOrders = orders
	insertOrders(t, selector, orders)
	verifyRoundRobinSharing(t, selector, 0, 3)

	// Add seven more orders to the database
	orders = randomOrders(t, 7)
	allOrders = append(allOrders, orders...)
	insertOrders(t, selector, orders)
	verifyRoundRobinSharing(t, selector, 3, 10)

	// Delete three orders from the database
	deleteOrders(t, selector, allOrders[0:3:3])
	verifyRoundRobinSharing(t, selector, selector.nextOffset, 5)

	// Add 12 more orders to the database
	orders = randomOrders(t, 12)
	allOrders = append(allOrders, orders...)
	insertOrders(t, selector, orders)
	verifyRoundRobinSharing(t, selector, selector.nextOffset, 7)
}

// Verify that the correct messages are shared by `GetMessagesToShare` given `orders`, a `nextOffset`, and `max`
func verifyRoundRobinSharing(t *testing.T, selector *OrderSelector, nextOffset int, max int) {
	notRemovedFilter := selector.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	ordersSnapshot, err := selector.db.Orders.GetSnapshot()

	// Get the number of orders in the database
	count, err := ordersSnapshot.NewQuery(notRemovedFilter).Count()
	require.NoError(t, err)

	expectedOrdersLength := min(max, count)
	expectedOrders := make([][]byte, expectedOrdersLength)

	// Get all of the orders in the database.
	var orderList []*meshdb.Order
	err = ordersSnapshot.NewQuery(notRemovedFilter).Offset(0).Max(count).Run(&orderList)

	// Update `nextOffset` to zero if it is larger than the number of orders that are stored
	if nextOffset > count {
		nextOffset = 0
	}

	// Calculate the orders that we expect to be shared
	for i := 0; i < expectedOrdersLength; i++ {
		encoding, err := encodeOrder(orderList[(nextOffset+i)%count].SignedOrder)
		require.NoError(t, err)

		expectedOrders[i] = encoding
	}

	// Get the actual list of orders that are shared
	actualOrders, err := selector.GetMessagesToShare(max)
	require.NoError(t, err)

	// Ensure that the result from `GetMessagesToShare` matches the expected result.
	assert.Equal(t, expectedOrders, actualOrders)
}

func deleteOrders(t *testing.T, selector *OrderSelector, orders []*meshdb.Order) {
	for _, order := range orders {
		err := selector.db.Orders.Delete(order.ID())
		require.NoError(t, err)
	}
}

func insertOrders(t *testing.T, selector *OrderSelector, orders []*meshdb.Order) {
	for _, order := range orders {
		err := selector.db.Orders.Insert(order)
		require.NoError(t, err)
	}
}

func randomOrders(t *testing.T, orderCount int) []*meshdb.Order {
	orders := make([]*meshdb.Order, orderCount)

	for i := 0; i < orderCount; i++ {
		orders[i] = randomOrder(t)
	}

	return orders
}

func randomOrder(t *testing.T) *meshdb.Order {
	signedOrder := &zeroex.SignedOrder{
		Order: zeroex.Order{
			MakerAddress:          constants.GanacheAccount0,
			TakerAddress:          common.HexToAddress(randomAddress(t)),
			SenderAddress:         common.HexToAddress(randomAddress(t)),
			FeeRecipientAddress:   common.HexToAddress(randomAddress(t)),
			MakerAssetData:        common.Hex2Bytes(randomAssetData(t)),
			TakerAssetData:        common.Hex2Bytes(randomAssetData(t)),
			Salt:                  big.NewInt(rand.Int63()),
			MakerFee:              big.NewInt(rand.Int63()),
			TakerFee:              big.NewInt(rand.Int63()),
			MakerAssetAmount:      big.NewInt(rand.Int63()),
			TakerAssetAmount:      big.NewInt(rand.Int63()),
			ExpirationTimeSeconds: big.NewInt(time.Now().Add(48 * time.Hour).Unix()),
			ExchangeAddress:       common.HexToAddress(randomAddress(t)),
		},
		Signature: []byte(randomHex(t, 65)),
	}
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	return &meshdb.Order{
		Hash:                     orderHash,
		SignedOrder:              signedOrder,
		LastUpdated:              time.Now().UTC(),
		FillableTakerAssetAmount: signedOrder.Order.TakerAssetAmount,
		IsRemoved:                false,
	}
}

func randomAddress(t *testing.T) string {
	return "0x" + randomHex(t, 20)
}

func randomAssetData(t *testing.T) string {
	// Note: Asset data must begin with a valid asset proxy id or parsing will fail
	return "f47261b0000000000000000000000000" + randomHex(t, 20)
}

func randomHex(t *testing.T, n int) string {
	bytes := make([]byte, n)
	_, err := cryptoRand.Read(bytes)
	require.NoError(t, err)
	return hex.EncodeToString(bytes)
}
