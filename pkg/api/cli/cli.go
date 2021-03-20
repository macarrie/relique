package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/macarrie/relique/internal/types/config/common"
	"github.com/macarrie/relique/internal/types/relique_job"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

func PingDaemon(config common.Configuration) error {
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

func ManualJobStart(config common.Configuration, params relique_job.JobSearchParams) (relique_job.ReliqueJob, error) {
	var job relique_job.ReliqueJob

	response, err := utils.PerformRequest(config,
		config.PublicAddress,
		config.Port,
		"POST",
		"/api/v1/job/start",
		params)
	if err != nil {
		return relique_job.ReliqueJob{}, errors.Wrap(err, "error when performing api request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return relique_job.ReliqueJob{}, errors.Wrap(err, "cannot read response body from api request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return relique_job.ReliqueJob{}, fmt.Errorf("cannot start job on client (%d response): '%s'", response.StatusCode, body)
	}

	if err := json.Unmarshal(body, &job); err != nil {
		return relique_job.ReliqueJob{}, errors.Wrap(err, "cannot parse started job returned from client")
	}

	return job, nil
}

func SearchJob(config common.Configuration, params relique_job.JobSearchParams) ([]relique_job.ReliqueJob, error) {
	var jobs []relique_job.ReliqueJob

	response, err := utils.PerformRequest(config,
		config.PublicAddress,
		config.Port,
		"POST",
		"/api/v1/backup/jobs",
		params)
	if err != nil {
		return []relique_job.ReliqueJob{}, errors.Wrap(err, "error when performing api request")
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
