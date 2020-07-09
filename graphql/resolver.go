package graphql

import (
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/graphql/gqltypes"
)

type resolver struct {
	app *core.App
}

func (r *resolver) Stats() (*gqltypes.Stats, error) {
	stats, err := r.app.GetStats()
	if err != nil {
		return nil, err
	}
	return gqltypes.StatsFromCommonType(stats), nil
}
