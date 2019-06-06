package db

import (
	"fmt"
	"reflect"
	"sync"
)

type colInfo struct {
	name      string
	modelType reflect.Type
	indexes   []*Index
	// indexMut protects the indexes slice.
	indexMut sync.RWMutex
}

func (c *colInfo) copy() *colInfo {
	c.indexMut.RLock()
	indexes := make([]*Index, len(c.indexes))
	copy(indexes, c.indexes)
	c.indexMut.RUnlock()
	return &colInfo{
		name:      c.name,
		modelType: c.modelType,
		indexes:   indexes,
	}
}

func (c *colInfo) prefix() []byte {
	return []byte(fmt.Sprintf("model:%s", c.name))
}

func (c *colInfo) primaryKeyForModel(model Model) []byte {
	return c.primaryKeyForID(model.ID())
}

func (c *colInfo) primaryKeyForID(id []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s", c.prefix(), escape(id)))
}

func (c *colInfo) primaryKeyForIDWithoutEscape(id []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s", c.prefix(), id))
}

func (c *colInfo) checkModelType(model Model) error {
	actualType := reflect.TypeOf(model)
	if c.modelType != actualType {
		if actualType.Kind() == reflect.Ptr {
			if c.modelType == actualType.Elem() {
				// Pointers to the expected type are allowed here.
				return nil
			}
		}
		return fmt.Errorf("for %q collection: incorrect type for model (expected %s but got %s)", c.name, c.modelType, actualType)
	}
	return nil
}

func (c *colInfo) checkModelsType(models interface{}) error {
	expectedType := reflect.PtrTo(reflect.SliceOf(c.modelType))
	actualType := reflect.TypeOf(models)
	if expectedType != actualType {
		return fmt.Errorf("for %q collection: incorrect type for models (expected %s but got %s)", c.name, expectedType, actualType)
	}
	return nil
}
