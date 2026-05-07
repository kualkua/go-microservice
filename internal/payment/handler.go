package payment

import (
	"encoding/json"
	"log"
	"payment-microservices/internal/contracts"
	"payment-microservices/pkg/config"
	"payment-microservices/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	conn    *amqp.Connection
	service *Service
}

func NewHandler(conn *amqp.Connection) *Handler {
	repo := NewRepository()
	service := NewService(repo)

	return &Handler{
		conn:    conn,
		service: service,
	}
}

func (h *Handler) Start() error {
	ch, err := h.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	conf := config.Load()
	err = ch.Qos(conf.MessagesLimit, 0, false)
	if err != nil {
		return err
	}

	_, err = rabbitmq.DeclarePaymentsQueue(ch)
	if err != nil {
		return err
	}
	_, err = rabbitmq.DeclarePaymentsProcessedQueue(ch)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		rabbitmq.PaymentsQueue,
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

	log.Println("Payments service is listening to queue payments")

	for msg := range msgs {
		var req contracts.PaymentRequest
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			log.Printf("Failed to decode payment message: %v", err)
			_ = msg.Nack(false, false)
			continue
		}

		log.Printf("Payment received: %+v", req)

		event, err := h.service.ProcessPayment(req)
		if err != nil {
			log.Printf("Failed to process payment: %v", err)
			_ = msg.Nack(false, true)
			continue
		}

		jsonData, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal payment processed event: %v", err)
			_ = msg.Nack(false, true)
			continue
		}
		err = ch.Publish(
			"",
			rabbitmq.PaymentProcessedQueue,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonData,
			},
		)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			_ = msg.Nack(false, true)
			continue
		}

		if err := msg.Ack(false); err != nil {
			log.Printf("Failed to ack message: %v", err)
		}
	}

	return nil
}
