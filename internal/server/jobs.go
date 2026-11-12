package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/job"
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

func webAPIPostJobRetry(c *gin.Context) {
	uuid := c.Param("uuid")

	if _, err := api.JobGet(uuid); err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("uuid", uuid),
		).Error("Cannot find job in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	retry, err := api.JobRetrySetup(uuid)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("uuid", uuid),
		).Error("Cannot retry job")
		if errors.Is(err, job.ErrRestoreNotRetryable) || errors.Is(err, job.ErrJobNotFinished) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusAccepted, retry)

	// The job can run for a long time: run it in background so the client can follow it from the new job details
	go func(j job.Job) {
		if err := api.JobRetryRun(&j); err != nil {
			j.GetLog().With(
				slog.Any("error", err),
				slog.String("original_uuid", uuid),
			).Error("Error during retried job")
		}
	}(retry)
}
