package posts

import (
	"database/sql"
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

func (post Post) Save(userID string) int64 {
	stmt, err := database.Db.Prepare("INSERT INTO posts (title, content, comments_disabled, user_id) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(post.Title, post.Content, post.CommentsDisabled, userID).Scan(&id)
	if err != nil {
		log.Fatal("Error executing statement:", err)
	}

	log.Print("Post inserted with ID:", id)
	return id
}

func GetPostByID(postID string) (*Post, error) {
	var post Post

	stmt, err := database.Db.Prepare("SELECT id, title, content, comments_disabled FROM posts WHERE id = $1")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(postID).Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Fatal("Error executing statement:", err)
		return nil, err
	}

	return &post, nil
}

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
		// If no cursor is provided, fetch the latest posts
		rows, err = database.Db.Query(
			`SELECT id, title, content, comments_disabled
             FROM posts
             ORDER BY id DESC
             LIMIT $1`, limit)
	}

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var posts []Post
	var lastPostID *int64

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
		if err != nil {
			return nil, nil, err
		}

		posts = append(posts, post)
		lastPostID = &post.ID
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return posts, lastPostID, nil
}
