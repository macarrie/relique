package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/utils"
)

func RestoreStart(targetClient client.Client, img image.Image, rawCustomPathRestore []string) error {
	if err := ClientSSHPing(targetClient); err != nil {
		return fmt.Errorf("cannot start restore on unreachable client:  %w", err)
	}

	restorePaths := utils.GenerateCustomRestorePaths(rawCustomPathRestore, img.Module.BackupPaths)

	j := job.NewRestore(img, targetClient, restorePaths)
	if err := j.SetupRestore(); err != nil {
		return fmt.Errorf("cannot setup job:  %w", err)
	}

	if err := j.Start(); err != nil {
		return fmt.Errorf("error encountered during job execution: %w", err)
	}

	return nil
}
