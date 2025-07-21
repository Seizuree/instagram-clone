package wizards

import (
	"encoding/json"
	"log"
	"notification-services/domains/notifications"
	"notification-services/events"
	"notification-services/infrastructures"
)

type NotificationConsumer struct {
	rmq *infrastructures.RabbitMQ
	uc  notifications.NotificationUseCase
}

func NewNotificationConsumer(rmq *infrastructures.RabbitMQ, uc notifications.NotificationUseCase) *NotificationConsumer {
	return &NotificationConsumer{rmq: rmq, uc: uc}
}

func (c *NotificationConsumer) Start() {
	queues := []string{
		"notification.like.created",
		"notification.comment.created",
		"notification.follow.created",
	}

	for _, queue := range queues {
		go c.consume(queue)
	}

	log.Println("Notification consumers started...")
	select {} // Block forever
}

func (c *NotificationConsumer) consume(queueName string) {
	ch, err := c.rmq.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue '%s': %v", queueName, err)
	}

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		log.Printf("Received message from queue '%s'", queueName)
		switch queueName {
		case "notification.like.created":
			var event events.LikeCreatedEvent
			json.Unmarshal(d.Body, &event)
			c.uc.CreateLikeNotification(&event)
		case "notification.comment.created":
			var event events.CommentCreatedEvent
			json.Unmarshal(d.Body, &event)
			c.uc.CreateCommentNotification(&event)
		case "notification.follow.created":
			var event events.FollowCreatedEvent
			json.Unmarshal(d.Body, &event)
			c.uc.CreateFollowNotification(&event)
		}
	}
}
