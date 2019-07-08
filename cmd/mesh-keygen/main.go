// +build !js

// mesh-keygen is a short program that can be used to generate private keys.
package main

import (
	"log"
	"os"

	"github.com/0xProject/0x-mesh/keys"
	"github.com/plaid/go-envvar/envvar"
)

type envVars struct {
	// PrivateKeyPath is the path where the private key will be written.
	PrivateKeyPath string `envvar:"PRIVATE_KEY_PATH" default:"0x_mesh/keys/privkey"`
}

func main() {
	env := envVars{}
	if err := envvar.Parse(&env); err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(env.PrivateKeyPath); !os.IsNotExist(err) {
		log.Fatalf("Key file: %s already exists. If you really want to overwrite it, delete the file and try again.", env.PrivateKeyPath)
	}
	if _, err := keys.GenerateAndSavePrivateKey(env.PrivateKeyPath); err != nil {
		log.Fatal(err)
	}
}
