package payment

import apppayment "payment-microservices/src/internal/app/payment"

type Handler struct {
	createPayment *apppayment.CreatePaymentUseCase
}

func NewHandler(createPayment *apppayment.CreatePaymentUseCase) *Handler {
	return &Handler{createPayment: createPayment}
}
