package graph

import (
	"comments-system/internal/models"
	"comments-system/internal/pubsub"
	"comments-system/internal/service"
	"comments-system/pkg/errors"
	"context"
	"fmt"
	"log/slog"
)

type Resolver struct {
	services *service.Service
	ps       *pubsub.PubSub
	log      *slog.Logger
}

func NewResolver(services *service.Service, ps *pubsub.PubSub, log *slog.Logger) *Resolver {
	return &Resolver{
		services: services,
		ps:       ps,
		log:      log,
	}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	r.log.Debug("Creating post", "input", input)

	post, err := r.services.PostService.CreatePost(ctx, input)
	if err != nil {
		r.log.Error("Failed to create post", "error", err, "input", input)
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	r.log.Info("Post created", "id", post.ID)
	return &post, nil
}
func (r *mutationResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (*models.Comment, error) {
	r.log.Debug("Creating comment", "input", input)

	// Проверяем разрешены ли комментарии для поста
	post, err := r.services.PostService.GetPost(ctx, input.PostID)
	if err != nil {
		r.log.Error("Failed to get post for comment", "postID", input.PostID, "error", err)
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if !post.CommentsEnabled {
		r.log.Warn("Comments disabled for post", "postID", input.PostID)
		return nil, errors.ErrCommentsDisabled
	}

	comment, err := r.services.CommentService.CreateComment(ctx, input)
	if err != nil {
		r.log.Error("Failed to create comment", "error", err, "input", input)
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Публикуем событие о новом комментарии
	r.ps.Publish(input.PostID, &comment)
	r.log.Info("Comment created", "id", comment.ID, "postID", input.PostID)

	return &comment, nil
}
func (r *mutationResolver) ToggleComments(ctx context.Context, postID string, enabled bool) (*models.Post, error) {
	r.log.Debug("Toggling comments", "postID", postID, "enabled", enabled)

	post, err := r.services.PostService.ToggleComments(ctx, postID, enabled)
	if err != nil {
		r.log.Error("Failed to toggle comments", "postID", postID, "error", err)
		return nil, fmt.Errorf("failed to toggle comments: %w", err)
	}

	r.log.Info("Comments toggled", "postID", postID, "enabled", enabled)
	return &post, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	l := 10
	if limit != nil {
		l = *limit
	}

	o := 0
	if offset != nil {
		o = *offset
	}

	r.log.Debug("Getting posts", "limit", l, "offset", o)

	posts, err := r.services.PostService.GetPosts(ctx, l, o)
	if err != nil {
		r.log.Error("Failed to get posts", "error", err)
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	// Конвертируем в указатели для GraphQL
	result := make([]*models.Post, len(posts))
	for i := range posts {
		result[i] = &posts[i]
	}

	r.log.Info("Posts retrieved", "count", len(posts))
	return result, nil
}
func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	r.log.Debug("Getting post", "id", id)

	post, err := r.services.PostService.GetPost(ctx, id)
	if err != nil {
		r.log.Error("Failed to get post", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	r.log.Info("Post retrieved", "id", id)
	return &post, nil
}
func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) (*models.CommentsPage, error) {
	l := 10
	if limit != nil {
		l = *limit
	}

	o := 0
	if offset != nil {
		o = *offset
	}

	r.log.Debug("Getting comments", "postID", postID, "limit", l, "offset", o)

	comments, total, err := r.services.CommentService.GetComments(ctx, postID, l, o)
	if err != nil {
		r.log.Error("Failed to get comments", "postID", postID, "error", err)
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	r.log.Info("Comments retrieved", "postID", postID, "count", len(comments), "total", total)
	return &models.CommentsPage{
		Total:    total,
		Comments: comments,
	}, nil
}
func (r *queryResolver) CommentReplies(ctx context.Context, parentID string) ([]*models.Comment, error) {
	r.log.Debug("Getting comment replies", "parentID", parentID)

	replies, err := r.services.CommentService.GetCommentReplies(ctx, parentID)
	if err != nil {
		r.log.Error("Failed to get comment replies", "parentID", parentID, "error", err)
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	// Конвертируем в указатели для GraphQL
	result := make([]*models.Comment, len(replies))
	for i := range replies {
		result[i] = &replies[i]
	}

	r.log.Info("Comment replies retrieved", "parentID", parentID, "count", len(replies))
	return result, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	r.log.Debug("Subscribing to comments", "postID", postID)

	ch, err := r.ps.Subscribe(ctx, postID)
	if err != nil {
		r.log.Error("Failed to subscribe to comments", "postID", postID, "error", err)
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	r.log.Info("Subscribed to comments", "postID", postID)
	return ch, nil
}
