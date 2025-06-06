package graph

import (
	"comments-system/internal/graph/generated"
)

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() generated.SubscriptionResolver {
	return &subscriptionResolver{r}
}
