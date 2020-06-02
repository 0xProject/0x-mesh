// +build js,wasm

package keys

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"syscall/js"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
)

// keyPrefix is a prefix applied to all entries in localStorage.
const keyPrefix = "0x-mesh-keys:"

// getKey returns the localStorage key corresponding to the given path.
func getKey(path string) string {
	return keyPrefix + path
}

func readFile(path string) ([]byte, error) {
	if isLocalStorageSupported() {
		return localStorageReadFile(path)
	}
	return ioutil.ReadFile(path)
}

func localStorageReadFile(path string) (data []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = convertRecoverErr(e)
		}
	}()
	key := getKey(path)
	rawData := js.Global().Get("localStorage").Call("getItem", key)
	if jsutil.IsNullOrUndefined(rawData) {
		return nil, os.ErrNotExist
	}
	return base64.StdEncoding.DecodeString(rawData.String())
}

func mkdirAll(dir string) error {
	if isLocalStorageSupported() {
		return localStorageMkdirall(dir)
	}
	return os.MkdirAll(dir, os.ModePerm)
}

func localStorageMkdirall(dir string) (err error) {
	// We don't need to mkdir in localStorage because each path corresponds
	// exactly to one key.
	return nil
}

func writeFile(path string, data []byte) error {
	if isLocalStorageSupported() {
		return localStorageWriteFile(path, data)
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}

func localStorageWriteFile(path string, data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = convertRecoverErr(e)
		}
	}()
	key := getKey(path)
	encodedData := base64.StdEncoding.EncodeToString(data)
	js.Global().Get("localStorage").Call("setItem", key, encodedData)
	return nil
}

// isLocalStorageSupported returns true if localStorage is supported. It does
// this by checking for the global "localStorage" object.
func isLocalStorageSupported() bool {
	return !jsutil.IsNullOrUndefined(js.Global().Get("localStorage"))
}

func convertRecoverErr(e interface{}) error {
	switch e := e.(type) {
	case error:
		return e
	default:
		return fmt.Errorf("recovered with non-error: %v", e)
	}
}
