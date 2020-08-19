// +build js, wasm

package db

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

// Ensure that we are implementing the ds.Batching interface.
var _ ds.Batching = &Datastore{}

// NOTE(jalextowle): Close is a noop in this implementation. We do not want a close
// operation to shut down the database connection.
func (d *Datastore) Close() error {
	return nil
}

// NOTE(jalextowle): Sync is a noop in this implementation. Operations
// such as Put and Delete are completed before a result is returned.
func (d *Datastore) Sync(ds.Key) error {
	return nil
}

// Datastore provides a Dexie implementation of the ds.Batching interface. The
// corresponding Javascript bindings can be found in
// packages/mesh-browser-lite/src/datastore.ts.
type Datastore struct {
	db         *DB
	ctx        context.Context
	dexieStore js.Value
}

type OperationType byte

const (
	PUT OperationType = iota
	DELETE
)

// Operation contains all of the data needed to communicate with the Javascript
// bindings that control access to the Dexie datastore. The Javascript bindings
// need to know what the operation should do (put or delete) and the data that
// should be used in the operation.
type Operation struct {
	operationType OperationType
	key           ds.Key
	value         []byte
}

func (o *Operation) JSValue() js.Value {
	uint8Array := js.Global().Get("Uint8Array").New(len(o.value))
	js.CopyBytesToJS(uint8Array, o.value)
	return js.ValueOf(map[string]interface{}{
		"operationType": int(o.operationType),
		"key":           o.key.String(),
		"value":         uint8Array,
	})
}

// Batch implements the ds.Batch interface, which allows Put and Delete operations
// to be queued and then committed all at once.
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

func (b *Batch) Put(key ds.Key, value []byte) error {
	b.operations = append(b.operations, &Operation{
		operationType: PUT,
		key:           key,
		value:         value,
	})
	return nil
}

func (b *Batch) Delete(key ds.Key) error {
	b.operations = append(b.operations, &Operation{
		operationType: DELETE,
		key:           key,
	})
	return nil
}

// Commit performs a batch of operations on the Dexie datastore. In this implementation,
// all of these operations occur in the same transactional context.
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

func (d *Datastore) Get(key ds.Key) ([]byte, error) {
	// FIXME - Remove Defer
	var jsResult js.Value
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(jsResult.String())
			panic(r)
		}
	}()
	var err error
	jsResult, err = jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("getAsync", key.String()))
	if err != nil {
		return nil, convertJSError(err)
	}
	result := make([]byte, jsResult.Get("length").Int())
	js.CopyBytesToGo(result, jsResult)
	return result, nil
}

func (d *Datastore) Has(key ds.Key) (bool, error) {
	jsResult, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("hasAsync", key.String()))
	if err != nil {
		return false, convertJSError(err)
	}
	return jsResult.Bool(), nil
}

func (d *Datastore) GetSize(key ds.Key) (size int, err error) {
	jsResult, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("getSizeAsync", key.String()))
	if err != nil {
		return -1, convertJSError(err)
	}
	return jsResult.Int(), nil
}

func (d *Datastore) Query(q dsq.Query) (dsq.Results, error) {
	jsResults, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("queryAsync", q.Prefix))
	if err != nil {
		return nil, convertJSError(err)
	}
	entries := make([]dsq.Entry, jsResults.Get("length").Int())
	for i := 0; i < jsResults.Get("length").Int(); i++ {
		jsResult := jsResults.Index(i)
		value := make([]byte, jsResult.Get("size").Int())
		js.CopyBytesToGo(value, jsResult.Get("value"))
		entries[i] = dsq.Entry{
			Key:   jsResult.Get("key").String(),
			Value: value,
			Size:  jsResult.Get("size").Int(),
		}
	}
	filteredEntries := []dsq.Entry{}
	for _, entry := range entries {
		passes := true
		for _, filter := range q.Filters {
			if !filter.Filter(entry) {
				passes = false
				break
			}
		}
		if passes {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	dsq.Sort(q.Orders, filteredEntries)
	if q.Offset > 0 && q.Offset <= len(filteredEntries) {
		filteredEntries = filteredEntries[q.Offset:]
	} else if q.Offset > len(filteredEntries) {
		filteredEntries = []dsq.Entry{}
	}
	if q.Limit > 0 && q.Limit <= len(filteredEntries) {
		filteredEntries = filteredEntries[:q.Limit]
	}
	return dsq.ResultsWithEntries(q, filteredEntries), nil
}

func (d *Datastore) Put(key ds.Key, value []byte) error {
	uint8Array := js.Global().Get("Uint8Array").New(len(value))
	js.CopyBytesToJS(uint8Array, value)
	_, err := jsutil.AwaitPromiseContext(d.ctx, d.dexieStore.Call("putAsync", key.String(), uint8Array))
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
