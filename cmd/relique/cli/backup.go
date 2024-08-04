package cli

import (
	"log/slog"
	"os"

	"github.com/macarrie/relique/api"
	"github.com/spf13/cobra"
)

var backupClient string
var backupModule string

func init() {
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup related commands",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := api.ConfigGet()
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get relique configuration")
				os.Exit(1)
			}
		},
	}
	backupCmd.Flags().StringVarP(&backupClient, "client", "", "", "Client to backup")
	backupCmd.Flags().StringVarP(&backupModule, "module", "m", "", "Module to use")
	backupCmd.MarkFlagRequired("client")

	rootCmd.AddCommand(backupCmd)
}
