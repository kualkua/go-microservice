package notification

import (
	"context"
	"payment-microservices/src/internal/domain/payment"
)

type PaymentProcessedNotifier interface {
	NotifyPaymentProcessed(ctx context.Context, event payment.ProcessedEvent) error
}
