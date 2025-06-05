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
	return &postService{
		storage: storage,
		log:     log,
	}
}

func (ps *postService) CreatePost(ctx context.Context, input models.CreatePostInput) (models.Post, error) {
	const op = "service.postService.CreatePost"
	log := ps.log.With(slog.String("op", op))

	post := models.Post{
		Title:           input.Title,
		Content:         input.Content,
		Author:          input.Author,
		CommentsEnabled: input.CommentsEnabled,
	}

	createdPost, err := ps.storage.CreatePost(ctx, post)
	if err != nil {
		log.Error("Failed to create post", sl.Err(err), "input", input)
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Post created", "id", createdPost.ID)
	return createdPost, nil
}

func (ps *postService) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {
	const op = "service.postService.GetPosts"
	log := ps.log.With(slog.String("op", op))

	posts, err := ps.storage.GetPosts(ctx, limit, offset)
	if err != nil {
		log.Error("Failed to get posts", sl.Err(err), "limit", limit, "offset", offset)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Posts retrieved", "count", len(posts))
	return posts, nil
}

func (ps *postService) GetPost(ctx context.Context, id string) (models.Post, error) {
	const op = "service.postService.GetPost"
	log := ps.log.With(slog.String("op", op))

	post, err := ps.storage.GetPost(ctx, id)
	if err != nil {
		log.Error("Failed to get post", sl.Err(err), "id", id)
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Post retrieved", "id", id)
	return post, nil
}

func (ps *postService) ToggleComments(ctx context.Context, postID string, enabled bool) (models.Post, error) {
	const op = "service.postService.ToggleComments"
	log := ps.log.With(slog.String("op", op))

	post, err := ps.storage.GetPost(ctx, postID)
	if err != nil {
		log.Error("Failed to get post", sl.Err(err), "id", postID)
		return models.Post{}, fmt.Errorf("%s: failed to get post: %w", op, err)
	}

	if post.CommentsEnabled == enabled {
		log.Info("Comments already in requested state", "enabled", enabled)
		return post, nil
	}

	post.CommentsEnabled = enabled
	if err := ps.storage.UpdatePost(ctx, post); err != nil {
		log.Error("Failed to update post", sl.Err(err), "id", postID)
		return models.Post{}, fmt.Errorf("%s: failed to update post: %w", op, err)
	}

	log.Info("Comments toggled", "id", postID, "enabled", enabled)
	return post, nil
}
