package infrastructures

import (
	"core-services/config"
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.Ch.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) PublishMessage(queueName string, body []byte) error {
	return r.Ch.Publish(
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
