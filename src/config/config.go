package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIAddress            string
	RabbitMQURL           string
	PaymentsQueue         string
	PaymentProcessedQueue string
	MessagesLimit         int
}

func Load() *Config {
	return &Config{
		APIAddress:            stringValue("API_GATEWAY_ADDR", ":8080"),
		RabbitMQURL:           stringValue("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		PaymentsQueue:         stringValue("PAYMENTS_QUEUE", "payments"),
		PaymentProcessedQueue: stringValue("PAYMENT_PROCESSED_QUEUE", "payment_processed"),
		MessagesLimit:         intValue("MESSAGE_LIMIT", 100),
	}
}

func stringValue(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func intValue(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return result
}
