package server

import (
	"github.com/argon-chat/sentinel/pkg/config"
	"github.com/gin-gonic/gin"
)

func Run() error {
	router := gin.Default()
	router.POST(config.Instance.Route, postHandler)
	return router.Run(config.Instance.Port)
}

func postHandler(c *gin.Context) {
	appID := c.GetHeader("app_id")
	if appID == "" {
		c.JSON(400, gin.H{"error": "app_id header is required"})
		return
	}

	dsn, ok := config.Instance.Projects[appID]
	if !ok {
		c.JSON(400, gin.H{"error": "invalid app_id"})
		return
	}
	_, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to read request body"})
		return
	}

}
