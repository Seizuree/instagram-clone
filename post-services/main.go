package main

import (
	"fmt"
	"post-services/domains/posts/entities"
	"post-services/wizards"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&entities.Post{},
	)

	wizards.StartUserDeletedConsumer(wizards.RabbitMQ, wizards.PostUseCase)

	router := gin.Default()
	router.MaxMultipartMemory = 64 << 20
	router.RedirectTrailingSlash = false

	wizards.RegisterServer(router)

	port := wizards.Config.Server.Port
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server starting on %s\n", addr)

	if err := router.Run(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
