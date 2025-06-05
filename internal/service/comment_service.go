package service

import (
	"comments-system/internal/models"
	"comments-system/internal/storage"
	"comments-system/pkg/logger/sl"
	"comments-system/pkg/utils"
	"context"
	"fmt"
	"log/slog"
)

type commentService struct {
	storage storage.Storage
	log     *slog.Logger
}

func NewCommentService(storage storage.Storage, log *slog.Logger) CommentService {
	const op = "service.comment_service.NewCommentService"
	log = log.With(slog.String("op", op))
	return &commentService{storage: storage, log: log}
}

func (c *commentService) CreateComment(ctx context.Context, input models.CreateCommentInput) (models.Comment, error) {
	const op = "service.comment_service.CreateComment"

	if err := utils.ValidateComment(input.Content); err != nil {
		c.log.Error("invalid comment content", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	comment := models.Comment{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		Author:   input.Author,
		Content:  input.Content,
	}

	createdComment, err := c.storage.CreateComment(ctx, comment)
	if err != nil {
		c.log.Error("failed to create comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: failed to create comment: %w", op, err)
	}

	c.log.Info("comment created")
	return createdComment, nil
}

// GetComments получает комментарии для поста с учетом лимита и смещения.
func (c *commentService) GetComments(ctx context.Context, postID string, limit, offset int) ([]models.Comment, int, error) {
	const op = "service.comment_service.GetComments"

	comments, err := c.storage.GetCommentsByPost(ctx, postID, limit, offset) // Передаем postID
	if err != nil {
		c.log.Error("failed to get comments", slog.String("op", op), slog.String("post_id", postID), slog.Any("error", err)) // Логирование ошибки получения
		return nil, 0, fmt.Errorf("%s: failed to get comments: %w", op, err)
	}

	total, err := c.storage.CountCommentsByPost(ctx, postID)
	if err != nil {
		c.log.Error("failed to count comments", sl.Err(err))
		return nil, 0, fmt.Errorf("%s: failed to count comments: %w", op, err)
	}

	c.log.Info("comments retrieved")
	return comments, total, nil
}

// GetCommentReplies получает ответы на комментарий.
func (c *commentService) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	const op = "service.comment_service.GetCommentReplies"

	replies, err := c.storage.GetCommentReplies(ctx, parentID)
	if err != nil {
		c.log.Error("failed to get comment replies", slog.String("op", op), slog.String("parent_id", parentID), slog.Any("error", err)) // Логирование ошибки получения
		return nil, fmt.Errorf("%s: failed to get comment replies: %w", op, err)
	}

	c.log.Info("comment replies retrieved", slog.String("op", op), slog.String("parent_id", parentID), slog.Int("count", len(replies))) // Логирование успешного получения
	return replies, nil
}

// GetComment получает комментарий по ID.
func (c *commentService) GetComment(ctx context.Context, id string) (models.Comment, error) {
	const op = "service.comment_service.GetComment"

	comment, err := c.storage.GetComment(ctx, id)
	if err != nil {
		c.log.Error("failed to get comment", slog.String("op", op), slog.String("comment_id", id), slog.Any("error", err)) // Логирование ошибки получения
		return models.Comment{}, fmt.Errorf("%s: failed to get comment: %w", op, err)
	}

	c.log.Info("comment retrieved", slog.String("op", op), slog.String("comment_id", id)) // Логирование успешного получения
	return comment, nil
}
