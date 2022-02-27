package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/macarrie/relique/internal/types/displayable"
	"github.com/macarrie/relique/internal/types/module"
	"github.com/spf13/cobra"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/macarrie/relique/internal/types/config/common"
	"github.com/macarrie/relique/internal/types/relique_job"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

type Args struct {
	Debug      bool
	ConfigPath string
	JSON       bool
}

var Params = Args{}

var ModuleInstallPath string
var ModuleInstallIsArchive bool
var ModuleInstallIsLocal bool
var ModuleInstallForce bool
var ModuleInstallSkipChown bool

func InitCommonParams() {
	if Params.JSON {
		displayable.DisplayMode = displayable.JSON
	} else {
		displayable.DisplayMode = displayable.TUI
	}
	log.SetupCliLogger(Params.Debug, Params.JSON)
}

func GetCommonCliCommands(rootCmd *cobra.Command) {
	// ROOT CMD
	rootCmd.PersistentFlags().BoolVar(&Params.JSON, "json", false, "Output content as JSON")
	rootCmd.PersistentFlags().BoolVarP(&Params.Debug, "verbose", "v", false, "verbose log output")
	rootCmd.PersistentFlags().StringVarP(&Params.ConfigPath, "config", "c", "", "Configuration file path")

	moduleCmd := &cobra.Command{
		Use:   "module",
		Short: "Module related commands",
	}
	moduleListCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed modules command",
		Run: func(cmd *cobra.Command, args []string) {
			if ModuleInstallPath != "" {
				module.MODULES_INSTALL_PATH = ModuleInstallPath
				module.ModulesInstallPathReadInConfig = true
			}
			installedModules, err := module.GetLocallyInstalled()
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot list installed modules")
				os.Exit(1)
			}

			disp := make([]displayable.Displayable, len(installedModules))
			for i, v := range installedModules {
				disp[i] = v
			}
			displayable.Table(disp)
		},
	}
	moduleInstallCmd := &cobra.Command{
		Use:   "install",
		Short: "Module install command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			moduleSource := args[0]
			if ModuleInstallPath != "" {
				module.MODULES_INSTALL_PATH = ModuleInstallPath
				module.ModulesInstallPathReadInConfig = true
			}
			err := module.Install(moduleSource, ModuleInstallIsLocal, ModuleInstallIsArchive, ModuleInstallForce, ModuleInstallSkipChown)
			if err != nil {
				log.WithFields(log.Fields{
					"err":    err,
					"module": moduleSource,
				}).Error("Cannot install relique module")
				os.Exit(1)
			}
		},
	}
	moduleRemoveCmd := &cobra.Command{
		Use:   "remove",
		Short: "Module uninstall command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			moduleName := args[0]
			if ModuleInstallPath != "" {
				module.MODULES_INSTALL_PATH = ModuleInstallPath
				module.ModulesInstallPathReadInConfig = true
			}
			err := module.Remove(moduleName)
			if err != nil {
				log.WithFields(log.Fields{
					"err":    err,
					"module": moduleName,
				}).Error("Cannot remove relique module")
				os.Exit(1)
			}
		},
	}

	// MODULE CMD
	rootCmd.AddCommand(moduleCmd)
	moduleCmd.PersistentFlags().StringVarP(&ModuleInstallPath, "install-path", "p", "", "Module install path")
	moduleCmd.AddCommand(moduleListCmd)
	moduleCmd.AddCommand(moduleInstallCmd)
	moduleCmd.AddCommand(moduleRemoveCmd)
	moduleInstallCmd.Flags().BoolVarP(&ModuleInstallIsArchive, "archive", "a", false, "Module to install is packaged into a tar.gz archive instead of being a git repository")
	moduleInstallCmd.Flags().BoolVarP(&ModuleInstallIsLocal, "local", "l", false, "Module to install is already available locally on disk (offline install)")
	moduleInstallCmd.Flags().BoolVarP(&ModuleInstallForce, "force", "f", false, "Force module install. If module is already installed, files with be overwritten")
	moduleInstallCmd.Flags().BoolVarP(&ModuleInstallSkipChown, "skip-chown", "", false, "Do not chown module files to relique user and group after install")
}

func ManualJobStart(config common.Configuration, params relique_job.JobSearchParams) (relique_job.ReliqueJob, error) {
	var job relique_job.ReliqueJob

	response, err := utils.PerformRequest(config,
		config.PublicAddress,
		config.Port,
		"POST",
		"/api/v1/jobs/start",
		params)
	if err != nil {
		return relique_job.ReliqueJob{}, errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return relique_job.ReliqueJob{}, errors.Wrap(err, "cannot read response body from api request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return relique_job.ReliqueJob{}, fmt.Errorf("cannot start job (%d response): '%s'", response.StatusCode, body)
	}

	if err := json.Unmarshal(body, &job); err != nil {
		return relique_job.ReliqueJob{}, errors.Wrap(err, "cannot parse started job returned from client")
	}

	return job, nil
}
