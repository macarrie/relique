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
)

func postBackupRegisterJob(c *gin.Context) {
	var job backup_job.BackupJob
	if err := c.BindJSON(&job); err != nil {
		log.Error("Cannot parse received job")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	job.GetLog().Info("Registering job")

	if previousJob, err := backup_job.GetPreviousJob(job); err != nil || previousJob.Uuid == "" {
		job.BackupType.Type = backup_type.Full
		job.GetLog().Info("No previous backup job found when registering job. This job backup type is now changed to 'full'")
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

	job, err := backup_job.GetByUuid(uuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot retrieve job from db")
	}

	if job.Uuid == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.GetLog().Info("Updating job status")
	job.Status = status
	if _, err := job.UpdateStatus(); err != nil {
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

	job, err := backup_job.GetByUuid(uuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot retrieve job from db")
	}

	if job.Uuid == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	job.Done = done
	job.EndTime = time.Now()
	job.GetLog().Info("Updating job done marker")
	if _, err := job.Save(); err != nil {
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot save job updated done marker")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func postBackupJobApplyDiff(c *gin.Context) {
	var bkpItem backup_item.BackupItem
	if err := c.BindJSON(&bkpItem); err != nil {
		log.Error("Cannot bind received backup item for diff apply")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if bkpItem.JobUuid == "" || bkpItem.Path == "" {
		log.Error("Empty job Uuid or item path received for diff apply request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := bkpItem.ApplyDiff(); err != nil {
		bkpItem.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot apply item diff")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func postBackupJobFile(c *gin.Context) {
	// you can bind multipart form with explicit binding declaration:
	// c.ShouldBindWith(&form, binding.Form)
	// or you can simply use autobinding with ShouldBind method:
	var form backup_item.BackupItemFile
	// in this case proper binding will be automatically selected
	if err := c.Bind(&form); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("cannot parse backup item file")
		c.String(http.StatusBadRequest, "cannot parse backup item file")
		return
	}

	if err := form.Item.SaveFile(form.File); err != nil {
		form.Item.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot save uploaded file")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "ok")
}

func postBackupJobChecksum(c *gin.Context) {
	var bkpItem backup_item.BackupItem
	if err := c.BindJSON(&bkpItem); err != nil {
		log.Error("Cannot bind received backup item for checksum computation")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if bkpItem.JobUuid == "" || bkpItem.Path == "" {
		log.Error("Empty job Uuid or item path received for checksum computation request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := bkpItem.ComputeChecksum(); err != nil {
		bkpItem.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get backup item checksum")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bkpItem)
}

func postBackupJobSignature(c *gin.Context) {
	var bkpItem backup_item.BackupItem
	if err := c.BindJSON(&bkpItem); err != nil {
		log.Error("Cannot bind received backup item for signature computation")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if bkpItem.JobUuid == "" || bkpItem.Path == "" {
		log.Error("Empty job Uuid or item path received for signature computation request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := bkpItem.GetSignature(); err != nil {
		bkpItem.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get backup item signature")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bkpItem)
}

func getBackupJob(c *gin.Context) {
	var params backup_job.JobSearchParams
	if err := c.BindJSON(&params); err != nil {
		log.Error("Cannot bind received job search parameters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	jobs, err := backup_job.Search(params)
	if err != nil {
		params.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot perform job search")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, jobs)
}
