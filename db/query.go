package db

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/albrow/stringset"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// TODO(albrow): Test filter, max, and reverse.

type Query struct {
	col     *Collection
	filter  *Filter
	max     int
	reverse bool
}

type Filter struct {
	index *Index
	slice *util.Range
}

func (c *Collection) NewQuery(filter *Filter) *Query {
	return &Query{
		col:    c,
		filter: filter,
	}
}

func (q *Query) Max(max int) *Query {
	q.max = max
	return q
}

func (q *Query) Reverse() *Query {
	q.reverse = true
	return q
}

// copy returns a shallow copy of the query, which can then be modified without
// affecting the original. Intended for testing only.
func (q *Query) copy() *Query {
	return &Query{
		col:     q.col,
		filter:  q.filter,
		max:     q.max,
		reverse: q.reverse,
	}
}

// ValueFilter returns a Filter which will match all models with the given value
// according to the index.
func (index *Index) ValueFilter(val []byte) *Filter {
	prefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(val)))
	return &Filter{
		index: index,
		slice: util.BytesPrefix(prefix),
	}
}

// RangeFilter returns a Filter which will match all models with a value >=
// start and < limit according to the index.
func (index *Index) RangeFilter(start []byte, limit []byte) *Filter {
	startWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(start)))
	limitWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(limit)))
	slice := &util.Range{Start: startWithPrefix, Limit: limitWithPrefix}
	return &Filter{
		index: index,
		slice: slice,
	}
}

// PrefixFilter returns a Filter which will match all models with a value that
// starts with the given prefix according to the index.
func (index *Index) PrefixFilter(prefix []byte) *Filter {
	keyPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(prefix)))
	return &Filter{
		index: index,
		slice: util.BytesPrefix(keyPrefix),
	}
}

// Run runs the query and scans the results into models. models
// should be a pointer to an empty slice of a concrete model type (e.g.
// *[]myModelType).
func (q *Query) Run(models interface{}) error {
	if err := q.col.checkModelsType(models); err != nil {
		return err
	}

	iter := q.col.db.ldb.NewIterator(q.filter.slice, nil)
	defer iter.Release()
	index := q.filter.index
	if q.reverse {
		return getModelsWithIteratorReverse(iter, index, q.max, models)
	}
	return getModelsWithIteratorForwards(iter, index, q.max, models)
}

func getModelsWithIteratorForwards(iter iterator.Iterator, index *Index, max int, models interface{}) error {
	// MultiIndexes can result in the same model being included more than once. To
	// prevent this, we keep track of the primaryKeys we have already seen using
	// pkSet.
	pkSet := stringset.New()
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		if err := getAndAppendModelIfUnique(index, pkSet, iter.Key(), modelsVal); err != nil {
			return err
		}
		if max != 0 && modelsVal.Len() >= max {
			return iter.Error()
		}
	}
	return iter.Error()
}

func getModelsWithIteratorReverse(iter iterator.Iterator, index *Index, max int, models interface{}) error {
	pkSet := stringset.New()
	modelsVal := reflect.ValueOf(models).Elem()
	// Move the iterator to the last key and then move backwards.
	iter.Last()
	iter.Next()
	for iter.Prev() {
		if err := getAndAppendModelIfUnique(index, pkSet, iter.Key(), modelsVal); err != nil {
			return err
		}
		if max != 0 && modelsVal.Len() >= max {
			return iter.Error()
		}
	}
	return iter.Error()
}

func getAndAppendModelIfUnique(index *Index, pkSet stringset.Set, key []byte, modelsVal reflect.Value) error {
	// We assume that each key in the iterator consists of an index prefix, the
	// value for a particular model, and the model ID. We can extract a primary
	// key from this key and use it to get the encoded data for the model
	// itself.
	pk := index.primaryKeyFromIndexKey(key)
	if pkSet.Contains(string(pk)) {
		return nil
	}
	pkSet.Add(string(pk))
	data, err := index.col.db.ldb.Get(pk, nil)
	if err != nil {
		return err
	}
	model := reflect.New(index.col.modelType)
	if err := json.Unmarshal(data, model.Interface()); err != nil {
		return err
	}
	modelsVal.Set(reflect.Append(modelsVal, model.Elem()))
	return nil
}
