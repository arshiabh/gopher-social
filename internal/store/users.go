package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserStore interface {
	Create(context.Context, *User) error
	GetByUserID(context.Context, int64) (*User, error)
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

func (s *PostgresUserStore) GetByUserID(ctx context.Context, userid int64) (*User, error) {
	query := `
	SELECT id, username, email, created_at FROM users
	WHERE id = ($1);
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	var user User
	res := s.db.QueryRowContext(ctx, query, userid)
	if err := res.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}


