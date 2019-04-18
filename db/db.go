package db

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/syndtr/goleveldb/leveldb"
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
	db   *DB
	name string
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
		return err
	} else if exists {
		return fmt.Errorf("%s model with given ID already exists in database: %s", c.name, hex.Dump(model.ID()))
	}
	if err := txn.Put(pk, data, nil); err != nil {
		return err
	}

	// TODO(albrow): Add/update indexes.
	return txn.Commit()
}

func (c *Collection) FindByID(id []byte, model Model) error {
	pk := c.primaryKeyForID(id)
	data, err := c.db.ldb.Get(pk, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, model)
}

func (c *Collection) FindAll(models interface{}) error {
	modelType, err := getModelTypeFromSlice(models)
	if err != nil {
		return err
	}
	modelsVal := reflect.ValueOf(models).Elem()
	prefixRange := util.BytesPrefix(c.prefix())
	iter := c.db.ldb.NewIterator(prefixRange, nil)
	defer iter.Release()
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
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
