// +build !js

// TODO(albrow): Some tests don't require any network calls and should be able
// to run in a Wasm/JavaScript environment.

package core

import (
	"container/heap"
	"math/big"
	"math/rand"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func generateTestOrders(t assert.TestingT, makerAddresses []common.Address, count int) []*testOrder {
	testOrders := make([]*testOrder, count)
	for i := 0; i < count; i++ {
		makerAddress := makerAddresses[i%len(makerAddresses)]
		testOrders[i] = newTestOrder().
			withMakerAddress(makerAddress).
			withSalt(big.NewInt(int64(i)))
	}
	return testOrders
}

func generateETHBackingHeap(t assert.TestingT, makerAddresses []common.Address) *ETHBackingHeap {
	ethBackings := make([]*meshdb.ETHBacking, len(makerAddresses))
	for i, makerAdress := range makerAddresses {
		ethBackings[i] = &meshdb.ETHBacking{
			MakerAddress: makerAdress,
			OrderCount:   rand.Intn(100),
			ETHAmount:    big.NewInt(rand.Int63n(1000)),
		}
	}
	ethBackingHeap := ETHBackingHeap(ethBackings)
	heap.Init(&ethBackingHeap)
	return &ethBackingHeap
}

var testAccounts = []common.Address{constants.GanacheAccount0, constants.GanacheAccount1, constants.GanacheAccount2, constants.GanacheAccount3, constants.GanacheAccount4}

func BenchmarkValidateETHBackings1Account100Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, testAccounts[0:1], 100)
}

func BenchmarkValidateETHBackings1Account1000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, testAccounts[0:1], 1000)
}
func BenchmarkValidateETHBackings1Account10000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, testAccounts[0:1], 10000)
}

func BenchmarkValidateETHBackings5Accounts100Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, testAccounts, 100)
}

func BenchmarkValidateETHBackings5Accounts1000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, testAccounts, 1000)
}
func BenchmarkValidateETHBackings5Accounts10000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, testAccounts, 10000)
}

func benchmarkValidateETHBackings(b *testing.B, makerAddresses []common.Address, count int) {
	orders := testOrdersToSignedOrders(b, generateTestOrders(b, makerAddresses, count))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ethBackingHeap := generateETHBackingHeap(b, makerAddresses)
		b.StartTimer()
		validateETHBackingsWithHeap(0, ethBackingHeap, orders)
	}
}
