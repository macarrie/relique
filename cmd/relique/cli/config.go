package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/macarrie/relique/api"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var configInitCfgPath string
var configInitModPath string
var configInitStoragePath string

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Show relique configuration",
	}

	configInitCmd := &cobra.Command{
		Use:   "init CFG_PATH",
		Short: "Initialize default relique configuration in CFG_PATH/relique folder",
		Run: func(cmd *cobra.Command, args []string) {
			if err := api.ConfigInit(configInitCfgPath, configInitModPath, configInitStoragePath); err != nil {
				slog.With(
					slog.String("cfg_root", configInitCfgPath),
					slog.String("module_root", configInitModPath),
					slog.Any("error", err),
				).Error("cannot initialize default relique config")
				os.Exit(1)
			}
		},
	}
	configInitCmd.Flags().StringVarP(&configInitCfgPath, "path", "", "", "Configuration folder path (default: '/etc/relique/')")
	configInitCmd.Flags().StringVarP(&configInitModPath, "module-install-path", "", "", "Module install path (default: '/var/lib/relique/modules')")
	configInitCmd.Flags().StringVarP(&configInitModPath, "storage-path", "", "", "Default repository storage path (default: '/var/lib/relique/storage')")

	configShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current relique config",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := api.ConfigGet()
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get relique configuration")
				os.Exit(1)
			}

			out, err := toml.Marshal(cfg)
			if err != nil {
				slog.Error("Cannot display current configuration", slog.Any("error", err))
				os.Exit(1)
			}
			fmt.Printf("\n%v\n", string(out))
		},
	}

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
}
