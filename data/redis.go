package data

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct{}

var rdb *redis.Client

func InitRedis() *RedisClient {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "redis",
		DB:       0,
	})
	log.Println("Connected to redis")
	return &RedisClient{}
}

func (rd *RedisClient) AddRefreshToken(userID string, token string, expiresIn time.Duration) error {
	var ctx = context.Background()
	err := rdb.Set(ctx, userID, token, expiresIn).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rd *RedisClient) RemoveRefreshToken(userID string) error {
	var ctx = context.Background()
	err := rdb.Del(ctx, userID).Err()

	if err != nil {
		return err
	}
	return nil
}

func (rd *RedisClient) CheckRefreshTokenValid(userID string, token string) bool {
	var ctx = context.Background()
	cmd := rdb.Get(ctx, userID)

	if cmd.Err() != nil {
		return false
	}

	return cmd.Val() == token
}
