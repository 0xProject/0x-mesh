// +build !js

package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/sqltypes"
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

	sqldb, err := sqlx.ConnectContext(connectCtx, "sqlite3", filepath.Join(path, "db.sqlite"))
	if err != nil {
		return nil, err
	}

	// Automatically close the database connection when the context is canceled.
	go func() {
		select {
		case <-ctx.Done():
			_ = sqldb.Close()
		}
	}()

	db := &DB{
		ctx:   ctx,
		sqldb: sqldb,
	}
	if err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
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

func (db *DB) AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	for _, order := range orders {
		result, err := txn.NamedExecContext(db.ctx, insertOrderQuery, sqltypes.OrderFromCommonType(order))
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

func (db *DB) GetOrder(hash common.Hash) (*types.OrderWithMetadata, error) {
	var order sqltypes.Order
	if err := db.sqldb.GetContext(db.ctx, &order, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return nil, err
	}
	return sqltypes.OrderToCommonType(&order), nil
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

func (db *DB) FindOrders(opts *FindOrdersOpts) ([]*types.OrderWithMetadata, error) {
	query, err := db.findOrdersQueryFromOpts(opts)
	if err != nil {
		return nil, err
	}
	var orders []*sqltypes.Order
	rawQuery, bindings := query.ToSQL(false)
	fmt.Println(rawQuery, bindings)
	if err := query.GetAllContext(db.ctx, &orders); err != nil {
		return nil, err
	}
	return sqltypes.OrdersToCommonType(orders), nil
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
		value := convertFilterValue(filterOpt.Value)
		switch filterOpt.Kind {
		case Equal:
			whereConditions[i] = sqlz.Eq(string(filterOpt.Field), value)
		case NotEqual:
			whereConditions[i] = sqlz.Not(sqlz.Eq(string(filterOpt.Field), value))
		case Less:
			whereConditions[i] = sqlz.Lt(string(filterOpt.Field), value)
		case Greater:
			whereConditions[i] = sqlz.Gt(string(filterOpt.Field), value)
		case LessOrEqual:
			whereConditions[i] = sqlz.Lte(string(filterOpt.Field), value)
		case GreaterOrEqual:
			whereConditions[i] = sqlz.Gte(string(filterOpt.Field), value)
		case Contains:
			// TODO(albrow): Value cannot contain special characters like "%".
			// TODO(albrow): Optimize this so it is easier to index.
			whereConditions[i] = sqlz.Like(string(filterOpt.Field), fmt.Sprintf("%%%s%%", value))
		default:
			return nil, fmt.Errorf("db.FindOrder: unknown FilterOpt.Kind: %s", filterOpt.Kind)
		}
	}
	return whereConditions, nil
}

func (db *DB) DeleteOrder(hash common.Hash) error {
	if _, err := db.sqldb.ExecContext(db.ctx, "DELETE FROM orders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return err
	}
	return nil
}

type DeleteOrdersOpts struct {
	Filters []OrderFilter
}

// TODO(albrow): Return orders that were deleted?
func (db *DB) DeleteOrders(opts *DeleteOrdersOpts) error {
	query, err := db.deleteOrdersQueryFromOpts(opts)
	if err != nil {
		return err
	}
	rawQuery, bindings := query.ToSQL(false)
	fmt.Println(rawQuery, bindings)
	if _, err := query.ExecContext(db.ctx); err != nil {
		return err
	}
	return nil
}

func (db *DB) deleteOrdersQueryFromOpts(opts *DeleteOrdersOpts) (*sqlz.DeleteStmt, error) {
	query := sqlz.Newx(db.sqldb).DeleteFrom("orders")
	if opts == nil {
		return query, nil
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

// TODO(albrow): Consider automatically setting LastUpdated?
func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) error {
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

	var existingOrder sqltypes.Order
	if err := txn.GetContext(db.ctx, &existingOrder, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return err
	}

	commonOrder := sqltypes.OrderToCommonType(&existingOrder)
	commonUpdatedOrder, err := updateFunc(commonOrder)
	if err != nil {
		return fmt.Errorf("db.UpdateOrders: updateFunc returned error")
	}
	updatedOrder := sqltypes.OrderFromCommonType(commonUpdatedOrder)
	_, err = txn.NamedExecContext(db.ctx, updateOrderQuery, updatedOrder)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func (db *DB) AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error) {
	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	for _, miniHeader := range miniHeaders {
		result, err := txn.NamedExecContext(db.ctx, insertMiniHeaderQuery, sqltypes.MiniHeaderFromCommonType(miniHeader))
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

func (db *DB) GetMiniHeader(hash common.Hash) (*types.MiniHeader, error) {
	var miniHeader sqltypes.MiniHeader
	if err := db.sqldb.GetContext(db.ctx, &miniHeader, "SELECT * FROM miniHeaders WHERE hash = $1", hash); err != nil {
		// TODO(albrow): Specifically handle not found error.
		// - Maybe wrap other types of errors for consistency with Dexie.js implementation?
		return nil, err
	}
	return sqltypes.MiniHeaderToCommonType(&miniHeader), nil
}

type MiniHeaderField string

const (
	MFHash      MiniHeaderField = "hash"
	MFParent    MiniHeaderField = "parent"
	MFNumber    MiniHeaderField = "number"
	MFTimestamp MiniHeaderField = "timestamp"
	MFLogs      MiniHeaderField = "logs"
)

type FindMiniHeadersOpts struct {
	Filters []MiniHeaderFilter
	Sort    []MiniHeaderSort
	Limit   uint
	Offset  uint
}

type MiniHeaderSort struct {
	Field     MiniHeaderField
	Direction SortDirection
}

type MiniHeaderFilter struct {
	Field MiniHeaderField
	Kind  FilterKind
	Value interface{}
}

// TODO(albrow): Add options for filtering, sorting, limit, and offset.
func (db *DB) FindMiniHeaders(opts *FindMiniHeadersOpts) ([]*types.MiniHeader, error) {
	query, err := db.findMiniHeadersQueryFromOpts(opts)
	if err != nil {
		return nil, err
	}
	rawQuery, bindings := query.ToSQL(false)
	fmt.Println(rawQuery, bindings)
	var miniHeaders []*sqltypes.MiniHeader
	if err := query.GetAllContext(db.ctx, &miniHeaders); err != nil {
		return nil, err
	}
	return sqltypes.MiniHeadersToCommonType(miniHeaders), nil
}

// TODO(albrow): Can this be de-duped?
func (db *DB) findMiniHeadersQueryFromOpts(opts *FindMiniHeadersOpts) (*sqlz.SelectStmt, error) {
	query := sqlz.Newx(db.sqldb).Select("*").From("miniHeaders")
	if opts == nil {
		return query, nil
	}

	ordering := orderingFromMiniHeaderSortOpts(opts.Sort)
	if len(ordering) != 0 {
		query.OrderBy(ordering...)
	}
	if opts.Limit != 0 {
		query.Limit(int64(opts.Limit))
	}
	if opts.Offset != 0 {
		if opts.Limit == 0 {
			return nil, errors.New("db.FindMiniHeaders: can't use Offset without Limit")
		}
		query.Offset(int64(opts.Offset))
	}
	whereConditions, err := whereConditionsFromMiniHeaderFilterOpts(opts.Filters)
	if err != nil {
		return nil, err
	}
	if len(whereConditions) != 0 {
		query.Where(whereConditions...)
	}

	return query, nil
}

// TODO(albrow): Can this be de-duped?
func orderingFromMiniHeaderSortOpts(opts []MiniHeaderSort) []sqlz.SQLStmt {
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

// TODO(albrow): Can this be de-duped?
func whereConditionsFromMiniHeaderFilterOpts(opts []MiniHeaderFilter) ([]sqlz.WhereCondition, error) {
	// TODO(albrow): Type-check on value? You can't use CONTAINS with numeric types.
	whereConditions := make([]sqlz.WhereCondition, len(opts))
	for i, filterOpt := range opts {
		value := convertFilterValue(filterOpt.Value)
		switch filterOpt.Kind {
		case Equal:
			whereConditions[i] = sqlz.Eq(string(filterOpt.Field), value)
		case NotEqual:
			whereConditions[i] = sqlz.Not(sqlz.Eq(string(filterOpt.Field), value))
		case Less:
			whereConditions[i] = sqlz.Lt(string(filterOpt.Field), value)
		case Greater:
			whereConditions[i] = sqlz.Gt(string(filterOpt.Field), value)
		case LessOrEqual:
			whereConditions[i] = sqlz.Lte(string(filterOpt.Field), value)
		case GreaterOrEqual:
			whereConditions[i] = sqlz.Gte(string(filterOpt.Field), value)
		case Contains:
			// TODO(albrow): Value cannot contain special characters like "%".
			// TODO(albrow): Optimize this so it is easier to index.
			whereConditions[i] = sqlz.Like(string(filterOpt.Field), fmt.Sprintf("%%%s%%", value))
		default:
			return nil, fmt.Errorf("db.FindMiniHeaders: unknown FilterOpt.Kind: %s", filterOpt.Kind)
		}
	}
	return whereConditions, nil
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

func convertFilterValue(value interface{}) interface{} {
	switch v := value.(type) {
	case *big.Int:
		return sqltypes.NewBigInt(v)
	}
	return value
}
