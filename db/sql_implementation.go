// +build !js

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/sqltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
	"github.com/ido50/sqlz"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var _ Database = (*DB)(nil)

// DB instantiates the DB connection and creates all the collections used by the application
type DB struct {
	ctx   context.Context
	sqldb *sqlz.DB
	opts  *Options
}

func defaultOptions() *Options {
	return &Options{
		DriverName:     "sqlite3",
		DataSourceName: "0x_mesh/db/db.sqlite",
		MaxOrders:      100000,
		MaxMiniHeaders: 20,
	}
}

// TestOptions returns a set of options suitable for testing.
func TestOptions() *Options {
	return &Options{
		DriverName:     "sqlite3",
		DataSourceName: filepath.Join("/tmp", "mesh_testing", uuid.New().String(), "db.sqlite"),
		MaxOrders:      100,
		MaxMiniHeaders: 20,
	}
}

// New creates a new connection to the database. The connection will be automatically closed
// when the given context is canceled.
func New(ctx context.Context, opts *Options) (*DB, error) {
	opts = parseOptions(opts)

	connectCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := os.MkdirAll(filepath.Dir(opts.DataSourceName), os.ModePerm); err != nil && err != os.ErrExist {
		return nil, err
	}

	sqldb, err := sqlx.ConnectContext(connectCtx, opts.DriverName, opts.DataSourceName)
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
		sqldb: sqlz.Newx(sqldb),
		opts:  opts,
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
	chainID                  TEXT NOT NULL,
	exchangeAddress          TEXT NOT NULL,
	makerAddress             TEXT NOT NULL,
	makerAssetData           TEXT NOT NULL,
	makerFeeAssetData        TEXT NOT NULL,
	makerAssetAmount         TEXT NOT NULL,
	makerFee                 TEXT NOT NULL,
	takerAddress             TEXT NOT NULL,
	takerAssetData           TEXT NOT NULL,
	takerFeeAssetData        TEXT NOT NULL,
	takerAssetAmount         TEXT NOT NULL,
	takerFee                 TEXT NOT NULL,
	senderAddress            TEXT NOT NULL,
	feeRecipientAddress      TEXT NOT NULL,
	expirationTimeSeconds    TEXT NOT NULL,
	salt                     TEXT NOT NULL,
	signature                TEXT NOT NULL,
	lastUpdated              DATETIME NOT NULL,
	fillableTakerAssetAmount TEXT NOT NULL,
	isRemoved                BOOLEAN NOT NULL,
	isPinned                 BOOLEAN NOT NULL,
	parsedMakerAssetData     TEXT NOT NULL,
	parsedMakerFeeAssetData  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS miniHeaders (
	hash      TEXT UNIQUE NOT NULL,
	number    TEXT UNIQUE NOT NULL,
	parent    TEXT NOT NULL,
	timestamp DATETIME NOT NULL,
	logs      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata (
	ethereumChainID                   BIGINT NOT NULL,
	maxExpirationTime                 TEXT NOT NULL,
	ethRPCRequestsSentInCurrentUTCDay BIGINT NOT NULL,
	startOfCurrentUTCDay              DATETIME NOT NULL
);
`

// TODO(albrow): Used prepared statement for inserts.
const insertOrderQuery = `INSERT INTO orders (
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
) ON CONFLICT DO NOTHING
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

const insertMiniHeaderQuery = `INSERT INTO miniHeaders (
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
) ON CONFLICT DO NOTHING`

const insertMetadataQuery = `INSERT INTO metadata (
	ethereumChainID,
	maxExpirationTime,
	ethRPCRequestsSentInCurrentUTCDay,
	startOfCurrentUTCDay
) VALUES (
	:ethereumChainID,
	:maxExpirationTime,
	:ethRPCRequestsSentInCurrentUTCDay,
	:startOfCurrentUTCDay
)`

const updateMetadataQuery = `UPDATE metadata SET
	ethereumChainID = :ethereumChainID,
	maxExpirationTime = :maxExpirationTime,
	ethRPCRequestsSentInCurrentUTCDay = :ethRPCRequestsSentInCurrentUTCDay,
	startOfCurrentUTCDay = :startOfCurrentUTCDay
`

func (db *DB) migrate() error {
	_, err := db.sqldb.ExecContext(db.ctx, schema)
	return convertErr(err)
}

func (db *DB) AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
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

func (db *DB) GetOrder(hash common.Hash) (order *types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	var foundOrder sqltypes.Order
	if err := db.sqldb.GetContext(db.ctx, &foundOrder, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		return nil, err
	}
	return sqltypes.OrderToCommonType(&foundOrder), nil
}

func (db *DB) FindOrders(opts *OrderQuery) (orders []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQuery(opts); err != nil {
		return nil, err
	}
	query, err := addOptsToSelectOrdersQuery(db.sqldb.Select("*").From("orders"), opts)
	if err != nil {
		return nil, err
	}
	var foundOrders []*sqltypes.Order
	if err := query.GetAllContext(db.ctx, &foundOrders); err != nil {
		return nil, err
	}
	return sqltypes.OrdersToCommonType(foundOrders), nil
}

func (db *DB) CountOrders(opts *OrderQuery) (count int, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQuery(opts); err != nil {
		return 0, err
	}
	query, err := addOptsToSelectOrdersQuery(db.sqldb.Select("COUNT(*)").From("orders"), opts)
	if err != nil {
		return 0, err
	}
	gotCount, err := query.GetCount()
	if err != nil {
		return 0, err
	}
	return int(gotCount), nil
}

type Selector interface {
	Select(cols ...string) *sqlz.SelectStmt
}

func addOptsToSelectOrdersQuery(query *sqlz.SelectStmt, opts *OrderQuery) (*sqlz.SelectStmt, error) {
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
		return convertErr(err)
	}
	return nil
}

// TODO(albrow): Test deleting with ORDER BY, LIMIT, and OFFSET.
func (db *DB) DeleteOrders(opts *OrderQuery) (deleted []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQuery(opts); err != nil {
		return nil, err
	}
	// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
	// for DELETE statements. It also doesn't support RETURNING. As a
	// workaround, we do a SELECT and DELETE inside a transaction.
	var ordersToDelete []*sqltypes.Order
	err = db.sqldb.TransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		query, err := addOptsToSelectOrdersQuery(txn.Select("*").From("orders"), opts)
		if err != nil {
			return err
		}
		if err := query.GetAllContext(db.ctx, &ordersToDelete); err != nil {
			return err
		}
		for _, order := range ordersToDelete {
			_, err := txn.DeleteFrom("orders").Where(sqlz.Eq(string(OFHash), order.Hash)).ExecContext(db.ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sqltypes.OrdersToCommonType(ordersToDelete), nil
}

// TODO(albrow): Consider automatically setting LastUpdated?
func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) (err error) {
	defer func() {
		err = convertErr(err)
	}()
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
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
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
	defer func() {
		err = convertErr(err)
	}()
	var miniHeadersToRemove []*sqltypes.MiniHeader
	err = db.sqldb.TransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		for _, miniHeader := range miniHeaders {
			result, err := txn.NamedExecContext(db.ctx, insertMiniHeaderQuery, sqltypes.MiniHeaderFromCommonType(miniHeader))
			if err != nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if affected > 0 {
				added = append(added, miniHeader)
			}
		}

		// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
		// for DELETE statements. It also doesn't support RETURNING. As a
		// workaround, we do a SELECT and DELETE inside a transaction.
		trimQuery := txn.Select("*").From("miniHeaders").OrderBy(sqlz.Desc(string(MFNumber))).Limit(99999999999).Offset(int64(db.opts.MaxMiniHeaders))
		if err := trimQuery.GetAllContext(db.ctx, &miniHeadersToRemove); err != nil {
			return err
		}
		for _, miniHeader := range miniHeadersToRemove {
			_, err := txn.DeleteFrom("miniHeaders").Where(sqlz.Eq(string(MFHash), miniHeader.Hash)).ExecContext(db.ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// TODO(albrow): Because of how the above code is written, a single
	// miniHeader could exist in both added and removed sets. Should we
	// remove such miniHeaders from both sets in this case?
	return added, sqltypes.MiniHeadersToCommonType(miniHeadersToRemove), nil
}

func (db *DB) GetMiniHeader(hash common.Hash) (miniHeader *types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()
	var sqlMiniHeader sqltypes.MiniHeader
	if err := db.sqldb.GetContext(db.ctx, &sqlMiniHeader, "SELECT * FROM miniHeaders WHERE hash = $1", hash); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return sqltypes.MiniHeaderToCommonType(&sqlMiniHeader), nil
}

func (db *DB) FindMiniHeaders(opts *MiniHeaderQuery) (miniHeaders []*types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()
	query, err := findMiniHeadersQueryFromOpts(db.sqldb, opts)
	if err != nil {
		return nil, err
	}
	var sqlMiniHeaders []*sqltypes.MiniHeader
	if err := query.GetAllContext(db.ctx, &sqlMiniHeaders); err != nil {
		return nil, err
	}
	return sqltypes.MiniHeadersToCommonType(sqlMiniHeaders), nil
}

// TODO(albrow): Can this be de-duped?
func findMiniHeadersQueryFromOpts(selector Selector, opts *MiniHeaderQuery) (*sqlz.SelectStmt, error) {
	query := selector.Select("*").From("miniHeaders")
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

func (db *DB) DeleteMiniHeader(hash common.Hash) error {
	if _, err := db.sqldb.ExecContext(db.ctx, "DELETE FROM miniHeaders WHERE hash = $1", hash); err != nil {
		return convertErr(err)
	}
	return nil
}

// TODO(albrow): Test deleting with ORDER BY, LIMIT, and OFFSET.
func (db *DB) DeleteMiniHeaders(opts *MiniHeaderQuery) (deleted []*types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()
	// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
	// for DELETE statements. It also doesn't support RETURNING. As a
	// workaround, we do a SELECT and DELETE inside a transaction.
	var miniHeadersToDelete []*sqltypes.MiniHeader
	err = db.sqldb.TransactionalContext(db.ctx, nil, func(tx *sqlz.Tx) error {
		query, err := findMiniHeadersQueryFromOpts(tx, opts)
		if err != nil {
			return err
		}
		if err := query.GetAllContext(db.ctx, &miniHeadersToDelete); err != nil {
			return err
		}
		for _, miniHeader := range miniHeadersToDelete {
			_, err := tx.DeleteFrom("miniHeaders").Where(sqlz.Eq(string(MFHash), miniHeader.Hash)).ExecContext(db.ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sqltypes.MiniHeadersToCommonType(miniHeadersToDelete), nil
}

// GetMetadata returns the metadata (or db.ErrNotFound if no metadata has been saved).
func (db *DB) GetMetadata() (*types.Metadata, error) {
	var metadata sqltypes.Metadata
	if err := db.sqldb.GetContext(db.ctx, &metadata, "SELECT * FROM metadata LIMIT 1"); err != nil {
		return nil, convertErr(err)
	}
	return sqltypes.MetadataToCommonType(&metadata), nil
}

// SaveMetadata inserts the metadata into the database, overwriting any existing
// metadata. It returns ErrMetadataAlreadyExists if the metadata has already been
// saved in the database.
func (db *DB) SaveMetadata(metadata *types.Metadata) (err error) {
	defer func() {
		err = convertErr(err)
	}()
	err = db.sqldb.TransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		query := db.sqldb.Select("COUNT(*)").From("metadata")
		count, err := query.GetCount()
		if err != nil {
			return err
		}
		if count != 0 {
			return ErrMetadataAlreadyExists
		}
		_, err = db.sqldb.NamedExecContext(db.ctx, insertMetadataQuery, sqltypes.MetadataFromCommonType(metadata))
		return err
	})
	return err
}

// UpdateMetadata updates the metadata in the database via a transaction. It
// accepts a callback function which will be provided with the old metadata and
// should return the new metadata to save.
func (db *DB) UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) (err error) {
	defer func() {
		err = convertErr(err)
	}()
	if updateFunc == nil {
		return errors.New("db.UpdateMetadata: updateFunc cannot be nil")
	}

	txn, err := db.sqldb.BeginTxx(db.ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	var existingMetadata sqltypes.Metadata
	if err := txn.GetContext(db.ctx, &existingMetadata, "SELECT * FROM metadata LIMIT 1"); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	commonMetadata := sqltypes.MetadataToCommonType(&existingMetadata)
	commonUpdatedMetadata := updateFunc(commonMetadata)
	updatedMetadata := sqltypes.MetadataFromCommonType(commonUpdatedMetadata)
	_, err = txn.NamedExecContext(db.ctx, updateMetadataQuery, updatedMetadata)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func convertFilterValue(value interface{}) interface{} {
	switch v := value.(type) {
	case *big.Int:
		return sqltypes.NewSortedBigInt(v)
	}
	return value
}

func assetDataIncludesTokenAddressAndTokenID(field OrderField, tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	filterValueJSON, err := canonicaljson.Marshal(sqltypes.SingleAssetData{
		Address: tokenAddress,
		TokenID: sqltypes.NewBigInt(tokenID),
	})
	if err != nil {
		// big.Int and common.Address types should never return an error when marshaling to JSON
		panic(err)
	}
	fmt.Println(string(filterValueJSON))
	return OrderFilter{
		Field: field,
		Kind:  Contains,
		Value: string(filterValueJSON),
	}
}

// convertErr converts from SQL specific errors to common error types.
func convertErr(err error) error {
	if err == nil {
		return nil
	}
	// Check if the error matches known exported values.
	switch err {
	case sql.ErrNoRows:
		return ErrNotFound
	case sql.ErrConnDone:
		return ErrClosed
	}
	// As a fallback, check the error string for errors which have no
	// exported type/value.
	switch err.Error() {
	case "sql: database is closed":
		return ErrClosed
	}
	return err
}
