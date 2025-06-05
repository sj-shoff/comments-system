package models

import (
	"time"
)

// Post представляет пост
type Post struct {
	ID              string    `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Content         string    `json:"content" db:"content"`
	Author          string    `json:"author" db:"author"`
	CommentsEnabled bool      `json:"commentsEnabled" db:"comments_enabled"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

// Comment представляет комментарий
type Comment struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"postId" db:"post_id"`
	ParentID  *string   `json:"parentId,omitempty" db:"parent_id"`
	Author    string    `json:"author" db:"author"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// CommentsPage представляет страницу с комментариями
type CommentsPage struct {
	Total    int       `json:"total"`
	Comments []Comment `json:"comments"`
}

// CreatePostInput представляет входные данные для создания поста
type CreatePostInput struct {
	Title           string `json:"title"`
	Content         string `json:"content"`
	Author          string `json:"author"`
	CommentsEnabled bool   `json:"commentsEnabled"`
}

// CreateCommentInput представляет входные данные для создания комментария
type CreateCommentInput struct {
	PostID   string  `json:"postId"`
	ParentID *string `json:"parentId,omitempty"`
	Author   string  `json:"author"`
	Content  string  `json:"content"`
}
