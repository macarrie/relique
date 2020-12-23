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

	"github.com/macarrie/relique/internal/types/backup_type"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/macarrie/relique/internal/types/backup_item"

	"github.com/macarrie/relique/internal/types/job_status"

	"github.com/macarrie/relique/internal/types/config/client_daemon_config"

	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

func RunJob(job *backup_job.BackupJob) error {
	// TODO: Update job status on errors
	job.GetLog().Info("Starting backup job")

	job.Status.Status = job_status.Active
	if err := RegisterJob(*job); err != nil {
		return errors.Wrap(err, "cannot not register job to relique server")
	}

	// TODO
	if err := job.StartPreBackupScript(); err != nil {
		return errors.Wrap(err, "error occurred during pre backup script execution")
	}

	// TODO: Do not return on error (post backup script may restart stopped service for example) ?
	if err := SendFiles(job); err != nil {
		return errors.Wrap(err, "error occurred when sending files to backup to server")
	}

	// TODO
	if err := job.StartPostBackupScript(); err != nil {
		return errors.Wrap(err, "error occurred during pre backup script execution")
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

func SendFiles(j *backup_job.BackupJob) error {
	// TODO: Send diff and signature with multipart streams
	// TODO: Send files in parallel for more performance
	// TODO: Do not compute diff for full backups to speed up backup process, send files directly
	var jobStatus uint8 = job_status.Active
	hasBackupPathsSuccess := false
	for _, path := range j.Module.BackupPaths {
		j.GetLog().WithFields(log.Fields{
			"path": path,
		}).Info("Starting module path backup")

		// Returning err from within the walk function stop further file processing.
		//Instead log and handle the error and return nil to continue backing up the other files
		err := filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"path": walkPath,
				}).Error("Cannot read file or directory")
				jobStatus = job_status.Incomplete
				return err
			}

			bkpItem := backup_item.New(j.Uuid, walkPath, info)

			if !info.IsDir() {
				if err := GetChecksum(bkpItem); err != nil {
					// If checksum cannot be computed, try sending the whole file directly instead to ensure file is saved anyway
					bkpItem.GetLog().WithFields(log.Fields{
						"err": err,
					}).Error("Cannot get checksum from server. Sending whole file instead")
					bkpItem.Exists = false
				}

				if !bkpItem.Exists || j.BackupType.Type == backup_type.Full {
					bkpItem.GetLog().Info("Item to backup does not exist yet on server. Sending file directly without computing diff")
					if err := SendRawFile(bkpItem); err != nil {
						bkpItem.GetLog().WithFields(log.Fields{
							"err": err,
						}).Error("Cannot send raw file to server")
						jobStatus = job_status.Incomplete
						return nil
					}
					// If success, backup done for this item
					return nil
				}

				localHash, err := bkpItem.GetLocalChecksum()
				if err != nil {
					bkpItem.GetLog().WithFields(log.Fields{
						"err": err,
					}).Warning("Cannot compute local checksum for file version comparison. Diff will be computed for this file even if it could have been unnecessary")
				}
				// Only send file if backup full or does not exist on server
				if bytes.Equal(localHash, bkpItem.Checksum) {
					bkpItem.GetLog().Info("Checksum match between server and client. Skipping file upload but applying diff to create hardlink")
					// TODO: Handle rights and file times
					bkpItem.CreateHardlink = true
					if err := ApplyDiff(bkpItem); err != nil {
						bkpItem.GetLog().WithFields(log.Fields{
							"err": err,
						}).Error("Cannot apply diff on server. Check server logs for more details")
						jobStatus = job_status.Incomplete
						return nil
					}
					return nil
				}

				if err := GetSignature(bkpItem); err != nil {
					bkpItem.GetLog().WithFields(log.Fields{
						"err": err,
					}).Error("Cannot get signature from server")
					jobStatus = job_status.Incomplete
					return nil
				}

				if bkpItem.Signature == nil && bkpItem.Exists {
					bkpItem.GetLog().Error("Got nil signature from server")
					jobStatus = job_status.Incomplete
					return nil
				}

				bkpItem.GetDiff()
			}

			if err := ApplyDiff(bkpItem); err != nil {
				bkpItem.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Cannot apply diff on server. Check server logs for more details")
				jobStatus = job_status.Incomplete
				return nil
			}

			hasBackupPathsSuccess = true
			return nil
		})
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"path": path,
			}).Error("Cannot back up path")
			jobStatus = job_status.Error
			continue
		}
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

func SendRawFile(b *backup_item.BackupItem) error {
	b.GetLog().Info("Sending raw file")

	backupItemJson, err := json.Marshal(b)
	if err != nil {
		return errors.Wrap(err, "cannot serialize backup item")
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		m.WriteField("item", string(backupItemJson))
		part, err := m.CreateFormFile("file", b.Path)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot create form file for multipart send")
			return
		}
		if !b.IsSymlink {
			file, err := os.Open(b.Path)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot open file for multipart send")
				return
			}
			defer file.Close()

			if _, err = io.Copy(part, file); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Cannot copy file for multipart send")
				return
			}
		}
	}()

	response, err := utils.SendMultipart(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		fmt.Sprintf("/api/v1/backup/jobs/%s/file", b.JobUuid),
		m.FormDataContentType(),
		r)
	if err != nil || b.JobUuid == "" {
		return errors.Wrap(err, "error when performing multipart api request")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot send backup item raw file to server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func RegisterJob(job backup_job.BackupJob) error {
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

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot register job to server (%d response): see server logs for more details", response.StatusCode)
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

func MarkAsDone(job backup_job.BackupJob) error {
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

func GetChecksum(item *backup_item.BackupItem) error {
	item.GetLog().Info("Getting item checksum from server")

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"POST",
		fmt.Sprintf("/api/v1/backup/jobs/%s/checksum", item.JobUuid),
		item)
	if err != nil || item.JobUuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body from api requets")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var i backup_item.BackupItem
		if err := json.Unmarshal(body, &i); err != nil {
			return errors.Wrap(err, "cannot parse item returned from server checksum computation request")
		}
		item.Checksum = i.Checksum
		item.Exists = i.Exists

		return nil
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot get item checksum from server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func GetSignature(item *backup_item.BackupItem) error {
	item.GetLog().Info("Getting item signature from server")

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"POST",
		fmt.Sprintf("/api/v1/backup/jobs/%s/signature", item.JobUuid),
		item)
	if err != nil || item.JobUuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body from api requets")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var i backup_item.BackupItem
		if err := json.Unmarshal(body, &i); err != nil {
			return errors.Wrap(err, "cannot parse item returned from server signature computation request")
		}
		item.Signature = i.Signature

		return nil
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot get item signature from server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func ApplyDiff(item *backup_item.BackupItem) error {
	item.GetLog().Info("Applying diff on server")

	response, err := utils.PerformRequest(client_daemon_config.Config,
		client_daemon_config.BackupConfig.ServerAddress,
		client_daemon_config.BackupConfig.ServerPort,
		"POST",
		fmt.Sprintf("/api/v1/backup/jobs/%s/apply_diff", item.JobUuid),
		item)
	if err != nil || item.JobUuid == "" {
		return errors.Wrap(err, "error when performing api request")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot apply item diff on server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}
