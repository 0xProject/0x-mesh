// +build js,wasm

package keys

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"syscall/js"
)

func readFile(path string) ([]byte, error) {
	if isBrowserFSSupported() {
		return browserFSReadFile(path)
	}
	return ioutil.ReadFile(path)
}

func browserFSReadFile(path string) (data []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			if jsErr, ok := e.(js.Error); ok {
				err = convertJSError(jsErr)
			}
		}
	}()
	options := map[string]interface{}{
		"encoding": "hex",
	}
	rawData := js.Global().Get("browserFS").Call("readFileSync", path, options)
	return hex.DecodeString(rawData.String())
}

func mkdirAll(dir string) error {
	if isBrowserFSSupported() {
		return browserFSMkdirall(dir)
	}
	return os.MkdirAll(dir, os.ModePerm)
}

func browserFSMkdirall(dir string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if jsErr, ok := e.(js.Error); ok {
				err = convertJSError(jsErr)
			}
		}
	}()
	// Note: mkdirAll is not supported by BrowserFS so we have to manually create
	// each directory.
	names := strings.Split(dir, string(os.PathSeparator))
	for i := range names {
		partialPath := filepath.Join(names[:i+1]...)
		if err := browserFSMkdir(partialPath); err != nil {
			if os.IsExist(err) {
				// If the directory already exists, that's fine.
				continue
			}
		}
	}
	return nil
}

func browserFSMkdir(dir string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if jsErr, ok := e.(js.Error); ok {
				err = convertJSError(jsErr)
			}
		}
	}()
	js.Global().Get("browserFS").Call("mkdirSync", dir, int(os.ModePerm))
	return nil
}

func writeFile(path string, data []byte) error {
	if isBrowserFSSupported() {
		return browserFSWriteFile(path, data)
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}

func browserFSWriteFile(path string, data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if jsErr, ok := e.(js.Error); ok {
				err = convertJSError(jsErr)
			}
		}
	}()
	// The naive approach of using `string(data)` for the data to write doesn't
	// work, regardless of the encoding used. Encoding to hex seems like the most
	// reliable way to do it.
	options := map[string]interface{}{
		"encoding": "hex",
	}
	js.Global().Get("browserFS").Call("writeFileSync", path, hex.EncodeToString(data), options)
	return nil
}

// isBrowserFSSupported returns true if BrowserFS is supported. It does this by
// checking for the global "browserFS" object.
func isBrowserFSSupported() bool {
	return js.Global().Get("browserFS") != js.Null() && js.Global().Get("browserFS") != js.Undefined()
}

// convertJSError converts an error returned by the BrowserFS API into a Go
// error. This is important because Go expects certain types of errors to be
// returned (e.g. ENOENT when a file doesn't exist) and programs often change
// their behavior depending on the type of error.
func convertJSError(err js.Error) error {
	if err.Value == js.Undefined() || err.Value == js.Null() {
		return nil
	}
	// TODO(albrow): Convert to os.PathError when possible/appropriate.
	if code := err.Get("code"); code != js.Undefined() && code != js.Null() {
		switch code.String() {
		case "ENOENT":
			return os.ErrNotExist
		case "EISDIR":
			return syscall.EISDIR
		case "EEXIST":
			return os.ErrExist
			// TODO(albrow): Fill in more codes here.
		}
	}
	return err
}
