package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/module"
)

func getModules(c *gin.Context) {
	installedModules, err := module.GetLocallyInstalled()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get locally installed module list")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, installedModules)
}

func getModule(c *gin.Context) {
	modName := c.Param("name")

	installedModules, err := module.GetLocallyInstalled()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get locally installed module list")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, m := range installedModules {
		if m.Name == modName {
			c.JSON(http.StatusOK, m)
			return
		}
	}

	c.AbortWithStatus(http.StatusNotFound)
}
