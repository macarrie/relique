package server_daemon_config

import (
	"github.com/macarrie/relique/internal/lib/rsync"
	"github.com/pkg/errors"

	"github.com/macarrie/relique/internal/db"
	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/schedule"

	"github.com/macarrie/relique/internal/types/config/common"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/relique/internal/logging"
)

var Config common.Configuration
var SyncTasks map[string][]*rsync.Rsync

func init() {
	SyncTasks = make(map[string][]*rsync.Rsync)
}

func SaveConfigObjectsInDb() {
	// Save clients
	for _, cl := range Config.Clients {
		if _, err := cl.Save(nil); err != nil {
			cl.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Cannot save client found in configuration into database")
		} else {
			cl.GetLog().Debug("Saved client found in configuration into database")
		}

		for _, mod := range cl.Modules {
			if _, err := mod.Save(nil); err != nil {
				mod.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Cannot save client module found in configuration into database")
			} else {
				mod.GetLog().Debug("Saved client module found in configuration into database")
			}
		}
	}
}

func Load(filePath string) error {
	if filePath != "" {
		common.UseFile(filePath)
	}
	conf, err := common.Load("server")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot load configuration")
		return err
	}
	conf.Version = uuid.New().String()

	schedules, err := schedule.LoadFromPath(conf.SchedulesCfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": conf.SchedulesCfgPath,
		}).Error("Cannot load schedules from configuration")
		return err
	}
	conf.Schedules = schedules

	clients, err := clientObject.LoadFromPath(conf.ClientCfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": conf.ClientCfgPath,
		}).Error("Cannot load clients from configuration")
		return err
	}

	clients, err = clientObject.FillSchedulesStruct(clients, schedules)
	if err != nil {
		return errors.Wrap(err, "cannot match schedules chosen in client definitions with schedules definitions")
	}

	clients = clientObject.FillServerPublicAddress(clients, conf.PublicAddress, conf.Port)
	clients = clientObject.FillConfigVersion(clients, conf.Version)
	clients = clientObject.InitAliveStatus(clients)

	conf.Clients = clients

	Config = conf

	// Set DB path
	db.DbPathReadInConfig = true
	db.DbPath = conf.DbPath

	return nil
}

// TODO: Configuration validity checks
func Check() error {
	var errorList *multierror.Error

	return errorList.ErrorOrNil()
}
