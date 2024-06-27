package redis

import (
	"context"
	"encoding/json"
	"erp-user-service/utils"
)

func (r *RedisClient) StoreSessionInfo(userID string, refreshToken string, deviceInfo *utils.DeviceInfo) error {
	ctx := context.Background()

	deviceInfoJson, _ := json.Marshal(deviceInfo)

	cmd := rdb.HSet(ctx, userID, refreshToken, deviceInfoJson)

	if cmd.Err() != nil {
		return utils.ErrorFactory.StoreSessionFailed()
	}
	return nil
}

func (*RedisClient) GetSessionInfo(userId string, refreshToken string) (*utils.DeviceInfo, error) {
	ctx := context.Background()

	getCmd := rdb.HGet(ctx, userId, refreshToken)

	if getCmd.Err() != nil {
		return nil, utils.ErrorFactory.NotFound("session data")
	}

	var deviceInfo utils.DeviceInfo

	err := json.Unmarshal([]byte(getCmd.Val()), &deviceInfo)

	if err != nil {
		return nil, utils.ErrorFactory.Malformatted("session data")
	}

	return &deviceInfo, nil
}

func (*RedisClient) GetAllSessionsOfUser(userId string) ([]*utils.DeviceInfo, error) {
	ctx := context.Background()

	getCmd := rdb.HGetAll(ctx, userId)

	if getCmd.Err() != nil {
		return nil, utils.ErrorFactory.NotFound("session data")
	}

	var deviceInfos []*utils.DeviceInfo
	for _, value := range getCmd.Val() {
		var deviceInfo utils.DeviceInfo
		err := json.Unmarshal([]byte(value), &deviceInfo)
		if err != nil {
			return nil, utils.ErrorFactory.Malformatted("session data")
		}
		deviceInfos = append(deviceInfos, &deviceInfo)
	}
	return deviceInfos, nil
}

func (*RedisClient) RevokeAllUserSessions(userId string) error {
	ctx := context.Background()
	cmd := rdb.Del(ctx, userId)

	if cmd.Err() != nil {
		return utils.ErrorFactory.Unexpected()
	}
	return nil
}

func (*RedisClient) RemoveSessionInfo(userID string, refreshToken string) error {
	var ctx = context.Background()
	err := rdb.HDel(ctx, userID, refreshToken).Err()

	if err != nil {
		return utils.ErrorFactory.Unexpected()
	}
	return nil
}
