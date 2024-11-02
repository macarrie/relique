package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var Version string

func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show relique version",
		Run: func(cmd *cobra.Command, args []string) {
			if Version == "" {
				slog.Error("Empty version flag")
				os.Exit(1)
			}

			fmt.Println(Version)
		},
	}

	rootCmd.AddCommand(versionCmd)
}
