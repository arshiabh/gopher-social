package cache

import (
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	User UserStore
}

func NewRedisStorage(db *redis.Client) Storage {
	return Storage{
		User: NewRedisUserStore(db),
	}
}
