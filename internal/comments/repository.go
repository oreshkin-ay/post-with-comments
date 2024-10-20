package comments

type CommentRepository interface {
	Save(comment Comment, userID string) int64
	GetCommentsByPostIDWithPagination(postID int64, cursor *int64, limit int, parentCommentID *string) ([]Comment, error)
}

type DBCommentRepository struct{}

func (r *DBCommentRepository) Save(comment Comment, userID string) int64 {
	return comment.Save(userID)
}

func (r *DBCommentRepository) GetCommentsByPostIDWithPagination(postID int64, cursor *int64, limit int, parentCommentID *string) ([]Comment, error) {
	return GetCommentsByPostIDWithPagination(postID, cursor, limit, parentCommentID)
}
