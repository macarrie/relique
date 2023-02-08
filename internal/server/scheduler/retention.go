package scheduler

import (
	"encoding/json"
	"io"
	"os"

	log "github.com/macarrie/relique/internal/logging"
	"github.com/macarrie/relique/internal/types/relique_job"
	"github.com/pkg/errors"
)

func LoadRetention(path string) (err error) {
	// TODO: If active jobs found in retention, it means these jobs are failed and should be marker as such in db
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
	if err != nil {
		return errors.Wrap(err, "cannot open retention file")
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = errors.Wrap(cerr, "cannot close file correctly")
		}
	}()

	byteVal, err := io.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "cannot read retention file contents")
	}

	var jobsFromRetention []relique_job.ReliqueJob
	if err := json.Unmarshal(byteVal, &jobsFromRetention); err != nil {
		return errors.Wrap(err, "cannot parse retention file")
	}

	currentJobs = jobsFromRetention
	return nil
}

func UpdateRetention(path string) error {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Updating jobs retention file")

	jsonData, err := json.MarshalIndent(currentJobs, "", " ")
	if err != nil {
		return errors.Wrap(err, "cannot form json from retention data")
	}

	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return errors.Wrap(err, "cannot write jobs to retention file")
	}

	return nil
}

func CleanRetention(path string) error {
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Cleaning jobs retention")

	// Check if retention has running jobs
	warn := false
	var runningJobs []relique_job.ReliqueJob
	for i := range currentJobs {
		if !currentJobs[i].Done {
			runningJobs = append(runningJobs, currentJobs[i])
			warn = true
		}
	}

	currentJobs = runningJobs
	if warn {
		log.Warning("Found running jobs in retention. These jobs are excluded from retention clean and will stay in retention")
	}

	return UpdateRetention(path)
}
