package db

import (
	"fmt"
	"reflect"
	"sync"
)

// colInfo is a set of information/metadata about a collection.
type colInfo struct {
	db        *DB
	name      string
	modelType reflect.Type
	indexes   []*Index
	// indexMut protects the indexes slice.
	indexMut sync.RWMutex
	// writeMut is used by transactions to prevent other goroutines from writing
	// until the transaction is committed or discarded. Needs to be a pointer so
	// that copies of this colInfo retain the same writeLock.
	writeMut *sync.Mutex
}

// copy returns a copy of the colInfo. Any changes made to the original (e.g.
// adding a new index) will not affect the copy. The copy and the original share
// the same writeMut.
func (info *colInfo) copy() *colInfo {
	info.indexMut.RLock()
	indexes := make([]*Index, len(info.indexes))
	copy(indexes, info.indexes)
	info.indexMut.RUnlock()
	return &colInfo{
		db:        info.db,
		name:      info.name,
		modelType: info.modelType,
		indexes:   indexes,
		writeMut:  info.writeMut,
	}
}

func (info *colInfo) prefix() []byte {
	return []byte(fmt.Sprintf("model:%s", escape([]byte(info.name))))
}

// countKey returns the key used to store a count of the number of models in the
// collection.
func (info *colInfo) countKey() []byte {
	return []byte(fmt.Sprintf("count:%s", escape([]byte(info.name))))
}

func (info *colInfo) primaryKeyForModel(model Model) []byte {
	return info.primaryKeyForID(model.ID())
}

func (info *colInfo) primaryKeyForID(id []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s", info.prefix(), escape(id)))
}

func (info *colInfo) primaryKeyForIDWithoutEscape(id []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s", info.prefix(), id))
}

func (info *colInfo) checkModelType(model Model) error {
	actualType := reflect.TypeOf(model)
	if info.modelType != actualType {
		if actualType.Kind() == reflect.Ptr {
			if info.modelType == actualType.Elem() {
				// Pointers to the expected type are allowed here.
				return nil
			}
		}
		return fmt.Errorf("for %q collection: incorrect type for model (expected %s but got %s)", info.name, info.modelType, actualType)
	}
	return nil
}

func (info *colInfo) checkModelsType(models interface{}) error {
	expectedType := reflect.PtrTo(reflect.SliceOf(info.modelType))
	actualType := reflect.TypeOf(models)
	if expectedType != actualType {
		return fmt.Errorf("for %q collection: incorrect type for models (expected %s but got %s)", info.name, expectedType, actualType)
	}
	return nil
}
