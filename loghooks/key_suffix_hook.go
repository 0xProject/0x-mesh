package loghooks

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var errNestedMapType = errors.New("nested map types are not supported")

// KeySuffixHook is a logger hook that adds suffixes to all keys based on their
// type.
type KeySuffixHook struct{}

// NewKeySuffixHook creates and returns a new KeySuffixHook.
func NewKeySuffixHook() *KeySuffixHook {
	return &KeySuffixHook{}
}

// Ensure that KeySuffixHook implements log.Hook.
var _ log.Hook = &KeySuffixHook{}

func (h *KeySuffixHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *KeySuffixHook) Fire(entry *log.Entry) error {
	newFields := log.Fields{}
	for key, value := range entry.Data {
		typ, err := getTypeForValue(value)
		if err != nil {
			if err == errNestedMapType {
				// We can't safely log nested map types, so replace the value with a
				// string.
				newKey := fmt.Sprintf("%s_json_string", key)
				mapString, err := json.Marshal(value)
				if err != nil {
					return err
				}
				newFields[newKey] = string(mapString)
				continue
			} else {
				return err
			}
		}
		newKey := fmt.Sprintf("%s_%s", key, typ)
		newFields[newKey] = value
	}
	entry.Data = newFields
	return nil
}

// getTypeForValue returns a string representation of the type of the given val.
func getTypeForValue(val interface{}) (string, error) {
	if _, ok := val.(json.Marshaler); ok {
		// If val implements json.Marshaler, return the type of json.Marshal(val)
		// instead of the type of val.
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(val); err != nil {
			return "", err
		}
		var holder interface{}
		if err := json.NewDecoder(buf).Decode(&holder); err != nil {
			return "", err
		}
		return getTypeForValue(holder)
	}
	if _, ok := val.(encoding.TextMarshaler); ok {
		// The json package always encodes values that implement
		// encoding.TextMarshaler as a string.
		return "string", nil
	}
	if _, ok := val.(error); ok {
		// The json package always encodes values that implement
		// error as a string.
		return "string", nil
	}

	underlyingType := getUnderlyingType(reflect.TypeOf(val))
	switch kind := underlyingType.Kind(); kind {
	case reflect.Ptr:
		reflectVal := reflect.ValueOf(val)
		if !reflectVal.IsNil() {
			return getTypeForValue(reflectVal.Elem())
		} else {
			return "null", nil
		}
	case reflect.Bool:
		return "bool", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return "number", nil
	case reflect.String, reflect.Complex64, reflect.Complex128, reflect.Func, reflect.Chan:
		return "string", nil
	case reflect.Array, reflect.Slice:
		return "array", nil
	case reflect.Map:
		// Nested map types can't be efficiently indexed because they allow for
		// arbitrary keys. We don't allow them.
		return "", errNestedMapType
	case reflect.Struct:
		return getSafeStructTypeName(underlyingType)
	default:
		return "", fmt.Errorf("cannot determine type suffix for kind: %s", kind)
	}
}

// getUnderlyingType returns the underlying type for the given type by
// recursively dereferencing pointer types.
func getUnderlyingType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		return getUnderlyingType(typ.Elem())
	}
	return typ
}

// getSafeStructTypeName replaces dots in the name of the given type with
// underscores. Elasticsearch does not allow dots in key names.
func getSafeStructTypeName(typ reflect.Type) (string, error) {
	unsafeTypeName := typ.String()
	safeTypeName := strings.ReplaceAll(unsafeTypeName, ".", "_")
	return safeTypeName, nil
}
