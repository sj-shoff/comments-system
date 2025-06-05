package service

import (
	"comments-system/internal/models"
	"comments-system/internal/storage"
	"comments-system/pkg/logger/sl"
	"context"
	"fmt"
	"log/slog"
)

type postService struct {
	storage storage.Storage
	log     *slog.Logger
}

func NewPostService(storage storage.Storage, log *slog.Logger) PostService {
	const op = "service.post_service.NewPostService"
	log = log.With(slog.String("op", op))
	return &postService{storage: storage, log: log}
}

func (ps *postService) CreatePost(ctx context.Context, input models.CreatePostInput) (models.Post, error) {
	const op = "service.post_service.CreatePost"

	post := models.Post{
		Title:           input.Title,
		Content:         input.Content,
		Author:          input.Author,
		CommentsEnabled: input.CommentsEnabled,
	}

	createdPost, err := ps.storage.CreatePost(ctx, post)
	if err != nil {
		ps.log.Error("failed to create post", sl.Err(err))
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	ps.log.Info("Post created")
	return createdPost, nil
}

func (ps *postService) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {
	const op = "service.post_service.GetPosts"

	posts, err := ps.storage.GetPosts(ctx)
	if err != nil {
		ps.log.Error("failed to get posts", slog.String("op", op), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Implement pagination if limit and offset are used.  For now, return all posts.
	if limit > 0 && offset >= 0 {
		//TODO Implement pagination
		ps.log.Warn("pagination is not implemented, returning all posts", slog.String("op", op), slog.Int("limit", limit), slog.Int("offset", offset))
	}

	ps.log.Info("posts retrieved", slog.String("op", op), slog.Int("count", len(posts)))
	return posts, nil
}

func (ps *postService) GetPost(ctx context.Context, id string) (models.Post, error) {
	const op = "service.post_service.GetPost"

	post, err := ps.storage.GetPost(ctx, id)
	if err != nil {
		ps.log.Error("failed to get post", slog.String("op", op), slog.String("post_id", id), slog.Any("error", err))
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	ps.log.Info("post retrieved", slog.String("op", op), slog.String("post_id", id))
	return post, nil
}

func (ps *postService) ToggleComments(ctx context.Context, postID string, enabled bool) (models.Post, error) {
	const op = "service.post_service.ToggleComments"

	post, err := ps.storage.GetPost(ctx, postID)
	if err != nil {
		ps.log.Error("Failed to get post for toggling comments", sl.Err(err))
		return models.Post{}, fmt.Errorf("%s: failed to get post: %w", op, err)
	}

	post.CommentsEnabled = enabled
	if err := ps.storage.UpdatePost(ctx, post); err != nil {
		ps.log.Error("Failed to update post (toggle comments)", sl.Err(err))
		return models.Post{}, fmt.Errorf("%s: failed to update post: %w", op, err)
	}

	ps.log.Info("Comments toggled")
	return post, nil
}
