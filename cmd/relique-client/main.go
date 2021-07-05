package main

import (
	"github.com/macarrie/relique/cmd/relique-client/cli"
	log "github.com/macarrie/relique/internal/logging"
)

func main() {
	// Set log without debug enabled as we do not have parsed cli params yet
	log.SetupCliLogger(false, false)

	cli.Init()

	// TODO: handle error
	cli.Execute()
}
