package cli

import (
	"os"

	"github.com/macarrie/relique/internal/types/config/common"

	"github.com/macarrie/relique/internal/db"

	"github.com/macarrie/relique/internal/types/client"
	cliApi "github.com/macarrie/relique/pkg/api/cli"
	serverApi "github.com/macarrie/relique/pkg/api/server"

	"github.com/macarrie/relique/internal/types/relique_job"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/server"
	"github.com/macarrie/relique/internal/types/displayable"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var config common.Configuration

var jobSearchParams relique_job.JobSearchParams
var manualJobParams relique_job.JobSearchParams

func Init() {
	rootCmd = &cobra.Command{
		Use:   "relique-server",
		Short: "rsync based backup utility main server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cliApi.InitCommonParams()
			if err := db.Open(false); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot open relique database")
				os.Exit(1)
			}
		},
	}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start relique server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Setup(cliApi.Params.Debug, "relique-server.log")
			log.Info("Starting relique-server")
			server.Run(cliApi.Params)
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
			jobs, err := relique_job.Search(jobSearchParams)
			if err != nil {
				jobSearchParams.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Cannot perform job search")
				os.Exit(1)
			}

			disp := make([]displayable.Displayable, len(jobs))
			for i, v := range jobs {
				disp[i] = v
			}
			displayable.Table(disp)
		},
	}

	pingClientCmd := &cobra.Command{
		Use:   "ping",
		Short: "Checks SSH connection from server and client",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			clientAddr := args[0]
			c := client.Client{
				Name:    clientAddr,
				Address: clientAddr,
			}
			if err := serverApi.PingSSHClient(c); err != nil {
				c.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Cannot ping client")
				os.Exit(1)
			}

			c.GetLog().Info("Ping successful")
		},
	}

	retentionCmd := &cobra.Command{
		Use:   "retention",
		Short: "Client jobs retention related commands",
	}

	retentionCleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean server jobs retention",
		Run: func(cmd *cobra.Command, args []string) {
			if config.Port == 0 {
				config.Port = 8433
			}
			err := serverApi.CleanRetention(config)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot clean retention on server")
				os.Exit(1)
			}
			log.Info("Server jobs retention cleaned successfully")
		},
	}

	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Perform backup related operations",
	}
	backupStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a manual backup job",
		Run: func(cmd *cobra.Command, args []string) {
			manualJobParams.JobType = "backup"
			job, err := cliApi.ManualJobStart(config, manualJobParams)
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
		Short: "Perform restore related operations",
	}
	restoreStartCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a manual restore on relique client",
		Run: func(cmd *cobra.Command, args []string) {
			manualJobParams.JobType = "restore"
			manualJobParams.BackupType = "restore"
			job, err := cliApi.ManualJobStart(config, manualJobParams)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot start manual restore")
				os.Exit(1)
			}
			displayable.Details(job)
		},
	}

	// COMMON COMMANDS (CLIENT AND SERVER)
	cliApi.GetCommonCliCommands(rootCmd)

	// DAEMON START
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&cliApi.Params.ConfigPath, "config", "c", "", "Configuration file path")
	startCmd.PersistentFlags().BoolVarP(&cliApi.Params.Debug, "debug", "d", false, "debug log output")

	// JOBS CMD
	rootCmd.AddCommand(jobsCmd)
	jobsCmd.AddCommand(jobListCmd)
	jobListCmd.Flags().StringVarP(&jobSearchParams.Module, "module", "m", "", "Module name")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Client, "client", "c", "", "Client name")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Status, "status", "s", "", "Job status")
	jobListCmd.Flags().StringVarP(&jobSearchParams.BackupType, "backup-type", "t", "", "Backup type (diff, full)")
	jobListCmd.Flags().StringVarP(&jobSearchParams.Uuid, "uuid", "u", "", "Job with UUID")
	jobListCmd.Flags().IntVarP(&jobSearchParams.Limit, "limit", "l", 0, "Limit job search to LIMIT items (0 corresponds to no limit)")

	// PING_SERVER CMD
	rootCmd.AddCommand(pingClientCmd)

	// BACKUP CMD
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&config.PublicAddress, "addr", "", "localhost", "Server address")
	backupCmd.Flags().Uint32VarP(&config.Port, "port", "p", 8433, "Server port")
	backupCmd.AddCommand(backupStartCmd)
	backupStartCmd.Flags().StringVarP(&manualJobParams.Client, "client", "", "", "Client name to backup from")
	backupStartCmd.Flags().StringVarP(&manualJobParams.Module, "module", "m", "", "Module name")
	backupStartCmd.Flags().StringVarP(&manualJobParams.BackupType, "backup-type", "t", "", "Backup type (diff, full)")
	backupStartCmd.MarkFlagRequired("client")
	backupStartCmd.MarkFlagRequired("module")
	backupStartCmd.MarkFlagRequired("backup-type")

	// RESTORE CMD
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&config.PublicAddress, "addr", "", "localhost", "Server address")
	restoreCmd.Flags().Uint32VarP(&config.Port, "port", "p", 8433, "Server port")
	restoreCmd.AddCommand(restoreStartCmd)
	restoreStartCmd.Flags().StringVarP(&manualJobParams.Client, "client", "", "", "Target client for restoration")
	restoreStartCmd.Flags().StringVarP(&manualJobParams.Module, "module", "m", "", "Module name")
	restoreStartCmd.Flags().StringVarP(&manualJobParams.RestoreJobUuid, "job", "j", "", "Job UUID to restore data from")
	restoreStartCmd.Flags().StringVarP(&manualJobParams.RestoreDestination, "destination", "d", "", "Alternate file restore destination")
	restoreStartCmd.MarkFlagRequired("client")
	restoreStartCmd.MarkFlagRequired("module")
	restoreStartCmd.MarkFlagRequired("job")

	// RETENTION CMD
	rootCmd.AddCommand(retentionCmd)
	retentionCmd.AddCommand(retentionCleanCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
