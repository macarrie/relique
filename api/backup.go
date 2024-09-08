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

	j := job.New(c, m, r)
	j.GetLog().Info("Starting job")

	if err := j.Setup(); err != nil {
		return fmt.Errorf("cannot setup job:  %w", err)
	}
	if err := j.Start(); err != nil {
		return fmt.Errorf("error encountered during job execution: %w", err)
	}

	return nil
}
