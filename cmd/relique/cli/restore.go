package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/utils"
	"github.com/spf13/cobra"
)

var imageId string
var restoreClient string
var assumeYes bool

func printRestoreRecap(img image.Image, c client.Client, customRestorePaths []string) {
	restorePaths := utils.GenerateCustomRestorePaths(customRestorePaths, img.Module.BackupPaths)
	data := struct {
		Image        image.Image
		Client       client.Client
		RestorePaths map[string]string
	}{
		Image:        img,
		Client:       c,
		RestorePaths: restorePaths,
	}
	restoreRecapTemplate := `# Data restore recap

## Image details

Source image: {{ .Image.Uuid }}

This image has been generated from a backup with the following parameters:

* Client '{{ .Image.Client.Name }}' ({{ .Image.Client.Address }})
* Module '{{ .Image.Module.Name }}'
	* Module type: {{ .Image.Module.ModuleType }}
* Backup paths
{{ range .Image.Module.BackupPaths }}
	* {{ . }}
{{ end }}

## Restore details

Data will be restored on client '{{ .Client.Name }}' ({{ .Client.Address }})
{{ if gt (len .RestorePaths) 0 }}
The following custom restore paths will be used

| Source | Destination on '{{ .Client.Name }}'|
| ------ | ----------- |
{{- range $source, $dest := .RestorePaths }}
| {{ $source }} | {{ $dest }} |
{{- end }}
{{ else }}
Original backup paths will be restored on client.

| Source | Destination on '{{ .Client.Name }}' |
| ------ | ----------- |
{{- range .Image.Module.BackupPaths }}
| {{ . }} | {{ . }} |
{{- end }}
{{ end }}
`

	render, err := utils.RenderTemplateToMarkdown("restore_recap", restoreRecapTemplate, data)
	if err != nil {
		slog.With(
			slog.Any("error", err),
		).Error("Cannot generate restore recap")
		os.Exit(1)
	}
	fmt.Print(render)
}

func init() {
	restoreCmd := &cobra.Command{
		Use:   "restore",
		Args:  cobra.ArbitraryArgs,
		Short: "Restore related commands",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := api.ConfigGet()
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get relique configuration")
				os.Exit(1)
			}

			if err := db.Init(config.GetDBPath()); err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot initialize database connection")
				os.Exit(1)
			}

			img, err := image.GetByUuid(imageId)
			if err != nil {
				slog.With(
					slog.Any("error", err),
					slog.Any("image_uuid", imageId),
				).Error("Cannot find image")
				os.Exit(1)
			}

			c, err := api.ClientGet(restoreClient)
			if err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", restoreClient),
				).Error("Cannot find client")
				os.Exit(1)
			}

			printRestoreRecap(img, c, args)
			if assumeYes {
				slog.Info("Skipping confirmation on user request (-y/--yes flag provided)")
			} else {
				if !utils.Confirm("Continue restore") {
					slog.Error("Restore process canceled")
					os.Exit(1)
				}
			}

			if err := api.RestoreStart(c, img, args); err != nil {
				slog.With(
					slog.Any("error", err),
					slog.String("client", c.Name),
					slog.String("image", img.Uuid),
					slog.String("restore_paths", strings.Join(args, ", ")),
				).Error("Error during restore job")
				os.Exit(1)
			}
		},
	}
	restoreCmd.Flags().StringVarP(&imageId, "image", "", "", "Reference image to restore")
	restoreCmd.Flags().StringVarP(&restoreClient, "to", "", "", "Target client for image restore")
	restoreCmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "Skip confirmation on restore")
	restoreCmd.MarkFlagRequired("image")
	restoreCmd.MarkFlagRequired("to")

	rootCmd.AddCommand(restoreCmd)
}
