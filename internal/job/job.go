package job

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/image"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/job_type"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	"github.com/macarrie/relique/internal/rsync_task"
	rsync_lib "github.com/macarrie/relique/internal/rsync_task/lib"
	"github.com/macarrie/relique/internal/utils"
)

func NewBackup(c client.Client, m module.Module, r repo.Repository) Job {
	return Job{
		Uuid:       uuid.New().String(),
		JobType:    job_type.New(job_type.Backup),
		Client:     c,
		Module:     m,
		Repository: r,
		Status:     job_status.New(job_status.Pending),
		BackupType: m.BackupType,
	}
}

func NewRestore(img image.Image, targetClient client.Client, restorePaths map[string]string) Job {
	if len(restorePaths) != 0 {
		img.Module.Name = "on-demand"
	}

	return Job{
		Uuid:               uuid.New().String(),
		JobType:            job_type.New(job_type.Restore),
		Client:             targetClient,
		Module:             img.Module,
		Repository:         img.Repository,
		Status:             job_status.New(job_status.Pending),
		RestoreImageUuid:   img.Uuid,
		CustomRestorePaths: restorePaths,
	}
}

func (j *Job) SetupBackup() error {
	j.GetLog().Debug("Starting job setup")

	if j.BackupType.Type == backup_type.Diff {
		j.GetLog().Debug("Looking for previous diff jobs to compute diff from")
		previousDiffJob, diffErr := GetPrevious(*j, backup_type.BackupType{Type: backup_type.Diff})
		if diffErr == nil {
			j.PreviousJobUuid = previousDiffJob.Uuid
			j.PreviousJob = &previousDiffJob
			j.GetLog().With(
				slog.String("previous_job_uuid", j.PreviousJobUuid),
			).Debug("Previous diff job found for diff computation")
		} else {
			// Drop back to cumulative diff if no previous diff or cumulative diff found
			j.GetLog().Debug("Previous diff job not found, looking for previous full job to compute diff from")
			previousFullJob, err := GetPrevious(*j, backup_type.BackupType{Type: backup_type.Full})
			if err == nil {
				j.PreviousJobUuid = previousFullJob.Uuid
				j.PreviousJob = &previousFullJob
				j.GetLog().With(
					slog.String("previous_job_uuid", j.PreviousJobUuid),
				).Info("No previous diff backup job found when registering job. The image generated during the previous full backup is used as reference.")
			} else {
				j.BackupType.Type = backup_type.Full
				j.GetLog().With(
					slog.String("previous_job_uuid", j.PreviousJobUuid),
				).Info("No previous diff or full backup job found when registering job. This job backup type is changed to 'full'")
			}
		}
	}

	j.GetLog().Debug("Creating job storage folder")
	jobFolderPath, err := j.GetStorageFolderPath()
	if err != nil {
		return fmt.Errorf("cannot determine job storage folder: %w", err)
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/_data", jobFolderPath), 0755); err != nil {
		return fmt.Errorf("cannot setup job data folder: %w", err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/_logs", jobFolderPath), 0755); err != nil {
		return fmt.Errorf("cannot setup job logs folder: %w", err)
	}

	var tasks []rsync_task.RsyncTask
	for _, backupPath := range j.Module.BackupPaths {
		switch j.BackupType.Type {
		case backup_type.Full:
			tasks = append(tasks, rsync_task.NewFullBackup(
				// Source
				fmt.Sprintf("%s@%s:%s", j.Client.SSHUser, j.Client.Address, backupPath),
				// Destination
				filepath.Clean(fmt.Sprintf("%s/_data/", jobFolderPath)),
				// Log folder
				filepath.Clean(fmt.Sprintf("%s/_logs/", jobFolderPath)),
				// Backup path
				backupPath,
				// Exclusions/inclusions
				j.Module.Exclude,
				j.Module.ExcludeCVS,
				j.Module.Include,
			))

		case backup_type.Diff:
			previousJobFolderPath, err := j.PreviousJob.GetStorageFolderPath()
			if err != nil {
				return fmt.Errorf("cannot determine job storage folder: %w", err)
			}
			tasks = append(tasks, rsync_task.NewDiffBackup(
				// Source
				fmt.Sprintf("%s@%s:%s", j.Client.SSHUser, j.Client.Address, backupPath),
				// Destination
				filepath.Clean(fmt.Sprintf("%s/_data/", jobFolderPath)),
				// Previous job folder for comparison
				filepath.Clean(fmt.Sprintf("%s/_data/", previousJobFolderPath)),
				// Log folder
				filepath.Clean(fmt.Sprintf("%s/_logs/", jobFolderPath)),
				// Backup path
				backupPath,
				// Exclusions/inclusions
				j.Module.Exclude,
				j.Module.ExcludeCVS,
				j.Module.Include,
			))
		default:
			return fmt.Errorf("unknown backup type '%s'", j.BackupType.String())
		}
	}
	j.Tasks = tasks

	j.GetLog().Debug("Creating job catalog folder")
	jobCatalogPath := j.GetCatalogPath()

	if err := os.MkdirAll(jobCatalogPath, 0755); err != nil {
		return fmt.Errorf("cannot setup job catalog folder: %w", err)
	}

	// Save module used to file in job folder path. Modules configuration files can change so we need to keep trace of the exact module used for backup
	if err := utils.SerializeToFile[module.Module](j.Module, fmt.Sprintf("%s/module.toml", jobCatalogPath)); err != nil {
		return fmt.Errorf("cannot export module to file: %w", err)
	}

	// Save client to file in job folder path. Client configuration files can change so we need to keep trace of the exact client used for backup for later reference
	if err := utils.SerializeToFile[client.Client](j.Client, fmt.Sprintf("%s/client.toml", jobCatalogPath)); err != nil {
		return fmt.Errorf("cannot export client to file: %w", err)
	}

	// Save repo to file in job folder path. Repo configuration files can change so we need to keep trace of the exact repo used for backup for later reference
	if err := utils.SerializeToFile[repo.Repository](j.Repository, fmt.Sprintf("%s/repo.toml", jobCatalogPath)); err != nil {
		return fmt.Errorf("cannot export repository to file: %w", err)
	}

	if _, err := j.Save(); err != nil {
		return fmt.Errorf("cannot save job info to database after setup complete: %w", err)
	}

	// TODO: Add job setup event
	return nil
}

func (j *Job) SetupRestore() error {
	j.GetLog().Debug("Starting job setup")

	if j.JobType.Type == job_type.Restore && j.RestoreImageUuid == "" {
		return fmt.Errorf("restore job has no target image UUID to restore data from")
	}

	j.GetLog().Debug("Creating job storage folder")
	jobFolderPath, err := j.GetStorageFolderPath()
	if err != nil {
		return fmt.Errorf("cannot determine job storage folder: %w", err)
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/_logs", jobFolderPath), 0755); err != nil {
		return fmt.Errorf("cannot setup job logs folder: %w", err)
	}

	restoreSourceJob, err := GetByUuid(j.RestoreImageUuid)
	if err != nil {
		return fmt.Errorf("cannot get restore source job from db: %w", err)
	}
	restoreSourceFolderPath, err := restoreSourceJob.GetStorageFolderPath()
	if err != nil {
		return fmt.Errorf("cannot determine job storage folder: %w", err)
	}

	var tasks []rsync_task.RsyncTask
	if len(j.CustomRestorePaths) == 0 {
		for _, backupPath := range j.Module.BackupPaths {
			tasks = append(tasks, rsync_task.NewRestore(
				// Source
				fmt.Sprintf("%s/_data/%s", restoreSourceFolderPath, backupPath),
				// Destination
				fmt.Sprintf("%s@%s:%s", j.Client.SSHUser, j.Client.Address, backupPath),
				// Log folder
				filepath.Clean(fmt.Sprintf("%s/_logs/", jobFolderPath)),
				// Backup path
				backupPath,
				// Inclusions/exclusions
				j.Module.Exclude,
				j.Module.ExcludeCVS,
				j.Module.Include,
			))
		}
	} else {
		for source, dest := range j.CustomRestorePaths {
			tasks = append(tasks, rsync_task.NewRestore(
				// Source
				fmt.Sprintf("%s/_data/%s", restoreSourceFolderPath, source),
				// Destination
				fmt.Sprintf("%s@%s:%s", j.Client.SSHUser, j.Client.Address, dest),
				// Log folder
				filepath.Clean(fmt.Sprintf("%s/_logs/", jobFolderPath)),
				// Backup path
				source,
				// Inclusions/exclusions
				j.Module.Exclude,
				j.Module.ExcludeCVS,
				j.Module.Include,
			))
		}
	}
	j.Tasks = tasks

	j.GetLog().Debug("Creating job catalog folder")
	jobCatalogPath := j.GetCatalogPath()

	if err := os.MkdirAll(jobCatalogPath, 0755); err != nil {
		return fmt.Errorf("cannot setup job catalog folder: %w", err)
	}

	// Save module used to file in job folder path. Modules configuration files can change so we need to keep trace of the exact module used for backup
	if err := utils.SerializeToFile[module.Module](j.Module, fmt.Sprintf("%s/module.toml", jobCatalogPath)); err != nil {
		return fmt.Errorf("cannot export module to file: %w", err)
	}

	// Save client to file in job folder path. Client configuration files can change so we need to keep trace of the exact client used for backup for later reference
	if err := utils.SerializeToFile[client.Client](j.Client, fmt.Sprintf("%s/client.toml", jobCatalogPath)); err != nil {
		return fmt.Errorf("cannot export client to file: %w", err)
	}

	// Save repo to file in job folder path. Repo configuration files can change so we need to keep trace of the exact repo used for backup for later reference
	if err := utils.SerializeToFile[repo.Repository](j.Repository, fmt.Sprintf("%s/repo.toml", jobCatalogPath)); err != nil {
		return fmt.Errorf("cannot export repository to file: %w", err)
	}

	if _, err := j.Save(); err != nil {
		return fmt.Errorf("cannot save job info to database after setup complete: %w", err)
	}

	// TODO: Add job setup event
	return nil
}

func (j *Job) Start() error {
	j.GetLog().Debug("Starting job sync tasks")
	j.Status.Status = job_status.Active
	j.StartTime = time.Now()

	if _, err := j.Save(); err != nil {
		return fmt.Errorf("cannot save job info to database before start: %w", err)
	}

	ticker := time.NewTicker(1 * time.Second)
	var wg sync.WaitGroup
	syncHasIncomplete := false
	syncHasError := false

	for i, _ := range j.Tasks {
		wg.Add(1)

		go func(task *rsync_task.RsyncTask) {
			for {
				<-ticker.C
				task.GetProgressLog().Info("File sync in progress")
			}
		}(&j.Tasks[i])

		go func(task *rsync_task.RsyncTask) {
			defer wg.Done()

			slog.With(
				slog.String("cmd", task.Task.Rsync.Cmd.String()),
			).Debug("Running rsync command")

			if err := task.Task.Run(); err != nil {
				slog.With(slog.Any("error", err)).Error("Error encountered during task run")

				// Rsync exit codes 23,24,25 mean that some files still may have been transferred even if exit code is != 0.
				// Treating those exit codes as partial success/error
				// 23 - Partial transfer due to error
				// 24 - Partial transfer due to vanished source files
				// 25 - The --max-delete limit stopped deletions
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.ExitCode() >= 23 || exitErr.ExitCode() <= 25 {
						syncHasIncomplete = true
					} else {
						syncHasError = true
					}
				}
			}

			logStruct := task.Task.Log()
			if err := os.WriteFile(task.LogFile, []byte(logStruct.Stdout), 0755); err != nil {
				slog.With(slog.Any("error", err)).Error("Cannot write task log to file")
			}
			if err := os.WriteFile(task.LogErrorFile, []byte(logStruct.Stderr), 0755); err != nil {
				slog.With(slog.Any("error", err)).Error("Cannot write task error log to file")
			}

			err := task.Task.Stats.GetFromRsyncLog(task.LogFile)
			if err != nil {
				slog.With(slog.Any("error", err)).Error("Cannot write task error log to file")
			}

		}(&j.Tasks[i])
	}

	wg.Wait()
	ticker.Stop()

	for i, _ := range j.Tasks {
		// Print progress at least once at the end with sync stats
		j.GetLog().With(
			slog.String("backup_path", j.Tasks[i].BackupPath),
			slog.Group("stats",
				slog.Int("elements_nb", j.Tasks[i].Task.Stats.NumberOfFiles),
				slog.Int("files_nb", j.Tasks[i].Task.Stats.NumberOfRegularFiles),
				slog.Int("folder_nb", j.Tasks[i].Task.Stats.NumberOfDirectories),
				slog.Int("total_file_size", int(j.Tasks[i].Task.Stats.TotalFileSize)),
				slog.Int("literal_data", int(j.Tasks[i].Task.Stats.LiteralData)),
				slog.Int("matched_data", int(j.Tasks[i].Task.Stats.MatchedData)),
				slog.Int("total_bytes_sent", int(j.Tasks[i].Task.Stats.TotalBytesSent)),
				slog.Int("total_bytes_received", int(j.Tasks[i].Task.Stats.TotalBytesReceived)),
				slog.Int("transfer_speed", int(j.Tasks[i].Task.Stats.TransferSpeed)),
				slog.Int("transfer_speedup", int(j.Tasks[i].Task.Stats.TransferSpeedup)),
			)).Info("File sync complete")
	}

	if syncHasIncomplete || syncHasError {
		if syncHasIncomplete {
			j.Status.Status = job_status.Incomplete
		} else {
			j.Status.Status = job_status.Error
		}
	} else {
		j.Status.Status = job_status.Success
	}

	j.Done = true
	j.EndTime = time.Now()
	if _, err := j.Save(); err != nil {
		return fmt.Errorf("cannot save job info to database after completion: %w", err)
	}

	catalogPath := j.GetCatalogPath()
	jobStats := rsync_task.MergeStats(j.Tasks)
	if err := utils.SerializeToFile[rsync_lib.Stats](jobStats, fmt.Sprintf("%s/stats.toml", catalogPath)); err != nil {
		return fmt.Errorf("cannot export job stats to file: %w", err)
	}

	if j.JobType.Type == job_type.Backup {
		if j.Status.Status == job_status.Success || j.Status.Status == job_status.Incomplete {
			j.GetLog().Info("Generating backup image from job")
			img := image.New(j.Client, j.Module, j.Repository)
			img.Uuid = j.Uuid
			if err := img.FillStats(jobStats, catalogPath); err != nil {
				return fmt.Errorf("cannot get image stats: %w", err)
			}
			if _, err := img.Save(); err != nil {
				slog.With(slog.Any("error", err)).Error("Cannot save generated image to database")
			}
		} else {
			j.GetLog().Info("No image generated for unsuccessful job")
		}
	}

	return nil
}
