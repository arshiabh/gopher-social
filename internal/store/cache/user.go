package cache

import (
	"context"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore interface{
	Get(context.Context, *store.User)
}

type RedisUserStore struct {
	db *redis.Client
}

func NewRedisUserStore(db *redis.Client) *RedisUserStore{
	return &RedisUserStore{
		db: db,
	}
}

func (s *RedisUserStore) Get(ctx context.Context, user *store.User) {
	
}