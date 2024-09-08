package cli

import (
	"log/slog"
	"os"

	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	"github.com/spf13/cobra"
)

var backupClient string
var backupModule string
var backupRepo string

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

			if err := db.Init(config.GetDBPath()); err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot initialize database connection")
				os.Exit(1)
			}

			c, err := api.ClientGet(backupClient)
			if err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", backupClient),
				).Error("Cannot find client")
				os.Exit(1)
			}

			mod, err := module.GetByName(c.Modules, backupModule)
			if err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("module", backupModule),
				).Error("Cannot find module on client")
				os.Exit(1)
			}

			var r repo.Repository
			if backupRepo == "" {
				r, err = repo.GetDefault(config.Current.Repositories)
				slog.Debug("Repository not provided. Looking for default repo in configuration")
				if err != nil {
					slog.With(
						slog.Any("error", err),
					).Error("Cannot find default repository in config")
					os.Exit(1)
				}
			} else {
				r, err = repo.GetByName(config.Current.Repositories, backupRepo)
				if err != nil {
					slog.With(
						slog.Any("error", err),
						slog.String("repository", backupRepo),
					).Error("Cannot find repository in config")
					os.Exit(1)
				}
			}

			// TODO: Get backup type from cli param or module
			if err := api.BackupStart(c, mod, r); err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", c.Name),
					slog.String("module", mod.Name),
					slog.String("repository", r.GetName()),
				).Error("Error during backup job")
				os.Exit(1)
			}
		},
	}
	backupCmd.Flags().StringVarP(&backupClient, "client", "", "", "Client to backup")
	backupCmd.Flags().StringVarP(&backupModule, "module", "m", "", "Module to use")
	backupCmd.Flags().StringVarP(&backupModule, "repo", "r", "", "Repository to use")
	backupCmd.MarkFlagRequired("client")

	rootCmd.AddCommand(backupCmd)
}
