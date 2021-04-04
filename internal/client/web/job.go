package web

import (
	"fmt"
	"net/http"

	"github.com/macarrie/relique/internal/types/job_type"

	"github.com/macarrie/relique/internal/client/scheduler"

	clientConfig "github.com/macarrie/relique/internal/types/config/client_daemon_config"
	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/types/backup_type"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/relique_job"
)

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
	var moduleFound bool = false

	// Check for module in client configuration
	for _, mod := range clientConfig.BackupConfig.Modules {
		if mod.Name == params.Module {
			mod.GetLog().Info("Using module found in client configuration for manual job start")
			moduleFound = true
			targetModule = mod
			targetModule.BackupType = bType
		}
	}

	if !moduleFound {
		log.Info("Module not found in client configuration for manual job start. Checking if a module with this name is installed on client")
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
			c.String(http.StatusBadRequest, "cannot load module default configuration")
			return
		}
	}

	job := relique_job.New(clientConfig.BackupConfig, targetModule, jType)
	job.RestoreJobUuid = params.RestoreJobUuid
	job.RestoreDestination = params.RestoreDestination
	scheduler.AddJob(job)

	c.JSON(http.StatusOK, job)
}
