package main

import (
	"context"
	"errors"
	"log"
	"payment-microservices/src/config"
	"payment-microservices/src/internal/adapters/logger"
	"payment-microservices/src/internal/adapters/rabbitmq"
	consumernotification "payment-microservices/src/internal/api/consumer/notification"
	appnotification "payment-microservices/src/internal/app/notification"
)

func main() {
	cfg := config.Load()

	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	notifier := logger.NewPaymentProcessedNotifier()
	notify := appnotification.NewNotifyPaymentProcessedUseCase(notifier)
	handler := consumernotification.NewHandler(cfg, conn, notify)

	if err := handler.Start(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Failed to start notification consumer: %v", err)
	}
}
