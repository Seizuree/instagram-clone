package wizards

import (
	"encoding/json"
	"log"
	"post-services/domains/posts"
	"post-services/infrastructures"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// UserDeletedEvent represents the message structure for a user deletion.
type UserDeletedEvent struct {
	UserID string `json:"user_id"`
}

// StartUserDeletedConsumer starts a consumer that listens for user deletion events.
func StartUserDeletedConsumer(rabbitMQ *infrastructures.RabbitMQ, postUseCase posts.PostUseCase) {
	queueName := "user.deleted"

	handler := func(d amqp.Delivery) {
		var event UserDeletedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("ERROR: Failed to unmarshal user.deleted event: %v", err)
			// We could potentially move this to a dead-letter queue.
			return
		}

		userID, err := uuid.Parse(event.UserID)
		if err != nil {
			log.Printf("ERROR: Invalid user ID format in user.deleted event: %v", err)
			return
		}

		if err := postUseCase.DeletePostsByUserID(userID); err != nil {
			log.Printf("ERROR: Failed to process user.deleted event for userID %s: %v", event.UserID, err)
			// In a production system, this might require manual intervention or a retry mechanism.
		}

		log.Printf("Successfully processed user.deleted event for userID: %s", event.UserID)
	}

	// The StartConsumer function (from infrastructures/rabbitmq.go) runs in a goroutine,
	// so this call is non-blocking.
	rabbitMQ.StartConsumer(queueName, handler)

	log.Printf("User Deletion Consumer has been started.")
}
