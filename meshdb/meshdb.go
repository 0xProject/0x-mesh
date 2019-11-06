package meshdb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

var ErrDBFilledWithPinnedOrders = errors.New("the database is full of pinned orders; no orders can be removed in order to make space")

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
	// IsPinned indicates whether or not the order is pinned. Pinned orders are
	// not removed from the database unless they become unfillable.
	IsPinned bool
}

// ID returns the Order's ID
func (o Order) ID() []byte {
	return o.Hash.Bytes()
}

// Metadata is the database representation of MeshDB instance metadata
type Metadata struct {
	EthereumChainID int
	MaxExpirationTime                 *big.Int
	EthRPCRequestsSentInCurrentUTCDay int
	StartOfCurrentUTCDay              time.Time
}

// ID returns the id used for the metadata collection (one per DB)
func (m Metadata) ID() []byte {
	return []byte{0}
}

// MeshDB instantiates the DB connection and creates all the collections used by the application
type MeshDB struct {
	database    *db.DB
	metadata    *MetadataCollection
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
	MakerAddressAndSaltIndex                     *db.Index
	MakerAddressTokenAddressTokenIDIndex         *db.Index
	MakerAddressMakerFeeAssetAddressTokenIDIndex *db.Index
	LastUpdatedIndex                             *db.Index
	IsRemovedIndex                               *db.Index
	ExpirationTimeIndex                  *db.Index
}

// MetadataCollection represents a DB collection used to store instance metadata
type MetadataCollection struct {
	*db.Collection
}

// New instantiates a new MeshDB instance
func New(path string) (*MeshDB, error) {
	database, err := db.Open(path)
	if err != nil {
		return nil, err
	}

	miniHeaders, err := setupMiniHeaders(database)
	if err != nil {
		return nil, err
	}

	orders, err := setupOrders(database)
	if err != nil {
		return nil, err
	}

	metadata, err := setupMetadata(database)
	if err != nil {
		return nil, err
	}

	return &MeshDB{
		database:    database,
		metadata:    metadata,
		MiniHeaders: miniHeaders,
		Orders:      orders,
	}, nil
}

func setupOrders(database *db.DB) (*OrdersCollection, error) {
	col, err := database.NewCollection("order", &Order{})
	if err != nil {
		return nil, err
	}
	lastUpdatedIndex := col.AddIndex("lastUpdated", func(m db.Model) []byte {
		index := []byte(m.(*Order).LastUpdated.UTC().Format(time.RFC3339Nano))
		return index
	})
	makerAddressAndSaltIndex := col.AddIndex("makerAddressAndSalt", func(m db.Model) []byte {
		// By default, the index is sorted in byte order. In order to sort by
		// numerical order, we need to pad with zeroes. The maximum length of an
		// unsigned 256 bit integer is 80, so we pad with zeroes such that the
		// length of the number is always 80.
		signedOrder := m.(*Order).SignedOrder
		index := []byte(fmt.Sprintf("%s|%s", signedOrder.MakerAddress.Hex(), uint256ToConstantLengthBytes(signedOrder.Salt)))
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
	makerAddressMakerFeeAssetAddressTokenIDIndex := col.AddMultiIndex("makerAddressMakerFeeAssetAddressTokenID", func(m db.Model) [][]byte {
		order := m.(*Order)

		if bytes.Equal(order.SignedOrder.MakerFeeAssetData, constants.NullBytes) {
			// MakerFeeAssetData is optional and the lack of a maker fee is indicated
			// by null bytes ("0x0"). We still want to index this value so we can look
			// up orders without a maker fee.
			return [][]byte{
				[]byte(order.SignedOrder.MakerAddress.Hex() + "|" + common.ToHex(constants.NullBytes) + "|"),
			}
		}

		singleAssetDatas, err := parseContractAddressesAndTokenIdsFromAssetData(order.SignedOrder.MakerFeeAssetData)
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

	isRemovedIndex := col.AddIndex("isRemoved", func(m db.Model) []byte {
		order := m.(*Order)
		// false = 0; true = 1
		if order.IsRemoved {
			return []byte{1}
		}
		return []byte{0}
	})

	expirationTimeIndex := col.AddIndex("expirationTime", func(m db.Model) []byte {
		order := m.(*Order)
		expTimeString := uint256ToConstantLengthBytes(order.SignedOrder.ExpirationTimeSeconds)
		// We separate pinned and non-pinned orders via a prefix that is either 0 or
		// 1.
		pinnedString := "0"
		if order.IsPinned {
			pinnedString = "1"
		}
		return []byte(fmt.Sprintf("%s|%s", pinnedString, expTimeString))
	})

	return &OrdersCollection{
		Collection:                                   col,
		MakerAddressTokenAddressTokenIDIndex:         makerAddressTokenAddressTokenIDIndex,
		MakerAddressMakerFeeAssetAddressTokenIDIndex: makerAddressMakerFeeAssetAddressTokenIDIndex,
		MakerAddressAndSaltIndex:                     makerAddressAndSaltIndex,
		LastUpdatedIndex:                             lastUpdatedIndex,
		IsRemovedIndex:                               isRemovedIndex,
		ExpirationTimeIndex:                  expirationTimeIndex,
	}, nil
}

func setupMiniHeaders(database *db.DB) (*MiniHeadersCollection, error) {
	col, err := database.NewCollection("miniHeader", &miniheader.MiniHeader{})
	if err != nil {
		return nil, err
	}
	numberIndex := col.AddIndex("number", func(model db.Model) []byte {
		// By default, the index is sorted in byte order. In order to sort by
		// numerical order, we need to pad with zeroes. The maximum length of an
		// unsigned 256 bit integer is 80, so we pad with zeroes such that the
		// length of the number is always 80.
		number := model.(*miniheader.MiniHeader).Number
		return uint256ToConstantLengthBytes(number)
	})

	return &MiniHeadersCollection{
		Collection:  col,
		numberIndex: numberIndex,
	}, nil
}

func setupMetadata(database *db.DB) (*MetadataCollection, error) {
	col, err := database.NewCollection("metadata", &Metadata{})
	if err != nil {
		return nil, err
	}
	return &MetadataCollection{col}, nil
}

// Close closes the database connection
func (m *MeshDB) Close() {
	m.database.Close()
}

// FindAllMiniHeadersSortedByNumber returns all MiniHeaders sorted by block number
func (m *MeshDB) FindAllMiniHeadersSortedByNumber() ([]*miniheader.MiniHeader, error) {
	miniHeaders := []*miniheader.MiniHeader{}
	query := m.MiniHeaders.NewQuery(m.MiniHeaders.numberIndex.All())
	if err := query.Run(&miniHeaders); err != nil {
		return nil, err
	}
	return miniHeaders, nil
}

// FindLatestMiniHeader returns the latest MiniHeader (i.e. the one with the
// largest block number), or nil if there are none in the database.
func (m *MeshDB) FindLatestMiniHeader() (*miniheader.MiniHeader, error) {
	miniHeaders := []*miniheader.MiniHeader{}
	query := m.MiniHeaders.NewQuery(m.MiniHeaders.numberIndex.All()).Reverse().Max(1)
	if err := query.Run(&miniHeaders); err != nil {
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
	if err := m.Orders.NewQuery(filter).Run(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersByMakerAddressTokenAddressAndTokenID finds all orders belonging to a particular maker
// address where makerAssetData encodes for a particular token contract and optionally a token ID
func (m *MeshDB) FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	prefix := []byte(makerAddress.Hex() + "|" + tokenAddress.Hex() + "|")
	if tokenID != nil {
		prefix = append(prefix, tokenID.Bytes()...)
	}
	filter := m.Orders.MakerAddressTokenAddressTokenIDIndex.PrefixFilter(prefix)
	orders := []*Order{}
	if err := m.Orders.NewQuery(filter).Run(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersByMakerAddressMakerFeeAssetAddressTokenID finds all orders belonging to
// a particular maker address where makerFeeAssetData encodes for a particular
// token contract and optionally a token ID. To find orders without a maker fee,
// use constants.NullAddress for makerFeeAssetAddress.
func (m *MeshDB) FindOrdersByMakerAddressMakerFeeAssetAddressAndTokenID(makerAddress, makerFeeAssetAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	var prefix []byte
	if makerFeeAssetAddress == constants.NullAddress {
		prefix = []byte(makerAddress.Hex() + "|" + common.ToHex(constants.NullBytes) + "|")
	} else {
		prefix = []byte(makerAddress.Hex() + "|" + makerFeeAssetAddress.Hex() + "|")
		if tokenID != nil {
			prefix = append(prefix, tokenID.Bytes()...)
		}
	}

	filter := m.Orders.MakerAddressMakerFeeAssetAddressTokenIDIndex.PrefixFilter(prefix)
	orders := []*Order{}
	if err := m.Orders.NewQuery(filter).Run(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersByMakerAddressAndMaxSalt finds all orders belonging to a particular maker address that
// also have a salt value less then or equal to X
func (m *MeshDB) FindOrdersByMakerAddressAndMaxSalt(makerAddress common.Address, salt *big.Int) ([]*Order, error) {
	// DB range queries exclude the limit value however the 0x protocol `cancelOrdersUpTo` method
	// is inclusive of the value supplied. In order to make this helper method more useful to our
	// particular use-case, we add 1 to the supplied salt (making the query inclusive instead)
	saltPlusOne := new(big.Int).Add(salt, big.NewInt(1))
	start := []byte(fmt.Sprintf("%s|%080s", makerAddress.Hex(), "0"))
	limit := []byte(fmt.Sprintf("%s|%s", makerAddress.Hex(), uint256ToConstantLengthBytes(saltPlusOne)))
	filter := m.Orders.MakerAddressAndSaltIndex.RangeFilter(start, limit)
	orders := []*Order{}
	if err := m.Orders.NewQuery(filter).Run(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// FindOrdersLastUpdatedBefore finds all orders where the LastUpdated time is less
// than X
func (m *MeshDB) FindOrdersLastUpdatedBefore(lastUpdated time.Time) ([]*Order, error) {
	start := []byte(time.Unix(0, 0).Format(time.RFC3339Nano))
	limit := []byte(lastUpdated.UTC().Format(time.RFC3339Nano))
	filter := m.Orders.LastUpdatedIndex.RangeFilter(start, limit)
	orders := []*Order{}
	if err := m.Orders.NewQuery(filter).Run(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// FindRemovedOrders finds all orders that have been flagged for removal
func (m *MeshDB) FindRemovedOrders() ([]*Order, error) {
	var removedOrders []*Order
	isRemovedFilter := m.Orders.IsRemovedIndex.ValueFilter([]byte{1})
	if err := m.Orders.NewQuery(isRemovedFilter).Run(&removedOrders); err != nil {
		return nil, err
	}
	return removedOrders, nil
}

// GetMetadata returns the metadata (or a db.NotFoundError if no metadata has been found).
func (m *MeshDB) GetMetadata() (*Metadata, error) {
	var metadata Metadata
	if err := m.metadata.FindByID([]byte{0}, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}

// SaveMetadata inserts the metadata into the database, overwriting any existing
// metadata.
func (m *MeshDB) SaveMetadata(metadata *Metadata) error {
	if err := m.metadata.Insert(metadata); err != nil {
		return err
	}
	return nil
}

// UpdateMetadata updates the metadata in the database via a transaction. It
// accepts a callback function which will be provided with the old metadata and
// should return the new metadata to save.
func (m *MeshDB) UpdateMetadata(updater func(oldmetadata Metadata) (newMetadata Metadata)) error {
	txn := m.metadata.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	oldMetadata, err := m.GetMetadata()
	if err != nil {
		return err
	}
	newMetadata := updater(*oldMetadata)
	if err := txn.Update(&newMetadata); err != nil {
		return err
	}

	return txn.Commit()
}

type singleAssetData struct {
	Address common.Address
	TokenID *big.Int
}

func parseContractAddressesAndTokenIdsFromAssetData(assetData []byte) ([]singleAssetData, error) {
	singleAssetDatas := []singleAssetData{}
	assetDataDecoder := zeroex.NewAssetDataDecoder()

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
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, id := range decodedAssetData.Ids {
			a := singleAssetData{
				Address: decodedAssetData.Address,
				TokenID: id,
			}
			singleAssetDatas = append(singleAssetDatas, a)
		}
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

func uint256ToConstantLengthBytes(v *big.Int) []byte {
	return []byte(fmt.Sprintf("%080s", v.String()))
}

// TrimOrdersByExpirationTime removes existing orders with the highest
// expiration time until the number of remaining orders is <= targetMaxOrders.
// It returns any orders that were removed and the new max expiration time that
// can be used to eliminate incoming orders that expire too far in the future.
func (m *MeshDB) TrimOrdersByExpirationTime(targetMaxOrders int) (newMaxExpirationTime *big.Int, removedOrders []*Order, err error) {
	txn := m.Orders.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	numOrders, err := m.Orders.Count()
	if err != nil {
		return nil, nil, err
	}
	if numOrders <= targetMaxOrders {
		// If the number of orders is less than the target, we don't need to remove
		// any orders. Return UnlimitedExpirationTime.
		return constants.UnlimitedExpirationTime, nil, nil
	}

	// Find the orders which we need to remove. We use a prefix filter of "0|: so
	// that we only remove non-pinned orders.
	filter := m.Orders.ExpirationTimeIndex.PrefixFilter([]byte("0|"))
	numOrdersToRemove := numOrders - targetMaxOrders
	if err := m.Orders.NewQuery(filter).Reverse().Max(numOrdersToRemove).Run(&removedOrders); err != nil {
		return nil, nil, err
	}

	// Remove those orders and commit the transaction.
	for _, order := range removedOrders {
		if err := txn.Delete(order.Hash.Bytes()); err != nil {
			return nil, nil, err
		}
	}
	if err := txn.Commit(); err != nil {
		return nil, nil, err
	}

	// If we could not remove numOrdersToRemove orders than it means the database
	// is full of pinned orders. We still remove as many orders as we can and then
	// return an error.
	if len(removedOrders) < numOrdersToRemove {
		return nil, nil, ErrDBFilledWithPinnedOrders
	}

	// The new max expiration time is simply the minimum expiration time of the
	// orders that were removed (i.e., the expiration time of the last order in
	// the slice). We add a buffer of -1 just to make sure we don't exceed
	// targetMaxOrders. This means it is technically possible that there are a
	// number of orders currently in the database that exceed the max expiration
	// time, but no new orders that exceed this time will be added.
	newMaxExpirationTime = removedOrders[len(removedOrders)-1].SignedOrder.ExpirationTimeSeconds
	newMaxExpirationTime = newMaxExpirationTime.Sub(newMaxExpirationTime, big.NewInt(1))
	return newMaxExpirationTime, removedOrders, nil
}

// CountPinnedOrders returns the number of pinned orders.
func (m *MeshDB) CountPinnedOrders() (int, error) {
	// We use a prefix filter of "1|" so that we only count pinned orders.
	filter := m.Orders.ExpirationTimeIndex.PrefixFilter([]byte("1|"))
	return m.Orders.NewQuery(filter).Count()
}
