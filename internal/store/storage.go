package store

import "database/sql"

type Storage struct {
	Users UserStore
	Posts PostStore
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users: NewPostgresUserStore(db),
		Posts: NewPostgresPostStore(db),
	}
}
