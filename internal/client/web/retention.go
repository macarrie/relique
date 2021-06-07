package web

import (
	"net/http"

	"github.com/macarrie/relique/internal/client/scheduler"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config/client_daemon_config"

	"github.com/gin-gonic/gin"
)

func postRetentionClean(c *gin.Context) {
	err := scheduler.CleanRetention(client_daemon_config.Config.RetentionPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot clean jobs retention")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
