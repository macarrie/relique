package cli

import (
	"os"

	"github.com/macarrie/relique/internal/client"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config/client_daemon_config"
	cliApi "github.com/macarrie/relique/pkg/api/cli"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func Init() {
	rootCmd = &cobra.Command{
		Use:   "relique-client",
		Short: "rsync based backup utility client",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cliApi.InitCommonParams()

			if err := client_daemon_config.Load(cliApi.Params.ConfigPath); err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"path": cliApi.Params.ConfigPath,
				}).Error("Cannot load configuration")
				os.Exit(1)
			}
		},
	}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start relique client",
		Run: func(cmd *cobra.Command, args []string) {
			log.Setup(cliApi.Params.Debug, "relique-client.log")
			log.Info("Starting relique-client")
			client.Run(cliApi.Params)
		},
	}

	// COMMON COMMANDS (CLIENT AND SERVER)
	cliApi.GetCommonCliCommands(rootCmd)

	// DAEMON START
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().BoolVarP(&cliApi.Params.Debug, "debug", "d", false, "debug log output")
	startCmd.PersistentFlags().StringVarP(&cliApi.Params.ConfigPath, "config", "c", "", "Configuration file path")
}

func Execute() error {
	return rootCmd.Execute()
}
