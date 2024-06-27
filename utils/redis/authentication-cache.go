package redis

import (
	"context"
	"erp-user-service/utils"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

func _loginOtpKey(userId string) string {
	return "login_otp:" + userId
}

func (*RedisClient) StoreUserLoginOtp(userId string) (string, error) {
	ctx := context.Background()

	random, _ := uuid.NewV4()
	otp := strings.ReplaceAll(random.String(), "-", "")
	cmd := rdb.Set(ctx, _loginOtpKey(userId), otp, 10*time.Minute)

	if cmd.Err() != nil {
		return "", utils.ErrorFactory.StoreSessionFailed()
	}
	return otp, nil
}
