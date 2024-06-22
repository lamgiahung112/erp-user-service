package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"mail_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"mail-queue",
		true,
		false,
		true,
		false,
		nil,
	)
}
