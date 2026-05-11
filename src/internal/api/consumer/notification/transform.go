package notification

import "payment-microservices/src/internal/domain/payment"

type processedEventMessage struct {
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
	Status         string  `json:"status"`
	Message        string  `json:"message"`
}

func (msg processedEventMessage) toDomain() payment.ProcessedEvent {
	return payment.ProcessedEvent{
		UserID:         msg.UserID,
		Amount:         msg.Amount,
		Currency:       msg.Currency,
		IdempotencyKey: msg.IdempotencyKey,
		Status:         msg.Status,
		Message:        msg.Message,
	}
}
