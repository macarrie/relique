package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/macarrie/relique/internal/db"

	config "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/server/scheduler"
	"github.com/macarrie/relique/internal/server/web"
)

type CliArgs struct {
	Debug      bool
	ConfigPath string
}

func Run(args CliArgs) {
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

	scheduler.Run()
	go web.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		switch sig := <-signalChannel; sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Info("Signal received. Shutting down...")
			web.Stop()
			os.Exit(0)
		}
	}
}
