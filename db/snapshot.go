package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Snapshot is a frozen snapshot of a DB state at a particular point in time.
type Snapshot struct {
	*readOnlyCollection
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
		readOnlyCollection: c.readOnlyCollection,
		snapshot:           snapshot,
	}, nil
}

// Release releases the snapshot. This will not release any ongoing queries,
// which will still finish unless the database is closed. Other methods should
// not be called after the snapshot has been released.
func (s *Snapshot) Release() {
	s.snapshot.Release()
}
