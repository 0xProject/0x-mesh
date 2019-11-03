// +build !js

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0xProject/0x-mesh/keys"
	"github.com/0xProject/0x-mesh/p2p"
	leveldbStore "github.com/ipfs/go-ds-leveldb"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoreds"
	sqlds "github.com/opaolini/go-ds-sql"
	log "github.com/sirupsen/logrus"
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

func getNewDHT(ctx context.Context, config Config, kadDHT *dht.IpfsDHT) func(h host.Host) (routing.PeerRouting, error) {
	switch config.DataStoreType {
	case leveldbDataStore:
		newDHT := func(h host.Host) (routing.PeerRouting, error) {
			var err error
			dhtDir := getDHTDir(config)
			kadDHT, err = p2p.NewDHT(ctx, dhtDir, h)
			if err != nil {
				log.WithField("error", err).Fatal("could not create DHT")
			}
			return kadDHT, err
		}

		return newDHT
	case postgresDataStore:
		newDHT := func(h host.Host) (routing.PeerRouting, error) {
			var err error
			sqlOpts := &sqlds.Options{
				Host:     config.DataDBHost,
				Port:     config.DataDBPort,
				User:     config.DataDBUser,
				Password: config.DataDBPassword,
				Database: config.DataDBDatabaseName,
				Table:    "dhtkv",
			}
			store, err := sqlOpts.CreatePostgres()
			if err != nil {
				log.WithField("error", err).Fatal("could not create postgres datastore")
			}

			kadDHT, err = NewDHTWithDatastore(ctx, store, h)
			if err != nil {
				log.WithField("error", err).Fatal("could not create DHT")
			}

			return kadDHT, err
		}

		return newDHT

	default:
		log.Fatalf("invalid datastore configured: %s. Expected either %s or %s", config.DataStoreType, leveldbDataStore, postgresDataStore)
		return nil

	}
}

func getNewPeerstore(ctx context.Context, config Config) (peerstore.Peerstore, error) {
	switch config.DataStoreType {
	case leveldbDataStore:
		// Set up the peerstore to use LevelDB.
		store, err := leveldbStore.NewDatastore(getPeerstoreDir(config), nil)
		if err != nil {
			return nil, err
		}

		pstore, err := pstoreds.NewPeerstore(ctx, store, pstoreds.DefaultOpts())
		if err != nil {
			return nil, err
		}

		return pstore, nil
	case postgresDataStore:
		sqlOpts := &sqlds.Options{
			Host:     config.DataDBHost,
			Port:     config.DataDBPort,
			User:     config.DataDBUser,
			Password: config.DataDBPassword,
			Database: config.DataDBDatabaseName,
			Table:    "peerStore",
		}
		store, err := sqlOpts.CreatePostgres()
		if err != nil {
			return nil, err
		}

		pstore, err := pstoreds.NewPeerstore(ctx, store, pstoreds.DefaultOpts())
		if err != nil {
			return nil, err
		}

		return pstore, nil
	default:
		return nil, fmt.Errorf("invalid datastore configured: %s. Expected either %s or %s", config.DataStoreType, leveldbDataStore, postgresDataStore)

	}
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
