package api

import (
	"fmt"
	"log/slog"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/config"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job"
)

func ImageList(p api_helpers.PaginationParams, s api_helpers.ImageSearch) (api_helpers.PaginatedResponse[image.Image], error) {
	imgCount, err := image.Count()
	if err != nil {
		return api_helpers.PaginatedResponse[image.Image]{}, fmt.Errorf("cannot count total images: %w", err)
	}

	imgs, err := image.Search(p, s, config.Current.ModuleInstallPath)
	if err != nil {
		return api_helpers.PaginatedResponse[image.Image]{}, fmt.Errorf("cannot get images from database: %w", err)
	}

	return api_helpers.PaginatedResponse[image.Image]{
		Count:      imgCount,
		Pagination: p,
		Data:       imgs,
	}, nil
}

func ImageGet(uuid string) (image.Image, error) {
	img, err := image.GetByUuid(uuid)
	if err != nil {
		return image.Image{}, fmt.Errorf("cannot get image from db: %w", err)
	}

	return img, nil
}

func ImageExists(uuid string) (bool, error) {
	return image.Exists(uuid)
}

// ImageRemove deletes an image from database and its data from the repository.
// The job that generated the image is kept unless cleanJob is set.
func ImageRemove(uuid string, cleanJob bool) error {
	img, err := image.GetByUuid(uuid)
	if err != nil {
		return fmt.Errorf("cannot get image from db: %w", err)
	}

	storagePath, err := img.GetStorageFolderPath()
	if err != nil {
		img.GetLog().With(
			slog.Any("error", err),
		).Warn("Cannot determine image storage path, data in repository will not be removed")
		storagePath = ""
	}

	if err := img.Delete(); err != nil {
		return fmt.Errorf("cannot delete image: %w", err)
	}
	img.GetLog().Info("Image removed")

	if cleanJob {
		jobExists, err := job.Exists(uuid)
		if err != nil {
			return fmt.Errorf("cannot check if job linked to image exists: %w", err)
		}
		if jobExists {
			j := job.Job{Uuid: uuid}
			if err := j.Delete(); err != nil {
				return fmt.Errorf("cannot delete job linked to image: %w", err)
			}
			slog.With(slog.String("uuid", uuid)).Info("Job linked to image removed")
		}
	}

	if err := cleanupArtifacts(uuid, storagePath); err != nil {
		return fmt.Errorf("cannot remove image files from disk: %w", err)
	}

	return nil
}
