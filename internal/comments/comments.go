package comments

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/oreshkin/posts/graph/model"
	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
)

type Comment struct {
	ID              int64
	PostID          int64
	Text            string
	ParentCommentID *int64
	CreatedAt       string
}

func (comment Comment) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO comments (post_id, text, parent_comment_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(comment.PostID, comment.Text, comment.ParentCommentID).Scan(&id)
	if err != nil {
		log.Fatal("Error executing statement:", err)
	}

	log.Print("Comment inserted with ID:", id)
	return id
}

func GetCommentsByPostIDWithPagination(postID int64, cursor *int64, limit int, parentCommentID *string) ([]Comment, error) {
	var rows *sql.Rows
	var err error

	if parentCommentID == nil {
		if cursor != nil {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id IS NULL AND id < $2
				 ORDER BY id DESC
				 LIMIT $3`, postID, *cursor, limit)
		} else {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id IS NULL
				 ORDER BY id DESC
				 LIMIT $2`, postID, limit)
		}
	} else {
		if cursor != nil {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id = $2 AND id < $3
				 ORDER BY id DESC
				 LIMIT $4`, postID, *parentCommentID, *cursor, limit)
		} else {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id = $2
				 ORDER BY id DESC
				 LIMIT $3`, postID, *parentCommentID, limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.Text, &comment.ParentCommentID, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)

	}

	return comments, nil
}

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
