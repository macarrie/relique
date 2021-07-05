package scheduler

import (
	"time"

	log "github.com/macarrie/relique/internal/logging"
	clientConfig "github.com/macarrie/relique/internal/types/config/client_daemon_config"
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
	if clientConfig.BackupConfig.Version == "" {
		log.Info("Waiting for configuration from relique server")
		return
	}
}
