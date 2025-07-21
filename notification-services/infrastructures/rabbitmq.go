package infrastructures

import (
	"fmt"
	"log"
	"notification-services/config"
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
