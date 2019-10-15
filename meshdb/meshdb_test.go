package meshdb

import (
	"math/big"
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

func TestOrderCRUDOperations(t *testing.T) {
	meshDB, err := New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

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
	// We need to call ResetHash so that unexported hash field is equal in later
	// assertions.
	signedOrder.ResetHash()

	// Find
	foundOrder := &Order{}
	require.NoError(t, meshDB.Orders.FindByID(order.ID(), foundOrder))
	assert.Equal(t, order, foundOrder)

	// Check Indexes
	orders, err := meshDB.FindOrdersByMakerAddressAndMaxSalt(makerAddress, salt)
	require.NoError(t, err)
	assert.Equal(t, []*Order{order}, orders)

	orders, err = meshDB.FindOrdersByMakerAddress(makerAddress)
	require.NoError(t, err)
	assert.Equal(t, []*Order{order}, orders)

	orders, err = meshDB.FindOrdersLastUpdatedBefore(fiveMinutesFromNow)
	require.NoError(t, err)
	assert.Equal(t, []*Order{order}, orders)

	// Update
	modifiedOrder := foundOrder
	modifiedOrder.FillableTakerAssetAmount = big.NewInt(0)
	require.NoError(t, meshDB.Orders.Update(modifiedOrder))
	foundModifiedOrder := &Order{}
	require.NoError(t, meshDB.Orders.FindByID(modifiedOrder.ID(), foundModifiedOrder))
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

func TestTrimOrdersByExpirationTime(t *testing.T) {
	meshDB, err := New("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)
	defer meshDB.Close()

	// TODO(albrow): Move these to top of file.
	contractAddresses, err := ethereum.GetContractAddressesForNetworkID(constants.TestNetworkID)
	require.NoError(t, err)
	makerAddress := constants.GanacheAccount0

	// Note: most of the fields in these orders are the same. For the purposes of
	// this test, the only thing that matters is the Salt and ExpirationTime.
	rawOrders := []*zeroex.Order{
		{
			MakerAddress:          makerAddress,
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
			ExpirationTimeSeconds: big.NewInt(100),
			ExchangeAddress:       contractAddresses.Exchange,
		},
		{
			MakerAddress:          makerAddress,
			TakerAddress:          constants.NullAddress,
			SenderAddress:         constants.NullAddress,
			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
			Salt:                  big.NewInt(1),
			MakerFee:              big.NewInt(0),
			TakerFee:              big.NewInt(0),
			MakerAssetAmount:      big.NewInt(3551808554499581700),
			TakerAssetAmount:      big.NewInt(1),
			ExpirationTimeSeconds: big.NewInt(200),
			ExchangeAddress:       contractAddresses.Exchange,
		},
		{
			MakerAddress:          makerAddress,
			TakerAddress:          constants.NullAddress,
			SenderAddress:         constants.NullAddress,
			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
			Salt:                  big.NewInt(2),
			MakerFee:              big.NewInt(0),
			TakerFee:              big.NewInt(0),
			MakerAssetAmount:      big.NewInt(3551808554499581700),
			TakerAssetAmount:      big.NewInt(1),
			ExpirationTimeSeconds: big.NewInt(200),
			ExchangeAddress:       contractAddresses.Exchange,
		},
		{
			MakerAddress:          makerAddress,
			TakerAddress:          constants.NullAddress,
			SenderAddress:         constants.NullAddress,
			FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
			TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
			MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
			Salt:                  big.NewInt(3),
			MakerFee:              big.NewInt(0),
			TakerFee:              big.NewInt(0),
			MakerAssetAmount:      big.NewInt(3551808554499581700),
			TakerAssetAmount:      big.NewInt(1),
			ExpirationTimeSeconds: big.NewInt(300),
			ExchangeAddress:       contractAddresses.Exchange,
		},
	}

	orders := make([]*Order, len(rawOrders))
	for i, order := range rawOrders {
		// Sign, compute order hash, and insert.
		signedOrder, err := zeroex.SignTestOrder(order)
		require.NoError(t, err)
		orderHash, err := order.ComputeOrderHash()
		require.NoError(t, err)

		order := &Order{
			Hash:                     orderHash,
			SignedOrder:              signedOrder,
			FillableTakerAssetAmount: big.NewInt(1),
			LastUpdated:              time.Now(),
			IsRemoved:                false,
		}
		orders[i] = order
		require.NoError(t, meshDB.Orders.Insert(order))
	}

	// Call CalculateNewMaxExpirationTimeAndTrimDatabase and check the results.
	targetMaxOrders := 2
	gotExpirationTime, gotRemovedOrders, err := meshDB.TrimOrdersByExpirationTime(targetMaxOrders)
	require.NoError(t, err)
	assert.Equal(t, "199", gotExpirationTime.String(), "newMaxExpirationTime")
	assert.Len(t, gotRemovedOrders, 2, "wrong number of orders removed")
	// Check that the expiration time of each removed order is >= the new max.
	for _, removedOrder := range gotRemovedOrders {
		expirationTimeOfRemovedOrder := removedOrder.SignedOrder.ExpirationTimeSeconds
		assert.True(t, expirationTimeOfRemovedOrder.Cmp(gotExpirationTime) != -1, "an order was removed with expiration time (%s) less than the new max (%s)", expirationTimeOfRemovedOrder, gotExpirationTime)
	}
	var remainingOrders []*Order
	require.NoError(t, meshDB.Orders.FindAll(&remainingOrders))
	assert.Len(t, remainingOrders, 2, "wrong number of orders remaining")
	// Check that the expiration time of each remaining order is <= the new max.
	for _, removedOrder := range remainingOrders {
		expirationTimeOfRemovedOrder := removedOrder.SignedOrder.ExpirationTimeSeconds
		newMaxPlusOne := big.NewInt(0).Add(gotExpirationTime, big.NewInt(1))
		assert.True(t, expirationTimeOfRemovedOrder.Cmp(newMaxPlusOne) != 1, "a remaining order had an expiration time (%s) greater than the new max + 1 (%s)", expirationTimeOfRemovedOrder, newMaxPlusOne)
	}
}
