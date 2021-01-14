// +build !js

package db

import (
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/sqltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ido50/sqlz"
)

func (db *DB) AddOrdersV4(orders []*types.OrderWithMetadata) (alreadyStored []common.Hash, added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()

	sqlOrders := sqltypes.OrdersFromCommonTypeV4(orders)
	addedMap := map[common.Hash]*types.OrderWithMetadata{}
	sqlRemoved := []*sqltypes.OrderV4{}

	err = db.ReadWriteTransactionalContext(db.ctx, nil, func(txn *sqlz.Tx) error {
		for i, order := range sqlOrders {
			result, err := txn.NamedExecContext(db.ctx, insertOrderQueryV4, order)
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
		removeQuery := txn.Select("*").From("ordersv4").
			OrderBy(sqlz.Desc(string(OV4FIsPinned)), sqlz.Asc(string(OV4FExpiry))).
			Limit(largeLimit).
			Offset(int64(db.opts.MaxOrders))
		var ordersToRemove []*sqltypes.OrderV4
		err = removeQuery.GetAllContext(db.ctx, &ordersToRemove)
		if err != nil {
			return err
		}

		for _, order := range ordersToRemove {
			_, err := txn.DeleteFrom("ordersv4").Where(sqlz.Eq(string(OV4FHash), order.Hash)).ExecContext(db.ctx)
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
	return alreadyStored, added, sqltypes.OrdersToCommonTypeV4(sqlRemoved), nil
}

func (db *DB) GetOrderV4(hash common.Hash) (order *types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	var foundOrder sqltypes.OrderV4
	db.mu.RLock()
	if err := db.sqldb.GetContext(db.ctx, &foundOrder, "SELECT * FROM ordersv4 WHERE hash = $1", hash); err != nil {
		db.mu.RUnlock()
		return nil, err
	}
	db.mu.RUnlock()
	return sqltypes.OrderToCommonTypeV4(&foundOrder), nil
}

func (db *DB) FindOrdersV4(query *OrderQueryV4) (orders []*types.OrderWithMetadata, err error) {
	defer func() {
		err = convertErr(err)
	}()
	if err := checkOrderQueryV4(query); err != nil {
		return nil, err
	}
	stmt, err := addOptsToSelectOrdersQueryV4(db.sqldb.Select("*").From("ordersv4"), query)
	if err != nil {
		return nil, err
	}
	var foundOrders []*sqltypes.OrderV4
	db.mu.RLock()
	err = stmt.GetAllContext(db.ctx, &foundOrders)
	db.mu.RUnlock()
	if err != nil {
		return nil, err
	}
	return sqltypes.OrdersToCommonTypeV4(foundOrders), nil
}

func addOptsToSelectOrdersQueryV4(stmt *sqlz.SelectStmt, opts *OrderQueryV4) (*sqlz.SelectStmt, error) {
	if opts == nil {
		return stmt, nil
	}

	ordering := orderingFromOrderSortOptsV4(opts.Sort)
	if len(ordering) != 0 {
		stmt.OrderBy(ordering...)
	}
	if opts.Limit != 0 {
		stmt.Limit(int64(opts.Limit))
	}
	if opts.Offset != 0 {
		stmt.Offset(int64(opts.Offset))
	}
	whereConditions, err := whereConditionsFromOrderFilterOptsV4(opts.Filters)
	if err != nil {
		return nil, err
	}
	if len(whereConditions) != 0 {
		stmt.Where(whereConditions...)
	}

	return stmt, nil
}

func orderingFromOrderSortOptsV4(sortOpts []OrderSortV4) []sqlz.SQLStmt {
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

func whereConditionsFromOrderFilterOptsV4(filterOpts []OrderFilterV4) ([]sqlz.WhereCondition, error) {
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

func checkOrderQueryV4(query *OrderQueryV4) error {
	if query == nil {
		return nil
	}
	if query.Offset != 0 && query.Limit == 0 {
		return errors.New("can't use Offset without Limit")
	}
	return nil

}
