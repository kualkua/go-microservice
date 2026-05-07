package rabbitmq

import (
	"payment-microservices/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PaymentsQueue         = "payments"
	PaymentProcessedQueue = "payment_processed"
)

func NewConnection(cfg *config.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DeclarePaymentsQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		PaymentsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
}

func DeclarePaymentsProcessedQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		PaymentProcessedQueue,
		true,
		false,
		false,
		false,
		nil,
	)
}
