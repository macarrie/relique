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

var clientListPageSize int
var clientListSearchModule string
var clientListSearchModuleType string

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
			page := api_helpers.PaginationParams{
				Limit:  uint64(clientListPageSize),
				Offset: 0,
			}
			search := api_helpers.ClientSearch{
				ModuleName: clientListSearchModule,
				ModuleType: clientListSearchModuleType,
			}
			clientList := api.ClientList(page, search)

			tab := tabular.New()
			tab.Col("name", "Name", 40)
			tab.Col("addr", "Address", 40)
			tab.Col("modules", "Modules", 40)

			format := tab.Print("name", "addr", "modules")
			for _, c := range clientList.Data {
				var moduleNames []string
				for _, mod := range c.Modules {
					moduleNames = append(moduleNames, mod.String())
				}
				fmt.Printf(format, c.Name, c.Address, strings.Join(moduleNames, ", "))
			}

			fmt.Printf("\nShowing %d out of %d records\n", len(clientList.Data), clientList.Count)
		},
	}
	utils.AddPaginationParams(clientListCmd, &clientListPageSize)
	clientListCmd.Flags().StringVarP(&clientListSearchModule, "module", "m", "", "Filter on module name")
	clientListCmd.Flags().StringVarP(&clientListSearchModuleType, "module-type", "", "", "Filter on module type")

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

			clientDetailsTemplate := `# Client details
-----
## Global
			
Name: 	{{.Name}}

Address: 	{{.Address}}

-----
## SSH connexion

User: 	{{.SSHUser}}

Port: 	{{.SSHPort}}


{{ range .Modules}}
-----
## Module __{{.Name}}__

| Parameter | Value |
| --------- | ----- |
| Name | {{ .Name }} |
| Module type | {{ .ModuleType }} |
| Backup type | {{ .BackupType }} |
| Variant | {{ if .Variant | eq "" }} default {{ else }}{{ .Variant }}{{ end }} |
| Available variants | {{ if .AvailableVariants | len | eq 0 }} default {{ else }}{{ .Variant }}{{ end }}{{ join .AvailableVariants ", " }} |
| Backup paths | {{ join .BackupPaths ", " }} |

{{ end }}
`

			render, err := utils.RenderTemplateToMarkdown("client_details", clientDetailsTemplate, cl)
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot display client info")
				os.Exit(1)
			}
			fmt.Print(render)
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
