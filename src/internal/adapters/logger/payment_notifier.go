package logger

import (
	"context"
	"log"
	"payment-microservices/src/internal/domain/payment"
)

type PaymentProcessedNotifier struct{}

func NewPaymentProcessedNotifier() *PaymentProcessedNotifier {
	return &PaymentProcessedNotifier{}
}

func (n *PaymentProcessedNotifier) NotifyPaymentProcessed(ctx context.Context, event payment.ProcessedEvent) error {
	log.Printf("Received payment processed event: %+v", event)
	return nil
}
