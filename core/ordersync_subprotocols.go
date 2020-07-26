package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/core/ordersync"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// ordersyncSubprotocolFactory is a function that can be used to create an ordersync.Subprotocol.
// Note(albrow): Using a factory here allows us to specify which subprotocols we want to use before
// the app is fully initialized. Factory functions won't actually be called until the app is done
// initializing.
type ordersyncSubprotocolFactory func(app *App, perPage int) ordersync.Subprotocol

// Ensure that FilteredPaginationSubProtocolV0 implements the Subprotocol interface.
var _ ordersync.Subprotocol = (*FilteredPaginationSubProtocolV0)(nil)

// FilteredPaginationSubProtocolV0 is an ordersync subprotocol which returns all orders by
// paginating through them. It involves sending multiple requests until pagination is
// finished and all orders have been returned. Version 0 of the subprotocol is deprecated
// but included for backwards-compatibility.
type FilteredPaginationSubProtocolV0 struct {
	app         *App
	orderFilter *orderfilter.Filter
	perPage     int
}

// NewFilteredPaginationSubprotocolV0 creates and returns a new FilteredPaginationSubprotocolV0
// which will respond with perPage orders for each individual request/response.
func NewFilteredPaginationSubprotocolV0(app *App, perPage int) ordersync.Subprotocol {
	return &FilteredPaginationSubProtocolV0{
		app:         app,
		orderFilter: app.orderFilter,
		perPage:     perPage,
	}
}

// FilteredPaginationRequestMetadataV0 is the request metadata for the
// FilteredPaginationSubProtocolV0. It keeps track of the current page and SnapshotID,
// which is expected to be an empty string on the first request.
type FilteredPaginationRequestMetadataV0 struct {
	Page       int    `json:"page"`
	SnapshotID string `json:"snapshotID"`
}

// FilteredPaginationResponseMetadataV0 is the response metadata for the
// FilteredPaginationSubProtocolV0. It keeps track of the current page and SnapshotID.
type FilteredPaginationResponseMetadataV0 struct {
	Page       int    `json:"page"`
	SnapshotID string `json:"snapshotID"`
}

// Name returns the name of the FilteredPaginationSubProtocolV0
func (p *FilteredPaginationSubProtocolV0) Name() string {
	return "/pagination-with-filter/version/0"
}

// HandleOrderSyncRequest returns the orders for one page, based on the page number
// and snapshotID corresponding to the given request. This is
// the implementation for the "provider" side of the subprotocol.
func (p *FilteredPaginationSubProtocolV0) HandleOrderSyncRequest(ctx context.Context, req *ordersync.Request) (*ordersync.Response, error) {
	var metadata *FilteredPaginationRequestMetadataV0
	if req.Metadata == nil {
		// Default metadata for the first request.
		metadata = &FilteredPaginationRequestMetadataV0{
			Page:       0,
			SnapshotID: "",
		}
	} else {
		var ok bool
		metadata, ok = req.Metadata.(*FilteredPaginationRequestMetadataV0)
		if !ok {
			return nil, fmt.Errorf("FilteredPaginationSubProtocolV0 received request with wrong metadata type (got %T)", req.Metadata)
		}
	}

	// Note(albrow): This version of Mesh does not support database snapshots. Instead, we use the SnapshotID
	// field as minOrderHash.
	var currentMinOrderHash common.Hash
	if metadata.SnapshotID != "" {
		if err := validateHexHash(metadata.SnapshotID); err != nil {
			return nil, fmt.Errorf("FilteredPaginationSubProtocolV0 could not decode snapshotID (%q) as hex: %s", metadata.SnapshotID, err.Error())
		}
		currentMinOrderHash = common.HexToHash(metadata.SnapshotID)
	}

	// It's possible that none of the orders in the current page match the filter.
	// We don't want to respond with zero orders, so keep iterating until we find
	// at least some orders that match the filter.
	filteredOrders := []*zeroex.SignedOrder{}
	var nextMinOrderHash common.Hash
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		// Get the orders for this page.
		ordersResp, err := p.app.GetOrders(p.perPage, currentMinOrderHash)
		if err != nil {
			return nil, err
		}
		if len(ordersResp.OrdersInfos) == 0 {
			// No more orders left.
			break
		}
		nextMinOrderHash = ordersResp.OrdersInfos[len(ordersResp.OrdersInfos)-1].OrderHash
		// Filter the orders for this page.
		for _, orderInfo := range ordersResp.OrdersInfos {
			if matches, err := p.orderFilter.MatchOrder(orderInfo.SignedOrder); err != nil {
				return nil, err
			} else if matches {
				filteredOrders = append(filteredOrders, orderInfo.SignedOrder)
			}
		}
		if len(filteredOrders) == 0 {
			// If none of the orders for this page match the filter, we continue
			// on to the next page.
			currentMinOrderHash = nextMinOrderHash
			continue
		} else {
			break
		}
	}

	return &ordersync.Response{
		Orders:   filteredOrders,
		Complete: len(filteredOrders) == 0,
		Metadata: &FilteredPaginationResponseMetadataV0{
			// Note(albrow): Page isn't actually used. Included for backwards compatibility only.
			Page:       metadata.Page + 1,
			SnapshotID: nextMinOrderHash.Hex(),
		},
	}, nil
}

// HandleOrderSyncResponse handles the orders for one page by validating them, storing them
// in the database, and firing the appropriate events. It also returns the next request to
// be sent. This is the implementation for the "requester" side of the subprotocol.
func (p *FilteredPaginationSubProtocolV0) HandleOrderSyncResponse(ctx context.Context, res *ordersync.Response) (*ordersync.Request, []*zeroex.SignedOrder, error) {
	if res.Metadata == nil {
		return nil, nil, errors.New("FilteredPaginationSubProtocolV0 received response with nil metadata")
	}
	metadata, ok := res.Metadata.(*FilteredPaginationResponseMetadataV0)
	if !ok {
		return nil, nil, fmt.Errorf("FilteredPaginationSubProtocolV0 received response with wrong metadata type (got %T)", res.Metadata)
	}
	filteredOrders := []*zeroex.SignedOrder{}
	for _, order := range res.Orders {
		if matches, err := p.orderFilter.MatchOrder(order); err != nil {
			return nil, nil, err
		} else if matches {
			filteredOrders = append(filteredOrders, order)
		} else if !matches {
			p.app.handlePeerScoreEvent(res.ProviderID, psReceivedOrderDoesNotMatchFilter)
		}
	}
	validationResults, err := p.app.orderWatcher.ValidateAndStoreValidOrders(ctx, filteredOrders, false, p.app.chainID)
	if err != nil {
		return nil, filteredOrders, err
	}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		if acceptedOrderInfo.IsNew {
			log.WithFields(map[string]interface{}{
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
				"protocol":  "ordersync",
			}).Info("received new valid order from peer")
			log.WithFields(map[string]interface{}{
				"order":     acceptedOrderInfo.SignedOrder,
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
				"protocol":  "ordersync",
			}).Trace("all fields for new valid order received from peer")
		}
	}

	return &ordersync.Request{
		Metadata: &FilteredPaginationRequestMetadataV0{
			Page:       metadata.Page + 1,
			SnapshotID: metadata.SnapshotID,
		},
	}, filteredOrders, nil
}

func (p *FilteredPaginationSubProtocolV0) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationRequestMetadataV0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (p *FilteredPaginationSubProtocolV0) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationResponseMetadataV0
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Ensure that FilteredPaginationSubProtocolV1 implements the Subprotocol interface.
var _ ordersync.Subprotocol = (*FilteredPaginationSubProtocolV1)(nil)

// FilteredPaginationSubProtocolV1 is an ordersync subprotocol which returns all orders by
// paginating through them. It involves sending multiple requests until pagination is
// finished and all orders have been returned. Version 1 was implemented in
// https://github.com/0xProject/0x-mesh/pull/793 after changing the database implementation
// from LevelDB to SQL and Dexie.js/IndexedDB.
type FilteredPaginationSubProtocolV1 struct {
	app         *App
	orderFilter *orderfilter.Filter
	perPage     int
}

// NewFilteredPaginationSubprotocolV1 creates and returns a new FilteredPaginationSubprotocolV1
// which will respond with perPage orders for each individual request/response.
func NewFilteredPaginationSubprotocolV1(app *App, perPage int) ordersync.Subprotocol {
	return &FilteredPaginationSubProtocolV1{
		app:         app,
		orderFilter: app.orderFilter,
		perPage:     perPage,
	}
}

// FilteredPaginationRequestMetadataV1 is the request metadata for the
// FilteredPaginationSubProtocolV1. It keeps track of the current
// minOrderHash, which is expected to be an empty string on the first request.
type FilteredPaginationRequestMetadataV1 struct {
	MinOrderHash common.Hash `json:"minOrderHash"`
}

// FilteredPaginationResponseMetadataV1 is the response metadata for the
// FilteredPaginationSubProtocolV1. It contains the minOrderHash to use for
// the next request.
type FilteredPaginationResponseMetadataV1 struct {
	NextMinOrderHash common.Hash `json:"nextMinOrderHash"`
}

// Name returns the name of the FilteredPaginationSubProtocolV1
func (p *FilteredPaginationSubProtocolV1) Name() string {
	return "/pagination-with-filter/version/1"
}

// HandleOrderSyncRequest returns the orders for one page, based on the page number
// and snapshotID corresponding to the given request. This is
// the implementation for the "provider" side of the subprotocol.
func (p *FilteredPaginationSubProtocolV1) HandleOrderSyncRequest(ctx context.Context, req *ordersync.Request) (*ordersync.Response, error) {
	var metadata *FilteredPaginationRequestMetadataV1
	if req.Metadata == nil {
		// Default metadata for the first request.
		metadata = &FilteredPaginationRequestMetadataV1{
			MinOrderHash: common.Hash{},
		}
	} else {
		var ok bool
		metadata, ok = req.Metadata.(*FilteredPaginationRequestMetadataV1)
		if !ok {
			return nil, fmt.Errorf("FilteredPaginationSubProtocolV1 received request with wrong metadata type (got %T)", req.Metadata)
		}
	}

	// It's possible that none of the orders in the current page match the filter.
	// We don't want to respond with zero orders, so keep iterating until we find
	// at least some orders that match the filter.
	filteredOrders := []*zeroex.SignedOrder{}
	currentMinOrderHash := metadata.MinOrderHash
	nextMinOrderHash := common.Hash{}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		// Get the orders for this page.
		ordersResp, err := p.app.GetOrders(p.perPage, currentMinOrderHash)
		if err != nil {
			return nil, err
		}
		if len(ordersResp.OrdersInfos) == 0 {
			// No more orders left.
			break
		}
		nextMinOrderHash = ordersResp.OrdersInfos[len(ordersResp.OrdersInfos)-1].OrderHash
		// Filter the orders for this page.
		for _, orderInfo := range ordersResp.OrdersInfos {
			if matches, err := p.orderFilter.MatchOrder(orderInfo.SignedOrder); err != nil {
				return nil, err
			} else if matches {
				filteredOrders = append(filteredOrders, orderInfo.SignedOrder)
			}
		}
		if len(filteredOrders) == 0 {
			// If none of the orders for this page match the filter, we continue
			// on to the next page.
			currentMinOrderHash = nextMinOrderHash
			continue
		} else {
			break
		}
	}

	return &ordersync.Response{
		Orders:   filteredOrders,
		Complete: len(filteredOrders) == 0,
		Metadata: &FilteredPaginationResponseMetadataV1{
			NextMinOrderHash: nextMinOrderHash,
		},
	}, nil
}

// HandleOrderSyncResponse handles the orders for one page by validating them, storing them
// in the database, and firing the appropriate events. It also returns the next request to
// be sent. This is the implementation for the "requester" side of the subprotocol.
func (p *FilteredPaginationSubProtocolV1) HandleOrderSyncResponse(ctx context.Context, res *ordersync.Response) (*ordersync.Request, []*zeroex.SignedOrder, error) {
	if res.Metadata == nil {
		return nil, nil, errors.New("FilteredPaginationSubProtocolV1 received response with nil metadata")
	}
	_, ok := res.Metadata.(*FilteredPaginationResponseMetadataV1)
	if !ok {
		return nil, nil, fmt.Errorf("FilteredPaginationSubProtocolV1 received response with wrong metadata type (got %T)", res.Metadata)
	}
	filteredOrders := []*zeroex.SignedOrder{}
	for _, order := range res.Orders {
		if matches, err := p.orderFilter.MatchOrder(order); err != nil {
			return nil, nil, err
		} else if matches {
			filteredOrders = append(filteredOrders, order)
		} else if !matches {
			p.app.handlePeerScoreEvent(res.ProviderID, psReceivedOrderDoesNotMatchFilter)
		}
	}
	validationResults, err := p.app.orderWatcher.ValidateAndStoreValidOrders(ctx, filteredOrders, false, p.app.chainID)
	if err != nil {
		return nil, filteredOrders, err
	}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		if acceptedOrderInfo.IsNew {
			log.WithFields(map[string]interface{}{
				"orderHash": acceptedOrderInfo.OrderHash.Hex(),
				"from":      res.ProviderID.Pretty(),
				"protocol":  "ordersync",
			}).Info("received new valid order from peer")
			log.WithFields(map[string]interface{}{
				"order":     acceptedOrderInfo.SignedOrder,
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
			return nil, filteredOrders, err
		}
		nextMinOrderHash = hash
	}
	return &ordersync.Request{
		Metadata: &FilteredPaginationRequestMetadataV1{
			MinOrderHash: nextMinOrderHash,
		},
	}, filteredOrders, nil
}

func (p *FilteredPaginationSubProtocolV1) ParseRequestMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationRequestMetadataV1
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (p *FilteredPaginationSubProtocolV1) ParseResponseMetadata(metadata json.RawMessage) (interface{}, error) {
	var parsed FilteredPaginationResponseMetadataV1
	if err := json.Unmarshal(metadata, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

// validateHexHash returns an error if s is not a valid hex hash. It supports
// encodings with or without the "0x" prefix.
// Note(albrow) This is based on unexported code in go-ethereum.
func validateHexHash(s string) error {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	_, err := hex.DecodeString(s)
	return err
}

// has0xPrefix returns true if the given hex string starts with "0x"
// Note(albrow) This is copied from go-ethereum, where it is unexported.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
