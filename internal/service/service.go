package service

import (
	"comments-system/internal/models"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=PostService --output=./mocks --case=underscore
type PostService interface {
	CreatePost(ctx context.Context, input models.CreatePostInput) (models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error)
	GetPost(ctx context.Context, id string) (models.Post, error)
	ToggleComments(ctx context.Context, postID string, enabled bool) (models.Post, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=CommentService --output=./mocks --case=underscore
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
