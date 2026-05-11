package rabbitmq

import (
	"context"
	"encoding/json"
	"payment-microservices/src/config"
	"payment-microservices/src/internal/domain/payment"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentPublisher struct {
	cfg  *config.Config
	conn *amqp.Connection
}

func NewPaymentPublisher(cfg *config.Config, conn *amqp.Connection) *PaymentPublisher {
	return &PaymentPublisher{
		cfg:  cfg,
		conn: conn,
	}
}

func (p *PaymentPublisher) PublishPayment(ctx context.Context, payment payment.Payment) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err := declareQueue(ch, p.cfg.PaymentsQueue); err != nil {
		return err
	}

	body, err := json.Marshal(paymentMessageFromDomain(payment))
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		"",
		p.cfg.PaymentsQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

type ProcessedEventPublisher struct {
	cfg  *config.Config
	conn *amqp.Connection
}

func NewProcessedEventPublisher(cfg *config.Config, conn *amqp.Connection) *ProcessedEventPublisher {
	return &ProcessedEventPublisher{
		cfg:  cfg,
		conn: conn,
	}
}

func (p *ProcessedEventPublisher) PublishProcessed(ctx context.Context, event payment.ProcessedEvent) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err := declareQueue(ch, p.cfg.PaymentProcessedQueue); err != nil {
		return err
	}

	body, err := json.Marshal(processedEventMessageFromDomain(event))
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		"",
		p.cfg.PaymentProcessedQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

type paymentMessage struct {
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func paymentMessageFromDomain(payment payment.Payment) paymentMessage {
	return paymentMessage{
		UserID:         payment.UserID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		IdempotencyKey: payment.IdempotencyKey,
	}
}

type processedEventMessage struct {
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
	Status         string  `json:"status"`
	Message        string  `json:"message"`
}

func processedEventMessageFromDomain(event payment.ProcessedEvent) processedEventMessage {
	return processedEventMessage{
		UserID:         event.UserID,
		Amount:         event.Amount,
		Currency:       event.Currency,
		IdempotencyKey: event.IdempotencyKey,
		Status:         event.Status,
		Message:        event.Message,
	}
}
