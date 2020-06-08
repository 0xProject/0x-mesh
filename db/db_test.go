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
			sortOpts: []OrderSortOpts{
				{
					Field:     MakerAssetAmount,
					Direction: Ascending,
				},
			},
			less: lessByMakerAssetAmountAsc,
		},
		{
			sortOpts: []OrderSortOpts{
				{
					Field:     MakerAssetAmount,
					Direction: Descending,
				},
			},
			less: lessByMakerAssetAmountDesc,
		},
		{
			sortOpts: []OrderSortOpts{
				{
					Field:     TakerAssetAmount,
					Direction: Ascending,
				},
				{
					Field:     MakerAssetAmount,
					Direction: Ascending,
				},
			},
			less: lessByTakerAssetAmountAscAndMakerAssetAmountAsc,
		},
		{
			sortOpts: []OrderSortOpts{
				{
					Field:     TakerAssetAmount,
					Direction: Descending,
				},
				{
					Field:     MakerAssetAmount,
					Direction: Descending,
				},
			},
			less: lessByTakerAssetAmountDescAndMakerAssetAmountDesc,
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("test case %d", i)
		t.Run(testCaseName, runFindOrderSortTestCase(t, db, originalOrders, testCase))
	}
}

type findOrdersSortTestCase struct {
	sortOpts []OrderSortOpts
	less     func([]*Order) func(i, j int) bool
}

func runFindOrderSortTestCase(t *testing.T, db *DB, originalOrders []*Order, testCase findOrdersSortTestCase) func(t *testing.T) {
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
	sortOrders(originalOrders)

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
			Sort: []OrderSortOpts{
				{
					Field:     Hash,
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
	// - Both MakerAssetAmount and TakerAssetAmount will e 0, 1, 2, etc.
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
			filters:        []OrderFilterOpts{},
			expectedOrders: originalOrders,
		},

		// Filter on MakerAssetAmount (type Uint256/NUMERIC)
		{
			name: "MakerAssetAmount = 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  Equal,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "MakerAssetAmount != 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedOrders: append(safeSubslice(originalOrders, 0, 5), safeSubslice(originalOrders, 6, 10)...),
		},
		{
			name: "MakerAssetAmount < 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  Less,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[:5],
		},
		{
			name: "MakerAssetAmount > 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  Greater,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[6:],
		},
		{
			name: "MakerAssetAmount <= 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  LessOrEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[:6],
		},
		{
			name: "MakerAssetAmount >= 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 5,
				},
			},
			expectedOrders: originalOrders[5:],
		},
		{
			name: "MakerAssetAmount < 10^76",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  Less,
					Value: NewUint256(math.BigPow(10, 76)),
				},
			},
			expectedOrders: originalOrders,
		},

		// Filter on MakerAssetData (type []byte/TEXT)
		{
			name: "MakerAssetData = f",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetData,
					Kind:  Equal,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[5:6],
		},
		{
			name: "MakerAssetData != f",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetData,
					Kind:  NotEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: append(safeSubslice(originalOrders, 0, 5), safeSubslice(originalOrders, 6, 10)...),
		},
		{
			name: "MakerAssetData < f",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetData,
					Kind:  Less,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[:5],
		},
		{
			name: "MakerAssetData > f",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetData,
					Kind:  Greater,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[6:],
		},
		{
			name: "MakerAssetData <= f",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetData,
					Kind:  LessOrEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[:6],
		},
		{
			name: "MakerAssetData >= f",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetData,
					Kind:  GreaterOrEqual,
					Value: []byte("f"),
				},
			},
			expectedOrders: originalOrders[5:],
		},

		// Filter on ParsedMakerAssetData (type ParsedAssetData/TEXT)
		{
			name: "ParsedMakerAssetData CONTAINS query that matches all",
			filters: []OrderFilterOpts{
				{
					Field: ParsedMakerAssetData,
					Kind:  Contains,
					Value: fmt.Sprintf(`"address":"%s","tokenID":"0"`, strings.ToLower(constants.GanacheDummyERC721TokenAddress.Hex())),
				},
			},
			expectedOrders: originalOrders,
		},
		{
			name: "ParsedMakerAssetData CONTAINS query that matches one",
			filters: []OrderFilterOpts{
				{
					Field: ParsedMakerAssetData,
					Kind:  Contains,
					Value: fmt.Sprintf(`"address":"%s","tokenID":"51"`, strings.ToLower(constants.GanacheDummyERC721TokenAddress.Hex())),
				},
			},
			expectedOrders: originalOrders[5:6],
		},

		// Combining two or more filters
		{
			name: "MakerAssetAmount >= 3 AND MakerAssetData < h",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: MakerAssetData,
					Kind:  Less,
					Value: []byte("h"),
				},
			},
			expectedOrders: originalOrders[3:7],
		},
		{
			name: "MakerAssetAmount >= 3 AND MakerAssetData < h AND TakerAssetAmount != 5",
			filters: []OrderFilterOpts{
				{
					Field: MakerAssetAmount,
					Kind:  GreaterOrEqual,
					Value: 3,
				},
				{
					Field: MakerAssetData,
					Kind:  Less,
					Value: []byte("h"),
				},
				{
					Field: TakerAssetAmount,
					Kind:  NotEqual,
					Value: 5,
				},
			},
			expectedOrders: append(safeSubslice(originalOrders, 3, 5), safeSubslice(originalOrders, 6, 7)...),
		},
	}
	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("%s (test case %d)", testCase.name, i)
		t.Run(testCaseName, runFindOrdersFilterTestCase(t, db, testCase))
	}
}

type findOrdersFilterTestCase struct {
	name           string
	filters        []OrderFilterOpts
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
		assertMiniHeaderSlicesAreEqual(t, miniHeaders, added)
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

	foundMiniHeaders, err := db.FindMiniHeaders()
	require.NoError(t, err)
	assertMiniHeaderSlicesAreEqual(t, originalMiniHeaders, foundMiniHeaders)
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

// returns a (shallow) subslice of orders without modifying the original slice. Uses the
// same semantics as slice expressions: low is inclusive, hi is exclusive. The returned
// slice still contains pointers, it just doesn't use the same underlying array.
func safeSubslice(orders []*Order, low, hi int) []*Order {
	result := make([]*Order, hi-low)
	for i := low; i < hi; i++ {
		result[i-low] = orders[i]
	}
	return result
}

func sortOrders(orders []*Order) {
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
	sortOrders(expectedCopy)
	actualCopy := make([]*Order, len(actual))
	copy(actualCopy, actual)
	sortOrders(actualCopy)
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

func sortMiniHeaders(miniHeaders []*MiniHeader) {
	sort.SliceStable(miniHeaders, func(i, j int) bool {
		return bytes.Compare(miniHeaders[i].Hash.Bytes(), miniHeaders[j].Hash.Bytes()) == -1
	})
}

func assertMiniHeaderSlicesAreEqual(t *testing.T, expected, actual []*MiniHeader) {
	assert.Len(t, actual, len(expected), "wrong number of miniheaders")
	sortMiniHeaders(expected)
	sortMiniHeaders(actual)
	for i, expectedMiniHeader := range expected {
		if i >= len(actual) {
			break
		}
		actualMiniHeader := expected[i]
		assertMiniHeadersAreEqual(t, expectedMiniHeader, actualMiniHeader)
	}
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

// func TestFindOrdersByMakerAddressMakerFeeAssetAddressTokenID(t *testing.T) {
// 	meshDB, err := New("/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
// 	require.NoError(t, err)
// 	defer meshDB.Close()

// 	makerAddress := constants.GanacheAccount0
// 	nextSalt := big.NewInt(1548619145450)

// 	zeroexOrders := []*zeroex.Order{
// 		// No Maker fee
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     constants.NullBytes,
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 		// ERC20 maker fee
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32"),
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 		// ERC721 maker fee with token id = 1
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 		// ERC721 maker fee with token id = 2
// 		&zeroex.Order{
// 			ChainID:               big.NewInt(constants.TestChainID),
// 			ExchangeAddress:       contractAddresses.Exchange,
// 			MakerAddress:          makerAddress,
// 			TakerAddress:          constants.NullAddress,
// 			SenderAddress:         constants.NullAddress,
// 			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
// 			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
// 			TakerFeeAssetData:     constants.NullBytes,
// 			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
// 			MakerFeeAssetData:     common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000002"),
// 			Salt:                  nextSalt.Add(nextSalt, big.NewInt(1)),
// 			MakerFee:              big.NewInt(0),
// 			TakerFee:              big.NewInt(0),
// 			MakerAssetAmount:      big.NewInt(3551808554499581700),
// 			TakerAssetAmount:      big.NewInt(1),
// 			ExpirationTimeSeconds: big.NewInt(1548619325),
// 		},
// 	}
// 	orders := make([]*Order, len(zeroexOrders))
// 	for i, o := range zeroexOrders {
// 		signedOrder, err := zeroex.SignTestOrder(o)
// 		require.NoError(t, err)
// 		orderHash, err := o.ComputeOrderHash()
// 		require.NoError(t, err)

// 		orders[i] = &Order{
// 			Hash:                     orderHash,
// 			SignedOrder:              signedOrder,
// 			FillableTakerAssetAmount: big.NewInt(1),
// 			LastUpdated:              time.Now().UTC(),
// 			IsRemoved:                false,
// 		}
// 		require.NoError(t, meshDB.Orders.Insert(orders[i]))
// 		// We need to call ResetHash so that unexported hash field is equal in later
// 		// assertions.
// 		signedOrder.ResetHash()
// 	}

// 	testCases := []struct {
// 		makerFeeAssetAddress common.Address
// 		makerFeeTokenID      *big.Int
// 		expectedOrders       []*Order
// 	}{
// 		{
// 			makerFeeAssetAddress: constants.NullAddress,
// 			makerFeeTokenID:      nil,
// 			expectedOrders:       orders[0:1],
// 		},
// 		{
// 			makerFeeAssetAddress: common.HexToAddress("0x38ae374ecf4db50b0ff37125b591a04997106a32"),
// 			makerFeeTokenID:      nil,
// 			expectedOrders:       orders[1:2],
// 		},
// 		{
// 			// Since no token id was specified, this query should match all token ids.
// 			makerFeeAssetAddress: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			makerFeeTokenID:      nil,
// 			expectedOrders:       orders[2:4],
// 		},
// 		{
// 			makerFeeAssetAddress: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			makerFeeTokenID:      big.NewInt(1),
// 			expectedOrders:       orders[2:3],
// 		},
// 		{
// 			makerFeeAssetAddress: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
// 			makerFeeTokenID:      big.NewInt(2),
// 			expectedOrders:       orders[3:4],
// 		},
// 	}
// 	for i, tc := range testCases {
// 		foundOrders, err := meshDB.FindOrdersByMakerAddressMakerFeeAssetAddressAndTokenID(makerAddress, tc.makerFeeAssetAddress, tc.makerFeeTokenID)
// 		require.NoError(t, err)
// 		assert.Equal(t, tc.expectedOrders, foundOrders, "test case %d", i)
// 	}
// }

// func insertRawOrders(t *testing.T, meshDB *MeshDB, rawOrders []*zeroex.Order, isPinned bool) []*Order {
// 	results := make([]*Order, len(rawOrders))
// 	for i, order := range rawOrders {
// 		// Sign, compute order hash, and insert.
// 		signedOrder, err := zeroex.SignTestOrder(order)
// 		require.NoError(t, err)
// 		orderHash, err := order.ComputeOrderHash()
// 		require.NoError(t, err)

// 		order := &Order{
// 			Hash:                     orderHash,
// 			SignedOrder:              signedOrder,
// 			FillableTakerAssetAmount: big.NewInt(1),
// 			LastUpdated:              time.Now(),
// 			IsRemoved:                false,
// 			IsPinned:                 isPinned,
// 		}
// 		results[i] = order
// 		require.NoError(t, meshDB.Orders.Insert(order))
// 	}
// 	return results
// }

// func TestPruneMiniHeadersAboveRetentionLimit(t *testing.T) {
// 	t.Parallel()

// 	meshDB, err := New("/tmp/meshdb_testing/"+uuid.New().String(), contractAddresses)
// 	require.NoError(t, err)
// 	defer meshDB.Close()

// 	txn := meshDB.MiniHeaders.OpenTransaction()
// 	defer func() {
// 		_ = txn.Discard()
// 	}()

// 	miniHeadersToAdd := miniHeadersMaxPerPage*2 + defaultMiniHeaderRetentionLimit + 1
// 	for i := 0; i < miniHeadersToAdd; i++ {
// 		miniHeader := &miniheader.MiniHeader{
// 			Hash:      common.BigToHash(big.NewInt(int64(i))),
// 			Number:    big.NewInt(int64(i)),
// 			Timestamp: time.Now().Add(time.Duration(i)*time.Second - 5*time.Hour),
// 		}
// 		require.NoError(t, txn.Insert(miniHeader))
// 	}
// 	require.NoError(t, txn.Commit())

// 	require.NoError(t, meshDB.PruneMiniHeadersAboveRetentionLimit())
// 	remainingMiniHeaders, err := meshDB.MiniHeaders.Count()
// 	assert.Equal(t, defaultMiniHeaderRetentionLimit, remainingMiniHeaders, "wrong number of MiniHeaders remaining")
// }
