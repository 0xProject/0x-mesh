// +build js, wasm

package db

import (
	"context"
	"syscall/js"

	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

var _ ds.Batching = &Datastore{}

type Datastore struct {
	db         *DB
	ctx        context.Context
	dexieStore js.Value
}

// io.Closer

// FIXME - Is this what we want?
func (d *Datastore) Close() error {
	// Noop
	return nil
}

// Sync

func (d *Datastore) Sync(ds.Key) error {
	// Noop
	return nil
}

/// Write

func (d *Datastore) Put(key ds.Key, value []byte) error {
	_, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("putAsync", key.String(), string(value)))
	if err != nil {
		return convertJSError(err)
	}
	return nil
}

func (d *Datastore) Delete(key ds.Key) error {
	_, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("deleteAsync", key.String()))
	if err != nil {
		return convertJSError(err)
	}
	return nil
}

// Read

func (d *Datastore) Get(key ds.Key) ([]byte, error) {
	jsResult, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("getAsync", key.String()))
	if err != nil {
		return nil, convertJSError(err)
	}
	return []byte(jsResult.String()), nil
}

func (d *Datastore) Has(key ds.Key) (bool, error) {
	jsResult, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("hasAsync", key.String()))
	if err != nil {
		return false, convertJSError(err)
	}
	return jsResult.Bool(), nil
}

func (d *Datastore) GetSize(key ds.Key) (int, error) {
	jsResult, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("getSizeAsync", key.String()))
	if err != nil {
		return 0, convertJSError(err)
	}
	return jsResult.Int(), nil
}

func (d *Datastore) Query(q dsq.Query) (dsq.Results, error) {
	jsResults, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("queryAsync", convertQueryToJS(q)))
	if err != nil {
		return nil, convertJSError(err)
	}
	entries := make([]dsq.Entry, jsResults.Get("length").Int())
	for i := 0; i < jsResults.Get("length").Int(); i++ {
		jsResult := jsResults.Index(i)
		entries[i] = dsq.Entry{
			Key:   jsResult.Get("key").String(),
			Value: []byte(jsResult.Get("value").String()),
			Size:  jsResult.Get("size").Int(),
		}
	}
	return dsq.ResultsWithEntries(q, entries), nil
}

/// Batching

type OperationType byte

const (
	ADDITION OperationType = iota
	// FIXME - Rename to DELETION
	REMOVAL
)

type Operation struct {
	operationType OperationType
	key           ds.Key
	value         []byte
}

func (o *Operation) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"operationType": int(o.operationType),
		"key":           o.key.String(),
		"value":         string(o.value),
	})
}

type Batch struct {
	ctx        context.Context
	dexieStore js.Value
	operations []*Operation
}

func (d *Datastore) Batch() (ds.Batch, error) {
	return &Batch{
		ctx:        d.ctx,
		dexieStore: d.dexieStore,
	}, nil
}

func (b *Batch) Commit() error {
	convertibleOperations := make([]interface{}, len(b.operations))
	for i, operation := range b.operations {
		convertibleOperations[i] = interface{}(operation)
	}
	_, err := jsutil.AwaitPromiseContext(b.ctx, b.dexieStore.Call("commitAsync", convertibleOperations))
	if err != nil {
		return convertJSError(err)
	}
	return nil
}

func (b *Batch) Put(key ds.Key, value []byte) error {
	b.operations = append(b.operations, &Operation{
		operationType: ADDITION,
		key:           key,
		value:         value,
	})
	return nil
}

func (b *Batch) Delete(key ds.Key) error {
	b.operations = append(b.operations, &Operation{
		operationType: REMOVAL,
		key:           key,
	})
	return nil
}

/// js conversions

// FIXME - length checks and code dedupe
func convertQueryToJS(q dsq.Query) js.Value {
	jsFilters := make([]interface{}, len(q.Filters))
	for i, filter := range q.Filters {
		jsFilters[i] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			entry := dsq.Entry{
				Key:   args[0].Get("key").String(),
				Value: []byte(args[0].Get("value").String()),
				Size:  args[0].Get("size").Int(),
			}
			return filter.Filter(entry)
		})
	}
	jsOrders := make([]interface{}, len(q.Orders))
	for i, order := range q.Orders {
		jsOrders[i] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			a := dsq.Entry{
				Key:   args[0].Get("key").String(),
				Value: []byte(args[0].Get("value").String()),
				Size:  args[0].Get("size").Int(),
			}
			b := dsq.Entry{
				Key:   args[1].Get("key").String(),
				Value: []byte(args[1].Get("value").String()),
				Size:  args[1].Get("size").Int(),
			}
			return order.Compare(a, b)
		})
	}
	return js.ValueOf(map[string]interface{}{
		"prefix":  q.Prefix,
		"filters": jsFilters,
		"orders":  jsOrders,
		"limit":   q.Limit,
		"offset":  q.Offset,
	})
}
