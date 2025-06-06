package graph

import (
	"comments-system/internal/graph/generated"
	"comments-system/internal/models"
	"comments-system/internal/pubsub"
	"comments-system/internal/service"
	"comments-system/pkg/errors"
	"context"
	"fmt"
	"log/slog"
)

var _ generated.ResolverRoot = (*Resolver)(nil)

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

func (r *mutationResolver) CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	const op = "resolver.mutationResolver.CreatePost"
	log := r.log.With(slog.String("op", op))

	log.Debug("Creating post requested", "input", input)

	post, err := r.services.PostService.CreatePost(ctx, input)
	if err != nil {
		log.Error("Create post failed", "error", err, "input", input)
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	log.Info("Create post completed", "id", post.ID)
	return &post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (*models.Comment, error) {
	const op = "resolver.mutationResolver.CreateComment"
	log := r.log.With(slog.String("op", op))

	log.Debug("Creating comment requested", "input", input)

	post, err := r.services.PostService.GetPost(ctx, input.PostID)
	if err != nil {
		log.Error("Failed to get post before creating comment", "error", err, "postID", input.PostID)
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if !post.CommentsEnabled {
		log.Warn("Comments disabled for post", "postID", input.PostID)
		return nil, errors.ErrCommentsDisabled
	}

	comment, err := r.services.CommentService.CreateComment(ctx, input)
	if err != nil {
		log.Error("Failed to create comment", "error", err, "input", input)
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	r.ps.Publish(input.PostID, &comment)
	log.Info("Comment created completed", "id", comment.ID, "postID", input.PostID)
	return &comment, nil
}

func (r *mutationResolver) ToggleComments(ctx context.Context, postID string, enabled bool) (*models.Post, error) {
	const op = "resolver.mutationResolver.ToggleComments"
	log := r.log.With(slog.String("op", op))

	log.Debug("Toggling comments requested", "postID", postID, "enabled", enabled)
	post, err := r.services.PostService.ToggleComments(ctx, postID, enabled)
	if err != nil {
		log.Error("Toggle comments failed", "error", err, "postID", postID, "enabled", enabled)
		return nil, fmt.Errorf("failed to toggle comments: %w", err)
	}

	log.Info("Comments toggled completed", "postID", postID, "enabled", enabled)
	return &post, nil
}

func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	const op = "resolver.queryResolver.Posts"
	log := r.log.With(slog.String("op", op))

	l := 10
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}

	log.Debug("Getting posts requested", "limit", l, "offset", o)
	posts, err := r.services.PostService.GetPosts(ctx, l, o)
	if err != nil {
		log.Error("Failed to get posts", "error", err, "limit", l, "offset", o)
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	result := make([]*models.Post, len(posts))
	for i := range posts {
		result[i] = &posts[i]
	}

	log.Info("Posts retrieved completed", "count", len(posts))
	return result, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	const op = "resolver.queryResolver.Post"
	log := r.log.With(slog.String("op", op))

	log.Debug("Getting post requested", "id", id)

	post, err := r.services.PostService.GetPost(ctx, id)
	if err != nil {
		log.Error("Failed to get post", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	log.Info("Post retrieved completed", "id", id)
	return &post, nil
}

func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) (*models.CommentsPage, error) {
	const op = "resolver.queryResolver.Comments"
	log := r.log.With(slog.String("op", op))

	l := 10
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}

	log.Debug("Getting comments requested", "postID", postID, "limit", l, "offset", o)

	comments, total, err := r.services.CommentService.GetComments(ctx, postID, l, o)
	if err != nil {
		log.Error("Failed to get comments", "error", err, "postID", postID, "limit", l, "offset", o)
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	log.Info("Comments retrieved completed", "postID", postID, "count", len(comments), "total", total)
	return &models.CommentsPage{
		Total:    total,
		Comments: comments,
	}, nil
}

func (r *queryResolver) CommentReplies(ctx context.Context, parentID string) ([]*models.Comment, error) {
	const op = "resolver.queryResolver.CommentReplies"
	log := r.log.With(slog.String("op", op))

	log.Debug("Getting comment replies requested", "parentID", parentID)

	replies, err := r.services.CommentService.GetCommentReplies(ctx, parentID)
	if err != nil {
		log.Error("Failed to get comment replies", "error", err, "parentID", parentID)
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}

	result := make([]*models.Comment, len(replies))
	for i := range replies {
		result[i] = &replies[i]
	}

	log.Info("Comment replies retrieved completed", "parentID", parentID, "count", len(replies))
	return result, nil
}

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	const op = "resolver.subscriptionResolver.CommentAdded"
	log := r.log.With(slog.String("op", op))

	log.Debug("Subscribing to comments requested", "postID", postID)

	ch, err := r.ps.Subscribe(ctx, postID)
	if err != nil {
		log.Error("Failed to subscribe to comments", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	log.Info("Subscribed to comments completed", "postID", postID)
	return ch, nil
}
