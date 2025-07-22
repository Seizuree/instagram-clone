package main

import (
	"fmt"
	"log"
	"notification-services/domains/notifications/entities"
	"notification-services/wizards"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(&entities.Notification{})

	// Start the WebSocket hub in a background goroutine to manage client connections.
	go wizards.Hub.Run()

	// Create a new RabbitMQ consumer and start it in a background goroutine.
	// It listens for events and passes them to the use case for processing.
	consumer := wizards.NewNotificationConsumer(wizards.RabbitMQ, wizards.NotificationUseCase)
	go consumer.Start()

	router := gin.Default()
	router.RedirectTrailingSlash = false

	wizards.RegisterServer(router)

	port := wizards.Config.Server.Port
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Notification service is running on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start notification service: %v", err)
	}
}
