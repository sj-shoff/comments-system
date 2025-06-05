package postgres

import (
	"comments-system/internal/config"
	"comments-system/internal/models"
	"comments-system/pkg/errors"
	"comments-system/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg config.Postgres) (*Storage, error) {

	const op = "storage.postgres.NewPostgresDB"

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: db.Ping error: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {

	const op = "storage.postgres.CreatePost"

	const query = `
		INSERT INTO posts (id, title, content, author, comments_enabled, created_at)
		VALUES (:id, :title, :content, :author, :comments_enabled, :created_at)
		RETURNING id
	`

	if post.ID == "" {
		post.ID = utils.GenerateID()
	}
	post.CreatedAt = time.Now()

	_, err := s.db.NamedExecContext(ctx, query, post)
	if err != nil {
		return models.Post{}, fmt.Errorf("%s: failed to create post: %w", op, err)
	}

	return post, nil
}

func (s *Storage) GetPosts(ctx context.Context) ([]models.Post, error) {

	const op = "storage.postgres.GetPosts"

	const query = `SELECT * FROM posts ORDER BY created_at DESC`

	var posts []models.Post
	err := s.db.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get posts: %w", op, err)
	}

	return posts, nil
}

func (s *Storage) GetPost(ctx context.Context, id string) (models.Post, error) {

	const op = "storage.postgres.GetPost"

	const query = `SELECT * FROM posts WHERE id = $1`

	var post models.Post
	err := s.db.GetContext(ctx, &post, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, errors.ErrNotFound
		}
		return models.Post{}, fmt.Errorf("%s: failed to get post: %w", op, err)
	}

	return post, nil
}

func (s *Storage) UpdatePost(ctx context.Context, post models.Post) error {

	const op = "storage.postgres.UpdatePost"

	const query = `
		UPDATE posts 
		SET title = :title, content = :content, comments_enabled = :comments_enabled
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, query, post)
	if err != nil {
		return fmt.Errorf("%s: failed to update post: %w", op, err)
	}

	return nil
}

func (s *Storage) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {

	const op = "storage.postgres.CreateComment"

	const query = `
		INSERT INTO comments (id, post_id, parent_id, author, content, created_at)
		VALUES (:id, :post_id, :parent_id, :author, :content, :created_at)
		RETURNING id
	`

	if comment.ParentID != nil {
		var parent models.Comment
		err := s.db.GetContext(ctx, &parent,
			"SELECT * FROM comments WHERE id = $1", *comment.ParentID)

		if err != nil || parent.PostID != comment.PostID {
			return models.Comment{}, errors.ErrParentNotFound
		}
	}

	comment.ID = utils.GenerateID()
	comment.CreatedAt = time.Now()

	_, err := s.db.NamedExecContext(ctx, query, comment)
	if err != nil {
		return models.Comment{}, fmt.Errorf("%s: failed to create comment: %w", op, err)
	}

	return comment, nil
}

func (s *Storage) GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]models.Comment, error) {

	const op = "storage.postgres.GetCommentsByPost"

	const query = `
		SELECT * FROM comments 
		WHERE post_id = $1 AND parent_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var comments []models.Comment
	err := s.db.SelectContext(ctx, &comments, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get comments by post: %w", op, err)
	}

	return comments, nil
}

func (s *Storage) CountCommentsByPost(ctx context.Context, postID string) (int, error) {

	const op = "storage.postgres.CountCommentsByPost"

	const query = `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND parent_id IS NULL`

	var count int
	err := s.db.GetContext(ctx, &count, query, postID)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to count comments: %w", op, err)
	}

	return count, nil
}

func (s *Storage) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {

	const op = "storage.postgres.GetCommentReplies"

	const query = `
		SELECT * FROM comments 
		WHERE parent_id = $1 
		ORDER BY created_at ASC
	`

	var replies []models.Comment
	err := s.db.SelectContext(ctx, &replies, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get comment replies: %w", op, err)
	}

	return replies, nil
}

func (s *Storage) GetComment(ctx context.Context, id string) (models.Comment, error) {
	const op = "storage.postgres.GetComment"

	const query = `SELECT * FROM comments WHERE id = $1`

	var comment models.Comment
	err := s.db.GetContext(ctx, &comment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Comment{}, errors.ErrNotFound
		}
		return models.Comment{}, fmt.Errorf("%s: failed to get comment: %w", op, err)
	}

	return comment, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
