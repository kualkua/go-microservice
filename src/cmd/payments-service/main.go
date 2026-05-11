package main

import (
	"context"
	"errors"
	"log"
	"payment-microservices/src/config"
	"payment-microservices/src/internal/adapters/memory"
	"payment-microservices/src/internal/adapters/rabbitmq"
	consumerpayment "payment-microservices/src/internal/api/consumer/payment"
	apppayment "payment-microservices/src/internal/app/payment"
)

func main() {
	cfg := config.Load()

	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	repository := memory.NewPaymentRepository()
	processedEventPublisher := rabbitmq.NewProcessedEventPublisher(cfg, conn)
	processPayment := apppayment.NewProcessPaymentUseCase(repository, processedEventPublisher)
	handler := consumerpayment.NewHandler(cfg, conn, processPayment)

	if err := handler.Start(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Failed to start payments consumer: %v", err)
	}
}
