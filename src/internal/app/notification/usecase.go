package notification

import (
	"context"
	domainnotification "payment-microservices/src/internal/domain/notification"
	"payment-microservices/src/internal/domain/payment"
)

type NotifyPaymentProcessedUseCase struct {
	notifier domainnotification.PaymentProcessedNotifier
}

func NewNotifyPaymentProcessedUseCase(notifier domainnotification.PaymentProcessedNotifier) *NotifyPaymentProcessedUseCase {
	return &NotifyPaymentProcessedUseCase{notifier: notifier}
}

func (u *NotifyPaymentProcessedUseCase) Execute(ctx context.Context, event payment.ProcessedEvent) error {
	return u.notifier.NotifyPaymentProcessed(ctx, event)
}
