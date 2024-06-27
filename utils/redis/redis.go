package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisClient struct{}

var rdb *redis.Client

func InitRedis() *RedisClient {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "redis",
			DB:       0,
			PoolSize: 10,
		})
	}
	return &RedisClient{}
}
