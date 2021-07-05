package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/server/scheduler"
	"github.com/macarrie/relique/internal/types/config/server_daemon_config"
)

func postRetentionClean(c *gin.Context) {
	err := scheduler.CleanRetention(server_daemon_config.Config.RetentionPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot clean jobs retention")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
