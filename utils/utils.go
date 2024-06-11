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

	return &AppUtilities{
		Jwt: &JwtUtilities{
			expirationPeriod: 5*time.Minute + 30*time.Second,
			key:              []byte(os.Getenv("JWT_KEY")),
		},
		IpLocation: &IpLocationUtils{},
		Redis:      InitRedis(),
		DeviceInfo: &DeviceInfoUtilities{},
	}
}
