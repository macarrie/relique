package cli

import (
	"log/slog"
	"os"

	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Relique web server related commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := api.ConfigGet()
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get relique configuration")
				os.Exit(1)
			}

			if err := db.Init(config.GetDBPath()); err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot initialize database connection")
				os.Exit(1)
			}
		},
	}

	serverStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start relique web server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start(debug, config.Current.WebUI.BindAddr, config.Current.WebUI.Port, config.Current.WebUI.SSLCert, config.Current.WebUI.SSLKey)
		},
	}

	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverStartCmd)
}
