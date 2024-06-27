package utils

import (
	"erp-user-service/factory"
	"erp-user-service/utils/rabbitmq"
	"erp-user-service/utils/redis"
	"os"
	"time"
)

type AppUtilities struct {
	Jwt          *JwtUtilities
	IpLocation   *IpLocationUtils
	Redis        *redis.RedisClient
	DeviceInfo   *DeviceInfoUtilities
	QR           *QRUtils
	EventEmitter *rabbitmq.EventEmitter
}

var ErrorFactory = &factory.ErrorFactory{}

func New() *AppUtilities {
	InitIpLocationUtils()
	rabbitmq.ConnectRabbitMQ()

	return &AppUtilities{
		Jwt: &JwtUtilities{
			expirationPeriod: 5*time.Minute + 30*time.Second,
			key:              []byte(os.Getenv("JWT_KEY")),
		},
		IpLocation:   &IpLocationUtils{},
		Redis:        redis.InitRedis(),
		DeviceInfo:   &DeviceInfoUtilities{},
		QR:           &QRUtils{},
		EventEmitter: rabbitmq.GetEventEmitter(),
	}
}
