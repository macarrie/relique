// API Methods used by server daemon
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/macarrie/relique/internal/types/job_type"

	"github.com/macarrie/relique/internal/types/backup_type"

	"github.com/macarrie/relique/internal/types/rsync"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/macarrie/relique/internal/types/job_status"

	"github.com/macarrie/relique/internal/types/config/client_daemon_config"

	"github.com/macarrie/relique/internal/types/relique_job"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

func RunJob(job *relique_job.ReliqueJob) error {
	job.GetLog().Info("Starting relique job")

	job.Status.Status = job_status.Active
	if err := SSHPing(client_daemon_config.BackupConfig.ServerAddress); err != nil {
		job.Status.Status = job_status.Error
		job.Done = true
		return errors.Wrap(err, "cannot connect to relique server via SSH")
	}

	if err := RegisterJob(job); err != nil {
		job.Status.Status = job_status.Error
		job.Done = true
		return errors.Wrap(err, "cannot not register job to relique server")
	}

	if job.JobType.Type == job_type.Backup {
		// TODO: Run script
		if err := job.StartPreBackupScript(); err != nil {
			return errors.Wrap(err, "error occurred during pre backup script execution")
		}

		if err := SendFiles(job); err != nil {
			return errors.Wrap(err, "error occurred when sending files to backup to server")
		}

		// TODO: Run script
		if err := job.StartPostBackupScript(); err != nil {
			return errors.Wrap(err, "error occurred during pre backup script execution")
		}
	} else if job.JobType.Type == job_type.Restore {
		// TODO: Run script
		if err := job.StartPreRestoreScript(); err != nil {
			return errors.Wrap(err, "error occurred during pre backup script execution")
		}

		if err := GetRestoreFileList(job); err != nil {
			return errors.Wrap(err, "error occurred when getting file list to restore to server")
		}

		if err := DownloadFiles(job); err != nil {
			return errors.Wrap(err, "error occurred when restoring files from server")
		}

		// TODO: Run script
		if err := job.StartPostRestoreScript(); err != nil {
			return errors.Wrap(err, "error occurred during pre backup script execution")
		}
	}

	if err := UpdateJobStatus(*job); err != nil {
		return errors.Wrap(err, "cannot update job status in relique server")
	}

	job.Done = true
	if err := MarkAsDone(*job); err != nil {
		return errors.Wrap(err, "cannot mark job as done in relique server")
	}

	return nil
}

func SSHPing(addr string) error {
	log.Debug("Performing SSH ping")

	cmd := exec.Command("ssh", addr, "/bin/true")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("could not ping address '%s' via SSH: '%s'", addr, string(out))
	}

	return nil
}

func RegisterJob(job *backup_job.BackupJob) error {
	job.GetLog().Info("Registering job to relique server")

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"POST",
		"/api/v1/backup/register_job",
		job)
	if err != nil || job.Uuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body from api request")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var j backup_job.BackupJob
		if err := json.Unmarshal(body, &j); err != nil {
			return errors.Wrap(err, "cannot parse job returned from server job register request")
		}
		if j.StorageDestination == "" {
			return fmt.Errorf("got empty storage destination from server job register request response")
		}

		job.StorageDestination = j.StorageDestination

		if j.BackupType.Type == backup_type.Diff {
			if j.PreviousJobStorageDestination == "" {
				return fmt.Errorf("got empty previous storage destination from server job register request response")
			}

			job.PreviousJobStorageDestination = j.PreviousJobStorageDestination
		}
	} else {
		return fmt.Errorf("cannot register job to server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func SendFiles(j *backup_job.BackupJob) error {
	// TODO: Execute rsync tasks in parallel for perf gainz. Use wggroup to sync partial rsync jobs ?
	var jobStatus uint8 = job_status.Active
	hasBackupPathsSuccess := false
	for _, path := range j.Module.BackupPaths {
		j.GetLog().WithFields(log.Fields{
			"path": path,
		}).Info("Starting module path backup")

		if _, err := os.Lstat(path); os.IsNotExist(err) {
			j.GetLog().WithFields(log.Fields{
				"path": path,
			}).Error("Backup path does not exist on client")
			jobStatus = job_status.Incomplete
			continue
		}

		syncTask := rsync.GetSyncTask(j, path)
		j.RSyncTasks = append(j.RSyncTasks, syncTask)

		if err := syncTask.Run(); err != nil {
			j.GetLog().WithFields(log.Fields{
				"err":  err,
				"path": path,
			}).Error("Error during backup path rsync backup")
			jobStatus = job_status.Incomplete
			continue
		}

		taskLog := syncTask.Log()
		rsyncLogFile, err := j.GetRsyncLogFile(path)
		defer rsyncLogFile.Close()
		if err != nil {
			j.GetLog().WithFields(log.Fields{
				"err":  err,
				"path": path,
			}).Error("Cannot create log file for rsync task")
		} else {
			if _, err := rsyncLogFile.WriteString(taskLog.Stdout); err != nil {
				j.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": path,
				}).Error("Cannot write rsync log to log file")
			}
		}
		rsyncErrorLogFile, err := j.GetRsyncErrorLogFile(path)
		defer rsyncErrorLogFile.Close()
		if err != nil {
			j.GetLog().WithFields(log.Fields{
				"err":  err,
				"path": path,
			}).Error("Cannot create error log file for rsync task")
		} else {
			if _, err := rsyncErrorLogFile.WriteString(taskLog.Stdout); err != nil {
				j.GetLog().WithFields(log.Fields{
					"err":  err,
					"path": path,
				}).Error("Cannot write rsync error log to log file")
			}
		}

		hasBackupPathsSuccess = true
	}

	if !hasBackupPathsSuccess && jobStatus == job_status.Incomplete {
		jobStatus = job_status.Error
	}

	j.Status.Status = jobStatus

	// If job has not been marked as Incomplete or Error yet and is still active, this means it's a success. Mark it as such
	if jobStatus == job_status.Active {
		j.Status.Status = job_status.Success
	}

	return nil
}

func UpdateJobStatus(job backup_job.BackupJob) error {
	job.GetLog().Info("Update job status to relique server")

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"PUT",
		fmt.Sprintf("/api/v1/backup/jobs/%s/status", job.Uuid),
		job.Status)
	if err != nil || job.Uuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot update job status to server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func MarkAsDone(job relique_job.ReliqueJob) error {
	job.GetLog().Info("Mark job as done in relique server")

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"PUT",
		fmt.Sprintf("/api/v1/backup/jobs/%s/done", job.Uuid),
		job.Done)
	if err != nil || job.Uuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot mark job as done on server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}
