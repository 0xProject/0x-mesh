package meshdb

import (
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
	ETHBalance   *big.Int
}

func (eb *ETHBacking) ID() []byte {
	return eb.MakerAddress.Bytes()
}

// MeshDB instantiates the DB connection and creates all the collections used by the application
type MeshDB struct {
	database    *db.DB
	maxOrders   int
	MiniHeaders *MiniHeadersCollection
	Orders      *OrdersCollection
	ETHBackings *ETHBackingsCollection
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
	ETHPerOrderIndex *db.Index
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
		database:    database,
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
	ethPerOrderIndex := col.AddIndex("ethPerOrder", func(model db.Model) []byte {
		ethBacking := model.(*ETHBacking)
		return ratToBytes(ethBacking.ethPerOrder())
	})

	return &ETHBackingsCollection{
		Collection:       col,
		ETHPerOrderIndex: ethPerOrderIndex,
	}, nil
}

func ratToBytes(rat *big.Rat) []byte {
	// Recall that we must ensure the length of the numbers is always the same
	// for indexes. Here, we have allowed 5 digits of precision after the
	// decimal point. We pad with zeroes such that the length of the number is
	// always 86. (80 characters before the decimal point, followed by the
	// decimal point itself, followed by 5 characters after the decimal point).
	return []byte(fmt.Sprintf("%080.5s", rat.FloatString(5)))
}

func (eb *ETHBacking) ethPerOrder() *big.Rat {
	ethPerOrder := big.NewRat(0, 0)
	ethPerOrder.SetFrac(eb.ETHBalance, big.NewInt(int64(eb.OrderCount)))
	return ethPerOrder
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

func (m *MeshDB) findEthBackingsWithLessEthPerOrder(target *big.Rat, max int) ([]*ETHBacking, error) {
	filter := m.ETHBackings.ETHPerOrderIndex.RangeFilter([]byte{}, ratToBytes(target))
	query := m.ETHBackings.NewQuery(filter).Max(max)
	var result []*ETHBacking
	if err := query.Run(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// ISSUES/QUESTIONS:
// - Difficult to maintain consistency because there are two collections involved. (Could fix by adding global transactions to db package?)
// - Either meshdb or orderwatch package needs to know how to retrieve balances. Violates separation of concerns.
// - Is there a way to batch balance requests?
// - Writes within the transaction don't take effect until after transaction is committed. Need to duplicate some work in memory.
// - How to count orders with IsRemoved = true?
// - Figuring out which orders to delete is insanely complicated. Each order we insert/delete changes the calculation for the ETH per order, so we have to recompute it every time.
//

// TODO(albrow):
// 1. Multi-collection or global transactions.
// 2. More effecient count method (cache in a key).
// 3. Check ETH balances upon receipt of the orders before doing on-chain validation.
// 4. Insert orders one at a time. This simplifies the algorithm at the cost of effeciency. We can sort incoming orders by ETH backing per order which should help.
// 5. In the core package, receive events from ETH balance watcher and update all ETHBackings in a single transaction.
// 6. Update ETHBackings when we remove an order.

func (m *MeshDB) InsertOrders(orders []*Order) error {
	ordersTxn := m.Orders.OpenTransaction()
	defer ordersTxn.Discard()
	ethBackingsTxn := m.ETHBackings.OpenTransaction()
	defer ethBackingsTxn.Discard()

	// TODO(albrow): This query has runtime of O(N) where N is the number of
	// orders stored. Could be optimized.
	totalExistingOrders, err := m.Orders.NewQuery(m.Orders.IsRemovedIndex.All()).Count()
	if err != nil {
		return err
	}

	// Create a map of maker address to orders for that maker address.
	ordersByMakerAddress := map[string][]*Order{}
	for _, order := range orders {
		makerAddress := order.SignedOrder.MakerAddress
		if orders, found := ordersByMakerAddress[makerAddress.String()]; found {
			orders = append(orders, order)
			ordersByMakerAddress[makerAddress.String()] = orders
		} else {
			ordersByMakerAddress[makerAddress.String()] = []*Order{order}
		}
	}

	if totalExistingOrders+len(orders) > m.maxOrders {
		// In this case, we are going to hit the storage limit. We cannot store all
		// the given orders without first deleting some existing ones.

		// Find all ETH backings for the given orders.
		makerAddressToEthBacking := map[string]*ETHBacking{}
		leastBacking := big.NewRat(math.MaxInt64, 1)
		for makerAddressStr, orders := range ordersByMakerAddress {
			makerAddress := common.HexToAddress(makerAddressStr)
			// Check if this maker address has an existing ETH backing.
			var existingETHBacking ETHBacking
			if err := m.ETHBackings.FindByID(makerAddress.Bytes(), &existingETHBacking); err != nil {
				if _, ok := err.(db.NotFoundError); !ok {
					// An unexpected error occurred. Return it.
					return err
				}
				// No backing has been stored for this maker address. Make a new one.
				// TODO(albrow): Get the actual initial balance.
				// TODO(albrow): Add this maker address to the ETH watcher.
				ethBacking := &ETHBacking{
					MakerAddress: makerAddress,
					OrderCount:   len(orders),
					ETHBalance:   big.NewInt(0),
				}
				makerAddressToEthBacking[makerAddressStr] = ethBacking
				if ethPerOrder := ethBacking.ethPerOrder(); ethPerOrder.Cmp(leastBacking) == -1 {
					leastBacking = ethPerOrder
				}
			} else {
				// A backing was stored for this maker address. Update it.
				ethBacking := &ETHBacking{
					MakerAddress: makerAddress,
					OrderCount:   existingETHBacking.OrderCount + len(orders),
					ETHBalance:   existingETHBacking.ETHBalance,
				}
				makerAddressToEthBacking[makerAddressStr] = ethBacking
				if ethPerOrder := ethBacking.ethPerOrder(); ethPerOrder.Cmp(leastBacking) == -1 {
					leastBacking = ethPerOrder
				}
			}
		}

		// Find any existing orders with a lower backing than the least backing in
		// the given orders. These need to be deleted.
		backingsToDelete, err := m.findEthBackingsWithLessEthPerOrder(leastBacking, len(orders))
		if err != nil {
			return err
		}
		numberOfOrdersToDelete := 0
		for _, backingToDelete := range backingsToDelete {
			numberOfOrdersToDelete += backingToDelete.OrderCount
		}
		if numberOfOrdersToDelete == 0 {
			// There are no orders with a lesser ETH backing than the given orders.
			// In other words, *none* of the given orders have enough ETH backing to
			// be stored, so we don't store any of them. Just return immediately.
			return nil
		} else if numberOfOrdersToDelete >= len(orders) {
			// The number of orders which could be deleted is *greater than* the
			// number of orders we are trying to store. This means we can store all
			// the given orders and delete up to len(orders) from the set of orders
			// that have a lesser ETH backing.

			// Queue operations to delete any existing orders with a lesser ETH
			// backing.
			// TODO(albrow): How can we do this correctly?
			// for _, backingToDelete := range backingsToDelete {
			//
			// }

			// Queue operations to store all the given orders and insert/update their
			// corresponding ETH backings.
			for _, order := range orders {
				if err := ordersTxn.Insert(order); err != nil {
					return err
				}
			}
			for _, ethBacking := range makerAddressToEthBacking {
				var eb ETHBacking
				if err := m.ETHBackings.FindByID(ethBacking.MakerAddress.Bytes(), &eb); err != nil {
					if _, notFound := err.(db.NotFoundError); notFound {
						// The backing *doesn't* exist in the db yet; we need to insert it.
						if err := ethBackingsTxn.Insert(ethBacking); err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					// The backing *does* exist in the db already; we need to update it.
					if err := ethBackingsTxn.Update(ethBacking); err != nil {
						return err
					}
				}
			}
		} else {
			// The number of orders which could be deleted is *less than* the number
			// of orders we are trying to store. This means we can store some of the
			// given orders but not all of them.
		}

	} else {
		// In this case, we can insert all the given orders without hitting the
		// storage limit.
		for makerAddressStr, orders := range ordersByMakerAddress {
			makerAddress := common.HexToAddress(makerAddressStr)
			// Check if this maker address has an existing ETH backing.
			var existingETHBacking ETHBacking
			if err := m.ETHBackings.FindByID(makerAddress.Bytes(), &existingETHBacking); err != nil {
				if _, ok := err.(db.NotFoundError); !ok {
					// An unexpected error occurred. Return it.
					return err
				}
				// No backing has been stored for this maker address. Insert a new one.
				// TODO(albrow): Get the actual initial balance.
				// TODO(albrow): Add this maker address to the ETH watcher.
				ethBackingsTxn.Insert(&ETHBacking{
					MakerAddress: makerAddress,
					OrderCount:   len(orders),
					ETHBalance:   big.NewInt(0),
				})
			} else {
				// A backing was found for this maker address. We need to update it.
				existingETHBacking.OrderCount += len(orders)
				ethBackingsTxn.Update(&existingETHBacking)
			}

			// Queue an operation to insert each order.
			for _, order := range orders {
				if err := ordersTxn.Insert(order); err != nil {
					return err
				}
			}
		}
	}

	// Commit each transaction.
	if err := ethBackingsTxn.Commit(); err != nil {
		// TODO(albrow): We lose consistency guarantees if this happens.
		panic(err)
	}
	if err := ordersTxn.Commit(); err != nil {
		// TODO(albrow): We lose consistency guarantees if this happens.
		panic(err)
	}
	return nil
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
