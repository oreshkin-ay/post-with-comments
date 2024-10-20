package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"
	"strconv"

	"github.com/oreshkin/posts/graph/model"
	"github.com/oreshkin/posts/internal/auth"
	"github.com/oreshkin/posts/internal/comments"
	"github.com/oreshkin/posts/internal/pkg/jwt"
	"github.com/oreshkin/posts/internal/posts"
	"github.com/oreshkin/posts/internal/users"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPostInput) (*model.Post, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Post{}, fmt.Errorf("access denied")
	}

	var post posts.Post
	post.Title = input.Title
	post.Content = input.Content
	post.CommentsDisabled = input.CommentsDisabled

	postID := r.PostRepository.Save(post, user.ID)

	graphqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}

	return &model.Post{
		ID:               strconv.FormatInt(postID, 10),
		Title:            post.Title,
		Content:          post.Content,
		CommentsDisabled: post.CommentsDisabled,
		User:             graphqlUser,
	}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewCommentInput) (*model.Comment, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Comment{}, fmt.Errorf("access denied")
	}

	postIDInt, err := strconv.ParseInt(input.PostID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %v", err)
	}

	post, err := r.PostRepository.GetPostByID(input.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post with ID %s not found", input.PostID)
	}

	if post.CommentsDisabled {
		return nil, fmt.Errorf("comments are disabled for this post")
	}

	var comment comments.Comment

	comment.PostID = postIDInt

	comment.Text = input.Text

	if input.ParentCommentID != nil {
		parentIDInt, err := strconv.ParseInt(*input.ParentCommentID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID: %v", err)
		}
		comment.ParentCommentID = &parentIDInt
	} else {
		comment.ParentCommentID = nil
	}

	commentID := r.CommentRepository.Save(comment, user.ID)

	var parentIDStr *string
	if comment.ParentCommentID != nil {
		parentIDStr = new(string)
		*parentIDStr = strconv.FormatInt(*comment.ParentCommentID, 10)
	}

	return &model.Comment{
		ID:              strconv.FormatInt(commentID, 10),
		PostID:          strconv.FormatInt(comment.PostID, 10),
		Text:            comment.Text,
		ParentCommentID: parentIDStr,
	}, nil
}

// UpdateCommentsDisabled is the resolver for the updateCommentsDisabled field.
func (r *mutationResolver) UpdateCommentsDisabled(ctx context.Context, input model.UpdateCommentsDisabledInput) (*model.Post, error) {
	post, err := r.PostRepository.GetPostByID(input.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post with ID %s not found", input.PostID)
	}

	err = r.PostRepository.UpdateCommentsDisabled(input.PostID, input.CommentsDisabled)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	post.CommentsDisabled = input.CommentsDisabled
	return &model.Post{
		ID:               strconv.FormatInt(post.ID, 10),
		Title:            post.Title,
		Content:          post.Content,
		CommentsDisabled: post.CommentsDisabled,
	}, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, cursor *string, limit *int) (*model.PostConnection, error) {
	defaultLimit := 10
	if limit != nil {
		defaultLimit = *limit
	}

	var parsedCursor *int64
	if cursor != nil {
		parsedCursorValue, err := strconv.ParseInt(*cursor, 10, 64)
		if err != nil {
			return &model.PostConnection{
				Edges:    []*model.PostEdge{},
				PageInfo: &model.PageInfo{EndCursor: "", HasNextPage: false},
			}, nil
		}
		parsedCursor = &parsedCursorValue
	}

	posts, lastPostID, err := r.PostRepository.GetPostsWithPagination(defaultLimit, parsedCursor)
	if err != nil {
		return &model.PostConnection{
			Edges:    []*model.PostEdge{},
			PageInfo: &model.PageInfo{EndCursor: "", HasNextPage: false},
		}, nil
	}

	var edges []*model.PostEdge
	for _, post := range posts {
		edges = append(edges, &model.PostEdge{
			Cursor: strconv.FormatInt(post.ID, 10),
			Node: &model.Post{
				ID:               strconv.FormatInt(post.ID, 10),
				Title:            post.Title,
				Content:          post.Content,
				CommentsDisabled: post.CommentsDisabled,
			},
		})
	}

	hasNextPage := len(posts) == defaultLimit

	var endCursor string
	if lastPostID != nil {
		endCursor = strconv.FormatInt(*lastPostID, 10)
	} else {
		endCursor = ""
	}

	return &model.PostConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string, parentCommentID *string, cursor *string, limit *int) (*model.Post, error) {
	defaultLimit := 10
	if limit != nil {
		defaultLimit = *limit
	}

	post, err := r.PostRepository.GetPostByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post by ID: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post with ID %s not found", id)
	}

	var parsedCursor *int64
	if cursor != nil {
		parsedCursorValue, err := strconv.ParseInt(*cursor, 10, 64)
		if err != nil {
			return comments.EmptyPostWithNoComments(post.ID, post.Title, post.Content, post.CommentsDisabled), nil
		}
		parsedCursor = &parsedCursorValue
	}

	fetchedComments, err := r.CommentRepository.GetCommentsByPostIDWithPagination(post.ID, parsedCursor, defaultLimit, parentCommentID)
	if err != nil || len(fetchedComments) == 0 {
		return comments.EmptyPostWithNoComments(post.ID, post.Title, post.Content, post.CommentsDisabled), nil
	}

	var commentEdges []*model.CommentEdge
	for _, comment := range fetchedComments {
		var parentIDStr *string
		if comment.ParentCommentID != nil {
			parentIDStr = new(string)
			*parentIDStr = strconv.FormatInt(*comment.ParentCommentID, 10)
		}

		commentEdges = append(commentEdges, &model.CommentEdge{
			Cursor: strconv.FormatInt(comment.ID, 10),
			Node: &model.Comment{
				ID:              strconv.FormatInt(comment.ID, 10),
				PostID:          strconv.FormatInt(comment.PostID, 10),
				Text:            comment.Text,
				ParentCommentID: parentIDStr,
				CreatedAt:       comment.CreatedAt,
			},
		})
	}

	var lastCommentCursor *string
	if len(fetchedComments) > 0 {
		lastCommentID := fetchedComments[len(fetchedComments)-1].ID
		lastCommentCursorStr := strconv.FormatInt(lastCommentID, 10)
		lastCommentCursor = &lastCommentCursorStr
	}

	hasNextPage := len(fetchedComments) == defaultLimit

	commentConnection := &model.CommentConnection{
		Edges: commentEdges,
		PageInfo: &model.PageInfo{
			EndCursor:   comments.LastCommentCursorOrDefault(lastCommentCursor),
			HasNextPage: hasNextPage,
		},
	}

	return &model.Post{
		ID:               strconv.FormatInt(post.ID, 10),
		Title:            post.Title,
		Content:          post.Content,
		CommentsDisabled: post.CommentsDisabled,
		Comments:         commentConnection,
	}, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
