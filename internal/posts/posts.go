package posts

import (
	"log"

	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
)

type Post struct {
	ID               int64
	Title            string
	Content          string
	CommentsDisabled bool
}

func (post Post) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO posts (title, content, comments_disabled) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(post.Title, post.Content, post.CommentsDisabled).Scan(&id)
	if err != nil {
		log.Fatal("Error executing statement:", err)
	}

	log.Print("Post inserted with ID:", id)
	return id
}
