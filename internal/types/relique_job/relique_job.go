package relique_job

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-multierror"

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
	PreScript = iota
	PostScript
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
	ClientID int64

	Uuid               string                 `json:"uuid"`
	Client             *clientObject.Client   `json:"client"`
	Module             module.Module          `json:"module"`
	Status             job_status.JobStatus   `json:"status"`
	Done               bool                   `json:"done"`
	BackupType         backup_type.BackupType `json:"backup_type"`
	JobType            job_type.JobType       `json:"job_type"`
	PreviousJobUuid    string                 `json:"previous_job_uuid"`
	RestoreJobUuid     string                 `json:"restore_job_uuid"`
	RestoreDestination string                 `json:"restore_destination"`
	StartTime          time.Time              `json:"start_time"`
	EndTime            time.Time              `json:"end_time"`
	StorageRoot        string                 `json:"storage_root"`
	ModuleType         string
	ClientName         string
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
		"storage_root":     j.StorageRoot,
	})
}

func (j *ReliqueJob) GetJobFolderPath() string {
	return filepath.Clean(fmt.Sprintf("%s/%s/", j.StorageRoot, j.Uuid))
}

func (j *ReliqueJob) CreateJobFolder() error {
	if j.StorageRoot == "" {
		return fmt.Errorf("cannot create storage destination folder because job storage is empty")
	}

	root := j.GetJobFolderPath()

	j.GetLog().WithFields(log.Fields{
		"path": root,
	}).Debug("Creating job storage folder")

	return os.MkdirAll(root, 0755)
}

func (j *ReliqueJob) CreateJobDataFolder() error {
	if j.StorageRoot == "" {
		return fmt.Errorf("cannot create storage destination folder because job storage is empty")
	}

	root := j.GetJobFolderPath()
	folder := os.MkdirAll(fmt.Sprintf("%s/_data", root), 0755)
	j.GetLog().WithFields(log.Fields{
		"path": root,
	}).Debug("Creating job data storage subfolder")

	return folder
}

func createLogFolder(j *ReliqueJob) error {
	path := filepath.Clean(fmt.Sprintf("%s/%s", log.GetLogRoot(), j.Uuid))
	return os.MkdirAll(path, 0755)
}

func (j *ReliqueJob) CreateLogFile(path string) (*os.File, error) {
	if err := createLogFolder(j); err != nil {
		return nil, errors.Wrap(err, "cannot create job log folder")
	}
	return os.Create(path)
}

func (j *ReliqueJob) GetFullLogPath(filename string) string {
	return filepath.Clean(fmt.Sprintf(
		"%s/%s/%s.log",
		log.GetLogRoot(),
		j.Uuid,
		sanitize.Accents(sanitize.BaseName(filename)),
	))
}

func (j *ReliqueJob) GetRsyncLogFilePath(path string) string {
	filename := fmt.Sprintf("rsync_log_%s", sanitize.Accents(sanitize.BaseName(path)))
	return j.GetFullLogPath(filename)
}

func (j *ReliqueJob) CreateRsyncLogFile(path string) (*os.File, error) {
	return j.CreateLogFile(j.GetRsyncLogFilePath(path))
}

func (j *ReliqueJob) GetRsyncErrorLogFilePath(path string) string {
	filename := fmt.Sprintf("rsync_error_log_%s", sanitize.Accents(sanitize.BaseName(path)))
	return j.GetFullLogPath(filename)
}

func (j *ReliqueJob) CreateRsyncErrorLogFile(path string) (*os.File, error) {
	return j.CreateLogFile(j.GetRsyncErrorLogFilePath(path))
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

	request := sq.Insert("jobs").SetMap(sq.Eq{
		"uuid":                db.GetNullString(j.Uuid),
		"status":              j.Status.Status,
		"backup_type":         j.BackupType.Type,
		"job_type":            j.JobType.Type,
		"done":                j.Done,
		"start_time":          j.StartTime,
		"end_time":            j.EndTime,
		"restore_job_uuid":    j.RestoreJobUuid,
		"restore_destination": j.RestoreDestination,
		"storage_root":        j.StorageRoot,
		"module_type":         j.Module.ModuleType,
		"client_name":         j.Client.Name,
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

	request := sq.Update("jobs").SetMap(sq.Eq{
		"status":              j.Status.Status,
		"backup_type":         j.BackupType.Type,
		"job_type":            j.JobType.Type,
		"done":                j.Done,
		"start_time":          j.StartTime,
		"end_time":            j.EndTime,
		"restore_job_uuid":    j.RestoreJobUuid,
		"restore_destination": j.RestoreDestination,
		"storage_root":        j.StorageRoot,
		"module_type":         j.ModuleType,
		"client_name":         j.ClientName,
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
		fmt.Sprintf("RELIQUE_MODULE_NAME=%s", j.Module.Name),
		fmt.Sprintf("RELIQUE_MODULE_TYPE=%s", j.Module.ModuleType),
		fmt.Sprintf("RELIQUE_RESTORE_JOB_UUID=%s", j.RestoreJobUuid),
		fmt.Sprintf("RELIQUE_RESTORE_DESTINATION=%s", j.RestoreDestination),
	)
	cmd.Env = append(cmd.Env, j.Module.ExtraParamsEnvVars("RELIQUE_MODULE_PARAM_")...)

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
	if path == "" || path == "none" {
		j.GetLog().WithFields(log.Fields{
			"script_type": scriptType,
		}).Info("No module script to launch")
		return nil
	}

	// Try different paths for script (module directory and parent module directory)
	scriptList := []string{j.Module.Name, j.Module.ModuleType}
	scriptFound := false
	scriptCompletePath := ""
	for _, option := range scriptList {
		fullpath := path
		if !filepath.IsAbs(path) {
			fullpath = j.Module.GetAbsScriptPath(option, path)
		}

		if _, err := os.Lstat(fullpath); !os.IsNotExist(err) {
			scriptFound = true
			scriptCompletePath = fullpath
			j.GetLog().WithFields(log.Fields{
				"path": scriptCompletePath,
			}).Info("Found job script file to use")
			break
		}

		j.GetLog().WithFields(log.Fields{
			"path": fullpath,
		}).Debug("Job script file not found")
	}

	if !scriptFound {
		j.Status.Status = job_status.Error
		j.Done = true
		return fmt.Errorf(
			"module script file (%s) does not exist. Check that the module '%s' (module type '%s') is correctly installed on relique client",
			path,
			j.Module.Name,
			j.Module.ModuleType,
		)
	}

	j.GetLog().WithFields(log.Fields{
		"path":        scriptCompletePath,
		"script_type": scriptType,
	}).Info("Starting module script")

	exitCode, err := j.runModuleScript(scriptCompletePath, logFile)
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

	logFile, err := j.CreateLogFile(j.GetFullLogPath(logFileName))
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

	logFile, err := j.CreateLogFile(j.GetFullLogPath(logFileName))
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

func (j *ReliqueJob) PreFlightCheck() error {
	var errorList *multierror.Error

	moduleIsInstalled, err := module.IsInstalled(j.Module.ModuleType)
	if err != nil {
		errorList = multierror.Append(errorList, errors.Wrapf(err, "cannot check if module '%s' is installed", j.Module.Name))
	} else {
		if !moduleIsInstalled {
			errorList = multierror.Append(errorList, fmt.Errorf("module '%s' is not installed", j.Module.Name))
		}
	}

	return errorList.ErrorOrNil()
}
