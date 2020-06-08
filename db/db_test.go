package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var contractAddresses = ethereum.GanacheAddresses

func TestAddOrders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numOrders := 10
	orders := []*Order{}
	for i := 0; i < numOrders; i++ {
		orders = append(orders, newTestOrder())
	}

	{
		added, removed, err := db.AddOrders(orders)
		require.NoError(t, err)
		assert.Len(t, removed, 0, "Expected no orders to be removed")
		assertOrderSlicesAreUnsortedEqual(t, orders, added)
	}
	{
		added, removed, err := db.AddOrders(orders)
		require.NoError(t, err)
		assert.Len(t, removed, 0, "Expected no orders to be removed")
		assert.Len(t, added, 0, "Expected no orders to be added (they should already exist)")
	}
}

func TestGetOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	added, _, err := db.AddOrders([]*Order{newTestOrder()})
	require.NoError(t, err)
	originalOrder := added[0]

	foundOrder, err := db.GetOrder(originalOrder.Hash)
	require.NoError(t, err)
	require.NotNil(t, foundOrder, "found order should not be nil")
	assertOrdersAreEqual(t, *originalOrder, *foundOrder)
}

func TestUpdateOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	// Note(albrow): We create more than one order to make sure that
	// UpdateOrder only updates one of them and does not affect the
	// others.
	numOrders := 3
	originalOrders := []*Order{}
	for i := 0; i < numOrders; i++ {
		originalOrders = append(originalOrders, newTestOrder())
	}
	_, _, err := db.AddOrders(originalOrders)
	require.NoError(t, err)

	orderToUpdate := originalOrders[0]
	updatedFillableAmount := NewUint256(big.NewInt(12345))
	err = db.UpdateOrder(orderToUpdate.Hash, func(existingOrder *Order) (*Order, error) {
		updatedOrder := existingOrder
		updatedOrder.FillableTakerAssetAmount = updatedFillableAmount
		return updatedOrder, nil
	})

	expectedOrders := originalOrders
	expectedOrders[0].FillableTakerAssetAmount = updatedFillableAmount
	foundOrders, err := db.FindOrders(nil)
	require.NoError(t, err)
	assertOrderSlicesAreUnsortedEqual(t, expectedOrders, foundOrders)
}

func TestFindOrders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numOrders := 10
	originalOrders := []*Order{}
	for i := 0; i < numOrders; i++ {
		originalOrders = append(originalOrders, newTestOrder())
	}
	_, _, err := db.AddOrders(originalOrders)
	require.NoError(t, err)

	foundOrders, err := db.FindOrders(nil)
	require.NoError(t, err)
	assertOrderSlicesAreUnsortedEqual(t, originalOrders, foundOrders)
}

func TestFindOrdersSort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	// Create some test orders with carefully chosen MakerAssetAmount
	// and TakerAssetAmount values for testing sorting.
	numOrders := 5
	originalOrders := []*Order{}
	for i := 0; i < numOrders; i++ {
		order := newTestOrder()
		order.MakerAssetAmount = NewUint256(big.NewInt(int64(i)))
		// It's important for some orders to have the same TakerAssetAmount
		// so that we can test secondary sorts (sorting on more than one
		// field).
		if i%2 == 0 {
			order.TakerAssetAmount = NewUint256(big.NewInt(100))
		} else {
			order.TakerAssetAmount = NewUint256(big.NewInt(200))
		}
		originalOrders = append(originalOrders, order)
	}
	_, _, err := db.AddOrders(originalOrders)
	require.NoError(t, err)

	testCases := []findOrdersSortTestCase{
		{
			sortOpts: []OrderSort{
				{
					Field:     OFMakerAssetAmount,
					Direction: Ascending,
				},
			},
			less: lessByMakerAssetAmountAsc,
		},
		{
			sortOpts: []OrderSort{
				{
					Field:     OFMakerAssetAmount,
					Direction: Descending,
				},
			},
			less: lessByMakerAssetAmountDesc,
		},
		{
			sortOpts: []OrderSort{
				{
					Field:     OFTakerAssetAmount,
					Direction: Ascending,
				},
				{
					Field:     OFMakerAssetAmount,
					Direction: Ascending,
				},
			},
			less: lessByTakerAssetAmountAscAndMakerAssetAmountAsc,
		},
		{
			sortOpts: []OrderSort{
				{
					Field:     OFTakerAssetAmount,
					Direction: Descending,
				},
				{
					Field:     OFMakerAssetAmount,
					Direction: Descending,
				},
			},
			less: lessByTakerAssetAmountDescAndMakerAssetAmountDesc,
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("test case %d", i)
		t.Run(testCaseName, runFindOrdersSortTestCase(t, db, originalOrders, testCase))
	}
}

type findOrdersSortTestCase struct {
	sortOpts []OrderSort
	less     func([]*Order) func(i, j int) bool
}

func runFindOrdersSortTestCase(t *testing.T, db *DB, originalOrders []*Order, testCase findOrdersSortTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		expectedOrders := make([]*Order, len(originalOrders))
		copy(expectedOrders, originalOrders)
		sort.Slice(expectedOrders, testCase.less(expectedOrders))
		findOpts := &FindOrdersOpts{
			Sort: testCase.sortOpts,
		}
		foundOrders, err := db.FindOrders(findOpts)
		require.NoError(t, err)
		assertOrderSlicesAreEqual(t, expectedOrders, foundOrders)
	}
}

func TestFindOrdersLimitAndOffset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numOrders := 10
	originalOrders := []*Order{}
	for i := 0; i < numOrders; i++ {
		originalOrders = append(originalOrders, newTestOrder())
	}
	_, _, err := db.AddOrders(originalOrders)
	require.NoError(t, err)
	sortOrdersByHash(originalOrders)

	testCases := []findOrdersLimitAndOffsetTestCase{
		{
			limit:          0,
			offset:         0,
			expectedOrders: originalOrders,
		},
		{
			limit:          3,
			offset:         0,
			expectedOrders: originalOrders[:3],
		},
		{
			limit:         0,
			offset:        3,
			expectedError: "can't use Offset without Limit",
		},
		{
			limit:          10,
			offset:         3,
			expectedOrders: originalOrders[3:],
		},
		{
			limit:          4,
			offset:         3,
			expectedOrders: originalOrders[3:7],
		},
		{
			limit:          10,
			offset:         10,
			expectedOrders: []*Order{},
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("test case %d", i)
		t.Run(testCaseName, runFindOrdersLimitAndOffsetTestCase(t, db, originalOrders, testCase))
	}
}

type findOrdersLimitAndOffsetTestCase struct {
	limit          uint
	offset         uint
	expectedOrders []*Order
	expectedError  string
}

func runFindOrdersLimitAndOffsetTestCase(t *testing.T, db *DB, originalOrders []*Order, testCase findOrdersLimitAndOffsetTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		findOpts := &FindOrdersOpts{
			Sort: []OrderSort{
				{
					Field:     OFHash,
					Direction: Ascending,
				},
			},
			Limit:  testCase.limit,
			Offset: testCase.offset,
		}

		foundOrders, err := db.FindOrders(findOpts)
		if testCase.expectedError != "" {
			require.Error(t, err, "expected an error but got nil")
			assert.Contains(t, err.Error(), testCase.expectedError, "wrong error message")
		} else {
			require.NoError(t, err)
			assertOrderSlicesAreEqual(t, testCase.expectedOrders, foundOrders)
		}
	}
}

func TestFindOrdersFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	// Create some test orders with very specific characteristics to make it easier to write tests.
	// - Both MakerAssetAmount and TakerAssetAmount will be 0, 1, 2, etc.
	// - MakerAssetData will be 'a', 'b', 'c', etc.
	// - ParsedMakerAssetData will always be for the ERC721Dummy contract, and each will contain
	//   two token ids: (0, 1), (0, 11), (0, 21), (0, 31) etc.
	numOrders := 10
	originalOrders := []*Order{}
	for i := 0; i < numOrders; i++ {
		order := newTestOrder()
		order.MakerAssetAmount = NewUint256(big.NewInt(int64(i)))
		order.TakerAssetAmount = NewUint256(big.NewInt(int64(i)))
		order.MakerAssetData = []byte{97 + byte(i)}
		parsedMakerAssetData := ParsedAssetData([]SingleAssetData{
			{
				Address: constants.GanacheDummyERC721TokenAddress,
				TokenID: NewUint256(big.NewInt(0)),
			},
			{
				Address: constants.GanacheDummyERC721TokenAddress,
				TokenID: NewUint256(big.NewInt(int64(i)*10 + 1)),
			},
		})
		order.ParsedMakerAssetData = &parsedMakerAssetData
		originalOrders = append(originalOrders, order)
	}
	_, _, err := db.AddOrders(originalOrders)
	require.NoError(t, err)

	testCases := []findOrdersFilterTestCase{
		{
			name:           "no filter",
			filters:        []OrderFilter{},
			expectedOrders: originalOrders,
		},

		// Filter on MakerAssetAmount (type Uint256/NUMERIC)
		{
			name: "MakerAssetAmount = 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Equal,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "MakerAssetAmount != 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 5), safeSubsliceOrders(originalOrders, 6, 10)...),
		},
		{
			name: "MakerAssetAmount < 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Less,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[:5],
		},
		{
			name: "MakerAssetAmount > 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Greater,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[6:],
		},
		{
			name: "MakerAssetAmount <= 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  LessOrEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[:6],
		},
		{
			name: "MakerAssetAmount >= 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[5:],
		},
		{
			name: "MakerAssetAmount < 10^76",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Less,
					Value: NewUint256(math.BigPow(10, 76)),
				},
			},
			expectedOrders: originalOrders,
		},

		// Filter on MakerAssetData (type []byte/TEXT)
		{
			name: "MakerAssetData = f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  Equal,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "MakerAssetData != f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  NotEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 5), safeSubsliceOrders(originalOrders, 6, 10)...),
		},
		{
			name: "MakerAssetData < f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  Less,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[:5],
		},
		{
			name: "MakerAssetData > f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  Greater,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[6:],
		},
		{
			name: "MakerAssetData <= f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  LessOrEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[:6],
		},
		{
			name: "MakerAssetData >= f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  GreaterOrEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[5:],
		},

		// Filter on ParsedMakerAssetData (type ParsedAssetData/TEXT)
		{
			name: "ParsedMakerAssetData CONTAINS query that matches all",
			filters: []OrderFilter{
				{
					Field: OFParsedMakerAssetData,
					Kind:  Contains,
					Value: fmt.Sprintf(`"address":"%s","tokenID":"0"`, strings.ToLower(constants.GanacheDummyERC721TokenAddress.Hex())),
				},
			},
			expectedOrders: originalOrders,
		},
		{
			name: "ParsedMakerAssetData CONTAINS with helper method query that matches one",
			filters: []OrderFilter{
				{
					Field: OFParsedMakerAssetData,
					Kind:  Contains,
					Value: fmt.Sprintf(`"address":"%s","tokenID":"51"`, strings.ToLower(constants.GanacheDummyERC721TokenAddress.Hex())),
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "ParsedMakerAssetData CONTAINS with helper method query that matches all",
			filters: []OrderFilter{
				IncludesMakerAssetData(constants.GanacheDummyERC721TokenAddress, big.NewInt(0)),
			},
			expectedOrders: originalOrders,
		},
		{
			name: "ParsedMakerAssetData CONTAINS with helper method query that matches one",
			filters: []OrderFilter{
				IncludesMakerAssetData(constants.GanacheDummyERC721TokenAddress, big.NewInt(51)),
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "ParsedMakerFeeAssetData CONTAINS with helper method query that matches all",
			filters: []OrderFilter{
				IncludesMakerFeeAssetData(constants.GanacheDummyERC1155MintableAddress, big.NewInt(567)),
			},
			expectedOrders: originalOrders,
		},

		// Combining two or more filters
		{
			name: "MakerAssetAmount >= 3 AND MakerAssetData < h",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: OFMakerAssetData,
					Kind:  Less,
					Value: []byte("h"),
				},
			},
			expectedOrders: originalOrders[3:7],
		},
		{
			name: "MakerAssetAmount >= 3 AND MakerAssetData < h AND TakerAssetAmount != 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: OFMakerAssetData,
					Kind:  Less,
					Value: []byte("h"),
				},
				{
					Field: OFTakerAssetAmount,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 3, 5), safeSubsliceOrders(originalOrders, 6, 7)...),
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("%s (test case %d)", testCase.name, i)
		t.Run(testCaseName, runFindOrdersFilterTestCase(t, db, testCase))
	}
}

type findOrdersFilterTestCase struct {
	name           string
	filters        []OrderFilter
	expectedOrders []*Order
	expectedError  string
}

func runFindOrdersFilterTestCase(t *testing.T, db *DB, testCase findOrdersFilterTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		findOpts := &FindOrdersOpts{
			Filters: testCase.filters,
		}

		foundOrders, err := db.FindOrders(findOpts)
		if testCase.expectedError != "" {
			require.Error(t, err, "expected an error but got nil")
			assert.Contains(t, err.Error(), testCase.expectedError, "wrong error message")
		} else {
			require.NoError(t, err)
			assertOrderSlicesAreUnsortedEqual(t, testCase.expectedOrders, foundOrders)
		}
	}
}

func TestDeleteOrdersFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	// Create some test orders with very specific characteristics to make it easier to write tests.
	// - Both MakerAssetAmount and TakerAssetAmount will be 0, 1, 2, etc.
	// - MakerAssetData will be 'a', 'b', 'c', etc.
	// - ParsedMakerAssetData will always be for the ERC721Dummy contract, and each will contain
	//   two token ids: (0, 1), (0, 11), (0, 21), (0, 31) etc.
	numOrders := 10
	originalOrders := []*Order{}
	for i := 0; i < numOrders; i++ {
		order := newTestOrder()
		order.MakerAssetAmount = NewUint256(big.NewInt(int64(i)))
		order.TakerAssetAmount = NewUint256(big.NewInt(int64(i)))
		order.MakerAssetData = []byte{97 + byte(i)}
		parsedMakerAssetData := ParsedAssetData([]SingleAssetData{
			{
				Address: constants.GanacheDummyERC721TokenAddress,
				TokenID: NewUint256(big.NewInt(0)),
			},
			{
				Address: constants.GanacheDummyERC721TokenAddress,
				TokenID: NewUint256(big.NewInt(int64(i)*10 + 1)),
			},
		})
		order.ParsedMakerAssetData = &parsedMakerAssetData
		originalOrders = append(originalOrders, order)
	}
	_, _, err := db.AddOrders(originalOrders)
	require.NoError(t, err)

	testCases := []findOrdersFilterTestCase{
		{
			name:           "no filter",
			filters:        []OrderFilter{},
			expectedOrders: []*Order{},
		},

		// Filter on MakerAssetAmount (type Uint256/NUMERIC)
		{
			name: "MakerAssetAmount = 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Equal,
					Value: 5,
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 5), safeSubsliceOrders(originalOrders, 6, 10)...),
		},
		{
			name: "MakerAssetAmount != 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "MakerAssetAmount < 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Less,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[5:],
		},
		{
			name: "MakerAssetAmount > 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Greater,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[:6],
		},
		{
			name: "MakerAssetAmount <= 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  LessOrEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[6:],
		},
		{
			name: "MakerAssetAmount >= 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[:5],
		},
		{
			name: "MakerAssetAmount < 10^76",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  Less,
					Value: NewUint256(math.BigPow(10, 76)),
				},
			},
			expectedOrders: []*Order{},
		},

		// Filter on MakerAssetData (type []byte/TEXT)
		{
			name: "MakerAssetData = f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  Equal,
					Value: []byte("f"),
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 5), safeSubsliceOrders(originalOrders, 6, 10)...),
		},
		{
			name: "MakerAssetData != f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  NotEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "MakerAssetData < f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  Less,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[5:],
		},
		{
			name: "MakerAssetData > f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  Greater,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[:6],
		},
		{
			name: "MakerAssetData <= f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  LessOrEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[6:],
		},
		{
			name: "MakerAssetData >= f",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetData,
					Kind:  GreaterOrEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[:5],
		},

		// Filter on ParsedMakerAssetData (type ParsedAssetData/TEXT)
		{
			name: "ParsedMakerAssetData CONTAINS query that matches all",
			filters: []OrderFilter{
				{
					Field: OFParsedMakerAssetData,
					Kind:  Contains,
					Value: fmt.Sprintf(`"address":"%s","tokenID":"0"`, strings.ToLower(constants.GanacheDummyERC721TokenAddress.Hex())),
				},
			},
			expectedOrders: []*Order{},
		},
		{
			name: "ParsedMakerAssetData CONTAINS with helper method query that matches one",
			filters: []OrderFilter{
				{
					Field: OFParsedMakerAssetData,
					Kind:  Contains,
					Value: fmt.Sprintf(`"address":"%s","tokenID":"51"`, strings.ToLower(constants.GanacheDummyERC721TokenAddress.Hex())),
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 5), safeSubsliceOrders(originalOrders, 6, 10)...),
		},
		{
			name: "ParsedMakerAssetData CONTAINS with helper method query that matches all",
			filters: []OrderFilter{
				IncludesMakerAssetData(constants.GanacheDummyERC721TokenAddress, big.NewInt(0)),
			},
			expectedOrders: []*Order{},
		},
		{
			name: "ParsedMakerAssetData CONTAINS with helper method query that matches one",
			filters: []OrderFilter{
				IncludesMakerAssetData(constants.GanacheDummyERC721TokenAddress, big.NewInt(51)),
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 5), safeSubsliceOrders(originalOrders, 6, 10)...),
		},
		{
			name: "ParsedMakerFeeAssetData CONTAINS with helper method query that matches all",
			filters: []OrderFilter{
				IncludesMakerFeeAssetData(constants.GanacheDummyERC1155MintableAddress, big.NewInt(567)),
			},
			expectedOrders: []*Order{},
		},

		// Combining two or more filters
		{
			name: "MakerAssetAmount >= 3 AND MakerAssetData < h",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: OFMakerAssetData,
					Kind:  Less,
					Value: []byte("h"),
				},
			},
			expectedOrders: append(safeSubsliceOrders(originalOrders, 0, 3), safeSubsliceOrders(originalOrders, 7, 10)...),
		},
		{
			name: "MakerAssetAmount >= 3 AND MakerAssetData < h AND TakerAssetAmount != 5",
			filters: []OrderFilter{
				{
					Field: OFMakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: OFMakerAssetData,
					Kind:  Less,
					Value: []byte("h"),
				},
				{
					Field: OFTakerAssetAmount,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedOrders: []*Order{
				originalOrders[0],
				originalOrders[1],
				originalOrders[2],
				originalOrders[5],
				originalOrders[7],
				originalOrders[8],
				originalOrders[9],
			},
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("%s (test case %d)", testCase.name, i)
		t.Run(testCaseName, runDeleteOrdersFilterTestCase(t, db, originalOrders, testCase))
	}
}

type deleteOrdersFilterTestCase struct {
	name           string
	filters        []OrderFilter
	expectedOrders []*Order
	expectedError  string
}

func runDeleteOrdersFilterTestCase(t *testing.T, db *DB, originalOrders []*Order, testCase findOrdersFilterTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		defer func() {
			// After each case, reset the state of the database by re-adding the original orders.
			_, _, err := db.AddOrders(originalOrders)
			require.NoError(t, err)
		}()

		deleteOpts := &DeleteOrdersOpts{
			Filters: testCase.filters,
		}

		err := db.DeleteOrders(deleteOpts)
		if testCase.expectedError != "" {
			require.Error(t, err, "expected an error but got nil")
			assert.Contains(t, err.Error(), testCase.expectedError, "wrong error message")
		} else {
			require.NoError(t, err)
			foundOrders, err := db.FindOrders(nil)
			require.NoError(t, err)
			assertOrderSlicesAreUnsortedEqual(t, testCase.expectedOrders, foundOrders)
		}
	}
}

func TestAddMiniHeaders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numMiniHeaders := 10
	miniHeaders := []*MiniHeader{}
	for i := 0; i < numMiniHeaders; i++ {
		miniHeaders = append(miniHeaders, newTestMiniHeader())
	}

	{
		added, removed, err := db.AddMiniHeaders(miniHeaders)
		require.NoError(t, err)
		assert.Len(t, removed, 0, "Expected no miniHeaders to be removed")
		assertMiniHeaderSlicesAreUnsortedEqual(t, miniHeaders, added)
	}
	{
		added, removed, err := db.AddMiniHeaders(miniHeaders)
		require.NoError(t, err)
		assert.Len(t, removed, 0, "Expected no miniHeaders to be removed")
		assert.Len(t, added, 0, "Expected no miniHeaders to be added (they should already exist)")
	}
}

func TestGetMiniHeader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	added, _, err := db.AddMiniHeaders([]*MiniHeader{newTestMiniHeader()})
	require.NoError(t, err)
	originalMiniHeader := added[0]

	foundMiniHeader, err := db.GetMiniHeader(originalMiniHeader.Hash)
	require.NoError(t, err)
	assertMiniHeadersAreEqual(t, originalMiniHeader, foundMiniHeader)
}

func TestFindMiniHeaders(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numMiniHeaders := 10
	originalMiniHeaders := []*MiniHeader{}
	for i := 0; i < numMiniHeaders; i++ {
		originalMiniHeaders = append(originalMiniHeaders, newTestMiniHeader())
	}
	_, _, err := db.AddMiniHeaders(originalMiniHeaders)
	require.NoError(t, err)

	foundMiniHeaders, err := db.FindMiniHeaders(nil)
	require.NoError(t, err)
	assertMiniHeaderSlicesAreUnsortedEqual(t, originalMiniHeaders, foundMiniHeaders)
}

func TestFindMiniHeadersSort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	// Create some test miniHeaders with carefully chosen Number and Timestamp
	// values for testing sorting.
	numMiniHeaders := 5
	originalMiniHeaders := []*MiniHeader{}
	for i := 0; i < numMiniHeaders; i++ {
		miniHeader := newTestMiniHeader()
		miniHeader.Number = NewUint256(big.NewInt(int64(i)))
		// It's important for some miniHeaders to have the same Timestamp
		// so that we can test secondary sorts (sorting on more than one
		// field).
		if i%2 == 0 {
			miniHeader.Timestamp = time.Unix(717793653, 0)
		} else {
			miniHeader.Timestamp = time.Unix(1588194484, 0)
		}
		originalMiniHeaders = append(originalMiniHeaders, miniHeader)
	}
	_, _, err := db.AddMiniHeaders(originalMiniHeaders)
	require.NoError(t, err)

	// TODO(albrow): Add test cases.
	testCases := []findMiniHeadersSortTestCase{
		{
			sortOpts: []MiniHeaderSort{
				{
					Field:     MFNumber,
					Direction: Ascending,
				},
			},
			less: lessByNumberAsc,
		},
		{
			sortOpts: []MiniHeaderSort{
				{
					Field:     MFNumber,
					Direction: Descending,
				},
			},
			less: lessByNumberDesc,
		},
		{
			sortOpts: []MiniHeaderSort{
				{
					Field:     MFTimestamp,
					Direction: Ascending,
				},
				{
					Field:     MFNumber,
					Direction: Ascending,
				},
			},
			less: lessByTimestampAscAndNumberAsc,
		},
		{
			sortOpts: []MiniHeaderSort{
				{
					Field:     MFTimestamp,
					Direction: Descending,
				},
				{
					Field:     MFNumber,
					Direction: Descending,
				},
			},
			less: lessByTimestampDescAndNumberDesc,
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("test case %d", i)
		t.Run(testCaseName, runFindMiniHeadersSortTestCase(t, db, originalMiniHeaders, testCase))
	}
}

type findMiniHeadersSortTestCase struct {
	sortOpts []MiniHeaderSort
	less     func([]*MiniHeader) func(i, j int) bool
}

func runFindMiniHeadersSortTestCase(t *testing.T, db *DB, originalMiniHeaders []*MiniHeader, testCase findMiniHeadersSortTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		expectedMiniHeaders := make([]*MiniHeader, len(originalMiniHeaders))
		copy(expectedMiniHeaders, originalMiniHeaders)
		sort.Slice(expectedMiniHeaders, testCase.less(expectedMiniHeaders))
		findOpts := &FindMiniHeadersOpts{
			Sort: testCase.sortOpts,
		}
		foundMiniHeaders, err := db.FindMiniHeaders(findOpts)
		require.NoError(t, err)
		assertMiniHeaderSlicesAreEqual(t, expectedMiniHeaders, foundMiniHeaders)
	}
}

func TestFindMiniHeadersLimitAndOffset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	numMiniHeaders := 10
	originalMiniHeaders := []*MiniHeader{}
	for i := 0; i < numMiniHeaders; i++ {
		originalMiniHeaders = append(originalMiniHeaders, newTestMiniHeader())
	}
	_, _, err := db.AddMiniHeaders(originalMiniHeaders)
	require.NoError(t, err)
	sortMiniHeadersByHash(originalMiniHeaders)

	testCases := []findMiniHeadersLimitAndOffsetTestCase{
		{
			limit:               0,
			offset:              0,
			expectedMiniHeaders: originalMiniHeaders,
		},
		{
			limit:               3,
			offset:              0,
			expectedMiniHeaders: originalMiniHeaders[:3],
		},
		{
			limit:         0,
			offset:        3,
			expectedError: "can't use Offset without Limit",
		},
		{
			limit:               10,
			offset:              3,
			expectedMiniHeaders: originalMiniHeaders[3:],
		},
		{
			limit:               4,
			offset:              3,
			expectedMiniHeaders: originalMiniHeaders[3:7],
		},
		{
			limit:               10,
			offset:              10,
			expectedMiniHeaders: []*MiniHeader{},
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("test case %d", i)
		t.Run(testCaseName, runFindMiniHeadersLimitAndOffsetTestCase(t, db, originalMiniHeaders, testCase))
	}
}

type findMiniHeadersLimitAndOffsetTestCase struct {
	limit               uint
	offset              uint
	expectedMiniHeaders []*MiniHeader
	expectedError       string
}

func runFindMiniHeadersLimitAndOffsetTestCase(t *testing.T, db *DB, originalMiniHeaders []*MiniHeader, testCase findMiniHeadersLimitAndOffsetTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		findOpts := &FindMiniHeadersOpts{
			Sort: []MiniHeaderSort{
				{
					Field:     MFHash,
					Direction: Ascending,
				},
			},
			Limit:  testCase.limit,
			Offset: testCase.offset,
		}

		foundMiniHeaders, err := db.FindMiniHeaders(findOpts)
		if testCase.expectedError != "" {
			require.Error(t, err, "expected an error but got nil")
			assert.Contains(t, err.Error(), testCase.expectedError, "wrong error message")
		} else {
			require.NoError(t, err)
			assertMiniHeaderSlicesAreEqual(t, testCase.expectedMiniHeaders, foundMiniHeaders)
		}
	}
}

func TestFindMiniHeadersFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestDB(t, ctx)

	// Create some test miniheaders with very specific characteristics to make it easier to write tests.
	// - Number will be 0, 1, 2, etc.
	// - Timestamp will be 0, 100, 200, etc. seconds since Unix Epoch
	// - Each log in Logs will have BlockNumber set to 0, 1, 2, etc.
	numMiniHeaders := 10
	originalMiniHeaders := []*MiniHeader{}
	for i := 0; i < numMiniHeaders; i++ {
		miniHeader := newTestMiniHeader()
		miniHeader.Number = NewUint256(big.NewInt(int64(i)))
		miniHeader.Timestamp = time.Unix(int64(i)*100, 0)
		for i := range miniHeader.Logs.Logs {
			miniHeader.Logs.Logs[i].BlockNumber = miniHeader.Number.Uint64()
		}
		originalMiniHeaders = append(originalMiniHeaders, miniHeader)
	}
	_, _, err := db.AddMiniHeaders(originalMiniHeaders)
	require.NoError(t, err)

	testCases := []findMiniHeadersFilterTestCase{
		{
			name:                "no filter",
			filters:             []MiniHeaderFilter{},
			expectedMiniHeaders: originalMiniHeaders,
		},

		// Filter on Number (type Uint256/NUMERIC)
		{
			name: "Number = 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  Equal,
					Value: 5,
				},
			},
			expectedMiniHeaders: originalMiniHeaders[5:6],
		},
		{
			name: "Number != 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedMiniHeaders: append(safeSubsliceMiniHeaders(originalMiniHeaders, 0, 5), safeSubsliceMiniHeaders(originalMiniHeaders, 6, 10)...),
		},
		{
			name: "Number < 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  Less,
					Value: 5,
				},
			},
			expectedMiniHeaders: originalMiniHeaders[:5],
		},
		{
			name: "Number > 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  Greater,
					Value: 5,
				},
			},
			expectedMiniHeaders: originalMiniHeaders[6:],
		},
		{
			name: "Number <= 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  LessOrEqual,
					Value: 5,
				},
			},
			expectedMiniHeaders: originalMiniHeaders[:6],
		},
		{
			name: "Number >= 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  GreaterOrEqual,
					Value: 5,
				},
			},
			expectedMiniHeaders: originalMiniHeaders[5:],
		},
		{
			name: "Number < 10^76",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  Less,
					Value: NewUint256(math.BigPow(10, 76)),
				},
			},
			expectedMiniHeaders: originalMiniHeaders,
		},

		// Filter on Timestamp (type time.Time/TIMESTAMP)
		{
			name: "Timestamp = 500",
			filters: []MiniHeaderFilter{
				{
					Field: MFTimestamp,
					Kind:  Equal,
					Value: time.Unix(500, 0),
				},
			},
			expectedMiniHeaders: originalMiniHeaders[5:6],
		},
		{
			name: "Timestamp != 500",
			filters: []MiniHeaderFilter{
				{
					Field: MFTimestamp,
					Kind:  NotEqual,
					Value: time.Unix(500, 0),
				},
			},
			expectedMiniHeaders: append(safeSubsliceMiniHeaders(originalMiniHeaders, 0, 5), safeSubsliceMiniHeaders(originalMiniHeaders, 6, 10)...),
		},
		{
			name: "Timestamp < 500",
			filters: []MiniHeaderFilter{
				{
					Field: MFTimestamp,
					Kind:  Less,
					Value: time.Unix(500, 0),
				},
			},
			expectedMiniHeaders: originalMiniHeaders[:5],
		},
		{
			name: "Timestamp > 500",
			filters: []MiniHeaderFilter{
				{
					Field: MFTimestamp,
					Kind:  Greater,
					Value: time.Unix(500, 0),
				},
			},
			expectedMiniHeaders: originalMiniHeaders[6:],
		},
		{
			name: "Timestamp <= 500",
			filters: []MiniHeaderFilter{
				{
					Field: MFTimestamp,
					Kind:  LessOrEqual,
					Value: time.Unix(500, 0),
				},
			},
			expectedMiniHeaders: originalMiniHeaders[:6],
		},
		{
			name: "Timestamp >= 500",
			filters: []MiniHeaderFilter{
				{
					Field: MFTimestamp,
					Kind:  GreaterOrEqual,
					Value: time.Unix(500, 0),
				},
			},
			expectedMiniHeaders: originalMiniHeaders[5:],
		},

		// Filter on Logs (type ParsedAssetData/TEXT)
		{
			name: "Logs CONTAINS query that matches all",
			filters: []MiniHeaderFilter{
				{
					Field: MFLogs,
					Kind:  Contains,
					Value: `"address":"0x21ab6c9fac80c59d401b37cb43f81ea9dde7fe34"`,
				},
			},
			expectedMiniHeaders: originalMiniHeaders,
		},
		{
			name: "Logs CONTAINS query that matches one",
			filters: []MiniHeaderFilter{
				{
					Field: MFLogs,
					Kind:  Contains,
					Value: `"blockNumber":"0x5"`,
				},
			},
			expectedMiniHeaders: originalMiniHeaders[5:6],
		},

		// Combining two or more filters
		{
			name: "Number >= 3 AND Timestamp < h",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: MFTimestamp,
					Kind:  Less,
					Value: time.Unix(700, 0),
				},
			},
			expectedMiniHeaders: originalMiniHeaders[3:7],
		},
		{
			name: "Number >= 3 AND Timestamp < 700 AND Number != 5",
			filters: []MiniHeaderFilter{
				{
					Field: MFNumber,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: MFTimestamp,
					Kind:  Less,
					Value: time.Unix(700, 0),
				},
				{
					Field: MFNumber,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedMiniHeaders: append(safeSubsliceMiniHeaders(originalMiniHeaders, 3, 5), safeSubsliceMiniHeaders(originalMiniHeaders, 6, 7)...),
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("%s (test case %d)", testCase.name, i)
		t.Run(testCaseName, runFindMiniHeadersFilterTestCase(t, db, testCase))
	}
}

type findMiniHeadersFilterTestCase struct {
	name                string
	filters             []MiniHeaderFilter
	expectedMiniHeaders []*MiniHeader
	expectedError       string
}

func runFindMiniHeadersFilterTestCase(t *testing.T, db *DB, testCase findMiniHeadersFilterTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		findOpts := &FindMiniHeadersOpts{
			Filters: testCase.filters,
		}

		foundMiniHeaders, err := db.FindMiniHeaders(findOpts)
		if testCase.expectedError != "" {
			require.Error(t, err, "expected an error but got nil")
			assert.Contains(t, err.Error(), testCase.expectedError, "wrong error message")
		} else {
			require.NoError(t, err)
			assertMiniHeaderSlicesAreUnsortedEqual(t, testCase.expectedMiniHeaders, foundMiniHeaders)
		}
	}
}

func TestParseContractAddressesAndTokenIdsFromAssetData(t *testing.T) {
	// ERC20 AssetData
	erc20AssetData := common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32")
	parsedAssetData, err := ParseContractAddressesAndTokenIdsFromAssetData(erc20AssetData, contractAddresses)
	require.NoError(t, err)
	assert.Len(t, parsedAssetData, 1)
	expectedAddress := common.HexToAddress("0x38ae374ecf4db50b0ff37125b591a04997106a32")
	assert.Equal(t, expectedAddress, parsedAssetData[0].Address)
	var expectedTokenID *Uint256 = nil
	assert.Equal(t, expectedTokenID, parsedAssetData[0].TokenID)

	// ERC721 AssetData
	erc721AssetData := common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001")
	parsedAssetData, err = ParseContractAddressesAndTokenIdsFromAssetData(erc721AssetData, contractAddresses)
	require.NoError(t, err)
	assert.Equal(t, 1, len(parsedAssetData))
	expectedAddress = common.HexToAddress("0x1dC4c1cEFEF38a777b15aA20260a54E584b16C48")
	assert.Equal(t, expectedAddress, parsedAssetData[0].Address)
	expectedTokenID = NewUint256(big.NewInt(1))
	assert.Equal(t, expectedTokenID, parsedAssetData[0].TokenID)

	// Multi AssetData
	multiAssetData := common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000x94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000")
	parsedAssetData, err = ParseContractAddressesAndTokenIdsFromAssetData(multiAssetData, contractAddresses)
	require.NoError(t, err)
	assert.Equal(t, 2, len(parsedAssetData))
	expectedParsedAssetData := []SingleAssetData{
		{
			Address: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		},
		{
			Address: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
			TokenID: NewUint256(big.NewInt(1)),
		},
	}
	for i, singleAssetData := range parsedAssetData {
		expectedSingleAssetData := expectedParsedAssetData[i]
		assert.Equal(t, expectedSingleAssetData.Address, singleAssetData.Address)
		assert.Equal(t, expectedSingleAssetData.TokenID, singleAssetData.TokenID)
	}
}

func newTestDB(t *testing.T, ctx context.Context) *DB {
	db, err := New(ctx, filepath.Join("tmp", "db_testing", uuid.New().String()))
	require.NoError(t, err)
	require.NoError(t, db.migrate())
	return db
}

// newTestOrder returns a new order with a random hash that is ready to insert
// into the database. Some computed fields (e.g. hash, signature) may not be
// correct, so the order will not pass 0x validation.
func newTestOrder() *Order {
	parsedMakerAssetData := ParsedAssetData([]SingleAssetData{
		{
			Address: constants.GanacheDummyERC721TokenAddress,
			TokenID: NewUint256(big.NewInt(10)),
		},
		{
			Address: constants.GanacheDummyERC721TokenAddress,
			TokenID: NewUint256(big.NewInt(20)),
		},
		{
			Address: constants.GanacheDummyERC721TokenAddress,
			TokenID: NewUint256(big.NewInt(30)),
		},
	})
	parsedMakerFeeAssetData := ParsedAssetData([]SingleAssetData{
		{
			Address: constants.GanacheDummyERC1155MintableAddress,
			TokenID: NewUint256(big.NewInt(567)),
		},
	})
	return &Order{
		Hash:                     common.BigToHash(big.NewInt(int64(rand.Int()))),
		ChainID:                  NewUint256(big.NewInt(constants.TestChainID)),
		MakerAddress:             constants.GanacheAccount1,
		TakerAddress:             constants.NullAddress,
		SenderAddress:            constants.NullAddress,
		FeeRecipientAddress:      constants.NullAddress,
		MakerAssetData:           constants.ZRXAssetData,
		MakerFeeAssetData:        constants.NullBytes,
		TakerAssetData:           constants.WETHAssetData,
		TakerFeeAssetData:        constants.NullBytes,
		Salt:                     NewUint256(big.NewInt(int64(time.Now().Nanosecond()))),
		MakerFee:                 NewUint256(big.NewInt(0)),
		TakerFee:                 NewUint256(big.NewInt(0)),
		MakerAssetAmount:         NewUint256(big.NewInt(100)),
		TakerAssetAmount:         NewUint256(big.NewInt(42)),
		ExpirationTimeSeconds:    NewUint256(big.NewInt(time.Now().Add(24 * time.Hour).Unix())),
		ExchangeAddress:          contractAddresses.Exchange,
		Signature:                []byte{1, 2, 255, 255},
		LastUpdated:              time.Now(),
		FillableTakerAssetAmount: NewUint256(big.NewInt(42)),
		IsRemoved:                false,
		IsPinned:                 true,
		ParsedMakerAssetData:     &parsedMakerAssetData,
		ParsedMakerFeeAssetData:  &parsedMakerFeeAssetData,
	}
}

func newTestMiniHeader() *MiniHeader {
	return &MiniHeader{
		Hash:      common.BigToHash(big.NewInt(int64(rand.Int()))),
		Parent:    common.BigToHash(big.NewInt(int64(rand.Int()))),
		Number:    NewUint256(big.NewInt(int64(rand.Int()))),
		Timestamp: time.Now(),
		Logs:      newTestEventLogs(),
	}
}

func newTestEventLogs() *EventLogs {
	return NewEventLogs([]types.Log{
		{
			Address: common.HexToAddress("0x21ab6c9fac80c59d401b37cb43f81ea9dde7fe34"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x0000000000000000000000004d8a4aa1f304f9632cf3877473445d85c577fe5d"),
				common.HexToHash("0x0000000000000000000000004bdd0d16cfa18e33860470fc4d65c6f5cee60959"),
			},
			Data:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000337ad34c0"),
			BlockNumber: 30,
			TxHash:      common.HexToHash("0xd9bb5f9e888ee6f74bedcda811c2461230f247c205849d6f83cb6c3925e54586"),
			TxIndex:     0,
			BlockHash:   common.HexToHash("0x6bbf9b6e836207ab25379c20e517a89090cbbaf8877746f6ed7fb6820770816b"),
			Index:       0,
			Removed:     false,
		},
		{
			Address: common.HexToAddress("0x21ab6c9fac80c59d401b37cb43f81ea9dde7fe34"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x0000000000000000000000004d8a4aa1f304f9632cf3877473445d85c577fe5d"),
				common.HexToHash("0x0000000000000000000000004bdd0d16cfa18e33860470fc4d65c6f5cee60959"),
			},
			Data:        common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000deadbeef"),
			BlockNumber: 31,
			TxHash:      common.HexToHash("0xd9bb5f9e888ee6f74bedcda811c2461230f247c205849d6f83cb6c3925e54586"),
			TxIndex:     1,
			BlockHash:   common.HexToHash("0x6bbf9b6e836207ab25379c20e517a89090cbbaf8877746f6ed7fb6820770816b"),
			Index:       2,
			Removed:     true,
		},
	})
}

// safeSubsliceOrders returns a (shallow) subslice of orders without modifying
// the original slice. Uses the same semantics as slice expressions: low is
// inclusive, hi is exclusive. The returned slice still contains pointers, it
// just doesn't use the same underlying array.
func safeSubsliceOrders(orders []*Order, low, hi int) []*Order {
	result := make([]*Order, hi-low)
	for i := low; i < hi; i++ {
		result[i-low] = orders[i]
	}
	return result
}

func sortOrdersByHash(orders []*Order) {
	sort.SliceStable(orders, func(i, j int) bool {
		return bytes.Compare(orders[i].Hash.Bytes(), orders[j].Hash.Bytes()) == -1
	})
}

func lessByMakerAssetAmountAsc(orders []*Order) func(i, j int) bool {
	return func(i, j int) bool {
		return orders[i].MakerAssetAmount.Cmp(orders[j].MakerAssetAmount.Int) == -1
	}
}

func lessByMakerAssetAmountDesc(orders []*Order) func(i, j int) bool {
	return func(i, j int) bool {
		return orders[i].MakerAssetAmount.Cmp(orders[j].MakerAssetAmount.Int) == 1
	}
}

func lessByTakerAssetAmountAscAndMakerAssetAmountAsc(orders []*Order) func(i, j int) bool {
	return func(i, j int) bool {
		switch orders[i].TakerAssetAmount.Cmp(orders[j].TakerAssetAmount.Int) {
		case -1:
			// Less
			return true
		case 1:
			// Greater
			return false
		default:
			// Equal. In this case we use MakerAssetAmount as a secondary sort
			// (i.e. a tie-breaker)
			return orders[i].MakerAssetAmount.Cmp(orders[j].MakerAssetAmount.Int) == -1
		}
	}
}

func lessByTakerAssetAmountDescAndMakerAssetAmountDesc(orders []*Order) func(i, j int) bool {
	return func(i, j int) bool {
		switch orders[i].TakerAssetAmount.Cmp(orders[j].TakerAssetAmount.Int) {
		case -1:
			// Less
			return false
		case 1:
			// Greater
			return true
		default:
			// Equal. In this case we use MakerAssetAmount as a secondary sort
			// (i.e. a tie-breaker)
			return orders[i].MakerAssetAmount.Cmp(orders[j].MakerAssetAmount.Int) == 1
		}
	}
}

func assertOrderSlicesAreEqual(t *testing.T, expected, actual []*Order) {
	assert.Equal(t, len(expected), len(actual), "wrong number of orders")
	for i, expectedOrder := range expected {
		if i >= len(actual) {
			break
		}
		actualOrder := actual[i]
		assertOrdersAreEqual(t, *expectedOrder, *actualOrder)
	}
	if t.Failed() {
		expectedJSON, err := json.MarshalIndent(expected, "", "  ")
		require.NoError(t, err)
		actualJSON, err := json.MarshalIndent(actual, "", "  ")
		require.NoError(t, err)
		t.Logf("\nexpected:\n%s\n\n", string(expectedJSON))
		t.Logf("\nactual:\n%s\n\n", string(actualJSON))
		assert.Equal(t, string(expectedJSON), string(actualJSON))
	}
}

func assertOrderSlicesAreUnsortedEqual(t *testing.T, expected, actual []*Order) {
	// Make a copy of the given orders so we don't mess up the original when sorting them.
	expectedCopy := make([]*Order, len(expected))
	copy(expectedCopy, expected)
	sortOrdersByHash(expectedCopy)
	actualCopy := make([]*Order, len(actual))
	copy(actualCopy, actual)
	sortOrdersByHash(actualCopy)
	assertOrderSlicesAreEqual(t, expectedCopy, actualCopy)
}

func assertOrdersAreEqual(t *testing.T, expected, actual Order) {
	if expected.LastUpdated.Equal(actual.LastUpdated) {
		// HACK(albrow): In this case, the two values represent the same time.
		// This is what we care about, but the assert package might consider
		// them unequal if some internal fields are different (there are
		// different ways of representing the same time). As a workaround,
		// we manually set actual.LastUpdated.
		actual.LastUpdated = expected.LastUpdated
	} else {
		assert.Equal(t, expected.LastUpdated, actual.LastUpdated, "order.LastUpdated was not equal")
	}
	// We can compare the rest of the fields normally.
	assert.Equal(t, expected, actual)
}

// safeSubsliceMiniHeaders returns a (shallow) subslice of mini headers without
// modifying the original slice. Uses the same semantics as slice expressions:
// low is inclusive, hi is exclusive. The returned slice still contains
// pointers, it just doesn't use the same underlying array.
func safeSubsliceMiniHeaders(miniHeaders []*MiniHeader, low, hi int) []*MiniHeader {
	result := make([]*MiniHeader, hi-low)
	for i := low; i < hi; i++ {
		result[i-low] = miniHeaders[i]
	}
	return result
}

func sortMiniHeadersByHash(miniHeaders []*MiniHeader) {
	sort.SliceStable(miniHeaders, func(i, j int) bool {
		return bytes.Compare(miniHeaders[i].Hash.Bytes(), miniHeaders[j].Hash.Bytes()) == -1
	})
}

func lessByNumberAsc(miniHeaders []*MiniHeader) func(i, j int) bool {
	return func(i, j int) bool {
		return miniHeaders[i].Number.Cmp(miniHeaders[j].Number.Int) == -1
	}
}

func lessByNumberDesc(miniHeaders []*MiniHeader) func(i, j int) bool {
	return func(i, j int) bool {
		return miniHeaders[i].Number.Cmp(miniHeaders[j].Number.Int) == 1
	}
}

func lessByTimestampAscAndNumberAsc(miniHeaders []*MiniHeader) func(i, j int) bool {
	return func(i, j int) bool {
		switch {
		case miniHeaders[i].Timestamp.Before(miniHeaders[j].Timestamp):
			// Less
			return true
		case miniHeaders[i].Timestamp.After(miniHeaders[j].Timestamp):
			// Greater
			return false
		default:
			// Equal. In this case we use Number as a secondary sort
			// (i.e. a tie-breaker)
			return miniHeaders[i].Number.Cmp(miniHeaders[j].Number.Int) == -1
		}
	}
}

func lessByTimestampDescAndNumberDesc(miniHeaders []*MiniHeader) func(i, j int) bool {
	return func(i, j int) bool {
		switch {
		case miniHeaders[i].Timestamp.Before(miniHeaders[j].Timestamp):
			// Less
			return false
		case miniHeaders[i].Timestamp.After(miniHeaders[j].Timestamp):
			// Greater
			return true
		default:
			// Equal. In this case we use Number as a secondary sort
			// (i.e. a tie-breaker)
			return miniHeaders[i].Number.Cmp(miniHeaders[j].Number.Int) == 1
		}
	}
}

func assertMiniHeaderSlicesAreEqual(t *testing.T, expected, actual []*MiniHeader) {
	assert.Len(t, actual, len(expected), "wrong number of miniheaders")
	for i, expectedMiniHeader := range expected {
		if i >= len(actual) {
			break
		}
		actualMiniHeader := expected[i]
		assertMiniHeadersAreEqual(t, expectedMiniHeader, actualMiniHeader)
	}
	if t.Failed() {
		expectedJSON, err := json.MarshalIndent(expected, "", "  ")
		require.NoError(t, err)
		actualJSON, err := json.MarshalIndent(actual, "", "  ")
		require.NoError(t, err)
		t.Logf("\nexpected:\n%s\n\n", string(expectedJSON))
		t.Logf("\nactual:\n%s\n\n", string(actualJSON))
		assert.Equal(t, string(expectedJSON), string(actualJSON))
	}
}

func assertMiniHeaderSlicesAreUnsortedEqual(t *testing.T, expected, actual []*MiniHeader) {
	// Make a copy of the given mini headers so we don't mess up the original when sorting them.
	expectedCopy := make([]*MiniHeader, len(expected))
	copy(expectedCopy, expected)
	sortMiniHeadersByHash(expectedCopy)
	actualCopy := make([]*MiniHeader, len(actual))
	copy(actualCopy, actual)
	sortMiniHeadersByHash(actualCopy)
	assertMiniHeaderSlicesAreEqual(t, expected, actual)
}

func assertMiniHeadersAreEqual(t *testing.T, expected, actual *MiniHeader) {
	if expected.Timestamp.Equal(actual.Timestamp) {
		// HACK(albrow): In this case, the two values represent the same time.
		// This is what we care about, but the assert package might consider
		// them unequal if some internal fields are different (there are
		// different ways of representing the same time). As a workaround,
		// we manually set actual.Timestamp.
		actual.Timestamp = expected.Timestamp
	} else {
		assert.Equal(t, expected.Timestamp, actual.Timestamp, "miniHeader.Timestamp was not equal")
	}
	// We can compare the rest of the fields normally.
	assert.Equal(t, expected, actual)
}
