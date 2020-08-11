package graphql

import "github.com/0xProject/0x-mesh/core"

type Resolver struct {
	app *core.App
}

func NewResolver(app *core.App) *Resolver {
	return &Resolver{
		app: app,
	}
}
