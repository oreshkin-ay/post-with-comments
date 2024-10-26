// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID                string             `json:"id"`
	PostID            string             `json:"postId"`
	Text              string             `json:"text"`
	ParentCommentID   *string            `json:"parentCommentId,omitempty"`
	CreatedAt         string             `json:"createdAt"`
	Children          *CommentConnection `json:"children"`
	ChildCommentCount int                `json:"childCommentCount"`
}

type CommentConnection struct {
	Edges    []*CommentEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

type CommentEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Comment `json:"node"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Mutation struct {
}

type NewCommentInput struct {
	PostID          string  `json:"postId"`
	Text            string  `json:"text"`
	ParentCommentID *string `json:"parentCommentId,omitempty"`
}

type NewPostInput struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	CommentsDisabled bool   `json:"commentsDisabled"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type Post struct {
	ID               string             `json:"id"`
	Title            string             `json:"title"`
	Content          string             `json:"content"`
	User             *User              `json:"user"`
	CommentsDisabled bool               `json:"commentsDisabled"`
	Comments         *CommentConnection `json:"comments"`
}

type PostConnection struct {
	Edges    []*PostEdge `json:"edges"`
	PageInfo *PageInfo   `json:"pageInfo"`
}

type PostEdge struct {
	Cursor string `json:"cursor"`
	Node   *Post  `json:"node"`
}

type Query struct {
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type Subscription struct {
}

type UpdateCommentsDisabledInput struct {
	PostID           string `json:"postId"`
	CommentsDisabled bool   `json:"commentsDisabled"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
