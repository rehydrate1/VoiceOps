package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/api/v1/webhook", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"text": "I'm working"})
	})

	log.Println("Server is starting...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error while starting server: %v", err)
		os.Exit(1)
	}
}