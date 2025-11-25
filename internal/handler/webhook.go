package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rehydrate1/VoiceOps/internal/config"
	"github.com/rehydrate1/VoiceOps/internal/models"
	"github.com/rehydrate1/VoiceOps/internal/service"
)

var urls = []string{
	"https://google.com",
	"https://github.com",
	"https://non-existent-site.ru",
}

type Handler struct {
	Cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{Cfg: cfg}
}

func (h *Handler) SberWebhook(c *gin.Context) {
	var req models.SberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})

		return
	}

	log.Printf("Received a request. MessageID: %d", req.MessageID)

	text := strings.ToLower(req.Payload.Message.OriginalText)

	var pronounceText string
	if text == "" || strings.Contains(text, "запусти") || strings.Contains(text, "открой") {
		pronounceText = "VoiceOps на связи. Скажите 'Проверь прод', чтобы начать диагностику."
	} else if strings.Contains(text, "проверь") {
		pronounceText = service.CheckSites(urls)
	} else {
		pronounceText = "Я вас не поняла. Скажите 'Проверь прод'."
	}

	response := gin.H{
		"messageName": "ANSWER_TO_USER",
		"sessionId":   req.SessionID,
		"messageId":   req.MessageID,
		"uuid":        req.Uuid,
		"payload": gin.H{
			"device": gin.H{
				"deviceId": "sber-boom-home",
			},
			"pronounceText": pronounceText,
			"items": []gin.H{
				{
					"bubble": gin.H{
						"text": pronounceText,
					},
				},
			},
			"finished":       false,
			"auto_listening": true,
		},
	}

	c.JSON(http.StatusOK, response)
}