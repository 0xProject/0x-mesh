// +build !js

package db

import "github.com/syndtr/goleveldb/leveldb"

// Open creates a new database using the given file path for permanent storage.
// It is not safe to have multiple DBs using the same file path.
func Open(path string) (*DB, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		ldb: ldb,
	}, nil
}
