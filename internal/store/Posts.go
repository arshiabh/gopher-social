package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostStore interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
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
	var post Post
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID, &post.Title, &post.UserID, &post.Content, &post.CreatedAt, pq.Array(&post.Tags), &post.UpdatedAt); err != nil {
		return nil, err
	}
	return &post, nil
}
