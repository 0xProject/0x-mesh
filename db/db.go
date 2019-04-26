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

// Note about the implementation:
//
// There are two types of keys used. A "primary key" is the main key for a
// particular model. It's value is the encoded data for that model. The format
// for a primary key is: `model:<collection name>:<model ID>`.
//
// An "index key" is used in the FindWithValue and FindWithRange methods to find
// models with specific indexed values. The format for an index key is:
// `index:<collection name>:<index name>:<value>:<model ID>`. Unlike primary
// keys, index keys have no values and don't store any actual data. Instead, the
// primary key can be extracted from an index key and then used to look up the
// data for the corresponding model.

// Model is any type which can be inserted and retrieved from the database. The
// only requirement is an ID method.
type Model interface {
	// ID returns a unique identifier for this model.
	ID() []byte
}

// DB is the top-level Database.
type DB struct {
	ldb *leveldb.DB
}

// Open creates a new database using the given file path for permanent storage.
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

// Collection represents a set of a specific type of model.
type Collection struct {
	db        *DB
	name      string
	modelType reflect.Type
	indexes   []*Index
}

// NewCollection creates and returns a new collection with the given name and
// model type. You should create exactly one collection for each model type. The
// collection should typically be created once at the start of your application
// and re-used.
func (db *DB) NewCollection(name string, typ Model) *Collection {
	return &Collection{
		db:        db,
		name:      name,
		modelType: reflect.TypeOf(typ),
	}
}

// Name returns the name of the collection.
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

func (c *Collection) checkModelType(model Model) error {
	actualType := reflect.TypeOf(model)
	if c.modelType != actualType {
		if actualType.Kind() == reflect.Ptr {
			if c.modelType == actualType.Elem() {
				// Pointers to the expected type are allowed here.
				return nil
			}
		}
		return fmt.Errorf("for %q collection: incorrect type for model (expected %s but got %s)", c.Name(), c.modelType, actualType)
	}
	return nil
}

func (c *Collection) checkModelsType(models interface{}) error {
	expectedType := reflect.PtrTo(reflect.SliceOf(c.modelType))
	actualType := reflect.TypeOf(models)
	if expectedType != actualType {
		return fmt.Errorf("for %q collection: incorrect type for models (expected %s but got %s)", c.Name(), expectedType, actualType)
	}
	return nil
}

// Insert inserts the given model into the database. It returns an error if a
// model with the same id already exists.
func (c *Collection) Insert(model Model) error {
	if len(model.ID()) == 0 {
		return errors.New("can't insert model with empty ID")
	}
	if err := c.checkModelType(model); err != nil {
		return err
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

// Update updates an existing model in the database. It returns an error if the
// given model doesn't already exist.
func (c *Collection) Update(model Model) error {
	if len(model.ID()) == 0 {
		return errors.New("can't update model with empty ID")
	}
	if err := c.checkModelType(model); err != nil {
		return err
	}
	newData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	txn, err := c.db.ldb.OpenTransaction()
	if err != nil {
		return err
	}
	pk := c.primaryKeyForModel(model)

	// Check if the model already exists and return an error if not.
	if exists, err := txn.Has(pk, nil); err != nil {
		txn.Discard()
		return err
	} else if !exists {
		txn.Discard()
		return fmt.Errorf("can't update %s model because ID doesn't exist in database: %s", c.name, hex.Dump(model.ID()))
	}

	// Get the existing data for the model and delete any (now outdated) indexes.
	existingModel, err := c.findExistingModelByPrimaryKey(txn, pk)
	if err != nil {
		txn.Discard()
		return err
	}
	if err := c.deleteIndexesForModel(txn, existingModel); err != nil {
		txn.Discard()
		return err
	}

	// Save the new data and add the new indexes.
	if err := txn.Put(pk, newData, nil); err != nil {
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

// FindByID finds the model with the given ID and scans the results into the
// given model. As in the Unmarshal and Decode methods in the encoding/json
// package, model must be settable via reflect. Typically, this means you should
// pass in a pointer.
func (c *Collection) FindByID(id []byte, model Model) error {
	pk := c.primaryKeyForID(id)
	return c.findByKey(pk, model)
}

func (c *Collection) findByKey(key []byte, model Model) error {
	if err := c.checkModelType(model); err != nil {
		return err
	}
	data, err := c.db.ldb.Get(key, nil)
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New("model not found")
	}
	return json.Unmarshal(data, model)
}

// FindAll finds all models for the collection and scans the results into the
// given models. models should be a pointer to an empty slice of a concrete
// model type (e.g. *[]myModelType).
func (c *Collection) FindAll(models interface{}) error {
	prefixRange := util.BytesPrefix(c.prefix())
	iter := c.db.ldb.NewIterator(prefixRange, nil)
	return c.findWithIterator(iter, models)
}

func (c *Collection) findWithIterator(iter iterator.Iterator, models interface{}) error {
	defer iter.Release()
	if err := c.checkModelsType(models); err != nil {
		return err
	}
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		// We assume that each value in the iterator is the encoded data for some
		// model.
		data := iter.Value()
		model := reflect.New(c.modelType)
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

// Delete deletes the model with the given ID from the database. It returns an
// error if the model doesn't exist in the database.
func (c *Collection) Delete(id []byte) error {
	if len(id) == 0 {
		return errors.New("can't delete model with empty ID")
	}
	txn, err := c.db.ldb.OpenTransaction()
	if err != nil {
		return err
	}

	// Check if the model already exists and return an error if not.
	pk := c.primaryKeyForID(id)
	if exists, err := txn.Has(pk, nil); err != nil {
		txn.Discard()
		return err
	} else if !exists {
		txn.Discard()
		return fmt.Errorf("can't delete %s model because ID doesn't exist in database: %s", c.name, hex.Dump(id))
	}

	// We need to get the latest data because the given model might be out of sync
	// with the actual data in the database.
	latest, err := c.findExistingModelByPrimaryKey(txn, pk)
	if err != nil {
		txn.Discard()
		return err
	}

	// Delete the primary key.
	if err := txn.Delete(pk, nil); err != nil {
		txn.Discard()
		return err
	}

	// Delete any index entries.
	if err := c.deleteIndexesForModel(txn, latest); err != nil {
		txn.Discard()
		return err
	}

	return txn.Commit()
}

// deleteIndexesForModel deletes any indexes computed from the given model. It
// *doesn't* discard the transaction if there is an error.
func (c *Collection) deleteIndexesForModel(txn *leveldb.Transaction, model Model) error {
	for _, index := range c.indexes {
		key := index.keyForModel(model)
		if err := txn.Delete(key, nil); err != nil {
			return err
		}
	}
	return nil
}

// findExistingModelByPrimaryKey gets the latest data for the given primary key.
// Useful in cases where the given model may be out of date with what is
// currently stored in the database. It *doesn't* discard the transaction if
// there is an error.
func (c *Collection) findExistingModelByPrimaryKey(txn *leveldb.Transaction, primaryKey []byte) (Model, error) {
	data, err := txn.Get(primaryKey, nil)
	if err != nil {
		return nil, err
	}
	// Use reflect to create a new reference for the model type.
	modelRef := reflect.New(c.modelType).Interface()
	if err := json.Unmarshal(data, modelRef); err != nil {
		return nil, err
	}
	model := reflect.ValueOf(modelRef).Elem().Interface().(Model)
	return model, nil
}

// Index can be used to search for specific values or specific ranges of values
// for a collection.
type Index struct {
	col    *Collection
	name   string
	getter func(m Model) []byte
}

// AddIndex creates and returns a new index. name is an arbitrary, unique name
// for the index. getter is a function that accepts a model and returns the
// value for this particular index. For example, if you wanted to add an index
// on a struct field, getter should return the value of that field. After
// AddIndex is called, any new models in this collection that are inserted will
// be indexed. Any models inserted prior to calling AddIndex will *not* be
// indexed.
func (c *Collection) AddIndex(name string, getter func(Model) []byte) *Index {
	index := &Index{
		col:    c,
		name:   name,
		getter: getter,
	}
	c.indexes = append(c.indexes, index)
	return index
}

// Name returns the name of the index.
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

// FindWithValue finds all models with the given value according to the index
// and scans the results into models. models should be a pointer to an empty
// slice of a concrete model type (e.g. *[]myModelType).
func (c *Collection) FindWithValue(index *Index, val []byte, models interface{}) error {
	prefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), val))
	prefixRange := util.BytesPrefix(prefix)
	iter := c.db.ldb.NewIterator(prefixRange, nil)
	return c.findWithIndexIterator(index, iter, models)
}

// FindWithRange finds all models with a value >= start and < limit according to
// the index and scans the results into models. models should be a pointer to an
// empty slice of a concrete model type (e.g. *[]myModelType).
func (c *Collection) FindWithRange(index *Index, start []byte, limit []byte, models interface{}) error {
	startWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), start))
	limitWithPrefix := []byte(fmt.Sprintf("%s:%s", index.prefix(), limit))
	r := &util.Range{Start: startWithPrefix, Limit: limitWithPrefix}
	iter := c.db.ldb.NewIterator(r, nil)
	return c.findWithIndexIterator(index, iter, models)
}

func (c *Collection) findWithIndexIterator(index *Index, iter iterator.Iterator, models interface{}) error {
	defer iter.Release()
	if err := c.checkModelsType(models); err != nil {
		return err
	}
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		// We assume that each key in the iterator consists of an index prefix, the
		// value for a particular model, and the model ID. We can extract a primary
		// key from this key and use it to get the encoded data for the model
		// itself.
		pk := index.primaryKeyFromKey(iter.Key())
		data, err := c.db.ldb.Get(pk, nil)
		if err != nil {
			return err
		}
		model := reflect.New(c.modelType)
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
