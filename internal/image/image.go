package image

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	rsync_lib "github.com/macarrie/relique/internal/rsync_task/lib"
	"github.com/macarrie/relique/internal/utils"
)

func New(c client.Client, m module.Module, r repo.Repository) Image {
	return Image{
		Uuid:       uuid.New().String(),
		Client:     c,
		Module:     m,
		Repository: r,
		CreatedAt:  time.Now(),
	}
}

func GetSizeOnDisk(rootPath string) (uint64, error) {
	var totalSize uint64

	// TODO: Handle error
	filepath.Walk(rootPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += uint64(info.Size())
		}
		return err
	})

	return totalSize, nil
}

func (img *Image) GetStorageFolderPath() (string, error) {
	return utils.GetStoragePath(img.Repository, img.RepoName, img.Uuid)
}

func (img *Image) FillStats(stats rsync_lib.Stats, rootPath string) error {
	sizeOnDisk, err := GetSizeOnDisk(rootPath)
	if err != nil {
		return fmt.Errorf("cannot get image size on disk: %w", err)
	}

	img.SizeOnDisk = sizeOnDisk
	img.NumberOfElements = stats.NumberOfFiles
	img.NumberOfFiles = stats.NumberOfRegularFiles
	img.NumberOfFolders = stats.NumberOfDirectories

	return nil
}
