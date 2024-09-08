package job

import (
	"database/sql"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/backup_type"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
)

func (j *Job) Save() (int64, error) {
	tx, err := db.Handler().Begin()
	// Defers are stacked, defer are executed in reverse order of stacking
	defer func() {
		if err != nil {
			j.GetLog().With(
				slog.Any("error", err),
			).Debug("Rollback job save")
			tx.Rollback()
		}
	}()

	if err != nil {
		return 0, fmt.Errorf("cannot start transaction to save job: %w", err)
	}

	if j.ID != 0 {
		id, err := j.Update(tx)
		if err != nil || id == 0 {
			return 0, fmt.Errorf("cannot update job: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return 0, fmt.Errorf("cannot commit job save transaction: %w", err)
		}

		return id, err
	}

	j.GetLog().Debug("Saving job into database")

	request := sq.Insert("jobs").SetMap(sq.Eq{
		"uuid":        j.Uuid,
		"status":      j.Status.Status,
		"backup_type": j.BackupType.Type,
		"job_type":    j.JobType.Type,
		"done":        j.Done,
		"start_time":  j.StartTime,
		"end_time":    j.EndTime,
		"module_type": j.Module.ModuleType,
		"client_name": j.Client.Name,
		"repo_name":   j.Repository.GetName(),
	})
	query, args, err := request.ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("cannot save job into db: %w", err)
	}

	j.ID, err = result.LastInsertId()
	if j.ID == 0 || err != nil {
		return 0, fmt.Errorf("cannot get last insert ID: %w", err)
	}

	j.GetLog().Debug("Commit job save transaction")
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("cannot commit job save transaction: %w", err)
	}

	return j.ID, nil
}

func (j *Job) Update(tx *sql.Tx) (int64, error) {
	j.GetLog().Debug("Updating job details into database")

	request := sq.Update("jobs").SetMap(sq.Eq{
		"status":      j.Status.Status,
		"backup_type": j.BackupType.Type,
		"job_type":    j.JobType.Type,
		"done":        j.Done,
		"start_time":  j.StartTime,
		"end_time":    j.EndTime,
		"module_type": j.Module.ModuleType,
		"client_name": j.Client.Name,
		"repo_name":   j.Repository.GetName(),
	}).Where(
		"uuid = ?",
		j.Uuid,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	var result sql.Result
	if tx == nil {
		result, err = db.Handler().Exec(query, args...)
	} else {
		result, err = tx.Exec(query, args...)
	}
	if err != nil {
		return 0, fmt.Errorf("cannot update job into db: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		return 0, fmt.Errorf("no rows affected: %w", err)
	}

	return j.ID, nil
}

func GetByUuid(uuid string) (Job, error) {
	slog.With(
		slog.String("uuid", uuid),
	).Debug("Looking for job in database")

	request := sq.Select(
		"id",
		"uuid",
		"status",
		"backup_type",
		"job_type",
		"done",
		"start_time",
		"end_time",
		"client_name",
		"module_type",
		"repo_name",
	).From("jobs").Where("uuid = ?", uuid)
	query, args, err := request.ToSql()
	if err != nil {
		return Job{}, fmt.Errorf("cannot build sql query: %w", err)
	}

	row := db.Handler().QueryRow(query, args...)

	var job Job
	if err := row.Scan(&job.ID,
		&job.Uuid,
		&job.Status.Status,
		&job.BackupType.Type,
		&job.JobType.Type,
		&job.Done,
		&job.StartTime,
		&job.EndTime,
		&job.ClientName,
		&job.ModuleType,
		&job.RepoName,
	); err == sql.ErrNoRows {
		return Job{}, fmt.Errorf("no job with UUID '%s' found in db", uuid)
	} else if err != nil {
		return Job{}, fmt.Errorf("cannot retrieve job from db: %w", err)
	}

	jobStorageFolderPath, err := job.GetStorageFolderPath()
	if err != nil {
		return Job{}, fmt.Errorf("cannot get job storage folder path: %w", err)
	}

	modFilePath := fmt.Sprintf("%s/module.toml", jobStorageFolderPath)
	mod, err := module.LoadFromFile(modFilePath)
	if err != nil {
		return Job{}, fmt.Errorf("linked module cannot be loaded from file: %w", err)
	}
	job.Module = mod

	clFilePath := fmt.Sprintf("%s/client.toml", jobStorageFolderPath)
	cl, err := client.LoadFromFile(clFilePath)
	if err != nil {
		return Job{}, fmt.Errorf("linked client cannot be loaded from file: %w", err)
	}
	job.Client = cl

	repoFilePath := fmt.Sprintf("%s/repo.toml", jobStorageFolderPath)
	r, err := repo.LoadFromFile(repoFilePath)
	if err != nil {
		return Job{}, fmt.Errorf("linked repo cannot be loaded from file: %w", err)
	}
	job.Repository = r

	return job, nil
}

func Search(p api_helpers.PaginationParams) ([]Job, error) {
	slog.Debug("Searching for jobs in db")
	var jobs []Job

	// TODO: Prepare request and clean data to avoid SQL injections
	// TODO: handle status and backup type
	request := sq.Select(
		"uuid",
	).From(
		"jobs",
	)
	if p.Limit > 0 {
		request = request.Limit(p.Limit)
	}
	if p.Offset > 0 {
		request = request.Offset(p.Offset)
	}

	// TODO/ Handle pagination
	// TODO/ Handle search parameters
	request = request.OrderBy("jobs.id DESC")

	query, args, err := request.ToSql()
	if err != nil {
		return []Job{}, fmt.Errorf("cannot build sql query: %w", err)
	}

	rows, err := db.Handler().Query(query, args...)
	if err == sql.ErrNoRows {
		return jobs, nil
	} else if err != nil {
		return jobs, fmt.Errorf("cannot search jobs IDs from db: %w", err)
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			slog.With(
				slog.Any("error", err),
			).Error("Cannot parse job uuid from db")
		}
		if u != "" {
			uuids = append(uuids, u)
		}
	}

	if len(uuids) == 0 {
		// No previous job found
		return jobs, nil
	}

	for _, jobUuid := range uuids {
		jobFromDB, err := GetByUuid(jobUuid)
		if err != nil {
			slog.With(
				slog.Any("error", err),
				slog.String("uuid", jobUuid),
			).Error("Cannot get job with uuid from db")
			continue
		}
		if jobFromDB.ID == 0 {
			slog.With(
				slog.Any("error", err),
				slog.String("uuid", jobUuid),
			).Error("No job with this uuid found in db")
			continue
		}

		jobs = append(jobs, jobFromDB)
	}

	return jobs, nil
}

// TODO: Handle search parameters to have a selective count
func Count() (uint64, error) {
	var count uint64

	request := sq.Select(
		"COUNT(*)",
	).From(
		"jobs",
	)

	request = request.OrderBy("jobs.id DESC")

	query, args, err := request.ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	queryErr := db.Handler().QueryRow(query, args...).Scan(&count)
	if queryErr == sql.ErrNoRows {
		return 0, nil
	} else if queryErr != nil {
		return 0, fmt.Errorf("cannot count jobs from db: %w", err)
	}

	return count, nil
}

func GetPrevious(j Job, backupType backup_type.BackupType) (Job, error) {
	j.GetLog().With(
		slog.String("backup_type", j.BackupType.String()),
	).Debug("Looking for previous backup job")

	request := sq.Select(
		"uuid",
	).From(
		"jobs",
	).Where(
		"jobs.backup_type = ?", backupType.Type,
	).Where(
		"jobs.done = ?", true,
	).Where(
		"jobs.client_name = ?", j.Client.Name,
	).Where(
		"jobs.module_type = ?", j.Module.ModuleType,
	).OrderBy(
		"jobs.id DESC",
	)
	query, args, err := request.ToSql()
	if err != nil {
		return Job{}, fmt.Errorf("cannot build sql query: %w", err)
	}

	rows, err := db.Handler().Query(query, args...)
	if err == sql.ErrNoRows {
		return Job{}, nil
	} else if err != nil {
		return Job{}, fmt.Errorf("cannot query previous full jobs IDs from db: %w", err)
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var jobUuid string
		if err := rows.Scan(&jobUuid); err != nil {
			slog.With(
				slog.Any("error", err),
			).Error("Cannot parse job uuid from db")
		}
		if jobUuid != "" {
			uuids = append(uuids, jobUuid)
		}
	}

	if len(uuids) == 0 {
		// No previous full job found
		return Job{}, fmt.Errorf("no job found with specified criteria")
	}

	// TODO: Filter on module. If not, diff has no use and will not work
	// Get first job since previous jobs are listed by id DESC
	jobUuid := uuids[0]
	jobFromDB, err := GetByUuid(jobUuid)
	if err != nil {
		slog.With(
			slog.Any("error", err),
			slog.String("uuid", jobUuid),
		).Error("Cannot get job with Uuid from db")
		return Job{}, fmt.Errorf("cannot get job with this uuid from db: %w", err)
	}
	if jobFromDB.ID == 0 {
		slog.With(
			slog.String("uuid", jobUuid),
		).Error("No job with this Uuid found in db")
		return Job{}, fmt.Errorf("no job with this Uuid found in db")
	}

	return jobFromDB, nil
}
