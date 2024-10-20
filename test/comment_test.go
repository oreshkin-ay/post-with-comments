package graph_test

import (
	"context"
	"testing"

	"github.com/oreshkin/posts/graph"
	"github.com/oreshkin/posts/graph/model"
	"github.com/oreshkin/posts/internal/auth"
	"github.com/oreshkin/posts/internal/posts"
	"github.com/oreshkin/posts/internal/users"
	"github.com/oreshkin/posts/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCommentWithoutParentID(t *testing.T) {
	ctx := context.Background()
	mockUser := &users.User{
		ID:       "123",
		Username: "testuser",
	}
	ctx = context.WithValue(ctx, auth.UserCtxKey, mockUser)

	mockPostRepository := new(mocks.MockPostRepository)
	mockCommentRepository := new(mocks.MockCommentRepository)

	mockPost := &posts.Post{
		ID:               1,
		Title:            "Первый пост",
		Content:          "Содержимое первого поста.",
		CommentsDisabled: false,
	}
	mockPostRepository.On("GetPostByID", "1").Return(mockPost, nil)

	mockCommentRepository.On("Save", mock.AnythingOfType("comments.Comment"), mockUser.ID).Return(int64(1))

	input := model.NewCommentInput{
		PostID:          "1",
		Text:            "Тестовый комментарий",
		ParentCommentID: nil,
	}

	resolver := &graph.Resolver{
		PostRepository:    mockPostRepository,
		CommentRepository: mockCommentRepository,
	}

	comment, err := resolver.Mutation().CreateComment(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "1", comment.PostID)
	assert.Equal(t, "Тестовый комментарий", comment.Text)
	assert.Nil(t, comment.ParentCommentID)
}

func TestCreateCommentWithParentID(t *testing.T) {
	ctx := context.Background()
	mockUser := &users.User{
		ID:       "123",
		Username: "testuser",
	}
	ctx = context.WithValue(ctx, auth.UserCtxKey, mockUser)

	mockPostRepository := new(mocks.MockPostRepository)
	mockCommentRepository := new(mocks.MockCommentRepository)

	mockPost := &posts.Post{
		ID:               1,
		Title:            "Первый пост",
		Content:          "Содержимое первого поста.",
		CommentsDisabled: false,
	}
	mockPostRepository.On("GetPostByID", "1").Return(mockPost, nil)

	mockCommentRepository.On("Save", mock.AnythingOfType("comments.Comment"), mockUser.ID).Return(int64(1))

	parentCommentID := "2"
	input := model.NewCommentInput{
		PostID:          "1",
		Text:            "Ответ на родительский комментарий",
		ParentCommentID: &parentCommentID,
	}

	resolver := &graph.Resolver{
		PostRepository:    mockPostRepository,
		CommentRepository: mockCommentRepository,
	}

	comment, err := resolver.Mutation().CreateComment(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "1", comment.PostID)
	assert.Equal(t, "Ответ на родительский комментарий", comment.Text)
	assert.NotNil(t, comment.ParentCommentID)
	assert.Equal(t, parentCommentID, *comment.ParentCommentID)

	mockCommentRepository.AssertCalled(t, "Save", mock.AnythingOfType("comments.Comment"), mockUser.ID)
}
