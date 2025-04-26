package image

import (
	"log/slog"
	"time"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
)

type Image struct {
	// Database IDs
	ID int64

	Uuid             string          `json:"uuid"`
	CreatedAt        time.Time       `json:"created_at"`
	Client           client.Client   `json:"client"`
	Module           module.Module   `json:"module"`
	Repository       repo.Repository `json:"repository"`
	NumberOfElements int             `json:"number_of_elements"`
	NumberOfFiles    int             `json:"number_of_files"`
	NumberOfFolders  int             `json:"number_of_folders"`
	SizeOnDisk       uint64          `json:"size_on_disk"`

	ClientName string
	ModuleName string
	RepoName   string
}

func (img *Image) GetLog() *slog.Logger {
	return slog.With(
		slog.String("uuid", img.Uuid),
		slog.String("client", img.Client.String()),
		slog.String("module", img.Module.String()),
		slog.String("repository", img.Repository.GetName()),
	)
}
