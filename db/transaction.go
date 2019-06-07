package db

import (
	"errors"
	"sync"
)

var (
	ErrDiscarded = errors.New("transaction has already been discarded")
	ErrCommitted = errors.New("transaction has already been committed")
)

// Transaction is an atomic database transaction which can be used to guarantee
// consistency.
type Transaction struct {
	mut             sync.Mutex
	colInfo         *colInfo
	batchWriter     dbBatchWriter
	readerWithBatch *readerWithBatchWriter
	committed       bool
	discarded       bool
}

// OpenTransaction opens and returns a new transaction. While the transaction is
// open, no other state changes (e.g. Insert, Update, or Delete) can be made to
// the database (but concurrent reads are still allowed).
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
	c.info.writeMut.Lock()
	return &Transaction{
		colInfo:         c.info.copy(),
		batchWriter:     c.ldb,
		readerWithBatch: newReaderWithBatchWriter(c.ldb),
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
	if err := txn.unsafeCheckState(); err != nil {
		txn.mut.Unlock()
		return err
	}
	txn.committed = true
	txn.mut.Unlock()
	defer txn.colInfo.writeMut.Unlock()
	return txn.batchWriter.Write(txn.readerWithBatch.batch, nil)
}

// Discard discards the transaction.
//
// Other methods should not be called after transaction has been discarded.
// However, it is safe to call Discard multiple times.
func (txn *Transaction) Discard() error {
	txn.mut.Lock()
	if txn.committed {
		return ErrCommitted
	}
	if txn.discarded {
		return nil
	}
	defer txn.mut.Unlock()
	txn.discarded = true
	txn.colInfo.writeMut.Unlock()
	return nil
}

// Insert queues an operation to insert the given model into the database. It
// returns an error if a model with the same id already exists. The model will
// not actually be inserted until the transaction is committed.
func (txn *Transaction) Insert(model Model) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return insertWithTransaction(txn.colInfo, txn, model)
}

// Update queues an operation to update an existing model in the database. It
// returns an error if the given model doesn't already exist. The model will
// not actually be updated until the transaction is committed.
func (txn *Transaction) Update(model Model) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return updateWithTransaction(txn.colInfo, txn, model)
}

// Delete queues an operation to delete the model with the given ID from the
// database. It returns an error if the model doesn't exist in the database. The
// model will not actually be deleted until the transaction is committed.
func (txn *Transaction) Delete(id []byte) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return deleteWithTransaction(txn.colInfo, txn, id)
}
