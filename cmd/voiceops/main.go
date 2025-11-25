package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rehydrate1/VoiceOps/internal/config"
	"github.com/rehydrate1/VoiceOps/internal/service"

	"github.com/gin-gonic/gin"
)

var urls = []string{
	"https://google.com",
	"https://github.com",
	"https://non-existent-site.ru",
}

var cfg *config.Config

type SberRequest struct {
	MessageID int64  `json:"messageId"`
	SessionID string `json:"sessionId"`
	Uuid      struct {
		UserID      interface{} `json:"userId"`
		UserChannel string      `json:"userChannel"`
		Sub         string      `json:"sub"`
	} `json:"uuid"`
	Payload struct {
		Message struct {
			OriginalText string `json:"original_text"`
		} `json:"message"`
	} `json:"payload"`
}

func main() {
	// init config
	var err error
	cfg, err = config.LoadConfig("./configs/config.yaml")
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
	router.POST("/api/v1/webhook", SberHandler)

	// start server
	log.Println("Server is starting...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error while starting server: %v", err)
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

	text := strings.ToLower(req.Payload.Message.OriginalText)

	var pronounceText string
	if text == "" || strings.Contains(text, "запусти") || strings.Contains(text, "открой") {
		pronounceText = "VoiceOps на связи. Скажите 'Проверь прод', чтобы начать диагностику."
	} else if strings.Contains(text, "проверь") {
		pronounceText = checkSites(urls)
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

func checkSites(urls []string) string {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	badSites := []string{}

	clean := func(u string) string {
		u = strings.Replace(u, "https://", "", 1)
		u = strings.Replace(u, "http://", "", 1)
		return u
	}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Failed to get %s: %v", url, err)
			badSites = append(badSites, clean(url))
			continue
		}

		if resp.StatusCode != 200 {
			badSites = append(badSites, clean(url))
		}

		resp.Body.Close()
	}

	if len(badSites) == 0 {
		return "Все системы работают штатно. Ошибок не обнаружено"
	}

	return fmt.Sprintf("Обнаружены проблемы с сайтами: %s", strings.Join(badSites, ", "))
}
