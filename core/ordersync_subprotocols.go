package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/core/ordersync"
	"github.com/0xProject/0x-mesh/zeroex"
)

var _ ordersync.Subprotocol = (*FilteredPaginationSubProtocol)(nil)

type FilteredPaginationSubProtocol struct {
	app     *App
	perPage int
}

func NewFilteredPaginationSubprotocol(app *App, perPage int) *FilteredPaginationSubProtocol {
	return &FilteredPaginationSubProtocol{
		app:     app,
		perPage: perPage,
	}
}

type FilteredPaginationRequestMetadata struct {
	Page       int    `json:"page"`
	SnapshotID string `json:"snapshotID"`
}

type FilteredPaginationResponseMetadata struct {
	Page       int    `json:"page"`
	SnapshotID string `json:"snapshotID"`
}

func (p *FilteredPaginationSubProtocol) Name() string {
	return "/pagination-with-filter/version/0"
}

func (p *FilteredPaginationSubProtocol) GetOrders(ctx context.Context, req *ordersync.Request) (*ordersync.Response, error) {
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

	ordersResp, err := p.app.GetOrders(metadata.Page, p.perPage, metadata.SnapshotID)
	if err != nil {
		return nil, err
	}
	orders := make([]*zeroex.SignedOrder, len(ordersResp.OrdersInfos))
	for i, orderInfo := range ordersResp.OrdersInfos {
		orders[i] = orderInfo.SignedOrder
	}
	// TODO(albrow): Filter orders
	return &ordersync.Response{
		Orders:   orders,
		Complete: len(orders) == 0,
		Metadata: &FilteredPaginationResponseMetadata{
			Page:       metadata.Page,
			SnapshotID: ordersResp.SnapshotID,
		},
	}, nil
}

func (p *FilteredPaginationSubProtocol) HandleOrders(ctx context.Context, res *ordersync.Response) (*ordersync.Request, error) {
	if res.Metadata == nil {
		return nil, errors.New("FilteredPaginationSubProtocol received response with nil metadata")
	}
	metadata, ok := res.Metadata.(*FilteredPaginationResponseMetadata)
	if !ok {
		return nil, fmt.Errorf("FilteredPaginationSubProtocol received response with wrong metadata type (got %T)", res.Metadata)
	}
	// TODO(albrow): Check that this order matches our current filter/topic
	_, err := p.app.orderWatcher.ValidateAndStoreValidOrders(ctx, res.Orders, false, p.app.chainID)
	if err != nil {
		return nil, err
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
