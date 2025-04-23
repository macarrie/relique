package api

import (
	"fmt"
	"time"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
)

func BackupStart(c client.Client, m module.Module, r repo.Repository) error {
	j := job.NewBackup(c, m, r)
	if err := j.SetupBackup(); err != nil {
		return fmt.Errorf("cannot setup job:  %w", err)
	}

	if err := ClientSSHPing(c); err != nil {
		j.EndTime = time.Now()
		j.Status.Status = job_status.Error
		j.Done = true

		if _, err := j.Save(); err != nil {
			return fmt.Errorf("cannot save job info after failed client ping: %w", err)
		}
		return fmt.Errorf("cannot start backup on unreachable client:  %w", err)
	}

	j.GetLog().Info("Starting job file sync")
	if err := j.Start(); err != nil {
		return fmt.Errorf("error encountered during job execution: %w", err)
	}

	return nil
}
