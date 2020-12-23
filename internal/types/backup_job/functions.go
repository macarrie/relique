package backup_job

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/macarrie/relique/internal/db"
	log "github.com/macarrie/relique/internal/logging"
	clientObject "github.com/macarrie/relique/internal/types/client"
	"github.com/macarrie/relique/internal/types/job_status"
	"github.com/macarrie/relique/internal/types/module"
	"github.com/pkg/errors"
)

func New(client clientObject.Client, module module.Module) BackupJob {
	return BackupJob{
		Uuid:       uuid.New().String(),
		Client:     client,
		Module:     module,
		Status:     job_status.New(job_status.Pending),
		BackupType: module.BackupType,
	}
}

func GetByUuid(uuid string) (BackupJob, error) {
	log.WithFields(log.Fields{
		"uuid": uuid,
	}).Trace("Looking for job in database")

	request := `SELECT 
		id,
		uuid, 
		status, 
		backup_type, 
		module_id,
		client_id,
        done,
        start_time,
        end_time
	FROM jobs 
	WHERE uuid = $1`
	row := db.Read().QueryRow(request, uuid)
	defer db.RUnlock()

	var job BackupJob
	err := row.Scan(&job.ID,
		&job.Uuid,
		&job.Status.Status,
		&job.BackupType.Type,
		&job.ModuleID,
		&job.ClientID,
		&job.Done,
		&job.StartTime,
		&job.EndTime,
	)
	if err == sql.ErrNoRows {
		return BackupJob{}, nil
	} else if err != nil {
		return BackupJob{}, errors.Wrap(err, "cannot retrieve job from db")
	}

	if job.ModuleID == 0 {
		return BackupJob{}, fmt.Errorf("db inconsistency: no module associated for this job in db")
	}
	if job.ClientID == 0 {
		return BackupJob{}, fmt.Errorf("db inconsistency: no client associated for this job in db")
	}

	mod, err := module.GetByID(job.ModuleID)
	if err == nil && mod.ID == 0 {
		return BackupJob{}, errors.Wrap(err, "job linked module not found in db")
	}
	if err != nil || mod.ID == 0 {
		return BackupJob{}, errors.Wrap(err, "cannot load job linked module")
	}
	job.Module = mod

	cl, err := clientObject.GetByID(job.ClientID)
	if err == nil && cl.ID == 0 {
		return BackupJob{}, errors.Wrap(err, "job linked client not found in db")
	}
	if err != nil || cl.ID == 0 {
		return BackupJob{}, errors.Wrap(err, "cannot load job linked client")
	}
	job.Client = cl

	return job, nil
}

func GetPreviousJob(job BackupJob) (BackupJob, error) {
	log.Trace("Looking for previous full backup job")

	request := `SELECT uuid
		FROM jobs 
		JOIN clients 
			ON jobs.client_id = clients.id 
		JOIN modules 
		    ON jobs.module_id = modules.id 
		WHERE modules.module_type = $1 
		    AND jobs.done = $3 
			AND clients.name = $4 
		ORDER BY jobs.id DESC`
	rows, err := db.Read().Query(request, job.Module.ModuleType, true, job.Client.Name)
	defer db.RUnlock()
	if err == sql.ErrNoRows {
		return BackupJob{}, nil
	} else if err != nil {
		return BackupJob{}, errors.Wrap(err, "cannot query previous full jobs IDs from db")
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot parse job Uuid from db")
		}
		if uuid != "" {
			uuids = append(uuids, uuid)
		}
	}

	if len(uuids) == 0 {
		// No previous full job found
		return BackupJob{}, nil
	}

	jobUuid := uuids[len(uuids)-1]
	jobFromDB, err := GetByUuid(jobUuid)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"uuid": jobUuid,
		}).Error("Cannot get job with Uuid from db")
		return BackupJob{}, errors.Wrap(err, "cannot get job with this Uuid from db")
	}
	if jobFromDB.ID == 0 {
		log.WithFields(log.Fields{
			"uuid": jobUuid,
		}).Error("No job with this Uuid found in db")
		return BackupJob{}, fmt.Errorf("no job with this Uuid found in db")
	}

	return jobFromDB, nil
}

func GetActiveJobs() ([]BackupJob, error) {
	log.Trace("Looking for active jobs in database")

	request := `SELECT uuid
		FROM jobs 
		WHERE status = $1`
	rows, err := db.Read().Query(request, job_status.Active)
	if err == sql.ErrNoRows {
		return []BackupJob{}, nil
	} else if err != nil {
		return []BackupJob{}, errors.Wrap(err, "cannot query active jobs IDs from db")
	}
	defer db.RUnlock()

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

	jobs := make([]BackupJob, 0)
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

func Search(params JobSearchParams) ([]BackupJob, error) {
	params.GetLog().Trace("Searching for jobs in db")
	var jobs []BackupJob

	// TODO: Handle search parameters
	request := `SELECT uuid
		FROM jobs 
		JOIN clients ON jobs.client_id = clients.id 
		JOIN modules ON jobs.module_id = modules.id 
		ORDER BY jobs.id DESC`
	rows, err := db.Read().Query(request)
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
		// No previous full job found
		return jobs, nil
	}

	for _, jobUuid := range uuids {
		jobFromDB, err := GetByUuid(jobUuid)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"uuid": jobUuid,
			}).Error("Cannot get job with Uuid from db")
			continue
		}
		if jobFromDB.ID == 0 {
			log.WithFields(log.Fields{
				"uuid": jobUuid,
			}).Error("No job with this Uuid found in db")
			continue
		}

		jobs = append(jobs, jobFromDB)
	}

	return jobs, nil
}
