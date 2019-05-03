package db

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/albrow/stringset"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Query is used to return certain results from the database.
type Query struct {
	col     *Collection
	filter  *Filter
	max     int
	reverse bool
}

// Filter determines which models to return in the query and what order to
// return them in.
type Filter struct {
	index *Index
	slice *util.Range
}

// New Query creates and returns a new query with the given filter. By default,
// a query will return all models that match the filter in ascending byte order
// according to their index values. The query offers methods that can be used to
// change this (e.g. Reverse and Max). The query is lazily executed, i.e. it
// does not actually touch the database until you call Run.
func (c *Collection) NewQuery(filter *Filter) *Query {
	return &Query{
		col:    c,
		filter: filter,
	}
}

// Max causes the query to only return up to max results.
func (q *Query) Max(max int) *Query {
	q.max = max
	return q
}

// Reverse causes the query to return models in descending byte order according
// to their index values instead of the default (ascending byte order).
func (q *Query) Reverse() *Query {
	q.reverse = true
	return q
}

// ValueFilter returns a Filter which will match all models with an index value
// equal to the given value.
func (index *Index) ValueFilter(val []byte) *Filter {
	prefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(val)))
	return &Filter{
		index: index,
		slice: util.BytesPrefix(prefix),
	}
}

// RangeFilter returns a Filter which will match all models with an index value
// >= start and < limit.
func (index *Index) RangeFilter(start []byte, limit []byte) *Filter {
	startWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(start)))
	limitWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(limit)))
	slice := &util.Range{Start: startWithPrefix, Limit: limitWithPrefix}
	return &Filter{
		index: index,
		slice: slice,
	}
}

// PrefixFilter returns a Filter which will match all models with an index value
// that starts with the given prefix.
func (index *Index) PrefixFilter(prefix []byte) *Filter {
	keyPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(prefix)))
	return &Filter{
		index: index,
		slice: util.BytesPrefix(keyPrefix),
	}
}

// All returns a Filter which will match all models. It is useful for when you
// want to retrieve models in sorted order without excluding any of them.
func (index *Index) All() *Filter {
	return index.PrefixFilter([]byte{})
}

// Run runs the query and scans the results into models. models should be a
// pointer to an empty slice of a concrete model type (e.g. *[]myModelType). It
// returns an error if models is the wrong type or there was a problem reading
// from the database. It does not return an error if no models match the query.
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
	// Move the iterator to the last key and then iterate backwards by calling
	// Prev instead of Next for each iteration of the for loop.
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
