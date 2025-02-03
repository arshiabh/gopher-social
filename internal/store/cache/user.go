package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore interface {
	Get(context.Context, int64) (*store.User, error)
	Set(context.Context, *store.User) error
}

type RedisUserStore struct {
	db *redis.Client
}

func NewRedisUserStore(db *redis.Client) *RedisUserStore {
	return &RedisUserStore{
		db: db,
	}
}

const UserTimeExp = time.Minute

func (s *RedisUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cachekey := fmt.Sprintf("user-%v", userID)
	data, err := s.db.Get(ctx, cachekey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *RedisUserStore) Set(ctx context.Context, user *store.User) error {
	cachekey := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	if err := s.db.SetEx(ctx, cachekey, json, UserTimeExp).Err(); err != nil {
		return err
	}
	return nil
}
