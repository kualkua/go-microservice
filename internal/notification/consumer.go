package notification

import (
	"encoding/json"
	"log"
	"payment-microservices/internal/contracts"
	"payment-microservices/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	conn *amqp.Connection
}

func NewHandler(conn *amqp.Connection) *Handler {
	return &Handler{conn: conn}
}

func (h *Handler) Start() error {
	ch, err := h.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = rabbitmq.DeclarePaymentsProcessedQueue(ch)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		rabbitmq.PaymentProcessedQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Notification service is listening to queue payment_processed")

	for msg := range msgs {
		var res contracts.PaymentProcessedEvent
		if err := json.Unmarshal(msg.Body, &res); err != nil {
			log.Printf("Failed to decode payment processed message: %v", err)
			_ = msg.Nack(false, false)
			continue
		}
		log.Printf("Received payment processed message: %+v", res)
		if err := msg.Ack(false); err != nil {
			log.Printf("Failed to ack notification message: %v", err)
		}
	}

	return nil
}
