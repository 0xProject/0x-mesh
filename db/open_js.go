// +build js,wasm

package db

import (
	"errors"
	"syscall/js"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

const (
	// browserFSLoadCheckInterval is frequently to check whether browserFS is
	// loaded.
	browserFSLoadCheckInterval = 50 * time.Millisecond
	// browserFSLoadTimeout is how long to wait for BrowserFS to finish loading
	// before giving up.
	browserFSLoadTimeout = 5 * time.Second
)

// Open creates a new database for js/wasm environments.
func Open(path string) (*DB, error) {
	// The global willLoadBrowserFS variable indicates whether browserFS will be
	// loaded. browserFS has to be explicitly loaded in by JavaScript (and
	// typically Webpack) and can't be loaded here.
	if willLoadBrowserFS := js.Global().Get("willLoadBrowserFS"); !willLoadBrowserFS.Equal(js.Undefined()) && willLoadBrowserFS.Bool() == true {
		return openBrowserFSDB(path)
	}
	// If browserFS is not going to be loaded, fallback to using an in-memory
	// database.
	return openInMemoryDB()
}

func openInMemoryDB() (*DB, error) {
	log.Warn("BrowserFS not detected. Using in-memory databse.")
	ldb, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		ldb: ldb,
	}, nil
}

func openBrowserFSDB(path string) (*DB, error) {
	log.Info("BrowserFS detected. Using BrowserFS-backed databse.")
	// Wait for browserFS to load.
	//
	// HACK(albrow): We do this by checking for the global browserFS
	// variable. This is definitely a bit of a hack and wastes some CPU resources,
	// but it is also extremely reliable. Given that we have a chicken and egg
	// problem with both Wasm and JavaScript code loading and executing at the
	// same time, it is difficult to match this level of reliability with something
	// like callback functions or events.
	start := time.Now()
	for {
		if time.Since(start) >= browserFSLoadTimeout {
			return nil, errors.New("timed out waiting for BrowserFS to load")
		}
		if !js.Global().Get("browserFS").Equal(js.Undefined()) && !js.Global().Get("browserFS").Equal(js.Null()) {
			log.Info("BrowserFS finished loading")
			break
		}
		time.Sleep(browserFSLoadCheckInterval)
	}
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		ldb: ldb,
	}, nil
}
