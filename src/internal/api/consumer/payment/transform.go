package payment

import domainpayment "payment-microservices/src/internal/domain/payment"

type paymentMessage struct {
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func (msg paymentMessage) toDomain() domainpayment.Payment {
	return domainpayment.Payment{
		UserID:         msg.UserID,
		Amount:         msg.Amount,
		Currency:       msg.Currency,
		IdempotencyKey: msg.IdempotencyKey,
	}
}
