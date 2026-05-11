package payment

import (
	"context"
	domainpayment "payment-microservices/src/internal/domain/payment"
)

type CreatePaymentUseCase struct {
	publisher domainpayment.Publisher
}

func NewCreatePaymentUseCase(publisher domainpayment.Publisher) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{publisher: publisher}
}

func (u *CreatePaymentUseCase) Execute(ctx context.Context, payment domainpayment.Payment) error {
	return u.publisher.PublishPayment(ctx, payment)
}

type ProcessPaymentUseCase struct {
	repository              domainpayment.Repository
	processedEventPublisher domainpayment.ProcessedEventPublisher
}

func NewProcessPaymentUseCase(
	repository domainpayment.Repository,
	processedEventPublisher domainpayment.ProcessedEventPublisher,
) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		repository:              repository,
		processedEventPublisher: processedEventPublisher,
	}
}

func (u *ProcessPaymentUseCase) Execute(ctx context.Context, payment domainpayment.Payment) error {
	if err := u.repository.Save(ctx, payment); err != nil {
		return err
	}

	event := domainpayment.NewProcessedEvent(payment)
	return u.processedEventPublisher.PublishProcessed(ctx, event)
}
