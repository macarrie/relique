package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd := &cobra.Command{
		Use:   "job",
		Short: "Backup job related commands",
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

	jobListCmd := &cobra.Command{
		Use:   "list",
		Short: "List backup and restore jobs",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Handle pagination
			jobList, err := api.JobList(api_helpers.PaginationParams{})
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get job list")
				os.Exit(1)
			}

			tab := tabular.New()
			tab.Col("uuid", "UUID", 40)
			tab.Col("client", "Client", 10)
			tab.Col("module", "Module", 15)
			tab.Col("status", "Status", 10)
			tab.Col("backup_type", "Backup type", 15)
			tab.Col("start_time", "Start time", 20)
			tab.Col("duration", "Duration", 10)

			format := tab.Print("uuid", "client", "module", "status", "backup_type", "start_time", "duration")
			for _, j := range jobList.Data {
				fmt.Printf(
					format,
					j.Uuid,
					j.Client.String(),
					j.Module.String(),
					j.Status.String(),
					j.BackupType.String(),
					utils.FormatDatetime(j.StartTime),
					utils.FormatDuration(j.Duration()),
				)
			}

			// TODO: Handle pagination
			fmt.Printf("\nShowing %d out of %d records\n", jobList.Count, jobList.Count)
		},
	}

	jobShowCmd := &cobra.Command{
		Use:   "show UUID",
		Short: "Show job details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cl, err := api.ClientGet(args[0])
			if err != nil {
				slog.With(
					slog.String("client", args[0]),
					slog.Any("error", err),
				).Error("Cannot get client details")
				os.Exit(1)
			}

			out, err := toml.Marshal(cl)
			if err != nil {
				slog.Error("Cannot display client details", slog.Any("error", err))
				os.Exit(1)
			}
			fmt.Printf("\n%v\n", string(out))
		},
	}

	rootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobShowCmd)
}
