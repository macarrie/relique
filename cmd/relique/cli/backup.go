package cli

import (
	"log/slog"
	"os"
	"strings"

	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	"github.com/spf13/cobra"
)

var backupClient string
var backupModule string
var backupRepo string
var backupInclusions []string
var backupExclusions []string
var backupExcludeCVS bool

func init() {
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup related commands",
		Args:  cobra.ArbitraryArgs,
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

			var mod module.Module
			if backupModule == "" {
				if len(args) == 0 {
					slog.Error("Path arguments are needed if no module specified. Use --module or add arguments to command to announce paths to backup")
					os.Exit(1)
				} else {
					mod = module.Module{
						Name:        "on-demand",
						BackupType:  backup_type.BackupType{Type: backup_type.Diff},
						ModuleType:  "generic",
						BackupPaths: args,
						Exclude:     backupExclusions,
						ExcludeCVS:  backupExcludeCVS,
					}
				}
			} else {
				mod, err = module.GetByName(c.Modules, backupModule)
				if err != nil {
					slog.With(
						slog.Any("error", err),
						slog.String("module", backupModule),
					).Error("Cannot find module on client")
					os.Exit(1)
				}

				// Override module with provided cli params
				if len(backupInclusions) > 0 {
					slog.With(
						slog.String("module", mod.Name),
						slog.String("cli_include", strings.Join(backupInclusions, ", ")),
						slog.String("mod_include", strings.Join(mod.Include, ", ")),
					).Info("Override module include list with the one provided in CLI params")
					mod.Include = backupInclusions
				}
				if len(backupExclusions) > 0 {
					slog.With(
						slog.String("module", mod.Name),
						slog.String("cli_exclude", strings.Join(backupExclusions, ", ")),
						slog.String("mod_exclude", strings.Join(mod.Exclude, ", ")),
					).Info("Override module exclusions list with the one provided in CLI params")
					mod.Exclude = backupExclusions
				}
				if backupExcludeCVS {
					slog.With(
						slog.String("module", mod.Name),
						slog.Bool("cli_exclude_cvs", backupExcludeCVS),
						slog.Bool("mod_exclude_cvs", mod.ExcludeCVS),
					).Info("Override module exclude_cvs parameter")
					mod.ExcludeCVS = backupExcludeCVS
				}
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
	backupCmd.Flags().StringSliceVarP(&backupInclusions, "include", "i", []string{}, "File inclusions")
	backupCmd.Flags().StringSliceVarP(&backupExclusions, "exclude", "e", []string{}, "File exclusions")
	backupCmd.Flags().BoolVarP(&backupExcludeCVS, "exclude-cvs", "", false, "Exclude CVS from file selection")
	backupCmd.MarkFlagRequired("client")

	rootCmd.AddCommand(backupCmd)
}
