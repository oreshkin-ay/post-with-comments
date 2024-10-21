package posts

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/oreshkin/posts/internal/comments"
	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
	"github.com/oreshkin/posts/internal/users"
)

type Post struct {
	ID               int64
	Title            string
	Content          string
	CommentsDisabled bool
	Comments         []comments.Comment
	User             *users.User
}

// Save saves a post to the database and returns the post's ID
func (post Post) Save(userID string) (int64, error) {
	stmt, err := database.Db.Prepare("INSERT INTO posts (title, content, comments_disabled, user_id) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(post.Title, post.Content, post.CommentsDisabled, userID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error executing statement: %w", err)
	}

	log.Printf("Post inserted with ID: %d", id)
	return id, nil
}

// GetPostByID returns a post by its ID
func GetPostByID(postID string) (*Post, error) {
	var post Post

	stmt, err := database.Db.Prepare("SELECT id, title, content, comments_disabled FROM posts WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(postID).Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	return &post, nil
}

// GetPostsWithPagination returns a list of posts with pagination and a cursor
func GetPostsWithPagination(limit int, cursor *int64) ([]Post, *int64, error) {
	var rows *sql.Rows
	var err error

	if cursor != nil {
		rows, err = database.Db.Query(
			`SELECT id, title, content, comments_disabled
             FROM posts
             WHERE id < $1
             ORDER BY id DESC
             LIMIT $2`, cursor, limit)
	} else {
		rows, err = database.Db.Query(
			`SELECT id, title, content, comments_disabled
             FROM posts
             ORDER BY id DESC
             LIMIT $1`, limit)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var posts []Post
	var lastPostID *int64

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %w", err)
		}

		posts = append(posts, post)
		lastPostID = &post.ID
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return posts, lastPostID, nil
}

// UpdateCommentsDisabled updates the commentsDisabled flag for a post
func UpdateCommentsDisabled(postID string, commentsDisabled bool) error {
	stmt, err := database.Db.Prepare("UPDATE posts SET comments_disabled = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("error preparing update statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(commentsDisabled, postID)
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	return nil
}
