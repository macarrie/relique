package web

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"

	rsync "github.com/macarrie/relique/internal/lib/rsync"

	"github.com/macarrie/relique/internal/types/job_type"
	"github.com/macarrie/relique/internal/types/sync_task"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/config/server_daemon_config"
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

	if job.JobType.Type == job_type.Restore && job.RestoreJobUuid == "" {
		job.GetLog().Error("Restore job has no target job UUID to restore data from")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if job.BackupType.Type == backup_type.Diff {
		previousJob, err := relique_job.GetPreviousJob(job)
		if err != nil || previousJob.Uuid == "" {
			job.BackupType.Type = backup_type.Full
			job.GetLog().Info("No previous backup job found when registering job. This job backup type is now changed to 'full'")
		}

		job.PreviousJobUuid = previousJob.Uuid
	}

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

func postBackupJobSync(c *gin.Context) {
	uuid := c.Param("uuid")

	job, err := relique_job.GetByUuid(uuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot retrieve job from db")

		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.GetLog().Info("Starting file sync")

	// TODO: Parse stats from output and save them in job
	for _, path := range job.Module.BackupPaths {
		var syncTask sync_task.SyncTask
		if job.JobType.Type == job_type.Backup {
			syncTask = rsync.GetBackupSyncTask(&job, path)
		} else {
			syncTask = rsync.GetRestoreSyncTask(&job, path)
		}
		tasksList := append(server_daemon_config.SyncTasks[job.Uuid], &syncTask)
		server_daemon_config.SyncTasks[job.Uuid] = tasksList

		// TODO: Get stats from output
		go func(task *sync_task.SyncTask, backupPath string) {
			var multiErr *multierror.Error
			defer func() {
				syncTask.Error = multiErr.ErrorOrNil()
				if _, err := job.Save(); err != nil {
					job.GetLog().Error("Could not save job to db during sync_task task update")
				}
			}()

			if err := syncTask.Task.Run(); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": backupPath,
				}).Error("Error during path sync")
				multiErr = multierror.Append(multiErr, err)
				job.Status.Status = job_status.Incomplete
				syncTask.Done = true
			}

			taskLog := syncTask.Task.Log()
			rsyncLogFile, err := job.GetRsyncLogFile(backupPath)
			defer rsyncLogFile.Close()
			if err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": backupPath,
				}).Error("Cannot create log file for sync_task task")
				multiErr = multierror.Append(multiErr, err)
			} else {
				if _, err := rsyncLogFile.WriteString(taskLog.Stdout); err != nil {
					job.GetLog().WithFields(log.Fields{
						"err":  err,
						"path": backupPath,
					}).Error("Cannot write sync_task log to log file")
					multiErr = multierror.Append(multiErr, err)
				}
			}
			rsyncErrorLogFile, err := job.GetRsyncErrorLogFile(backupPath)
			defer rsyncErrorLogFile.Close()
			if err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": backupPath,
				}).Error("Cannot create error log file for sync_task task")
				multiErr = multierror.Append(multiErr, err)
			} else {
				if _, err := rsyncErrorLogFile.WriteString(taskLog.Stderr); err != nil {
					job.GetLog().WithFields(log.Fields{
						"err":  err,
						"path": backupPath,
					}).Error("Cannot write sync_task error log to log file")
					multiErr = multierror.Append(multiErr, err)
				}
			}

			syncTask.Done = true
		}(&syncTask, path)
	}

	c.JSON(http.StatusOK, gin.H{})
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

func getBackupJobSyncProgress(c *gin.Context) {
	uuid := c.Param("uuid")

	job, err := relique_job.GetByUuid(uuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot retrieve job from db")

		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if job.Uuid == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.GetLog().Info("Retrieving job sync progress")
	tasks, ok := server_daemon_config.SyncTasks[job.Uuid]
	if !ok {
		job.GetLog().Error("No sync tasks found for job")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var progress []sync_task.SyncTaskProgress
	for i, _ := range tasks {
		progress = append(progress, sync_task.ProgressFromSyncTask(*tasks[i]))
	}

	c.JSON(http.StatusOK, progress)
}

func putBackupJobDone(c *gin.Context) {
	// TODO: Remove job SyncTask object
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
