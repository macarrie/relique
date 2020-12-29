package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/config/client_daemon_config"
)

func postConfig(c *gin.Context) {
	var clientConfig clientObject.Client
	if err := c.ShouldBind(&clientConfig); err != nil {
		log.Error("Cannot bind config received")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"version": clientConfig.Version,
	}).Info("Received new configuration from server")

	client_daemon_config.BackupConfig = clientConfig
	c.JSON(http.StatusOK, gin.H{})
}

func getConfigVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": client_daemon_config.BackupConfig.Version})
}
