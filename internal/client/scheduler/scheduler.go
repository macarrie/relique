package scheduler

import (
	"time"

	"github.com/macarrie/relique/internal/types/job_status"
	clientApi "github.com/macarrie/relique/pkg/api/client"

	"github.com/macarrie/relique/internal/types/backup_job"

	log "github.com/macarrie/relique/internal/logging"
	clientConfig "github.com/macarrie/relique/internal/types/config/client_daemon_config"
)

var RunTicker *time.Ticker

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
			if err := clientApi.RunJob(job); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Error during job execution")
				// TODO: Mark job as error and add message description field with error reason
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

	for _, module := range clientConfig.BackupConfig.Modules {
		// TODO: Check active schedules

		// Create new job only if a job for this module does not already exist
		if clientConfig.JobExists(module) {
			continue
		}

		job := backup_job.New(clientConfig.BackupConfig, module)
		clientConfig.Jobs = append(clientConfig.Jobs, job)
	}

	return nil
}

func CleanJobs() {
	// TODO: Clean done jobs at the end of schedule
}
