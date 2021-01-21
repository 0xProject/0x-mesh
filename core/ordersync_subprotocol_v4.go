package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/core/ordersync"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// Ensure that PaginationSubProtocolV4 implements the Subprotocol interface.
var _ ordersync.Subprotocol = (*PaginationSubProtocolV4)(nil)

type PaginationSubProtocolV4 struct {
	app     *App
	perPage int
}

// NewPaginationSubProtocolV4 creates and returns a new PaginationSubProtocolV4
// which will respond with perPage orders for each individual request/response.
func NewPaginationSubProtocolV4(app *App, perPage int) ordersync.Subprotocol {
	return &PaginationSubProtocolV4{
		app:     app,
		perPage: perPage,
	}
}

type PaginationRequestMetadataV4 struct {
	MinOrderHash common.Hash `json:"minOrderHash"`
}

// Name returns the name of the PaginationSubProtocolV4
func (p *PaginationSubProtocolV4) Name() string {
	return "/pagination/version/4"
}

// HandleOrderSyncRequest returns the orders for one page, based on the page number
// and snapshotID corresponding to the given request. This is
// the implementation for the "provider" side of the subprotocol.
func (p *PaginationSubProtocolV4) HandleOrderSyncRequest(ctx context.Context, req *ordersync.Request) (*ordersync.Response, error) {
	var metadata *PaginationRequestMetadataV4
	if req.Metadata == nil {
		// Default metadata for the first request.
		metadata = &PaginationRequestMetadataV4{
			MinOrderHash: common.Hash{},
		}
	} else {
		var ok bool
		metadata, ok = req.Metadata.(*PaginationRequestMetadataV4)
		if !ok {
			return nil, fmt.Errorf("PaginationSubProtocolV4 received request with wrong metadata type (got %T)", req.Metadata)
		}
	}

	orders := []*zeroex.SignedOrderV4{}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	// Get the orders for this page.
	ordersResp, err := p.app.GetOrdersV4(p.perPage, metadata.MinOrderHash)
	if err != nil {
		return nil, err
	}
	for _, orderInfo := range ordersResp.OrdersInfos {
		orders = append(orders, orderInfo.SignedOrderV4)
	}
	return &ordersync.ResponseV4{
		Orders: orders,
	}, nil
}

// HandleOrderSyncResponse handles the orders for one page by validating them, storing them
// in the database, and firing the appropriate events. It also returns the next request to
// be sent. This is the implementation for the "requester" side of the subprotocol.
func (p *PaginationSubProtocolV4) HandleOrderSyncResponse(ctx context.Context, res *ordersync.Response) (*ordersync.Request, int, error) {
	validationResults, err := p.app.orderWatcher.ValidateAndStoreValidOrders(ctx, res.Orders, p.app.chainID, false, &types.AddOrdersOpts{})
	if err != nil {
		return nil, len(res.Orders), err
	}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		if acceptedOrderInfo.IsNew {
			log.WithFields(map[string]interface{}{
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
				"protocol":  "ordersync",
			}).Info("received new valid order from peer")
			log.WithFields(map[string]interface{}{
				"order":     acceptedOrderInfo.SignedOrderV4,
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
				"protocol":  "ordersync",
			}).Trace("all fields for new valid order received from peer")
		}
	}

	// Calculate the next min order hash to send in our next request.
	// This is equal to the maximum order hash we have received so far.
	var nextMinOrderHash common.Hash
	if len(res.Orders) > 0 {
		hash, err := res.Orders[len(res.Orders)-1].ComputeOrderHash()
		if err != nil {
			return nil, len(res.Orders), err
		}
		nextMinOrderHash = hash
	}
	return &ordersync.Request{
		Metadata: &PaginationRequestMetadataV4{
			MinOrderHash: nextMinOrderHash,
		},
	}, len(res.Orders), nil
}

func (p *PaginationSubProtocolV4) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationRequestMetadataV0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (p *PaginationSubProtocolV4) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationResponseMetadataV0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (p *PaginationSubProtocolV4) GenerateFirstRequestMetadata() (json.RawMessage, error) {
	return json.Marshal(FilteredPaginationRequestMetadataV0{
		OrderFilter: p.orderFilter,
		Page:        0,
		SnapshotID:  "",
	})
}
