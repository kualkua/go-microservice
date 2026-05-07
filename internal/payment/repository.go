package payment

import (
	"log"
	"payment-microservices/internal/contracts"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Save(req contracts.PaymentRequest) error {
	log.Printf("Payment saved: %+v", req)
	return nil
}
