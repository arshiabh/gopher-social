package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type PostStore interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Patch(context.Context, *Post) error
}

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Comments  []Comments `json:"comments"`
}

type PostgresPostStore struct {
	db *sql.DB
}

func NewPostgresPostStore(db *sql.DB) *PostgresPostStore {
	return &PostgresPostStore{
		db: db,
	}
}

func (s *PostgresPostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags) 
	VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at 
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresPostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT * FROM posts WHERE id = ($1);
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var post Post
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID, &post.Title, &post.UserID, &post.Content, &post.CreatedAt, pq.Array(&post.Tags), &post.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("no data found")
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostgresPostStore) Delete(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts
	WHERE id = ($1);
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresPostStore) Patch(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts SET title = $1, content = $2, updated_at = now() 
	WHERE id = $3;
	`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post.Title, post.Content, post.ID)
	return err
}
