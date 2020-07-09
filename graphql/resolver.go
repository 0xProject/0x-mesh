package graphql

import (
	"errors"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/ethereum/go-ethereum/common"
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

type orderArgs struct {
	Hash string
}

func (r *resolver) Order(args orderArgs) (*gqltypes.OrderWithMetadata, error) {
	hash := common.HexToHash(args.Hash)
	if len(hash) == 0 {
		return nil, errors.New("invalid order hash")
	}
	order, err := r.app.GetOrder(hash)
	if err != nil {
		return nil, err
	}
	return gqltypes.OrderWithMetadataFromCommonType(order), nil
}
