// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	models "comments-system/internal/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CommentService is an autogenerated mock type for the CommentService type
type CommentService struct {
	mock.Mock
}

// CreateComment provides a mock function with given fields: ctx, input
func (_m *CommentService) CreateComment(ctx context.Context, input models.CreateCommentInput) (models.Comment, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for CreateComment")
	}

	var r0 models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateCommentInput) (models.Comment, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateCommentInput) models.Comment); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(models.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.CreateCommentInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetComment provides a mock function with given fields: ctx, id
func (_m *CommentService) GetComment(ctx context.Context, id string) (models.Comment, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetComment")
	}

	var r0 models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.Comment, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Comment); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCommentReplies provides a mock function with given fields: ctx, parentID
func (_m *CommentService) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	ret := _m.Called(ctx, parentID)

	if len(ret) == 0 {
		panic("no return value specified for GetCommentReplies")
	}

	var r0 []models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]models.Comment, error)); ok {
		return rf(ctx, parentID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.Comment); ok {
		r0 = rf(ctx, parentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, parentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetComments provides a mock function with given fields: ctx, postID, limit, offset
func (_m *CommentService) GetComments(ctx context.Context, postID string, limit int, offset int) ([]models.Comment, int, error) {
	ret := _m.Called(ctx, postID, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetComments")
	}

	var r0 []models.Comment
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]models.Comment, int, error)); ok {
		return rf(ctx, postID, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []models.Comment); ok {
		r0 = rf(ctx, postID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int); ok {
		r1 = rf(ctx, postID, limit, offset)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, postID, limit, offset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewCommentService creates a new instance of CommentService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommentService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentService {
	mock := &CommentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
