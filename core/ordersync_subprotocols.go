package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/core/ordersync"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/zeroex"
	log "github.com/sirupsen/logrus"
)

// Ensure that FilteredPaginationSubProtocol implements the Subprotocol interface.
var _ ordersync.Subprotocol = (*FilteredPaginationSubProtocol)(nil)

// FilteredPaginationSubProtocol is an ordersync subprotocol which returns all orders by
// paginating through them. It involves sending multiple requests until pagination is
// finished and all orders have been returned.
type FilteredPaginationSubProtocol struct {
	app         *App
	orderFilter *orderfilter.Filter
	perPage     int
}

// NewFilteredPaginationSubprotocol creates and returns a new FilteredPaginationSubprotocol
// which will respond with perPage orders for each individual request/response.
func NewFilteredPaginationSubprotocol(app *App, perPage int) *FilteredPaginationSubProtocol {
	return &FilteredPaginationSubProtocol{
		app:         app,
		orderFilter: app.orderFilter,
		perPage:     perPage,
	}
}

// FilteredPaginationRequestMetadata is the request metadata for the
// FilteredPaginationSubProtocol. It keeps track of the current page and SnapshotID,
// which is expected to be an empty string on the first request.
type FilteredPaginationRequestMetadata struct {
	Page       int    `json:"page"`
	SnapshotID string `json:"snapshotID"`
}

// FilteredPaginationResponseMetadata is the response metadata for the
// FilteredPaginationSubProtocol. It keeps track of the current page and SnapshotID.
type FilteredPaginationResponseMetadata struct {
	Page       int    `json:"page"`
	SnapshotID string `json:"snapshotID"`
}

// Name returns the name of the FilteredPaginationSubProtocol
func (p *FilteredPaginationSubProtocol) Name() string {
	return "/pagination-with-filter/version/0"
}

// HandleOrderSyncRequest returns the orders for one page, based on the page number
// and snapshotID corresponding to the given request. This is
// the implementation for the "provider" side of the subprotocol.
func (p *FilteredPaginationSubProtocol) HandleOrderSyncRequest(ctx context.Context, req *ordersync.Request) (*ordersync.Response, error) {
	var metadata *FilteredPaginationRequestMetadata
	if req.Metadata == nil {
		// Default metadata for the first request.
		metadata = &FilteredPaginationRequestMetadata{
			Page:       0,
			SnapshotID: "",
		}
	} else {
		var ok bool
		metadata, ok = req.Metadata.(*FilteredPaginationRequestMetadata)
		if !ok {
			return nil, fmt.Errorf("FilteredPaginationSubProtocol received request with wrong metadata type (got %T)", req.Metadata)
		}
	}
	// It's possible that none of the orders in the current page match the filter.
	// We don't want to respond with zero orders, so keep iterating until we find
	// at least some orders that match the filter.
	filteredOrders := []*zeroex.SignedOrder{}
	var snapshotID string
	var currentPage int
	for currentPage = metadata.Page; len(filteredOrders) == 0; currentPage += 1 {
		// Get the orders for this page.
		ordersResp, err := p.app.GetOrders(currentPage, p.perPage, metadata.SnapshotID)
		if err != nil {
			return nil, err
		}
		snapshotID = ordersResp.SnapshotID
		if len(ordersResp.OrdersInfos) == 0 {
			// No more orders left.
			break
		}
		// Filter the orders for this page. If none of them match the filter, we continue
		// on to the next page.
		for _, orderInfo := range ordersResp.OrdersInfos {
			if matches, err := p.orderFilter.MatchOrder(orderInfo.SignedOrder); err != nil {
				return nil, err
			} else if matches {
				filteredOrders = append(filteredOrders, orderInfo.SignedOrder)
			}
		}
	}

	return &ordersync.Response{
		Orders:   filteredOrders,
		Complete: len(filteredOrders) == 0,
		Metadata: &FilteredPaginationResponseMetadata{
			Page:       currentPage,
			SnapshotID: snapshotID,
		},
	}, nil
}

// HandleOrderSyncResponse handles the orders for one page by validating them, storing them
// in the database, and firing the appropriate events. It also returns the next request to
// be sent. This is the implementation for the "requester" side of the subprotocol.
func (p *FilteredPaginationSubProtocol) HandleOrderSyncResponse(ctx context.Context, res *ordersync.Response) (*ordersync.Request, error) {
	if res.Metadata == nil {
		return nil, errors.New("FilteredPaginationSubProtocol received response with nil metadata")
	}
	metadata, ok := res.Metadata.(*FilteredPaginationResponseMetadata)
	if !ok {
		return nil, fmt.Errorf("FilteredPaginationSubProtocol received response with wrong metadata type (got %T)", res.Metadata)
	}
	filteredOrders := []*zeroex.SignedOrder{}
	for _, order := range res.Orders {
		if matches, err := p.orderFilter.MatchOrder(order); err != nil {
			return nil, err
		} else if matches {
			filteredOrders = append(filteredOrders, order)
		} else if !matches {
			p.app.handlePeerScoreEvent(res.ProviderID, psReceivedOrderDoesNotMatchFilter)
		}
	}
	validationResults, err := p.app.orderWatcher.ValidateAndStoreValidOrders(ctx, filteredOrders, false, p.app.chainID)
	if err != nil {
		return nil, err
	}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		if acceptedOrderInfo.IsNew {
			log.WithFields(map[string]interface{}{
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
			}).Info("received new valid order from peer via ordersync")
			log.WithFields(map[string]interface{}{
				"order":     acceptedOrderInfo.SignedOrder,
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
			}).Trace("all fields for new valid order received from peer via ordersync")
		}
	}

	return &ordersync.Request{
		Metadata: &FilteredPaginationRequestMetadata{
			Page:       metadata.Page + 1,
			SnapshotID: metadata.SnapshotID,
		},
	}, nil
}

func (p *FilteredPaginationSubProtocol) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationRequestMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (p *FilteredPaginationSubProtocol) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationResponseMetadata
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}
