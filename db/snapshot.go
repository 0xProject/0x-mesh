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
	c.indexMut.RLock()
	indexes := make([]*Index, len(c.indexes))
	copy(indexes, c.indexes)
	c.indexMut.RUnlock()
	return &Snapshot{
		readOnlyCollection: &readOnlyCollection{
			reader:    snapshot,
			name:      c.name,
			modelType: c.modelType,
			indexes:   indexes,
		},
		snapshot: snapshot,
	}, nil
}

// Release releases the snapshot. This will not release any ongoing queries,
// which will still finish unless the database is closed. Other methods should
// not be called after the snapshot has been released.
func (s *Snapshot) Release() {
	s.snapshot.Release()
}
