// +build !js

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"sync"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/sqltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
	"github.com/ido50/sqlz"
	ds "github.com/ipfs/go-datastore"
	sqlds "github.com/ipfs/go-ds-sql"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

// largeLimit is used as a workaround due to the fact that SQL does not allow limit without offset.
const largeLimit = math.MaxInt64

var _ Database = (*DB)(nil)

// DB instantiates the DB connection and creates all the collections used by the application
type DB struct {
	ctx   context.Context
	sqldb *sqlz.DB
	opts  *Options

	// mu is used to protect all reads and writes to the database. This is a solution to the
	// `database is locked` error that appears on SQLite. https://github.com/mattn/go-sqlite3/issues/607
	// TODO(albrow): Make this mutex optional since not all SQL implementations need it.
	mu sync.RWMutex
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
		<-ctx.Done()
		_ = sqldb.Close()
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

func (db *DB) DHTStore() ds.Batching {
	return sqlds.NewDatastore(db.sqldb.DB.DB, NewSqliteQueriesForTable("dhtstore"))
}

func (db *DB) PeerStore() ds.Batching {
	return sqlds.NewDatastore(db.sqldb.DB.DB, NewSqliteQueriesForTable("peerstore"))
}

// TODO(albrow): Use a proper migration tool. We don't technically need this
// now but it will be necessary if we ever change the database schema.
// Note(albrow): If needed, we can optimize this by adding indexes to the
// orders and miniHeaders tables.
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
	isUnfillable             BOOLEAN NOT NULL,
	isExpired                BOOLEAN NOT NULL,
	parsedMakerAssetData     TEXT NOT NULL,
	parsedMakerFeeAssetData  TEXT NOT NULL,
	lastValidatedBlockNumber TEXT NOT NULL,
	lastValidatedBlockHash   TEXT NOT NULL,
	keepCancelled            BOOLEAN NOT NULL,
	keepExpired              BOOLEAN NOT NULL,
	keepFullyFilled          BOOLEAN NOT NULL,
	keepUnfunded             BOOLEAN NOT NULL
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
	ethRPCRequestsSentInCurrentUTCDay BIGINT NOT NULL,
	startOfCurrentUTCDay              DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS peerstore (
	key  TEXT NOT NULL UNIQUE,
	data BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS dhtstore (
	key  TEXT NOT NULL UNIQUE,
	data BYTEA NOT NULL
);
`

// Note(albrow): If needed, we can optimize this by using prepared
// statements for inserts instead of just a string.
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
	isUnfillable,
	isExpired,
	parsedMakerAssetData,
	parsedMakerFeeAssetData,
	lastValidatedBlockNumber,
	lastValidatedBlockHash,
	keepCancelled,
	keepExpired,
	keepFullyFilled,
	keepUnfunded
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
	:isUnfillable,
	:isExpired,
	:parsedMakerAssetData,
	:parsedMakerFeeAssetData,
	:lastValidatedBlockNumber,
	:lastValidatedBlockHash,
	:keepCancelled,
	:keepExpired,
	:keepFullyFilled,
	:keepUnfunded
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
	isUnfillable = :isUnfillable,
	isExpired = :isExpired,
	parsedMakerAssetData = :parsedMakerAssetData,
	parsedMakerFeeAssetData = :parsedMakerFeeAssetData,
	lastValidatedBlockNumber = :lastValidatedBlockNumber,
	lastValidatedBlockHash = :lastValidatedBlockHash,
	keepCancelled = :keepCancelled,
	keepExpired = :keepExpired,
	keepFullyFilled = :keepFullyFilled,
	keepUnfunded = :keepUnfunded
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
	ethRPCRequestsSentInCurrentUTCDay,
	startOfCurrentUTCDay
) VALUES (
	:ethereumChainID,
	:ethRPCRequestsSentInCurrentUTCDay,
	:startOfCurrentUTCDay
)`

const updateMetadataQuery = `UPDATE metadata SET
	ethereumChainID = :ethereumChainID,
	ethRPCRequestsSentInCurrentUTCDay = :ethRPCRequestsSentInCurrentUTCDay,
	startOfCurrentUTCDay = :startOfCurrentUTCDay
`

func (db *DB) migrate() error {
	_, err := db.sqldb.ExecContext(db.ctx, schema)
	return convertErr(err)
}

// ReadWriteTransactionalContext acquires a write lock, executes the transaction, then immediately releases the lock.
func (db *DB) ReadWriteTransactionalContext(ctx context.Context, opts *sql.TxOptions, f func(tx *sqlz.Tx) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.sqldb.TransactionalContext(ctx, opts, f)
}

func (db *DB) AddOrders(orders []*types.OrderWithMetadata) (alreadyStored []common.Hash, added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()

	sqlOrders := sqltypes.OrdersFromCommonType(orders)
	addedMap := map[common.Hash]*types.OrderWithMetadata{}
	sqlRemoved := []*sqltypes.Order{}

	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		for i, order := range sqlOrders {
			result, err := txn.NamedExecContext(db.ctx, insertOrderQuery, order)
			if err != nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if affected > 0 {
				addedMap[order.Hash] = orders[i]
			} else {
				alreadyStored = append(alreadyStored, order.Hash)
			}
		}

		// Remove orders with an expiration time too far in the future.
		// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
		// for DELETE statements. It also doesn't support RETURNING. As a
		// workaround, we do a SELECT and DELETE inside a transaction.
		// HACK(albrow): SQL doesn't support limit without offset. As a
		// workaround, we set the limit to an extremely large number.
		removeQuery := txn.Select("*").From("orders").
			OrderBy(sqlz.Desc(string(OFIsPinned)), sqlz.Asc(string(OFExpirationTimeSeconds))).
			Limit(largeLimit).
			Offset(int64(db.opts.MaxOrders))
		var ordersToRemove []*sqltypes.Order
		err = removeQuery.GetAllContext(db.ctx, &ordersToRemove)
		if err != nil {
			return err
		}
		for _, order := range ordersToRemove {
			_, err := txn.DeleteFrom("orders").Where(sqlz.Eq(string(OFHash), order.Hash)).ExecContext(db.ctx)
			if err != nil {
				return err
			}
			if _, found := addedMap[order.Hash]; found {
				// If the order was previously added, remove it from
				// the added set and don't add it to the removed set.
				delete(addedMap, order.Hash)
			} else {
				sqlRemoved = append(sqlRemoved, order)
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	for _, order := range addedMap {
		added = append(added, order)
	}
	return alreadyStored, added, sqltypes.OrdersToCommonType(sqlRemoved), nil
}

func (db *DB) GetOrder(hash common.Hash) (order *types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	var foundOrder sqltypes.Order
	db.mu.RLock()
	if err := db.sqldb.GetContext(db.ctx, &foundOrder, "SELECT * FROM orders WHERE hash = $1", hash); err != nil {
		db.mu.RUnlock()
		return nil, err
	}
	db.mu.RUnlock()
	return sqltypes.OrderToCommonType(&foundOrder), nil
}

func (db *DB) GetOrderStatuses(hashes []common.Hash) (statuses []*StoredOrderStatus, err error) {
	defer func() {
		err = convertErr(err)
	}()
	orderStatuses := make([]*StoredOrderStatus, len(hashes))
	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		for i, hash := range hashes {
			var foundOrder sqltypes.Order
			err := db.sqldb.GetContext(db.ctx, &foundOrder, "SELECT isRemoved, isUnfillable, fillableTakerAssetAmount FROM orders WHERE hash = $1", hash)
			if err != nil {
				if err == sql.ErrNoRows {
					orderStatuses[i] = &StoredOrderStatus{
						IsStored:                 false,
						IsMarkedRemoved:          false,
						IsMarkedUnfillable:       false,
						FillableTakerAssetAmount: nil,
					}
				} else {
					return err
				}
			} else {
				orderStatuses[i] = &StoredOrderStatus{
					IsStored:                 true,
					IsMarkedRemoved:          foundOrder.IsRemoved,
					IsMarkedUnfillable:       foundOrder.IsUnfillable,
					FillableTakerAssetAmount: foundOrder.FillableTakerAssetAmount.Int,
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return orderStatuses, nil
}

func (db *DB) FindOrders(query *OrderQuery) (orders []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQuery(query); err != nil {
		return nil, err
	}
	stmt, err := addOptsToSelectOrdersQuery(db.sqldb.Select("*").From("orders"), query)
	if err != nil {
		return nil, err
	}
	var foundOrders []*sqltypes.Order
	db.mu.RLock()
	err = stmt.GetAllContext(db.ctx, &foundOrders)
	db.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return sqltypes.OrdersToCommonType(foundOrders), nil
}

func (db *DB) CountOrders(query *OrderQuery) (count int, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQuery(query); err != nil {
		return 0, err
	}
	stmt, err := addOptsToSelectOrdersQuery(db.sqldb.Select("COUNT(*)").From("orders"), query)
	if err != nil {
		return 0, err
	}
	db.mu.RLock()
	gotCount, err := stmt.GetCountContext(db.ctx)
	db.mu.RUnlock()
	if err != nil {
		return 0, err
	}
	return int(gotCount), nil
}

type Selector interface {
	Select(cols ...string) *sqlz.SelectStmt
}

func addOptsToSelectOrdersQuery(stmt *sqlz.SelectStmt, opts *OrderQuery) (*sqlz.SelectStmt, error) {
	if opts == nil {
		return stmt, nil
	}

	ordering := orderingFromOrderSortOpts(opts.Sort)
	if len(ordering) != 0 {
		stmt.OrderBy(ordering...)
	}
	if opts.Limit != 0 {
		stmt.Limit(int64(opts.Limit))
	}
	if opts.Offset != 0 {
		stmt.Offset(int64(opts.Offset))
	}
	whereConditions, err := whereConditionsFromOrderFilterOpts(opts.Filters)
	if err != nil {
		return nil, err
	}
	if len(whereConditions) != 0 {
		stmt.Where(whereConditions...)
	}

	return stmt, nil
}

func orderingFromOrderSortOpts(sortOpts []OrderSort) []sqlz.SQLStmt {
	ordering := []sqlz.SQLStmt{}
	for _, sortOpt := range sortOpts {
		if sortOpt.Direction == Ascending {
			ordering = append(ordering, sqlz.Asc(string(sortOpt.Field)))
		} else {
			ordering = append(ordering, sqlz.Desc(string(sortOpt.Field)))
		}
	}
	return ordering
}

func whereConditionsFromOrderFilterOpts(filterOpts []OrderFilter) ([]sqlz.WhereCondition, error) {
	whereConditions := make([]sqlz.WhereCondition, len(filterOpts))
	for i, filterOpt := range filterOpts {
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
			// Note(albrow): If needed, we can optimize this so it is easier to index.
			// LIKE queries are notoriously slow.
			whereConditions[i] = sqlz.Like(string(filterOpt.Field), fmt.Sprintf("%%%s%%", value))
		default:
			return nil, fmt.Errorf("db.FindOrder: unknown FilterOpt.Kind: %s", filterOpt.Kind)
		}
	}
	return whereConditions, nil
}

func (db *DB) DeleteOrder(hash common.Hash) error {
	db.mu.Lock()
	_, err := db.sqldb.ExecContext(db.ctx, "DELETE FROM orders WHERE hash = $1", hash)
	db.mu.Unlock()
	if err != nil {
		return convertErr(err)
	}
	return nil
}

func (db *DB) DeleteOrders(query *OrderQuery) (deleted []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQuery(query); err != nil {
		return nil, err
	}
	// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
	// for DELETE statements. It also doesn't support RETURNING. As a
	// workaround, we do a SELECT and DELETE inside a transaction.
	var ordersToDelete []*sqltypes.Order
	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		stmt, err := addOptsToSelectOrdersQuery(txn.Select("*").From("orders"), query)
		if err != nil {
			return err
		}
		if err := stmt.GetAllContext(db.ctx, &ordersToDelete); err != nil {
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

func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) (err error) {
	defer func() {
		err = convertErr(err)
	}()
	if updateFunc == nil {
		return errors.New("db.UpdateOrders: updateFunc cannot be nil")
	}

	return db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
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
		return nil
	})
}

func (db *DB) AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()

	sqlMiniHeaders := sqltypes.MiniHeadersFromCommonType(miniHeaders)
	addedMap := map[common.Hash]*types.MiniHeader{}
	sqlRemoved := []*sqltypes.MiniHeader{}

	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		for i, miniHeader := range sqlMiniHeaders {
			result, err := txn.NamedExecContext(db.ctx, insertMiniHeaderQuery, miniHeader)
			if err != nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if affected > 0 {
				addedMap[miniHeader.Hash] = miniHeaders[i]
			}
		}

		// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
		// for DELETE statements. It also doesn't support RETURNING. As a
		// workaround, we do a SELECT and DELETE inside a transaction.
		// HACK(albrow): SQL doesn't support limit without offset. As a
		// workaround, we set the limit to an extremely large number.
		removeQuery := txn.Select("*").From("miniHeaders").OrderBy(sqlz.Desc(string(MFNumber))).Limit(largeLimit).Offset(int64(db.opts.MaxMiniHeaders))
		var miniHeadersToRemove []*sqltypes.MiniHeader
		if err := removeQuery.GetAllContext(db.ctx, &miniHeadersToRemove); err != nil {
			return err
		}
		for _, miniHeader := range miniHeadersToRemove {
			_, err := txn.DeleteFrom("miniHeaders").Where(sqlz.Eq(string(MFHash), miniHeader.Hash)).ExecContext(db.ctx)
			if err != nil {
				return err
			}
			if _, found := addedMap[miniHeader.Hash]; found {
				// If the miniHeader was previously added, remove it from
				// the added set and don't add it to the removed set.
				delete(addedMap, miniHeader.Hash)
			} else {
				sqlRemoved = append(sqlRemoved, miniHeader)
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	for _, miniHeader := range addedMap {
		added = append(added, miniHeader)
	}
	return added, sqltypes.MiniHeadersToCommonType(sqlRemoved), nil
}

// ResetMiniHeaders deletes all of the existing miniheaders and then stores new
// miniheaders in the database.
func (db *DB) ResetMiniHeaders(newMiniHeaders []*types.MiniHeader) (err error) {
	defer func() {
		err = convertErr(err)
	}()

	sqlNewMiniHeaders := sqltypes.MiniHeadersFromCommonType(newMiniHeaders)
	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		_, err := txn.DeleteFrom("miniHeaders").ExecContext(db.ctx)
		if err != nil {
			return err
		}
		for _, newMiniHeader := range sqlNewMiniHeaders {
			_, err := txn.NamedExecContext(db.ctx, insertMiniHeaderQuery, newMiniHeader)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (db *DB) GetMiniHeader(hash common.Hash) (miniHeader *types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()
	var sqlMiniHeader sqltypes.MiniHeader
	db.mu.RLock()
	err = db.sqldb.GetContext(db.ctx, &sqlMiniHeader, "SELECT * FROM miniHeaders WHERE hash = $1", hash)
	db.mu.RUnlock()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return sqltypes.MiniHeaderToCommonType(&sqlMiniHeader), nil
}

func (db *DB) FindMiniHeaders(query *MiniHeaderQuery) (miniHeaders []*types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()
	stmt, err := findMiniHeadersQueryFromOpts(db.sqldb, query)
	if err != nil {
		return nil, err
	}
	var sqlMiniHeaders []*sqltypes.MiniHeader
	db.mu.RLock()
	err = stmt.GetAllContext(db.ctx, &sqlMiniHeaders)
	db.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return sqltypes.MiniHeadersToCommonType(sqlMiniHeaders), nil
}

func findMiniHeadersQueryFromOpts(selector Selector, query *MiniHeaderQuery) (*sqlz.SelectStmt, error) {
	stmt := selector.Select("*").From("miniHeaders")
	if query == nil {
		return stmt, nil
	}

	ordering := orderingFromMiniHeaderSortOpts(query.Sort)
	if len(ordering) != 0 {
		stmt.OrderBy(ordering...)
	}
	if query.Limit != 0 {
		stmt.Limit(int64(query.Limit))
	}
	if query.Offset != 0 {
		if query.Limit == 0 {
			return nil, errors.New("db.FindMiniHeaders: can't use Offset without Limit")
		}
		stmt.Offset(int64(query.Offset))
	}
	whereConditions, err := whereConditionsFromMiniHeaderFilterOpts(query.Filters)
	if err != nil {
		return nil, err
	}
	if len(whereConditions) != 0 {
		stmt.Where(whereConditions...)
	}

	return stmt, nil
}

func orderingFromMiniHeaderSortOpts(sortOpts []MiniHeaderSort) []sqlz.SQLStmt {
	ordering := []sqlz.SQLStmt{}
	for _, sortOpt := range sortOpts {
		if sortOpt.Direction == Ascending {
			ordering = append(ordering, sqlz.Asc(string(sortOpt.Field)))
		} else {
			ordering = append(ordering, sqlz.Desc(string(sortOpt.Field)))
		}
	}
	return ordering
}

func whereConditionsFromMiniHeaderFilterOpts(filterOpts []MiniHeaderFilter) ([]sqlz.WhereCondition, error) {
	whereConditions := make([]sqlz.WhereCondition, len(filterOpts))
	for i, filterOpt := range filterOpts {
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
			// Note(albrow): If needed, we can optimize this so it is easier to index.
			// LIKE queries are notoriously slow.
			whereConditions[i] = sqlz.Like(string(filterOpt.Field), fmt.Sprintf("%%%s%%", value))
		default:
			return nil, fmt.Errorf("db.FindMiniHeaders: unknown FilterOpt.Kind: %s", filterOpt.Kind)
		}
	}
	return whereConditions, nil
}

func (db *DB) DeleteMiniHeader(hash common.Hash) error {
	db.mu.Lock()
	_, err := db.sqldb.ExecContext(db.ctx, "DELETE FROM miniHeaders WHERE hash = $1", hash)
	db.mu.Unlock()
	if err != nil {
		return convertErr(err)
	}
	return nil
}

func (db *DB) DeleteMiniHeaders(query *MiniHeaderQuery) (deleted []*types.MiniHeader, err error) {
	defer func() {
		err = convertErr(err)
	}()
	// HACK(albrow): sqlz doesn't support ORDER BY, LIMIT, and OFFSET
	// for DELETE statements. It also doesn't support RETURNING. As a
	// workaround, we do a SELECT and DELETE inside a transaction.
	var miniHeadersToDelete []*sqltypes.MiniHeader
	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(tx *sqlz.Tx) error {
		stmt, err := findMiniHeadersQueryFromOpts(tx, query)
		if err != nil {
			return err
		}
		if err := stmt.GetAllContext(db.ctx, &miniHeadersToDelete); err != nil {
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
	db.mu.RLock()
	err := db.sqldb.GetContext(db.ctx, &metadata, "SELECT * FROM metadata LIMIT 1")
	db.mu.RUnlock()
	if err != nil {
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
	sqlMetadata := sqltypes.MetadataFromCommonType(metadata)
	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		query := db.sqldb.Select("COUNT(*)").From("metadata")
		count, err := query.GetCount()
		if err != nil {
			return err
		}
		if count != 0 {
			return ErrMetadataAlreadyExists
		}
		_, err = db.sqldb.NamedExecContext(db.ctx, insertMetadataQuery, sqlMetadata)
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

	return db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
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

		return nil
	})
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
