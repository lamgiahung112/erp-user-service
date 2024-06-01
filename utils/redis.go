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
	if rdb != nil {
		return &RedisClient{}
	}
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

	if cmd.Err() != nil {
		return ErrorFactory.StoreSessionFailed()
	}
	return nil
}

func (*RedisClient) GetSessionInfo(userId string, refreshToken string) (*DeviceInfo, error) {
	ctx := context.Background()

	getCmd := rdb.HGet(ctx, userId, refreshToken)

	if getCmd.Err() != nil {
		return nil, ErrorFactory.NotFound("session data")
	}

	var deviceInfo DeviceInfo

	err := json.Unmarshal([]byte(getCmd.Val()), &deviceInfo)

	if err != nil {
		return nil, ErrorFactory.Malformatted("session data")
	}

	return &deviceInfo, nil
}

func (*RedisClient) RemoveRefreshToken(userID string, refreshToken string) error {
	var ctx = context.Background()
	err := rdb.HDel(ctx, userID, refreshToken).Err()

	if err != nil {
		return ErrorFactory.Unexpected()
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
