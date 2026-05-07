package main

import (
	"log"
	"payment-microservices/internal/notification"
	"payment-microservices/pkg/config"
	"payment-microservices/pkg/rabbitmq"
)

func main() {
	cfg := config.Load()
	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	handler := notification.NewHandler(conn)
	if err := handler.Start(); err != nil {
		log.Fatalf("Failed to start notification handler: %v", err)
	}
}
