package posts

import "fmt"

type PostRepository interface {
	Save(post Post, userID string) (int64, error)
	GetPostByID(postID string) (*Post, error)
	GetPostsWithPagination(limit int, cursor *int64, commentsLimit *int) ([]Post, *int64, error)
	UpdateCommentsDisabled(postID string, commentsDisabled bool) error
}

type DBPostRepository struct{}

func (r *DBPostRepository) Save(post Post, userID string) (int64, error) {
	id, err := post.Save(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to save post: %w", err)
	}
	return id, nil
}

func (r *DBPostRepository) GetPostByID(postID string) (*Post, error) {
	return GetPostByID(postID)
}

func (r *DBPostRepository) GetPostsWithPagination(limit int, cursor *int64, commentsLimit *int) ([]Post, *int64, error) {
	return GetPostsWithPagination(limit, cursor, commentsLimit)
}

func (r *DBPostRepository) UpdateCommentsDisabled(postID string, commentsDisabled bool) error {
	return UpdateCommentsDisabled(postID, commentsDisabled)
}
