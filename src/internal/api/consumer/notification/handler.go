package notification

import (
	"context"
	"encoding/json"
	"log"
	"payment-microservices/src/config"
	appnotification "payment-microservices/src/internal/app/notification"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	cfg    *config.Config
	conn   *amqp.Connection
	notify *appnotification.NotifyPaymentProcessedUseCase
}

func NewHandler(
	cfg *config.Config,
	conn *amqp.Connection,
	notify *appnotification.NotifyPaymentProcessedUseCase,
) *Handler {
	return &Handler{
		cfg:    cfg,
		conn:   conn,
		notify: notify,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	ch, err := h.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err := declareQueue(ch, h.cfg.PaymentProcessedQueue); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		h.cfg.PaymentProcessedQueue,
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

	log.Printf("Notification service is listening to queue %s", h.cfg.PaymentProcessedQueue)

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
	var event processedEventMessage
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to decode payment processed message: %v", err)
		_ = msg.Nack(false, false)
		return
	}

	if err := h.notify.Execute(ctx, event.toDomain()); err != nil {
		log.Printf("Failed to notify payment processed event: %v", err)
		_ = msg.Nack(false, true)
		return
	}

	if err := msg.Ack(false); err != nil {
		log.Printf("Failed to ack notification message: %v", err)
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
