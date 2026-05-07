package payment

import "payment-microservices/internal/contracts"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ProcessPayment(req contracts.PaymentRequest) (contracts.PaymentProcessedEvent, error) {
	if err := s.repo.Save(req); err != nil {
		return contracts.PaymentProcessedEvent{}, err
	}

	return contracts.PaymentProcessedEvent{
		UserId:         req.UserId,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
		Status:         "processed",
		Message:        "Payment processed successfully",
	}, nil
}
