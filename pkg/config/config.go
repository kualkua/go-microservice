package config

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL   string
	MessagesLimit int
}

func Load() *Config {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	messagesLimit := os.Getenv("MESSAGE_LIMIT")
	if messagesLimit == "" {
		messagesLimit = "100"
	}

	messagesLimitInt, err := strconv.Atoi(messagesLimit)
	if err != nil {
		messagesLimitInt = 100
	}

	return &Config{
		RabbitMQURL:   rabbitMQURL,
		MessagesLimit: messagesLimitInt,
	}
}
