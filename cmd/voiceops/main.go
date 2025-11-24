package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SberRequest struct {
	MessageID int64  `json:"messageId"`
	SessionID string `json:"sessionId"`
	Uuid      struct {
		UserID      interface{} `json:"userId"`
		UserChannel string      `json:"userChannel"`
		Sub         string      `json:"sub"`
	} `json:"uuid"`
	Payload   struct {
		Message struct {
			OriginalText string `json:"original_text"`
		} `json:"message"`
	} `json:"payload"`
}

func main() {
	router := gin.Default()

	router.POST("/api/v1/webhook", SberHandler)

	log.Println("Server is starting...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error while starting server: %v", err)
		os.Exit(1)
	}
}

func SberHandler(c *gin.Context) {
	var req SberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})

		return
	}

	log.Printf("Received a request. MessageID: %d", req.MessageID)

	response := gin.H{
		"messageName": "ANSWER_TO_USER",
		"sessionId": req.SessionID,
		"messageId": req.MessageID,
		"uuid":        req.Uuid,
		"payload": gin.H{
			"device": gin.H{
				"deviceId": "sber-boom-home",
			},
			"pronounceText": "Связь установлена! Я готова управлять сервером.",
			"items": []gin.H{
				{
					"bubble": gin.H{
						"text": "Связь установлена! Я готова управлять сервером.",
					},
				},
			},
			"finished": false,
			"auto_listening": true,
		},
	}

	c.JSON(http.StatusOK, response)
}
