// +build js,wasm

package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db/dexietypes"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
	ds "github.com/ipfs/go-datastore"
	"github.com/sirupsen/logrus"
)

const (
	// slowQueryDebugDuration is the minimum duration used to determine whether to log slow queries.
	// Any query that takes longer than this will be logged at the Debug level.
	slowQueryDebugDuration = 1 * time.Second
	// slowQueryWarnDuration is the minimum duration used to determine whether to log slow queries.
	// Any query that takes longer than this will be logged at the Warning level.
	slowQueryWarnDuration = 5 * time.Second
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
			err = jsutil.RecoverError(r)
		}
	}()
	newDexieDatabase := js.Global().Get("__mesh_dexie_newDatabase__")
	if jsutil.IsNullOrUndefined(newDexieDatabase) {
		return nil, errors.New("could not detect Dexie.js")
	}
	opts = parseOptions(opts)
	dexie := newDexieDatabase.Invoke(opts)

	// Automatically close the database connection when the context is canceled.
	go func() {
		select {
		case <-ctx.Done():
			_ = dexie.Call("close")
		}
	}()

	return &DB{
		ctx:   ctx,
		dexie: dexie,
		opts:  opts,
	}, nil
}

func (db *DB) PeerStore() ds.Batching {
	dexieStore := db.dexie.Call("peerStore")
	return &Datastore{
		ctx:        db.ctx,
		dexieStore: dexieStore,
	}
}

func (db *DB) DHTStore() ds.Batching {
	dexieStore := db.dexie.Call("dhtStore")
	return &Datastore{
		ctx:        db.ctx,
		dexieStore: dexieStore,
	}
}

func (db *DB) AddOrders(orders []*types.OrderWithMetadata) (alreadyStored []common.Hash, added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("AddOrders with %d orders", len(orders)))
	jsOrders, err := jsutil.InefficientlyConvertToJS(dexietypes.OrdersFromCommonType(orders))
	if err != nil {
		return nil, nil, nil, err
	}
	jsResult, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("addOrdersAsync", jsOrders))
	if err != nil {
		return nil, nil, nil, convertJSError(err)
	}
	jsAlreadyStored := jsResult.Get("alreadyStored")
	if !jsutil.IsNullOrUndefined(jsAlreadyStored) {
		for i := 0; i < jsAlreadyStored.Length(); i++ {
			hashString := jsAlreadyStored.Index(i).String()
			alreadyStored = append(alreadyStored, common.HexToHash(hashString))
		}
	}
	jsAdded := jsResult.Get("added")
	var dexieAdded []*dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsAdded, &dexieAdded); err != nil {
		return nil, nil, nil, err
	}
	jsRemoved := jsResult.Get("removed")
	var dexieRemoved []*dexietypes.Order
	if err := jsutil.InefficientlyConvertFromJS(jsRemoved, &dexieRemoved); err != nil {
		return nil, nil, nil, err
	}
	return alreadyStored, dexietypes.OrdersToCommonType(dexieAdded), dexietypes.OrdersToCommonType(dexieRemoved), nil
}

func (db *DB) GetOrder(hash common.Hash) (order *types.OrderWithMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "GetOrder")
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

func (db *DB) GetOrderStatuses(hashes []common.Hash) (statuses []*StoredOrderStatus, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("GetOrderStatuses with %d hashes", len(hashes)))
	stringHashes := make([]interface{}, len(hashes))
	for i, hash := range hashes {
		stringHashes[i] = hash.Hex()
	}
	jsResults, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("getOrderStatusesAsync", stringHashes))
	if err != nil {
		return nil, convertJSError(err)
	}
	statuses = make([]*StoredOrderStatus, jsResults.Length())
	for i := 0; i < len(statuses); i++ {
		jsStatus := jsResults.Index(i)
		var fillableAmount *big.Int
		jsAmount := jsStatus.Get("fillableTakerAssetAmount")
		if !jsutil.IsNullOrUndefined(jsAmount) {
			fillableAmount, _ = big.NewInt(0).SetString(jsAmount.String(), 10)
		}
		statuses[i] = &StoredOrderStatus{
			IsStored:                 jsStatus.Get("isStored").Bool(),
			IsMarkedRemoved:          jsStatus.Get("isMarkedRemoved").Bool(),
			FillableTakerAssetAmount: fillableAmount,
		}
	}
	return statuses, nil
}

func (db *DB) FindOrders(query *OrderQuery) (orders []*types.OrderWithMetadata, err error) {
	if err := checkOrderQuery(query); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("FindOrders %s", spew.Sdump(query)))
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("CountOrders %s", spew.Sdump(query)))
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "DeleteOrder")
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("DeleteOrders %s", spew.Sdump(query)))
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "UpdateOrder")
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("AddMiniHeaders with %d miniHeaders", len(miniHeaders)))
	jsMiniHeaders := dexietypes.MiniHeadersFromCommonType(miniHeaders)
	jsResult, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("addMiniHeadersAsync", jsMiniHeaders))
	if err != nil {
		return nil, nil, convertJSError(err)
	}
	jsAdded := jsResult.Get("added")
	jsRemoved := jsResult.Get("removed")
	return dexietypes.MiniHeadersToCommonType(jsAdded), dexietypes.MiniHeadersToCommonType(jsRemoved), nil
}

// ResetMiniHeaders deletes all of the existing miniheaders and then stores new
// miniheaders in the database.
func (db *DB) ResetMiniHeaders(newMiniHeaders []*types.MiniHeader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("ResetMiniHeaders with %d newMiniHeaders", len(newMiniHeaders)))
	jsNewMiniHeaders := dexietypes.MiniHeadersFromCommonType(newMiniHeaders)
	_, err = jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("resetMiniHeadersAsync", jsNewMiniHeaders))
	if err != nil {
		return convertJSError(err)
	}
	return nil
}

func (db *DB) GetMiniHeader(hash common.Hash) (miniHeader *types.MiniHeader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "GetMiniHeader")
	jsMiniHeader, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("getMiniHeaderAsync", hash.Hex()))
	if err != nil {
		return nil, convertJSError(err)
	}
	return dexietypes.MiniHeaderToCommonType(jsMiniHeader), nil
}

func (db *DB) FindMiniHeaders(query *MiniHeaderQuery) (miniHeaders []*types.MiniHeader, err error) {
	if err := checkMiniHeaderQuery(query); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("FindMiniHeaders %s", spew.Sdump(query)))
	query = formatMiniHeaderQuery(query)
	jsMiniHeaders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("findMiniHeadersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	return dexietypes.MiniHeadersToCommonType(jsMiniHeaders), nil
}

func (db *DB) DeleteMiniHeader(hash common.Hash) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "DeleteMiniHeader")
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, fmt.Sprintf("DeleteMiniHeaders %s", spew.Sdump(query)))
	query = formatMiniHeaderQuery(query)
	jsMiniHeaders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("deleteMiniHeadersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	return dexietypes.MiniHeadersToCommonType(jsMiniHeaders), nil
}

func (db *DB) GetMetadata() (metadata *types.Metadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "GetMetadata")
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "SaveMetadata")
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
			err = jsutil.RecoverError(r)
		}
	}()
	start := time.Now()
	defer logQueryIfSlow(start, "UpdateMetadata")
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
		case ds.ErrNotFound.Error():
			return ds.ErrNotFound
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

func logQueryIfSlow(start time.Time, msg string) {
	duration := time.Since(start)
	if duration > slowQueryDebugDuration {
		logWithFields := logrus.WithFields(logrus.Fields{
			"message":  msg,
			"duration": fmt.Sprint(duration),
		})
		if duration > slowQueryWarnDuration {
			logWithFields.Warn("slow query")
		} else {
			logWithFields.Debug("slow query")
		}
	}
}
