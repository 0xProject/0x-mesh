package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type dbReader interface {
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Has(key []byte, ro *opt.ReadOptions) (ret bool, err error)
	NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator
}

type dbWriter interface {
	Delete(key []byte, wo *opt.WriteOptions) error
	Put(key, value []byte, wo *opt.WriteOptions) error
}

type dbTransactor interface {
	OpenTransaction() (*leveldb.Transaction, error)
}

type dbWriterTransactor interface {
	dbWriter
	dbTransactor
}
