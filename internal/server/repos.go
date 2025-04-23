package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
)

func webAPIListRepos(c *gin.Context) {
	page := getPagination(c)

	repos := api.RepoList(page, api_helpers.RepoSearch{})
	c.JSON(http.StatusOK, repos)
}

func webAPIGetRepo(c *gin.Context) {
	name := c.Param("name")
	repo, err := api.RepoGet(name)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("name", name),
		).Error("Cannot find repo in config")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, repo)
}
