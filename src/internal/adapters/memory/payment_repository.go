package memory

import (
	"context"
	"log"
	"payment-microservices/src/internal/domain/payment"
)

type PaymentRepository struct{}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (r *PaymentRepository) Save(ctx context.Context, payment payment.Payment) error {
	log.Printf("Payment saved: %+v", payment)
	return nil
}
