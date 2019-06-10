package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type readerWithBatchWriter struct {
	reader dbReader
	batch  *leveldb.Batch
}

func newReaderWithBatchWriter(reader dbReader) *readerWithBatchWriter {
	return &readerWithBatchWriter{
		reader: reader,
		batch:  &leveldb.Batch{},
	}
}

var _ dbReadWriter = &readerWithBatchWriter{}

func (readWriter *readerWithBatchWriter) Get(key []byte, ro *opt.ReadOptions) ([]byte, error) {
	return readWriter.reader.Get(key, ro)
}

func (readWriter *readerWithBatchWriter) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return readWriter.reader.NewIterator(slice, ro)
}

func (readWriter *readerWithBatchWriter) Has(key []byte, ro *opt.ReadOptions) (bool, error) {
	return readWriter.reader.Has(key, ro)
}

func (readWriter *readerWithBatchWriter) Delete(key []byte, wo *opt.WriteOptions) error {
	readWriter.batch.Delete(key)
	return nil
}

func (readWriter *readerWithBatchWriter) Put(key, value []byte, wo *opt.WriteOptions) error {
	readWriter.batch.Put(key, value)
	return nil
}
