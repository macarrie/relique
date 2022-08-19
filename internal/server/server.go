package server

import (
	"os"
	"os/signal"
	"syscall"

	cliApi "github.com/macarrie/relique/pkg/api/cli"

	"github.com/macarrie/relique/internal/db"

	config "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/server/scheduler"
	"github.com/macarrie/relique/internal/server/web"
)

func Run(args cliApi.Args) {
	if err := config.Load(args.ConfigPath); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
	}
	if err := db.Init(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot init database")
	}

	config.SaveConfigObjectsInDb()

	if err := scheduler.LoadRetention(config.Config.RetentionPath); err != nil {
		log.WithFields(log.Fields{
			"path": config.Config.RetentionPath,
			"err":  err,
		}).Error("Cannot load relique jobs retention. Relique will start without previous done jobs in memory, jobs previously already performed might be restarted")
	}

	scheduler.Run()
	go web.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		switch sig := <-signalChannel; sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Info("Signal received. Shutting down...")

			if err := scheduler.UpdateRetention(config.Config.RetentionPath); err != nil {
				log.WithFields(log.Fields{
					"path": config.Config.RetentionPath,
					"err":  err,
				}).Error("Cannot update jobs retention. Done jobs will not be remembered and might be restarted at relique client restart")
			}

			web.Stop()
			os.Exit(0)
		}
	}
}
