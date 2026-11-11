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
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/utils"
	"github.com/spf13/cobra"
)

var imageListPageSize int
var imageListSearchModule string
var imageListSearchClient string

var imageRemoveOlderThan string
var imageRemoveSearchModule string
var imageRemoveSearchClient string
var imageRemoveCleanJob bool

// buildImageRemovalPlan selects the images to remove and the jobs linked to them, without modifying anything
func buildImageRemovalPlan(args []string) (removalPlan, error) {
	var selected []image.Image
	if len(args) == 1 {
		img, err := api.ImageGet(args[0])
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot get image details: %w", err)
		}
		selected = append(selected, img)
	} else {
		before, cutoff, err := olderThanCutoff(imageRemoveOlderThan)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot parse --older-than: %w", err)
		}
		search := api_helpers.ImageSearch{
			ModuleName: imageRemoveSearchModule,
			ClientName: imageRemoveSearchClient,
			Before:     before,
		}
		fmt.Printf("Selecting images created before %s\n", utils.FormatDatetime(cutoff))

		imageList, err := api.ImageList(api_helpers.PaginationParams{Limit: 0}, search)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot get image list: %w", err)
		}
		selected = imageList.Data
	}

	plan := removalPlan{ImagesToRemove: selected}
	for _, img := range selected {
		exists, err := api.JobExists(img.Uuid)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot check if image '%s' has a linked job: %w", img.Uuid, err)
		}
		if !exists {
			continue
		}
		j, err := api.JobGet(img.Uuid)
		if err != nil {
			return removalPlan{}, fmt.Errorf("cannot get job linked to image '%s': %w", img.Uuid, err)
		}
		if imageRemoveCleanJob {
			plan.JobsToRemove = append(plan.JobsToRemove, j)
		} else {
			plan.JobsKept = append(plan.JobsKept, j)
		}
	}

	return plan, nil
}

func init() {
	imageCmd := &cobra.Command{
		Use:     "image",
		Aliases: []string{"img"},
		Short:   "Backup image related commands",
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
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List images generated from backups",
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

	imageRemoveCmd := &cobra.Command{
		Use:     "remove [UUID]",
		Aliases: []string{"rm"},
		Short:   "Remove images and their data, by UUID or by age",
		Long: `Remove an image by UUID, or every image older than a given age with --older-than.
Image data stored in the repository is deleted.

The job that generated the image is kept for history unless --clean-job is used.
A recap of everything that will be removed is displayed and must be confirmed before anything is deleted.`,
		Example: `  relique image remove 1b4e28ba-2fa1-11d2-883f-0016d3cca427
  relique image remove 1b4e28ba-2fa1-11d2-883f-0016d3cca427 --clean-job
  relique image remove --older-than 90d --module my-module --yes`,
		Args: validateRemoveArgs("client", "module"),
		Run: func(cmd *cobra.Command, args []string) {
			plan, err := buildImageRemovalPlan(args)
			if err != nil {
				slog.With(
					slog.Any("error", err),
				).Error("Cannot select images to remove")
				os.Exit(1)
			}

			if !confirmRemoval(plan) {
				return
			}

			uuids := make([]string, 0, len(plan.ImagesToRemove))
			for _, img := range plan.ImagesToRemove {
				uuids = append(uuids, img.Uuid)
			}
			runRemoval("image", uuids, func(uuid string) error {
				return api.ImageRemove(uuid, imageRemoveCleanJob)
			})
		},
	}
	imageRemoveCmd.Flags().StringVarP(&imageRemoveOlderThan, "older-than", "", "", "Remove images created before this age (e.g. 90d, 2w, 12h)")
	imageRemoveCmd.Flags().StringVarP(&imageRemoveSearchClient, "client", "", "", "With --older-than: filter on client name")
	imageRemoveCmd.Flags().StringVarP(&imageRemoveSearchModule, "module", "m", "", "With --older-than: filter on module name")
	imageRemoveCmd.Flags().BoolVarP(&imageRemoveCleanJob, "clean-job", "", false, "Also remove the job that generated the image")
	imageRemoveCmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "Skip confirmation (the recap is still displayed)")

	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imageShowCmd)
	imageCmd.AddCommand(imageRemoveCmd)
}
