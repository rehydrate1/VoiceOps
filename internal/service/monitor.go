package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func CheckSites(urls []string) string {
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
