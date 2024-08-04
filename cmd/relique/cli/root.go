package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/logger"
)

var configPath string
var debug bool

var rootCmd = &cobra.Command{
	Use:   "relique",
	Short: "RSync based backup tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(debug)
		if configPath != "" {
			config.UseFile(configPath)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.EnableTraverseRunHooks = true
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Configuration file path")
}
