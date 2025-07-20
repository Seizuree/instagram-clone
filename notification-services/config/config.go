package config

import (
	"log"
	"os"
)

type Config struct {
	RabbitMQURL string
}

func Load() *Config {
	rmqURL := os.Getenv("RABBITMQ_URL")
	if rmqURL == "" {
		rmqURL = "amqp://guest:guest@localhost:5672/"
		log.Println("⚠️  Using default RabbitMQ URL")
	}

	return &Config{
		RabbitMQURL: rmqURL,
	}
}
