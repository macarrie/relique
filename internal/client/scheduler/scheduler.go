package scheduler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/macarrie/relique/internal/types/job_type"

	"github.com/macarrie/relique/internal/types/job_status"
	clientApi "github.com/macarrie/relique/pkg/api/client"

	"github.com/macarrie/relique/internal/types/relique_job"

	log "github.com/macarrie/relique/internal/logging"
	clientConfig "github.com/macarrie/relique/internal/types/config/client_daemon_config"
)

var RunTicker *time.Ticker
var previousLoopIterHasActiveSchedules bool
var currentDay time.Weekday

func Run() {
	// Set weekday to current day to avoid cleaning done jobs on startup
	currentDay = time.Now().Weekday()

	RunTicker = time.NewTicker(10 * time.Second)
	go func() {
		log.Debug("Starting main daemon loop")
		for {
			poll()
			<-RunTicker.C
		}
	}()
}

func poll() {
	if clientConfig.BackupConfig.Version == "" {
		log.Info("Waiting for configuration from relique server")
		return
	}

	if err := CreateBackupJobs(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot create backup jobs")
	}

	for i, _ := range clientConfig.Jobs {
		job := &clientConfig.Jobs[i]
		if job.Status.Status == job_status.Pending {
			// TODO: Run jobs in goroutines to avoid locking main loop -> Create job pool
			if err := clientApi.RunJob(job); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Error during job execution")
				// TODO: Add message description field with error reason
				job.Status.Status = job_status.Error
			}
		}
		job.GetLog().WithFields(log.Fields{
			"status": job.Status.String(),
		}).Info("Backup job currently handled by relique client")
	}
}

func CreateBackupJobs() error {
	if len(clientConfig.BackupConfig.Modules) == 0 {
		log.Info("No backup modules defined. No backup jobs to start")
		return nil
	}

	// Clean jobs on day change
	if currentDay != time.Now().Weekday() {
		CleanJobs()
		currentDay = time.Now().Weekday()
	}

	for _, m := range clientConfig.BackupConfig.Modules {
		var activeSchedulesNames []string
		hasActiveSchedule := false
		for _, schedule := range m.Schedules {
			if schedule.Active(time.Now()) {
				activeSchedulesNames = append(activeSchedulesNames, schedule.Name)
				hasActiveSchedule = true
			}
		}

		if hasActiveSchedule {
			if previousLoopIterHasActiveSchedules {
				log.WithFields(log.Fields{
					"schedules": strings.Join(activeSchedulesNames, ","),
					"nb":        len(activeSchedulesNames),
				}).Debug("Active schedules")
			} else {
				log.WithFields(log.Fields{
					"schedules": strings.Join(activeSchedulesNames, ","),
					"nb":        len(activeSchedulesNames),
				}).Info("Entering schedule")
			}

			// Create new job only if a job for this module does not already exist
			if clientConfig.JobExists(m) {
				previousLoopIterHasActiveSchedules = hasActiveSchedule
				continue
			}

			job := relique_job.New(clientConfig.BackupConfig, m, job_type.JobType{Type: job_type.Backup})
			AddJob(job)
		} else {
			if previousLoopIterHasActiveSchedules {
				log.Debug("Exiting active schedules")
				// Clean done jobs on schedule exit
				CleanJobs()
			} else {
				log.Debug("No active schedules")
			}
		}

		previousLoopIterHasActiveSchedules = hasActiveSchedule
	}

	return nil
}

func AddJob(job relique_job.ReliqueJob) {
	clientConfig.Jobs = append(clientConfig.Jobs, job)

	if err := UpdateRetention(clientConfig.Config.RetentionPath); err != nil {
		log.WithFields(log.Fields{
			"path": clientConfig.Config.RetentionPath,
			"err":  err,
		}).Error("Cannot update jobs retention. Done jobs will not be remembered and might be restarted at relique client restart")
	}
}

func CleanJobs() {
	log.Debug("Cleaning done jobs from internal queue")
	var jobs []relique_job.ReliqueJob
	for _, job := range clientConfig.Jobs {
		if !job.Done {
			jobs = append(jobs, job)
		}
	}

	clientConfig.Jobs = jobs
}

func LoadRetention(path string) error {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Loading jobs retention file")

	if _, err := os.Lstat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"path": path,
		}).Info("Jobs retention file does not exist. Nothing to load")
		return nil
	}

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return errors.Wrap(err, "cannot open retention file")
	}

	byteVal, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "cannot read retention file contents")
	}

	var jobsFromRetention []relique_job.ReliqueJob
	if err := json.Unmarshal(byteVal, &jobsFromRetention); err != nil {
		return errors.Wrap(err, "cannot parse retention file")
	}

	clientConfig.Jobs = jobsFromRetention
	return nil
}

func UpdateRetention(path string) error {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Updating jobs retention file")

	jsonData, err := json.MarshalIndent(clientConfig.Jobs, "", " ")
	if err != nil {
		return errors.Wrap(err, "cannot form json from retention data")
	}

	if err := ioutil.WriteFile(path, jsonData, 0644); err != nil {
		return errors.Wrap(err, "cannot write jobs to retention file")
	}

	return nil
}
