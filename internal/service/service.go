package service

import (
	"comments-system/internal/models"
	"context"
)

type PostService interface {
	CreatePost(ctx context.Context, input models.CreatePostInput) (models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error)
	GetPost(ctx context.Context, id string) (models.Post, error)
	ToggleComments(ctx context.Context, postID string, enabled bool) (models.Post, error)
}

type CommentService interface {
	CreateComment(ctx context.Context, input models.CreateCommentInput) (models.Comment, error)
	GetComments(ctx context.Context, postID string, limit, offset int) ([]models.Comment, int, error)
	GetComment(ctx context.Context, id string) (models.Comment, error)
	GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error)
}

type Service struct {
	PostService
	CommentService
}
