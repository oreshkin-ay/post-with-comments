package comments

import (
	"database/sql"
	"log"

	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
)

type Comment struct {
	ID                int64
	PostID            int64
	Text              string
	ParentCommentID   *int64
	CreatedAt         string
	ChildCommentCount int
}

func (comment Comment) Save(userID string) int64 {
	stmt, err := database.Db.Prepare("INSERT INTO comments (post_id, text, parent_comment_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(comment.PostID, comment.Text, comment.ParentCommentID, userID).Scan(&id)
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
				`SELECT id, post_id, text, parent_comment_id, created_at,
				 (SELECT COUNT(*) FROM comments AS child_comments WHERE child_comments.parent_comment_id = comments.id) AS child_comment_count
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id IS NULL AND id < $2
				 ORDER BY id DESC
				 LIMIT $3`, postID, *cursor, limit)
		} else {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at,
				 (SELECT COUNT(*) FROM comments AS child_comments WHERE child_comments.parent_comment_id = comments.id) AS child_comment_count
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id IS NULL
				 ORDER BY id DESC
				 LIMIT $2`, postID, limit)
		}
	} else {
		if cursor != nil {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at,
				 (SELECT COUNT(*) FROM comments AS child_comments WHERE child_comments.parent_comment_id = comments.id) AS child_comment_count
				 FROM comments
				 WHERE post_id = $1 AND parent_comment_id = $2 AND id < $3
				 ORDER BY id DESC
				 LIMIT $4`, postID, *parentCommentID, *cursor, limit)
		} else {
			rows, err = database.Db.Query(
				`SELECT id, post_id, text, parent_comment_id, created_at,
				 (SELECT COUNT(*) FROM comments AS child_comments WHERE child_comments.parent_comment_id = comments.id) AS child_comment_count
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
		var childCommentCount int
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.Text, &comment.ParentCommentID, &comment.CreatedAt, &childCommentCount)
		if err != nil {
			return nil, err
		}

		comment.ChildCommentCount = childCommentCount
		comments = append(comments, comment)
	}

	return comments, nil
}
