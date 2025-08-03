package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/InVisionApp/tabular"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/utils"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var clientListPageSize int
var clientListSearchModule string
var clientListSearchModuleType string

var clientCreateAddress string
var clientCreateSSHUser string
var clientCreateSSHPort int

var clientModuleType string
var clientModuleBackupType string
var clientModuleVariant string
var clientModuleBackupPaths []string
var clientModuleInclude []string
var clientModuleExclude []string
var clientModuleExcludeCVS bool

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
| Inclusions | {{ join .Include ", " }} |
| Exclusions | {{ join .Exclude ", " }} |
| Exclude CVS | {{ .ExcludeCVS }} |

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
			clientName := args[0]
			var clientAddress string

			if clientCreateAddress == "" {
				clientAddress = clientName
			} else {
				clientAddress = clientCreateAddress
			}

			if err := api.ClientCreate(clientName, clientAddress, clientCreateSSHUser, clientCreateSSHPort); err != nil {
				slog.With(
					slog.String("name", clientName),
					slog.String("address", clientAddress),
					slog.String("ssh_user", clientCreateSSHUser),
					slog.Int("ssh_port", clientCreateSSHPort),
					slog.Any("error", err),
				).Error("Cannot create client")
				os.Exit(1)
			}
			slog.With(
				slog.String("name", clientName),
				slog.String("address", clientAddress),
				slog.String("ssh_user", clientCreateSSHUser),
				slog.Int("ssh_port", clientCreateSSHPort),
			).Info("Created client configuration file")
		},
	}
	clientCreateCmd.Flags().StringVarP(&clientCreateAddress, "address", "", "", "Client address (using client name if empty)")
	clientCreateCmd.Flags().StringVarP(&clientCreateSSHUser, "ssh-user", "u", "", "Client SSH port")
	clientCreateCmd.Flags().IntVarP(&clientCreateSSHPort, "ssh-port", "p", 0, "Client SSH port")

	clientModuleCmd := &cobra.Command{
		Use:   "module",
		Short: "Handle client modules",
	}
	clientModuleAddCmd := &cobra.Command{
		Use:   "add CLIENT_NAME MODULE_NAME",
		Short: "Add module to client",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			clientName := args[0]
			moduleName := args[1]

			backupType := backup_type.FromString(clientModuleBackupType)
			if backupType.Type == backup_type.Unknown {
				slog.With(
					slog.String("value_from_cli", clientModuleBackupType),
				).Error("Invalid value for backup type")
				os.Exit(1)
			}

			c, err := api.ClientGet(clientName)
			if err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", clientName),
				).Error("Cannot find client")
				os.Exit(1)
			}

			mod := module.Module{
				Name:        moduleName,
				ModuleType:  clientModuleType,
				BackupType:  backupType,
				Variant:     clientModuleVariant,
				BackupPaths: clientModuleBackupPaths,
				Include:     clientModuleInclude,
				Exclude:     clientModuleExclude,
				ExcludeCVS:  clientModuleExcludeCVS,
			}

			// TODO: Check if module already exists
			_, alreadyExists := lo.Find(c.Modules, func(m module.Module) bool {
				return m.Name == moduleName
			})
			if alreadyExists {
				slog.With(
					slog.String("module_name", moduleName),
				).Error("Module already exists on client")
				os.Exit(1)
			}

			c.Modules = append(c.Modules, mod)
			if err := api.ClientSave(c); err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", clientName),
				).Error("Cannot save client configuration to file")
				os.Exit(1)
			}
			c.GetLog().Info("Saved new module into client configuration file")
		},
	}
	clientModuleAddCmd.Flags().StringVarP(&clientModuleType, "type", "t", "generic", "Module type")
	clientModuleAddCmd.Flags().StringVarP(&clientModuleBackupType, "backup-type", "", "diff", "Backup type")
	clientModuleAddCmd.Flags().StringVarP(&clientModuleVariant, "variant", "", "default", "Backup module variant")
	clientModuleAddCmd.Flags().StringSliceVarP(&clientModuleBackupPaths, "path", "p", []string{}, "Backup path")
	clientModuleAddCmd.Flags().StringSliceVarP(&clientModuleExclude, "exclude", "e", []string{}, "File exclusions")
	clientModuleAddCmd.Flags().StringSliceVarP(&clientModuleInclude, "include", "i", []string{}, "File inclusions")
	clientModuleAddCmd.Flags().BoolVarP(&clientModuleExcludeCVS, "exclude-cvs", "", false, "Exclude CVS from file selections")

	clientModuleRmCmd := &cobra.Command{
		Use:   "rm CLIENT_NAME MODULE_NAME",
		Short: "Remove module from client",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			clientName := args[0]
			moduleName := args[1]

			c, err := api.ClientGet(clientName)
			if err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", clientName),
				).Error("Cannot find client")
				os.Exit(1)
			}

			_, alreadyExists := lo.Find(c.Modules, func(m module.Module) bool {
				return m.Name == moduleName
			})
			if !alreadyExists {
				slog.With(
					slog.String("module_name", moduleName),
				).Error("Module not found on client")
				os.Exit(1)
			}

			if assumeYes {
				slog.Info("Skipping confirmation on user request (-y/--yes flag provided)")
			} else {
				if !utils.Confirm("Confirm module removal from client") {
					slog.Error("Module removal canceled")
					os.Exit(1)
				}
			}

			c.Modules = lo.Filter(c.Modules, func(m module.Module, _ int) bool {
				return m.Name != moduleName
			})
			if err := api.ClientSave(c); err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", clientName),
				).Error("Cannot save client configuration to file")
				os.Exit(1)
			}
			c.GetLog().With(
				slog.String("module_name", moduleName),
			).Info("Removed module from client configuration file")
		},
	}
	clientModuleRmCmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "Skip confirmation on module delete")

	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientListCmd)
	clientCmd.AddCommand(clientShowCmd)
	clientCmd.AddCommand(clientPingCmd)
	clientCmd.AddCommand(clientCreateCmd)
	clientCmd.AddCommand(clientModuleCmd)
	clientModuleCmd.AddCommand(clientModuleAddCmd)
	clientModuleCmd.AddCommand(clientModuleRmCmd)
}
