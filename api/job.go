package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/job"
)

func JobList(p api_helpers.PaginationParams) (api_helpers.PaginatedResponse[job.Job], error) {
	// TODO: Handle pagination
	jobCount, err := job.Count()
	if err != nil {
		return api_helpers.PaginatedResponse[job.Job]{}, fmt.Errorf("cannot count total jobs: %w", err)
	}

	jobs, err := job.Search(p)
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
