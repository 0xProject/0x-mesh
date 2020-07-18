package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/graphql/generated"
	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/ethereum/go-ethereum/common"
)

func (r *mutationResolver) AddOrders(ctx context.Context, orders []*gqltypes.NewOrder, pinned *bool) (*gqltypes.AddOrdersResults, error) {
	isPinned := false
	if pinned != nil {
		isPinned = (*pinned)
	}
	signedOrders := gqltypes.NewOrdersToSignedOrders(orders)
	results, err := r.app.AddOrders(ctx, signedOrders, isPinned)
	if err != nil {
		return nil, err
	}
	return gqltypes.AddOrdersResultsFromValidationResults(results)
}

func (r *queryResolver) Order(ctx context.Context, hash gqltypes.Hash) (*gqltypes.OrderWithMetadata, error) {
	order, err := r.app.GetOrder(common.Hash(hash))
	if err != nil {
		return nil, err
	}
	return gqltypes.OrderWithMetadataFromCommonType(order), nil
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
		filterValue, err := gqltypes.ConvertFilterValue(filter)
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

func (r *queryResolver) Stats(ctx context.Context) (*gqltypes.Stats, error) {
	stats, err := r.app.GetStats()
	if err != nil {
		return nil, err
	}
	return gqltypes.StatsFromCommonType(stats), nil
}

func (r *subscriptionResolver) OrderEvents(ctx context.Context) (<-chan []*gqltypes.OrderEvent, error) {
	panic("not implemented")
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
