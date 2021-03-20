package relique_job

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kennygrant/sanitize"

	"github.com/macarrie/relique/internal/types/job_type"

	sq "github.com/Masterminds/squirrel"

	"github.com/macarrie/relique/internal/db"
	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/backup_type"
	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/module"
	"github.com/pkg/errors"
)

type ReliqueJob struct {
	// Database IDs
	ID       int64
	ModuleID int64
	ClientID int64

	Uuid               string
	Client             clientObject.Client
	Module             module.Module
	Status             job_status.JobStatus
	Done               bool
	BackupType         backup_type.BackupType
	JobType            job_type.JobType
	PreviousJobUuid    string
	RestoreJobUuid     string
	RestoreDestination string
	StartTime          time.Time
	EndTime            time.Time
}

func (j *ReliqueJob) GetLog() *log.Entry {
	return log.WithFields(log.Fields{
		"uuid":             j.Uuid,
		"client":           j.Client.String(),
		"module":           j.Module.String(),
		"backup_type":      j.BackupType.String(),
		"job_type":         j.JobType.String(),
		"restore_job_uuid": j.RestoreJobUuid,
		"status":           j.Status.String(),
		"done":             j.Done,
	})
}

func createLogFolder(j *ReliqueJob) error {
	path := filepath.Clean(fmt.Sprintf("%s/%s", log.GetLogRoot(), j.Uuid))
	return os.MkdirAll(path, 0755)
}

func (j *ReliqueJob) GetRsyncLogFile(path string) (*os.File, error) {
	if err := createLogFolder(j); err != nil {
		return nil, errors.Wrap(err, "cannot create job log folder")
	}

	logFilePath := filepath.Clean(fmt.Sprintf("%s/%s/rsync_log_%s.log", log.GetLogRoot(), j.Uuid, sanitize.Accents(sanitize.BaseName(path))))
	return os.Create(logFilePath)
}

func (j *ReliqueJob) GetRsyncErrorLogFile(path string) (*os.File, error) {
	if err := createLogFolder(j); err != nil {
		return nil, errors.Wrap(err, "cannot create job log folder")
	}

	logFilePath := filepath.Clean(fmt.Sprintf("%s/%s/rsync_error_log_%s.log", log.GetLogRoot(), j.Uuid, sanitize.Accents(sanitize.BaseName(path))))
	return os.Create(logFilePath)
}

func (j *ReliqueJob) Save() (int64, error) {
	tx, err := db.Write().Begin()
	// Defers are stacked, defer are executed in reverse order of stacking
	defer db.Unlock()
	defer func() {
		if err != nil {
			j.GetLog().Debug("Rollback job save")
			tx.Rollback()
		}
	}()

	if err != nil {
		return 0, errors.Wrap(err, "cannot start transaction to save job")
	}

	if j.ID != 0 {
		id, err := j.Update(tx)
		if err != nil || id == 0 {
			return 0, errors.Wrap(err, "cannot update job")
		}

		if err := tx.Commit(); err != nil {
			return 0, errors.Wrap(err, "cannot commit job save transaction")
		}

		return id, err
	}

	j.GetLog().Debug("Saving job into database")

	moduleId, err := j.Module.Save(tx)
	if err != nil || moduleId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner module")
	}

	clientId, err := j.Client.Save(tx)
	if err != nil || clientId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner client")
	}

	request := sq.Insert("jobs").SetMap(sq.Eq{
		"uuid":                db.GetNullString(j.Uuid),
		"status":              j.Status.Status,
		"backup_type":         j.BackupType.Type,
		"job_type":            j.JobType.Type,
		"done":                j.Done,
		"module_id":           db.GetNullInt32(uint32(moduleId)),
		"client_id":           db.GetNullInt32(uint32(clientId)),
		"start_time":          j.StartTime,
		"end_time":            j.EndTime,
		"restore_job_uuid":    j.RestoreJobUuid,
		"restore_destination": j.RestoreDestination,
	})
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "cannot save job into db")
	}

	j.ID, err = result.LastInsertId()
	if j.ID == 0 || err != nil {
		return 0, errors.Wrap(err, "cannot get last insert ID")
	}

	j.GetLog().Debug("Commit job save transaction")
	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "cannot commit job save transaction")
	}

	return j.ID, nil
}

func (j *ReliqueJob) Update(tx *sql.Tx) (int64, error) {
	j.GetLog().Debug("Updating job details into database")

	var moduleId int64
	var clientId int64

	moduleId, err := j.Module.Save(tx)
	if err != nil || moduleId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner module")
	}
	clientId, err = j.Client.Save(tx)
	if err != nil || clientId == 0 {
		return 0, errors.Wrap(err, "cannot save job inner client")
	}

	request := sq.Update("jobs").SetMap(sq.Eq{
		"status":              j.Status.Status,
		"backup_type":         j.BackupType.Type,
		"job_type":            j.JobType.Type,
		"module_id":           db.GetNullInt32(uint32(moduleId)),
		"client_id":           db.GetNullInt32(uint32(clientId)),
		"done":                j.Done,
		"start_time":          j.StartTime,
		"end_time":            j.EndTime,
		"restore_job_uuid":    j.RestoreJobUuid,
		"restore_destination": j.RestoreDestination,
	}).Where(
		"uuid = ?",
		j.Uuid,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	var result sql.Result
	if tx == nil {
		result, err = db.Write().Exec(query, args...)
		defer db.Unlock()
	} else {
		result, err = tx.Exec(query, args...)
	}
	if err != nil {
		return 0, errors.Wrap(err, "cannot update job into db")
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		return 0, errors.Wrap(err, "no rows affected")
	}

	return j.ID, nil
}

func (j *ReliqueJob) StartPreBackupScript() error {
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

func (j *ReliqueJob) StartPostBackupScript() error {
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

func (j *ReliqueJob) StartPreRestoreScript() error {
	if j.Module.PreRestoreScript == "" {
		j.GetLog().Info("No pre restore script to launch")
	} else {
		j.GetLog().WithFields(log.Fields{
			"script": j.Module.PreRestoreScript,
		}).Info("Starting pre restore script")
	}

	// TODO
	return nil
}

func (j *ReliqueJob) StartPostRestoreScript() error {
	if j.Module.PostBackupScript == "" {
		j.GetLog().Info("No post restore script to launch")
	} else {
		j.GetLog().WithFields(log.Fields{
			"script": j.Module.PostRestoreScript,
		}).Info("Starting post restore script")
	}

	// TODO
	return nil
}

func (j *ReliqueJob) Duration() time.Duration {
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
