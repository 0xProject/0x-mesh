package db

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/albrow/stringset"
)

var (
	ErrDiscarded = errors.New("transaction has already been discarded")
	ErrCommitted = errors.New("transaction has already been committed")
)

// ConflictingOperationsError is returned when two conflicting operations are attempted within the same
// transaction
type ConflictingOperationsError struct {
	operation string
}

func (e ConflictingOperationsError) Error() string {
	return fmt.Sprintf("error on %s: cannot perform more than one operation on the same model within a transaction", e.operation)
}

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
	// internalCount keeps track of the number of models inserted/deleted within
	// the transaction. An Insert increments internalCount and a Delete decrements
	// it. When the transaction is committed, internalCount is added to the
	// current count.
	internalCount int64
	// affectedIDs keeps track of the model ids that are affected by this
	// transaction.
	affectedIDs stringset.Set
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
//     (2) the transaction will fail or be discarded, in which case *none* of
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
		affectedIDs: stringset.New(),
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
// discarded. A new transaction must be created if you wish to retry the
// operations.
//
// Other methods should not be called after transaction has been committed.
func (txn *Transaction) Commit() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if err := txn.unsafeCheckState(); err != nil {
		return err
	}
	// Right before we commit, we need to update the count with txn.internalCount.
	if err := updateCountWithTransaction(txn.colInfo, txn.readWriter, int(txn.internalCount)); err != nil {
		_ = txn.Discard()
		return err
	}
	if err := txn.batchWriter.Write(txn.readWriter.batch, nil); err != nil {
		_ = txn.Discard()
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
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if err := txn.unsafeCheckState(); err != nil {
		return err
	}
	if txn.affectedIDs.Contains(string(model.ID())) {
		return ConflictingOperationsError{operation: "insert"}
	}
	if err := insertWithTransaction(txn.colInfo, txn.readWriter, model); err != nil {
		return err
	}
	txn.updateInternalCount(1)
	txn.affectedIDs.Add(string(model.ID()))
	return nil
}

// Update queues an operation to update an existing model in the database. It
// returns an error if the given model doesn't already exist. The model will
// not actually be updated until the transaction is committed.
func (txn *Transaction) Update(model Model) error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if err := txn.unsafeCheckState(); err != nil {
		return err
	}
	if txn.affectedIDs.Contains(string(model.ID())) {
		return ConflictingOperationsError{operation: "update"}
	}
	if err := updateWithTransaction(txn.colInfo, txn.readWriter, model); err != nil {
		return err
	}
	txn.affectedIDs.Add(string(model.ID()))
	return nil
}

// Delete queues an operation to delete the model with the given ID from the
// database. It returns an error if the model doesn't exist in the database. The
// model will not actually be deleted until the transaction is committed.
func (txn *Transaction) Delete(id []byte) error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if err := txn.unsafeCheckState(); err != nil {
		return err
	}
	if txn.affectedIDs.Contains(string(id)) {
		return ConflictingOperationsError{operation: "delete"}
	}
	if err := deleteWithTransaction(txn.colInfo, txn.readWriter, id); err != nil {
		return err
	}
	txn.updateInternalCount(-1)
	txn.affectedIDs.Add(string(id))
	return nil
}

func (txn *Transaction) updateInternalCount(diff int64) {
	atomic.AddInt64(&txn.internalCount, diff)
}
