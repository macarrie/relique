package api

import (
	"errors"
	"testing"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/job_type"
)

func jobCount(t *testing.T) uint64 {
	t.Helper()
	count, err := job.Count(api_helpers.JobSearch{})
	if err != nil {
		t.Fatal(err)
	}
	return count
}

func TestJobRetryUnknown(t *testing.T) {
	setupRemoveTestEnv(t)
	if _, err := JobRetry("does-not-exist"); err == nil {
		t.Errorf("JobRetry() of an unknown job should return an error")
	}
	if got := jobCount(t); got != 0 {
		t.Errorf("JobRetry() of an unknown job created %d jobs", got)
	}
}

func TestJobRetryNotDone(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Backup, false, false)

	if _, err := JobRetry(a.uuid); !errors.Is(err, job.ErrJobNotFinished) {
		t.Errorf("JobRetry() of a running job error = %v, want %v", err, job.ErrJobNotFinished)
	}
	if got := jobCount(t); got != 1 {
		t.Errorf("job count = %d after refused retry, want 1", got)
	}
}

func TestJobRetryRestoreJob(t *testing.T) {
	r := setupRemoveTestEnv(t)
	a := createJob(t, r, job_type.Restore, true, false)

	_, err := JobRetry(a.uuid)
	if !errors.Is(err, job.ErrRestoreNotRetryable) {
		t.Errorf("JobRetry() of a restore job error = %v, want %v", err, job.ErrRestoreNotRetryable)
	}
	if got := jobCount(t); got != 1 {
		t.Errorf("job count = %d after refused retry, want 1", got)
	}
}
