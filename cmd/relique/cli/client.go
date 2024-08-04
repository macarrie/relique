package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func init() {
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "Backup client related commands",
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

	clientListCmd := &cobra.Command{
		Use:   "list",
		Short: "List configured backup clients",
		Run: func(cmd *cobra.Command, args []string) {
			clientList := api.ClientList()

			tab := tabular.New()
			tab.Col("name", "Name", 40)
			tab.Col("addr", "Address", 40)
			tab.Col("modules", "Modules", 40)

			format := tab.Print("name", "addr", "modules")
			for _, c := range clientList {
				var moduleNames []string
				for _, mod := range c.Modules {
					moduleNames = append(moduleNames, mod.String())
				}
				fmt.Printf(format, c.Name, c.Address, strings.Join(moduleNames, ", "))
			}
		},
	}

	clientShowCmd := &cobra.Command{
		Use:   "show CLIENT_NAME",
		Short: "Show backup client details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cl, err := api.ClientGet(args[0])
			if err != nil {
				slog.With(
					slog.String("client", args[0]),
					slog.Any("error", err),
				).Error("Cannot get client details")
				os.Exit(1)
			}

			out, err := toml.Marshal(cl)
			if err != nil {
				slog.Error("Cannot display client details", slog.Any("error", err))
				os.Exit(1)
			}
			fmt.Printf("\n%v\n", string(out))
		},
	}

	clientPingCmd := &cobra.Command{
		Use:   "ping CLIENT_NAME",
		Short: "Ping client via SSH",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cl, pingErr := api.ClientGet(args[0])
			if pingErr != nil {
				slog.With(
					slog.String("client", args[0]),
					slog.Any("error", pingErr),
				).Error("Cannot get client details")
				os.Exit(1)
			}

			if err := api.ClientSSHPing(cl); err != nil {
				slog.Error("Cannot ping client", slog.Any("error", err))
				os.Exit(1)
			}
		},
	}

	clientCreateCmd := &cobra.Command{
		Use:   "create CLIENT_NAME",
		Short: "Create a new backup client",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: Create backup client")
		},
	}

	clientModifyCmd := &cobra.Command{
		Use:   "modify CLIENT_NAME",
		Short: "Modify an existing backup client",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: Modify backup client")
		},
	}

	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientListCmd)
	clientCmd.AddCommand(clientShowCmd)
	clientCmd.AddCommand(clientPingCmd)
	clientCmd.AddCommand(clientCreateCmd)
	clientCmd.AddCommand(clientModifyCmd)
}
