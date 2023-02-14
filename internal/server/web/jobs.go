package web

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	"github.com/macarrie/relique/internal/types/client"

	"github.com/macarrie/relique/internal/server/scheduler"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/job_type"
	"github.com/macarrie/relique/internal/types/module"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/relique_job"
)

func getJob(c *gin.Context) {
	uuid := c.Param("uuid")
	job, getJobErr := relique_job.GetByUuid(uuid)
	if getJobErr != nil {
		log.WithFields(log.Fields{
			"uuid": uuid,
		}).Error("Cannot find job in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, job)
}

func searchJob(c *gin.Context) {
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

func postJobStart(c *gin.Context) {
	var params relique_job.JobSearchParams
	if err := c.BindJSON(&params); err != nil {
		c.String(http.StatusBadRequest, "cannot parse received job start parameters")
		return
	}
	jType := job_type.FromString(params.JobType)
	bType := backup_type.FromString(params.BackupType)
	if bType.Type == backup_type.Unknown && jType.Type == job_type.Backup {
		c.String(http.StatusBadRequest, "unknown backup type received")
		return
	}
	var targetModule module.Module
	moduleFound := false

	// Check that client exists
	var targetClient client.Client
	clientFound := false
	for _, cl := range serverConfig.Config.Clients {
		if cl.Name == params.Client {
			clientFound = true
			targetClient = cl
			break
		}
	}

	if !clientFound {
		c.String(http.StatusBadRequest, fmt.Sprintf("Cannot find client '%s' in relique server configuration", params.Client))
		return
	}

	targetClient.GetLog().Info("Found requested client in server configuration for manual job start")

	// Check for module in client configuration
	for _, mod := range targetClient.Modules {
		if mod.Name == params.Module {
			mod.GetLog().Info("Using module found in configuration for manual job start")
			moduleFound = true
			targetModule = mod
			targetModule.BackupType = bType
		}
	}

	if !moduleFound {
		log.Info("Module not found in configuration for manual job start. Checking if a module with this name is installed on server")
		if jType.Type == job_type.Backup {
			targetModule = module.Module{
				ModuleType: params.Module,
				Name:       fmt.Sprintf("ondemand-%s-%s-%s", params.Module, jType.String(), bType.String()),
				BackupType: bType,
			}
		} else {
			targetModule = module.Module{
				ModuleType: params.Module,
				Name:       fmt.Sprintf("ondemand-%s-%s", params.Module, jType.String()),
				BackupType: bType,
			}
		}
		if err := targetModule.LoadDefaultConfiguration(); err != nil {
			c.String(
				http.StatusBadRequest,
				fmt.Sprintf("Cannot load module configuration. Check that the module '%s' is correctly installed on relique server and client", params.Module),
			)
			return
		}
	}

	targetModule.Variant = params.ModuleVariant
	job := relique_job.New(&targetClient, targetModule, jType)
	job.RestoreJobUuid = params.RestoreJobUuid
	job.RestoreDestination = params.RestoreDestination
	job.StorageRoot = serverConfig.Config.BackupStoragePath
	scheduler.AddJob(job)

	c.JSON(http.StatusOK, job)
}

func getJobLogs(c *gin.Context) {
	uuid := c.Param("uuid")
	job, getJobErr := relique_job.GetByUuid(uuid)
	if getJobErr != nil {
		log.WithFields(log.Fields{
			"uuid": uuid,
		}).Error("Cannot find job in database")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	backupPathSelectionEnc := c.Query("bp")
	backupPathSelection, urlDecodeErr := url.QueryUnescape(backupPathSelectionEnc)
	if urlDecodeErr != nil {
		log.WithFields(log.Fields{
			"err": urlDecodeErr,
		}).Error("Cannot decode backup path from URL parameter")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	bpFound := false
	for _, path := range job.Module.BackupPaths {
		if backupPathSelection == path {
			bpFound = true
			break
		}
	}

	if !bpFound {
		log.Error("Cannot find requested backup path in job module")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// TODO: Filter from backup path
	logPath := job.GetRsyncLogFilePath(backupPathSelection)
	f, err := os.Open(logPath)
	if err != nil {
		log.WithFields(log.Fields{
			"uuid": uuid,
			"err":  err,
		}).Error("Cannot open log file for reading")
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Log file does not exist")
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fileStats, err := f.Stat()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get log file stats")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.DataFromReader(http.StatusOK, fileStats.Size(), "text", f, nil)
}
