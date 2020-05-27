// +build js,wasm

package db

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"syscall/js"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/ethereum/go-ethereum/common"
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
	jsOrders, err := jsutil.InefficientlyConvertToJS(orders)
	if err != nil {
		return nil, nil, err
	}
	jsResult, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("addOrdersAsync", jsOrders))
	if err != nil {
		return nil, nil, convertJSError(err)
	}
	jsAdded := jsResult.Get("added")
	if err := jsutil.InefficientlyConvertFromJS(jsAdded, &added); err != nil {
		return nil, nil, err
	}
	jsRemoved := jsResult.Get("removed")
	if err := jsutil.InefficientlyConvertFromJS(jsRemoved, &removed); err != nil {
		return nil, nil, err
	}
	return added, removed, nil
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
	order = &types.OrderWithMetadata{}
	if err := jsutil.InefficientlyConvertFromJS(jsOrder, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (db *DB) FindOrders(query *OrderQuery) (orders []*types.OrderWithMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsOrders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("findOrdersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	if err := jsutil.InefficientlyConvertFromJS(jsOrders, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (db *DB) CountOrders(query *OrderQuery) (count int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
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
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsOrders, err := jsutil.AwaitPromiseContext(db.ctx, db.dexie.Call("deleteOrdersAsync", query))
	if err != nil {
		return nil, convertJSError(err)
	}
	if err := jsutil.InefficientlyConvertFromJS(jsOrders, &deletedOrders); err != nil {
		return nil, err
	}
	return deletedOrders, nil
}

func (db *DB) UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()
	jsUpdateFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		jsExistingOrder := args[0]
		var existingOrder types.OrderWithMetadata
		if err := jsutil.InefficientlyConvertFromJS(jsExistingOrder, &existingOrder); err != nil {
			panic(err)
		}
		orderToUpdate, err := updateFunc(&existingOrder)
		if err != nil {
			panic(err)
		}
		jsOrderToUpdate, err := jsutil.InefficientlyConvertToJS(orderToUpdate)
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
	return nil, nil, errors.New("not yet implemented")
}

func (db *DB) GetMiniHeader(hash common.Hash) (*types.MiniHeader, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) FindMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) DeleteMiniHeader(hash common.Hash) error {
	return errors.New("not yet implemented")
}

func (db *DB) DeleteMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) GetMetadata() (*types.Metadata, error) {
	return nil, errors.New("not yet implemented")
}

func (db *DB) SaveMetadata(metadata *types.Metadata) error {
	return errors.New("not yet implemented")
}

func (db *DB) UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error {
	return errors.New("not yet implemented")
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
		}
	}
	return e
}
