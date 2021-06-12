package web

import (
	"net/http"

	clientApi "github.com/macarrie/relique/pkg/api/client"

	"github.com/macarrie/relique/internal/types/client"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
)

func getRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ping)
		v1.POST("/check_server_connection", postCheckServerConnection)
		v1.POST("/config", postConfig)
		v1.GET("/config/version", getConfigVersion)
		v1.POST("/job/start", postJobStart)
		v1.POST("/retention/clean", postRetentionClean)
	}

	return router
}

func ping(c *gin.Context) {
	c.Status(http.StatusOK)
}

func postCheckServerConnection(c *gin.Context) {
	var serverPingParams client.ServerPingParams
	if err := c.ShouldBind(&serverPingParams); err != nil {
		log.Error("Cannot bind server ping parameters received")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status := http.StatusOK
	ret, err := clientApi.CheckServerConnection(serverPingParams)
	if err != nil {
		status = http.StatusBadRequest
	}

	c.JSON(status, ret)
}
