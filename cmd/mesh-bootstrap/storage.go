// +build !js

package main

import (
	"os"
	"path/filepath"

	"github.com/0xProject/0x-mesh/keys"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq" // postgres driver
)

func getPrivateKeyPath(config Config) string {
	return filepath.Join(config.LevelDBDataDir, "keys", "privkey")
}

func getDHTDir(config Config) string {
	return filepath.Join(config.LevelDBDataDir, "p2p", "dht")
}

func getPeerstoreDir(config Config) string {
	return filepath.Join(config.LevelDBDataDir, "p2p", "peerstore")
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
