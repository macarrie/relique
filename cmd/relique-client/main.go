package main

import (
	"github.com/macarrie/relique/cmd/relique-client/cli"
	log "github.com/macarrie/relique/internal/logging"
)

func main() {
	cli.Init()

	// Set log without debug enabled as we do not have parsed cli params yet
	log.Setup(false, "relique-client.log")

	// TODO: handle error
	cli.Execute()
}
