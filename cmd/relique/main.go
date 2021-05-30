package main

import (
	"github.com/macarrie/relique/cmd/relique/cli"
	log "github.com/macarrie/relique/internal/logging"
)

func main() {
	cli.Init()

	// Setup specific formatted logger for cli display
	log.SetupCliLogger(false)

	// TODO: handle error
	cli.Execute()
}
