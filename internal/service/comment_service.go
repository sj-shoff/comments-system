package service

import (
	"comments-system/internal/models"
	"comments-system/internal/storage"
	"comments-system/pkg/errors"
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
	return &commentService{
		storage: storage,
		log:     log,
	}
}

func (cs *commentService) CreateComment(ctx context.Context, input models.CreateCommentInput) (models.Comment, error) {
	const op = "service.commentService.CreateComment"
	log := cs.log.With(slog.String("op", op))

	if err := utils.ValidateComment(input.Content); err != nil {
		log.Error("Invalid comment content", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	post, err := cs.storage.GetPost(ctx, input.PostID)
	if err != nil {
		log.Error("Failed to get post", sl.Err(err), "postID", input.PostID)
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	if !post.CommentsEnabled {
		log.Warn("Comments disabled for post", "postID", input.PostID)
		return models.Comment{}, errors.ErrCommentsDisabled
	}

	if input.ParentID != nil {
		_, err := cs.storage.GetComment(ctx, *input.ParentID)
		if err != nil {
			log.Error("Parent comment not found", sl.Err(err), "parentID", *input.ParentID)
			return models.Comment{}, fmt.Errorf("%s: %w", op, errors.ErrParentNotFound)
		}
	}

	comment := models.Comment{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		Author:   input.Author,
		Content:  input.Content,
	}

	createdComment, err := cs.storage.CreateComment(ctx, comment)
	if err != nil {
		log.Error("Failed to create comment", sl.Err(err), "input", input)
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Comment created", "id", createdComment.ID, "postID", input.PostID)
	return createdComment, nil
}

func (cs *commentService) GetComments(ctx context.Context, postID string, limit, offset int) ([]models.Comment, int, error) {
	const op = "service.commentService.GetComments"
	log := cs.log.With(slog.String("op", op))

	comments, err := cs.storage.GetCommentsByPost(ctx, postID, limit, offset)
	if err != nil {
		log.Error("Failed to get comments", sl.Err(err), "postID", postID)
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	total, err := cs.storage.CountCommentsByPost(ctx, postID)
	if err != nil {
		log.Error("Failed to count comments", sl.Err(err), "postID", postID)
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Comments retrieved", "postID", postID, "count", len(comments), "total", total)
	return comments, total, nil
}

func (cs *commentService) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	const op = "service.commentService.GetCommentReplies"
	log := cs.log.With(slog.String("op", op))

	replies, err := cs.storage.GetCommentReplies(ctx, parentID)
	if err != nil {
		log.Error("Failed to get comment replies", sl.Err(err), "parentID", parentID)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Comment replies retrieved", "parentID", parentID, "count", len(replies))
	return replies, nil
}

func (cs *commentService) GetComment(ctx context.Context, id string) (models.Comment, error) {
	const op = "service.commentService.GetComment"
	log := cs.log.With(slog.String("op", op))

	comment, err := cs.storage.GetComment(ctx, id)
	if err != nil {
		log.Error("Failed to get comment", sl.Err(err), "id", id)
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Comment retrieved", "id", id)
	return comment, nil
}
