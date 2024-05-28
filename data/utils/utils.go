package utils

import (
	"erp-user-service/factory"
	"os"
	"time"
)

type AppUtilities struct {
	Jwt        *JwtUtilities
	IpLocation *IpLocationUtils
	Redis      *RedisClient
	DeviceInfo *DeviceInfoUtilities
}

var ErrorFactory = &factory.ErrorFactory{}

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
