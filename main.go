package main

import (
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

type meshEnvVars struct {
	// Logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:VERBOSITY"`
}

func main() {
	// TODO: Main entry point for the 0x Mesh node
	vars := meshEnvVars{}
	if err := envvar.Parse(&vars); err != nil {
		log.Fatal(err)
	}

	// Logging settings
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(vars.Verbosity))
}
