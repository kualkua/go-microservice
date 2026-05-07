package contracts

type PaymentRequest struct {
	UserId         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

type PaymentProcessedEvent struct {
	UserId         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	IdempotencyKey string  `json:"idempotency_key"`
	Status         string  `json:"status"`
	Message        string  `json:"message"`
}
