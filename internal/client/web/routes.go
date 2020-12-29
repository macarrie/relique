package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)
		v1.POST("/config", postConfig)
		v1.GET("/config/version", getConfigVersion)
		v1.POST("/backup/start", postBackupStart)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
