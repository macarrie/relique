package client

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/macarrie/relique/internal/client/scheduler"
	"github.com/macarrie/relique/internal/client/web"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config/client_daemon_config"
)

type CliArgs struct {
	Debug      bool
	ConfigPath string
}

func Run(args CliArgs) {
	if err := client_daemon_config.Load(args.ConfigPath); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
	}

	if err := scheduler.LoadRetention(client_daemon_config.Config.RetentionPath); err != nil {
		log.WithFields(log.Fields{
			"path": client_daemon_config.Config.RetentionPath,
			"err":  err,
		}).Error("Cannot load relique client retention. Relique will start without previous done jobs in memory, jobs previously already performed might be restarted")
	}

	scheduler.Run()
	go web.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		switch sig := <-signalChannel; sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Info("Signal received. Shutting down...")
			if err := scheduler.UpdateRetention(client_daemon_config.Config.RetentionPath); err != nil {
				log.WithFields(log.Fields{
					"path": client_daemon_config.Config.RetentionPath,
					"err":  err,
				}).Error("Cannot update jobs retention. Done jobs will not be remembered and might be restarted at relique client restart")
			}
			web.Stop()
			os.Exit(0)
		}
	}
}
