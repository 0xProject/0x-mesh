package db

import (
	"sync"
)

// GlobalTransaction is an atomic database transaction across all collections
// which can be used to guarantee consistency.
type GlobalTransaction struct {
	db          *DB
	mut         sync.Mutex
	batchWriter dbBatchWriter
	readWriter  *readerWithBatchWriter
	committed   bool
	discarded   bool
	// internalCounts keeps track of the number of models inserted/deleted within
	// the transaction for each collection. An Insert increments the count and
	// a Delete decrements it. When the transaction is committed, the
	// internal count is added to the current count for each collection.
	internalCounts map[*Collection]int
}

// OpenGlobalTransaction opens and returns a new global transaction. While the
// transaction is open, no other state changes (e.g. Insert, Update, or Delete)
// can be made to the database (but concurrent reads are still allowed). This
// includes all collections.
//
// No new collections can be created while the global transaction is open.
// Calling NewCollection while the transaction is open will block until the
// transaction is committed or discarded.
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
func (db *DB) OpenGlobalTransaction() *GlobalTransaction {
	// Note we acquire a Lock on the global write mutex. We're not really a
	// "writer" but we behave like one in the context of an RWMutex. Up to one
	// write lock for each collection can be held, or one global write lock can be
	// held at any given time.
	db.colLock.Lock()
	db.globalWriteLock.Lock()
	return &GlobalTransaction{
		db:             db,
		batchWriter:    db.ldb,
		readWriter:     newReaderWithBatchWriter(db.ldb),
		internalCounts: map[*Collection]int{},
	}
}

// checkState acquires a lock on txn.mut and then calls unsafeCheckState.
func (txn *GlobalTransaction) checkState() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	return txn.unsafeCheckState()
}

// unsafeCheckState checks the state of the transaction, assuming the caller has
// already acquired a lock. It returns an error if the transaction has already
// been committed or discarded.
func (txn *GlobalTransaction) unsafeCheckState() error {
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
func (txn *GlobalTransaction) Commit() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if err := txn.unsafeCheckState(); err != nil {
		return err
	}
	// Right before we commit, we need to update the count for each collection
	// that was touched.
	for col, internalCount := range txn.internalCounts {
		if err := updateCountWithTransaction(col.info, txn.readWriter, int(internalCount)); err != nil {
			_ = txn.Discard()
			return err
		}
	}
	if err := txn.batchWriter.Write(txn.readWriter.batch, nil); err != nil {
		_ = txn.Discard()
		return err
	}
	txn.committed = true
	txn.db.globalWriteLock.Unlock()
	txn.db.colLock.Unlock()
	return nil
}

// Discard discards the transaction.
//
// Other methods should not be called after transaction has been discarded.
// However, it is safe to call Discard multiple times.
func (txn *GlobalTransaction) Discard() error {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if txn.committed {
		return ErrCommitted
	}
	if txn.discarded {
		return nil
	}
	txn.discarded = true
	txn.db.globalWriteLock.Unlock()
	txn.db.colLock.Unlock()
	return nil
}

// Insert queues an operation to insert the given model into the given
// collection. It returns an error if a model with the same id already exists.
// The model will not actually be inserted until the transaction is committed.
func (txn *GlobalTransaction) Insert(col *Collection, model Model) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	if err := insertWithTransaction(col.info, txn.readWriter, model); err != nil {
		return err
	}
	txn.updateInternalCount(col, 1)
	return nil
}

// Update queues an operation to update an existing model in the given
// collection. It returns an error if the given model doesn't already exist. The
// model will not actually be updated until the transaction is committed.
func (txn *GlobalTransaction) Update(col *Collection, model Model) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	return updateWithTransaction(col.info, txn.readWriter, model)
}

// Delete queues an operation to delete the model with the given ID from the
// given collection. It returns an error if the model doesn't exist in the
// database. The model will not actually be deleted until the transaction is
// committed.
func (txn *GlobalTransaction) Delete(col *Collection, id []byte) error {
	if err := txn.checkState(); err != nil {
		return err
	}
	if err := deleteWithTransaction(col.info, txn.readWriter, id); err != nil {
		return err
	}
	txn.updateInternalCount(col, -1)
	return nil
}

func (txn *GlobalTransaction) updateInternalCount(col *Collection, diff int) {
	txn.mut.Lock()
	defer txn.mut.Unlock()
	if existingCount, found := txn.internalCounts[col]; found {
		txn.internalCounts[col] = existingCount + diff
	} else {
		txn.internalCounts[col] = diff
	}
}
