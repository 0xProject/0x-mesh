package core

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageSharingIsolated(t *testing.T) {
	ethClient, err := ethrpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	signer := ethereum.NewEthRPCSigner(ethClient)
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	defer meshDB.Close()
	require.NoError(t, err)

	testCases := []struct {
		orders     []*meshdb.Order
		nextOffset int
		max        int
	}{
		{
			orders:     []*meshdb.Order{randomOrder(t, signer)},
			nextOffset: 0,
			max:        1,
		},
		{
			orders:     randomOrders(t, signer, 2),
			nextOffset: 0,
			max:        1,
		},
		{
			orders:     randomOrders(t, signer, 2),
			nextOffset: 1,
			max:        1,
		},
		{
			orders:     randomOrders(t, signer, 2),
			nextOffset: 1,
			max:        2,
		},
		{
			orders:     randomOrders(t, signer, 10),
			nextOffset: 5,
			max:        2,
		},
		{
			orders:     randomOrders(t, signer, 10),
			nextOffset: 1,
			max:        2,
		},
		{
			orders:     randomOrders(t, signer, 3),
			nextOffset: 1,
			max:        4,
		},
		{
			orders:     randomOrders(t, signer, 3),
			nextOffset: 20,
			max:        4,
		},
	}

	for _, testCase := range testCases {
		selector := &OrderSelector{
			nextOffset: testCase.nextOffset,
			db:         meshDB,
		}

		// Insert the orders into the database
		for _, order := range testCase.orders {
			err := selector.db.Orders.Insert(order)
			require.NoError(t, err)
		}

		var orderList []*meshdb.Order
		notRemovedFilter := selector.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
		ordersSnapshot, err := selector.db.Orders.GetSnapshot()
		err = ordersSnapshot.NewQuery(notRemovedFilter).Offset(0).Max(len(testCase.orders)).Run(&orderList)
		ordersLength := len(testCase.orders)
		expectedOrdersLength := min(testCase.max, ordersLength)
		expectedOrders := make([][]byte, expectedOrdersLength)

		// Update `nextOffset` to zero if it is larger than the number of orders that are stored
		if testCase.nextOffset > ordersLength {
			testCase.nextOffset = 0
		}

		// Calculate the orders that we expect to be shared
		for i := 0; i < expectedOrdersLength; i++ {
			encoding, err := encodeOrder(orderList[(testCase.nextOffset+i)%ordersLength].SignedOrder)
			require.NoError(t, err)

			expectedOrders[i] = encoding
		}

		// Get the actual list of orders that are shared
		actualOrders, err := selector.GetMessagesToShare(testCase.max)
		require.NoError(t, err)

		// Ensure that the result from `GetMessagesToShare` matches the expected result.
		assert.Equal(t, expectedOrders, actualOrders)

		// Delete the orders from the database
		for _, order := range testCase.orders {
			err := selector.db.Orders.Delete(order.ID())
			require.NoError(t, err)
		}
	}
}

func newOrderSelector(t *testing.T) *OrderSelector {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	return &OrderSelector{
		nextOffset: 0,
		db:         meshDB,
	}
}

func randomOrders(t *testing.T, signer ethereum.Signer, count int) []*meshdb.Order {
	orders := make([]*meshdb.Order, count)

	for i := 0; i < count; i++ {
		orders[i] = randomOrder(t, signer)
	}

	return orders
}

func randomOrder(t *testing.T, signer ethereum.Signer) *meshdb.Order {
	order := &zeroex.Order{
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
	}
	randomOrder, err := zeroex.SignOrder(signer, order)
	require.NoError(t, err)
	orderHash, err := randomOrder.ComputeOrderHash()
	require.NoError(t, err)
	return &meshdb.Order{
		Hash:                     orderHash,
		SignedOrder:              randomOrder,
		LastUpdated:              time.Now().UTC(),
		FillableTakerAssetAmount: randomOrder.Order.TakerAssetAmount,
		IsRemoved:                false,
	}
}

func randomAddress(t *testing.T) string {
	return "0x" + randomHex(t, 20)
}

func randomAssetData(t *testing.T) string {
	// Note: Asset data must begin with a valid asset proxy id or else
	//       parsing will fail
	return "f47261b0000000000000000000000000" + randomHex(t, 20)
}

func randomHex(t *testing.T, n int) string {
	bytes := make([]byte, n)
	_, err := cryptoRand.Read(bytes)
	require.NoError(t, err)
	return hex.EncodeToString(bytes)
}
