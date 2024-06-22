package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventEmitter struct {
	conn *amqp.Connection
}

type MailRequestPayload struct {
	MailType string `json:"name"`
	Data     string `json:"data"`
}

var emitter *EventEmitter

func GetEventEmitter() *EventEmitter {
	if emitter == nil {
		emitter = &EventEmitter{
			conn: conn,
		}
		emitter.setup()
	}
	return emitter
}

func (e *EventEmitter) setup() {
	channel, err := e.conn.Channel()

	if err != nil {
		panic(err)
	}
	defer channel.Close()
	_ = declareExchange(channel)
}
