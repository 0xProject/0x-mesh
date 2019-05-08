package meshdb

import (
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
var ganacheExchangeAddress = common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788")

func TestOrderCRUDOperations(t *testing.T) {
	meshDB, err := NewMeshDB("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)

	makerAddress := common.HexToAddress("0x6924a03bb710eaf199ab6ac9f2bb148215ae9b5d")
	salt := big.NewInt(1548619145450)
	signedOrder := &zeroex.SignedOrder{
		MakerAddress:          makerAddress,
		TakerAddress:          nullAddress,
		SenderAddress:         nullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
		MakerAssetData:        common.Hex2Bytes("025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001"),
		Salt:                  salt,
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(3551808554499581700),
		TakerAssetAmount:      big.NewInt(1),
		ExpirationTimeSeconds: big.NewInt(1548619325),
		ExchangeAddress:       ganacheExchangeAddress,
	}
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	// Insert
	order := &Order{
		Hash:                     orderHash,
		SignedOrder:              signedOrder,
		FillableTakerAssetAmount: big.NewInt(1),
	}
	require.NoError(t, meshDB.Orders.Insert(order))

	// Find
	foundOrder := &Order{}
	require.NoError(t, meshDB.Orders.FindByID(order.ID(), foundOrder))
	assert.Equal(t, order, foundOrder)

	// Check Indexes
	filter := meshDB.Orders.SaltIndex.ValueFilter(salt.Bytes())
	orders := []*Order{}
	require.NoError(t, meshDB.Orders.NewQuery(filter).Run(&orders))
	assert.Equal(t, []*Order{order}, orders)

	orders, err = meshDB.FindOrdersByMakerAddress(makerAddress)
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
	assert.Equal(t, "leveldb: not found", err.Error())
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
