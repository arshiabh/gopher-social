package store

import (
	"context"
	"database/sql"
)

type UserStore interface {
	Create(context.Context, *User) error
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, email, password) 
	VALUES ($1,$2,$3) RETURNING id, created_at 
	`
	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
