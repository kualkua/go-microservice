package payment

import "context"

type Publisher interface {
	PublishPayment(ctx context.Context, payment Payment) error
}

type ProcessedEventPublisher interface {
	PublishProcessed(ctx context.Context, event ProcessedEvent) error
}
