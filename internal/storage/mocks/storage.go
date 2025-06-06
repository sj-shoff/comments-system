// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	models "comments-system/internal/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *Storage) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountCommentsByPost provides a mock function with given fields: ctx, postID
func (_m *Storage) CountCommentsByPost(ctx context.Context, postID string) (int, error) {
	ret := _m.Called(ctx, postID)

	if len(ret) == 0 {
		panic("no return value specified for CountCommentsByPost")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, postID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, postID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateComment provides a mock function with given fields: ctx, comment
func (_m *Storage) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	ret := _m.Called(ctx, comment)

	if len(ret) == 0 {
		panic("no return value specified for CreateComment")
	}

	var r0 models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Comment) (models.Comment, error)); ok {
		return rf(ctx, comment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Comment) models.Comment); ok {
		r0 = rf(ctx, comment)
	} else {
		r0 = ret.Get(0).(models.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Comment) error); ok {
		r1 = rf(ctx, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePost provides a mock function with given fields: ctx, post
func (_m *Storage) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	ret := _m.Called(ctx, post)

	if len(ret) == 0 {
		panic("no return value specified for CreatePost")
	}

	var r0 models.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Post) (models.Post, error)); ok {
		return rf(ctx, post)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Post) models.Post); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Get(0).(models.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Post) error); ok {
		r1 = rf(ctx, post)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetComment provides a mock function with given fields: ctx, id
func (_m *Storage) GetComment(ctx context.Context, id string) (models.Comment, error) {
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
func (_m *Storage) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
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

// GetCommentsByPost provides a mock function with given fields: ctx, postID, limit, offset
func (_m *Storage) GetCommentsByPost(ctx context.Context, postID string, limit int, offset int) ([]models.Comment, error) {
	ret := _m.Called(ctx, postID, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetCommentsByPost")
	}

	var r0 []models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]models.Comment, error)); ok {
		return rf(ctx, postID, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []models.Comment); ok {
		r0 = rf(ctx, postID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, postID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPost provides a mock function with given fields: ctx, id
func (_m *Storage) GetPost(ctx context.Context, id string) (models.Post, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetPost")
	}

	var r0 models.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.Post, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Post); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPosts provides a mock function with given fields: ctx, limit, offset
func (_m *Storage) GetPosts(ctx context.Context, limit int, offset int) ([]models.Post, error) {
	ret := _m.Called(ctx, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetPosts")
	}

	var r0 []models.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) ([]models.Post, error)); ok {
		return rf(ctx, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []models.Post); ok {
		r0 = rf(ctx, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePost provides a mock function with given fields: ctx, post
func (_m *Storage) UpdatePost(ctx context.Context, post models.Post) error {
	ret := _m.Called(ctx, post)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePost")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
