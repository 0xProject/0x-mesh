package meshdb

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

// MiniHeader is the database representation of a succinct Ethereum block headers
type MiniHeader struct {
	Hash   common.Hash
	Parent common.Hash
	Number *big.Int
	Logs   []types.Log
}

// ID returns the MiniHeader's ID
func (m *MiniHeader) ID() []byte {
	return m.Hash.Bytes()
}

// Order is the database representation a 0x order along with some relevant metadata
type Order struct {
	Hash        common.Hash
	SignedOrder *zeroex.SignedOrder
	// When was this order last validated
	LastUpdated time.Time
	// How much of this order can still be filled
	FillableTakerAssetAmount *big.Int
}

// ID returns the Order's ID
func (o Order) ID() []byte {
	return o.Hash.Bytes()
}

// MeshDB instantiates the DB connection and creates all the collections used by the application
type MeshDB struct {
	database    *db.DB
	MiniHeaders *MiniHeadersCollection
	Orders      *OrdersCollection
}

// MiniHeadersCollection represents a DB collection of mini Ethereum block headers
type MiniHeadersCollection struct {
	*db.Collection
	numberIndex *db.Index
}

// OrdersCollection represents a DB collection of 0x orders
type OrdersCollection struct {
	*db.Collection
	SaltIndex                            *db.Index
	MakerAddressTokenAddressTokenIDIndex *db.Index
}

// NewMeshDB instantiates a new MeshDB instance
func NewMeshDB(path string) (*MeshDB, error) {
	database, err := db.Open(path)
	if err != nil {
		return nil, err
	}

	miniHeaders, err := setupMiniHeaders(database)
	if err != nil {
		return nil, err
	}

	orders := setupOrders(database)

	return &MeshDB{
		database:    database,
		MiniHeaders: miniHeaders,
		Orders:      orders,
	}, nil
}

func setupOrders(database *db.DB) *OrdersCollection {
	col := database.NewCollection("order", &Order{})
	saltIndex := col.AddIndex("salt", func(m db.Model) []byte {
		return m.(*Order).SignedOrder.Salt.Bytes()
	})
	makerAddressTokenAddressTokenIDIndex := col.AddMultiIndex("makerAddressTokenAddressTokenId", func(m db.Model) [][]byte {
		order := m.(*Order)
		singleAssetDatas, err := parseContractAddressesAndTokenIdsFromAssetData(order.SignedOrder.MakerAssetData)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Panic("Parsing assetData failed")
		}

		indexValues := make([][]byte, len(singleAssetDatas))
		for i, singleAssetData := range singleAssetDatas {
			indexValue := []byte(order.SignedOrder.MakerAddress.Hex() + "|" + singleAssetData.Address.Hex() + "|")
			if singleAssetData.TokenID != nil {
				indexValue = append(indexValue, singleAssetData.TokenID.Bytes()...)
			}
			indexValues[i] = indexValue
		}
		return indexValues
	})

	return &OrdersCollection{
		Collection:                           col,
		MakerAddressTokenAddressTokenIDIndex: makerAddressTokenAddressTokenIDIndex,
		SaltIndex:                            saltIndex,
	}
}

func setupMiniHeaders(database *db.DB) (*MiniHeadersCollection, error) {
	col := database.NewCollection("miniHeader", &MiniHeader{})
	numberIndex := col.AddIndex("number", func(model db.Model) []byte {
		// By default, the index is sorted in byte order. In order to sort by
		// numerical order, we need to pad with zeroes. The maximum length of an
		// unsigned 256 bit integer is 80, so we pad with zeroes such that the
		// length of the number is always 80.
		number := model.(*MiniHeader).Number
		return []byte(fmt.Sprintf("%80s", number.String()))
	})

	return &MiniHeadersCollection{
		Collection:  col,
		numberIndex: numberIndex,
	}, nil
}

// Close closes the database connection
func (m *MeshDB) Close() {
	m.database.Close()
}

// FindAllMiniHeadersSortedByNumber returns all MiniHeaders sorted by block number
func (m *MeshDB) FindAllMiniHeadersSortedByNumber() ([]*MiniHeader, error) {
	miniHeaders := []*MiniHeader{}
	query := m.MiniHeaders.NewQuery(m.MiniHeaders.numberIndex.All())
	err := query.Run(&miniHeaders)
	if err != nil {
		return nil, err
	}
	return miniHeaders, nil
}

// FindLatestMiniHeader returns the latest MiniHeader (i.e. the one with the
// largest block number), or nil if there are none in the database.
func (m *MeshDB) FindLatestMiniHeader() (*MiniHeader, error) {
	miniHeaders := []*MiniHeader{}
	query := m.MiniHeaders.NewQuery(m.MiniHeaders.numberIndex.All()).Reverse().Max(1)
	err := query.Run(&miniHeaders)
	if err != nil {
		return nil, err
	}
	if len(miniHeaders) == 0 {
		return nil, nil
	}
	return miniHeaders[0], nil
}

type singleAssetData struct {
	Address common.Address
	TokenID *big.Int
}

func parseContractAddressesAndTokenIdsFromAssetData(assetData []byte) ([]singleAssetData, error) {
	singleAssetDatas := []singleAssetData{}
	assetDataDecoder, err := zeroex.NewAssetDataDecoder()
	if err != nil {
		return nil, err
	}

	assetDataName, err := assetDataDecoder.GetName(assetData)
	if err != nil {
		return nil, err
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := singleAssetData{
			Address: decodedAssetData.Address,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := singleAssetData{
			Address: decodedAssetData.Address,
			TokenID: decodedAssetData.TokenId,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, assetData := range decodedAssetData.NestedAssetData {
			as, err := parseContractAddressesAndTokenIdsFromAssetData(assetData)
			if err != nil {
				return nil, err
			}
			singleAssetDatas = append(singleAssetDatas, as...)
		}
	default:
		return nil, fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return singleAssetDatas, nil
}
