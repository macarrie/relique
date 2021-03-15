package cli

import (
	"os"

	"github.com/macarrie/relique/internal/types/displayable"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/config/common"
	"github.com/macarrie/relique/internal/types/relique_job"
	"github.com/macarrie/relique/pkg/api/cli"
	"github.com/spf13/cobra"
)

var jobSearchParams relique_job.JobSearchParams
var manualJobParams relique_job.JobSearchParams
var config common.Configuration
var rootCmd *cobra.Command
var jsonOutput bool

func cliInitParams() {
	if jsonOutput {
		displayable.DisplayMode = displayable.JSON
	} else {
		displayable.DisplayMode = displayable.TUI
	}

	err := cli.PingDaemon(config)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot contact relique daemon. Relique daemon must be started and available")
		os.Exit(1)
	}
}

func Init() {
	rootCmd = &cobra.Command{
		Use:   "relique",
		Short: "rsync based backup utility command line interface",
	}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Server related commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if config.Port == 0 {
				config.Port = 8433
			}

			cliInitParams()
		},
	}
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "Client related commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if config.Port == 0 {
				config.Port = 8434
			}

			cliInitParams()
		},
	}

	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "Perform job related operations",
	}
	jobListCmd := &cobra.Command{
		Use:   "list",
		Short: "List jobs on relique server",
		Run: func(cmd *cobra.Command, args []string) {
			jobs, err := cli.SearchJob(config, jobSearchParams)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot list jobs from server")
				os.Exit(1)
			}

			disp := make([]displayable.Displayable, len(jobs))
			for i, v := range jobs {
				disp[i] = v
			}
			displayable.Table(disp)
		},
	}

	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Perform backup related operations on relique client",
	}
	backupStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a manual backup on relique client",
		Run: func(cmd *cobra.Command, args []string) {
			if config.Port == 0 {
				config.Port = 8434
			}
			manualJobParams.JobType = "backup"
			job, err := cli.ManualJobStart(config, manualJobParams)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot start manual backup")
				os.Exit(1)
			}
			displayable.Details(job)
		},
	}

	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Perform restore related operations on relique client",
	}
	restoreStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a manual restore on relique client",
		Run: func(cmd *cobra.Command, args []string) {
			if config.Port == 0 {
				config.Port = 8434
			}
			manualJobParams.JobType = "restore"
			manualJobParams.BackupType = "restore"
			job, err := cli.ManualJobStart(config, manualJobParams)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot start manual restore")
				os.Exit(1)
			}
			displayable.Details(job)
		},
	}

	// ROOT CMD
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output content as JSON")

	// SERVER CMD
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&config.PublicAddress, "address", "localhost", "Relique server address")
	serverCmd.PersistentFlags().Uint32VarP(&config.Port, "port", "p", 0, "Relique server port")

	//// JOBS CMD
	serverCmd.AddCommand(jobsCmd)
	jobsCmd.AddCommand(jobListCmd)
	jobListCmd.Flags().StringVarP(&jobSearchParams.Module, "module", "m", "", "Module name")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Client, "client", "c", "", "Client name")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Status, "status", "s", "", "Job status")
	jobListCmd.Flags().StringVarP(&jobSearchParams.BackupType, "backup-type", "t", "", "Backup type (diff, full)")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Uuid, "uuid", "u", "", "Job with UUID")
	jobListCmd.Flags().IntVar(&jobSearchParams.Limit, "limit", 0, "Limit job search to LIMIT items (0 corresponds to no limit)")

	// CLIENT CMD
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringVar(&config.PublicAddress, "address", "localhost", "Relique client address")
	clientCmd.PersistentFlags().Uint32VarP(&config.Port, "port", "p", 0, "Relique client port")

	//// BACKUP CMD
	clientCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(backupStartCmd)
	backupStartCmd.Flags().StringVarP(&manualJobParams.Module, "module", "m", "", "Module name")
	backupStartCmd.Flags().StringVarP(&manualJobParams.BackupType, "backup-type", "t", "", "Backup type (diff, full)")
	backupStartCmd.MarkFlagRequired("module")
	backupStartCmd.MarkFlagRequired("backup-type")

	//// RESTORE CMD
	clientCmd.AddCommand(restoreCmd)
	restoreCmd.AddCommand(restoreStartCmd)
	restoreStartCmd.Flags().StringVarP(&manualJobParams.Module, "module", "m", "", "Module name")
	restoreStartCmd.Flags().StringVarP(&manualJobParams.RestoreJobUuid, "job", "j", "", "Job UUID to restore data from")
	restoreStartCmd.Flags().StringVarP(&manualJobParams.RestoreDestination, "destination", "d", "", "Alternate file restore destination")
	restoreStartCmd.MarkFlagRequired("module")
	restoreStartCmd.MarkFlagRequired("job")
	// TODO: Add option to restore files to alternate location. Allows to avoid running relique daemon as root to restore root owned files

}

func Execute() error {
	return rootCmd.Execute()
}
