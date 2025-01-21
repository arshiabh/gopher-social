package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type PostStore interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetaData, error)
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
	User      User       `json:"user"`
}

type PostWithMetaData struct {
	Post          Post
	CommentsCount int64 `json:"comments_count"`
}

type PostgresPostStore struct {
	db *sql.DB
}

func NewPostgresPostStore(db *sql.DB) *PostgresPostStore {
	return &PostgresPostStore{
		db: db,
	}
}

func (s *PostgresPostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
	SELECT p.id, p.user_id, u.username, p.title, p.content, p.tags, p.created_at,
	COUNT(c.id) AS comments_count FROM posts p
	JOIN users u ON u.id = p.user_id
	JOIN followers f ON f.follower_id = u.id OR p.user_id = ($1)
	LEFT JOIN comments c ON c.post_id = p.id GROUP BY (p.id,u.username) 
	ORDER BY p.created_at ` + fq.Order + `
	LIMIT ($2) OFFSET ($3);	
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var FeedList []PostWithMetaData
	for rows.Next() {
		var feed PostWithMetaData
		feed.Post = Post{}
		if err := rows.Scan(&feed.Post.ID,
			&feed.Post.UserID,
			&feed.Post.User.Username,
			&feed.Post.Title,
			&feed.Post.Content,
			pq.Array(&feed.Post.Tags), &feed.Post.CreatedAt, &feed.CommentsCount); err != nil {
			return nil, err
		}
		FeedList = append(FeedList, feed)
	}
	return FeedList, nil
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
			return nil, ErrNotFound
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
