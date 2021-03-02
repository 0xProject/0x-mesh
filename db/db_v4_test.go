package db

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrdersV4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numOrders := 10
	orderHashes := []common.Hash{}
	orders := []*types.OrderWithMetadata{}
	for i := 0; i < numOrders; i++ {
		order := newTestOrderV4()
		orders = append(orders, order)
		orderHashes = append(orderHashes, order.Hash)
	}

	{
		alreadyStored, added, removed, err := db.AddOrdersV4(orders)
		require.NoError(t, err)
		assert.Len(t, alreadyStored, 0, "Expected no orders to be already stored")
		assert.Len(t, removed, 0, "Expected no orders to be removed")
		assertOrderSlicesAreUnsortedEqual(t, orders, added)
	}
	{
		alreadyStored, added, removed, err := db.AddOrdersV4(orders)
		require.NoError(t, err)
		assert.Len(t, alreadyStored, 10, "Expected 10 orders to be already stored")
		for _, expectedHash := range orderHashes {
			assert.Contains(t, alreadyStored, expectedHash, "Expected already stored to contain order hash")
		}
		assert.Len(t, removed, 0, "Expected no orders to be removed")
		assert.Len(t, added, 0, "Expected no orders to be added (they should already exist)")
	}

	storedOrders, err := db.FindOrdersV4(nil)
	require.NoError(t, err)
	assert.Len(t, storedOrders, numOrders)
}

func TestGetOrderV4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	_, added, _, err := db.AddOrdersV4([]*types.OrderWithMetadata{newTestOrderV4()})
	require.NoError(t, err)
	originalOrder := added[0]

	foundOrder, err := db.GetOrderV4(originalOrder.Hash)
	require.NoError(t, err)
	require.NotNil(t, foundOrder, "found order should not be nil")
	assertOrdersAreEqual(t, originalOrder, foundOrder)

	_, err = db.GetOrderV4(common.Hash{})
	assert.EqualError(t, err, ErrNotFound.Error(), "calling GetOrder with a hash that doesn't exist should return ErrNotFound")
}

func TestGetOrderStatusesV4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	removedOrder := newTestOrderV4()
	removedOrder.IsRemoved = true
	unfillableOrder := newTestOrderV4()
	unfillableOrder.IsUnfillable = true
	_, _, _, err := db.AddOrdersV4([]*types.OrderWithMetadata{removedOrder, unfillableOrder})
	require.NoError(t, err)

	hashes := []common.Hash{
		common.HexToHash("0xace746910c6a8a4730878e6e8a4abb328844c0b58f0cdfbb5b6ad28ee0bae347"),
		removedOrder.Hash,
		unfillableOrder.Hash,
	}
	actualStatuses, err := db.GetOrderStatusesV4(hashes)
	require.NoError(t, err)
	expectedStatuses := []*StoredOrderStatus{
		{
			IsStored:                 false,
			IsMarkedRemoved:          false,
			IsMarkedUnfillable:       false,
			FillableTakerAssetAmount: nil,
		},
		{
			IsStored:                 true,
			IsMarkedRemoved:          true,
			IsMarkedUnfillable:       false,
			FillableTakerAssetAmount: removedOrder.FillableTakerAssetAmount,
		},
		{
			IsStored:                 true,
			IsMarkedRemoved:          false,
			IsMarkedUnfillable:       true,
			FillableTakerAssetAmount: unfillableOrder.FillableTakerAssetAmount,
		},
	}
	assert.Equal(t, expectedStatuses, actualStatuses)
}

func TestUpdateOrderV4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	err := db.UpdateOrderV4(common.Hash{}, func(existingOrder *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		return existingOrder, nil
	})
	assert.EqualError(t, err, ErrNotFound.Error(), "calling UpdateOrder with a hash that doesn't exist should return ErrNotFound")

	// Note(albrow): We create more than one order to make sure that
	// UpdateOrder only updates one of them and does not affect the
	// others.
	numOrders := 3
	originalOrders := []*types.OrderWithMetadata{}
	for i := 0; i < numOrders; i++ {
		originalOrders = append(originalOrders, newTestOrderV4())
	}
	_, _, _, err = db.AddOrdersV4(originalOrders)
	require.NoError(t, err)

	orderToUpdate := originalOrders[0]
	updatedFillableAmount := big.NewInt(12345)
	err = db.UpdateOrderV4(orderToUpdate.Hash, func(existingOrder *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		updatedOrder := existingOrder
		updatedOrder.FillableTakerAssetAmount = updatedFillableAmount
		return updatedOrder, nil
	})
	require.NoError(t, err)

	expectedOrders := originalOrders
	expectedOrders[0].FillableTakerAssetAmount = updatedFillableAmount
	foundOrders, err := db.FindOrdersV4(nil)
	require.NoError(t, err)
	assertOrderSlicesAreUnsortedEqual(t, expectedOrders, foundOrders)
}

func TestFindOrdersV4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numOrders := 10
	originalOrders := []*types.OrderWithMetadata{}
	for i := 0; i < numOrders; i++ {
		originalOrders = append(originalOrders, newTestOrderV4())
	}
	_, _, _, err := db.AddOrdersV4(originalOrders)
	require.NoError(t, err)

	foundOrders, err := db.FindOrdersV4(nil)
	require.NoError(t, err)
	assertOrderSlicesAreUnsortedEqual(t, originalOrders, foundOrders)
}

func TestFindOrdersFilterSortLimitAndOffsetV4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)
	storedOrders := createAndStoreOrdersForFilterTestsV4(t, db)

	query := &OrderQueryV4{
		Filters: []OrderFilterV4{
			{
				Field: OV4FMakerAmount,
				Kind:  GreaterOrEqual,
				Value: big.NewInt(3),
			},
		},
		Sort: []OrderSortV4{
			{
				Field:     OV4FMakerAmount,
				Direction: Ascending,
			},
		},
		Limit:  3,
		Offset: 2,
	}
	expectedOrders := storedOrders[5:8]
	actualOrders, err := db.FindOrdersV4(query)
	require.NoError(t, err)
	assertOrderSlicesAreEqual(t, expectedOrders, actualOrders)
}
