package graphql

import (
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/graphql/data"
	graphql "github.com/graph-gophers/graphql-go"
)

//go:generate go-bindata -pkg data -o ./data/bindata.go ./schema/...
func NewSchema(app *core.App) (*graphql.Schema, error) {
	schema, err := data.Asset("schema/schema.graphql")
	if err != nil {
		return nil, err
	}
	// TODO(albrow): Look into more schema options.
	var opts = []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.UseStringDescriptions()}
	return graphql.ParseSchema(string(schema), &resolver{app: app}, opts...)
}
