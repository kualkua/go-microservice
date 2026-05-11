package rabbitmq

import (
	"payment-microservices/src/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConnection(cfg *config.Config) (*amqp.Connection, error) {
	return amqp.Dial(cfg.RabbitMQURL)
}

func declareQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}
