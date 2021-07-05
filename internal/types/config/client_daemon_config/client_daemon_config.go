package client_daemon_config

import (
	"github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/relique_job"

	"github.com/macarrie/relique/internal/types/config/common"

	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
)

var Config common.Configuration
var BackupConfig client.Client
var Jobs []relique_job.ReliqueJob

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

// TODO: Configuration validity checks
func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
