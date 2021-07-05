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

		v1.POST("/jobs/", getBackupJob)
		v1.POST("/jobs/start", postJobStart)

		v1.POST("/retention/clean/", postRetentionClean)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
