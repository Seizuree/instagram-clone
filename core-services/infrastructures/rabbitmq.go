package infrastructures

import (
	"context"
	"core-services/config"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ holds the connection and channel for RabbitMQ operations.
type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

var (
	rabbitOnce sync.Once
	rabbitInst *RabbitMQ
)

// NewRabbitMQ establishes a connection to RabbitMQ and returns a singleton instance.
func NewRabbitMQ(conf *config.Config) *RabbitMQ {
	rabbitOnce.Do(func() {
		dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
			conf.Rabbit.User, conf.Rabbit.Password, conf.Rabbit.Host, conf.Rabbit.Port)

		conn, err := amqp.Dial(dsn)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
		}

		rabbitInst = &RabbitMQ{
			Conn: conn,
			Ch:   ch,
		}
		log.Println("RabbitMQ connected successfully.")
	})

	return rabbitInst
}

// PublishJSON publishes a message to a specific queue in JSON format.
// It declares the queue to ensure it exists.
func (r *RabbitMQ) PublishJSON(ctx context.Context, queueName string, data interface{}) error {
	_, err := r.Ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.Ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// StartConsumer starts a consumer on a specified queue.
// It takes a handler function to process incoming messages.
func (r *RabbitMQ) StartConsumer(queueName string, handler func(d amqp.Delivery)) {
	q, err := r.Ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue for consumer: %s", err)
	}

	msgs, err := r.Ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message from queue %s", queueName)
			handler(d)
		}
	}()

	log.Printf("Consumer started on queue: %s. Waiting for messages.", q.Name)
}

// Close gracefully closes the RabbitMQ channel and connection.
func (r *RabbitMQ) Close() {
	if r.Ch != nil {
		r.Ch.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
	log.Println("RabbitMQ connection closed.")
}