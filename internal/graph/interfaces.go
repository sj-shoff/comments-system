package graph

import (
	"comments-system/internal/models"
	"context"
)

type MutationResolver interface {
	CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error)
	CreateComment(ctx context.Context, input models.CreateCommentInput) (*models.Comment, error)
	ToggleComments(ctx context.Context, postID string, enabled bool) (*models.Post, error)
}

type QueryResolver interface {
	Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error)
	Post(ctx context.Context, id string) (*models.Post, error)
	Comments(ctx context.Context, postID string, limit *int, offset *int) (*models.CommentsPage, error)
	CommentReplies(ctx context.Context, parentID string) ([]*models.Comment, error)
}

type SubscriptionResolver interface {
	CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error)
}
