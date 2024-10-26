package graph

import (
	"github.com/oreshkin/posts/graph/model"
	"github.com/oreshkin/posts/internal/comments"
	"github.com/oreshkin/posts/internal/posts"
)

type Resolver struct {
	PostRepository    posts.PostRepository
	CommentRepository comments.CommentRepository

	CommentChannel chan *model.Comment
}
