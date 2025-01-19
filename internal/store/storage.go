package store

import (
	"database/sql"
	"errors"
)

type Storage struct {
	Users    UserStore
	Posts    PostStore
	Comments CommentStore
}

var (
	ErrNotFound = errors.New("data not founddd")
)

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users:    NewPostgresUserStore(db),
		Posts:    NewPostgresPostStore(db),
		Comments: NewPostgresCommentStore(db),
	}
}
