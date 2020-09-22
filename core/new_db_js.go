// +build js,wasm

package core

import (
	"context"
	"path/filepath"

	"github.com/0xProject/0x-mesh/db"
)

func newDB(ctx context.Context, config Config) (*db.DB, error) {
	databasePath := filepath.Join(config.DataDir, "mesh_dexie_db")
	return db.New(ctx, &db.Options{
		DriverName:     "dexie",
		DataSourceName: databasePath,
		MaxOrders:      config.MaxOrdersInStorage,
	})
}
