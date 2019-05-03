package db

import (
	"fmt"
	"strings"
)

// Index can be used to search for specific values or specific ranges of values
// for a collection.
type Index struct {
	col    *Collection
	name   string
	getter func(m Model) [][]byte
}

// AddIndex creates and returns a new index. name is an arbitrary, unique name
// for the index. getter is a function that accepts a model and returns the
// value for this particular index. For example, if you wanted to add an index
// on a struct field, getter should return the value of that field. After
// AddIndex is called, any new models in this collection that are inserted will
// be indexed. Any models inserted prior to calling AddIndex will *not* be
// indexed. Note that in order to function correctly, indexes must be based on
// data that is actually saved to the database (e.g. exported struct fields).
func (c *Collection) AddIndex(name string, getter func(Model) []byte) *Index {
	// Internally, all indexes are treated as MultiIndexes. We wrap the given
	// getter function so that it returns [][]byte instead of just []byte.
	wrappedGetter := func(model Model) [][]byte {
		return [][]byte{getter(model)}
	}
	return c.AddMultiIndex(name, wrappedGetter)
}

// AddMultiIndex is like AddIndex but has the ability to index multiple values
// for the same model. For methods like FindWithRange and FindWithValue, the
// model will be included in the results if *any* of the values returned by the
// getter function satisfy the constraints. It is useful for representing
// one-to-many relationships. Any models inserted prior to calling AddMultiIndex
// will *not* be indexed. Note that in order to function correctly, indexes must
// be based on data that is actually saved to the database (e.g. exported struct fields).
func (c *Collection) AddMultiIndex(name string, getter func(Model) [][]byte) *Index {
	index := &Index{
		col:    c,
		name:   name,
		getter: getter,
	}
	c.indexes = append(c.indexes, index)
	return index
}

// Name returns the name of the index.
func (index *Index) Name() string {
	return index.name
}

func (index *Index) prefix() []byte {
	return []byte(fmt.Sprintf("index:%s:%s", index.col.name, index.name))
}

func (index *Index) keysForModel(model Model) [][]byte {
	values := index.getter(model)
	indexKeys := make([][]byte, len(values))
	for i, value := range values {
		indexKeys[i] = []byte(fmt.Sprintf("%s:%s:%s", index.prefix(), escape(value), escape(model.ID())))
	}
	return indexKeys
}

// primaryKeyFromIndexKey extracts and returns the primary key from the given index
// key.
func (index *Index) primaryKeyFromIndexKey(key []byte) []byte {
	pkAndVal := strings.TrimPrefix(string(key), string(index.prefix()))
	split := strings.Split(pkAndVal, ":")
	return index.col.primaryKeyForIDWithoutEscape([]byte(split[2]))
}
