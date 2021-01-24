package client

import (
	"context"

	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/machinebox/graphql"

	log "github.com/sirupsen/logrus"
)

const (
	orderQueryV4 = `query OrderV4($hash: String!) {
		orderv4(hash: $hash) {
                        hash
			chainId
			exchangeAddress
			makerToken
			takerToken
			makerAmount
			takerAmount
			takerTokenFeeAmount
			maker
			taker
			sender
			feeRecipient
			pool
			expiry
			salt
			signatureType
			signatureV
			signatureR
			signatureS
			fillableTakerAssetAmount
		}
	}`
	ordersQueryV4 = `query OrdersV4($filters: [OrderFilterV4!] = [], $sort: [OrderSortV4!] = [{ field: hash, direction: ASC }], $limit: Int = 100) {
		ordersv4(filters: $filters, sort: $sort, limit: $limit) {
                        hash
			chainId
			exchangeAddress
			makerToken
			takerToken
			makerAmount
			takerAmount
			takerTokenFeeAmount
			maker
			taker
			sender
			feeRecipient
			pool
			expiry
			salt
			signatureType
			signatureV
			signatureR
			signatureS
			fillableTakerAssetAmount
		}
	}`
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
		AddOrders gqltypes.AddOrdersResultsV4 `json:"addOrdersV4"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}

	log.Infof("client results are: %+v", resp.AddOrders)
	returnResult := addOrdersResultsFromGQLTypeV4(&resp.AddOrders)
	return returnResult, nil
}

func (c *Client) GetOrderV4(ctx context.Context, hash common.Hash) (*OrderWithMetadataV4, error) {
	req := graphql.NewRequest(orderQueryV4)
	req.Var("hash", hash.Hex())

	var resp struct {
		Order *gqltypes.OrderV4WithMetadata `json:"orderv4"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}

	log.Infof("getOrder client results are: %+v", resp.Order)
	if resp.Order == nil {
		return nil, nil
	}
	return orderWithMetadataFromGQLTypeV4(resp.Order), nil
}

func (c *Client) FindOrdersV4(ctx context.Context, opts ...FindOrdersOpts) ([]*OrderWithMetadataV4, error) {
	req := graphql.NewRequest(ordersQueryV4)

	if len(opts) > 0 {
		opts := opts[0]
		if len(opts.Filters) > 0 {
			// Convert each filter value from the native Go type to a JSON-compatible type.
			for i, filter := range opts.Filters {
				jsonCompatibleValue, err := gqltypes.FilterValueToJSON(filter)
				if err != nil {
					return nil, err
				}
				opts.Filters[i].Value = jsonCompatibleValue
			}
			req.Var("filters", opts.Filters)
		}
		if len(opts.Sort) > 0 {
			req.Var("sort", opts.Sort)
		}
		if opts.Limit != 0 {
			req.Var("limit", opts.Limit)
		}
	}

	var resp struct {
		Orders []*gqltypes.OrderV4WithMetadata `json:"ordersv4"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}
	return ordersWithMetadataFromGQLTypeV4(resp.Orders), nil
}
