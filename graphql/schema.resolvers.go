package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/graphql/generated"
	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func (r *mutationResolver) AddOrders(ctx context.Context, orders []*gqltypes.NewOrder, pinned *bool, opts *gqltypes.AddOrdersOpts) (*gqltypes.AddOrdersResults, error) {
	isPinned := false
	if pinned != nil {
		isPinned = (*pinned)
	}
	signedOrders := gqltypes.NewOrdersToSignedOrders(orders)
	commonTypeOpts := gqltypes.AddOrderOptsToCommonType(opts)
	results, err := r.app.AddOrders(ctx, signedOrders, isPinned, commonTypeOpts)
	if err != nil {
		return nil, err
	}
	return gqltypes.AddOrdersResultsFromValidationResults(results)
}

func (r *mutationResolver) AddOrdersV4(ctx context.Context, orders []*gqltypes.NewOrderV4, pinned *bool, opts *gqltypes.AddOrdersOpts) (*gqltypes.AddOrdersResults, error) {
	isPinned := false
	if pinned != nil {
		isPinned = (*pinned)
	}
	signedOrders := gqltypes.NewOrdersToSignedOrdersV4(orders)
	commonTypeOpts := gqltypes.AddOrderOptsToCommonType(opts)
	results, err := r.app.AddOrdersV4(ctx, signedOrders, isPinned, commonTypeOpts)
	if err != nil {
		return nil, err
	}
	return gqltypes.AddOrdersResultsFromValidationResults(results)
}

func (r *queryResolver) Order(ctx context.Context, hash string) (*gqltypes.OrderWithMetadata, error) {
	order, err := r.app.GetOrder(common.HexToHash(hash))
	if err != nil {
		if err == db.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return gqltypes.OrderWithMetadataFromCommonType(order), nil
}

func (r *queryResolver) Orderv4(ctx context.Context, hash string) (*gqltypes.OrderV4WithMetadata, error) {
	order, err := r.app.GetOrderV4(common.HexToHash(hash))
	if err != nil {
		if err == db.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return gqltypes.OrderWithMetadataFromCommonTypeV4(order), nil
}

func (r *queryResolver) Orders(ctx context.Context, sort []*gqltypes.OrderSort, filters []*gqltypes.OrderFilter, limit *int) ([]*gqltypes.OrderWithMetadata, error) {
	// TODO(albrow): More validation of query args. We can assume
	//               basic structure is correct but may need to validate
	//               some of the semantics.
	query := &db.OrderQuery{
		// We never include orders that are marked as removed.
		Filters: []db.OrderFilter{
			{
				Field: db.OFIsRemoved,
				Kind:  db.Equal,
				Value: false,
			},
		},
	}
	if limit != nil {
		query.Limit = uint(*limit)
	}
	for _, filter := range filters {
		kind, err := gqltypes.FilterKindToDBType(filter.Kind)
		if err != nil {
			return nil, err
		}
		filterValue, err := gqltypes.FilterValueFromJSON(*filter)
		if err != nil {
			return nil, err
		}
		query.Filters = append(query.Filters, db.OrderFilter{
			Field: db.OrderField(filter.Field),
			Kind:  kind,
			Value: filterValue,
		})
	}
	for _, sort := range sort {
		direction, err := gqltypes.SortDirectionToDBType(sort.Direction)
		if err != nil {
			return nil, err
		}
		query.Sort = append(query.Sort, db.OrderSort{
			Field:     db.OrderField(sort.Field),
			Direction: direction,
		})
	}

	orders, err := r.app.FindOrders(query)
	if err != nil {
		return nil, err
	}
	return gqltypes.OrdersWithMetadataFromCommonType(orders), nil
}

func (r *queryResolver) Ordersv4(ctx context.Context, sort []*gqltypes.OrderSortV4, filters []*gqltypes.OrderFilterV4, limit *int) ([]*gqltypes.OrderV4WithMetadata, error) {
	query := &db.OrderQueryV4{
		// We never include orders that are marked as removed.
		Filters: []db.OrderFilterV4{
			{
				Field: db.OV4FIsRemoved,
				Kind:  db.Equal,
				Value: false,
			},
		},
	}
	if limit != nil {
		query.Limit = uint(*limit)
	}
	for _, filter := range filters {
		kind, err := gqltypes.FilterKindToDBType(filter.Kind)
		if err != nil {
			return nil, err
		}
		filterValue, err := gqltypes.FilterValueFromJSONV4(*filter)
		if err != nil {
			return nil, err
		}
		query.Filters = append(query.Filters, db.OrderFilterV4{
			Field: db.OrderFieldV4(filter.Field),
			Kind:  kind,
			Value: filterValue,
		})
	}
	for _, sort := range sort {
		direction, err := gqltypes.SortDirectionToDBType(sort.Direction)
		if err != nil {
			return nil, err
		}
		query.Sort = append(query.Sort, db.OrderSortV4{
			Field:     db.OrderFieldV4(sort.Field),
			Direction: direction,
		})
	}

	orders, err := r.app.FindOrdersV4(query)
	if err != nil {
		return nil, err
	}
	return gqltypes.OrdersWithMetadataFromCommonTypeV4(orders), nil
}

func (r *queryResolver) Stats(ctx context.Context) (*gqltypes.Stats, error) {
	stats, err := r.app.GetStats()
	if err != nil {
		return nil, err
	}
	return gqltypes.StatsFromCommonType(stats), nil
}

func (r *subscriptionResolver) OrderEvents(ctx context.Context) (<-chan []*gqltypes.OrderEvent, error) {
	zeroExChan := make(chan []*zeroex.OrderEvent, orderEventBufferSize)
	gqlChan := make(chan []*gqltypes.OrderEvent, orderEventBufferSize)
	subscription := r.app.SubscribeToOrderEvents(zeroExChan)
	// TODO(albrow): Call subscription.Unsubscribe for slow or disconnected clients.
	go func() {
		for {
			select {
			case <-ctx.Done():
				subscription.Unsubscribe()
				return
			case err := <-subscription.Err():
				// TODO(albrow): Can we handle this better?
				if err != nil {
					subscription.Unsubscribe()
					panic(err)
				}
			case orderEvents := <-zeroExChan:
				select {
				case gqlChan <- gqltypes.OrderEventsFromZeroExType(orderEvents):
					log.Debugf("sent %d orders to subscriber", len(orderEvents))
				case <-time.After(r.config.SlowSubscriberTimeout):
					log.Debug("subscriber is slow or disconnected, unsubscribing")
					subscription.Unsubscribe()
					close(gqlChan)
					return
				}

			}
		}
	}()
	return gqlChan, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
