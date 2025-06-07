package service_test

import (
	"comments-system/internal/models"
	"comments-system/internal/service"
	"comments-system/internal/storage/mocks"
	"comments-system/pkg/errors"
	"comments-system/pkg/logger/slogdiscard"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_CreateComment_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	input := models.CreateCommentInput{
		PostID:  "post1",
		Author:  "user1",
		Content: "Valid content",
	}

	storageMock.On("GetPost", mock.Anything, "post1").Return(models.Post{
		ID:              "post1",
		CommentsEnabled: true,
	}, nil)
	storageMock.On("CreateComment", mock.Anything, mock.Anything).Return(models.Comment{
		ID:        "comment1",
		PostID:    "post1",
		Author:    "user1",
		Content:   "Valid content",
		CreatedAt: time.Now(),
	}, nil)

	comment, err := svc.CreateComment(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, "comment1", comment.ID)
	storageMock.AssertExpectations(t)
}

func TestCommentService_CreateComment_CommentsDisabled(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	input := models.CreateCommentInput{
		PostID:  "post1",
		Author:  "user1",
		Content: "Invalid content",
	}

	storageMock.On("GetPost", mock.Anything, "post1").Return(models.Post{
		CommentsEnabled: false,
	}, nil)

	_, err := svc.CreateComment(context.Background(), input)

	assert.ErrorIs(t, err, errors.ErrCommentsDisabled)
	storageMock.AssertExpectations(t)
}

func TestCommentService_GetComments_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	comments := []models.Comment{
		{ID: "comment1", PostID: "post1"},
		{ID: "comment2", PostID: "post1"},
	}

	storageMock.On("GetCommentsByPost", mock.Anything, "post1", 10, 0).Return(comments, nil)
	storageMock.On("CountCommentsByPost", mock.Anything, "post1").Return(2, nil)

	result, total, err := svc.GetComments(context.Background(), "post1", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
	storageMock.AssertExpectations(t)
}

func TestCommentService_GetComment_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	commentID := "comment1"
	expectedComment := models.Comment{ID: commentID, PostID: "post1", Content: "Test comment"}

	storageMock.On("GetComment", mock.Anything, commentID).Return(expectedComment, nil)

	comment, err := svc.GetComment(context.Background(), commentID)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	storageMock.AssertExpectations(t)
}

func TestCommentService_GetCommentReplies_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	parentID := "comment1"
	expectedReplies := []models.Comment{
		{ID: "reply1", PostID: "post1"},
		{ID: "reply2", PostID: "post1"},
	}

	storageMock.On("GetCommentReplies", mock.Anything, parentID).Return(expectedReplies, nil)

	replies, err := svc.GetCommentReplies(context.Background(), parentID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReplies, replies)
	storageMock.AssertExpectations(t)
}

func TestCommentService_GetComment_NotFound(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	commentID := "comment1"

	storageMock.On("GetComment", mock.Anything, commentID).Return(models.Comment{}, errors.ErrNotFound)

	_, err := svc.GetComment(context.Background(), commentID)

	assert.ErrorIs(t, err, errors.ErrNotFound)
	storageMock.AssertExpectations(t)
}

func TestCommentService_GetCommentReplies_NoReplies(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewCommentService(storageMock, log)

	parentID := "comment1"

	storageMock.On("GetCommentReplies", mock.Anything, parentID).Return([]models.Comment{}, nil)

	replies, err := svc.GetCommentReplies(context.Background(), parentID)

	assert.NoError(t, err)
	assert.Empty(t, replies)
	storageMock.AssertExpectations(t)
}
