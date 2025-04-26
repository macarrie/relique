package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
)

func webAPIListJobs(c *gin.Context) {
	page := getPagination(c)
	search := getJobSearchParams(c)

	jobs, err := api.JobList(page, search)
	if err != nil {
		slog.With(
			slog.Any("error", err),
		).Error("Cannot get job list")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, jobs)
}

func webAPIGetJob(c *gin.Context) {
	uuid := c.Param("uuid")
	job, err := api.JobGet(uuid)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("uuid", uuid),
		).Error("Cannot find job in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, job)
}
