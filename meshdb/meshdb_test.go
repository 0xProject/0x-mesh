package meshdb

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The max number of orders to store for each MeshDB instance throughout these
// tests.
const testingMaxOrders = 100

func TestOrderCRUDOperations(t *testing.T) {
	meshDB, err := NewMeshDB("/tmp/meshdb_testing/"+uuid.New().String(), testingMaxOrders)
	require.NoError(t, err)

	contractAddresses, err := ethereum.GetContractAddressesForNetworkID(constants.TestNetworkID)
	require.NoError(t, err)

	makerAddress := constants.GanacheAccount0
	salt := big.NewInt(1548619145450)
	o := &zeroex.Order{
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
		MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
		Salt:                  salt,
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(3551808554499581700),
		TakerAssetAmount:      big.NewInt(1),
		ExpirationTimeSeconds: big.NewInt(1548619325),
		ExchangeAddress:       contractAddresses.Exchange,
	}
	signedOrder, err := zeroex.SignTestOrder(o)
	require.NoError(t, err)

	orderHash, err := o.ComputeOrderHash()
	require.NoError(t, err)

	currentTime := time.Now().UTC()
	fiveMinutesFromNow := currentTime.Add(5 * time.Minute)

	// Insert
	order := &Order{
		Hash:                     orderHash,
		SignedOrder:              signedOrder,
		FillableTakerAssetAmount: big.NewInt(1),
		LastUpdated:              currentTime,
		IsRemoved:                false,
	}
	require.NoError(t, meshDB.Orders.Insert(order))

	// Find
	foundOrder := &Order{}
	require.NoError(t, meshDB.Orders.FindByID(order.ID(), foundOrder))
	// HACK(albrow): We need to call ComputeOrderHash in order to populate the
	// unexported hash field.
	_, _ = foundOrder.SignedOrder.ComputeOrderHash()
	assert.Equal(t, order, foundOrder)

	// Check Indexes
	orders, err := meshDB.FindOrdersByMakerAddressAndMaxSalt(makerAddress, salt)
	require.NoError(t, err)
	for _, foundOrder := range orders {
		// HACK(albrow): We need to call ComputeOrderHash in order to populate the
		// unexported hash field.
		_, _ = foundOrder.SignedOrder.ComputeOrderHash()
	}
	assert.Equal(t, []*Order{order}, orders)

	orders, err = meshDB.FindOrdersByMakerAddress(makerAddress)
	require.NoError(t, err)
	for _, foundOrder := range orders {
		// HACK(albrow): We need to call ComputeOrderHash in order to populate the
		// unexported hash field.
		_, _ = foundOrder.SignedOrder.ComputeOrderHash()
	}
	assert.Equal(t, []*Order{order}, orders)

	orders, err = meshDB.FindOrdersLastUpdatedBefore(fiveMinutesFromNow)
	require.NoError(t, err)
	for _, foundOrder := range orders {
		// HACK(albrow): We need to call ComputeOrderHash in order to populate the
		// unexported hash field.
		_, _ = foundOrder.SignedOrder.ComputeOrderHash()
	}
	assert.Equal(t, []*Order{order}, orders)

	// Update
	modifiedOrder := foundOrder
	modifiedOrder.FillableTakerAssetAmount = big.NewInt(0)
	// HACK(albrow): We need to call ComputeOrderHash in order to populate the
	// unexported hash field.
	_, _ = modifiedOrder.SignedOrder.ComputeOrderHash()
	require.NoError(t, meshDB.Orders.Update(modifiedOrder))
	foundModifiedOrder := &Order{}
	require.NoError(t, meshDB.Orders.FindByID(modifiedOrder.ID(), foundModifiedOrder))
	// HACK(albrow): We need to call ComputeOrderHash in order to populate the
	// unexported hash field.
	_, _ = foundModifiedOrder.SignedOrder.ComputeOrderHash()
	assert.Equal(t, modifiedOrder, foundModifiedOrder)

	// Delete
	require.NoError(t, meshDB.Orders.Delete(foundModifiedOrder.ID()))
	nonExistentOrder := &Order{}
	err = meshDB.Orders.FindByID(foundModifiedOrder.ID(), nonExistentOrder)
	assert.IsType(t, db.NotFoundError{}, err)
}

func TestParseContractAddressesAndTokenIdsFromAssetData(t *testing.T) {
	// ERC20 AssetData
	erc20AssetData := common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32")
	singleAssetDatas, err := parseContractAddressesAndTokenIdsFromAssetData(erc20AssetData)
	require.NoError(t, err)
	assert.Len(t, singleAssetDatas, 1)
	expectedAddress := common.HexToAddress("0x38ae374ecf4db50b0ff37125b591a04997106a32")
	assert.Equal(t, expectedAddress, singleAssetDatas[0].Address)
	var expectedTokenID *big.Int
	assert.Equal(t, expectedTokenID, singleAssetDatas[0].TokenID)

	// ERC721 AssetData
	erc721AssetData := common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001")
	singleAssetDatas, err = parseContractAddressesAndTokenIdsFromAssetData(erc721AssetData)
	require.NoError(t, err)
	assert.Equal(t, 1, len(singleAssetDatas))
	expectedAddress = common.HexToAddress("0x1dC4c1cEFEF38a777b15aA20260a54E584b16C48")
	assert.Equal(t, expectedAddress, singleAssetDatas[0].Address)
	expectedTokenID = big.NewInt(1)
	assert.Equal(t, expectedTokenID, singleAssetDatas[0].TokenID)

	// Multi AssetData
	multiAssetData := common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000x94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000")
	singleAssetDatas, err = parseContractAddressesAndTokenIdsFromAssetData(multiAssetData)
	require.NoError(t, err)
	assert.Equal(t, 2, len(singleAssetDatas))
	expectedSingleAssetDatas := []singleAssetData{
		singleAssetData{
			Address: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		},
		singleAssetData{
			Address: common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
			TokenID: big.NewInt(1),
		},
	}
	for i, singleAssetData := range singleAssetDatas {
		expectedSingleAssetData := expectedSingleAssetDatas[i]
		assert.Equal(t, expectedSingleAssetData.Address, singleAssetData.Address)
		assert.Equal(t, expectedSingleAssetData.TokenID, singleAssetData.TokenID)
	}
}

type insertOrdersTestCase struct {
	maxOrders           int
	initialOrders       []*testOrder
	initialETHBackings  []*ETHBacking
	incomingOrder       *testOrder
	expectedOrders      []*testOrder
	expectedETHBackings []*ETHBacking
}

func TestInsertOrder(t *testing.T) {
	t.Parallel()

	testCases := []insertOrdersTestCase{
		{
			// No orders should be inserted because we are already at the max and the
			// incoming order has the same maker address as the order with the lowest
			// ETH backing per order.
			maxOrders: 1,
			initialOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			initialETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
			},
			incomingOrder: newTestOrder().
				withMakerAddress(constants.GanacheAccount0).
				withSalt(big.NewInt(1)),
			expectedOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			expectedETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
			},
		},
		{
			// One order should be inserted because we are not yet at the max. All ETH
			// backings should be updated.
			maxOrders: 2,
			initialOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			initialETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
			},
			incomingOrder: newTestOrder().
				withMakerAddress(constants.GanacheAccount0).
				withSalt(big.NewInt(1)),
			expectedOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(1)),
			},
			expectedETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   2,
					ETHAmount:    big.NewInt(100),
				},
			},
		},
		{
			// One order should be inserted because we are not yet at the max. All ETH
			// backings should be updated. For this case we use different maker
			// addresses for the two orders so we expect two different ETH backings.
			maxOrders: 2,
			initialOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			initialETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    big.NewInt(100),
				},
			},
			incomingOrder: newTestOrder().
				withMakerAddress(constants.GanacheAccount1).
				withSalt(big.NewInt(1)),
			expectedOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount1).
					withSalt(big.NewInt(1)),
			},
			expectedETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
			},
		},
		{
			// One order should be inserted and one order should be removed. All ETH
			// backings should be updated accordingly. Account0 has a lower ETH
			// backing than Account1 so we will replace its order with the incoming
			// order.
			maxOrders: 1,
			initialOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			initialETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(50),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    big.NewInt(100),
				},
			},
			incomingOrder: newTestOrder().
				withMakerAddress(constants.GanacheAccount1).
				withSalt(big.NewInt(1)),
			expectedOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount1).
					withSalt(big.NewInt(1)),
			},
			expectedETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   0,
					ETHAmount:    big.NewInt(50),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
			},
		},
		{
			// No orders should be inserted/removed and no ETH backings should be
			// changed. Account0 has a higher ETH backing than Account1 so the
			// incoming order should not replace its order.
			maxOrders: 1,
			initialOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			initialETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    big.NewInt(50),
				},
			},
			incomingOrder: newTestOrder().
				withMakerAddress(constants.GanacheAccount1).
				withSalt(big.NewInt(1)),
			expectedOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
			},
			expectedETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    big.NewInt(100),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    big.NewInt(50),
				},
			},
		},
		{
			// A more complicated case with more orders and more accounts. One order
			// should be inserted and one order should be removed. All ETH
			// backings should be updated accordingly.
			maxOrders: 5,
			initialOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(1)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount1).
					withSalt(big.NewInt(2)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount1).
					withSalt(big.NewInt(3)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount2).
					withSalt(big.NewInt(4)),
			},
			initialETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   2,
					ETHAmount:    big.NewInt(100),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   2,
					ETHAmount:    big.NewInt(50),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount2,
					OrderCount:   1,
					ETHAmount:    big.NewInt(75),
				},
			},
			incomingOrder: newTestOrder().
				withMakerAddress(constants.GanacheAccount2).
				withSalt(big.NewInt(5)),
			expectedOrders: []*testOrder{
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(0)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount0).
					withSalt(big.NewInt(1)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount1).
					withSalt(big.NewInt(3)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount2).
					withSalt(big.NewInt(4)),
				newTestOrder().
					withMakerAddress(constants.GanacheAccount2).
					withSalt(big.NewInt(5)),
			},
			expectedETHBackings: []*ETHBacking{
				&ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   2,
					ETHAmount:    big.NewInt(100),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   1,
					ETHAmount:    big.NewInt(50),
				},
				&ETHBacking{
					MakerAddress: constants.GanacheAccount2,
					OrderCount:   2,
					ETHAmount:    big.NewInt(75),
				},
			},
		},
	}
	for caseNumber, testCase := range testCases {
		testInsertOrdersCase(t, testCase, caseNumber)
	}
}

func testInsertOrdersCase(t *testing.T, testCase insertOrdersTestCase, caseNumber int) {
	t.Helper()
	testInfo := fmt.Sprintf("(test case %d)", caseNumber)
	meshDB, err := NewMeshDB("/tmp/meshdb_testing/"+uuid.New().String(), testCase.maxOrders)
	require.NoError(t, err, testInfo)

	// Set up the initial state.
	ordersTxn := meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersTxn.Discard()
	}()
	for _, order := range testCase.initialOrders {
		dbOrder := order.toDBOrder(t)
		require.NoError(t, ordersTxn.Insert(dbOrder), testInfo)
	}
	require.NoError(t, ordersTxn.Commit(), testInfo)

	backingsTxn := meshDB.ETHBackings.OpenTransaction()
	defer func() {
		_ = backingsTxn.Discard()
	}()
	for _, backing := range testCase.initialETHBackings {
		require.NoError(t, backingsTxn.Insert(backing), testInfo)
	}
	require.NoError(t, backingsTxn.Commit(), testInfo)

	// Insert the new order.
	dbOrder := testCase.incomingOrder.toDBOrder(t)
	require.NoError(t, meshDB.InsertOrder(dbOrder), testInfo)

	// Make sure that the state after inserting the order is equal to the expected
	// state. Here we sort all the results so we can do a direct comparison.
	var actualOrders []*Order
	require.NoError(t, meshDB.Orders.FindAll(&actualOrders), testInfo)
	sort.Sort(ordersByHash(actualOrders))
	expectedDBOrders := make([]*Order, len(testCase.expectedOrders))
	for i, order := range testCase.expectedOrders {
		expectedDBOrders[i] = order.toDBOrder(t)
	}
	sort.Sort(ordersByHash(expectedDBOrders))
	for _, foundOrder := range actualOrders {
		// HACK(albrow): We need to call ComputeOrderHash in order to populate the
		// unexported hash field.
		_, _ = foundOrder.SignedOrder.ComputeOrderHash()
	}
	assert.Equal(t, expectedDBOrders, actualOrders, testInfo)

	var actualBackings []*ETHBacking
	require.NoError(t, meshDB.ETHBackings.FindAll(&actualBackings), testInfo)
	sort.Sort(ethBackingsByMakerAddress(actualBackings))
	sort.Sort(ethBackingsByMakerAddress(testCase.expectedETHBackings))
	assert.Equal(t, testCase.expectedETHBackings, actualBackings, testInfo)
}

type ordersByHash []*Order

// Len is the number of elements in the collection.
func (orders ordersByHash) Len() int {
	return len(orders)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (orders ordersByHash) Less(i, j int) bool {
	return bytes.Compare(orders[i].Hash.Bytes(), orders[j].Hash.Bytes()) == -1
}

// Swap swaps the elements with indexes i and j.
func (orders ordersByHash) Swap(i, j int) {
	orders[i], orders[j] = orders[j], orders[i]
}

type ethBackingsByMakerAddress []*ETHBacking

// Len is the number of elements in the collection.
func (backings ethBackingsByMakerAddress) Len() int {
	return len(backings)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (backings ethBackingsByMakerAddress) Less(i, j int) bool {
	return bytes.Compare(backings[i].MakerAddress.Bytes(), backings[j].MakerAddress.Bytes()) == -1
}

// Swap swaps the elements with indexes i and j.
func (backings ethBackingsByMakerAddress) Swap(i, j int) {
	backings[i], backings[j] = backings[j], backings[i]
}

type testOrder zeroex.Order

func newTestOrder() *testOrder {
	return (*testOrder)(&zeroex.Order{
		MakerAddress:          constants.NullAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
		MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
		Salt:                  big.NewInt(0),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(3551808554499581700),
		TakerAssetAmount:      big.NewInt(1),
		ExpirationTimeSeconds: big.NewInt(1548619325),
		ExchangeAddress:       constants.NullAddress,
	})
}

func (order *testOrder) withMakerAddress(address common.Address) *testOrder {
	order.MakerAddress = address
	return order
}

func (order *testOrder) withSalt(salt *big.Int) *testOrder {
	order.Salt = salt
	return order
}

func (order *testOrder) toDBOrder(t *testing.T) *Order {
	t.Helper()
	signedOrder, err := zeroex.SignTestOrder((*zeroex.Order)(order))
	require.NoError(t, err)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	return &Order{
		Hash:                     orderHash,
		SignedOrder:              signedOrder,
		FillableTakerAssetAmount: big.NewInt(1),
		LastUpdated:              time.Date(1992, time.September, 29, 8, 0, 0, 0, time.UTC),
		IsRemoved:                false,
	}
}
