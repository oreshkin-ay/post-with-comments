package mocks

import (
	"github.com/oreshkin/posts/internal/comments"
	"github.com/oreshkin/posts/internal/posts"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Save(post posts.Post, userID string) (int64, error) {
	args := m.Called(post, userID)
	return args.Get(0).(int64), nil
}

func (m *MockPostRepository) GetPostByID(postID string) (*posts.Post, error) {
	args := m.Called(postID)
	return args.Get(0).(*posts.Post), args.Error(1)
}

func (m *MockPostRepository) GetPostsWithPagination(limit int, cursor *int64) ([]posts.Post, *int64, error) {
	args := m.Called(limit, cursor)
	return args.Get(0).([]posts.Post), args.Get(1).(*int64), args.Error(2)
}

func (m *MockPostRepository) UpdateCommentsDisabled(postID string, commentsDisabled bool) error {
	args := m.Called(postID, commentsDisabled)
	return args.Error(0)
}

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Save(comment comments.Comment, userID string) (int64, error) {
	args := m.Called(comment, userID)
	return args.Get(0).(int64), nil
}

func (m *MockCommentRepository) GetCommentsByPostIDWithPagination(postID int64, cursor *int64, limit int, parentCommentID *string) ([]comments.Comment, error) {
	args := m.Called(postID, cursor, limit, parentCommentID)
	return args.Get(0).([]comments.Comment), args.Error(1)
}
