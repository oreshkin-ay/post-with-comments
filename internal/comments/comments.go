package comments

import (
	"log"

	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
)

type Comment struct {
	ID        int64
	PostID    int64
	Text      string
	ParentID  *int64
	CreatedAt string
}

func (comment Comment) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO comments (post_id, text, parent_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(comment.PostID, comment.Text, comment.ParentID).Scan(&id)
	if err != nil {
		log.Fatal("Error executing statement:", err)
	}

	log.Print("Comment inserted with ID:", id)
	return id
}
