package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// dbReader is an interface that encapsulates read-only functionality.
type dbReader interface {
	leveldb.Reader
	Has(key []byte, ro *opt.ReadOptions) (bool, error)
}

// dbWriter is an interface that encapsulates write/update functionality.
type dbWriter interface {
	Delete(key []byte, wo *opt.WriteOptions) error
	Put(key, value []byte, wo *opt.WriteOptions) error
}

type dbBatchWriter interface {
	Write(batch *leveldb.Batch, ro *opt.WriteOptions) error
}

type dbReadWriter interface {
	dbReader
	dbWriter
}

type dbReadBatchWriter interface {
	dbReadWriter
	dbBatchWriter
}
