package db

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (db *DB) CheckIntegrity() error {
	db.colLock.Lock()
	defer db.colLock.Unlock()
	for _, col := range db.collections {
		if err := db.checkCollectionIntegrity(col); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) checkCollectionIntegrity(col *Collection) error {
	col.info.indexMut.RLock()
	defer col.info.indexMut.RUnlock()

	snapshot, err := col.GetSnapshot()
	if err != nil {
		return err
	}
	defer snapshot.Release()

	slice := util.BytesPrefix([]byte(fmt.Sprintf("%s:", col.info.prefix())))
	iter := snapshot.snapshot.NewIterator(slice, nil)
	defer iter.Release()
	for iter.Next() {
		// Check that the model data can be unmarshaled into the expected type.
		data := iter.Value()
		modelVal := reflect.New(col.info.modelType)
		if err := json.Unmarshal(data, modelVal.Interface()); err != nil {
			return fmt.Errorf("integritiy check failed for collection %s: could not unmarshal model data for primary key %s: %s", col.Name(), iter.Key(), err.Error())
		}
		model := modelVal.Elem().Interface().(Model)

		// Check that the index entries exist for this model.
		for _, index := range col.info.indexes {
			indexKeys := index.keysForModel(model)
			for _, indexKey := range indexKeys {
				indexKeyExists, err := snapshot.snapshot.Has(indexKey, nil)
				if err != nil {
					return err
				}
				if !indexKeyExists {
					return fmt.Errorf("integritiy check failed for index %s.%s: indexKey %s does not exist", col.Name(), index.Name(), indexKey)
				}
			}
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}

	// Check the integrity of each index.
	for _, index := range col.info.indexes {
		if err := db.checkIndexIntegrity(snapshot, col, index); err != nil {
			return err
		}
	}

	return nil
}

// checkIndexIntegrity checks that each key in the index corresponds to model
// data that exists and is valid (can be unmarshaled into a model of the
// expected type).
func (db *DB) checkIndexIntegrity(snapshot *Snapshot, col *Collection, index *Index) error {
	slice := util.BytesPrefix([]byte(fmt.Sprintf("%s:", index.prefix())))
	iter := snapshot.snapshot.NewIterator(slice, nil)
	defer iter.Release()
	for iter.Next() {
		pk := index.primaryKeyFromIndexKey(iter.Key())
		data, err := snapshot.snapshot.Get(pk, nil)
		if err != nil {
			if err == leveldb.ErrNotFound {
				return fmt.Errorf("integritiy check failed for index %s.%s: key exists in index but could not find corresponding model data for primary key: %s", col.Name(), index.Name(), pk)
			} else {
				return err
			}
		}
		modelVal := reflect.New(col.info.modelType)
		if err := json.Unmarshal(data, modelVal.Interface()); err != nil {
			return fmt.Errorf("integritiy check failed for index %s.%s: could not unmarshal model data: %s", col.Name(), index.Name(), err.Error())
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}
