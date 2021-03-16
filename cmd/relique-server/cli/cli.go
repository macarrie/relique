package cli

import (
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/server"
	"github.com/spf13/cobra"
)

var Params = server.CliArgs{}
var rootCmd *cobra.Command

func Init() {
	rootCmd = &cobra.Command{
		Use:   "relique-server",
		Short: "rsync based backup utility main server",
	}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start relique server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Setup(Params.Debug, "relique-server.log")
			log.Info("Starting relique-server")
			server.Run(Params)
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.PersistentFlags().BoolVarP(&Params.Debug, "debug", "d", false, "debug log output")
	rootCmd.PersistentFlags().StringVarP(&Params.ConfigPath, "config", "c", "", "Configuration file path")
}

func Execute() error {
	return rootCmd.Execute()
}
