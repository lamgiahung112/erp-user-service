package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func ConnectRabbitMQ() {
	if conn != nil {
		return
	}
	c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
	if err != nil {
		panic(err)
	}
	conn = c
}
