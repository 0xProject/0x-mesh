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
	colInfo *colInfo
	reader  dbReader
	filter  *Filter
	max     int
	offset  int
	reverse bool
}

// Filter determines which models to return in the query and what order to
// return them in.
type Filter struct {
	index *Index
	slice *util.Range
}

func newQuery(colInfo *colInfo, reader dbReader, filter *Filter) *Query {
	return &Query{
		colInfo: colInfo,
		reader:  reader,
		filter:  filter,
	}
}

// Max causes the query to only return up to max results. It is the analog of
// the LIMIT keyword in SQL:
// https://www.postgresql.org/docs/current/queries-limit.html
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

// Offset causes the query to skip offset models when iterating through models
// that match the query. Note that queries which use an offset have a runtime
// of O(max(K, offset) + N), where N is the number of models returned by the
// query and K is the total number of keys in the corresponding index. Queries
// with a high offset can take a long time to run, regardless of the number of
// models returned. This is due to limitations of the underlying database.
// Offset is the analog of the OFFSET keyword in SQL:
// https://www.postgresql.org/docs/current/queries-limit.html
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// ValueFilter returns a Filter which will match all models with an index value
// equal to the given value.
func (index *Index) ValueFilter(val []byte) *Filter {
	prefix := []byte(fmt.Sprintf("%s:%s:", index.prefix(), escape(val)))
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
	if err := q.colInfo.checkModelsType(models); err != nil {
		return err
	}

	iter := q.reader.NewIterator(q.filter.slice, nil)
	defer iter.Release()
	if q.reverse {
		return q.getModelsWithIteratorReverse(iter, models)
	}
	return q.getModelsWithIteratorForward(iter, models)
}

// Count returns the number of unique models that match the query. It does not
// return an error if no models match the query. Note that this method *does*
// respect q.Max. If the number of models that match the filter is greater than
// q.Max, it will stop counting and return q.Max.
func (q *Query) Count() (int, error) {
	iter := q.reader.NewIterator(q.filter.slice, nil)
	defer iter.Release()
	pkSet := stringset.New()
	for i := 0; iter.Next() && iter.Error() == nil; i++ {
		if i < q.offset {
			continue
		}
		pk := q.filter.index.primaryKeyFromIndexKey(iter.Key())
		pkSet.Add(string(pk))
		if q.max != 0 && len(pkSet) >= q.max {
			break
		}
	}
	if iter.Error() != nil {
		return 0, iter.Error()
	}
	return len(pkSet), nil
}

func (q *Query) getModelsWithIteratorForward(iter iterator.Iterator, models interface{}) error {
	// MultiIndexes can result in the same model being included more than once. To
	// prevent this, we keep track of the primaryKeys we have already seen using
	// pkSet.
	pkSet := stringset.New()
	modelsVal := reflect.ValueOf(models).Elem()
	for i := 0; iter.Next() && iter.Error() == nil; i++ {
		if i < q.offset {
			continue
		}
		if err := q.getAndAppendModelIfUnique(q.filter.index, pkSet, iter.Key(), modelsVal); err != nil {
			return err
		}
		if q.max != 0 && modelsVal.Len() >= q.max {
			return iter.Error()
		}
	}
	return iter.Error()
}

func (q *Query) getModelsWithIteratorReverse(iter iterator.Iterator, models interface{}) error {
	pkSet := stringset.New()
	modelsVal := reflect.ValueOf(models).Elem()
	// Move the iterator to the last key and then iterate backwards by calling
	// Prev instead of Next for each iteration of the for loop.
	iter.Last()
	iter.Next()
	for i := 0; iter.Prev() && iter.Error() == nil; i++ {
		if i < q.offset {
			continue
		}
		if err := q.getAndAppendModelIfUnique(q.filter.index, pkSet, iter.Key(), modelsVal); err != nil {
			return err
		}
		if q.max != 0 && modelsVal.Len() >= q.max {
			return iter.Error()
		}
	}
	return iter.Error()
}

func (q *Query) getAndAppendModelIfUnique(index *Index, pkSet stringset.Set, key []byte, modelsVal reflect.Value) error {
	// We assume that each key in the iterator consists of an index prefix, the
	// value for a particular model, and the model ID. We can extract a primary
	// key from this key and use it to get the encoded data for the model
	// itself.
	pk := index.primaryKeyFromIndexKey(key)
	if pkSet.Contains(string(pk)) {
		return nil
	}
	pkSet.Add(string(pk))
	data, err := q.reader.Get(pk, nil)
	if err != nil {
		return err
	}
	model := reflect.New(q.colInfo.modelType)
	if err := json.Unmarshal(data, model.Interface()); err != nil {
		return err
	}
	modelsVal.Set(reflect.Append(modelsVal, model.Elem()))
	return nil
}
