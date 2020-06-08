// +build js,wasm

package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"syscall/js"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/dexietypes"
	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
)

var _ Database = (*DB)(nil)

type DB struct {
	ctx   context.Context
	dexie js.Value
	opts  *Options
}

func TestOptions() *Options {
	return &Options{
		DriverName:     "dexie",
		DataSourceName: filepath.Join("mesh_testing", uuid.New().String()),
		MaxOrders:      100,
		MaxMiniHeaders: 20,
	}
}

func defaultOptions() *Options {
	return &Options{
		DriverName:     "dexie",
		DataSourceName: "mesh_dexie_database",
		MaxOrders:      100000,
		MaxMiniHeaders: 20,
	}
}

// New creates a new connection to the database. The connection will be automatically closed
// when the given context is canceled.
func New(ctx context.Context, opts *Options) (database *DB, err error) {
	if opts != nil && opts.DriverName != "dexie" {
		return nil, fmt.Errorf(`unexpected driver name for js/wasm: %q (only "dexie" is supported)`, opts.DriverName)
	}
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	dexieClass := js.Global().Get("__mesh_dexie_newDatabase__")
	if jsutil.IsNullOrUndefined(dexieClass) {
		return nil, errors.New("could not detect Dexie.js")
	}
	opts = parseOptions(opts)
	dexie := dexieClass.Invoke(opts)
	return &DB{
		ctx:   ctx,
		dexie: dexie,
		opts:  opts,
	}, nil
}

func (db *DB) AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsOrders, err := jsutil.InefficientlyConvertToJS(dexietypes.OrdersFromCommonType(orders))
	if err != nil {
		return nil, nil, err
	}
	jsResult, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("addOrdersAsync", jsOrders))
	if err != nil {
		return nil, nil, convertJSError(err)
	}
	jsAdded := jsResult.Get("added")
	var dexieAdded []*dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsAdded, &dexieAdded); err != nil {
		return nil, nil, err
	}
	jsRemoved := jsResult.Get("removed")
	var dexieRemoved []*dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsRemoved, &dexieRemoved); err != nil {
		return nil, nil, err
	}
	return dexietypes.OrdersToCommonType(dexieAdded), dexietypes.OrdersToCommonType(dexieRemoved), nil
}

func (db *DB) GetOrder(hash common.Hash) (order *types.OrderWithMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsOrder, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("getOrderAsync", hash.Hex()))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieOrder dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsOrder, &dexieOrder); err != nil {
		return nil, err
	}
	return dexietypes.OrderToCommonType(&dexieOrder), nil
}

func (db *DB) FindOrders(query *OrderQuery) (orders []*types.OrderWithMetadata, err error) {
	if err := checkOrderQuery(query); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	query = formatOrderQuery(query)
	jsOrders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("findOrdersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieOrders []*dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsOrders, &dexieOrders); err != nil {
		return nil, err
	}
	return dexietypes.OrdersToCommonType(dexieOrders), nil
}

func (db *DB) CountOrders(query *OrderQuery) (count int, err error) {
	if err := checkOrderQuery(query); err != nil {
		return 0, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	query = formatOrderQuery(query)
	jsCount, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("countOrdersAsync", query))
	if err != nil {
		return 0, convertJSError(err)
	}
	return jsCount.Int(), nil
}

func (db *DB) DeleteOrder(hash common.Hash) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	_, jsErr := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("deleteOrderAsync", hash.Hex()))
	if jsErr != nil {
		return convertJSError(jsErr)
	}
	return nil
}

func (db *DB) DeleteOrders(query *OrderQuery) (deletedOrders []*types.OrderWithMetadata, err error) {
	if err := checkOrderQuery(query); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	query = formatOrderQuery(query)
	jsOrders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("deleteOrdersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieOrders []*dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsOrders, &dexieOrders); err != nil {
		return nil, err
	}
	return dexietypes.OrdersToCommonType(dexieOrders), nil
}

func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsUpdateFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		jsExistingOrder := args[0]
		var dexieExistingOrder dexietypes.Order
		if err := jsutil.InefficientlyConvertFromJS(jsExistingOrder, &dexieExistingOrder); err != nil {
			panic(err)
		}
		orderToUpdate, err := updateFunc(dexietypes.OrderToCommonType(&dexieExistingOrder))
		if err != nil {
			panic(err)
		}
		dexieOrderToUpdate := dexietypes.OrderFromCommonType(orderToUpdate)
		jsOrderToUpdate, err := jsutil.InefficientlyConvertToJS(dexieOrderToUpdate)
		if err != nil {
			panic(err)
		}
		return jsOrderToUpdate
	})
	defer jsUpdateFunc.Release()
	_, jsErr := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("updateOrderAsync", hash.Hex(), jsUpdateFunc))
	if jsErr != nil {
		return convertJSError(jsErr)
	}
	return nil
}

func (db *DB) AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsMiniHeaders, err := jsutil.InefficientlyConvertToJS(dexietypes.MiniHeadersFromCommonType(miniHeaders))
	if err != nil {
		return nil, nil, err
	}
	jsResult, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("addMiniHeadersAsync", jsMiniHeaders))
	if err != nil {
		return nil, nil, convertJSError(err)
	}
	jsAdded := jsResult.Get("added")
	var dexieAdded []*dexietypes.MiniHeader
	if err := jsutil.InefficientlyConvertFromJS(jsAdded, &dexieAdded); err != nil {
		return nil, nil, err
	}
	jsRemoved := jsResult.Get("removed")
	var dexieRemoved []*dexietypes.MiniHeader
	if err := jsutil.InefficientlyConvertFromJS(jsRemoved, &dexieRemoved); err != nil {
		return nil, nil, err
	}
	return dexietypes.MiniHeadersToCommonType(dexieAdded), dexietypes.MiniHeadersToCommonType(dexieRemoved), nil
}

func (db *DB) GetMiniHeader(hash common.Hash) (miniHeader *types.MiniHeader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsMiniHeader, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("getMiniHeaderAsync", hash.Hex()))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieMiniHeader dexietypes.MiniHeader
	if err := jsutil.InefficientlyConvertFromJS(jsMiniHeader, &dexieMiniHeader); err != nil {
		return nil, err
	}
	return dexietypes.MiniHeaderToCommonType(&dexieMiniHeader), nil
}

func (db *DB) FindMiniHeaders(query *MiniHeaderQuery) (miniHeaders []*types.MiniHeader, err error) {
	if err := checkMiniHeaderQuery(query); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	query = formatMiniHeaderQuery(query)
	jsMiniHeaders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("findMiniHeadersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieMiniHeaders []*dexietypes.MiniHeader
	if err := jsutil.InefficientlyConvertFromJS(jsMiniHeaders, &dexieMiniHeaders); err != nil {
		return nil, err
	}
	return dexietypes.MiniHeadersToCommonType(dexieMiniHeaders), nil
}

func (db *DB) DeleteMiniHeader(hash common.Hash) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	_, jsErr := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("deleteMiniHeaderAsync", hash.Hex()))
	if jsErr != nil {
		return convertJSError(jsErr)
	}
	return nil
}

func (db *DB) DeleteMiniHeaders(query *MiniHeaderQuery) (deleted []*types.MiniHeader, err error) {
	if err := checkMiniHeaderQuery(query); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	query = formatMiniHeaderQuery(query)
	jsMiniHeaders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("deleteMiniHeadersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieMiniHeaders []*dexietypes.MiniHeader
	if err := jsutil.InefficientlyConvertFromJS(jsMiniHeaders, &dexieMiniHeaders); err != nil {
		return nil, err
	}
	return dexietypes.MiniHeadersToCommonType(dexieMiniHeaders), nil
}

func (db *DB) GetMetadata() (metadata *types.Metadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsMetadata, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("getMetadataAsync"))
	if err != nil {
		return nil, convertJSError(err)
	}
	var dexieMetadata dexietypes.Metadata
	if err := jsutil.InefficientlyConvertFromJS(jsMetadata, &dexieMetadata); err != nil {
		return nil, err
	}
	return dexietypes.MetadataToCommonType(&dexieMetadata), nil
}

func (db *DB) SaveMetadata(metadata *types.Metadata) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	dexieMetadata := dexietypes.MetadataFromCommonType(metadata)
	jsMetadata, err := jsutil.InefficientlyConvertToJS(dexieMetadata)
	if err != nil {
		return err
	}
	_, err = jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("saveMetadataAsync", jsMetadata))
	if err != nil {
		return convertJSError(err)
	}
	return nil
}

func (db *DB) UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsUpdateFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		jsExistingMetadata := args[0]
		var dexieExistingMetadata dexietypes.Metadata
		if err := jsutil.InefficientlyConvertFromJS(jsExistingMetadata, &dexieExistingMetadata); err != nil {
			panic(err)
		}
		metadataToUpdate := updateFunc(dexietypes.MetadataToCommonType(&dexieExistingMetadata))
		dexieMetadataToUpdate := dexietypes.MetadataFromCommonType(metadataToUpdate)
		jsMetadataToUpdate, err := jsutil.InefficientlyConvertToJS(dexieMetadataToUpdate)
		if err != nil {
			panic(err)
		}
		return jsMetadataToUpdate
	})
	defer jsUpdateFunc.Release()
	_, jsErr := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("updateMetadataAsync", jsUpdateFunc))
	if jsErr != nil {
		return convertJSError(jsErr)
	}
	return nil
}

func recoverError(e interface{}) error {
	switch e := e.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	default:
		return fmt.Errorf("unexpected JavaScript error: (%T) %v", e, e)
	}
}

func convertJSError(e error) error {
	switch e := e.(type) {
	case js.Error:
		if jsutil.IsNullOrUndefined(e.Value) {
			return e
		}
		if jsutil.IsNullOrUndefined(e.Value.Get("message")) {
			return e
		}
		switch e.Value.Get("message").String() {
		// TOOD(albrow): Handle more error messages here
		case ErrNotFound.Error():
			return ErrNotFound
		case ErrMetadataAlreadyExists.Error():
			return ErrMetadataAlreadyExists
		case ErrDBFilledWithPinnedOrders.Error():
			return ErrDBFilledWithPinnedOrders
		}
	}
	return e
}

func formatOrderQuery(query *OrderQuery) *OrderQuery {
	if query == nil {
		return nil
	}
	for i, filter := range query.Filters {
		query.Filters[i].Value = convertFilterValue(filter.Value)
	}
	return query
}

func formatMiniHeaderQuery(query *MiniHeaderQuery) *MiniHeaderQuery {
	if query == nil {
		return nil
	}
	for i, filter := range query.Filters {
		query.Filters[i].Value = convertFilterValue(filter.Value)
	}
	return query
}

func convertFilterValue(value interface{}) interface{} {
	switch v := value.(type) {
	case *big.Int:
		return dexietypes.NewSortedBigInt(v)
	case bool:
		return dexietypes.BoolToUint8(v)
	}
	return value
}

func assetDataIncludesTokenAddressAndTokenID(field OrderField, tokenAddress common.Address, tokenID *big.Int) OrderFilter {
	filterValueJSON, err := canonicaljson.Marshal(dexietypes.SingleAssetData{
		Address: tokenAddress,
		TokenID: dexietypes.NewBigInt(tokenID),
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
