package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"

	consts "github.com/macarrie/relique/internal/types"
	"github.com/pelletier/go-toml"

	"github.com/macarrie/relique/internal/lib/rsync"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	"github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/config/server_daemon_config"
	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/job_type"
	"github.com/macarrie/relique/internal/types/relique_job"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

func RunJob(job *relique_job.ReliqueJob) error {
	job.GetLog().Info("Starting relique job")

	if err := job.PreFlightCheck(); err != nil {
		return errors.Wrap(err, "cannot start relique job due to incorrect configuration parameters or relique installation")
	}

	job.Status.Status = job_status.Active
	if err := RegisterJob(job); err != nil {
		job.Status.Status = job_status.Error
		job.Done = true
		if _, err := job.Save(); err != nil {
			return errors.Wrap(err, "cannot save job")
		}

		return errors.Wrap(err, "cannot not register job to relique server")
	}

	if err := PingSSHClient(job.Client); err != nil {
		job.Status.Status = job_status.Error
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot ping client via SSH. Aborting job")

		if err := MarkAsDone(job); err != nil {
			return errors.Wrap(err, "cannot mark job as done")
		}
		return errors.Wrap(err, "cannot ping client via SSH, aborting job")
	}

	setupSuccess := true
	if err := LaunchJobSetupOnClient(job); err != nil {
		job.Status.Status = job_status.Error
		setupSuccess = false
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Error encountered during relique job setup on client")

		if err := MarkAsDone(job); err != nil {
			return errors.Wrap(err, "cannot mark job as done")
		}

		return errors.Wrap(err, "cannot perform job setup on client, aborting job")
	}

	preJobScriptSuccess := true
	if err := StartPreScript(job); err != nil {
		preJobScriptSuccess = false
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Error encountered during module pre script execution")
	}

	if setupSuccess && preJobScriptSuccess && (job.Client.SSHAlive == consts.OK) {
		if err := SyncFiles(job); err != nil {
			return errors.Wrap(err, "error occurred when sending files to backup to server")
		}

		if err := StartPostScript(job); err != nil {
			job.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Error encountered during module post script execution")
		}
	}

	if err := MarkAsDone(job); err != nil {
		return errors.Wrap(err, "cannot mark job as done")
	}

	return nil
}

func RegisterJob(j *relique_job.ReliqueJob) error {
	j.GetLog().Info("Registering job")

	if j.JobType.Type == job_type.Restore && j.RestoreJobUuid == "" {
		return fmt.Errorf("restore job has no target job UUID to restore data from")
	}

	switch j.BackupType.Type {
	case backup_type.Diff:
		j.GetLog().Debug("Looking for previous diff jobs to compute diff from")
		previousDiffJob, diffErr := relique_job.PreviousJob(*j, backup_type.BackupType{Type: backup_type.Diff})
		previousCDiffJob, cDiffErr := relique_job.PreviousJob(*j, backup_type.BackupType{Type: backup_type.CumulativeDiff})
		if diffErr == nil {
			j.PreviousJobUuid = previousDiffJob.Uuid
			j.GetLog().WithFields(log.Fields{
				"previous_job_uuid": j.PreviousJobUuid,
			}).Debug("Previous diff job found for diff computation")
		} else {
			if cDiffErr == nil {
				j.PreviousJobUuid = previousCDiffJob.Uuid
				j.GetLog().WithFields(log.Fields{
					"previous_job_uuid": j.PreviousJobUuid,
				}).Info("Previous cumulative diff job found for diff computation")

			} else {
				// Drop back to cumulative diff if no previous diff or cumulative diff found
				j.GetLog().Debug("Previous diff job not found, looking for previous full job to compute diff from")
				previousFullJob, err := relique_job.PreviousJob(*j, backup_type.BackupType{Type: backup_type.Full})
				if err == nil {
					j.PreviousJobUuid = previousFullJob.Uuid
					j.BackupType.Type = backup_type.CumulativeDiff
					j.GetLog().WithFields(log.Fields{
						"previous_job_uuid": j.PreviousJobUuid,
					}).Info("No previous diff backup job found when registering job. This job backup type is changed to 'cumulative_diff'")
				} else {
					j.BackupType.Type = backup_type.Full
					j.GetLog().WithFields(log.Fields{
						"previous_job_uuid": j.PreviousJobUuid,
					}).Info("No previous diff or full backup job found when registering job. This job backup type is changed to 'full'")
				}
			}
		}
	case backup_type.CumulativeDiff:
		j.GetLog().Debug("Looking for previous full jobs to compute cumulative diff from")
		previousFullJob, err := relique_job.PreviousJob(*j, backup_type.BackupType{Type: backup_type.Full})
		if err == nil {
			j.PreviousJobUuid = previousFullJob.Uuid
			j.GetLog().WithFields(log.Fields{
				"previous_job_uuid": j.PreviousJobUuid,
			}).Debug("Previous full job found for cumulative diff computation")
		} else {
			j.BackupType.Type = backup_type.Full
			j.GetLog().Info("No previous backup job found when registering job. This job backup type is changed to 'full'")
		}
	}

	j.StartTime = time.Now()

	if err := j.CreateJobFolder(); err != nil {
		return errors.Wrap(err, "cannot create job folder")
	}

	if j.JobType.Type == job_type.Backup {
		if err := j.CreateJobDataFolder(); err != nil {
			return errors.Wrap(err, "cannot create job data subfolder")
		}
	}

	jobStorageRoot := j.GetJobFolderPath()
	// Save module used to file in job folder path. Modules configuration files can change so we need to keep trace of the exact module used for backup
	moduleExportFile, moduleExportErr := toml.Marshal(j.Module)
	if moduleExportErr != nil {
		return errors.Wrap(moduleExportErr, "cannot serialize job module info to toml data")
	}
	if err := os.WriteFile(fmt.Sprintf("%s/module.toml", jobStorageRoot), moduleExportFile, 0644); err != nil {
		return errors.Wrap(err, "cannot export job module info to file")
	}

	// Save client to file in job folder path. Client configuration files can change so we need to keep trace of the exact client used for backup for later reference
	clientExportFile, clientExportErr := toml.Marshal(j.Client)
	if clientExportErr != nil {
		return errors.Wrap(moduleExportErr, "cannot serialize job client info to toml data")
	}
	if err := os.WriteFile(fmt.Sprintf("%s/client.toml", jobStorageRoot), clientExportFile, 0644); err != nil {
		return errors.Wrap(err, "cannot export job client info to file")
	}

	if _, err := j.Save(); err != nil {
		return errors.Wrap(err, "cannot save job during register")
	}

	return nil
}

func PingSSHClient(c *client.Client) error {
	c.GetLog().Info("Checking SSH connexion with client")

	currentUser, err := user.Current()
	if err != nil || currentUser == nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warning("Cannot get current user to issue warning if trying to ping client from an account different than 'relique'")
	}

	if currentUser != nil && currentUser.Username != "relique" {
		log.WithFields(log.Fields{
			"current_user": currentUser.Username,
		}).Warning("Relique server usually runs as user 'relique' but you are trying to ping client with another user account (probably from cli). SSH ping check can possibly yield false results")
	}
	sshPingCmd := exec.Command("ssh", "-f", "-o BatchMode=yes", fmt.Sprintf("relique@%s", c.Address), "echo 'ping'")
	log.WithFields(log.Fields{
		"cmd":    sshPingCmd,
		"client": c.Name,
	}).Debug("Trying to ping client with following command")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	sshPingCmd.Stdout = &stdout
	sshPingCmd.Stderr = &stderr

	if err := sshPingCmd.Run(); err != nil {
		errorMessage := fmt.Sprintf("cannot ping client via ssh:'%s'", stderr.String())
		c.SSHAlive = consts.CRITICAL
		c.SSHAliveMessage = errorMessage
		return errors.Wrap(err, errorMessage)
	}

	if stderr.String() != "" {
		errorMessage := fmt.Errorf("cannot ping client via ssh:'%s'", stderr.String())
		c.SSHAlive = consts.CRITICAL
		c.SSHAliveMessage = errorMessage.Error()
		return errorMessage
	}

	c.SSHAlive = consts.OK
	return nil
}

func LaunchJobSetupOnClient(job *relique_job.ReliqueJob) error {
	job.GetLog().Debug("Asking client to launch job setup")

	response, err := utils.PerformRequest(
		server_daemon_config.Config,
		job.Client.Address,
		job.Client.Port,
		"POST",
		"/api/v1/job/setup",
		job)
	if err != nil {
		job.Status.Status = job_status.Error
		return errors.Wrap(err, "error when performing api request")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		job.Status.Status = job_status.Error
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "cannot read response body")
		}

		return fmt.Errorf("error during job setup on client, status code '%d': '%s'", response.StatusCode, string(bodyBytes))
	}

	return nil
}

func StartPreScript(job *relique_job.ReliqueJob) error {
	return startModuleScript(job, relique_job.PreScript)
}

func StartPostScript(job *relique_job.ReliqueJob) error {
	return startModuleScript(job, relique_job.PostScript)
}

func startModuleScript(job *relique_job.ReliqueJob, scriptType int) error {
	job.GetLog().WithFields(log.Fields{
		"script_type": scriptType,
	}).Debug("Asking client to launch module script")

	response, err := utils.PerformRequest(
		server_daemon_config.Config,
		job.Client.Address,
		job.Client.Port,
		"POST",
		fmt.Sprintf("/api/v1/job/launch_script/%d", scriptType),
		job)
	if err != nil {
		job.Status.Status = job_status.Error
		return errors.Wrap(err, "error when performing api request")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		job.Status.Status = job_status.Error
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "cannot read response body")
		}

		return fmt.Errorf("error during module script execution, status code '%d': '%s'", response.StatusCode, string(bodyBytes))
	}

	return nil
}

func SyncFiles(job *relique_job.ReliqueJob) error {
	job.GetLog().Info("Starting file sync")
	var wg sync.WaitGroup
	syncHasIncomplete := false
	syncHasError := false

	// TODO: Parse stats from output and save them in job
	for _, path := range job.Module.BackupPaths {
		wg.Add(1)

		var rsyncTask *rsync.Rsync
		if job.JobType.Type == job_type.Backup {
			rsyncTask = rsync.GetBackupTask(job, path)
		} else {
			rsyncTask = rsync.GetRestoreTask(job, path)
		}
		tasksList := append(server_daemon_config.SyncTasks[job.Uuid], rsyncTask)
		server_daemon_config.SyncTasks[job.Uuid] = tasksList

		// TODO: Get stats from output
		go func(task *rsync.Rsync, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
				if _, err := job.Save(); err != nil {
					job.GetLog().Error("Could not save job to db during sync_task task update")
				}
			}()

			rsyncLogFile, err := job.CreateRsyncLogFile(task.Path)
			defer func() {
				if err := rsyncLogFile.Close(); err != nil {
					job.GetLog().WithFields(log.Fields{
						"err": err,
					}).Error("Cannot close regular rsync log file")
				}
			}()
			if err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": task.Path,
				}).Error("Cannot create standard log file for rsync task")
			} else {
				rsyncTask.Cmd.Stdout = rsyncLogFile
				rsyncTask.Cmd.Stderr = rsyncLogFile
			}

			job.GetLog().WithFields(log.Fields{
				"command": task.Cmd.String(),
				"path":    task.Path,
			}).Debug("Starting backup path rsync")

			if err := rsyncTask.Run(); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": task.Path,
				}).Error("Error during path sync")

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

			logPath := job.GetRsyncLogFilePath(task.Path)
			if err := rsyncTask.Stats.GetFromRsyncLog(logPath); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err":      err,
					"log_file": logPath,
				}).Error("Cannot get stats from rsync job log file")
			}
		}(rsyncTask, &wg)
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-ticker.C
			if err := PrintSyncProgress(job); err != nil {
				job.GetLog().Warning("Cannot print job sync tasks progress")
			}
		}
	}()

	wg.Wait()
	ticker.Stop()

	// Print progress at end of job to display 100% progress
	if err := PrintSyncProgress(job); err != nil {
		job.GetLog().Warning("Cannot print job sync tasks progress")
	}

	if syncHasIncomplete || syncHasError {
		if syncHasIncomplete {
			job.Status.Status = job_status.Incomplete
		} else {
			job.Status.Status = job_status.Error
		}
	} else {
		job.Status.Status = job_status.Success
	}

	return nil
}

func PrintSyncProgress(job *relique_job.ReliqueJob) error {
	tasks, ok := server_daemon_config.SyncTasks[job.Uuid]
	if !ok {
		return fmt.Errorf("no sync tasks found for job")
	}

	for _, task := range tasks {
		if err := task.Progress.GetFromRsyncLog(job.GetRsyncLogFilePath(task.Path)); err != nil {
			job.GetLog().WithFields(log.Fields{
				"err":  err,
				"path": task.Path,
			}).Error("Cannot get progress from rsync task")
			continue
		}

		job.GetLog().WithFields(log.Fields{
			"backup_path": task.Path,
			"source":      task.Source,
			"destination": task.Destination,
			"percent":     task.Progress.Percent,
			"total":       task.Progress.Total,
			"current":     task.Progress.Current,
			"remaining":   task.Progress.Remaining,
			"speed":       task.Progress.Speed,
		}).Info("Sync progress")
	}

	return nil
}

func MarkAsDone(job *relique_job.ReliqueJob) error {
	// If job has not been marked incomplete or error, it is successful
	if job.Status.Status == job_status.Active {
		job.Status.Status = job_status.Success
	}

	job.Done = true
	job.EndTime = time.Now()
	job.GetLog().WithFields(log.Fields{
		"updated_done_marker": true,
	}).Info("Updating job done marker")

	if _, err := job.Save(); err != nil {
		return errors.Wrap(err, "cannot save job updated done marker")
	}

	return nil
}
