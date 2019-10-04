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
		TakerFeeAssetData:     constants.NullBytes,
		MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
		MakerFeeAssetData:     constants.NullBytes,
		Salt:                  salt,
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(3551808554499581700),
		TakerAssetAmount:      big.NewInt(1),
		ExpirationTimeSeconds: big.NewInt(1548619325),
		DomainHash:            constants.NetworkIDToDomainHash[constants.TestNetworkID],
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
