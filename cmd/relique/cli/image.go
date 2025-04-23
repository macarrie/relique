package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/InVisionApp/tabular"
	"github.com/dustin/go-humanize"
	"github.com/macarrie/relique/api"
	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/utils"
	"github.com/spf13/cobra"
)

var imageListPageSize int
var imageListSearchModule string
var imageListSearchClient string

func init() {
	imageCmd := &cobra.Command{
		Use:   "image",
		Short: "Backup image related commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
		},
	}

	imageListCmd := &cobra.Command{
		Use:   "list",
		Short: "List images generated from backups",
		Run: func(cmd *cobra.Command, args []string) {
			page := api_helpers.PaginationParams{
				Limit:  uint64(imageListPageSize),
				Offset: 0,
			}
			search := api_helpers.ImageSearch{
				ModuleName: imageListSearchModule,
				ClientName: imageListSearchClient,
			}
			imageList, err := api.ImageList(page, search)
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot get image list")
				os.Exit(1)
			}

			tab := tabular.New()
			tab.Col("uuid", "UUID", 40)
			tab.Col("client", "Client", 25)
			tab.Col("module", "Module", 15)
			tab.Col("date", "Date", 20)
			tab.Col("size", "Size", 20)

			format := tab.Print("uuid", "client", "module", "date", "size")
			for _, img := range imageList.Data {
				fmt.Printf(
					format,
					img.Uuid,
					img.Client.String(),
					img.Module.String(),
					utils.FormatDatetime(img.CreatedAt),
					humanize.Bytes(img.SizeOnDisk),
				)
			}

			fmt.Printf("\nShowing %d out of %d records\n", len(imageList.Data), imageList.Count)
		},
	}
	utils.AddPaginationParams(imageListCmd, &imageListPageSize)
	imageListCmd.Flags().StringVarP(&imageListSearchClient, "client", "", "", "Filter on client name")
	imageListCmd.Flags().StringVarP(&imageListSearchModule, "module", "m", "", "Filter on module name")

	imageShowCmd := &cobra.Command{
		Use:   "show UUID",
		Short: "Show image details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			img, err := api.ImageGet(args[0])
			if err != nil {
				slog.With(
					slog.String("image", args[0]),
					slog.Any("error", err),
				).Error("Cannot get image details")
				os.Exit(1)
			}

			imageDetailsTemplate := `# Image details
-----
## Global
			
Created on {{ datetime .CreatedAt }}

UUID: {{ .Uuid }}

This image has been generated from the job with the same UUID

Client: 		{{ .Client.Name }}

Repository: 		{{ .Repository.Name }}

Repository type: 	{{ .Repository.Type }}


## Module

| Parameter | Value |
| --------- | ----- |
| Name | {{ .Module.Name }} |
| Module type | {{ .Module.ModuleType }} |
| Backup paths | {{ join .Module.BackupPaths ", " }} |

## Statistics

Size on disk: {{ file_size .SizeOnDisk}}

Number of elements: {{ .NumberOfElements}}

Files: {{ .NumberOfFiles}}

Directories: {{ .NumberOfFolders}}
`

			render, err := utils.RenderTemplateToMarkdown("image_details", imageDetailsTemplate, img)
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot display image details")
				os.Exit(1)
			}
			fmt.Print(render)
		},
	}

	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imageShowCmd)
}
