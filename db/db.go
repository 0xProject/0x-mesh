package db

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Model interface {
	ID() []byte
}

type DB struct {
	ldb *leveldb.DB
}

func Open(path string) (*DB, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		ldb: ldb,
	}, nil
}

func (db *DB) Close() error {
	return db.ldb.Close()
}

type Collection struct {
	db      *DB
	name    string
	indexes []*Index
}

func (db *DB) NewCollection(name string) *Collection {
	return &Collection{
		db:   db,
		name: name,
	}
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) prefix() []byte {
	return []byte(fmt.Sprintf("model:%s", c.name))
}

func (c *Collection) primaryKeyForModel(model Model) []byte {
	return c.primaryKeyForID(model.ID())
}

func (c *Collection) primaryKeyForID(id []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s", c.prefix(), id))
}

func (c *Collection) Insert(model Model) error {
	if len(model.ID()) == 0 {
		return errors.New("can't insert model with empty ID")
	}
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}
	txn, err := c.db.ldb.OpenTransaction()
	if err != nil {
		return err
	}
	pk := c.primaryKeyForModel(model)
	if exists, err := txn.Has(pk, nil); err != nil {
		txn.Discard()
		return err
	} else if exists {
		txn.Discard()
		return fmt.Errorf("%s model with given ID already exists in database: %s", c.name, hex.Dump(model.ID()))
	}
	if err := txn.Put(pk, data, nil); err != nil {
		txn.Discard()
		return err
	}
	for _, index := range c.indexes {
		key := index.keyForModel(model)
		if err := txn.Put(key, nil, nil); err != nil {
			txn.Discard()
			return err
		}
	}

	return txn.Commit()
}

func (c *Collection) FindByID(id []byte, model Model) error {
	pk := c.primaryKeyForID(id)
	return c.findByKey(pk, model)
}

func (c *Collection) findByKey(key []byte, model Model) error {
	data, err := c.db.ldb.Get(key, nil)
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New("model not found")
	}
	return json.Unmarshal(data, model)
}

func (c *Collection) FindAll(models interface{}) error {
	prefixRange := util.BytesPrefix(c.prefix())
	iter := c.db.ldb.NewIterator(prefixRange, nil)
	return c.findWithIterator(iter, models)
}

func (c *Collection) findWithIterator(iter iterator.Iterator, models interface{}) error {
	defer iter.Release()
	modelType, err := getModelTypeFromSlice(models)
	if err != nil {
		return err
	}
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		data := iter.Value()
		model := reflect.New(modelType)
		if err := json.Unmarshal(data, model.Interface()); err != nil {
			return err
		}
		modelsVal.Set(reflect.Append(modelsVal, model.Elem()))
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}

func (c *Collection) findWithIndexIterator(index *Index, iter iterator.Iterator, models interface{}) error {
	defer iter.Release()
	modelType, err := getModelTypeFromSlice(models)
	if err != nil {
		return err
	}
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		pk := index.primaryKeyFromKey(iter.Key())
		data, err := c.db.ldb.Get(pk, nil)
		if err != nil {
			return err
		}
		model := reflect.New(modelType)
		if err := json.Unmarshal(data, model.Interface()); err != nil {
			return err
		}
		modelsVal.Set(reflect.Append(modelsVal, model.Elem()))
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Delete(model Model) error {
	if len(model.ID()) == 0 {
		return errors.New("can't delete model with empty ID")
	}
	txn, err := c.db.ldb.OpenTransaction()
	if err != nil {
		return err
	}

	// Get the latest data for the Model. Required because the given model might
	// be out of sync with the actual data in the database.
	pk := c.primaryKeyForModel(model)
	data, err := txn.Get(pk, nil)
	if err != nil {
		txn.Discard()
		return err
	}
	// TODO(albrow): Be more safe here. Handle pointers and non-pointers.
	updatedRef := reflect.New(reflect.TypeOf(model)).Interface().(Model)
	if err := json.Unmarshal(data, updatedRef); err != nil {
		txn.Discard()
		return err
	}
	updated := reflect.ValueOf(updatedRef).Elem().Interface().(Model)

	// Delete the primary key.
	if err := txn.Delete(pk, nil); err != nil {
		txn.Discard()
		return err
	}

	// Delete any index entries.
	for _, index := range c.indexes {
		key := index.keyForModel(updated)
		if err := txn.Delete(key, nil); err != nil {
			txn.Discard()
			return err
		}
	}

	return txn.Commit()
}

type Index struct {
	col    *Collection
	name   string
	getter func(m Model) []byte
}

func (c *Collection) AddIndex(name string, getter func(Model) []byte) *Index {
	index := &Index{
		col:    c,
		name:   name,
		getter: getter,
	}
	c.indexes = append(c.indexes, index)
	return index
}

func (index *Index) Name() string {
	return index.name
}

func (index *Index) prefix() []byte {
	return []byte(fmt.Sprintf("index:%s:%s", index.col.name, index.name))
}

func (index *Index) keyForModel(model Model) []byte {
	value := index.getter(model)
	return []byte(fmt.Sprintf("%s:%s:%s", index.prefix(), value, model.ID()))
}

// primaryKeyFromKey extracts and returns the primary key from the given index
// key.
func (index *Index) primaryKeyFromKey(key []byte) []byte {
	pkAndVal := strings.TrimPrefix(string(key), string(index.prefix()))
	split := strings.Split(pkAndVal, ":")
	return index.col.primaryKeyForID([]byte(split[2]))
}

func (c *Collection) FindWithValue(index *Index, val []byte, models interface{}) error {
	prefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), val))
	prefixRange := util.BytesPrefix(prefix)
	iter := c.db.ldb.NewIterator(prefixRange, nil)
	return c.findWithIndexIterator(index, iter, models)
}

func (c *Collection) FindWithRange(index *Index, start []byte, limit []byte, models interface{}) error {
	startWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), start))
	limitWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), limit))
	r := &util.Range{Start: startWithPrefix, Limit: limitWithPrefix}
	iter := c.db.ldb.NewIterator(r, nil)
	return c.findWithIndexIterator(index, iter, models)
}

var modelInterfaceType = reflect.TypeOf([]Model{}).Elem()

func getModelTypeFromSlice(models interface{}) (reflect.Type, error) {
	ptrType := reflect.TypeOf(models)
	if ptrType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("type should be a pointer to a slice of models but got: %T (not a pointer)", models)
	}
	sliceType := ptrType.Elem()
	if sliceType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("type should be a pointer to a slice of models but got: %T (not a slice)", models)
	}
	elemType := sliceType.Elem()
	if !elemType.Implements(modelInterfaceType) {
		return nil, fmt.Errorf("type should be a pointer to a slice of models but got: %T (element type doesn't implement Model)", models)
	}
	return elemType, nil
}
