package api

import (
	"fmt"

	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job"
	"github.com/macarrie/relique/internal/utils"
)

func RestoreStart(targetClient client.Client, img image.Image, rawCustomPathRestore []string) error {
	restorePaths := utils.GenerateCustomRestorePaths(rawCustomPathRestore, img.Module.BackupPaths)

	j := job.NewRestore(img, targetClient, restorePaths)
	if err := j.SetupRestore(); err != nil {
		err = fmt.Errorf("cannot setup job:  %w", err)
		j.RecordError(err)
		return err
	}

	return startJob(&j)
}
