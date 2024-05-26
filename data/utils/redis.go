package utils

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct{}

var rdb *redis.Client

func InitRedis() *RedisClient {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "redis",
		DB:       0,
		PoolSize: 10,
	})
	log.Println("Connected to redis")
	return &RedisClient{}
}

func (*RedisClient) StoreSessionInfo(userID string, refreshToken string, deviceInfo *DeviceInfo) error {
	ctx := context.Background()

	deviceInfoJson, _ := json.Marshal(deviceInfo)

	cmd := rdb.HSet(ctx, userID, refreshToken, deviceInfoJson)

	return cmd.Err()
}

func (*RedisClient) GetSessionInfo(userId string, refreshToken string) (*DeviceInfo, error) {
	ctx := context.Background()

	getCmd := rdb.HGet(ctx, userId, refreshToken)

	if getCmd.Err() != nil {
		return nil, getCmd.Err()
	}

	var deviceInfo DeviceInfo

	err := json.Unmarshal([]byte(getCmd.Val()), &deviceInfo)

	if err != nil {
		return nil, err
	}

	return &deviceInfo, nil
}

func (*RedisClient) RemoveRefreshToken(userID string) error {
	var ctx = context.Background()
	err := rdb.Del(ctx, userID).Err()

	if err != nil {
		return err
	}
	return nil
}

func (*RedisClient) CheckRefreshTokenValid(userID string, token string) bool {
	var ctx = context.Background()
	cmd := rdb.Get(ctx, userID)

	if cmd.Err() != nil {
		return false
	}

	return cmd.Val() == token
}
