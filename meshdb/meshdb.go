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
	// Was this order flagged for removal? Due to the possibility of block-reorgs, instead
	// of immediately removing an order when FillableTakerAssetAmount becomes 0, we instead
	// flag it for removal. After this order isn't updated for X time and has IsRemoved = true,
	// the order can be permanently deleted.
	IsRemoved bool
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
	MakerAddressAndSaltIndex             *db.Index
	MakerAddressTokenAddressTokenIDIndex *db.Index
	LastUpdatedIndex                     *db.Index
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
	lastUpdatedIndex := col.AddIndex("lastUpdated", func(m db.Model) []byte {
		index := []byte(m.(*Order).LastUpdated.Format(time.RFC3339))
		return index
	})
	makerAddressAndSaltIndex := col.AddIndex("makerAddressAndSalt", func(m db.Model) []byte {
		// By default, the index is sorted in byte order. In order to sort by
		// numerical order, we need to pad with zeroes. The maximum length of an
		// unsigned 256 bit integer is 80, so we pad with zeroes such that the
		// length of the number is always 80.
		signedOrder := m.(*Order).SignedOrder
		index := []byte(fmt.Sprintf("%s|%080s", signedOrder.MakerAddress.Hex(), signedOrder.Salt.String()))
		return index
	})
	// TODO(fabio): Optimize this index callback since it gets called many times under-the-hood.
	// We might want to parse the assetData once and store it's components in the DB. The trade-off
	// here is compute time for storage space.
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
		MakerAddressAndSaltIndex:             makerAddressAndSaltIndex,
		LastUpdatedIndex:                     lastUpdatedIndex,
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
		return []byte(fmt.Sprintf("%080s", number.String()))
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

// FindOrdersByMakerAddress finds all orders belonging to a particular maker address
func (m *MeshDB) FindOrdersByMakerAddress(makerAddress common.Address) ([]*Order, error) {
	prefix := []byte(makerAddress.Hex() + "|")
	filter := m.Orders.MakerAddressTokenAddressTokenIDIndex.PrefixFilter(prefix)
	orders := []*Order{}
	err := m.Orders.NewQuery(filter).Run(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersByMakerAddressTokenAddressAndTokenID finds all orders belonging to a particular maker
// address where makerAssetData encodes for a particular token contract & token ID
func (m *MeshDB) FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	prefix := []byte(makerAddress.Hex() + "|" + tokenAddress.Hex() + "|")
	if tokenID != nil {
		prefix = append(prefix, tokenID.Bytes()...)
	}
	filter := m.Orders.MakerAddressTokenAddressTokenIDIndex.PrefixFilter(prefix)
	orders := []*Order{}
	err := m.Orders.NewQuery(filter).Run(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersByMakerAddressAndMaxSalt finds all orders belonging to a particular maker address that
// also have a salt value less then or equal to X
func (m *MeshDB) FindOrdersByMakerAddressAndMaxSalt(makerAddress common.Address, salt *big.Int) ([]*Order, error) {
	// DB range queries exclude the limit value however for the 0x protocol `cancelOrdersUpTo` method
	// is inclusive of the value supplied. In order to make this helper method more useful to our
	// particular use-case, we add 1 to the supplied salt (making the query include the value supplied)
	saltPlusOne := new(big.Int).Add(salt, big.NewInt(1))
	start := []byte(fmt.Sprintf("%s|%080s", makerAddress.Hex(), "0"))
	limit := []byte(fmt.Sprintf("%s|%080s", makerAddress.Hex(), saltPlusOne.String()))
	filter := m.Orders.MakerAddressAndSaltIndex.RangeFilter(start, limit)
	orders := []*Order{}
	err := m.Orders.NewQuery(filter).Run(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersLastUpdatedBefore finds all orders where the LastUpdated time is less
// than X
func (m *MeshDB) FindOrdersLastUpdatedBefore(lastUpdated time.Time) ([]*Order, error) {
	start := []byte(time.Unix(0, 0).Format(time.RFC3339))
	limit := []byte(lastUpdated.Format(time.RFC3339))
	filter := m.Orders.LastUpdatedIndex.RangeFilter(start, limit)
	orders := []*Order{}
	err := m.Orders.NewQuery(filter).Run(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
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
