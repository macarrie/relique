package job

import (
	"log/slog"
	"time"

	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/job_status"
	"github.com/macarrie/relique/internal/job_type"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
	"github.com/macarrie/relique/internal/rsync_task"
	rsync_lib "github.com/macarrie/relique/internal/rsync_task/lib"
	"github.com/macarrie/relique/internal/utils"
)

const (
	OK = iota
	Warning
	Critical
	Unknown
)

type Job struct {
	// Database ID
	ID int64 `json:"-" toml:"-"`

	Uuid       string                 `json:"uuid" toml:"uuid"`
	Client     client.Client          `json:"client" toml:"client"`
	Module     module.Module          `json:"module" toml:"module"`
	Status     job_status.JobStatus   `json:"status" toml:"status"`
	Done       bool                   `json:"done" toml:"done"`
	BackupType backup_type.BackupType `json:"backup_type" toml:"backup_type"`
	JobType    job_type.JobType       `json:"job_type" toml:"job_type"`
	StartTime  time.Time              `json:"start_time" toml:"start_time"`
	EndTime    time.Time              `json:"end_time" toml:"end_time"`
	Repository repo.Repository        `json:"repository" toml:"repository"`

	Tasks              []rsync_task.RsyncTask `json:"-" toml:"-"`
	PreviousJobUuid    string                 `json:"previous_job_uuid" toml:"previous_job_uuid"`
	RestoreImageUuid   string                 `json:"restore_image_uuid" toml:"restore_image_uuid"`
	PreviousJob        *Job                   `json:"previous_job" toml:"previous_job"`
	Stats              rsync_lib.Stats        `json:"stats" toml:"stats"`
	CustomRestorePaths map[string]string      `json:"custom_restore_paths" toml:"custom_restore_paths"`

	// For DB storage
	ClientName string `json:"-" toml:"-"`
	ModuleName string `json:"-" toml:"-"`
	RepoName   string `json:"-" toml:"-"`
}

func (j *Job) GetLog() *slog.Logger {
	return slog.With(
		slog.String("uuid", j.Uuid),
		slog.String("client", j.Client.String()),
		slog.String("module", j.Module.String()),
		slog.String("backup_type", j.BackupType.String()),
		slog.String("job_type", j.JobType.String()),
		slog.String("status", j.Status.String()),
		slog.Bool("done", j.Done),
		slog.String("repository", j.Repository.GetName()),
	)
}

func (j *Job) Duration() time.Duration {
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

func (j *Job) GetStorageFolderPath() (string, error) {
	return utils.GetStoragePath(j.Repository, j.RepoName, j.Uuid)
}

func (j *Job) GetCatalogPath() string {
	return utils.GetCatalogPath(j.Uuid)
}
