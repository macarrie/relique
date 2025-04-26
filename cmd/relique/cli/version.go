package cli

import (
	"fmt"

	"github.com/macarrie/relique/api"
	"github.com/spf13/cobra"
)

func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show relique version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(api.ConfigGetVersion())
		},
	}

	rootCmd.AddCommand(versionCmd)
}
