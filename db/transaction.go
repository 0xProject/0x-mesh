package db

import (
	"errors"
	"sync"
)

var (
	ErrDiscarded = errors.New("transaction has already been discarded")
	ErrCommitted = errors.New("transaction has already been committed")
)

// Transaction is an atomic database transaction for a single collection which
// can be used to guarantee consistency.
type Transaction struct {
	db          *DB
	mut         sync.Mutex
	colInfo     *colInfo
	batchWriter dbBatchWriter
	readWriter  *readerWithBatchWriter
	committed   bool
	discarded   bool
}

// OpenTransaction opens and returns a new transaction for the collection. While
// the transaction is open, no other state changes (e.g. Insert, Update, or
// Delete) can be made to the collection (but concurrent reads are still
// allowed).
//
// Transactions are atomic, meaning that either:
//
//     (1) The transaction will succeed and *all* queued operations will be
//     applied, or
//     (2) the transaction will be fail or be discarded, in which case *none* of
//     the queued operations will be applied.
//
// The transaction must be closed once done, either by committing or discarding
// the transaction. No changes will be made to the database state until the
// transaction is committed.
func (c *Collection) OpenTransaction() *Transaction {
	// Note we acquire an RLock on the global write mutex. We're not really a
	// "reader" but we behave like one in the context of an RWMutex. Up to one
	// write lock for each collection can be held, or one global write lock can be
	// held at any given time.
	c.info.db.globalWriteLock.RLock()
	c.info.writeMut.Lock()
	return &Transaction{
		db:          c.info.db,
		colInfo:     c.info.copy(),
		batchWriter: c.ldb,
		readWriter:  newReaderWithBatchWriter(c.ldb),
	}
}

// checkState acquires a lock on txn.mut and then calls unsafeCheckState.
func (txn *Transaction) checkState() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	return txn.unsafeCheckState()
}

// unsafeCheckState checks the state of the transaction, assuming the caller has
// already acquired a lock. It returns an error if the transaction has already
// been committed or discarded.
func (txn *Transaction) unsafeCheckState() error {
	if txn.discarded {
		return ErrDiscarded
	} else if txn.committed {
		return ErrCommitted
	}
	return nil
}

// Commit commits the transaction. If error is not nil, then the transaction is
// not committed, it can then either be retried or discarded.
//
// Other methods should not be called after transaction has been committed.
func (txn *Transaction) Commit() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if err := txn.unsafeCheckState(); err != nil {
		return err
	}
	if err := txn.batchWriter.Write(txn.readWriter.batch, nil); err != nil {
		return err
	}
	txn.committed = true
	txn.colInfo.writeMut.Unlock()
	txn.db.globalWriteLock.RUnlock()
	return nil
}

// Discard discards the transaction.
//
// Other methods should not be called after transaction has been discarded.
// However, it is safe to call Discard multiple times.
func (txn *Transaction) Discard() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if txn.committed {
		return ErrCommitted
	}
	if txn.discarded {
		return nil
	}
	txn.discarded = true
	txn.colInfo.writeMut.Unlock()
	txn.db.globalWriteLock.RUnlock()
	return nil
}

// Insert queues an operation to insert the given model into the database. It
// returns an error if a model with the same id already exists. The model will
// not actually be inserted until the transaction is committed.
func (txn *Transaction) Insert(model Model) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return insertWithTransaction(txn.colInfo, txn.readWriter, model)
}

// Update queues an operation to update an existing model in the database. It
// returns an error if the given model doesn't already exist. The model will
// not actually be updated until the transaction is committed.
func (txn *Transaction) Update(model Model) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return updateWithTransaction(txn.colInfo, txn.readWriter, model)
}

// Delete queues an operation to delete the model with the given ID from the
// database. It returns an error if the model doesn't exist in the database. The
// model will not actually be deleted until the transaction is committed.
func (txn *Transaction) Delete(id []byte) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return deleteWithTransaction(txn.colInfo, txn.readWriter, id)
}
