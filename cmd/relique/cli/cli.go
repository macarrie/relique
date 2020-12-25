package cli

import (
	"os"

	"github.com/macarrie/relique/internal/tui"
	"github.com/macarrie/relique/internal/types/displayable"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/internal/types/config/common"
	"github.com/macarrie/relique/pkg/api/cli"
	"github.com/spf13/cobra"
)

var jobSearchParams backup_job.JobSearchParams
var config common.Configuration
var rootCmd *cobra.Command
var jsonOutput bool

func Init() {
	rootCmd = &cobra.Command{
		Use:   "relique",
		Short: "rsync based backup utility command line interface",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if jsonOutput {
				displayable.DisplayMode = displayable.JSON
			} else {
				displayable.DisplayMode = displayable.TUI
				tui.Init()
			}

			err := cli.PingServer(config)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot contact relique server. Relique server must be started and available")
				os.Exit(1)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !jsonOutput {
				defer tui.Close()
			}
		},
	}
	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "Perform job related operations",
	}
	jobListCmd := &cobra.Command{
		Use:   "list",
		Short: "List jobs",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Handle error
			jobs, err := cli.SearchJob(config, jobSearchParams)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot list jobs from server")
			}

			disp := make([]displayable.Displayable, len(jobs))
			for i, v := range jobs {
				disp[i] = v
			}
			displayable.Table(disp)
		},
	}

	rootCmd.PersistentFlags().StringVar(&config.PublicAddress, "server", "localhost", "Relique server address")
	rootCmd.PersistentFlags().Uint32VarP(&config.Port, "port", "p", 8433, "Relique server port")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output content as JSON")

	rootCmd.AddCommand(jobsCmd)

	jobsCmd.AddCommand(jobListCmd)
	jobListCmd.Flags().StringVarP(&jobSearchParams.Module, "module", "m", "", "Module name")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Client, "client", "c", "", "Client name")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Status, "status", "s", "", "Job status")
	jobListCmd.Flags().StringVarP(&jobSearchParams.BackupType, "backup_type", "t", "", "Backup type (diff, full)")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Uuid, "uuid", "u", "", "Job with UUID")
	jobListCmd.Flags().IntVar(&jobSearchParams.Limit, "limit", 0, "Limit job search to LIMIT items (0 corresponds to no limit)")

}

func Execute() error {
	return rootCmd.Execute()
}
