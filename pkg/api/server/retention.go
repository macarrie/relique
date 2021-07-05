package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/macarrie/relique/internal/types/config/common"
	"github.com/macarrie/relique/pkg/api/utils"
	"github.com/pkg/errors"
)

func CleanRetention(config common.Configuration) error {
	response, err := utils.PerformRequest(config,
		config.PublicAddress,
		config.Port,
		"POST",
		"/api/v1/retention/clean",
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
		return fmt.Errorf("cannot clean retention on server. See server logs for more details")
	}

	return nil
}
