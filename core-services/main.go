package main

import (
	followEntities "core-services/domains/follow/entities"
	userEntities "core-services/domains/users/entities"
	"core-services/wizards"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&userEntities.User{},
		&followEntities.Follow{},
	)

	router := gin.Default()

	wizards.RegisterServer(router)

	port := wizards.Config.Server.Port
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server starting on %s\n", addr)

	if err := router.Run(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
