// +build js

package main

import (
	"math/big"
	"reflect"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/qunit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/norunners/vert"
)

func main() {
	qunit.Module("meshdb")
	qunit.Test("Order CRUD operations", func(assert qunit.QUnitAssert) {
		dbPath := "meshdb-testing-" + uuid.New().String()
		meshDB, err := meshdb.New(dbPath)
		assertNoError(assert, err, "")
		defer meshDB.Close()

		contractAddresses, err := ethereum.GetContractAddressesForNetworkID(constants.TestNetworkID)
		assertNoError(assert, err, "")

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
		assertNoError(assert, err, "")

		orderHash, err := o.ComputeOrderHash()
		assertNoError(assert, err, "")

		currentTime := time.Now().UTC()
		fiveMinutesFromNow := currentTime.Add(5 * time.Minute)

		// Insert
		order := &meshdb.Order{
			Hash:                     orderHash,
			SignedOrder:              signedOrder,
			FillableTakerAssetAmount: big.NewInt(1),
			LastUpdated:              currentTime,
			IsRemoved:                false,
		}
		assertNoError(assert, meshDB.Orders.Insert(order), "")

		// Find
		foundOrder := &meshdb.Order{}
		assertNoError(assert, meshDB.Orders.FindByID(order.ID(), foundOrder), "")
		assertEqual(assert, order, foundOrder, "")

		// Check Indexes
		orders, err := meshDB.FindOrdersByMakerAddressAndMaxSalt(makerAddress, salt)
		assertNoError(assert, err, "")
		assertEqual(assert, []*meshdb.Order{order}, orders, "")

		orders, err = meshDB.FindOrdersByMakerAddress(makerAddress)
		assertNoError(assert, err, "")
		assertEqual(assert, []*meshdb.Order{order}, orders, "")

		orders, err = meshDB.FindOrdersLastUpdatedBefore(fiveMinutesFromNow)
		assertNoError(assert, err, "")
		assertEqual(assert, []*meshdb.Order{order}, orders, "")

		// Update
		modifiedOrder := foundOrder
		modifiedOrder.FillableTakerAssetAmount = big.NewInt(0)
		assertNoError(assert, meshDB.Orders.Update(modifiedOrder), "")
		foundModifiedOrder := &meshdb.Order{}
		assertNoError(assert, meshDB.Orders.FindByID(modifiedOrder.ID(), foundModifiedOrder), "")
		assertEqual(assert, modifiedOrder, foundModifiedOrder, "")

		// Delete
		assertNoError(assert, meshDB.Orders.Delete(foundModifiedOrder.ID()), "")
		nonExistentOrder := &meshdb.Order{}
		err = meshDB.Orders.FindByID(foundModifiedOrder.ID(), nonExistentOrder)
		assertEqual(assert, reflect.TypeOf(db.NotFoundError{}).String(), reflect.TypeOf(err).String(), "")
	})

	qunit.Start()

	select {}
}

func assertNoError(assert qunit.QUnitAssert, err error, msg string) {
	if err != nil {
		assert.Ok(false, "unexpected error: "+err.Error())
	}
}

func assertEqual(assert qunit.QUnitAssert, expected interface{}, actual interface{}, msg string) {
	assert.DeepEqual(vert.ValueOf(actual), vert.ValueOf(expected), msg)
}
