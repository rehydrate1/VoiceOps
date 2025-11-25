package main

import (
	"log"

	"github.com/rehydrate1/VoiceOps/internal/config"
	"github.com/rehydrate1/VoiceOps/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// init config
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// init router
	router := gin.Default()

	// init handlers
	h := handler.NewHandler(cfg)

	// init routes
	router.POST("/api/v1/webhook", h.SberWebhook)

	// start server
	log.Println("Server is starting...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}
