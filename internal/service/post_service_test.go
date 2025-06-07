package service_test

import (
	"comments-system/internal/models"
	"comments-system/internal/service"
	"comments-system/internal/storage/mocks"
	"comments-system/pkg/errors"
	"comments-system/pkg/logger/slogdiscard"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostService_CreatePost_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewPostService(storageMock, log)

	input := models.CreatePostInput{
		Title:   "Test Post",
		Content: "Content",
		Author:  "author1",
	}

	storageMock.On("CreatePost", mock.Anything, mock.Anything).Return(models.Post{
		ID:      "post1",
		Title:   "Test Post",
		Content: "Content",
		Author:  "author1",
	}, nil)

	post, err := svc.CreatePost(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, "post1", post.ID)
	storageMock.AssertExpectations(t)
}

func TestPostService_ToggleComments_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewPostService(storageMock, log)

	originalPost := models.Post{
		ID:              "post1",
		CommentsEnabled: false,
	}

	storageMock.On("GetPost", mock.Anything, "post1").Return(originalPost, nil)
	storageMock.On("UpdatePost", mock.Anything, mock.MatchedBy(func(p models.Post) bool {
		return p.CommentsEnabled == true
	})).Return(nil)

	updated, err := svc.ToggleComments(context.Background(), "post1", true)

	assert.NoError(t, err)
	assert.True(t, updated.CommentsEnabled)
	storageMock.AssertExpectations(t)
}

func TestPostService_ToggleComments_NoChange(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewPostService(storageMock, log)

	post := models.Post{
		ID:              "post1",
		CommentsEnabled: true,
	}

	storageMock.On("GetPost", mock.Anything, "post1").Return(post, nil)

	result, err := svc.ToggleComments(context.Background(), "post1", true)

	assert.NoError(t, err)
	assert.True(t, result.CommentsEnabled)
	storageMock.AssertNotCalled(t, "UpdatePost")
}

func TestPostService_GetPost_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewPostService(storageMock, log)

	postID := "post1"
	expectedPost := models.Post{ID: postID, Title: "Test Post"}

	storageMock.On("GetPost", mock.Anything, postID).Return(expectedPost, nil)

	post, err := svc.GetPost(context.Background(), postID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	storageMock.AssertExpectations(t)
}

func TestPostService_GetPosts_Success(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewPostService(storageMock, log)

	limit := 10
	offset := 0
	expectedPosts := []models.Post{
		{ID: "post1", Title: "Test Post 1"},
		{ID: "post2", Title: "Test Post 2"},
	}

	storageMock.On("GetPosts", mock.Anything, limit, offset).Return(expectedPosts, nil)

	posts, err := svc.GetPosts(context.Background(), limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	storageMock.AssertExpectations(t)
}

func TestPostService_GetPost_NotFound(t *testing.T) {
	storageMock := &mocks.Storage{}
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewPostService(storageMock, log)

	postID := "post1"

	storageMock.On("GetPost", mock.Anything, postID).Return(models.Post{}, errors.ErrNotFound)

	_, err := svc.GetPost(context.Background(), postID)

	assert.ErrorIs(t, err, errors.ErrNotFound)
	storageMock.AssertExpectations(t)
}
