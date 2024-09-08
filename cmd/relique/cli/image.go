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
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

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
			// TODO: Handle pagination
			imageList, err := api.ImageList(api_helpers.PaginationParams{})
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

			// TODO: Handle pagination
			fmt.Printf("\nShowing %d out of %d records\n", imageList.Count, imageList.Count)
		},
	}

	imageShowCmd := &cobra.Command{
		Use:   "show UUID",
		Short: "Show image details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cl, err := api.ImageGet(args[0])
			if err != nil {
				slog.With(
					slog.String("image", args[0]),
					slog.Any("error", err),
				).Error("Cannot get image details")
				os.Exit(1)
			}

			out, err := toml.Marshal(cl)
			if err != nil {
				slog.Error("Cannot display image details", slog.Any("error", err))
				os.Exit(1)
			}
			fmt.Printf("\n%v\n", string(out))
		},
	}

	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imageShowCmd)
}
