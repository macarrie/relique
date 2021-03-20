package scheduler

import (
	"time"

	"github.com/macarrie/relique/internal/types/relique_job"

	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	log "github.com/macarrie/relique/internal/logging"
	serverApi "github.com/macarrie/relique/pkg/api/server"
)

var RunTicker *time.Ticker

func Run() {
	RunTicker = time.NewTicker(10 * time.Second)
	go func() {
		log.Debug("Starting main daemon loop")
		for {
			poll()
			<-RunTicker.C
		}
	}()
}

func poll() {
	if len(serverConfig.Config.Clients) == 0 {
		log.Info("No clients found in configuration")
		return
	}

	for _, client := range serverConfig.Config.Clients {
		if err := serverApi.SendConfiguration(client); err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"client": client.Name,
			}).Error("Cannot send configuration to client")
		}
	}

	activeJobs, err := relique_job.GetActiveJobs()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get currently active jobs")
	}
	if len(activeJobs) == 0 {
		log.Info("No active backup jobs on clients")
	} else {
		log.WithFields(log.Fields{
			"nb": len(activeJobs),
		}).Info("Active jobs on clients")
		for _, job := range activeJobs {
			job.GetLog().Info("Active job on client currently being handled")
		}
	}
}
