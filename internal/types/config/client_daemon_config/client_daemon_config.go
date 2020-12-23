package client_daemon_config

import (
	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/types/config/common"

	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
)

var customConfigFilePath string
var customConfigFile bool

var Config common.Configuration
var BackupConfig client.Client
var Jobs []backup_job.BackupJob

func Load(filePath string) error {
	if filePath != "" {
		common.UseFile(filePath)
	}
	conf, err := common.Load("client")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Cannot load configuration")
		return err
	}

	Config = conf

	return nil
}

func JobExists(module module.Module) bool {
	for _, backupJob := range Jobs {
		if backupJob.Module.Name == module.Name && backupJob.Module.ModuleType == module.ModuleType {
			return true
		}
	}

	return false
}

func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
