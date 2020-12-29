package web

import (
	"fmt"
	"net/http"

	"github.com/macarrie/relique/internal/client/scheduler"

	clientConfig "github.com/macarrie/relique/internal/types/config/client_daemon_config"
	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/types/backup_type"

	"github.com/gin-gonic/gin"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_job"
)

func postBackupStart(c *gin.Context) {
	var params backup_job.JobSearchParams
	if err := c.BindJSON(&params); err != nil {
		c.String(http.StatusBadRequest, "cannot parse received job start parameters")
		return
	}
	log.Info("TODO: Start backup job")
	bType := backup_type.FromString(params.BackupType)
	if bType.Type == backup_type.Unknown {
		c.String(http.StatusBadRequest, "unknown backup type received")
		return
	}
	// TODO: Validate module
	// IF mod present in configuration, use this one
	var targetModule module.Module
	var moduleFound bool = false

	// Check for module in client configuration
	for _, mod := range clientConfig.BackupConfig.Modules {
		if mod.Name == params.Module {
			mod.GetLog().Info("Using module found in client configuration for manual backup")
			moduleFound = true
			targetModule = mod
			targetModule.BackupType = bType
		}
	}

	if !moduleFound {
		log.Info("Module not found in client configuration for manual backup. Checking if a module with this name is installed on client")
		targetModule = module.Module{
			ModuleType: params.Module,
			Name:       fmt.Sprintf("ondemand-%s-%s", params.Module, bType.String()),
			BackupType: bType,
		}
		if err := targetModule.LoadDefaultConfiguration(); err != nil {
			c.String(http.StatusBadRequest, "cannot load module default configuration")
			return
		}
	}

	job := scheduler.AddJob(targetModule)

	c.JSON(http.StatusOK, job)
}
