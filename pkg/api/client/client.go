// API Methods used by server daemon
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/macarrie/relique/internal/types/client"

	"github.com/macarrie/relique/internal/types/sync_task"

	"github.com/macarrie/relique/internal/types/job_type"

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
	if err := RegisterJob(job); err != nil {
		job.Status.Status = job_status.Error
		job.Done = true
		return errors.Wrap(err, "cannot not register job to relique server")
	}

	pingSuccess := true
	job.GetLog().Info("Checking SSH connection to relique server")
	if _, err := CheckServerConnection(client.ServerPingParams{
		UseSSH:     true,
		ServerAddr: client_daemon_config.BackupConfig.ServerAddress,
		ServerPort: client_daemon_config.BackupConfig.ServerPort,
	}); err != nil {
		pingSuccess = false
		job.Status.Status = job_status.Error
		job.Done = true
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot confirm if relique server can ping client via SSH. Aborting job")
	}

	preJobScriptSuccess := true
	if err := job.StartPreScript(); err != nil {
		preJobScriptSuccess = false
		job.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Error encountered during module pre script execution")
	}

	if preJobScriptSuccess && pingSuccess {
		if err := SyncFiles(job); err != nil {
			return errors.Wrap(err, "error occurred when sending files to backup to server")
		}

		if err := WaitForSyncCompletion(job); err != nil {
			return errors.Wrap(err, "error during sync completion wait")
		}

		if err := job.StartPostScript(); err != nil {
			job.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Error encountered during module post script execution")
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

func CheckServerConnection(serverPingParams client.ServerPingParams) (client.ServerPingParams, error) {
	params := serverPingParams
	params.UseSSH = serverPingParams.UseSSH
	params.ServerAddr = client_daemon_config.BackupConfig.ServerAddress
	params.ServerPort = client_daemon_config.BackupConfig.ServerPort

	// TODO: If config version is nil, cannot get server address from client conf and server address must be passed
	if client_daemon_config.BackupConfig.Version == "" {
		if serverPingParams.ServerAddr == "" || serverPingParams.ServerPort == 0 {
			params.Message = "client has not received configuration from relique server yet. Server port and address must be specified manually to check client/server connectivity"
			return params, fmt.Errorf(params.Message)
		}
	}

	if serverPingParams.ServerAddr != "" {
		params.ServerAddr = serverPingParams.ServerAddr
	}
	if serverPingParams.ServerPort != 0 {
		params.ServerPort = serverPingParams.ServerPort
	}

	response, err := utils.PerformRequest(client_daemon_config.Config,
		params.ServerAddr,
		params.ServerPort,
		"POST",
		"/api/v1/ping_client",
		serverPingParams)
	if err != nil {
		return params, errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return params, errors.Wrap(err, "cannot read response body from api request")
	}
	defer response.Body.Close()

	var returnParams client.ServerPingParams
	if err := json.Unmarshal(body, &returnParams); err != nil {
		return returnParams, errors.Wrap(err, "cannot read params used for server ping from server")
	}

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"use_ssh":     returnParams.UseSSH,
			"server_addr": returnParams.ServerAddr,
			"server_port": returnParams.ServerPort,
			"err":         returnParams.Message,
		}).Error("Cannot establish connection with relique server")
		return returnParams, fmt.Errorf("cannot ping server: %v", returnParams.Message)
	}
	return returnParams, nil
}

func RegisterJob(job *relique_job.ReliqueJob) error {
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
		var j relique_job.ReliqueJob
		if err := json.Unmarshal(body, &j); err != nil {
			return errors.Wrap(err, "cannot parse job returned from server job register request")
		}
	} else {
		return fmt.Errorf("cannot register job to server (%d response): see server logs for more details", response.StatusCode)
	}

	job.Status.Status = job_status.Active

	return nil
}

func SyncFiles(j *relique_job.ReliqueJob) error {
	for _, path := range j.Module.BackupPaths {
		j.GetLog().WithFields(log.Fields{
			"path": path,
		}).Info("Starting module path backup")

		// Create restore destination directory on restore before starting sync
		if j.JobType.Type == job_type.Restore {
			var targetRestoreDir string
			if j.RestoreDestination == "" {
				targetRestoreDir = filepath.Clean(path)
			} else {
				targetRestoreDir = filepath.Clean(fmt.Sprintf("%s/%s", j.RestoreDestination, path))
			}
			if err := os.MkdirAll(targetRestoreDir, 0755); err != nil {
				j.GetLog().WithFields(log.Fields{
					"path": targetRestoreDir,
					"err":  err,
				}).Error("Cannot create restore destination directory before starting file sync")
				j.Status.Status = job_status.Incomplete
				continue
			}
		} else {
			if _, err := os.Lstat(path); os.IsNotExist(err) {
				j.GetLog().WithFields(log.Fields{
					"path": path,
				}).Error("Backup path does not exist on client")
				j.Status.Status = job_status.Incomplete
				continue
			}
		}
	}

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"POST",
		fmt.Sprintf("/api/v1/backup/jobs/%s/sync", j.Uuid),
		nil)
	if err != nil || j.Uuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot start files sync on server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func WaitForSyncCompletion(job *relique_job.ReliqueJob) error {
	job.GetLog().Info("Waiting for sync tasks completion")

	var ticker *time.Ticker

	hasSuccess := false
	ticker = time.NewTicker(10 * time.Second)
	for {
		<-ticker.C
		progress, err := GetJobProgress(*job)
		if err != nil {
			job.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Cannot get job progress from server")
			continue
		}

		allDone := true
		for _, task := range progress {
			job.GetLog().WithFields(log.Fields{
				"path":      task.Path,
				"remaining": task.Progress.Remain,
				"total":     task.Progress.Total,
				"progress":  task.Progress.Progress,
				"speed":     task.Progress.Speed,
			}).Info("Sync progress")

			if task.Error == "" {
				hasSuccess = true
			} else {
				job.GetLog().WithFields(log.Fields{
					"err": task.Error,
				}).Error("Backup path sync error. Check log file on server for more info")
				job.Status.Status = job_status.Incomplete
			}

			if !task.Done {
				allDone = false
			}
		}

		if allDone {
			job.GetLog().Info("All sync tasks completed")
			break
		}
	}

	// Job has not been marked incomplete and all backup paths sync are done -> Success !
	if job.Status.Status == job_status.Active {
		job.Status.Status = job_status.Success
	}
	if job.Status.Status == job_status.Incomplete && !hasSuccess {
		job.Status.Status = job_status.Error
	}

	return nil
}

func GetJobProgress(job relique_job.ReliqueJob) ([]sync_task.SyncTaskProgress, error) {
	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"GET",
		fmt.Sprintf("/api/v1/backup/jobs/%s/sync_progress", job.Uuid),
		nil)
	if err != nil || job.Uuid == "" {
		return []sync_task.SyncTaskProgress{}, errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []sync_task.SyncTaskProgress{}, errors.Wrap(err, "cannot read response body from api request")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var progress []sync_task.SyncTaskProgress
		if err := json.Unmarshal(body, &progress); err != nil {
			return []sync_task.SyncTaskProgress{}, errors.Wrap(err, "cannot parse sync progress returned from server")
		}

		return progress, nil
	} else {
		return []sync_task.SyncTaskProgress{}, fmt.Errorf("cannot get sync_task task completion status from server (%d response): see server logs for more details", response.StatusCode)
	}
}

func UpdateJobStatus(job relique_job.ReliqueJob) error {
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
