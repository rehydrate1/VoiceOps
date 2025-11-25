package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rehydrate1/VoiceOps/internal/config"
	"github.com/rehydrate1/VoiceOps/internal/models"
	"github.com/rehydrate1/VoiceOps/internal/service"
)

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
	log.Printf("User command: %s", text)

	var pronounceText string
	commandFound := false

	// TODO: научить различать stdout и stderr
	for _, cmd := range h.Cfg.Commands {
		if strings.Contains(text, strings.ToLower(cmd.Phrase)) {
			output, err := service.RemoteExec(
				h.Cfg.SSH.Host,
				h.Cfg.SSH.User,
				h.Cfg.SSH.KeyPath,
				cmd.Script,
			)

			if err != nil {
				log.Printf("Command failed: %v", err)
				pronounceText = fmt.Sprintf("Ошибка выполнения комманды: %v", err)
			} else {
				output = strings.TrimSpace(output)

				if strings.Contains(cmd.Response, "%s") {
					pronounceText = fmt.Sprintf(cmd.Response, output)
				} else {
					pronounceText = cmd.Response
				}
			}

			commandFound = true
			break
		}
	}

	if !commandFound {
		if text == "" || strings.Contains(text, "запусти") || strings.Contains(text, "открой") {
			pronounceText = "VoiceOps на связи. Скажите 'Проверь прод', чтобы начать диагностику."
		} else if strings.Contains(text, "проверь") {
			pronounceText = service.CheckSites(h.Cfg.Monitoring.URLs)
		} else {
			pronounceText = "Я вас не поняла. Скажите 'Проверь прод' или 'Перезагрузи бота'."
		}
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