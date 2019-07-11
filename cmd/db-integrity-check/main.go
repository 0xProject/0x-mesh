// +build !js

// package db-integrity-check is an executable that can be used to check
// the integrity of the database used internally by 0x Mesh.
package main

import (
	"log"

	"github.com/0xProject/0x-mesh/db"
	"github.com/plaid/go-envvar/envvar"
)

type envVars struct {
	// DatabaseDir is the directory where the database files are persisted.
	DatabaseDir string `envvar:"DATABASE_DIR" default:"0x_mesh/db"`
}

func main() {
	env := envVars{}
	if err := envvar.Parse(&env); err != nil {
		log.Fatal(err)
	}
	database, err := db.Open(env.DatabaseDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.CheckIntegrity(); err != nil {
		log.Fatal(err)
	}
	log.Print("Integrity check passed ✓")
}
