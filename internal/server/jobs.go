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

func webAPIDeleteJob(c *gin.Context) {
	uuid := c.Param("uuid")
	cleanImage := c.Query("clean_image") == "true"

	j, err := api.JobGet(uuid)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("uuid", uuid),
		).Error("Cannot find job in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !j.Done {
		j.GetLog().Error("Cannot remove job that is not finished")
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "job is not finished"})
		return
	}

	if err := api.JobRemove(uuid, cleanImage); err != nil {
		j.GetLog().With(
			slog.Any("error", err),
			slog.Bool("clean_image", cleanImage),
		).Error("Cannot remove job")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
