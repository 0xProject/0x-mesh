// +build !js

package core

import (
	"context"
	"path/filepath"

	"github.com/0xProject/0x-mesh/db"
)

func newDB(ctx context.Context, config Config) (*db.DB, error) {
	databasePath := filepath.Join(config.DataDir, "sqlite-db", "db.sqlite?_mutex=full&_journal=WAL")
	return db.New(ctx, &db.Options{
		DriverName:     "sqlite3",
		DataSourceName: databasePath,
		MaxOrders:      config.MaxOrdersInStorage,
	})
}
