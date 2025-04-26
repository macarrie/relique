package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
)

func webAPIListModules(c *gin.Context) {
	page := getPagination(c)

	mods, err := api.ModuleList(page)
	if err != nil {
		slog.With(
			slog.Any("error", err),
		).Error("Cannot get module list")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, mods)
}

func webAPIGetModule(c *gin.Context) {
	name := c.Param("name")
	mod, err := api.ModuleGet(name)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("name", name),
		).Error("Cannot find module in config")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, mod)
}
