package main

import (
	"log"
	"net/http"
	"payment-microservices/src/config"
	"payment-microservices/src/internal/adapters/rabbitmq"
	httppayment "payment-microservices/src/internal/api/http/payment"
	apppayment "payment-microservices/src/internal/app/payment"
)

func main() {
	cfg := config.Load()

	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	publisher := rabbitmq.NewPaymentPublisher(cfg, conn)
	createPayment := apppayment.NewCreatePaymentUseCase(publisher)
	handler := httppayment.NewHandler(createPayment)

	mux := http.NewServeMux()
	mux.HandleFunc("/payments", handler.CreatePayment)

	log.Printf("API Gateway is running on %s", cfg.APIAddress)
	if err := http.ListenAndServe(cfg.APIAddress, mux); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
