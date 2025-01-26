package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Create(context.Context, *sql.Tx, *User) error
	CreateAndInvite(context.Context, *User, time.Duration, string) error
	GetByUserID(context.Context, int64) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	Activate(context.Context, string) error
}

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	IsActive  bool     `json:"is_active"`
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

func (s *PostgresUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	INSERT INTO users (username, email, password) 
	VALUES ($1,$2,$3) RETURNING id, created_at 
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err := tx.QueryRowContext(ctx, query,
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

func (s *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
	SELECT id,username, email FROM users
	WHERE email = ($1) AND is_active = true ; 
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	user := &User{}
	row := s.db.QueryRowContext(ctx, query, email)
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresUserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `
	INSERT INTO user_invitation (token, user_id, expiry)
	VALUES ($1,$2,$3) 
	`
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	return err
}

func (s *PostgresUserStore) CreateAndInvite(ctx context.Context, user *User, exp time.Duration, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		if err := s.createUserInvitation(ctx, tx, token, exp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *PostgresUserStore) getUserFromToken(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
	select u.id, u.username, u.is_active from users u
	join user_invitation i on u.id = i.user_id 
	where i.token = ($1) AND i.expiry > ($2)
	`
	hash := sha256.Sum256([]byte(token))
	hashtoken := hex.EncodeToString(hash[:])
	user := User{}
	row := tx.QueryRowContext(ctx, query, hashtoken, time.Now())
	if err := row.Scan(&user.ID, &user.Username, &user.IsActive); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	update users
	set is_active = ($1)
	where id = ($2);
	`
	_, err := tx.ExecContext(ctx, query, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
	delete from user_invitation where user_id = $1;
	`
	if _, err := tx.ExecContext(ctx, query, userID); err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		user, err := s.getUserFromToken(ctx, tx, token)
		if err != nil {
			return err
		}
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}
