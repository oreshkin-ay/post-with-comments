package comments

import (
	"strconv"

	"github.com/oreshkin/posts/graph/model"
)

func LastCommentCursorOrDefault(cursor *string) string {
	if cursor != nil {
		return *cursor
	}
	return ""
}

func EmptyPostWithNoComments(postID int64, title, content string, commentsDisabled bool) *model.Post {
	return &model.Post{
		ID:               strconv.FormatInt(postID, 10),
		Title:            title,
		Content:          content,
		CommentsDisabled: commentsDisabled,
		Comments: &model.CommentConnection{
			Edges:    []*model.CommentEdge{},
			PageInfo: &model.PageInfo{EndCursor: "", HasNextPage: false},
		},
	}
}
