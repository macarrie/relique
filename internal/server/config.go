package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/config"
)

func webAPIPostConfigReload(c *gin.Context) {
	if err := api.ConfigLoad(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func webAPIGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Current)
}

func webAPIGetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": api.ConfigGetVersion(),
	})
}
