package main

import (
	"fmt"
	commentEntities "interaction-services/domains/comments/entities"
	likeEntities "interaction-services/domains/likes/entities"
	timelineEntities "interaction-services/domains/timelines/entities"
	"interaction-services/wizards"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&commentEntities.Comment{},
		&likeEntities.Like{},
		&timelineEntities.Timeline{},
	)

	router := gin.Default()
	router.RedirectTrailingSlash = false

	wizards.RegisterServer(router)

	port := wizards.Config.Server.Port
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server starting on %s\n", addr)

	if err := router.Run(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
