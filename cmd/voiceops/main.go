package main

import (
	"log"

	"github.com/rehydrate1/VoiceOps/internal/config"
	"github.com/rehydrate1/VoiceOps/internal/handler"
	"github.com/rehydrate1/VoiceOps/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// init config
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// test ssh
	output, err := service.RemoteExec(
		cfg.SSH.Host,
		cfg.SSH.User,
		cfg.SSH.KeyPath,
		"uptime",
	)

	if err != nil {
		log.Printf("SSH Test Failed: %v", err)
	} else {
		log.Printf("SSH Test Success! Server uptime: %s", output)
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
