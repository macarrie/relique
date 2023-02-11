package scheduler

import (
	"strings"
	"time"

	consts "github.com/macarrie/relique/internal/types"

	"github.com/macarrie/relique/internal/types/job_status"

	"github.com/macarrie/relique/internal/types/module"

	"github.com/macarrie/relique/internal/types/client"

	"github.com/macarrie/relique/internal/types/job_type"

	"github.com/macarrie/relique/internal/types/relique_job"

	serverConfig "github.com/macarrie/relique/internal/types/config/server_daemon_config"

	log "github.com/macarrie/relique/internal/logging"
	serverApi "github.com/macarrie/relique/pkg/api/server"
)

var RunTicker *time.Ticker
var currentJobs []relique_job.ReliqueJob
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
	if len(serverConfig.Config.Clients) == 0 {
		log.Info("No clients found in configuration")
		return
	}

	for i := range serverConfig.Config.Clients {
		cl := &serverConfig.Config.Clients[i]
		if err := serverApi.SendConfiguration(cl); err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"client": serverConfig.Config.Clients[i].Name,
			}).Error("Cannot send configuration to client")
		}

		if cl.APIAlive != consts.OK {
			cl.GetLog().Warning("Client is not alive. Jobs for this client will not be started until relique client can be pinged")
			continue
		}

		if err := CreateJobs(cl); err != nil {
			cl.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Cannot create backup jobs")
		}
	}

	if err := StartJobs(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Some backup jobs could not be started correctly")
	}

	activeJobs, err := relique_job.GetActiveJobs()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot get currently active jobs")
	}
	if len(activeJobs) == 0 {
		log.Info("No active backup jobs")
	} else {
		log.WithFields(log.Fields{
			"nb": len(activeJobs),
		}).Info("Active jobs")
		for _, job := range activeJobs {
			job.GetLog().Info("Active job currently being handled")
		}
	}
}

func CreateJobs(cl *client.Client) error {
	if len(cl.Modules) == 0 {
		cl.GetLog().Info("No backup modules defined. No backup jobs to start")
		return nil
	}

	// Clean jobs on day change
	if currentDay != time.Now().Weekday() {
		CleanJobs()
		currentDay = time.Now().Weekday()
	}

	for _, m := range cl.Modules {
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
			if jobExists(m) {
				previousLoopIterHasActiveSchedules = hasActiveSchedule
				continue
			}

			job := relique_job.New(cl, m, job_type.JobType{Type: job_type.Backup})
			job.StorageRoot = serverConfig.Config.BackupStoragePath
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

func StartJobs() error {
	for i := range currentJobs {
		job := &currentJobs[i]
		if job.Status.Status == job_status.Pending {
			// TODO: Run jobs in goroutines to avoid locking main loop -> Create job pool
			if err := serverApi.RunJob(job); err != nil {
				job.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Error during job execution")
				// TODO: Add message description field with error reason
				job.Status.Status = job_status.Error
			}
		}
		job.GetLog().WithFields(log.Fields{
			"status": job.Status.String(),
		}).Info("Backup job currently handled")
	}

	return nil
}

func AddJob(job relique_job.ReliqueJob) {
	currentJobs = append(currentJobs, job)

	if err := UpdateRetention(serverConfig.Config.RetentionPath); err != nil {
		log.WithFields(log.Fields{
			"path": serverConfig.Config.RetentionPath,
			"err":  err,
		}).Error("Cannot update jobs retention. Done jobs will not be remembered and might be restarted at relique server restart")
	}
}

func CleanJobs() {
	log.Debug("Cleaning done jobs from internal queue")
	var jobs []relique_job.ReliqueJob
	for _, job := range currentJobs {
		if !job.Done {
			jobs = append(jobs, job)
		}
	}

	currentJobs = jobs
}

func jobExists(module module.Module) bool {
	for _, backupJob := range currentJobs {
		if backupJob.Module.Name == module.Name && backupJob.Module.ModuleType == module.ModuleType {
			return true
		}
	}

	return false
}
