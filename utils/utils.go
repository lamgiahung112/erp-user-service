package utils

import (
	"erp-user-service/factory"
	"erp-user-service/utils/rabbitmq"
	"os"
	"time"
)

type AppUtilities struct {
	Jwt          *JwtUtilities
	IpLocation   *IpLocationUtils
	Redis        *RedisClient
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
		Redis:        InitRedis(),
		DeviceInfo:   &DeviceInfoUtilities{},
		QR:           &QRUtils{},
		EventEmitter: rabbitmq.GetEventEmitter(),
	}
}
