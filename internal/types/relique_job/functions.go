package relique_job

import (
	"database/sql"
	"fmt"
	"github.com/macarrie/relique/internal/types/pagination"

	"github.com/macarrie/relique/internal/types/job_type"

	"github.com/macarrie/relique/internal/types/backup_type"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/macarrie/relique/internal/db"
	log "github.com/macarrie/relique/internal/logging"
	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/module"
	"github.com/pkg/errors"
)

func New(client *clientObject.Client, module module.Module, jobType job_type.JobType) ReliqueJob {
	return ReliqueJob{
		Uuid:       uuid.New().String(),
		Client:     client,
		Module:     module,
		Status:     job_status.New(job_status.Pending),
		BackupType: module.BackupType,
		JobType:    jobType,
	}
}

func GetByUuid(uuid string) (ReliqueJob, error) {
	// TODO: Look in job file also and merge results: get client and module info from file and other info from db
	log.WithFields(log.Fields{
		"uuid": uuid,
	}).Trace("Looking for job in database")

	request := sq.Select(
		"id",
		"uuid",
		"status",
		"backup_type",
		"job_type",
		"done",
		"start_time",
		"end_time",
		"restore_job_uuid",
		"restore_destination",
		"storage_root",
		"module_type",
		"client_name",
	).From("jobs").Where("uuid = ?", uuid)
	query, args, err := request.ToSql()
	if err != nil {
		return ReliqueJob{}, errors.Wrap(err, "cannot build sql query")
	}

	row := db.Read().QueryRow(query, args...)
	defer db.RUnlock()

	var job ReliqueJob
	if err := row.Scan(&job.ID,
		&job.Uuid,
		&job.Status.Status,
		&job.BackupType.Type,
		&job.JobType.Type,
		&job.Done,
		&job.StartTime,
		&job.EndTime,
		&job.RestoreJobUuid,
		&job.RestoreDestination,
		&job.StorageRoot,
		&job.ModuleType,
		&job.ClientName,
	); err == sql.ErrNoRows {
		return ReliqueJob{}, errors.Wrap(err, fmt.Sprintf("no job with UUID '%s' found in db", uuid))
	} else if err != nil {
		return ReliqueJob{}, errors.Wrap(err, "cannot retrieve job from db")
	}

	modFilePath := fmt.Sprintf("%s/module.toml", job.GetJobFolderPath())
	mod, err := module.LoadFromFile(modFilePath)
	if err != nil {
		return ReliqueJob{}, errors.Wrap(err, "job linked module cannot be loaded from file")
	}
	job.Module = mod

	clFilePath := fmt.Sprintf("%s/client.toml", job.GetJobFolderPath())
	cl, err := clientObject.LoadFromFile(clFilePath)
	if err != nil {
		return ReliqueJob{}, errors.Wrap(err, "job linked client cannot be loaded from file")
	}
	job.Client = &cl

	return job, nil
}

func PreviousJob(job ReliqueJob, backupType backup_type.BackupType) (ReliqueJob, error) {
	job.GetLog().WithFields(log.Fields{
		"backup_type": backupType.String(),
	}).Trace("Looking for previous backup job")

	request := sq.Select(
		"uuid",
	).From(
		"jobs",
	).Where(
		"jobs.backup_type = ?", backupType.Type,
	).Where(
		"jobs.done = ?", true,
	).Where(
		"jobs.client_name = ?", job.Client.Name,
	).Where(
		"jobs.module_type = ?", job.Module.ModuleType,
	).OrderBy(
		"jobs.id DESC",
	)
	query, args, err := request.ToSql()
	if err != nil {
		return ReliqueJob{}, errors.Wrap(err, "cannot build sql query")
	}

	rows, err := db.Read().Query(query, args...)
	defer db.RUnlock()
	if err == sql.ErrNoRows {
		return ReliqueJob{}, nil
	} else if err != nil {
		return ReliqueJob{}, errors.Wrap(err, "cannot query previous full jobs IDs from db")
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var jobUuid string
		if err := rows.Scan(&jobUuid); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot parse job Uuid from db")
		}
		if jobUuid != "" {
			uuids = append(uuids, jobUuid)
		}
	}

	if len(uuids) == 0 {
		// No previous full job found
		return ReliqueJob{}, fmt.Errorf("no job found with specified criteria")
	}

	// TODO: Filter on module. If not, diff has no use and will not work
	// Get first job since previous jobs are listed by id DESC
	jobUuid := uuids[0]
	jobFromDB, err := GetByUuid(jobUuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"uuid": jobUuid,
		}).Error("Cannot get job with Uuid from db")
		return ReliqueJob{}, errors.Wrap(err, "cannot get job with this Uuid from db")
	}
	if jobFromDB.ID == 0 {
		log.WithFields(log.Fields{
			"uuid": jobUuid,
		}).Error("No job with this Uuid found in db")
		return ReliqueJob{}, fmt.Errorf("no job with this Uuid found in db")
	}

	return jobFromDB, nil
}

func GetActiveJobs() ([]ReliqueJob, error) {
	log.Trace("Looking for active jobs in database")

	request := sq.Select(
		"uuid",
	).From(
		"jobs",
	).Where(
		"status = ?", job_status.Active,
	).OrderBy(
		"jobs.id DESC",
	)
	query, args, err := request.ToSql()
	if err != nil {
		return []ReliqueJob{}, errors.Wrap(err, "cannot build sql query")
	}

	rows, err := db.Read().Query(query, args...)
	defer db.RUnlock()
	if err == sql.ErrNoRows {
		return []ReliqueJob{}, nil
	} else if err != nil {
		return []ReliqueJob{}, errors.Wrap(err, "cannot query active jobs IDs from db")
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot parse active job Uuid from db")
		}
		if uuid != "" {
			uuids = append(uuids, uuid)
		}
	}

	jobs := make([]ReliqueJob, 0)
	for _, uuid := range uuids {
		job, err := GetByUuid(uuid)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"uuid": uuid,
			}).Error("Cannot get job with Uuid from db")
			continue
		}
		if job.ID == 0 {
			log.WithFields(log.Fields{
				"uuid": uuid,
			}).Error("No job with this Uuid found in db")
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func getRequestFromJobSearchParams(request sq.SelectBuilder, params JobSearchParams) sq.SelectBuilder {
	if params.BackupType != "" {
		bType := backup_type.FromString(params.BackupType)
		request = request.Where("jobs.backup_type = ?", bType.Type)
	}
	if params.JobType != "" {
		bType := job_type.FromString(params.JobType)
		request = request.Where("jobs.job_type = ?", bType.Type)
	}
	if params.Status != "" {
		status := job_status.FromString(params.Status)
		request = request.Where("jobs.status = ?", status.Status)
	}
	if params.Module != "" {
		request = request.Where("jobs.module_type = ?", params.Module)
	}
	if params.Uuid != "" {
		request = request.Where("jobs.uuid = ?", params.Uuid)
	}
	if params.Client != "" {
		request = request.Where("jobs.client_name = ?", params.Client)
	}

	return request
}

func Count(params JobSearchParams) (uint64, error) {
	params.GetLog().Trace("Counting jobs in db")
	var count uint64

	// TODO: Prepare request and clean data to avoid SQL injections
	// TODO: handle status and backup type
	request := sq.Select(
		"COUNT(*)",
	).From(
		"jobs",
	)

	request = getRequestFromJobSearchParams(request, params)
	request = request.OrderBy("jobs.id DESC")

	query, args, err := request.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "cannot build sql query")
	}

	queryErr := db.Read().QueryRow(query, args...).Scan(&count)
	defer db.RUnlock()
	if queryErr == sql.ErrNoRows {
		return 0, nil
	} else if queryErr != nil {
		return 0, errors.Wrap(err, "cannot count jobs from db")
	}

	return count, nil
}

func Search(params JobSearchParams, p pagination.Pagination) ([]ReliqueJob, error) {
	params.GetLog().Trace("Searching for jobs in db")
	var jobs []ReliqueJob

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

	request = getRequestFromJobSearchParams(request, params)
	request = request.OrderBy("jobs.id DESC")

	query, args, err := request.ToSql()
	if err != nil {
		return []ReliqueJob{}, errors.Wrap(err, "cannot build sql query")
	}

	rows, err := db.Read().Query(query, args...)
	defer db.RUnlock()
	if err == sql.ErrNoRows {
		return jobs, nil
	} else if err != nil {
		return jobs, errors.Wrap(err, "cannot search jobs IDs from db")
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot parse job Uuid from db")
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
			log.WithFields(log.Fields{
				"err":  err,
				"uuid": jobUuid,
			}).Error("Cannot get job with uuid from db")
			continue
		}
		if jobFromDB.ID == 0 {
			log.WithFields(log.Fields{
				"uuid": jobUuid,
			}).Error("No job with this uuid found in db")
			continue
		}

		jobs = append(jobs, jobFromDB)
	}

	return jobs, nil
}
