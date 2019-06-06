package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Snapshot is a frozen, read-only snapshot of a DB state at a particular point
// in time.
type Snapshot struct {
	colInfo  *colInfo
	snapshot *leveldb.Snapshot
}

// GetSnapshot returns a latest snapshot of the underlying DB. The content of
// snapshot are guaranteed to be consistent. The snapshot must be released after
// use, by calling Release method.
func (c *Collection) GetSnapshot() (*Snapshot, error) {
	snapshot, err := c.ldb.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		colInfo:  c.info.copy(),
		snapshot: snapshot,
	}, nil
}

// Release releases the snapshot. This will not release any ongoing queries,
// which will still finish unless the database is closed. Other methods should
// not be called after the snapshot has been released.
func (s *Snapshot) Release() {
	s.snapshot.Release()
}

// FindByID finds the model with the given ID and scans the results into the
// given model. As in the Unmarshal and Decode methods in the encoding/json
// package, model must be settable via reflect. Typically, this means you should
// pass in a pointer.
func (s *Snapshot) FindByID(id []byte, model Model) error {
	return findByID(s.colInfo, s.snapshot, id, model)
}

// FindAll finds all models for the collection and scans the results into the
// given models. models should be a pointer to an empty slice of a concrete
// model type (e.g. *[]myModelType).
func (s *Snapshot) FindAll(models interface{}) error {
	return findAll(s.colInfo, s.snapshot, models)
}

// New Query creates and returns a new query with the given filter. By default,
// a query will return all models that match the filter in ascending byte order
// according to their index values. The query offers methods that can be used to
// change this (e.g. Reverse and Max). The query is lazily executed, i.e. it
// does not actually touch the database until they are run. In general, queries
// have a runtime of O(N) where N is the number of models that are returned by
// the query, but using some features may significantly change this.
func (s *Snapshot) NewQuery(filter *Filter) *Query {
	return newQuery(s.colInfo, s.snapshot, filter)
}
