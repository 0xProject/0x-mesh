package client

import (
	"context"

	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/machinebox/graphql"

	log "github.com/sirupsen/logrus"
)

// AddOrders adds v4 orders to 0x Mesh and broadcasts them throughout the 0x Mesh network.
func (c *Client) AddOrdersV4(ctx context.Context, orders []*zeroex.SignedOrderV4, opts ...AddOrdersOpts) (*AddOrdersResultsV4, error) {
	req := graphql.NewRequest(addOrdersMutationV4)

	// Set up args
	newOrders := gqltypes.NewOrdersFromSignedOrdersV4(orders)
	req.Var("orders", newOrders)

	// Only set the pinned variable if opts were provided.
	if len(opts) > 0 {
		req.Var("pinned", opts[0].Pinned)
		req.Var("keepCancelled", opts[0].KeepCancelled)
		req.Var("keepExpired", opts[0].KeepExpired)
		req.Var("keepFullyFilled", opts[0].KeepFullyFilled)
		req.Var("keepUnfunded", opts[0].KeepUnfunded)
	}

	var resp struct {
		AddOrders gqltypes.AddOrdersResultsV4 `json:"addOrders"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}

	log.Infof("client results are: %+v", resp.AddOrders)
	returnResult := addOrdersResultsFromGQLTypeV4(&resp.AddOrders)
	return returnResult, nil
}
