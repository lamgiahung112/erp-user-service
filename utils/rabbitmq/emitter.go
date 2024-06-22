package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventEmitter struct {
	conn *amqp.Connection
}

type mailRequestPayload struct {
	MailType string `json:"name"`
	Data     any    `json:"data"`
}

type MailType string

func (e MailType) String() string {
	return string(e)
}

const (
	LoginOTP      = MailType("login_otp")
	VerifyAccount = MailType("verify_account")
)

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

func (e *EventEmitter) pushEmailRequest(mailType MailType, payload any) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	jsonPayload, err := json.Marshal(&mailRequestPayload{
		MailType: string(mailType),
		Data:     payload,
	})
	if err != nil {
		return err
	}

	err = ch.Publish(
		"mail_topic",
		"mail-queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(jsonPayload),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventEmitter) setup() {
	channel, err := e.conn.Channel()

	if err != nil {
		panic(err)
	}
	defer channel.Close()
	_ = declareExchange(channel)
}
