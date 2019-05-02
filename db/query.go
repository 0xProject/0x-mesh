package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/albrow/stringset"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Query struct {
	Filter *Filter
	Limit  *int
	Offset *int
	// TOOD(albrow): Add option for ASC or DESC order.
}

type Filter struct {
	index *Index
	slice *util.Range
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

// RunQuery runs the given query and scans the results into models. models
// should be a pointer to an empty slice of a concrete model type (e.g.
// *[]myModelType).
func (c *Collection) RunQuery(query *Query, models interface{}) error {

	// TODO(albrow): Respect query.Limit and query.Offset.

	if err := c.checkModelsType(models); err != nil {
		return err
	}

	// Get the appropriate iterator and index.
	var iter iterator.Iterator
	var index *Index
	if query.Filter != nil {
		iter = c.db.ldb.NewIterator(query.Filter.slice, nil)
		index = query.Filter.index
	} else {
		// TODO(albrow): Use a default iterator and index.
		return errors.New("Query.Filter is required")
	}
	defer iter.Release()

	// MultiIndexes can result in the same model being included more than once. To
	// prevent this, we keep track of the primaryKeys we have already seen.
	pkSet := stringset.New()
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		// We assume that each key in the iterator consists of an index prefix, the
		// value for a particular model, and the model ID. We can extract a primary
		// key from this key and use it to get the encoded data for the model
		// itself.
		pk := index.primaryKeyFromIndexKey(iter.Key())
		if pkSet.Contains(string(pk)) {
			continue
		}
		pkSet.Add(string(pk))
		data, err := c.db.ldb.Get(pk, nil)
		if err != nil {
			return err
		}
		model := reflect.New(c.modelType)
		if err := json.Unmarshal(data, model.Interface()); err != nil {
			return err
		}
		modelsVal.Set(reflect.Append(modelsVal, model.Elem()))
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}
