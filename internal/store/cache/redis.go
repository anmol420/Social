package cache

import "github.com/redis/go-redis/v9"

func NewRedisClient(add string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: add,
		DB:   db,
	})
}
