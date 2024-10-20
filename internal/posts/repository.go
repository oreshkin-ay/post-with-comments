package posts

type PostRepository interface {
	Save(post Post, userID string) int64
	GetPostByID(postID string) (*Post, error)
	GetPostsWithPagination(limit int, cursor *int64) ([]Post, *int64, error)
	UpdateCommentsDisabled(postID string, commentsDisabled bool) error
}

type DBPostRepository struct{}

func (r *DBPostRepository) Save(post Post, userID string) int64 {
	return post.Save(userID)
}

func (r *DBPostRepository) GetPostByID(postID string) (*Post, error) {
	return GetPostByID(postID)
}

func (r *DBPostRepository) GetPostsWithPagination(limit int, cursor *int64) ([]Post, *int64, error) {
	return GetPostsWithPagination(limit, cursor)
}

func (r *DBPostRepository) UpdateCommentsDisabled(postID string, commentsDisabled bool) error {
	return UpdateCommentsDisabled(postID, commentsDisabled)
}
