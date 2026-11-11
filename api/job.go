package api

import (
	"fmt"
	"log/slog"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job"
)

func JobList(p api_helpers.PaginationParams, s api_helpers.JobSearch) (api_helpers.PaginatedResponse[job.Job], error) {
	jobCount, err := job.Count(s)
	if err != nil {
		return api_helpers.PaginatedResponse[job.Job]{}, fmt.Errorf("cannot count total jobs: %w", err)
	}

	jobs, err := job.Search(p, s)
	if err != nil {
		return api_helpers.PaginatedResponse[job.Job]{}, fmt.Errorf("cannot get jobs from database: %w", err)
	}

	return api_helpers.PaginatedResponse[job.Job]{
		Count:      jobCount,
		Pagination: p,
		Data:       jobs,
	}, nil
}

func JobGet(uuid string) (job.Job, error) {
	j, err := job.GetByUuid(uuid)
	if err != nil {
		return job.Job{}, fmt.Errorf("cannot get job from db: %w", err)
	}

	return j, nil
}

func JobExists(uuid string) (bool, error) {
	return job.Exists(uuid)
}

// JobRemove deletes a finished job from database and the files it owns on disk.
// The image generated from a backup job is kept unless cleanImage is set.
func JobRemove(uuid string, cleanImage bool) error {
	j, err := job.GetByUuid(uuid)
	if err != nil {
		return fmt.Errorf("cannot get job from db: %w", err)
	}

	if !j.Done {
		return fmt.Errorf("cannot remove job '%s' with status '%s' because it is not finished", j.Uuid, j.Status.String())
	}

	storagePath, err := j.GetStorageFolderPath()
	if err != nil {
		j.GetLog().With(
			slog.Any("error", err),
		).Warn("Cannot determine job storage path, data in repository will not be removed")
		storagePath = ""
	}

	if err := j.Delete(); err != nil {
		return fmt.Errorf("cannot delete job: %w", err)
	}
	j.GetLog().Info("Job removed")

	if cleanImage {
		imgExists, err := image.Exists(uuid)
		if err != nil {
			return fmt.Errorf("cannot check if image linked to job exists: %w", err)
		}
		if imgExists {
			img := image.Image{Uuid: uuid}
			if err := img.Delete(); err != nil {
				return fmt.Errorf("cannot delete image linked to job: %w", err)
			}
			slog.With(slog.String("uuid", uuid)).Info("Image linked to job removed")
		}
	}

	if err := cleanupArtifacts(uuid, storagePath); err != nil {
		return fmt.Errorf("cannot remove job files from disk: %w", err)
	}

	return nil
}
