package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/job_type"
	"github.com/macarrie/relique/internal/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var jobListPageSize int
var jobListSearchModule string
var jobListSearchClient string
var jobListSearchType string
var jobListSearchBackupType string
var jobListSearchStatus string

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
			page := api_helpers.PaginationParams{
				Limit:  uint64(jobListPageSize),
				Offset: 0,
			}
			search := api_helpers.JobSearch{
				ModuleName: jobListSearchModule,
				ClientName: jobListSearchClient,
				JobType:    0,
				BackupType: 0,
				Status:     0,
			}
			if jobListSearchType != "" {
				search.JobType = job_type.FromString(jobListSearchType).Type
			}
			if jobListSearchBackupType != "" {
				search.BackupType = backup_type.FromString(jobListSearchBackupType).Type
			}
			if jobListSearchStatus != "" {
				search.Status = job_status.FromString(jobListSearchStatus).Status
			}

			jobList, err := api.JobList(page, search)
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
			tab.Col("type", "Type", 15)
			tab.Col("start_time", "Start time", 20)
			tab.Col("duration", "Duration", 10)

			format := tab.Print("uuid", "client", "module", "status", "type", "start_time", "duration")
			for _, j := range jobList.Data {
				var jobType string
				if j.JobType.Type == job_type.Backup {
					jobType = fmt.Sprintf("backup/%s", j.BackupType.String())
				} else if j.JobType.Type == job_type.Restore {
					jobType = "restore"
				} else {
					jobType = j.JobType.String()
				}

				fmt.Printf(
					format,
					j.Uuid,
					j.Client.String(),
					j.Module.String(),
					j.Status.String(),
					jobType,
					utils.FormatDatetime(j.StartTime),
					utils.FormatDuration(j.Duration()),
				)
			}

			fmt.Printf("\nShowing %d out of %d records\n", len(jobList.Data), jobList.Count)
		},
	}
	utils.AddPaginationParams(jobListCmd, &jobListPageSize)
	jobListCmd.Flags().StringVarP(&jobListSearchClient, "client", "", "", "Filter on client name")
	jobListCmd.Flags().StringVarP(&jobListSearchModule, "module", "m", "", "Filter on module name")
	jobListCmd.Flags().StringVarP(&jobListSearchStatus, "status", "s", "", "Filter on job status")
	jobListCmd.Flags().StringVarP(&jobListSearchType, "type", "t", "", "Filter on job type")
	jobListCmd.Flags().StringVarP(&jobListSearchBackupType, "backup-type", "", "", "Filter on backup type")

	jobShowCmd := &cobra.Command{
		Use:   "show UUID",
		Short: "Show job details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			j, err := api.JobGet(args[0])
			if err != nil {
				slog.With(
					slog.String("job", args[0]),
					slog.Any("error", err),
				).Error("Cannot get job details")
				os.Exit(1)
			}

			out, err := toml.Marshal(j)
			if err != nil {
				slog.Error("Cannot display job details", slog.Any("error", err))
				os.Exit(1)
			}
			fmt.Printf("\n%v\n", string(out))
		},
	}

	rootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobShowCmd)
}
