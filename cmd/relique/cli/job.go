package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/job"
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

var jobRemoveOlderThan string
var jobRemoveSearchModule string
var jobRemoveSearchClient string
var jobRemoveSearchType string
var jobRemoveSearchBackupType string
var jobRemoveSearchStatus string
var jobRemoveCleanImage bool

func buildJobSearch(moduleName, clientName, jobType, backupType, status string) api_helpers.JobSearch {
	search := api_helpers.JobSearch{
		ModuleName: moduleName,
		ClientName: clientName,
	}
	if jobType != "" {
		search.JobType = job_type.FromString(jobType).Type
	}
	if backupType != "" {
		search.BackupType = backup_type.FromString(backupType).Type
	}
	if status != "" {
		search.Status = job_status.FromString(status).Status
	}

	return search
}

// buildJobRemovalPlan selects the jobs to remove and the images linked to them, without modifying anything
func buildJobRemovalPlan(args []string) (removalPlan, error) {
	var selected []job.Job
	if len(args) == 1 {
		j, err := api.JobGet(args[0])
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot get job details: %w", err)
		}
		if !j.Done {
			return removalPlan{}, fmt.Errorf("cannot remove job with status '%s' because it is not finished", j.Status.String())
		}
		selected = append(selected, j)
	} else {
		before, cutoff, err := olderThanCutoff(jobRemoveOlderThan)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot parse --older-than: %w", err)
		}
		search := buildJobSearch(jobRemoveSearchModule, jobRemoveSearchClient, jobRemoveSearchType, jobRemoveSearchBackupType, jobRemoveSearchStatus)
		search.Before = before
		fmt.Printf("Selecting jobs started before %s\n", utils.FormatDatetime(cutoff))

		jobList, err := api.JobList(api_helpers.PaginationParams{Limit: 0}, search)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot get job list: %w", err)
		}
		for _, j := range jobList.Data {
			if !j.Done {
				j.GetLog().Warn("Job is not finished, skipping it")
				continue
			}
			selected = append(selected, j)
		}
	}

	plan := removalPlan{JobsToRemove: selected}
	for _, j := range selected {
		if j.JobType.Type != job_type.Backup {
			continue
		}
		exists, err := api.ImageExists(j.Uuid)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot check if job '%s' has an image: %w", j.Uuid, err)
		}
		if !exists {
			continue
		}
		img, err := api.ImageGet(j.Uuid)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot get image linked to job '%s': %w", j.Uuid, err)
		}
		if jobRemoveCleanImage {
			plan.ImagesToRemove = append(plan.ImagesToRemove, img)
		} else {
			plan.ImagesKept = append(plan.ImagesKept, img)
		}
	}

	return plan, nil
}

// printRestoreRetryHelp explains why a restore job is not retried and displays the command to run it again manually
func printRestoreRetryHelp(uuid string) {
	j, err := api.JobGet(uuid)
	if err != nil {
		slog.With(
			slog.String("job", uuid),
			slog.Any("error", err),
		).Error("Cannot get job details")
		os.Exit(1)
	}

	fmt.Printf(`Job '%s' is a restore job. Restores overwrite data on the target client, so relique does not retry them automatically.
To run the same restore, start it yourself. You are responsible for any side effects on the target client:

    %s

If the original restore used custom restore paths, add them as arguments (source:destination) because they are not recorded with the job.
`, j.Uuid, j.RestoreCommand())
}

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
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List backup and restore jobs",
		Run: func(cmd *cobra.Command, args []string) {
			page := api_helpers.PaginationParams{
				Limit:  uint64(jobListPageSize),
				Offset: 0,
			}
			search := buildJobSearch(jobListSearchModule, jobListSearchClient, jobListSearchType, jobListSearchBackupType, jobListSearchStatus)

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
				fmt.Printf(
					format,
					j.Uuid,
					j.Client.String(),
					j.Module.String(),
					j.Status.String(),
					jobTypeString(j),
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

	jobRemoveCmd := &cobra.Command{
		Use:     "remove [UUID]",
		Aliases: []string{"rm"},
		Short:   "Remove finished jobs, by UUID or by age",
		Long: `Remove a finished job by UUID, or every finished job older than a given age with --older-than.

The image generated by a backup job is kept unless --clean-image is used.
A recap of everything that will be removed is displayed and must be confirmed before anything is deleted.`,
		Example: `  relique job remove 1b4e28ba-2fa1-11d2-883f-0016d3cca427
  relique job remove 1b4e28ba-2fa1-11d2-883f-0016d3cca427 --clean-image
  relique job remove --older-than 90d --client my-client --yes`,
		Args: validateRemoveArgs("client", "module", "status", "type", "backup-type"),
		Run: func(cmd *cobra.Command, args []string) {
			plan, err := buildJobRemovalPlan(args)
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot select jobs to remove")
				os.Exit(1)
			}

			if !confirmRemoval(plan) {
				return
			}

			uuids := make([]string, 0, len(plan.JobsToRemove))
			for _, j := range plan.JobsToRemove {
				uuids = append(uuids, j.Uuid)
			}
			runRemoval("job", uuids, func(uuid string) error {
				return api.JobRemove(uuid, jobRemoveCleanImage)
			})
		},
	}
	jobRemoveCmd.Flags().StringVarP(&jobRemoveOlderThan, "older-than", "", "", "Remove finished jobs started before this age (e.g. 90d, 2w, 12h)")
	jobRemoveCmd.Flags().StringVarP(&jobRemoveSearchClient, "client", "", "", "With --older-than: filter on client name")
	jobRemoveCmd.Flags().StringVarP(&jobRemoveSearchModule, "module", "m", "", "With --older-than: filter on module name")
	jobRemoveCmd.Flags().StringVarP(&jobRemoveSearchStatus, "status", "s", "", "With --older-than: filter on job status")
	jobRemoveCmd.Flags().StringVarP(&jobRemoveSearchType, "type", "t", "", "With --older-than: filter on job type")
	jobRemoveCmd.Flags().StringVarP(&jobRemoveSearchBackupType, "backup-type", "", "", "With --older-than: filter on backup type")
	jobRemoveCmd.Flags().BoolVarP(&jobRemoveCleanImage, "clean-image", "", false, "Also remove the image generated by the job")
	jobRemoveCmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "Skip confirmation (the recap is still displayed)")

	jobRetryCmd := &cobra.Command{
		Use:   "retry UUID",
		Short: "Retry a finished backup job with the same parameters",
		Long: `Start a new backup job with the same client, module and repository as a finished job.

Restore jobs are never retried automatically because they overwrite data on the target client.
For a restore job, the equivalent 'relique restore' command is displayed instead.`,
		Example: `  relique job retry 1b4e28ba-2fa1-11d2-883f-0016d3cca427`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			retry, err := api.JobRetry(args[0])
			if errors.Is(err, job.ErrRestoreNotRetryable) {
				printRestoreRetryHelp(args[0])
				os.Exit(1)
			}
			if err != nil {
				log := slog.With(
					slog.String("job", args[0]),
					slog.Any("error", err),
				)
				if retry.Uuid != "" {
					log = log.With(slog.String("retry_job", retry.Uuid))
				}
				log.Error("Cannot retry job")
				os.Exit(1)
			}

			fmt.Printf("Job '%s' retried as new job '%s' (status: %s)\n", args[0], retry.Uuid, retry.Status.String())
		},
	}

	rootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobShowCmd)
	jobCmd.AddCommand(jobRemoveCmd)
	jobCmd.AddCommand(jobRetryCmd)
}
