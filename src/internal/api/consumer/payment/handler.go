package payment

import (
	"context"
	"encoding/json"
	"log"
	"payment-microservices/src/config"
	apppayment "payment-microservices/src/internal/app/payment"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	cfg            *config.Config
	conn           *amqp.Connection
	processPayment *apppayment.ProcessPaymentUseCase
}

func NewHandler(
	cfg *config.Config,
	conn *amqp.Connection,
	processPayment *apppayment.ProcessPaymentUseCase,
) *Handler {
	return &Handler{
		cfg:            cfg,
		conn:           conn,
		processPayment: processPayment,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	ch, err := h.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.Qos(h.cfg.MessagesLimit, 0, false); err != nil {
		return err
	}

	if _, err := declareQueue(ch, h.cfg.PaymentsQueue); err != nil {
		return err
	}
	if _, err := declareQueue(ch, h.cfg.PaymentProcessedQueue); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		h.cfg.PaymentsQueue,
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

	log.Printf("Payments service is listening to queue %s", h.cfg.PaymentsQueue)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			h.handleMessage(ctx, msg)
		}
	}
}

func (h *Handler) handleMessage(ctx context.Context, msg amqp.Delivery) {
	var req paymentMessage
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		log.Printf("Failed to decode payment message: %v", err)
		_ = msg.Nack(false, false)
		return
	}

	payment := req.toDomain()
	log.Printf("Payment received: %+v", payment)

	if err := h.processPayment.Execute(ctx, payment); err != nil {
		log.Printf("Failed to process payment: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	if err := msg.Ack(false); err != nil {
		log.Printf("Failed to ack payment message: %v", err)
	}
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
