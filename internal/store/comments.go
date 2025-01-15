package store

import (
	"context"
	"database/sql"
)

type CommentStore interface {
	GetByPostID(context.Context, int64) ([]Comments, error)
}

type Comments struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type PostgresCommentStore struct {
	db *sql.DB
}

func NewPostgresCommentStore(db *sql.DB) *PostgresCommentStore {
	return &PostgresCommentStore{
		db: db,
	}
}

func (s *PostgresCommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comments, error) {
	query := `
	SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, users.id, users.username 
	FROM comments c LEFT JOIN users ON c.user_id = users.id
	WHERE post_id = ($1)
	`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	var comments []Comments
	for rows.Next() {
		var c Comments
		c.User = User{}
		err := rows.Scan(&c.ID, &c.UserID, &c.PostID,&c.Content, &c.CreatedAt,
			&c.User.ID, &c.User.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
