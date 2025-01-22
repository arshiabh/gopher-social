package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Create(context.Context, *User) error
	CreateAndInvite(context.Context, *User, string) error
	GetByUserID(context.Context, int64) (*User, error)
}

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password.hash,
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

func (s *PostgresUserStore) CreateAndInvite(ctx context.Context, user *User, token string) error {
	return nil
}
