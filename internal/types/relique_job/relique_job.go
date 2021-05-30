package relique_job

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
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

const (
	OK = iota
	Warning
	Critical
	Unknown
)

const (
	PreBackup = iota
	PostBackup
	PreRestore
	PostRestore
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

func (j *ReliqueJob) GetLogFile(name string) (*os.File, error) {
	if err := createLogFolder(j); err != nil {
		return nil, errors.Wrap(err, "cannot create job log folder")
	}

	logFilePath := filepath.Clean(fmt.Sprintf(
		"%s/%s/%s.log",
		log.GetLogRoot(),
		j.Uuid,
		sanitize.Accents(sanitize.BaseName(name)),
	))
	return os.Create(logFilePath)
}

func (j *ReliqueJob) GetRsyncLogFile(path string) (*os.File, error) {
	name := fmt.Sprintf("rsync_log_%s", sanitize.Accents(sanitize.BaseName(path)))
	return j.GetLogFile(name)
}

func (j *ReliqueJob) GetRsyncErrorLogFile(path string) (*os.File, error) {
	name := fmt.Sprintf("rsync_error_log_%s", sanitize.Accents(sanitize.BaseName(path)))
	return j.GetLogFile(name)
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

func (j *ReliqueJob) runModuleScript(path string, logFile *os.File) (int, error) {
	cmd := exec.Command(path)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("RELIQUE_JOB_UUID=%s", j.Uuid),
		fmt.Sprintf("RELIQUE_JOB_TYPE=%s", j.JobType.String()),
		fmt.Sprintf("RELIQUE_JOB_BACKUP_TYPE=%s", j.BackupType.String()),
		fmt.Sprintf("RELIQUE_JOB_BACKUP_TYPE=%s", j.BackupType.String()),
		fmt.Sprintf("RELIQUE_MODULE_NAME=%s", j.Module.Name),
		fmt.Sprintf("RELIQUE_MODULE_TYPE=%s", j.Module.ModuleType),
		fmt.Sprintf("RELIQUE_RESTORE_JOB_UUID=%s", j.RestoreJobUuid),
		fmt.Sprintf("RELIQUE_RESTORE_DESTINATION=%s", j.RestoreDestination),
	)

	if err := cmd.Start(); err != nil {
		return Critical, errors.Wrap(err, "cannot start command")
	}

	var returnErr error
	var returnStatus int
	if err := cmd.Wait(); err == nil {
		returnErr = nil
		returnStatus = OK
	} else {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnErr = errors.Wrap(err, "error during module script execution")
			returnStatus = exitError.ExitCode()
		} else {
			returnErr = errors.Wrap(err, "could not determine exit code returned from module script execution")
			returnStatus = Critical
		}
	}

	if logFile != nil {
		if _, err := logFile.Write(output.Bytes()); err != nil {
			j.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Cannot write script log to log file")
		}
	}

	return returnStatus, returnErr
}

func (j *ReliqueJob) startJobScript(path string, logFile *os.File, scriptType int) error {
	if path == "" {
		j.GetLog().WithFields(log.Fields{
			"script_type": scriptType,
		}).Info("No module script to launch")
		return nil
	}

	if !filepath.IsAbs(path) {
		path = filepath.Clean(fmt.Sprintf("%s/%s/scripts/%s", module.MODULES_INSTALL_PATH, j.Module.Name, path))
	}

	if _, err := os.Lstat(path); os.IsNotExist(err) {
		j.Status.Status = job_status.Error
		j.Done = true
		return fmt.Errorf(
			"module script file (%s) does not exist. Check that the module '%s' (module type '%s') is correctly installed on relique client",
			path,
			j.Module.Name,
			j.Module.ModuleType,
		)
	} else {
		j.GetLog().WithFields(log.Fields{
			"path":        path,
			"script_type": scriptType,
		}).Info("Starting module script")

		exitCode, err := j.runModuleScript(path, logFile)
		if err != nil {
			switch exitCode {
			case Critical:
				j.Status.Status = job_status.Error
				j.Done = true
				j.GetLog().WithFields(log.Fields{
					"err":         err,
					"script_type": scriptType,
				}).Error("Critical exit code returned from module script. Check client job logs for more details")
				return fmt.Errorf("critical return code from module script. Check client job logs for more details")
			case Warning, Unknown:
				j.GetLog().WithFields(log.Fields{
					"err": err,
				}).Warning("Warning or Unknown exit code returned from module script. Check client job logs for more details")
				j.Status.Status = job_status.Incomplete
				return nil
			default:
				j.GetLog().WithFields(log.Fields{
					"err":       err,
					"exit_code": exitCode,
				}).Error("Unknown exit code returned from module script")
				j.Status.Status = job_status.Error
				j.Done = true
				return fmt.Errorf("unknown return code from module script. Check client job logs for more details")
			}
		}
	}

	return nil
}

func (j *ReliqueJob) StartPreScript() error {
	var path string
	var logFileName string
	var scriptType int
	if j.JobType.Type == job_type.Backup {
		logFileName = "prebackup"
		path = j.Module.PreBackupScript
		scriptType = PreBackup
	} else {
		logFileName = "prerestore"
		path = j.Module.PreRestoreScript
		scriptType = PreRestore
	}

	logFile, err := j.GetLogFile(logFileName)
	if err != nil {
		j.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot create log file for module pre script")
		logFile = nil
	}
	defer logFile.Close()

	return j.startJobScript(path, logFile, scriptType)
}

func (j *ReliqueJob) StartPostScript() error {
	var path string
	var logFileName string
	var scriptType int
	if j.JobType.Type == job_type.Backup {
		logFileName = "postbackup"
		path = j.Module.PostBackupScript
		scriptType = PostBackup
	} else {
		logFileName = "postrestore"
		path = j.Module.PostRestoreScript
		scriptType = PostRestore
	}

	logFile, err := j.GetLogFile(logFileName)
	if err != nil {
		j.GetLog().WithFields(log.Fields{
			"err": err,
		}).Error("Cannot create log file for module post script")
		logFile = nil
	}
	defer logFile.Close()

	return j.startJobScript(path, logFile, scriptType)
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
