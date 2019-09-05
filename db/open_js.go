// +build js,wasm

package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

// Open creates a new database using in-memory storage. The given file path is
// ignored.
func Open(path string) (*DB, error) {
	ldb, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		ldb: ldb,
	}, nil
}
