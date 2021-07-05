package client

import (
	"os"
	"os/signal"
	"syscall"

	cliApi "github.com/macarrie/relique/pkg/api/cli"

	"github.com/macarrie/relique/internal/client/scheduler"
	"github.com/macarrie/relique/internal/client/web"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config/client_daemon_config"
)

func Run(args cliApi.Args) {
	if err := client_daemon_config.Load(args.ConfigPath); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Cannot load configuration file")
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
