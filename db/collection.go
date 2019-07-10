package db

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

// Collection represents a set of a specific type of model.
type Collection struct {
	info *colInfo
	ldb  *leveldb.DB
}

// NewCollection creates and returns a new collection with the given name and
// model type. You should create exactly one collection for each model type. The
// collection should typically be created once at the start of your application
// and re-used. NewCollection returns an error if a collection has already been
// created with the given name for this db.
func (db *DB) NewCollection(name string, typ Model) (*Collection, error) {
	col := &Collection{
		info: &colInfo{
			db:        db,
			name:      name,
			modelType: reflect.TypeOf(typ),
			writeMut:  &sync.Mutex{},
		},
		ldb: db.ldb,
	}
	db.colLock.Lock()
	defer db.colLock.Unlock()
	for _, existingCol := range db.collections {
		if existingCol.info.name == name {
			return nil, fmt.Errorf("a collection with the name %q already exists", name)
		}
	}
	db.collections = append(db.collections, col)
	return col, nil
}

// Name returns the name of the collection.
func (c *Collection) Name() string {
	return c.info.name
}

// FindByID finds the model with the given ID and scans the results into the
// given model. As in the Unmarshal and Decode methods in the encoding/json
// package, model must be settable via reflect. Typically, this means you should
// pass in a pointer.
func (c *Collection) FindByID(id []byte, model Model) error {
	return findByID(c.info, c.ldb, id, model)
}

// FindAll finds all models for the collection and scans the results into the
// given models. models should be a pointer to an empty slice of a concrete
// model type (e.g. *[]myModelType).
func (c *Collection) FindAll(models interface{}) error {
	return findAll(c.info, c.ldb, models)
}

// Count returns the number of models in the collection.
func (c *Collection) Count() (int, error) {
	return count(c.info, c.ldb)
}

// Insert inserts the given model into the database. It returns an error if a
// model with the same id already exists.
func (c *Collection) Insert(model Model) error {
	txn := c.OpenTransaction()
	if err := insertWithTransaction(c.info, txn.readWriter, model); err != nil {
		_ = txn.Discard()
		return err
	}
	txn.updateInternalCount(1)
	if err := txn.Commit(); err != nil {
		_ = txn.Discard()
		return err
	}
	return nil
}

// Update updates an existing model in the database. It returns an error if the
// given model doesn't already exist.
func (c *Collection) Update(model Model) error {
	txn := c.OpenTransaction()
	if err := updateWithTransaction(c.info, txn.readWriter, model); err != nil {
		_ = txn.Discard()
		return err
	}
	if err := txn.Commit(); err != nil {
		_ = txn.Discard()
		return err
	}
	return nil
}

// Delete deletes the model with the given ID from the database. It returns an
// error if the model doesn't exist in the database.
func (c *Collection) Delete(id []byte) error {
	txn := c.OpenTransaction()
	if err := deleteWithTransaction(c.info, txn.readWriter, id); err != nil {
		_ = txn.Discard()
		return err
	}
	txn.updateInternalCount(-1)
	if err := txn.Commit(); err != nil {
		_ = txn.Discard()
		return err
	}
	return nil
}

// New Query creates and returns a new query with the given filter. By default,
// a query will return all models that match the filter in ascending byte order
// according to their index values. The query offers methods that can be used to
// change this (e.g. Reverse and Max). The query is lazily executed, i.e. it
// does not actually touch the database until they are run. In general, queries
// have a runtime of O(N) where N is the number of models that are returned by
// the query, but using some features may significantly change this.
func (c *Collection) NewQuery(filter *Filter) *Query {
	return newQuery(c.info, c.ldb, filter)
}
