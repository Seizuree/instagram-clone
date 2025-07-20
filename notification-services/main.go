package main

import (
	"log"
	"notification-service/config"
	"notification-service/http"
	"notification-service/infrastructures"
	"notification-service/rabbitmq"
	"notification-service/repositories"
	"notification-service/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	rmq := infrastructures.NewRabbitMQ(cfg)
	repo := repositories.NewNotificationRepository()
	uc := usecases.NewNotificationUseCase(repo)

	go rabbitmq.ConsumeNotifications(rmq, uc)

	router := gin.Default()
	http.RegisterWebSocketRoutes(router, uc)

	log.Println("Notification service running on :8083")
	router.Run(":8083")
}
