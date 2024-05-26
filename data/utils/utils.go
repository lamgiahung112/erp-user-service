package utils

import (
	"os"
	"time"
)

type AppUtilities struct {
	Jwt        *JwtUtilities
	IpLocation *IpLocationUtils
	Redis      *RedisClient
	DeviceInfo *DeviceInfoUtilities
}

func New() *AppUtilities {
	InitIpLocationUtils()
	InitRedis()

	return &AppUtilities{
		Jwt: &JwtUtilities{
			expirationPeriod: 7 * 24 * time.Hour,
			key:              []byte(os.Getenv("JWT_KEY")),
		},
		IpLocation: &IpLocationUtils{},
		Redis:      &RedisClient{},
		DeviceInfo: &DeviceInfoUtilities{},
	}
}
