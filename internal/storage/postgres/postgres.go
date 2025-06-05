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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: db.Ping error: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	const op = "storage.postgres.CreatePost"

	if post.ID == "" {
		post.ID = utils.GenerateID()
	}
	post.CreatedAt = time.Now()

	query := `
		INSERT INTO posts (id, title, content, author, comments_enabled, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		post.ID, post.Title, post.Content, post.Author, post.CommentsEnabled, post.CreatedAt)
	if err != nil {
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (s *Storage) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {
	const op = "storage.postgres.GetPosts"

	query := `
		SELECT * FROM posts 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var posts []models.Post
	err := s.db.SelectContext(ctx, &posts, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

func (s *Storage) GetPost(ctx context.Context, id string) (models.Post, error) {
	const op = "storage.postgres.GetPost"

	query := `SELECT * FROM posts WHERE id = $1`

	var post models.Post
	err := s.db.GetContext(ctx, &post, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, errors.ErrNotFound
		}
		return models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (s *Storage) UpdatePost(ctx context.Context, post models.Post) error {
	const op = "storage.postgres.UpdatePost"

	query := `
		UPDATE posts 
		SET title = $1, content = $2, comments_enabled = $3
		WHERE id = $4
	`

	result, err := s.db.ExecContext(ctx, query,
		post.Title, post.Content, post.CommentsEnabled, post.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (s *Storage) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	const op = "storage.postgres.CreateComment"

	if _, err := s.GetPost(ctx, comment.PostID); err != nil {
		return models.Comment{}, err
	}

	if comment.ParentID != nil {
		var parent models.Comment
		err := s.db.GetContext(ctx, &parent,
			"SELECT * FROM comments WHERE id = $1", *comment.ParentID)
		if err != nil || parent.PostID != comment.PostID {
			return models.Comment{}, errors.ErrParentNotFound
		}
	}

	if comment.ID == "" {
		comment.ID = utils.GenerateID()
	}
	comment.CreatedAt = time.Now()

	query := `
		INSERT INTO comments (id, post_id, parent_id, author, content, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		comment.ID, comment.PostID, comment.ParentID, comment.Author, comment.Content, comment.CreatedAt)
	if err != nil {
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (s *Storage) GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]models.Comment, error) {
	const op = "storage.postgres.GetCommentsByPost"

	query := `
		SELECT * FROM comments 
		WHERE post_id = $1 AND parent_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var comments []models.Comment
	err := s.db.SelectContext(ctx, &comments, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (s *Storage) CountCommentsByPost(ctx context.Context, postID string) (int, error) {
	const op = "storage.postgres.CountCommentsByPost"

	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND parent_id IS NULL`

	var count int
	err := s.db.GetContext(ctx, &count, query, postID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (s *Storage) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	const op = "storage.postgres.GetCommentReplies"

	query := `
		SELECT * FROM comments 
		WHERE parent_id = $1 
		ORDER BY created_at ASC
	`

	var replies []models.Comment
	err := s.db.SelectContext(ctx, &replies, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return replies, nil
}

func (s *Storage) GetComment(ctx context.Context, id string) (models.Comment, error) {
	const op = "storage.postgres.GetComment"

	query := `SELECT * FROM comments WHERE id = $1`

	var comment models.Comment
	err := s.db.GetContext(ctx, &comment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Comment{}, errors.ErrNotFound
		}
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
