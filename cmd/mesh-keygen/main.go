// +build !js

// mesh-keygen is a short program that can be used to generate private keys.
package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	p2pcrypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/plaid/go-envvar/envvar"
)

type envVars struct {
	// PrivateKeyPath is the path where the private key will be written.
	PrivateKeyPath string `envvar:"PRIVATE_KEY_PATH" default:"0x_mesh/keys/privkey"`
}

func main() {
	env := envVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}
	if _, err := os.Stat(env.PrivateKeyPath); !os.IsNotExist(err) {
		log.Fatalf("Key file: %s already exists. If you really want to overwrite it, delete the file and try again.", env.PrivateKeyPath)
	}
	dir := filepath.Dir(env.PrivateKeyPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	privKey, _, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	keyBytes, err := p2pcrypto.MarshalPrivateKey(privKey)
	if err != nil {
		log.Fatal(err)
	}
	encodedKey := p2pcrypto.ConfigEncodeKey(keyBytes)
	if err := ioutil.WriteFile(env.PrivateKeyPath, []byte(encodedKey), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
