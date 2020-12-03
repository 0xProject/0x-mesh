// +build !js

package core

import (
	"context"
	"path/filepath"

	"github.com/0xProject/0x-mesh/db"
)

func newDB(ctx context.Context, config Config) (*db.DB, error) {
	meshDatabasePath := filepath.Join(config.DataDir, "db", "db.sqlite?_journal=WAL")
	peerStoreDatabasePath := filepath.Join(config.DataDir, "db", "peerstore.sqlite?_journal=WAL")
	dhtDatabasePath := filepath.Join(config.DataDir, "db", "dht.sqlite?_journal=WAL")

	return db.New(ctx, &db.Options{
		DriverName:              "sqlite3",
		DataSourceName:          meshDatabasePath,
		DataSourcePeerStoreName: peerStoreDatabasePath,
		DataSourceDHTName:       dhtDatabasePath,
		MaxOrders:               config.MaxOrdersInStorage,
	})
}
