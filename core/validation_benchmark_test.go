// +build !js

// TODO(albrow): Some tests don't require any network calls and should be able
// to run in a Wasm/JavaScript environment.

package core

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
)

func generateTestOrders(makerAddresses []common.Address, count int) []*testOrder {
	testOrders := make([]*testOrder, count)
	for i := 0; i < count; i++ {
		makerAddress := makerAddresses[i%len(makerAddresses)]
		testOrders[i] = newTestOrder().
			withMakerAddress(makerAddress).
			withSalt(big.NewInt(int64(i)))
	}
	return testOrders
}

func generateETHBackings(makerAddresses []common.Address) []*meshdb.ETHBacking {
	ethBackings := make([]*meshdb.ETHBacking, len(makerAddresses))
	for i, makerAdress := range makerAddresses {
		ethBackings[i] = &meshdb.ETHBacking{
			MakerAddress: makerAdress,
			OrderCount:   rand.Intn(100000),
			AmountInWei:  float64(rand.Intn(1000000)),
		}
	}
	return ethBackings
}

func generateMakerAddresses(count int) []common.Address {
	addresses := make([]common.Address, count)
	for i := 0; i < count; i++ {
		addressHex := fmt.Sprintf("%040x", i)
		address := common.HexToAddress(addressHex)
		addresses = append(addresses, address)
	}
	return addresses
}

func BenchmarkValidateETHBackings1Address100Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 1, 100)
}

func BenchmarkValidateETHBackings1Address1000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 1, 1000)
}

func BenchmarkValidateETHBackings10Addresses100Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 10, 100)
}

func BenchmarkValidateETHBackings10Addresses1000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 10, 1000)
}

func BenchmarkValidateETHBackings100Addresses100Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 100, 100)
}

func BenchmarkValidateETHBackings100Addresses1000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 100, 1000)
}

func BenchmarkValidateETHBackings1000Addresses100Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 1000, 100)
}

func BenchmarkValidateETHBackings1000Addresses1000Orders(b *testing.B) {
	benchmarkValidateETHBackings(b, 1000, 1000)
}

func benchmarkValidateETHBackings(b *testing.B, addressCount int, orderCount int) {
	makerAddresses := generateMakerAddresses(addressCount)
	orders := testOrdersToFakeSignedOrders(generateTestOrders(makerAddresses, orderCount))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ethBackings := generateETHBackings(makerAddresses)
		b.StartTimer()
		validateETHBackingsWithHeap(0, ethBackings, orders)
	}
}

// converts a testOrder to a *zeroex.SignedOrder with an empty signature. This
// won't pass signature validation but is fine for benchmarking.
func testOrdersToFakeSignedOrders(testOrders []*testOrder) []*zeroex.SignedOrder {
	signedOrders := make([]*zeroex.SignedOrder, len(testOrders))
	for i, testOrder := range testOrders {
		signedOrders[i] = &zeroex.SignedOrder{
			Order:     (zeroex.Order)(*testOrder),
			Signature: []byte{},
		}
	}
	return signedOrders
}
