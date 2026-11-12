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

// JobRetry starts a new job with the same parameters as the finished job identified by uuid and waits for its end.
// The new job is returned even if its run fails, so callers can reference it
func JobRetry(uuid string) (job.Job, error) {
	retry, err := JobRetrySetup(uuid)
	if err != nil {
		return retry, err
	}

	if err := JobRetryRun(&retry); err != nil {
		return retry, err
	}

	return retry, nil
}

// JobRetrySetup prepares a retry of the finished job identified by uuid.
// The new job is set up and saved in database but not started: use JobRetryRun to run it
func JobRetrySetup(uuid string) (job.Job, error) {
	original, err := job.GetByUuid(uuid)
	if err != nil {
		return job.Job{}, fmt.Errorf("cannot get job from db: %w", err)
	}

	retry, err := original.NewRetry()
	if err != nil {
		return job.Job{}, fmt.Errorf("cannot prepare job retry: %w", err)
	}

	retry.GetLog().With(
		slog.String("original_uuid", original.Uuid),
	).Info("Retrying job")
	if err := retry.SetupBackup(); err != nil {
		return retry, fmt.Errorf("cannot setup retried job: %w", err)
	}

	return retry, nil
}

// JobRetryRun runs a job prepared by JobRetrySetup
func JobRetryRun(retry *job.Job) error {
	if err := startBackup(retry); err != nil {
		return fmt.Errorf("error during retried job: %w", err)
	}

	return nil
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
