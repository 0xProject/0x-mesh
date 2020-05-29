// +build !js

package core

import (
	"context"
	"path/filepath"

	"github.com/0xProject/0x-mesh/db"
)

func newDB(config Config) (*db.DB, error) {
	databasePath := filepath.Join(config.DataDir, "sqlite-db", "db.sqlite")
	// TOOD(albrow): Create and pass through context.
	return db.New(context.TODO(), &db.Options{
		DriverName:     "sqlite3",
		DataSourceName: databasePath,
		MaxOrders:      config.MaxOrdersInStorage,
	})
}
