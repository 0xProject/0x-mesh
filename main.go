package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	// TODO: Main entry point for the 0x Mesh node
	var (
		verbosity = flag.Int("verbosity", int(log.InfoLevel), "Logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace")
	)
	flag.Parse()

	// Logging settings
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(*verbosity))
}
