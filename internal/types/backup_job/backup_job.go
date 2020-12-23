package backup_job

import (
	"time"

	"github.com/macarrie/relique/internal/db"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/module"
	"github.com/pkg/errors"
)

type BackupJob struct {
	// Database IDs
	ID       int64
	ModuleID int64
	ClientID int64

	Uuid       string
	Client     clientObject.Client
	Module     module.Module
	Status     job_status.JobStatus
	Done       bool
	BackupType backup_type.BackupType
	StartTime  time.Time
	EndTime    time.Time
}

func (j *BackupJob) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"uuid":        j.Uuid,
		"client":      j.Client.String(),
		"module":      j.Module.String(),
		"backup_type": j.BackupType.String(),
		"status":      j.Status.String(),
		"done":        j.Done,
	})
}

func (j *BackupJob) Save() (int64, error) {
	if j.ID != 0 {
		return j.Update()
	}

	j.GetLog().Debug("Saving job into database")

	moduleId, err := j.Module.Save()
	if err != nil || moduleId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner module")
	}

	clientId, err := j.Client.Save()
	if err != nil || clientId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner client")
	}

	sql := "INSERT INTO jobs (uuid, status, backup_type, done, module_id, client_id, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	result, err := db.Write().Exec(
		sql,
		j.Uuid,
		j.Status.Status,
		j.BackupType.Type,
		j.Done,
		moduleId,
		clientId,
		j.StartTime,
		j.EndTime,
	)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot save job into db")
	}

	j.ID, _ = result.LastInsertId()

	return j.ID, nil
}

func (j *BackupJob) UpdateStatus() (int64, error) {
	j.GetLog().Debug("Updating job status into database")

	request := `UPDATE jobs  
		SET status = $1
		WHERE uuid = $2`

	// TODO: Save get module and client ID
	result, err := db.Write().Exec(
		request,
		j.Status.Status,
		j.Uuid,
	)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot update job status into db")
	}

	j.ID, _ = result.LastInsertId()

	return j.ID, nil
}

func (j *BackupJob) Update() (int64, error) {
	// TODO: Save job in DB
	j.GetLog().Debug("Updating job details into database")

	var moduleId int64
	var clientId int64

	moduleId, err := j.Module.Save()
	if err != nil || moduleId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner module")
	}
	clientId, err = j.Client.Save()
	if err != nil || clientId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner client")
	}

	sql := `UPDATE jobs  
SET status = $1, 
	backup_type = $2, 
	module_id = $3, 
	client_id = $4,
    done = $5,
    start_time = $6,
    end_time = $7
WHERE uuid = $8;`

	result, err := db.Write().Exec(
		sql,
		j.Status.Status,
		j.BackupType.Type,
		moduleId,
		clientId,
		j.Done,
		j.StartTime,
		j.EndTime,
		j.Uuid,
	)
	defer db.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "cannot update job into db")
	}

	j.ID, _ = result.LastInsertId()

	return j.ID, nil
}

func (j *BackupJob) StartPreBackupScript() error {
	if j.Module.PreBackupScript == "" {
		j.GetLog().Info("No pre backup script to launch")
	} else {
		j.GetLog().WithFields(log.Fields{
			"script": j.Module.PreBackupScript,
		}).Info("Starting pre backup script")
	}

	// TODO
	return nil
}

func (j *BackupJob) StartPostBackupScript() error {
	if j.Module.PostBackupScript == "" {
		j.GetLog().Info("No post backup script to launch")
	} else {
		j.GetLog().WithFields(log.Fields{
			"script": j.Module.PostBackupScript,
		}).Info("Starting post backup script")
	}

	// TODO
	return nil
}

func (j *BackupJob) Duration() time.Duration {
	if j.StartTime.IsZero() {
		return time.Time{}.Sub(time.Time{})
	}

	start := j.StartTime
	var end time.Time
	if j.EndTime.IsZero() {
		end = time.Now()
	} else {
		end = j.EndTime
	}

	return end.Sub(start).Truncate(time.Second)
}
