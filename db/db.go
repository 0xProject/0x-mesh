package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Note about the implementation:
//
// There are two types of keys used. A "primary key" is the main key for a
// particular model. It's value is the encoded data for that model. The format
// for a primary key is: `model:<collection name>:<model ID>`.
//
// An "index key" is used in queries to find models with specific indexed
// values. The format for an index key is:
// `index:<collection name>:<index name>:<value>:<model ID>`. Unlike primary
// keys, index keys have no values and don't store any actual data. Instead, the
// primary key can be extracted from an index key and then used to look up the
// data for the corresponding model.

// Model is any type which can be inserted and retrieved from the database. The
// only requirement is an ID method. Because the db package uses reflect to
// encode/decode models, only exported struct fields will be saved and retrieved
// from the database.
type Model interface {
	// ID returns a unique identifier for this model.
	ID() []byte
}

// DB is the top-level Database.
type DB struct {
	ldb *leveldb.DB
}

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

// Close closes the database. It is not safe to call Close if there are any
// other methods that have not yet returned. It is safe to call Close multiple
// times.
func (db *DB) Close() error {
	return db.ldb.Close()
}
