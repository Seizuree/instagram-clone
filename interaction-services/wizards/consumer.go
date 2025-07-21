package wizards

import (
	"encoding/json"
	"interaction-services/domains/timelines"
	"interaction-services/events"
	"interaction-services/infrastructures"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// StartPostCreatedConsumer starts a consumer for post creation events.
func StartPostCreatedConsumer(rabbitMQ *infrastructures.RabbitMQ, timelineUseCase timelines.TimelineUseCase) {
	queueName := "post.created"

	handler := func(d amqp091.Delivery) {
		var event events.PostCreatedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("ERROR: Failed to unmarshal post.created event: %v", err)
			return
		}

		// The timeline use case will contain the logic to fan-out the post.
		if err := timelineUseCase.AddPostToFollowerTimelines(&event); err != nil {
			log.Printf("ERROR: Failed to process post.created event for postID %s: %v", event.PostID, err)
		}

		log.Printf("Successfully processed post.created event for postID: %s", event.PostID)
	}

	rabbitMQ.StartConsumer(queueName, handler)
	log.Printf("Post Creation Consumer has been started.")
}
