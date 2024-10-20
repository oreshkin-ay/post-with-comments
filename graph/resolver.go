package graph

import (
	"github.com/oreshkin/posts/internal/comments"
	"github.com/oreshkin/posts/internal/posts"
)

type Resolver struct {
	PostRepository    posts.PostRepository
	CommentRepository comments.CommentRepository
}
