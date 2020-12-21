package graphql

import (
	"time"

	"github.com/0xProject/0x-mesh/core"
)

type ResolverConfig struct {
	SlowSubscriberTimeout time.Duration
}

type Resolver struct {
	app    *core.App
	config *ResolverConfig
}

func NewResolver(app *core.App, config *ResolverConfig) *Resolver {
	return &Resolver{
		app:    app,
		config: config,
	}
}
