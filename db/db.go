package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB instantiates the DB connection and creates all the collections used by the application
type DB struct {
	ctx   context.Context
	sqldb *sqlx.DB
	// MiniHeaderRetentionLimit int
}

// New creates a new connection to the database. The connection will be automatically closed
// when the given context is canceled.
func New(ctx context.Context, path string) (*DB, error) {
	connectCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := os.MkdirAll(path, os.ModePerm); err != nil && err != os.ErrExist {
		return nil, err
	}

	db, err := sqlx.ConnectContext(connectCtx, "sqlite3", filepath.Join(path, "db.sqlite"))
	if err != nil {
		return nil, err
	}

	// Automatically close the database connection when the context is canceled.
	go func() {
		select {
		case <-ctx.Done():
			_ = db.Close()
		}
	}()

	return &DB{
		ctx:   ctx,
		sqldb: db,
	}, nil
}

// TODO(albrow): Use a proper migration tool.
// TODO(albrow): Add indexes.
const schema = `
CREATE TABLE IF NOT EXISTS orders (
	hash                     TEXT UNIQUE NOT NULL,
	chainID                  NUMERIC(78, 0) NOT NULL,
	exchangeAddress          TEXT NOT NULL,
	makerAddress             TEXT NOT NULL,
	makerAssetData           TEXT NOT NULL,
	makerFeeAssetData        TEXT NOT NULL,
	makerAssetAmount         NUMERIC(78, 0) NOT NULL,
	makerFee                 NUMERIC(78, 0) NOT NULL,
	takerAddress             TEXT NOT NULL,
	takerAssetData           TEXT NOT NULL,
	takerFeeAssetData        TEXT NOT NULL,
	takerAssetAmount         NUMERIC(78, 0) NOT NULL,
	takerFee                 NUMERIC(78, 0) NOT NULL,
	senderAddress            TEXT NOT NULL,
	feeRecipientAddress      TEXT NOT NULL,
	expirationTimeSeconds    NUMERIC(78, 0) NOT NULL,
	salt                     NUMERIC(78, 0) NOT NULL,
	signature                TEXT NOT NULL,
	lastUpdated              DATETIME NOT NULL,
	fillableTakerAssetAmount NUMERIC(78, 0) NOT NULL,
	isRemoved                BOOLEAN NOT NULL,
	isPinned                 BOOLEAN NOT NULL,
	parsedMakerAssetData    TEXT NOT NULL,
	parsedMakerFeeAssetData TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS miniHeaders (
	hash      TEXT UNIQUE NOT NULL,
	parent    TEXT NOT NULL,
	number    NUMERIC(78, 0) NOT NULL,
	timestamp DATETIME NOT NULL,
	logs      TEXT NOT NULL
);
`

// TODO(albrow): Used prepared statement for inserts.
const insertOrderQuery = `INSERT OR IGNORE INTO orders (
	hash,
	chainID,
	exchangeAddress,
	makerAddress,
	makerAssetData,
	makerFeeAssetData,
	makerAssetAmount,
	makerFee,
	takerAddress,
	takerAssetData,
	takerFeeAssetData,
	takerAssetAmount,
	takerFee,
	senderAddress,
	feeRecipientAddress,
	expirationTimeSeconds,
	salt,
	signature,
	lastUpdated,
	fillableTakerAssetAmount,
	isRemoved,
	isPinned,
	parsedMakerAssetData,
	parsedMakerFeeAssetData
) VALUES (
	:hash,
	:chainID,
	:exchangeAddress,
	:makerAddress,
	:makerAssetData,
	:makerFeeAssetData,
	:makerAssetAmount,
	:makerFee,
	:takerAddress,
	:takerAssetData,
	:takerFeeAssetData,
	:takerAssetAmount,
	:takerFee,
	:senderAddress,
	:feeRecipientAddress,
	:expirationTimeSeconds,
	:salt,
	:signature,
	:lastUpdated,
	:fillableTakerAssetAmount,
	:isRemoved,
	:isPinned,
	:parsedMakerAssetData,
	:parsedMakerFeeAssetData
)
`

const updateOrderQuery = `UPDATE orders SET
	chainID = :chainID,
	exchangeAddress = :exchangeAddress,
	makerAddress = :makerAddress,
	makerAssetData = :makerAssetData,
	makerFeeAssetData = :makerFeeAssetData,
	makerAssetAmount = :makerAssetAmount,
	makerFee = :makerFee,
	takerAddress = :takerAddress,
	takerAssetData = :takerAssetData,
	takerFeeAssetData = :takerFeeAssetData,
	takerAssetAmount = :takerAssetAmount,
	takerFee = :takerFee,
	senderAddress = :senderAddress,
	feeRecipientAddress = :feeRecipientAddress,
	expirationTimeSeconds = :expirationTimeSeconds,
	salt = :salt,
	signature = :signature,
	lastUpdated = :lastUpdated,
	fillableTakerAssetAmount = :fillableTakerAssetAmount,
	isRemoved = :isRemoved,
	isPinned = :isPinned,
	parsedMakerAssetData = :parsedMakerAssetData,
	parsedMakerFeeAssetData = :parsedMakerFeeAssetData
WHERE orders.hash = :hash
`

const insertMiniHeaderQuery = `INSERT OR IGNORE INTO miniHeaders (
	hash,
	parent,
	number,
	timestamp,
	logs
) VALUES (
	:hash,
	:parent,
	:number,
	:timestamp,
	:logs
)`

func (db *DB) migrate() error {
	_, err := db.sqldb.ExecContext(db.ctx, schema)
	return err
}

func (db *DB) Close() error {
	panic(errors.New("Not implemented. Cancel the context instead."))
}

func (db *DB) AddOrders(orders []*Order) (added []*Order, removed []*Order, err error) {
	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	for _, order := range orders {
		result, err := txn.NamedExecContext(db.ctx, insertOrderQuery, order)
		if err != nil {
			return nil, nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return nil, nil, err
		}
		if affected > 0 {
			added = append(added, order)
		}
	}
	if err := txn.Commit(); err != nil {
		return nil, nil, err
	}

	// TODO(albrow): Remove orders with longest expiration time.
	// TODO(albrow): Return appropriate values for added, removed.
	return added, nil, nil
}

func (db *DB) GetOrder(hash common.Hash) (*Order, error) {
	var order Order
	if err := db.sqldb.GetContext(db.ctx, &order, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return nil, err
	}
	return &order, nil
}

type SortDirection string

const (
	Ascending  SortDirection = "ASC"
	Descending SortDirection = "DESC"
)

type FilterKind string

const (
	Equal          FilterKind = "="
	NotEqual       FilterKind = "!="
	Less           FilterKind = "<"
	Greater        FilterKind = ">"
	LessOrEqual    FilterKind = "<="
	GreaterOrEqual FilterKind = ">="
	Contains       FilterKind = "CONTAINS"
	// TODO(albrow): Starts with?
)

type OrderField string

const (
	OFHash                     OrderField = "hash"
	OFChainID                  OrderField = "chainID"
	OFExchangeAddress          OrderField = "exchangeAddress"
	OFMakerAddress             OrderField = "makerAddress"
	OFMakerAssetData           OrderField = "makerAssetData"
	OFMakerFeeAssetData        OrderField = "makerFeeAssetData"
	OFMakerAssetAmount         OrderField = "makerAssetAmount"
	OFMakerFee                 OrderField = "makerFee"
	OFTakerAddress             OrderField = "takerAddress"
	OFTakerAssetData           OrderField = "takerAssetData"
	OFTakerFeeAssetData        OrderField = "takerFeeAssetData"
	OFTakerAssetAmount         OrderField = "takerAssetAmount"
	OFTakerFee                 OrderField = "takerFee"
	OFSenderAddress            OrderField = "senderAddress"
	OFFeeRecipientAddress      OrderField = "feeRecipientAddress"
	OFExpirationTimeSeconds    OrderField = "expirationTimeSeconds"
	OFSalt                     OrderField = "salt"
	OFSignature                OrderField = "signature"
	OFLastUpdated              OrderField = "lastUpdated"
	OFFillableTakerAssetAmount OrderField = "fillableTakerAssetAmount"
	OFIsRemoved                OrderField = "isRemoved"
	OFIsPinned                 OrderField = "isPinned"
	OFParsedMakerAssetData     OrderField = "parsedMakerAssetData"
	OFParsedMakerFeeAssetData  OrderField = "parsedMakerFeeAssetData"
)

type FindOrdersOpts struct {
	Filters []OrderFilter
	Sort    []OrderSort
	Limit   uint
	Offset  uint
}

type OrderSort struct {
	Field     OrderField
	Direction SortDirection
}

type OrderFilter struct {
	Field OrderField
	Kind  FilterKind
	Value interface{}
}

// IncludesMakerAssetData is a helper method which returns a filter that will match orders
// that include the given asset data in MakerAssetData.
func IncludesMakerAssetData(tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	return OrderFilter{
		Field: OFParsedMakerAssetData,
		Kind:  Contains,
		Value: fmt.Sprintf(`{"address":"%s","tokenID":"%s"}`, strings.ToLower(tokenAddress.Hex()), tokenID.String()),
	}
}

// IncludesMakerFeeAssetData is a helper method which returns a filter that will match orders
// that include the given asset data in MakerFeeAssetData.
func IncludesMakerFeeAssetData(tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	return OrderFilter{
		Field: OFParsedMakerFeeAssetData,
		Kind:  Contains,
		Value: fmt.Sprintf(`{"address":"%s","tokenID":"%s"}`, strings.ToLower(tokenAddress.Hex()), tokenID.String()),
	}
}

func (db *DB) FindOrders(opts *FindOrdersOpts) ([]*Order, error) {
	query, err := db.findOrdersQueryFromOpts(opts)
	if err != nil {
		return nil, err
	}
	var orders []*Order
	rawQuery, bindings := query.ToSQL(false)
	fmt.Println(rawQuery, bindings)
	if err := query.GetAllContext(db.ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (db *DB) findOrdersQueryFromOpts(opts *FindOrdersOpts) (*sqlz.SelectStmt, error) {
	query := sqlz.Newx(db.sqldb).Select("*").From("orders")
	if opts == nil {
		return query, nil
	}

	ordering := orderingFromOrderSortOpts(opts.Sort)
	if len(ordering) != 0 {
		query.OrderBy(ordering...)
	}
	if opts.Limit != 0 {
		query.Limit(int64(opts.Limit))
	}
	if opts.Offset != 0 {
		if opts.Limit == 0 {
			return nil, errors.New("db.FindOrders: can't use Offset without Limit")
		}
		query.Offset(int64(opts.Offset))
	}
	whereConditions, err := whereConditionsFromOrderFilterOpts(opts.Filters)
	if err != nil {
		return nil, err
	}
	if len(whereConditions) != 0 {
		query.Where(whereConditions...)
	}

	return query, nil
}

func orderingFromOrderSortOpts(opts []OrderSort) []sqlz.SQLStmt {
	ordering := []sqlz.SQLStmt{}
	for _, sortOpt := range opts {
		if sortOpt.Direction == Ascending {
			ordering = append(ordering, sqlz.Asc(string(sortOpt.Field)))
		} else {
			ordering = append(ordering, sqlz.Desc(string(sortOpt.Field)))
		}
	}
	return ordering
}

func whereConditionsFromOrderFilterOpts(opts []OrderFilter) ([]sqlz.WhereCondition, error) {
	// TODO(albrow): Type-check on value? You can't use CONTAINS with numeric types.
	whereConditions := make([]sqlz.WhereCondition, len(opts))
	for i, filterOpt := range opts {
		switch filterOpt.Kind {
		case Equal:
			whereConditions[i] = sqlz.Eq(string(filterOpt.Field), filterOpt.Value)
		case NotEqual:
			whereConditions[i] = sqlz.Not(sqlz.Eq(string(filterOpt.Field), filterOpt.Value))
		case Less:
			whereConditions[i] = sqlz.Lt(string(filterOpt.Field), filterOpt.Value)
		case Greater:
			whereConditions[i] = sqlz.Gt(string(filterOpt.Field), filterOpt.Value)
		case LessOrEqual:
			whereConditions[i] = sqlz.Lte(string(filterOpt.Field), filterOpt.Value)
		case GreaterOrEqual:
			whereConditions[i] = sqlz.Gte(string(filterOpt.Field), filterOpt.Value)
		case Contains:
			// TODO(albrow): Value cannot contain special characters like "%".
			// TODO(albrow): Optimize this so it is easier to index.
			whereConditions[i] = sqlz.Like(string(filterOpt.Field), fmt.Sprintf("%%%s%%", filterOpt.Value))
		default:
			return nil, fmt.Errorf("db.FindOrder: unknown FilterOpt.Kind: %s", filterOpt.Kind)
		}
	}
	return whereConditions, nil
}

func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *Order) (updatedOrder *Order, err error)) error {
	if updateFunc == nil {
		return errors.New("db.UpdateOrders: updateFunc cannot be nil")
	}

	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	var existingOrder Order
	if err := txn.GetContext(db.ctx, &existingOrder, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return err
	}

	updatedOrder, err := updateFunc(&existingOrder)
	if err != nil {
		return fmt.Errorf("db.UpdateOrders: updateFunc returned error")
	}
	_, err = txn.NamedExecContext(db.ctx, updateOrderQuery, updatedOrder)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func (db *DB) AddMiniHeaders(miniHeaders []*MiniHeader) (added []*MiniHeader, removed []*MiniHeader, err error) {
	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	for _, miniHeader := range miniHeaders {
		result, err := txn.NamedExecContext(db.ctx, insertMiniHeaderQuery, miniHeader)
		if err != nil {
			return nil, nil, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return nil, nil, err
		}
		if affected > 0 {
			added = append(added, miniHeader)
		}
	}
	if err := txn.Commit(); err != nil {
		return nil, nil, err
	}

	// TODO(albrow): Remove miniheaders to keep count low.
	// TODO(albrow): Return appropriate values for added, removed.
	return added, nil, nil
}

func (db *DB) GetMiniHeader(hash common.Hash) (*MiniHeader, error) {
	var miniHeader MiniHeader
	if err := db.sqldb.GetContext(db.ctx, &miniHeader, "SELECT * FROM miniHeaders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return nil, err
	}
	return &miniHeader, nil
}

// TODO(albrow): Add options for filtering, sorting, limit, and offset.
func (db *DB) FindMiniHeaders() ([]*MiniHeader, error) {
	var miniHeaders []*MiniHeader
	if err := db.sqldb.SelectContext(db.ctx, &miniHeaders, "SELECT * FROM miniHeaders"); err != nil {
		return nil, err
	}
	return miniHeaders, nil
}

// FindAllMiniHeadersSortedByNumber returns all MiniHeaders sorted in ascending block number order
func (db *DB) FindAllMiniHeadersSortedByNumber() ([]*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// FindLatestMiniHeader returns the latest MiniHeader (i.e. the one with the
// largest block number). It returns nil, MiniHeaderCollectionEmptyError if there
// are no MiniHeaders in the database.
func (db *DB) FindLatestMiniHeader() (*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// FindMiniHeaderByBlockNumber returns the MiniHeader with the specified block number
func (db *DB) FindMiniHeaderByBlockNumber(blockNumber *big.Int) (*miniheader.MiniHeader, error) {
	return nil, errors.New("Not yet implemented")
}

// UpdateMiniHeaderRetentionLimit updates the MiniHeaderRetentionLimit. This is only used by tests in order
// to set the retention limit to a smaller size, making the tests shorter in length
func (db *DB) UpdateMiniHeaderRetentionLimit(limit int) error {
	return errors.New("Not yet implemented")
}

// PruneMiniHeadersAboveRetentionLimit prunes miniHeaders from the DB that are above the retention limit
func (db *DB) PruneMiniHeadersAboveRetentionLimit() error {
	return errors.New("Not yet implemented")
}

// ClearAllMiniHeaders removes all stored MiniHeaders from the database.
func (db *DB) ClearAllMiniHeaders() error {
	return errors.New("Not yet implemented")
}

// ClearOldMiniHeaders removes all stored MiniHeaders with a block number less then
// the given minBlockNumber.
func (db *DB) ClearOldMiniHeaders(minBlockNumber *big.Int) error {
	return errors.New("Not yet implemented")
}

// FindOrdersByMakerAddress finds all orders belonging to a particular maker address
// ✅
func (db *DB) FindOrdersByMakerAddress(makerAddress common.Address) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressTokenAddressAndTokenID finds all orders belonging to a particular maker
// address where makerAssetData encodes for a particular token contract and optionally a token ID
// ✅
func (db *DB) FindOrdersByMakerAddressTokenAddressAndTokenID(makerAddress, tokenAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressMakerFeeAssetAddressTokenID finds all orders belonging to
// a particular maker address where makerFeeAssetData encodes for a particular
// token contract and optionally a token ID. To find orders without a maker fee,
// use constants.NullAddress for makerFeeAssetAddress.
// ✅
func (db *DB) FindOrdersByMakerAddressMakerFeeAssetAddressAndTokenID(makerAddress, makerFeeAssetAddress common.Address, tokenID *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersByMakerAddressAndMaxSalt finds all orders belonging to a particular maker address that
// also have a salt value less then or equal to X
// ✅
func (db *DB) FindOrdersByMakerAddressAndMaxSalt(makerAddress common.Address, salt *big.Int) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindOrdersLastUpdatedBefore finds all orders where the LastUpdated time is less
// than X
// ✅
func (db *DB) FindOrdersLastUpdatedBefore(lastUpdated time.Time) ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// FindRemovedOrders finds all orders that have been flagged for removal
// ✅
func (db *DB) FindRemovedOrders() ([]*Order, error) {
	return nil, errors.New("Not yet implemented")
}

// GetMetadata returns the metadata (or a db.NotFoundError if no metadata has been found).
func (db *DB) GetMetadata() (*Metadata, error) {
	return nil, errors.New("Not yet implemented")
}

// SaveMetadata inserts the metadata into the database, overwriting any existing
// metadata.
func (db *DB) SaveMetadata(metadata *Metadata) error {
	return errors.New("Not yet implemented")
}

// UpdateMetadata updates the metadata in the database via a transaction. It
// accepts a callback function which will be provided with the old metadata and
// should return the new metadata to save.
func (db *DB) UpdateMetadata(updater func(oldmetadata Metadata) (newMetadata Metadata)) error {
	return errors.New("Not yet implemented")
}

// CountPinnedOrders returns the number of pinned orders.
func (db *DB) CountPinnedOrders() (int, error) {
	return 0, errors.New("Not yet implemented")
}
