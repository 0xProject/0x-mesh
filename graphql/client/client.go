package client

import (
	"context"
	"fmt"

	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/machinebox/graphql"
)

type Request = graphql.Request

// NewRequest can be used to run arbitrary GraphQL queries.
var NewRequest func(q string) *Request = graphql.NewRequest

// Client is a client for the 0x Mesh GraphQL API.
type Client struct {
	*graphql.Client
}

const (
	addOrdersMutation = `mutation AddOrders($orders: [NewOrder!]!, $pinned: Boolean = true) {
		addOrders(orders: $orders, pinned: $pinned) {
			accepted {
				order {
					hash
					chainId
					exchangeAddress
					makerAddress
					makerAssetData
					makerAssetAmount
					makerFeeAssetData
					makerFee
					takerAddress
					takerAssetData
					takerAssetAmount
					takerFeeAssetData
					takerFee
					senderAddress
					feeRecipientAddress
					expirationTimeSeconds
					salt
					signature
					fillableTakerAssetAmount
				}
				isNew
			}
			rejected {
				hash
				code
				message
				order {
					chainId
					exchangeAddress
					makerAddress
					makerAssetData
					makerAssetAmount
					makerFeeAssetData
					makerFee
					takerAddress
					takerAssetData
					takerAssetAmount
					takerFeeAssetData
					takerFee
					senderAddress
					feeRecipientAddress
					expirationTimeSeconds
					salt
					signature
				}
			}
		}
	}`

	ordersQuery = `query Orders($filters: [OrderFilter!] = [], $sort: [OrderSort!] = [{ field: hash, direction: ASC }], $limit: Int = 100) {
		orders(filters: $filters, sort: $sort, limit: $limit) {
			hash
			chainId
			exchangeAddress
			makerAddress
			makerAssetData
			makerAssetAmount
			makerFeeAssetData
			makerFee
			takerAddress
			takerAssetData
			takerAssetAmount
			takerFeeAssetData
			takerFee
			senderAddress
			feeRecipientAddress
			expirationTimeSeconds
			salt
			signature
			fillableTakerAssetAmount
		}
	}`

	orderQuery = `query Order($hash: Hash!) {
		order(hash: $hash) {
			hash
			chainId
			exchangeAddress
			makerAddress
			makerAssetData
			makerAssetAmount
			makerFeeAssetData
			makerFee
			takerAddress
			takerAssetData
			takerAssetAmount
			takerFeeAssetData
			takerFee
			senderAddress
			feeRecipientAddress
			expirationTimeSeconds
			salt
			signature
			fillableTakerAssetAmount
		}
	}`
)

// New creates a new client which points to the given URL.
func New(url string) *Client {
	client := graphql.NewClient(url)
	// TODO(albrow): Remove this.
	client.Log = func(s string) { fmt.Println(s) }
	return &Client{
		Client: client,
	}
}

// AddOrdersOpts is a set of options for the AddOrders method. They can
// be omitted in order to use the defaults.
type AddOrdersOpts struct {
	// Pinned determines whether or not the added orders should be pinned. Pinned
	// orders will not be affected by any DDoS prevention or incentive mechanisms
	// and will always stay in storage until they are no longer fillable. Defaults
	// to true.
	Pinned bool `json:"pinned"`
}

// AddOrders adds orders to 0x Mesh and broadcasts them throughout the 0x Mesh network.
func (c *Client) AddOrders(ctx context.Context, orders []*zeroex.SignedOrder, opts ...AddOrdersOpts) (*AddOrdersResults, error) {
	req := NewRequest(addOrdersMutation)

	// Set up args
	newOrders := gqltypes.NewOrdersFromSignedOrders(orders)
	req.Var("orders", newOrders)

	// Only set the pinned variable if opts were provided.
	if len(opts) > 0 {
		req.Var("pinned", opts[0].Pinned)
	}

	var resp struct {
		AddOrders gqltypes.AddOrdersResults `json:"addOrders"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}
	return addOrdersResultsFromGQLType(&resp.AddOrders), nil
}

func (c *Client) GetOrder(ctx context.Context, hash common.Hash) (*OrderWithMetadata, error) {
	req := NewRequest(orderQuery)
	req.Var("hash", hash.Hex())

	var resp struct {
		Order *gqltypes.OrderWithMetadata `json:"order"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}
	return orderWithMetadataFromGQLType(resp.Order), nil
}

// GetOrdersOpts is a set of options for the GetOrders method. They can
// be omitted in order to use the defaults.
type GetOrdersOpts struct {
	// TODO(albrow): Document fields.
	Filters []OrderFilter
	Sort    []OrderSort
	Limit   int
}

func (c *Client) GetOrders(ctx context.Context, opts ...GetOrdersOpts) ([]*OrderWithMetadata, error) {
	req := NewRequest(ordersQuery)

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
		Orders []*gqltypes.OrderWithMetadata `json:"orders"`
	}
	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}
	return ordersWithMetadataFromGQLType(resp.Orders), nil
}

// func (c *Client) GetStats() (*GetStatsResponse, error)
