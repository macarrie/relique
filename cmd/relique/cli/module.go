package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/utils"
	"github.com/spf13/cobra"
)

var moduleListPageSize int

func init() {
	moduleCmd := &cobra.Command{
		Use:   "module",
		Short: "Module related commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := api.ConfigGet()
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get relique configuration")
				os.Exit(1)
			}
		},
	}

	moduleListCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed modules command",
		Run: func(cmd *cobra.Command, args []string) {
			page := api_helpers.PaginationParams{
				Limit:  uint64(moduleListPageSize),
				Offset: 0,
			}
			mods, err := api.ModuleList(page)
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get installed modules")
				os.Exit(1)
			}

			tab := tabular.New()
			tab.Col("name", "Name", 30)
			tab.Col("variant", "Variant", 30)
			tab.Col("available_variants", "Available variants", 30)
			tab.Col("backup_paths", "Backup paths", 30)

			format := tab.Print("name", "variant", "available_variants", "backup_paths")
			for _, m := range mods.Data {
				var variant string = m.Variant
				if variant == "" {
					variant = "default"
				}
				fmt.Printf(format,
					m.Name,
					variant,
					strings.Join(m.AvailableVariants, ", "),
					strings.Join(m.BackupPaths, ", "),
				)
			}

			fmt.Printf("\nShowing %d out of %d records\n", len(mods.Data), mods.Count)
		},
	}
	utils.AddPaginationParams(moduleListCmd, &moduleListPageSize)

	moduleShowCmd := &cobra.Command{
		Use:   "show MODULE_NAME",
		Short: "Show detailed information about installed module",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("TODO: Show module details")
		},
	}

	moduleInstallCmd := &cobra.Command{
		Use:   "install MODULE_URL_OR_PATH",
		Short: "Module install command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("TODO: Install module")
		},
	}

	moduleRemoveCmd := &cobra.Command{
		Use:   "remove MODULE_NAME",
		Short: "Module uninstall command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("TODO: Remove module")
		},
	}

	rootCmd.AddCommand(moduleCmd)
	moduleCmd.AddCommand(moduleListCmd)
	moduleCmd.AddCommand(moduleShowCmd)
	moduleCmd.AddCommand(moduleInstallCmd)
	moduleCmd.AddCommand(moduleRemoveCmd)
}
