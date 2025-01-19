package store

import (
	"context"
	"database/sql"
	"time"
)

type FollowerStore interface {
	Follow(context.Context, int64, *User) error
	UnFollow(context.Context, int64, *User) error
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAT  string `json:"created_at"`
}

type PostgresFollowerStore struct {
	db *sql.DB
}

func NewPostgresFollowerStore(db *sql.DB) *PostgresFollowerStore {
	return &PostgresFollowerStore{
		db: db,
	}
}

func (s *PostgresFollowerStore) Follow(ctx context.Context, userID int64, follower *User) error {
	query := `
	INSERT INTO followers (follower_id, user_id)
	VALUES ($1,$2); 
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, follower.ID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresFollowerStore) UnFollow(ctx context.Context, userID int64, follower *User) error {
	query := `
	DELETE FROM followers 
	WHERE user_id = ($1) AND follower_id = ($2);
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, follower.ID)
	if err != nil {
		return err
	}
	return nil
}
