package payment

import (
	"errors"
	domainpayment "payment-microservices/src/internal/domain/payment"
)

type createPaymentRequest struct {
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type paymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func validateCreatePaymentRequest(req createPaymentRequest) error {
	if req.UserID == "" {
		return errors.New("user_id is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if req.Currency == "" {
		return errors.New("currency is required")
	}
	if req.IdempotencyKey == "" {
		return errors.New("idempotency_key is required")
	}
	return nil
}

func (req createPaymentRequest) toDomain() domainpayment.Payment {
	return domainpayment.Payment{
		UserID:         req.UserID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
	}
}
