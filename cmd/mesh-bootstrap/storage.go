// +build !js

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0xProject/0x-mesh/keys"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq" // postgres driver
)

const (
	dhtTableName       = "dhtkv"
	peerStoreTableName = "peerStore"
)

func getPrivateKeyPath(config Config) string {
	return filepath.Join(config.DataDir, "keys", "privkey")
}

func getDHTDir(config Config) string {
	return filepath.Join(config.DataDir, "p2p", "dht")
}

func getPeerstoreDir(config Config) string {
	return filepath.Join(config.DataDir, "p2p", "peerstore")
}

func getSQLDatabase(config Config) (*sql.DB, error) {
	// Currently we only support the postgres driver.
	if config.SQLDBEngine != "postgres" {
		return nil, errors.New("sqld currently only supports postgres driver")
	}

	if config.SQLDBConnectionString != "" {
		return sql.Open(config.SQLDBEngine, config.SQLDBConnectionString)
	}

	fmtStr := "postgresql:///%s?host=%s&port=%s&user=%s&password=%s&sslmode=disable"
	connstr := fmt.Sprintf(fmtStr, config.SQLDBName, config.SQLDBHost, config.SQLDBPort, config.SQLDBUser, config.SQLDBPassword)

	return sql.Open(config.SQLDBEngine, connstr)
}

func prepareSQLDatabase(db *sql.DB) error {
	createTableString := "CREATE TABLE IF NOT EXISTS %s (key TEXT NOT NULL UNIQUE, data BYTEA NOT NULL)"
	createDHTTable := fmt.Sprintf(createTableString, dhtTableName)
	createPeerStoreTable := fmt.Sprintf(createTableString, peerStoreTableName)

	_, err := db.Exec(createDHTTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createPeerStoreTable)
	if err != nil {
		return err
	}

	return nil
}

func initPrivateKey(path string) (p2pcrypto.PrivKey, error) {
	privKey, err := keys.GetPrivateKeyFromPath(path)
	if err == nil {
		return privKey, nil
	} else if os.IsNotExist(err) {
		// If the private key doesn't exist, generate one.
		log.Info("No private key found. Generating a new one.")
		return keys.GenerateAndSavePrivateKey(path)
	}

	// For any other type of error, return it.
	return nil, err
}
