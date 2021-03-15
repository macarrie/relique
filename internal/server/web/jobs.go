package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_item"
	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/relique_job"
)

func postBackupRegisterJob(c *gin.Context) {
	var job relique_job.ReliqueJob
	if err := c.BindJSON(&job); err != nil {
		log.Error("Cannot parse received job")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	job.GetLog().Info("Registering job")

	if job.BackupType.Type == backup_type.Diff {
		previousJob, err := backup_job.GetPreviousJob(job)
		if err != nil || previousJob.Uuid == "" {
			job.BackupType.Type = backup_type.Full
			job.PreviousJobStorageDestination = ""
			job.GetLog().Info("No previous backup job found when registering job. This job backup type is now changed to 'full'")
		} else {
			job.PreviousJobStorageDestination = filepath.Clean(fmt.Sprintf("%s/%s", server_daemon_config.Config.BackupStoragePath, previousJob.Uuid))
		}

	}

	job.StorageDestination = filepath.Clean(fmt.Sprintf("%s/%s", server_daemon_config.Config.BackupStoragePath, job.Uuid))
	job.StartTime = time.Now()

	_, err := job.Save()
	if err != nil {
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot save registered job")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, job)
}

func putBackupJobStatus(c *gin.Context) {
	uuid := c.Param("uuid")

	var status job_status.JobStatus
	if err := c.ShouldBind(&status); err != nil {
		log.Error("Cannot parse received status")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	job, err := relique_job.GetByUuid(uuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot retrieve job from db")
	}

	if job.Uuid == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.GetLog().WithFields(log.Fields{
		"updated_status": status,
	}).Info("Updating job status")
	job.Status = status
	if _, err := job.Save(); err != nil {
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot save job updated status")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func putBackupJobDone(c *gin.Context) {
	uuid := c.Param("uuid")

	var done bool
	if err := c.BindJSON(&done); err != nil {
		log.Error("Cannot parse received bool")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	job, err := relique_job.GetByUuid(uuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot retrieve job from db")
	}

	if job.Uuid == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.GetLog().WithFields(log.Fields{
		"updated_done_marker": done,
	}).Info("Updating job done marker")
	job.Done = done
	job.EndTime = time.Now()
	if _, err := job.Save(); err != nil {
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot save job updated done marker")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func getBackupJob(c *gin.Context) {
	var params relique_job.JobSearchParams
	if err := c.BindJSON(&params); err != nil {
		log.Error("Cannot bind received job search parameters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	jobs, err := relique_job.Search(params)
	if err != nil {
		params.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot perform job search")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, jobs)
}
