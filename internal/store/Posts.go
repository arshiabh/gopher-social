package store

import "database/sql"

type PostStore interface {
	Get() error
}

type PostgresPostStore struct {
	db *sql.DB
}

func NewPostgresPostStore(db *sql.DB) *PostgresPostStore {
	return &PostgresPostStore{
		db: db,
	}
}

func (s *PostgresPostStore) Get() error {
	return nil
}
