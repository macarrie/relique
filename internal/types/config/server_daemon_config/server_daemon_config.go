package server_daemon_config

import (
	client2 "github.com/macarrie/relique/internal/types/client"

	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/macarrie/relique/internal/types/config/common"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
)

var customConfigFilePath string
var customConfigFile bool

var Config common.Configuration

func Load(filePath string) error {
	if filePath != "" {
		common.UseFile(filePath)
	}
	conf, err := common.Load("server")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Cannot load configuration")
		return err
	}

	clients, err := client2.LoadFromPath(conf.ClientCfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": conf.ClientCfgPath,
		}).Fatal("Cannot load clients from configuration")
		return err
	}
	conf.Clients = clients

	schedules, err := schedule.LoadFromPath(conf.SchedulesCfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": conf.SchedulesCfgPath,
		}).Fatal("Cannot load schedules from configuration")
		return err
	}
	conf.Schedules = schedules

	conf.Version = uuid.New().String()

	Config = conf

	return nil
}

func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
