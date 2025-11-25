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

	for _, cmd := range h.Cfg.Commands {
		if strings.Contains(text, strings.ToLower(cmd.Phrase)) {
			stdout, stderr, err := service.RemoteExec(
				h.Cfg.SSH.Host,
				h.Cfg.SSH.User,
				h.Cfg.SSH.KeyPath,
				cmd.Script,
			)
			stdout = strings.TrimSpace(stdout)
			stderr = strings.TrimSpace(stderr)

			if err != nil {
				log.Printf("Command failed. Error: %v | Stderr: %s", err, stderr)
				pronounceText = fmt.Sprintf("Ошибка выполнения комманды: %v", err)

				if stderr != "" {
					pronounceText = fmt.Sprintf("Ошибка выполнения: %s", stderr)
				} else {
					pronounceText = fmt.Sprintf("Команда упала: %v", err)
				}
			} else {
				log.Printf("Command Success. Stdout: %s", stdout)
				if strings.Contains(cmd.Response, "%s") {
					if stdout == "" {
						stdout = "нет данных"
					}
					pronounceText = fmt.Sprintf(cmd.Response, stdout)
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
			pronounceText = "VoiceOps на связи"
		} else if strings.Contains(text, "проверь") {
			pronounceText = service.CheckSites(h.Cfg.Monitoring.URLs)
		} else {
			pronounceText = "Я вас не поняла. Скажите 'Список команд' чтобы узнать команды."
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