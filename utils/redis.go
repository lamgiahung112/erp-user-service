package utils

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
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

func (r *RedisClient) StoreSessionInfo(userID string, refreshToken string, deviceInfo *DeviceInfo) error {
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

func (*RedisClient) GetAllSessionsOfUser(userId string) ([]*DeviceInfo, error) {
	ctx := context.Background()

	getCmd := rdb.HGetAll(ctx, userId)

	if getCmd.Err() != nil {
		return nil, ErrorFactory.NotFound("session data")
	}

	var deviceInfos []*DeviceInfo
	for _, value := range getCmd.Val() {
		var deviceInfo DeviceInfo
		err := json.Unmarshal([]byte(value), &deviceInfo)
		if err != nil {
			return nil, ErrorFactory.Malformatted("session data")
		}
		deviceInfos = append(deviceInfos, &deviceInfo)
	}
	return deviceInfos, nil
}

func (*RedisClient) RevokeAllUserSessions(userId string) error {
	ctx := context.Background()
	cmd := rdb.Del(ctx, userId)

	if cmd.Err() != nil {
		return ErrorFactory.Unexpected()
	}
	return nil
}

func (*RedisClient) RemoveSessionInfo(userID string, refreshToken string) error {
	var ctx = context.Background()
	err := rdb.HDel(ctx, userID, refreshToken).Err()

	if err != nil {
		return ErrorFactory.Unexpected()
	}
	return nil
}

func loginOtpKey(userId string) string {
	return "login_otp:" + userId
}

func (*RedisClient) StoreUserLoginOtp(userId string) (string, error) {
	ctx := context.Background()

	random, _ := uuid.NewV4()
	otp := strings.ToUpper(random.String())[0:6]
	cmd := rdb.Set(ctx, loginOtpKey(userId), otp, 10*time.Minute)

	if cmd.Err() != nil {
		return "", ErrorFactory.StoreSessionFailed()
	}
	return otp, nil
}

func (*RedisClient) GetUserLoginOtp(userId string) (string, error) {
	ctx := context.Background()

	cmd := rdb.Get(ctx, loginOtpKey(userId))
	if cmd.Err() != nil {
		return "", ErrorFactory.TimeOut("Otp")
	}
	return cmd.Val(), nil
}

func (*RedisClient) RemoveUserLoginOtp(userId string) error {
	ctx := context.Background()

	cmd := rdb.Del(ctx, loginOtpKey(userId))
	if cmd.Err() != nil {
		return ErrorFactory.Unexpected()
	}
	return nil
}
