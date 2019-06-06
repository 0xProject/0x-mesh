package db

import "github.com/syndtr/goleveldb/leveldb"

type Transaction struct {
	colInfo *colInfo
	txn     *leveldb.Transaction
}

func (c *Collection) OpenTransaction() (*Transaction, error) {
	txn, err := c.ldb.OpenTransaction()
	if err != nil {
		return nil, err
	}
	return &Transaction{
		colInfo: c.info.copy(),
		txn:     txn,
	}, nil
}

func (txn *Transaction) Commit() error {
	return txn.txn.Commit()
}

func (txn *Transaction) Discard() {
	txn.txn.Discard()
}

// FindByID finds the model with the given ID and scans the results into the
// given model. As in the Unmarshal and Decode methods in the encoding/json
// package, model must be settable via reflect. Typically, this means you should
// pass in a pointer.
func (txn *Transaction) FindByID(id []byte, model Model) error {
	return findByID(txn.colInfo, txn.txn, id, model)
}

// FindAll finds all models for the collection and scans the results into the
// given models. models should be a pointer to an empty slice of a concrete
// model type (e.g. *[]myModelType).
func (txn *Transaction) FindAll(models interface{}) error {
	return findAll(txn.colInfo, txn.txn, models)
}

// New Query creates and returns a new query with the given filter. By default,
// a query will return all models that match the filter in ascending byte order
// according to their index values. The query offers methods that can be used to
// change this (e.g. Reverse and Max). The query is lazily executed, i.e. it
// does not actually touch the database until they are run. In general, queries
// have a runtime of O(N) where N is the number of models that are returned by
// the query, but using some features may significantly change this.
func (txn *Transaction) NewQuery(filter *Filter) *Query {
	return newQuery(txn.colInfo, txn.txn, filter)
}

// Insert inserts the given model into the database. It returns an error if a
// model with the same id already exists.
func (txn *Transaction) Insert(model Model) error {
	return insertWithTransaction(txn.colInfo, txn.txn, model)
}

// Update updates an existing model in the database. It returns an error if the
// given model doesn't already exist.
func (txn *Transaction) Update(model Model) error {
	return updateWithTransaction(txn.colInfo, txn.txn, model)
}

// Delete deletes the model with the given ID from the database. It returns an
// error if the model doesn't exist in the database.
func (txn *Transaction) Delete(id []byte) error {
	return deleteWithTransaction(txn.colInfo, txn.txn, id)
}
