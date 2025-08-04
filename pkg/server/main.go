package server

import (
	"github.com/argon-chat/sentinel/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() error {
	router := gin.Default()

	router.GET(config.Instance.Route, postHandler)

	return router.Run(config.Instance.Port)
}

func postHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
