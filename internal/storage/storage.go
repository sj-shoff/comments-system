package storage

import (
	"comments-system/internal/models"
	"context"
)

type PostStorage interface {
	CreatePost(ctx context.Context, post models.Post) (models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error)
	GetPost(ctx context.Context, id string) (models.Post, error)
	UpdatePost(ctx context.Context, post models.Post) error
}

type CommentStorage interface {
	CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error)
	GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]models.Comment, error)
	GetComment(ctx context.Context, id string) (models.Comment, error)
	CountCommentsByPost(ctx context.Context, postID string) (int, error)
	GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error)
}

type Storage interface {
	PostStorage
	CommentStorage
	Close() error
}
