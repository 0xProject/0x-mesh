package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/0xProject/0x-mesh/graphql/generated"
	"github.com/0xProject/0x-mesh/graphql/types"
)

type Resolver struct{}

func (r *mutationResolver) AddOrders(ctx context.Context, orders []*types.NewOrder, pinned *bool) (*types.AddOrdersResults, error) {
	panic("not implemented")
}

func (r *queryResolver) Order(ctx context.Context, hash types.Hash) (*types.OrderWithMetadata, error) {
	panic("not implemented")
}

func (r *queryResolver) Orders(ctx context.Context, sort []*types.OrderSort, filters []*types.OrderFilter, limit *int) ([]*types.OrderWithMetadata, error) {
	panic("not implemented")
}

func (r *queryResolver) Stats(ctx context.Context) (*types.Stats, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) OrderEvents(ctx context.Context) (<-chan []*types.OrderEvent, error) {
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
