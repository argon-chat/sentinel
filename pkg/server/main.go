package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/argon-chat/sentinel/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() error {
	router := gin.Default()
	corsConfig := cors.Config{
		AllowOrigins: config.Instance.AllowedOrigins,
		AllowMethods: []string{"POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", config.Instance.Header},
	}
	router.Use(cors.New(corsConfig))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.POST(config.Instance.Server.Route, postHandler)
	return router.Run(fmt.Sprintf(":%d", config.Instance.Server.Port))
}

func postHandler(c *gin.Context) {
	appID := "test" //c.GetHeader(config.Instance.Header)
	if appID == "" {
		c.JSON(400, gin.H{"error": "Sec-Ner header is required"})
		return
	}
	project, ok := config.Instance.Projects[appID]
	if !ok {
		c.JSON(400, gin.H{"error": "invalid app_id"})
		return
	}
	envelope, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to read request body"})
		return
	}
	upstreamSentryURL := fmt.Sprintf("%s/api/%s/envelope/?sentry_key=%s", config.Instance.SentryUrl, project.SentryProjectId, project.SentryKey)
	resp, err := post(upstreamSentryURL, envelope)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to send request to sentry" + err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %s", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "sentry did not return OK"})
	}
	c.Status(200)
}

func post(url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	return client.Do(req)
}
