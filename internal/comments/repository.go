package comments

import "fmt"

type CommentRepository interface {
	Save(comment Comment, userID string) (int64, error)
	GetCommentsByPostIDWithPagination(postID int64, cursor *int64, limit int, parentCommentID *string) ([]Comment, error)
}

type DBCommentRepository struct{}

func (r *DBCommentRepository) Save(comment Comment, userID string) (int64, error) {
	id, err := comment.Save(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to save comment: %w", err)
	}
	return id, nil
}

func (r *DBCommentRepository) GetCommentsByPostIDWithPagination(postID int64, cursor *int64, limit int, parentCommentID *string) ([]Comment, error) {
	return GetCommentsByPostIDWithPagination(postID, cursor, limit, parentCommentID)
}
