package main

import (
	"log"
	"token-hub/config"
	"token-hub/internal/repository"
	"token-hub/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	repository.InitDatabase()

	gin.SetMode(config.AppConfig.Server.Mode)

	r := router.SetupRouter()

	log.Printf("Server starting on port " + config.AppConfig.Server.Port)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
