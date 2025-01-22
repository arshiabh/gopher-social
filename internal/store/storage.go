package store

import (
	"context"
	"database/sql"
	"errors"
)

type Storage struct {
	Users     UserStore
	Posts     PostStore
	Comments  CommentStore
	Followers FollowerStore
}

var (
	ErrNotFound = errors.New("data not founddd")
)

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users:     NewPostgresUserStore(db),
		Posts:     NewPostgresPostStore(db),
		Comments:  NewPostgresCommentStore(db),
		Followers: NewPostgresFollowerStore(db),
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
	}

	return tx.Commit()
}
