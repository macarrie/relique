package main

import (
	"github.com/macarrie/relique/cmd/relique/cli"
	log "github.com/macarrie/relique/internal/logging"
)

func main() {
	cli.Init()

	// For relique cli, do not write into log file
	log.Setup(false, "")

	// TODO: handle error
	cli.Execute()
}
