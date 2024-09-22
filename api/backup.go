package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
)

func BackupStart(c client.Client, m module.Module, r repo.Repository) error {
	if err := ClientSSHPing(c); err != nil {
		return fmt.Errorf("cannot start backup on unreachable client:  %w", err)
	}

	j := job.NewBackup(c, m, r)
	if err := j.SetupBackup(); err != nil {
		return fmt.Errorf("cannot setup job:  %w", err)
	}

	j.GetLog().Info("Starting job")
	if err := j.Start(); err != nil {
		return fmt.Errorf("error encountered during job execution: %w", err)
	}

	return nil
}
