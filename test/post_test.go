package graph_test

import (
	"context"
	"testing"

	"github.com/oreshkin/posts/graph"
	"github.com/oreshkin/posts/graph/model"
	"github.com/oreshkin/posts/internal/auth"
	"github.com/oreshkin/posts/internal/posts"
	"github.com/oreshkin/posts/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Save(post posts.Post, userID string) int64 {
	args := m.Called(post, userID)
	return args.Get(0).(int64)
}

func (m *MockPostRepository) GetPostByID(postID string) (*posts.Post, error) {
	return nil, nil
}

func (m *MockPostRepository) GetPostsWithPagination(limit int, cursor *int64) ([]posts.Post, *int64, error) {
	return nil, nil, nil
}

func (m *MockPostRepository) UpdateCommentsDisabled(postID string, commentsDisabled bool) error {
	return nil
}

func TestCreatePost(t *testing.T) {
	ctx := context.Background()
	mockUser := &users.User{
		ID:       "123",
		Username: "testuser",
	}
	ctx = context.WithValue(ctx, auth.UserCtxKey, mockUser)

	mockPostRepository := new(MockPostRepository)

	mockPostRepository.On("Save", mock.AnythingOfType("posts.Post"), mockUser.ID).Return(int64(1))

	input := model.NewPostInput{
		Title:            "Первый пост",
		Content:          "Содержимое первого поста.",
		CommentsDisabled: false,
	}

	resolver := &graph.Resolver{
		PostRepository: mockPostRepository,
	}

	post, err := resolver.Mutation().CreatePost(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "Первый пост", post.Title)
	assert.Equal(t, "Содержимое первого поста.", post.Content)
	assert.Equal(t, false, post.CommentsDisabled)

	mockPostRepository.AssertCalled(t, "Save", mock.AnythingOfType("posts.Post"), mockUser.ID)
}
