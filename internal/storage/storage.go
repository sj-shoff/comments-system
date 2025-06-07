package storage

import (
	"comments-system/internal/models"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=PostStorage --output=./mocks --case=underscore
type PostStorage interface {
	CreatePost(ctx context.Context, post models.Post) (models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error)
	GetPost(ctx context.Context, id string) (models.Post, error)
	UpdatePost(ctx context.Context, post models.Post) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=CommentStorage --output=./mocks --case=underscore
type CommentStorage interface {
	CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error)
	GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]models.Comment, error)
	GetComment(ctx context.Context, id string) (models.Comment, error)
	CountCommentsByPost(ctx context.Context, postID string) (int, error)
	GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=Storage --output=./mocks --case=underscore
type Storage interface {
	PostStorage
	CommentStorage
	Close() error
}
