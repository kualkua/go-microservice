package payment

const StatusProcessed = "processed"

type Payment struct {
	UserID         string
	Amount         float64
	Currency       string
	IdempotencyKey string
}

type ProcessedEvent struct {
	UserID         string
	Amount         float64
	Currency       string
	IdempotencyKey string
	Status         string
	Message        string
}

func NewProcessedEvent(payment Payment) ProcessedEvent {
	return ProcessedEvent{
		UserID:         payment.UserID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		IdempotencyKey: payment.IdempotencyKey,
		Status:         StatusProcessed,
		Message:        "Payment processed successfully",
	}
}
