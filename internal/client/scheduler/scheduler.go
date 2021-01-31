package scheduler

import (
	"strings"
	"time"

	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/types/job_status"
	clientApi "github.com/macarrie/relique/pkg/api/client"

	"github.com/macarrie/relique/internal/types/backup_job"

	log "github.com/macarrie/relique/internal/logging"
	clientConfig "github.com/macarrie/relique/internal/types/config/client_daemon_config"
)

var RunTicker *time.Ticker
var previousLoopIterHasActiveSchedules bool
var currentDay time.Weekday

func Run() {
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

	// TODO: Check active schedules
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
				continue
			}

			AddJob(m)
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

func AddJob(m module.Module) backup_job.BackupJob {
	job := backup_job.New(clientConfig.BackupConfig, m)
	clientConfig.Jobs = append(clientConfig.Jobs, job)

	return job
}

func CleanJobs() {
	log.Debug("Cleaning done jobs from internal queue")
	var jobs []backup_job.BackupJob
	for _, job := range clientConfig.Jobs {
		if !job.Done {
			jobs = append(jobs, job)
		}
	}

	clientConfig.Jobs = jobs
}
