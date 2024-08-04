package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var repoCreateName string
var repoCreateIsDefault bool
var repoCreateLocalPath string

func init() {
	repoCmd := &cobra.Command{
		Use:   "repo",
		Short: "Backup repository related commands",
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

	repoListCmd := &cobra.Command{
		Use:   "list",
		Short: "List configured backup repositories",
		Run: func(cmd *cobra.Command, args []string) {
			repoList := api.RepoList()

			tab := tabular.New()
			tab.Col("name", "Name", 40)
			tab.Col("type", "Type", 40)
			tab.Col("default", "Default", 40)

			format := tab.Print("name", "type", "default")
			for _, r := range repoList {
				fmt.Printf(format, r.GetName(), r.GetType(), r.IsDefault())
			}
		},
	}

	repoShowCmd := &cobra.Command{
		Use:   "show REPO_NAME",
		Short: "Show backup repository details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			r, err := api.RepoGet(args[0])
			if err != nil {
				slog.With(
					slog.String("repository", args[0]),
					slog.Any("error", err),
				).Error("Cannot get repository details")
				os.Exit(1)
			}

			out, err := toml.Marshal(r)
			if err != nil {
				slog.Error("Cannot display repository details", slog.Any("error", err))
				os.Exit(1)
			}
			fmt.Printf("\n%v\n", string(out))
		},
	}

	repoCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new backup repository",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("TODO: Create backup repo")
		},
	}
	repoCreateCmd.PersistentFlags().StringVarP(&repoCreateName, "name", "n", "", "Repository name")
	repoCreateCmd.PersistentFlags().BoolVarP(&repoCreateIsDefault, "default", "", false, "Set repository as default")
	repoCreateCmd.MarkFlagRequired("name")

	repoCreateLocalCmd := &cobra.Command{
		Use:   "local",
		Short: "Create a new local backup repository",
		Run: func(cmd *cobra.Command, args []string) {
			if err := api.RepoCreateLocal(repoCreateName, repoCreateLocalPath, repoCreateIsDefault); err != nil {
				slog.With(
					slog.String("name", repoCreateName),
					slog.String("path", repoCreateLocalPath),
					slog.Bool("default", repoCreateIsDefault),
					slog.Any("error", err),
				).Error("cannot create local repository")
				os.Exit(1)
			}

			slog.With(
				slog.String("name", repoCreateName),
				slog.String("path", repoCreateLocalPath),
				slog.Bool("default", repoCreateIsDefault),
			).Info("Successfully created local repository")
		},
	}
	repoCreateCmd.AddCommand(repoCreateLocalCmd)
	repoCreateLocalCmd.Flags().StringVarP(&repoCreateLocalPath, "path", "p", "", "Local repository data storage path")
	repoCreateLocalCmd.MarkFlagRequired("path")

	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoShowCmd)
	repoCmd.AddCommand(repoCreateCmd)
}
