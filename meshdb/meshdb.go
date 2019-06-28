package meshdb

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

// MeshDB instantiates the DB connection and creates all the collections used by the application
type MeshDB struct {
	Database    *db.DB
	maxOrders   int
	MiniHeaders *MiniHeadersCollection
	Orders      *OrdersCollection
	ETHBackings *ETHBackingsCollection
}

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

type ETHBacking struct {
	MakerAddress common.Address
	OrderCount   int
	// Note(albrow): The maximum value for float64 is ~1.8e308. The amount of Wei
	// in circulation is ~1e26. Technically we might lose some precision by
	// representing amounts in float64 (especially if we allow backing with other
	// tokens in the future), but it should be good enough for the purposes of ETH
	// backing validation.
	AmountInWei float64
}

func (eb *ETHBacking) ID() []byte {
	return eb.MakerAddress.Bytes()
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
	IsRemovedIndex                       *db.Index
}

type ETHBackingsCollection struct {
	*db.Collection
	WeiPerOrderIndex *db.Index
}

// NewMeshDB instantiates a new MeshDB instance
func NewMeshDB(path string, maxOrders int) (*MeshDB, error) {
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

	ethBackings, err := setupETHBackings(database)
	if err != nil {
		return nil, err
	}

	return &MeshDB{
		Database:    database,
		maxOrders:   maxOrders,
		MiniHeaders: miniHeaders,
		Orders:      orders,
		ETHBackings: ethBackings,
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

	isRemovedIndex := col.AddIndex("isRemoved", func(m db.Model) []byte {
		order := m.(*Order)
		// false = 0; true = 1
		if order.IsRemoved {
			return []byte{1}
		}
		return []byte{0}
	})

	return &OrdersCollection{
		Collection:                           col,
		MakerAddressTokenAddressTokenIDIndex: makerAddressTokenAddressTokenIDIndex,
		MakerAddressAndSaltIndex:             makerAddressAndSaltIndex,
		LastUpdatedIndex:                     lastUpdatedIndex,
		IsRemovedIndex:                       isRemovedIndex,
	}, nil
}

func setupMiniHeaders(database *db.DB) (*MiniHeadersCollection, error) {
	col, err := database.NewCollection("miniHeader", &MiniHeader{})
	if err != nil {
		return nil, err
	}
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

func setupETHBackings(database *db.DB) (*ETHBackingsCollection, error) {
	col, err := database.NewCollection("ethBackings", &ETHBacking{})
	if err != nil {
		return nil, err
	}
	weiPerOrderIndex := col.AddIndex("weiPerOrder", func(model db.Model) []byte {
		ethBacking := model.(*ETHBacking)
		return float64ToBytes(ethBacking.WeiPerOrder())
	})

	return &ETHBackingsCollection{
		Collection:       col,
		WeiPerOrderIndex: weiPerOrderIndex,
	}, nil
}

func float64ToBytes(number float64) []byte {
	// Recall that we must ensure the length of the numbers is always the same
	// for indexes. Here, we have allowed 5 digits of precision after the
	// decimal point. We pad with zeroes such that the length of the number is
	// always 86. (80 characters before the decimal point, followed by the
	// decimal point itself, followed by 5 characters after the decimal point).
	return []byte(fmt.Sprintf("%080.5f", number))
}

func (eb *ETHBacking) WeiPerOrder() float64 {
	if eb.OrderCount == 0 {
		// We can't divide by zero. Instead return max float.
		return math.MaxFloat64
	}
	return float64(eb.AmountInWei) / float64(eb.OrderCount)
}

// Close closes the database connection
func (m *MeshDB) Close() {
	m.Database.Close()
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

// InsertOrder atomically inserts the given order and updates any relevant ETH
// backings. It also removes the order with the lowest ETH backing if needed in
// order to make room. If order storage is full and the given order does not
// have a high enough ETH backing, it will not be stored (this is not considered
// an error). InsertOrder will return an error if any corresponding ETH backings
// cannot be found.
// TODO: return hash of order that was deleted.
func (m *MeshDB) InsertOrder(order *Order) error {
	txn := m.Database.OpenGlobalTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	totalExistingOrders, err := m.Orders.Count()
	if err != nil {
		return err
	}
	if totalExistingOrders > m.maxOrders {
		// This should never happen and indicates a bug in meshdb or the db package.
		// It is hard to fix the problem because we can't efficiently remove more
		// than one order with the lowest ETH backing per order.
		return fmt.Errorf("invalid database state: total number of orders (%d) is greater than max orders (%d)", totalExistingOrders, m.maxOrders)
	}

	var ethBackingForIncomingOrder ETHBacking
	if err := m.ETHBackings.FindByID(order.SignedOrder.MakerAddress.Bytes(), &ethBackingForIncomingOrder); err != nil {
		if _, ok := err.(db.NotFoundError); ok {
			// This should never happen because the ETHBacking should have already
			// been stored. We can't fix the problem here because meshdb doesn't know
			// about ETHWatcher and can't send the call to get new balances (if it did
			// it would violate separation of concerns).
			return fmt.Errorf("invalid database state: could not find ETHBacking for maker address: %s", order.SignedOrder.MakerAddress.Hex())
		} else {
			return err
		}
	}
	// updatedETHBacking is the hypothetical new ETHBacking for the incoming
	// order. It includes the order itself in the OrderCount.
	updatedETHBacking := &ETHBacking{
		MakerAddress: ethBackingForIncomingOrder.MakerAddress,
		OrderCount:   ethBackingForIncomingOrder.OrderCount + 1,
		AmountInWei:  ethBackingForIncomingOrder.AmountInWei,
	}

	if totalExistingOrders == m.maxOrders {
		log.WithField("maxOrders", m.maxOrders).Trace("Maximum order limit reached; deleting an order to make room")
		// In this case, we are going to hit the storage limit. Delete the order
		// with the lowest ETH backing in order to make room.
		lowestETHBacking, err := m.GetETHBackingWithLowestETHPerOrder()
		if err != nil {
			return err
		}
		if updatedETHBacking.WeiPerOrder() <= lowestETHBacking.WeiPerOrder() {
			// If the incoming order is associated with an ETH backing that is not
			// greater than the current lowest ETH backing, don't store it.
			return nil
		}
		log.WithField("makerAddress", lowestETHBacking.MakerAddress.Hex()).Trace("Found maker with lowest ETH backing per order")
		randomOrderForMaker, err := m.getRandomOrderForMaker(lowestETHBacking.MakerAddress)
		if err != nil {
			return err
		}
		// Delete the order and decrement the OrderCount for the corresponding ETH
		// backing.
		if err := txn.Delete(m.Orders.Collection, randomOrderForMaker.ID()); err != nil {
			return err
		}
		lowestETHBacking.OrderCount = lowestETHBacking.OrderCount - 1
		if err := txn.Update(m.ETHBackings.Collection, lowestETHBacking); err != nil {
			return err
		}
	}

	// Insert the order and update the ETH backing for the corresponding maker
	if err := txn.Insert(m.Orders.Collection, order); err != nil {
		return err
	}
	if err := txn.Update(m.ETHBackings.Collection, updatedETHBacking); err != nil {
		return err
	}

	// Commit the transaction.
	return txn.Commit()
}

func (m *MeshDB) GetETHBackingWithLowestETHPerOrder() (*ETHBacking, error) {
	filter := m.ETHBackings.WeiPerOrderIndex.All()
	queryResults := []*ETHBacking{}
	if err := m.ETHBackings.NewQuery(filter).Max(1).Run(&queryResults); err != nil {
		return nil, err
	}
	if len(queryResults) == 0 {
		return nil, errors.New("query for lowest ETHBacking by ETH per order returned no results")
	}
	return queryResults[0], nil
}

func (m *MeshDB) GetETHBackingsWithLowestETHPerOrder(count int) ([]*ETHBacking, error) {
	filter := m.ETHBackings.WeiPerOrderIndex.All()
	queryResults := []*ETHBacking{}
	if err := m.ETHBackings.NewQuery(filter).Max(1).Run(&queryResults); err != nil {
		return nil, err
	}
	return queryResults, nil
}

func (m *MeshDB) getRandomOrderForMaker(makerAddress common.Address) (*Order, error) {
	// TODO(albrow): Currently this always returns the order with the lowest hash,
	// but ideally it should be completely random.
	query := m.newQueryForOrdersByMakerAddress(makerAddress)
	orders := []*Order{}
	if err := query.Max(1).Run(&orders); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("query for orders by maker addressreturned no results")
	}
	return orders[0], nil
}

// FindOrdersByMakerAddress finds all orders belonging to a particular maker address
func (m *MeshDB) FindOrdersByMakerAddress(makerAddress common.Address) ([]*Order, error) {
	query := m.newQueryForOrdersByMakerAddress(makerAddress)
	orders := []*Order{}
	if err := query.Run(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (m *MeshDB) newQueryForOrdersByMakerAddress(makerAddress common.Address) *db.Query {
	prefix := []byte(makerAddress.Hex() + "|")
	filter := m.Orders.MakerAddressTokenAddressTokenIDIndex.PrefixFilter(prefix)
	return m.Orders.NewQuery(filter)
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
	err := m.Orders.NewQuery(filter).Run(&orders)
	if err != nil {
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
	start := []byte(time.Unix(0, 0).Format(time.RFC3339Nano))
	limit := []byte(lastUpdated.UTC().Format(time.RFC3339Nano))
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
