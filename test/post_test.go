package graph_test

import (
	"context"
	"testing"

	"github.com/oreshkin/posts/graph"
	"github.com/oreshkin/posts/graph/model"
	"github.com/oreshkin/posts/internal/auth"
	"github.com/oreshkin/posts/internal/comments"
	"github.com/oreshkin/posts/internal/posts"
	"github.com/oreshkin/posts/internal/users"
	"github.com/oreshkin/posts/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePostResolver(t *testing.T) {
	ctx := context.Background()
	mockUser := &users.User{
		ID:       "123",
		Username: "testuser",
	}
	ctx = context.WithValue(ctx, auth.UserCtxKey, mockUser)

	mockPostRepository := new(mocks.MockPostRepository)

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

func TestPostsResolver(t *testing.T) {
	ctx := context.Background()

	mockPostRepository := new(mocks.MockPostRepository)
	mockPosts := []posts.Post{
		{
			ID:               1,
			Title:            "Тестовый пост 1",
			Content:          "Контент тестового поста 1",
			CommentsDisabled: false,
			Comments: []comments.Comment{
				{ID: 1, Text: "Комментарий 1 к посту 1", ParentCommentID: nil},
				{ID: 2, Text: "Комментарий 2 к посту 1", ParentCommentID: nil},
			},
		},
		{
			ID:               2,
			Title:            "Тестовый пост 2",
			Content:          "Контент тестового поста 2",
			CommentsDisabled: false,
			Comments: []comments.Comment{
				{ID: 3, Text: "Комментарий 1 к посту 2", ParentCommentID: nil},
			},
		},
	}

	mockPostRepository.On("GetPostsWithPagination", 10, (*int64)(nil)).Return(mockPosts, (*int64)(nil), nil)

	resolver := &graph.Resolver{
		PostRepository: mockPostRepository,
	}

	result, err := resolver.Query().Posts(ctx, nil, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, result.Edges, 2)

	assert.Equal(t, "1", result.Edges[0].Cursor)
	assert.Equal(t, "Тестовый пост 1", result.Edges[0].Node.Title)
	assert.Equal(t, "Контент тестового поста 1", result.Edges[0].Node.Content)

	assert.Len(t, result.Edges[0].Node.Comments.Edges, 2)
	assert.Equal(t, "Комментарий 1 к посту 1", result.Edges[0].Node.Comments.Edges[0].Node.Text)
	assert.Equal(t, "Комментарий 2 к посту 1", result.Edges[0].Node.Comments.Edges[1].Node.Text)

	assert.Equal(t, "2", result.Edges[1].Cursor)
	assert.Equal(t, "Тестовый пост 2", result.Edges[1].Node.Title)
	assert.Equal(t, "Контент тестового поста 2", result.Edges[1].Node.Content)

	assert.Len(t, result.Edges[1].Node.Comments.Edges, 1)
	assert.Equal(t, "Комментарий 1 к посту 2", result.Edges[1].Node.Comments.Edges[0].Node.Text)

	assert.False(t, result.PageInfo.HasNextPage)
	assert.Equal(t, "", result.PageInfo.EndCursor)

	mockPostRepository.AssertExpectations(t)
}

func TestPostResolver(t *testing.T) {
	ctx := context.Background()

	mockPost := &posts.Post{
		ID:               1,
		Title:            "Новый пост",
		Content:          "Содержимое нового поста",
		CommentsDisabled: true,
	}

	mockComments := []comments.Comment{
		{
			ID:        5,
			PostID:    1,
			Text:      "Тестовый комментарий",
			CreatedAt: "2024-10-20T07:12:29.210862Z",
		},
		{
			ID:        3,
			PostID:    1,
			Text:      "Тестовый комментарий",
			CreatedAt: "2024-10-20T06:43:55.252381Z",
		},
		{
			ID:        2,
			PostID:    1,
			Text:      "Тестовый комментарий",
			CreatedAt: "2024-10-20T06:42:08.396116Z",
		},
	}

	mockPostRepository := new(mocks.MockPostRepository)
	mockCommentRepository := new(mocks.MockCommentRepository)

	mockPostRepository.On("GetPostByID", "1").Return(mockPost, nil)
	mockCommentRepository.On("GetCommentsByPostIDWithPagination", int64(1), (*int64)(nil), 10, (*string)(nil)).Return(mockComments, nil)

	resolver := &graph.Resolver{
		PostRepository:    mockPostRepository,
		CommentRepository: mockCommentRepository,
	}

	result, err := resolver.Query().Post(ctx, "1", nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "1", result.ID)
	assert.Equal(t, "Новый пост", result.Title)
	assert.Equal(t, "Содержимое нового поста", result.Content)
	assert.True(t, result.CommentsDisabled)

	assert.Len(t, result.Comments.Edges, 3)
	assert.Equal(t, "5", result.Comments.Edges[0].Cursor)
	assert.Equal(t, "Тестовый комментарий", result.Comments.Edges[0].Node.Text)
	assert.Equal(t, "2024-10-20T07:12:29.210862Z", result.Comments.Edges[0].Node.CreatedAt)

	assert.Equal(t, "3", result.Comments.Edges[1].Cursor)
	assert.Equal(t, "2", result.Comments.PageInfo.EndCursor)
	assert.False(t, result.Comments.PageInfo.HasNextPage)

	mockPostRepository.AssertExpectations(t)
	mockCommentRepository.AssertExpectations(t)
}
