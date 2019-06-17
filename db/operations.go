package db

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func findByID(info *colInfo, reader dbReader, id []byte, model Model) error {
	if err := info.checkModelType(model); err != nil {
		return err
	}
	pk := info.primaryKeyForID(id)
	data, err := reader.Get(pk, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return NotFoundError{ID: id}
		}
		return err
	}
	return json.Unmarshal(data, model)
}

func findAll(info *colInfo, reader dbReader, models interface{}) error {
	prefixRange := util.BytesPrefix(info.prefix())
	iter := reader.NewIterator(prefixRange, nil)
	return findWithIterator(info, iter, models)
}

func findWithIterator(info *colInfo, iter iterator.Iterator, models interface{}) error {
	defer iter.Release()
	if err := info.checkModelsType(models); err != nil {
		return err
	}
	modelsVal := reflect.ValueOf(models).Elem()
	for iter.Next() {
		// We assume that each value in the iterator is the encoded data for some
		// model.
		data := iter.Value()
		model := reflect.New(info.modelType)
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

// findExistingModelByPrimaryKeyWithTransaction gets the latest data for the
// given primary key. Useful in cases where the given model may be out of date
// with what is currently stored in the database. It *doesn't* discard the
// transaction if there is an error.
func findExistingModelByPrimaryKeyWithTransaction(info *colInfo, txn *Transaction, primaryKey []byte) (Model, error) {
	data, err := txn.readerWithBatch.Get(primaryKey, nil)
	if err != nil {
		return nil, err
	}
	// Use reflect to create a new reference for the model type.
	modelRef := reflect.New(info.modelType).Interface()
	if err := json.Unmarshal(data, modelRef); err != nil {
		return nil, err
	}
	model := reflect.ValueOf(modelRef).Elem().Interface().(Model)
	return model, nil
}

func insertWithTransaction(info *colInfo, txn *Transaction, model Model) error {
	if len(model.ID()) == 0 {
		return errors.New("can't insert model with empty ID")
	}
	if err := info.checkModelType(model); err != nil {
		return err
	}
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}
	pk := info.primaryKeyForModel(model)
	if exists, err := txn.readerWithBatch.Has(pk, nil); err != nil {
		return err
	} else if exists {
		return AlreadyExistsError{ID: model.ID()}
	}
	if err := txn.readerWithBatch.Put(pk, data, nil); err != nil {
		return err
	}
	if err := saveIndexesWithTransaction(info, txn, model); err != nil {
		return err
	}
	if err := updateCountWithTransaction(info, txn, 1); err != nil {
		return err
	}
	return nil
}

func updateWithTransaction(info *colInfo, txn *Transaction, model Model) error {
	if len(model.ID()) == 0 {
		return errors.New("can't update model with empty ID")
	}
	if err := info.checkModelType(model); err != nil {
		return err
	}

	// Check if the model already exists and return an error if not.
	pk := info.primaryKeyForModel(model)
	if exists, err := txn.readerWithBatch.Has(pk, nil); err != nil {
		return err
	} else if !exists {
		return NotFoundError{ID: model.ID()}
	}

	// Get the existing data for the model and delete any (now outdated) indexes.
	existingModel, err := findExistingModelByPrimaryKeyWithTransaction(info, txn, pk)
	if err != nil {
		return err
	}
	if err := deleteIndexesWithTransaction(info, txn, existingModel); err != nil {
		return err
	}

	// Save the new data and add the new indexes.
	newData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	if err := txn.readerWithBatch.Put(pk, newData, nil); err != nil {
		return err
	}
	if err := saveIndexesWithTransaction(info, txn, model); err != nil {
		return err
	}
	return nil
}

func deleteWithTransaction(info *colInfo, txn *Transaction, id []byte) error {
	if len(id) == 0 {
		return errors.New("can't delete model with empty ID")
	}

	// We need to get the latest data because the given model might be out of sync
	// with the actual data in the database.
	pk := info.primaryKeyForID(id)
	latest, err := findExistingModelByPrimaryKeyWithTransaction(info, txn, pk)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return NotFoundError{ID: id}
		}
		return err
	}

	// Delete the primary key.
	if err := txn.readerWithBatch.Delete(pk, nil); err != nil {
		return err
	}

	// Delete any index entries.
	if err := deleteIndexesWithTransaction(info, txn, latest); err != nil {
		return err
	}

	// Decrement the model count by 1.
	if err := updateCountWithTransaction(info, txn, -1); err != nil {
		return err
	}

	return nil
}

func saveIndexesWithTransaction(info *colInfo, txn *Transaction, model Model) error {
	info.indexMut.RLock()
	defer info.indexMut.RUnlock()
	for _, index := range info.indexes {
		keys := index.keysForModel(model)
		for _, key := range keys {
			if err := txn.readerWithBatch.Put(key, nil, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

// deleteIndexesForModel deletes any indexes computed from the given model. It
// *doesn't* discard the transaction if there is an error.
func deleteIndexesWithTransaction(info *colInfo, txn *Transaction, model Model) error {
	info.indexMut.RLock()
	defer info.indexMut.RUnlock()
	for _, index := range info.indexes {
		keys := index.keysForModel(model)
		for _, key := range keys {
			if err := txn.readerWithBatch.Delete(key, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func count(info *colInfo, reader dbReader) (int, error) {
	encodedCount, err := reader.Get(info.countKey(), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			// If countKey doesn't exist, assume no models have been inserted and
			// return a count of 0.
			return 0, nil
		}
		return 0, err
	}
	count, err := decodeInt(encodedCount)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func updateCountWithTransaction(info *colInfo, txn *Transaction, diff int) error {
	existingCount, err := count(info, txn.readerWithBatch)
	if err != nil {
		return err
	}
	newCount := existingCount + diff
	if newCount == 0 {
		return txn.readerWithBatch.Delete(info.countKey(), nil)
	} else {
		return txn.readerWithBatch.Put(info.countKey(), encodeInt(newCount), nil)
	}
}

func encodeInt(i int) []byte {
	// TODO(albrow): Could potentially be optimized.
	return []byte(strconv.Itoa(i))
}

func decodeInt(b []byte) (int, error) {
	// TODO(albrow): Could potentially be optimized.
	return strconv.Atoi(string(b))
}
