package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/macarrie/relique/internal/types/backup_job"
	"github.com/macarrie/relique/internal/types/config/common"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

func PingServer(config common.Configuration) error {
	response, err := utils.PerformRequest(config,
		config.PublicAddress,
		config.Port,
		"GET",
		"/api/v1/ping",
		nil)
	if err != nil {
		return errors.Wrap(err, "error when performing api request")
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body from api request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot ping server (%d response): see server logs for more details", response.StatusCode)
	}

	return nil
}

func SearchJob(config common.Configuration, params backup_job.JobSearchParams) ([]backup_job.BackupJob, error) {
	var jobs []backup_job.BackupJob

	response, err := utils.PerformRequest(config,
		config.PublicAddress,
		config.Port,
		"POST",
		"/api/v1/backup/jobs",
		params)
	if err != nil {
		return []backup_job.BackupJob{}, errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return jobs, errors.Wrap(err, "cannot read response body from api requets")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return jobs, fmt.Errorf("cannot get jobs from server (%d response): see server logs for more details", response.StatusCode)
	}

	if err := json.Unmarshal(body, &jobs); err != nil {
		return jobs, errors.Wrap(err, "cannot parse jobs from search results")
	}

	return jobs, nil
}
