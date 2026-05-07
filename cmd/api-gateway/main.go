package main

import (
	"log"
	"net/http"
	"payment-microservices/internal/gateway"
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

	handler := gateway.NewHandler(conn)

	http.HandleFunc("/payments", handler.CreatePaymentsHandler)
	log.Println("API Gateway is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
