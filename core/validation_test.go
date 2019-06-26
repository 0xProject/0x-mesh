// +build !js

// TODO(albrow): Some tests don't require any network calls and should be able
// to run in a Wasm/JavaScript environment.
package core

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validateETHBackingTestCase struct {
	spareCapacity    int
	incomingOrders   []*testOrder
	ethBackings      []*meshdb.ETHBacking
	expectedValid    []*testOrder
	expectedRejected []*zeroex.RejectedOrderInfo
}

var testOrders = map[common.Address][]*testOrder{
	constants.GanacheAccount0: []*testOrder{
		newTestOrder().
			withMakerAddress(constants.GanacheAccount0).
			withSalt(big.NewInt(0)),
		newTestOrder().
			withMakerAddress(constants.GanacheAccount0).
			withSalt(big.NewInt(1)),
	},
	constants.GanacheAccount1: []*testOrder{
		newTestOrder().
			withMakerAddress(constants.GanacheAccount1).
			withSalt(big.NewInt(0)),
		newTestOrder().
			withMakerAddress(constants.GanacheAccount1).
			withSalt(big.NewInt(1)),
	},
	constants.GanacheAccount2: []*testOrder{
		newTestOrder().
			withMakerAddress(constants.GanacheAccount2).
			withSalt(big.NewInt(0)),
		newTestOrder().
			withMakerAddress(constants.GanacheAccount2).
			withSalt(big.NewInt(1)),
	},
}

func TestValidateETHBackings(t *testing.T) {
	testCases := []validateETHBackingTestCase{
		{
			// No orders should be considered valid because we have no spare capacity
			// and the incoming order does not have a greater ETH backing than the
			// current lowest ETH backing (they come from the same maker).
			spareCapacity: 0,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    100,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
			},
			expectedValid: []*testOrder{},
			expectedRejected: []*zeroex.RejectedOrderInfo{
				{
					OrderHash:   testOrders[constants.GanacheAccount0][0].hash(t),
					SignedOrder: testOrders[constants.GanacheAccount0][0].toSignedOrder(t),
					Kind:        MeshValidation,
					Status:      ROInsufficientETHBacking,
				},
			},
		},
		{
			// Same as above but with two different maker addresses. No orders should
			// be considered valid because we have no spare capacity and the incoming
			// order does not have a greater ETH backing than the current lowest ETH
			// backing.
			spareCapacity: 0,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    100,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    50,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount1][0],
			},
			expectedValid: []*testOrder{},
			expectedRejected: []*zeroex.RejectedOrderInfo{
				{
					OrderHash:   testOrders[constants.GanacheAccount1][0].hash(t),
					SignedOrder: testOrders[constants.GanacheAccount1][0].toSignedOrder(t),
					Kind:        MeshValidation,
					Status:      ROInsufficientETHBacking,
				},
			},
		},
		{
			// The incoming order should replace the existing order because we have no
			// spare capacity and it has a greater ETH backing.
			spareCapacity: 0,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    100,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    150,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount1][0],
			},
			expectedValid: []*testOrder{
				testOrders[constants.GanacheAccount1][0],
			},
			expectedRejected: []*zeroex.RejectedOrderInfo{},
		},
		{
			// One order should be valid because we have one spare capacity.
			spareCapacity: 1,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    100,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
			},
			expectedValid: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
			},
			expectedRejected: []*zeroex.RejectedOrderInfo{},
		},
		{
			// Similar to above, but with two different maker addresses. One order
			// should be valid because we have one spare capacity.
			spareCapacity: 1,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   1,
					ETHAmount:    100,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   0,
					ETHAmount:    100,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount1][0],
			},
			expectedValid: []*testOrder{
				testOrders[constants.GanacheAccount1][0],
			},
			expectedRejected: []*zeroex.RejectedOrderInfo{},
		},
		{
			// A more complicated case with more orders and more accounts. One order
			// should be valid and one order should be rejected. We still have a spare
			// capacity of 0.
			spareCapacity: 0,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   2,
					ETHAmount:    100,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   2,
					ETHAmount:    50,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount2,
					OrderCount:   1,
					ETHAmount:    75,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount1][0],
				testOrders[constants.GanacheAccount2][0],
			},
			expectedValid: []*testOrder{
				testOrders[constants.GanacheAccount2][0],
			},
			expectedRejected: []*zeroex.RejectedOrderInfo{
				{
					OrderHash:   testOrders[constants.GanacheAccount1][0].hash(t),
					SignedOrder: testOrders[constants.GanacheAccount1][0].toSignedOrder(t),
					Kind:        MeshValidation,
					Status:      ROInsufficientETHBacking,
				},
			},
		},
		{
			// A more complicated case with more orders and more accounts. This time
			// we do have a spare capacity of 1. Two orders should be valid and one
			// order should be rejected.
			spareCapacity: 1,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   2,
					ETHAmount:    100,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   2,
					ETHAmount:    50,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount2,
					OrderCount:   1,
					ETHAmount:    75,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
				testOrders[constants.GanacheAccount1][0],
				testOrders[constants.GanacheAccount2][0],
			},
			expectedValid: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
				testOrders[constants.GanacheAccount2][0],
			},
			expectedRejected: []*zeroex.RejectedOrderInfo{
				{
					OrderHash:   testOrders[constants.GanacheAccount1][0].hash(t),
					SignedOrder: testOrders[constants.GanacheAccount1][0].toSignedOrder(t),
					Kind:        MeshValidation,
					Status:      ROInsufficientETHBacking,
				},
			},
		},
		{
			// The most complicated test case of all.
			//
			// spare capacity = 2
			// total capacity = 7
			// # incoming orders = 5
			//
			// solution:
			//
			//     account0 100 ETH / 3 orders = 33.3
			//     account1  40 ETH / 1 order = 40
			//     account2  80 ETH / 3 orders = 26.7
			//
			// 6 incoming. 3 should be valid. 3 should be invalid.
			//
			// - 1 order from account0 is valid
			// - 0 orders from account1 are valid
			// - 2 orders from account2 is valid
			//
			spareCapacity: 2,
			ethBackings: []*meshdb.ETHBacking{
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount0,
					OrderCount:   2,
					ETHAmount:    100,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount1,
					OrderCount:   2,
					ETHAmount:    40,
				},
				&meshdb.ETHBacking{
					MakerAddress: constants.GanacheAccount2,
					OrderCount:   1,
					ETHAmount:    80,
				},
			},
			incomingOrders: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
				testOrders[constants.GanacheAccount1][0],
				testOrders[constants.GanacheAccount1][1],
				testOrders[constants.GanacheAccount2][0],
				testOrders[constants.GanacheAccount2][1],
			},
			expectedValid: []*testOrder{
				testOrders[constants.GanacheAccount0][0],
				testOrders[constants.GanacheAccount2][0],
				testOrders[constants.GanacheAccount2][1],
			},
			expectedRejected: []*zeroex.RejectedOrderInfo{
				{
					OrderHash:   testOrders[constants.GanacheAccount1][0].hash(t),
					SignedOrder: testOrders[constants.GanacheAccount1][0].toSignedOrder(t),
					Kind:        MeshValidation,
					Status:      ROInsufficientETHBacking,
				},
				{
					OrderHash:   testOrders[constants.GanacheAccount1][1].hash(t),
					SignedOrder: testOrders[constants.GanacheAccount1][1].toSignedOrder(t),
					Kind:        MeshValidation,
					Status:      ROInsufficientETHBacking,
				},
			},
		},
	}

	for i, testCase := range testCases {
		testValidateETHBackingCase(t, testCase, i)
	}
}

func testValidateETHBackingCase(t *testing.T, testCase validateETHBackingTestCase, caseNumber int) {
	testInfo := fmt.Sprintf("(test case %d) Addresses: [%s, %s, %s]", caseNumber, constants.GanacheAccount0.Hex(), constants.GanacheAccount1.Hex(), constants.GanacheAccount2.Hex())

	// Build the arguments we need to pass to the validation function.
	ethBackingHeap := ETHBackingHeap(testCase.ethBackings)
	heap.Init(&ethBackingHeap)
	incomingOrders := testOrdersToSignedOrders(t, testCase.incomingOrders)

	// Call the validateETHBackingsWithHeap function. and check the results
	actualValid, actualRejected := validateETHBackingsWithHeap(testCase.spareCapacity, &ethBackingHeap, incomingOrders)

	sort.Sort(ordersByHash(actualValid))
	expectedValid := testOrdersToSignedOrders(t, testCase.expectedValid)
	sort.Sort(ordersByHash(expectedValid))
	assert.Equal(t, expectedValid, actualValid, testInfo)

	sort.Sort(rejectedOrderInfosByHash(actualRejected))
	sort.Sort(rejectedOrderInfosByHash(testCase.expectedRejected))
	assert.Equal(t, testCase.expectedRejected, actualRejected, testInfo)
}

type ordersByHash []*zeroex.SignedOrder

// Len is the number of elements in the collection.
func (orders ordersByHash) Len() int {
	return len(orders)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (orders ordersByHash) Less(i, j int) bool {
	orderHashI, _ := orders[i].ComputeOrderHash()
	orderHashJ, _ := orders[j].ComputeOrderHash()
	return bytes.Compare(orderHashI.Bytes(), orderHashJ.Bytes()) == -1
}

// Swap swaps the elements with indexes i and j.
func (orders ordersByHash) Swap(i, j int) {
	orders[i], orders[j] = orders[j], orders[i]
}

type rejectedOrderInfosByHash []*zeroex.RejectedOrderInfo

// Len is the number of elements in the collection.
func (orderInfos rejectedOrderInfosByHash) Len() int {
	return len(orderInfos)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (orderInfos rejectedOrderInfosByHash) Less(i, j int) bool {
	return bytes.Compare(orderInfos[i].OrderHash.Bytes(), orderInfos[j].OrderHash.Bytes()) == -1
}

// Swap swaps the elements with indexes i and j.
func (orderInfos rejectedOrderInfosByHash) Swap(i, j int) {
	orderInfos[i], orderInfos[j] = orderInfos[j], orderInfos[i]
}

type ethBackingsByMakerAddress []*meshdb.ETHBacking

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
	(*zeroex.Order)(order).ResetOrderHash()
	return order
}

func (order *testOrder) withSalt(salt *big.Int) *testOrder {
	order.Salt = salt
	(*zeroex.Order)(order).ResetOrderHash()
	return order
}

func (order *testOrder) hash(t require.TestingT) common.Hash {
	orderHash, err := order.toSignedOrder(t).ComputeOrderHash()
	require.NoError(t, err)
	return orderHash
}

func (order *testOrder) toSignedOrder(t require.TestingT) *zeroex.SignedOrder {
	signedOrder, err := zeroex.SignTestOrder((*zeroex.Order)(order))
	require.NoError(t, err)
	return signedOrder
}

func testOrdersToSignedOrders(t require.TestingT, testOrders []*testOrder) []*zeroex.SignedOrder {
	signedOrders := make([]*zeroex.SignedOrder, len(testOrders))
	for i, testOrder := range testOrders {
		signedOrders[i] = testOrder.toSignedOrder(t)
	}
	return signedOrders
}

func (order *testOrder) toDBOrder(t require.TestingT) *meshdb.Order {
	signedOrder := order.toSignedOrder(t)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	return &meshdb.Order{
		Hash:                     orderHash,
		SignedOrder:              signedOrder,
		FillableTakerAssetAmount: big.NewInt(1),
		LastUpdated:              time.Date(1992, time.September, 29, 8, 0, 0, 0, time.UTC),
		IsRemoved:                false,
	}
}
