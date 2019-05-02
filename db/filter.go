package db

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Filter struct {
	index    *Index
	iterator iterator.Iterator
}

// WithValue returns a Filter which will match all models with the given value
// according to the index.
func (index *Index) WithValue(val []byte) *Filter {
	prefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(val)))
	prefixRange := util.BytesPrefix(prefix)
	iterator := index.col.db.ldb.NewIterator(prefixRange, nil)
	return &Filter{
		index:    index,
		iterator: iterator,
	}
}

// WithRange returns a Filter which will match all models with a value >=
// start and < limit according to the index.
func (index *Index) WithRange(start []byte, limit []byte) *Filter {
	startWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(start)))
	limitWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(limit)))
	r := &util.Range{Start: startWithPrefix, Limit: limitWithPrefix}
	iterator := index.col.db.ldb.NewIterator(r, nil)
	return &Filter{
		index:    index,
		iterator: iterator,
	}
}

// WithPrefix returns a Filter which will match all models with a value that
// starts with the given prefix according to the index.
func (index *Index) WithPrefix(prefix []byte) *Filter {
	keyPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), escape(prefix)))
	r := util.BytesPrefix(keyPrefix)
	iterator := index.col.db.ldb.NewIterator(r, nil)
	return &Filter{
		index:    index,
		iterator: iterator,
	}
}
