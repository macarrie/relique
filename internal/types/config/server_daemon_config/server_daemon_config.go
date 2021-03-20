package server_daemon_config

import (
	"github.com/macarrie/relique/internal/types/sync_task"
	"github.com/pkg/errors"

	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/macarrie/relique/internal/types/config/common"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
)

var Config common.Configuration
var SyncTasks map[string][]*sync_task.SyncTask

func init() {
	SyncTasks = make(map[string][]*sync_task.SyncTask)
}

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

	schedules, err := schedule.LoadFromPath(conf.SchedulesCfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": conf.SchedulesCfgPath,
		}).Fatal("Cannot load schedules from configuration")
		return err
	}
	conf.Schedules = schedules

	clients, err := clientObject.LoadFromPath(conf.ClientCfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": conf.ClientCfgPath,
		}).Fatal("Cannot load clients from configuration")
		return err
	}

	clients, err = clientObject.FillSchedulesStruct(clients, schedules)
	if err != nil {
		return errors.Wrap(err, "cannot match schedules chosen in client definitions with schedules definitions")
	}

	conf.Clients = clients

	conf.Version = uuid.New().String()

	Config = conf

	return nil
}

// TODO: Configuration validity checks
func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
