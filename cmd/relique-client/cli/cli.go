package cli

import (
	"github.com/macarrie/relique/internal/client"
	log "github.com/macarrie/relique/internal/logging"

	"github.com/spf13/cobra"
)

var Params = client.CliArgs{}
var rootCmd *cobra.Command

func Init() {
	rootCmd = &cobra.Command{
		Use:   "relique-client",
		Short: "rsync based backup utility client agent",
	}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start relique client",
		Run: func(cmd *cobra.Command, args []string) {
			log.Setup(Params.Debug, "relique-client.log")
			log.Info("Starting relique-client")
			client.Run(Params)
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.PersistentFlags().BoolVarP(&Params.Debug, "debug", "d", false, "debug log output")
	rootCmd.PersistentFlags().StringVarP(&Params.ConfigPath, "config", "c", "", "Configuration file path")
}

func Execute() error {
	return rootCmd.Execute()
}
